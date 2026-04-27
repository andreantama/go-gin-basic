// Package database — file ini berisi implementasi runner untuk golang-migrate.
//
// RunMigrations menjalankan semua file SQL migration yang belum dieksekusi.
// Migration dilacak di tabel "schema_migrations" yang dibuat otomatis oleh golang-migrate.
//
// Keunggulan golang-migrate dibanding GORM AutoMigrate:
//   - Versi migration yang eksplisit dan dapat di-audit (file SQL bernomor urut).
//   - Mendukung rollback (DOWN migration) untuk membatalkan perubahan skema.
//   - Migration hanya dijalankan sekali; status disimpan di tabel schema_migrations.
//   - File SQL mudah di-review, di-version-control, dan dipahami tanpa ORM.
package database

import (
	// errors digunakan untuk membandingkan error sentinel dari golang-migrate.
	"errors"

	// fmt digunakan untuk memformat pesan error yang deskriptif.
	"fmt"

	// log digunakan untuk mencatat status migration ke console.
	"log"

	// migrations adalah package internal yang mengekspor embed.FS berisi file SQL.
	"github.com/andreantama/go-gin-basic/migrations"

	// migrate adalah package utama golang-migrate untuk menjalankan migration.
	migrate "github.com/golang-migrate/migrate/v4"

	// migmysql adalah driver golang-migrate untuk database MySQL.
	migmysql "github.com/golang-migrate/migrate/v4/database/mysql"

	// migpostgres adalah driver golang-migrate untuk database PostgreSQL.
	migpostgres "github.com/golang-migrate/migrate/v4/database/postgres"

	// iofs adalah source driver golang-migrate untuk membaca file dari embed.FS.
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// gorm diperlukan untuk mendapatkan *sql.DB yang mendasari koneksi GORM.
	"gorm.io/gorm"
)

// RunMigrations menjalankan semua pending SQL migration menggunakan golang-migrate.
//
// Fungsi ini:
//  1. Mendapatkan koneksi *sql.DB yang mendasari *gorm.DB.
//  2. Membuat source driver dari embed.FS yang berisi file SQL.
//  3. Membuat database driver sesuai jenis database (mysql/postgres).
//  4. Menjalankan semua migration yang belum dieksekusi (UP).
//
// Jika semua migration sudah dieksekusi sebelumnya (ErrNoChange),
// fungsi ini tetap dianggap sukses.
//
// Parameter:
//   - db:     instance *gorm.DB yang sudah terkoneksi ke database.
//   - driver: jenis database — "mysql" atau "postgres".
//   - dbName: nama database yang digunakan (diperlukan untuk PostgreSQL agar migration
//     berjalan di database yang benar sesuai konfigurasi DB_NAME di .env).
//
// Mengembalikan error jika migration gagal dijalankan.
func RunMigrations(db *gorm.DB, driver string, dbName string) error {
	// Dapatkan koneksi *sql.DB yang mendasari GORM.
	// Diperlukan karena golang-migrate bekerja langsung dengan database/sql.
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan SQL DB: %w", err)
	}

	switch driver {
	case "mysql":
		// Buat database driver golang-migrate untuk MySQL menggunakan koneksi yang sudah ada.
		// WithInstance menghindari pembuatan koneksi baru dan menggunakan pool yang sama.
		dbDriver, err := migmysql.WithInstance(sqlDB, &migmysql.Config{})
		if err != nil {
			return fmt.Errorf("gagal membuat driver migrasi MySQL: %w", err)
		}

		// Buat source driver dari embed.FS, mengarah ke sub-direktori "mysql".
		sourceDriver, err := iofs.New(migrations.MigrationFiles, "mysql")
		if err != nil {
			return fmt.Errorf("gagal membuat source migrasi MySQL: %w", err)
		}

		// Buat instance migrate dengan source dan database driver yang sudah dibuat.
		m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
		if err != nil {
			return fmt.Errorf("gagal menginisialisasi migrasi MySQL: %w", err)
		}

		// Jalankan semua migration yang belum dieksekusi.
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("gagal menjalankan migrasi MySQL: %w", err)
		}

	case "postgres":
		// Buat database driver golang-migrate untuk PostgreSQL menggunakan koneksi yang sudah ada.
		// DatabaseName wajib diisi agar golang-migrate menggunakan database yang benar
		// sesuai dengan DB_NAME yang dikonfigurasi di .env.
		dbDriver, err := migpostgres.WithInstance(sqlDB, &migpostgres.Config{DatabaseName: dbName})
		if err != nil {
			return fmt.Errorf("gagal membuat driver migrasi PostgreSQL: %w", err)
		}

		// Buat source driver dari embed.FS, mengarah ke sub-direktori "postgres".
		sourceDriver, err := iofs.New(migrations.MigrationFiles, "postgres")
		if err != nil {
			return fmt.Errorf("gagal membuat source migrasi PostgreSQL: %w", err)
		}

		// Buat instance migrate dengan source dan database driver yang sudah dibuat.
		m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
		if err != nil {
			return fmt.Errorf("gagal menginisialisasi migrasi PostgreSQL: %w", err)
		}

		// Jalankan semua migration yang belum dieksekusi.
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("gagal menjalankan migrasi PostgreSQL: %w", err)
		}

	default:
		return fmt.Errorf("driver '%s' tidak didukung untuk migrasi", driver)
	}

	// Cetak pesan sukses ke console.
	log.Println("Database migration berhasil dijalankan")
	return nil
}
