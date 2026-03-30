// Package database berisi konfigurasi dan inisialisasi koneksi database.
//
// File ini berisi factory function untuk membuat koneksi database
// berdasarkan driver yang dipilih di konfigurasi (DB_DRIVER).
// Saat ini mendukung: MySQL dan PostgreSQL.
package database

import (
	// fmt digunakan untuk memformat pesan error.
	"fmt"

	// config berisi konfigurasi aplikasi termasuk jenis database driver.
	"github.com/andreantama/go-gin-basic/config"

	// gorm adalah ORM yang digunakan untuk berinteraksi dengan database.
	"gorm.io/gorm"
)

// NewDatabaseConnection adalah factory function yang membuat koneksi database
// berdasarkan driver yang dikonfigurasi di DB_DRIVER.
//
// Driver yang didukung:
//   - "mysql": koneksi ke database MySQL (default).
//   - "postgres": koneksi ke database PostgreSQL.
//
// Parameter:
//   - cfg: konfigurasi aplikasi yang berisi jenis driver dan DSN database.
//
// Mengembalikan:
//   - *gorm.DB: instance koneksi database yang siap digunakan.
//   - error: error jika driver tidak didukung atau koneksi gagal.
//
// Contoh penggunaan:
//
//	db, err := database.NewDatabaseConnection(cfg)
//	if err != nil {
//	    log.Fatal("Gagal terhubung ke database:", err)
//	}
func NewDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {
	// Pilih driver database berdasarkan konfigurasi.
	switch cfg.DBDriver {
	case "mysql":
		// Gunakan koneksi MySQL.
		return NewMySQLConnection(cfg)
	case "postgres":
		// Gunakan koneksi PostgreSQL.
		return NewPostgresConnection(cfg)
	default:
		// Kembalikan error jika driver tidak dikenal.
		return nil, fmt.Errorf("database driver '%s' tidak didukung, gunakan 'mysql' atau 'postgres'", cfg.DBDriver)
	}
}
