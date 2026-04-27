// Package database berisi konfigurasi dan inisialisasi koneksi database.
//
// File ini berisi implementasi koneksi untuk database PostgreSQL menggunakan GORM.
// PostgreSQL adalah database relasional open-source yang populer dan powerful.
//
// Catatan: SQL migration tidak dijalankan otomatis saat koneksi dibuat.
// Jalankan migration secara terpisah menggunakan perintah:
//
//	go run cmd/migrate/main.go
package database

import (
	// fmt digunakan untuk memformat pesan log.
	"fmt"

	// log digunakan untuk logging informasi koneksi database.
	"log"

	// time digunakan untuk mengkonfigurasi timeout koneksi database.
	"time"

	// config berisi konfigurasi aplikasi termasuk DSN database.
	"github.com/andreantama/go-gin-basic/config"

	// gorm adalah ORM yang digunakan untuk berinteraksi dengan database.
	"gorm.io/gorm"

	// logger adalah konfigurasi logging untuk GORM.
	"gorm.io/gorm/logger"

	// postgres adalah driver GORM untuk database PostgreSQL.
	"gorm.io/driver/postgres"
)

// NewPostgresConnection membuat dan mengembalikan koneksi database PostgreSQL menggunakan GORM.
// Fungsi ini juga mengkonfigurasi connection pool dan menjalankan SQL migration.
//
// Parameter:
//   - cfg: konfigurasi aplikasi yang berisi DSN database.
//
// Mengembalikan:
//   - *gorm.DB: instance koneksi database yang siap digunakan.
//   - error: error jika koneksi gagal dibuat.
//
// Contoh penggunaan:
//
//	db, err := database.NewPostgresConnection(cfg)
//	if err != nil {
//	    log.Fatal("Gagal terhubung ke database:", err)
//	}
func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	// Konfigurasi GORM — tentukan level logging berdasarkan environment.
	gormConfig := &gorm.Config{}

	// Di environment development, aktifkan logging query SQL untuk debugging.
	if cfg.AppEnv == "development" {
		// Logger.Default menggunakan logger bawaan GORM dengan output ke stdout.
		// logger.Info berarti semua query SQL akan dicatat di log.
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		// Di production/staging, hanya catat query yang lambat atau error.
		// logger.Warn hanya mencatat warning dan error, mengurangi noise di log.
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	}

	// Buka koneksi ke PostgreSQL menggunakan DSN dari konfigurasi.
	// postgres.Open() membuat dialector untuk PostgreSQL.
	// gorm.Open() menggunakan dialector tersebut untuk membuka koneksi.
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), gormConfig)
	if err != nil {
		// Kembalikan error yang lebih deskriptif jika koneksi gagal.
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	// Ambil koneksi database SQL yang mendasari GORM untuk konfigurasi connection pool.
	// GORM menyembunyikan detail ini di balik abstraksi, tapi kita perlu mengaksesnya
	// untuk mengkonfigurasi connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		// Kembalikan error jika gagal mendapatkan instance SQL DB.
		return nil, fmt.Errorf("gagal mendapatkan SQL DB: %w", err)
	}

	// Konfigurasi connection pool untuk performa optimal.
	// Connection pool adalah sekumpulan koneksi database yang siap digunakan kembali,
	// sehingga tidak perlu membuat koneksi baru untuk setiap request.

	// SetMaxIdleConns mengatur jumlah maksimum koneksi idle dalam pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns mengatur jumlah maksimum koneksi yang bisa dibuka ke database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime mengatur berapa lama koneksi boleh digunakan sebelum ditutup.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Log pesan sukses ke console.
	log.Println("Berhasil terhubung ke database PostgreSQL")

	// Kembalikan instance koneksi database yang sudah dikonfigurasi.
	return db, nil
}
