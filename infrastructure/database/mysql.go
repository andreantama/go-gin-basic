// Package database berisi konfigurasi dan inisialisasi koneksi database.
//
// Lapisan ini termasuk dalam "Frameworks & Drivers" dalam Clean Architecture.
// Package ini bertanggung jawab untuk:
//   - Membuka koneksi ke database MySQL menggunakan GORM.
//   - Mengkonfigurasi connection pool.
//   - Menjalankan SQL migration menggunakan golang-migrate.
//
// Dengan memisahkan setup database ke package tersendiri, kita bisa
// mengganti database atau ORM tanpa mengubah lapisan bisnis.
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

	// mysql adalah driver GORM untuk database MySQL.
	"gorm.io/driver/mysql"
)

// NewMySQLConnection membuat dan mengembalikan koneksi database MySQL menggunakan GORM.
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
//	db, err := database.NewMySQLConnection(cfg)
//	if err != nil {
//	    log.Fatal("Gagal terhubung ke database:", err)
//	}
func NewMySQLConnection(cfg *config.Config) (*gorm.DB, error) {
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

	// Buka koneksi ke MySQL menggunakan DSN dari konfigurasi.
	// mysql.Open() membuat dialector untuk MySQL.
	// gorm.Open() menggunakan dialector tersebut untuk membuka koneksi.
	db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
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
	// Koneksi idle adalah koneksi yang terbuka tapi tidak sedang digunakan.
	// Nilai 10 adalah titik awal yang baik untuk sebagian besar aplikasi.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns mengatur jumlah maksimum koneksi yang bisa dibuka ke database.
	// Ini mencegah overload pada database server.
	// Nilai 100 cocok untuk aplikasi dengan traffic sedang.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime mengatur berapa lama koneksi boleh digunakan sebelum ditutup.
	// Ini mencegah masalah dengan koneksi yang sudah tidak valid karena timeout dari sisi database.
	// Nilai 1 jam adalah nilai yang umum digunakan.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Jalankan SQL migration menggunakan golang-migrate.
	// Migration membaca file SQL dari folder migrations/mysql/ yang di-embed ke binary.
	// Setiap migration hanya dijalankan sekali; status dilacak di tabel schema_migrations.
	if err := RunMigrations(db, "mysql"); err != nil {
		// Kembalikan error jika migration gagal.
		return nil, fmt.Errorf("gagal menjalankan database migration: %w", err)
	}

	// Log pesan sukses ke console.
	log.Println("Berhasil terhubung ke database MySQL")

	// Kembalikan instance koneksi database yang sudah dikonfigurasi.
	return db, nil
}
