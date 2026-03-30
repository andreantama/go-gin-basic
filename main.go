// main.go adalah entry point (titik masuk) aplikasi Go.
// File ini bertanggung jawab untuk:
//   1. Memuat konfigurasi dari environment variable.
//   2. Menginisialisasi koneksi database.
//   3. Melakukan dependency injection (menghubungkan semua lapisan).
//   4. Mendaftarkan middleware dan route HTTP.
//   5. Menjalankan server HTTP.
//
// Prinsip "Dependency Injection" sangat penting di sini:
// main.go adalah satu-satunya tempat di mana semua komponen "dirakit" bersama.
// Setiap lapisan (repository, usecase, handler) menerima dependency-nya melalui constructor.
package main

import (
	// fmt digunakan untuk mencetak pesan ke console.
	"fmt"

	// log digunakan untuk logging fatal error yang menghentikan aplikasi.
	"log"

	// config berisi fungsi untuk memuat konfigurasi dari environment.
	"github.com/andreantama/go-gin-basic/config"

	// database berisi fungsi untuk membuat koneksi ke MySQL.
	"github.com/andreantama/go-gin-basic/infrastructure/database"

	// delivery berisi HTTP handler untuk endpoint API.
	deliveryHTTP "github.com/andreantama/go-gin-basic/internal/delivery/http"

	// middleware berisi middleware HTTP (auth, CORS, dll).
	"github.com/andreantama/go-gin-basic/internal/middleware"

	// repository berisi implementasi repository untuk akses database.
	"github.com/andreantama/go-gin-basic/internal/repository"

	// scheduler berisi fitur penjadwalan tugas otomatis (mirip Laravel Scheduler).
	"github.com/andreantama/go-gin-basic/internal/scheduler"

	// usecase berisi implementasi business logic.
	"github.com/andreantama/go-gin-basic/internal/usecase"

	// gin adalah HTTP web framework yang digunakan.
	"github.com/gin-gonic/gin"
)

// main adalah fungsi utama yang dipanggil saat aplikasi dijalankan.
// Fungsi ini mengorkestrasi semua inisialisasi dan menjalankan server HTTP.
func main() {
	// ─── LANGKAH 1: MUAT KONFIGURASI ────────────────────────────────────────
	// Muat semua konfigurasi dari file .env atau environment variable sistem.
	cfg, err := config.LoadConfig()
	if err != nil {
		// log.Fatal mencetak pesan error dan menghentikan aplikasi.
		// Aplikasi tidak bisa berjalan tanpa konfigurasi yang valid.
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// ─── LANGKAH 2: SET GIN MODE ────────────────────────────────────────────
	// Set mode GIN berdasarkan environment untuk mengoptimalkan performa.
	if cfg.AppEnv == "production" {
		// Mode production: GIN menonaktifkan debug logging dan mengoptimalkan performa.
		gin.SetMode(gin.ReleaseMode)
	} else {
		// Mode development: GIN menampilkan route yang terdaftar dan debug info.
		gin.SetMode(gin.DebugMode)
	}

	// ─── LANGKAH 3: INISIALISASI DATABASE ───────────────────────────────────
	// Buat koneksi ke database MySQL menggunakan konfigurasi yang sudah dimuat.
	// Fungsi ini juga menjalankan auto-migration untuk membuat tabel yang diperlukan.
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		// Hentikan aplikasi jika tidak bisa terhubung ke database.
		// Aplikasi tidak bisa berfungsi tanpa database.
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// ─── LANGKAH 4: DEPENDENCY INJECTION ────────────────────────────────────
	// Inisialisasi setiap lapisan dengan dependency yang dibutuhkan.
	// Urutan: Repository → Use Case → Handler
	// (dari lapisan dalam ke lapisan luar)

	// Buat instance repository dengan menginject koneksi database.
	// Repository bertanggung jawab untuk semua operasi database.
	userRepo := repository.NewUserRepository(db)

	// Buat instance use case dengan menginject repository dan konfigurasi.
	// Use case bertanggung jawab untuk semua business logic.
	userUsecase := usecase.NewUserUsecase(userRepo, cfg)

	// Buat instance HTTP handler dengan menginject use case.
	// Handler bertanggung jawab untuk menangani HTTP request/response.
	userHandler := deliveryHTTP.NewUserHandler(userUsecase)

	// ─── LANGKAH 5: SETUP ROUTER GIN ────────────────────────────────────────
	// Buat instance router GIN.
	// gin.New() membuat router tanpa middleware default.
	// Kita tambahkan middleware secara eksplisit di bawah.
	router := gin.New()

	// Tambahkan middleware Logger untuk mencatat setiap request HTTP.
	// Logger middleware menampilkan: method, path, status code, dan waktu proses.
	router.Use(gin.Logger())

	// Tambahkan middleware Recovery untuk menangkap panic dan mengembalikan 500.
	// Tanpa ini, satu panic akan menghentikan seluruh server.
	router.Use(gin.Recovery())

	// Tambahkan middleware CORS untuk mengizinkan request dari browser (cross-origin).
	// Diperlukan jika frontend berada di domain yang berbeda dari API.
	router.Use(middleware.CORSMiddleware())

	// ─── LANGKAH 6: DAFTARKAN ROUTE ─────────────────────────────────────────
	// Daftarkan semua route untuk setiap handler.

	// Daftarkan route pengguna (GET /users, POST /auth/register, dll).
	// Beberapa route dilindungi oleh JWT middleware yang dikonfigurasi di dalam RegisterRoutes.
	userHandler.RegisterRoutes(router)

	// Contoh route yang dilindungi JWT (cara alternatif menggunakan group):
	// protected := router.Group("/api/v1")
	// protected.Use(middleware.AuthMiddleware(cfg))
	// {
	//     protected.GET("/profile", profileHandler.GetProfile)
	// }

	// ─── LANGKAH 7: INISIALISASI SCHEDULER ─────────────────────────────────
	// Buat instance scheduler untuk menjalankan tugas-tugas terjadwal.
	// Scheduler ini mirip dengan Laravel Task Scheduler.
	taskScheduler := scheduler.NewScheduler()

	// Daftarkan tugas-tugas terjadwal.
	// Contoh: Health check setiap menit untuk memastikan server berjalan normal.
	err = taskScheduler.EveryMinute("health-check", func() error {
		log.Println("💓 [Health Check] Server berjalan normal")
		return nil
	})
	if err != nil {
		log.Printf("⚠️ Gagal mendaftarkan task health-check: %v", err)
	}

	// Contoh: Membersihkan data sementara setiap hari pukul 02:00.
	// Anda bisa menambahkan logika pembersihan database, file cache, dll.
	err = taskScheduler.DailyAt("cleanup-temp-data", 2, 0, func() error {
		log.Println("🧹 [Cleanup] Membersihkan data sementara...")
		// Tambahkan logika pembersihan di sini.
		// Contoh: db.Where("created_at < ?", time.Now().AddDate(0, 0, -30)).Delete(&SomeModel{})
		return nil
	})
	if err != nil {
		log.Printf("⚠️ Gagal mendaftarkan task cleanup-temp-data: %v", err)
	}

	// Jalankan scheduler di background (non-blocking).
	taskScheduler.Start()
	// Hentikan scheduler saat aplikasi shutdown.
	defer taskScheduler.Stop()

	// ─── LANGKAH 8: JALANKAN SERVER ─────────────────────────────────────────
	// Cetak informasi server yang berjalan ke console.
	fmt.Printf("🚀 Server %s berjalan di http://localhost:%s\n", cfg.AppName, cfg.AppPort)
	fmt.Printf("📌 Environment: %s\n", cfg.AppEnv)

	// Jalankan server HTTP di port yang dikonfigurasi.
	// router.Run() memblokir dan terus mendengarkan request sampai server dihentikan.
	if err := router.Run(":" + cfg.AppPort); err != nil {
		// Hentikan aplikasi jika server gagal dijalankan.
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
