// Package http berisi HTTP handler yang menangani request dan response HTTP.
//
// Dalam Clean Architecture, handler termasuk dalam lapisan "Interface Adapters".
// Handler bertanggung jawab untuk:
//   - Membaca data dari HTTP request (parameter URL, query string, request body).
//   - Memanggil use case yang sesuai dengan business logic.
//   - Mengkonversi hasil dari use case ke HTTP response (JSON).
//
// Handler TIDAK boleh mengandung business logic — tugasnya hanya mengadaptasi
// antara format HTTP dengan format yang dimengerti use case.
package http

import (
	// net/http menyediakan konstanta HTTP status code (http.StatusOK, dll).
	"net/http"

	// strconv digunakan untuk mengkonversi string (dari URL param) ke uint.
	"strconv"

	// domain berisi interface UserUsecase yang digunakan oleh handler.
	"github.com/andreantama/go-gin-basic/internal/domain"

	// pkgerrors berisi tipe error kustom dan helper untuk penanganan error.
	pkgerrors "github.com/andreantama/go-gin-basic/pkg/errors"

	// response berisi helper untuk format respons JSON yang konsisten.
	"github.com/andreantama/go-gin-basic/pkg/response"

	// gin adalah framework HTTP yang digunakan untuk routing dan handling request.
	"github.com/gin-gonic/gin"
)

// UserHandler adalah struct yang menangani semua request HTTP terkait pengguna.
// Struct ini menyimpan reference ke use case yang berisi business logic.
type UserHandler struct {
	// userUsecase adalah interface use case pengguna.
	// Handler bergantung pada INTERFACE, bukan implementasi konkret,
	// sehingga mudah untuk di-mock saat testing.
	userUsecase domain.UserUsecase
}

// NewUserHandler adalah constructor function untuk membuat UserHandler baru.
// Menerima interface UserUsecase melalui dependency injection.
//
// Parameter:
//   - userUsecase: implementasi use case pengguna.
//
// Contoh penggunaan:
//
//	handler := http.NewUserHandler(userUsecase)
//	handler.RegisterRoutes(router)
func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	// Buat dan kembalikan instance UserHandler dengan use case yang diinjeksikan.
	return &UserHandler{userUsecase: userUsecase}
}

// RegisterRoutes mendaftarkan semua route HTTP untuk endpoint pengguna ke router GIN.
// Route dikelompokkan berdasarkan prefix "/users" dan "/auth".
//
// Parameter:
//   - router: instance *gin.Engine tempat route akan didaftarkan.
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	// Grup route untuk autentikasi — tidak memerlukan JWT token.
	auth := router.Group("/auth")
	{
		// POST /auth/register — Endpoint untuk registrasi pengguna baru.
		auth.POST("/register", h.Register)

		// POST /auth/login — Endpoint untuk login dan mendapatkan JWT token.
		auth.POST("/login", h.Login)
	}

	// Grup route untuk manajemen pengguna — memerlukan JWT token (ditambah middleware di main.go).
	users := router.Group("/users")
	{
		// GET /users — Endpoint untuk mengambil semua pengguna.
		users.GET("", h.GetAllUsers)

		// GET /users/:id — Endpoint untuk mengambil pengguna berdasarkan ID.
		users.GET("/:id", h.GetUserByID)

		// PUT /users/:id — Endpoint untuk memperbarui data pengguna.
		users.PUT("/:id", h.UpdateUser)

		// DELETE /users/:id — Endpoint untuk menghapus pengguna.
		users.DELETE("/:id", h.DeleteUser)
	}
}

// registerRequest adalah struct yang merepresentasikan body request untuk registrasi.
// Tag `json` digunakan untuk memetakan field JSON ke field struct.
// Tag `binding:"required"` memastikan field wajib diisi (divalidasi oleh GIN).
type registerRequest struct {
	// Name adalah nama lengkap pengguna, wajib diisi.
	Name string `json:"name" binding:"required"`

	// Email adalah alamat email, wajib diisi dan harus format email yang valid.
	Email string `json:"email" binding:"required,email"`

	// Password adalah password pengguna, wajib diisi dan minimal 8 karakter.
	Password string `json:"password" binding:"required,min=8"`
}

// loginRequest adalah struct untuk body request login.
type loginRequest struct {
	// Email adalah alamat email yang akan digunakan untuk login.
	Email string `json:"email" binding:"required,email"`

	// Password adalah password pengguna untuk verifikasi.
	Password string `json:"password" binding:"required"`
}

// updateUserRequest adalah struct untuk body request update pengguna.
// Field tidak menggunakan `binding:"required"` karena update bersifat opsional (partial update).
type updateUserRequest struct {
	// Name adalah nama baru pengguna (opsional).
	Name string `json:"name"`

	// Password adalah password baru (opsional). Akan di-hash sebelum disimpan.
	Password string `json:"password" binding:"omitempty,min=8"`
}

// Register menangani request HTTP POST /auth/register untuk mendaftarkan pengguna baru.
// Flow:
//  1. Parse dan validasi request body.
//  2. Panggil use case RegisterUser.
//  3. Kembalikan response JSON dengan data pengguna yang baru dibuat.
func (h *UserHandler) Register(c *gin.Context) {
	// Deklarasi variabel untuk menampung data request body.
	var req registerRequest

	// Bind dan validasi request body JSON ke struct registerRequest.
	// ShouldBindJSON mengembalikan error jika:
	//   - Body bukan JSON yang valid.
	//   - Field dengan `binding:"required"` tidak diisi.
	//   - Validasi lain gagal (format email, minimum length, dll).
	if err := c.ShouldBindJSON(&req); err != nil {
		// Kembalikan response 400 Bad Request dengan pesan error validasi.
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// Konversi request ke struct domain.User untuk dikirim ke use case.
	// Handler tidak bisa langsung melempar struct request ke use case
	// karena use case bekerja dengan domain objects, bukan HTTP structs.
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Password akan di-hash oleh use case.
	}

	// Panggil use case untuk mendaftarkan pengguna baru.
	createdUser, err := h.userUsecase.RegisterUser(user)
	if err != nil {
		// Tentukan HTTP status code berdasarkan tipe error.
		statusCode := pkgerrors.HTTPStatus(err)
		// Kembalikan response error dengan status code yang sesuai.
		c.JSON(statusCode, response.Error(err.Error()))
		return
	}

	// Kembalikan response 201 Created dengan data pengguna yang baru dibuat.
	c.JSON(http.StatusCreated, response.Success("Registrasi berhasil", createdUser))
}

// Login menangani request HTTP POST /auth/login untuk autentikasi pengguna.
// Flow:
//  1. Parse dan validasi request body.
//  2. Panggil use case LoginUser.
//  3. Kembalikan response JSON dengan JWT token.
func (h *UserHandler) Login(c *gin.Context) {
	// Deklarasi variabel untuk menampung data request body.
	var req loginRequest

	// Bind dan validasi request body JSON ke struct loginRequest.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Kembalikan 400 Bad Request jika validasi gagal.
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// Panggil use case untuk memvalidasi kredensial dan mendapatkan JWT token.
	token, err := h.userUsecase.LoginUser(req.Email, req.Password)
	if err != nil {
		// Kembalikan 401 Unauthorized jika email/password salah.
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	// Kembalikan response 200 OK dengan JWT token dalam body response.
	c.JSON(http.StatusOK, response.Success("Login berhasil", gin.H{
		// Sertakan token dalam respons untuk digunakan oleh client.
		"token": token,
	}))
}

// GetAllUsers menangani request HTTP GET /users untuk mengambil semua pengguna.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Panggil use case untuk mengambil semua data pengguna.
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		// Kembalikan 500 Internal Server Error jika terjadi kesalahan server.
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// Kembalikan response 200 OK dengan daftar pengguna.
	c.JSON(http.StatusOK, response.Success("Daftar pengguna berhasil diambil", users))
}

// GetUserByID menangani request HTTP GET /users/:id untuk mengambil satu pengguna.
// Parameter `:id` dalam URL path diambil menggunakan c.Param("id").
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Ambil parameter "id" dari URL path (contoh: GET /users/1 → id = "1").
	idStr := c.Param("id")

	// Konversi string ID ke uint karena domain.UserRepository menerima uint.
	// ParseUint(string, base, bitSize) — base 10, 64 bit.
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// Kembalikan 400 Bad Request jika ID bukan angka yang valid.
		c.JSON(http.StatusBadRequest, response.Error("ID tidak valid"))
		return
	}

	// Panggil use case untuk mengambil pengguna berdasarkan ID.
	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		// Tentukan status code berdasarkan tipe error (404 jika tidak ditemukan).
		statusCode := pkgerrors.HTTPStatus(err)
		c.JSON(statusCode, response.Error(err.Error()))
		return
	}

	// Kembalikan response 200 OK dengan data pengguna.
	c.JSON(http.StatusOK, response.Success("Data pengguna berhasil diambil", user))
}

// UpdateUser menangani request HTTP PUT /users/:id untuk memperbarui data pengguna.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Ambil dan konversi ID dari URL path parameter.
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// Kembalikan 400 Bad Request jika ID tidak valid.
		c.JSON(http.StatusBadRequest, response.Error("ID tidak valid"))
		return
	}

	// Deklarasi variabel untuk menampung data request body update.
	var req updateUserRequest

	// Bind request body ke struct updateUserRequest.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Kembalikan 400 Bad Request jika validasi gagal.
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// Buat struct domain.User dengan ID dari URL dan data dari request body.
	user := &domain.User{
		ID:       uint(id), // ID dari URL path.
		Name:     req.Name,
		Password: req.Password, // Password baru (akan di-hash oleh use case).
	}

	// Panggil use case untuk memperbarui data pengguna.
	updatedUser, err := h.userUsecase.UpdateUser(user)
	if err != nil {
		// Tentukan status code berdasarkan tipe error.
		statusCode := pkgerrors.HTTPStatus(err)
		c.JSON(statusCode, response.Error(err.Error()))
		return
	}

	// Kembalikan response 200 OK dengan data pengguna yang sudah diperbarui.
	c.JSON(http.StatusOK, response.Success("Data pengguna berhasil diperbarui", updatedUser))
}

// DeleteUser menangani request HTTP DELETE /users/:id untuk menghapus pengguna.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Ambil dan konversi ID dari URL path parameter.
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// Kembalikan 400 Bad Request jika ID tidak valid.
		c.JSON(http.StatusBadRequest, response.Error("ID tidak valid"))
		return
	}

	// Panggil use case untuk menghapus pengguna berdasarkan ID.
	if err := h.userUsecase.DeleteUser(uint(id)); err != nil {
		// Tentukan status code berdasarkan tipe error.
		statusCode := pkgerrors.HTTPStatus(err)
		c.JSON(statusCode, response.Error(err.Error()))
		return
	}

	// Kembalikan response 200 OK dengan pesan konfirmasi penghapusan.
	c.JSON(http.StatusOK, response.Success("Pengguna berhasil dihapus", nil))
}
