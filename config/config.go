// Package config bertanggung jawab untuk membaca dan menyimpan semua konfigurasi aplikasi.
// Konfigurasi dibaca dari environment variable (atau file .env).
// Dengan memisahkan konfigurasi ke package tersendiri, seluruh bagian aplikasi
// dapat mengakses konfigurasi dengan cara yang konsisten dan mudah diuji.
package config

import (
	// fmt digunakan untuk memformat string error.
	"fmt"

	// os digunakan untuk membaca environment variable dari sistem operasi.
	"os"

	// strconv digunakan untuk mengkonversi string ke tipe data lain (contoh: string ke int).
	"strconv"

	// godotenv digunakan untuk memuat variabel dari file .env ke environment variable.
	"github.com/joho/godotenv"
)

// Config adalah struct yang menampung seluruh konfigurasi aplikasi.
// Setiap field mewakili satu nilai konfigurasi yang dibaca dari environment.
type Config struct {
	// AppName adalah nama aplikasi, digunakan untuk logging dan identifikasi.
	AppName string

	// AppEnv adalah lingkungan aplikasi: "development", "staging", atau "production".
	AppEnv string

	// AppPort adalah port HTTP yang digunakan server, contoh: "8080".
	AppPort string

	// DBHost adalah alamat host database MySQL/PostgreSQL.
	DBHost string

	// DBPort adalah port database, contoh: "3306" untuk MySQL.
	DBPort string

	// DBUser adalah username untuk koneksi database.
	DBUser string

	// DBPassword adalah password untuk koneksi database.
	DBPassword string

	// DBName adalah nama database yang digunakan.
	DBName string

	// JWTSecret adalah kunci rahasia untuk menandatangani token JWT.
	// Simpan nilai ini dengan aman dan jangan commit ke repository.
	JWTSecret string

	// JWTExpireHours adalah masa berlaku token JWT dalam satuan jam.
	JWTExpireHours int
}

// LoadConfig memuat konfigurasi dari file .env (jika ada) dan environment variable.
// Fungsi ini mengembalikan pointer ke Config dan error jika ada masalah.
// Contoh penggunaan:
//
//	cfg, err := config.LoadConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadConfig() (*Config, error) {
	// Coba muat file .env dari direktori saat ini.
	// Jika file tidak ditemukan, kita tetap lanjut karena env var bisa di-set langsung.
	_ = godotenv.Load()

	// Baca nilai JWTExpireHours dan konversi dari string ke int.
	// Jika nilai tidak valid atau kosong, gunakan default 24 jam.
	jwtExpire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		// Jika konversi gagal, gunakan nilai default 24 jam.
		jwtExpire = 24
	}

	// Buat instance Config dengan nilai dari environment variable.
	cfg := &Config{
		// Baca nama aplikasi, default "go-gin-clean-arch".
		AppName: getEnv("APP_NAME", "go-gin-clean-arch"),

		// Baca environment, default "development".
		AppEnv: getEnv("APP_ENV", "development"),

		// Baca port server, default "8080".
		AppPort: getEnv("APP_PORT", "8080"),

		// Baca konfigurasi database.
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "go_gin_db"),

		// Baca konfigurasi JWT.
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-in-production"),
		JWTExpireHours: jwtExpire,
	}

	// Validasi field wajib yang tidak boleh kosong di lingkungan production.
	if cfg.AppEnv == "production" && cfg.JWTSecret == "default-secret-change-in-production" {
		// Kembalikan error jika JWT secret masih menggunakan nilai default di production.
		return nil, fmt.Errorf("JWT_SECRET harus diset di lingkungan production")
	}

	// Kembalikan konfigurasi yang sudah diisi dan nil untuk error (tidak ada error).
	return cfg, nil
}

// DSN menghasilkan Data Source Name (DSN) untuk koneksi ke database MySQL.
// DSN adalah string koneksi yang digunakan oleh GORM untuk terhubung ke database.
// Format: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func (c *Config) DSN() string {
	// Bangun DSN dari field-field konfigurasi database.
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,     // Username database.
		c.DBPassword, // Password database.
		c.DBHost,     // Host database.
		c.DBPort,     // Port database.
		c.DBName,     // Nama database.
	)
}

// getEnv adalah helper function untuk membaca environment variable.
// Jika variable tidak ditemukan atau kosong, kembalikan nilai default (fallback).
// Fungsi ini bersifat private (huruf kecil) sehingga hanya bisa digunakan di package ini.
func getEnv(key, fallback string) string {
	// Coba baca nilai dari environment variable menggunakan key yang diberikan.
	if value, ok := os.LookupEnv(key); ok && value != "" {
		// Jika ditemukan dan tidak kosong, kembalikan nilainya.
		return value
	}

	// Jika tidak ditemukan, kembalikan nilai fallback/default.
	return fallback
}
