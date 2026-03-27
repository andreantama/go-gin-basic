// Package middleware berisi fungsi middleware HTTP untuk GIN.
//
// Middleware adalah fungsi yang dijalankan SEBELUM handler utama dieksekusi.
// Middleware digunakan untuk:
//   - Autentikasi dan otorisasi (JWT validation).
//   - Logging request dan response.
//   - Rate limiting (membatasi jumlah request).
//   - CORS (Cross-Origin Resource Sharing).
//   - Request ID tracking.
//
// Setiap middleware memanggil c.Next() untuk melanjutkan ke middleware/handler berikutnya,
// atau c.Abort() untuk menghentikan proses dan mengembalikan response segera.
package middleware

import (
	// fmt digunakan untuk memformat pesan log dan error.
	"fmt"

	// net/http menyediakan konstanta HTTP status code.
	"net/http"

	// strings digunakan untuk memanipulasi string (misalnya memisahkan "Bearer token").
	"strings"

	// config berisi konfigurasi aplikasi termasuk JWT secret.
	"github.com/andreantama/go-gin-basic/config"

	// response berisi helper untuk format respons JSON yang konsisten.
	"github.com/andreantama/go-gin-basic/pkg/response"

	// gin adalah framework HTTP yang digunakan.
	"github.com/gin-gonic/gin"

	// jwt adalah library untuk memvalidasi JSON Web Token.
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah middleware untuk validasi JWT token.
// Middleware ini memeriksa apakah request memiliki token JWT yang valid di header Authorization.
// Jika valid, informasi pengguna dari token di-set ke context GIN untuk digunakan oleh handler.
// Jika tidak valid, request ditolak dengan respons 401 Unauthorized.
//
// Parameter:
//   - cfg: konfigurasi aplikasi yang berisi JWT secret untuk validasi.
//
// Penggunaan:
//
//	protectedGroup := router.Group("/api")
//	protectedGroup.Use(middleware.AuthMiddleware(cfg))
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Kembalikan fungsi yang akan dipanggil GIN untuk setiap request.
	return func(c *gin.Context) {
		// Ambil nilai header "Authorization" dari request.
		// Format yang diharapkan: "Bearer <jwt_token>"
		authHeader := c.GetHeader("Authorization")

		// Periksa apakah header Authorization ada dan tidak kosong.
		if authHeader == "" {
			// Jika tidak ada, kembalikan 401 dan hentikan proses request.
			c.JSON(http.StatusUnauthorized, response.Error("Token autentikasi diperlukan"))
			// c.Abort() menghentikan eksekusi middleware dan handler berikutnya.
			c.Abort()
			return
		}

		// Pisahkan header menjadi dua bagian: "Bearer" dan token-nya.
		// strings.SplitN memisahkan string dengan " " maksimal 2 bagian.
		parts := strings.SplitN(authHeader, " ", 2)

		// Validasi format header: harus "Bearer <token>".
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Format header tidak valid, kembalikan 401.
			c.JSON(http.StatusUnauthorized, response.Error("Format token tidak valid. Gunakan: Bearer <token>"))
			c.Abort()
			return
		}

		// Ekstrak string token dari bagian kedua header.
		tokenString := parts[1]

		// Parse dan validasi token JWT menggunakan secret key dari konfigurasi.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma signing yang digunakan adalah HMAC (HS256).
			// Ini mencegah serangan "algorithm confusion" dimana attacker mengganti
			// algoritma ke "none" untuk bypass validasi.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Kembalikan error jika algoritma tidak sesuai.
				return nil, fmt.Errorf("algoritma signing tidak valid: %v", token.Header["alg"])
			}
			// Kembalikan secret key sebagai []byte untuk validasi tanda tangan.
			return []byte(cfg.JWTSecret), nil
		})

		// Periksa apakah parsing dan validasi token berhasil.
		if err != nil || !token.Valid {
			// Token tidak valid atau kadaluarsa, kembalikan 401.
			c.JSON(http.StatusUnauthorized, response.Error("Token tidak valid atau sudah kadaluarsa"))
			c.Abort()
			return
		}

		// Ekstrak claims (payload) dari token yang sudah divalidasi.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// Gagal mengekstrak claims, kembalikan 401.
			c.JSON(http.StatusUnauthorized, response.Error("Gagal membaca informasi token"))
			c.Abort()
			return
		}

		// Simpan informasi pengguna dari token ke context GIN.
		// Context GIN digunakan untuk berbagi data antar middleware dan handler.

		// Ambil user ID dari claim "sub" (subject).
		// JWT menyimpan angka sebagai float64, jadi kita perlu konversi ke uint.
		if sub, ok := claims["sub"].(float64); ok {
			// Set user ID ke context dengan key "userID".
			c.Set("userID", uint(sub))
		}

		// Ambil email dari claims dan set ke context.
		if email, ok := claims["email"].(string); ok {
			// Set email ke context dengan key "email".
			c.Set("email", email)
		}

		// Ambil role dari claims dan set ke context.
		if role, ok := claims["role"].(string); ok {
			// Set role ke context dengan key "role".
			c.Set("role", role)
		}

		// Lanjutkan ke middleware atau handler berikutnya dalam chain.
		c.Next()
	}
}

// AdminOnly adalah middleware otorisasi yang hanya mengizinkan pengguna dengan role "admin".
// Middleware ini harus digunakan SETELAH AuthMiddleware karena memerlukan data dari context.
//
// Penggunaan:
//
//	adminGroup := router.Group("/admin")
//	adminGroup.Use(middleware.AuthMiddleware(cfg), middleware.AdminOnly())
func AdminOnly() gin.HandlerFunc {
	// Kembalikan fungsi middleware.
	return func(c *gin.Context) {
		// Ambil role dari context yang sudah di-set oleh AuthMiddleware.
		role, exists := c.Get("role")
		if !exists {
			// Jika role tidak ada di context, tolak request.
			c.JSON(http.StatusForbidden, response.Error("Akses ditolak"))
			c.Abort()
			return
		}

		// Periksa apakah role adalah "admin".
		if role != "admin" {
			// Jika bukan admin, kembalikan 403 Forbidden.
			c.JSON(http.StatusForbidden, response.Error("Hanya admin yang diizinkan mengakses endpoint ini"))
			c.Abort()
			return
		}

		// Role adalah admin, lanjutkan ke handler berikutnya.
		c.Next()
	}
}

// CORSMiddleware mengkonfigurasi header CORS untuk mengizinkan cross-origin request.
// CORS diperlukan agar aplikasi frontend (di domain berbeda) bisa mengakses API ini.
//
// Penggunaan:
//
//	router.Use(middleware.CORSMiddleware())
func CORSMiddleware() gin.HandlerFunc {
	// Kembalikan fungsi middleware.
	return func(c *gin.Context) {
		// Set header Access-Control-Allow-Origin untuk mengizinkan semua origin.
		// Di production, ganti "*" dengan domain frontend yang spesifik.
		c.Header("Access-Control-Allow-Origin", "*")

		// Izinkan header-header yang umum digunakan dalam request.
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		// Izinkan HTTP method yang digunakan dalam API.
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Jika ini adalah preflight request (OPTIONS), kembalikan 204 No Content.
		// Browser mengirim preflight request sebelum request sebenarnya untuk memeriksa CORS.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Lanjutkan ke handler berikutnya untuk request non-OPTIONS.
		c.Next()
	}
}
