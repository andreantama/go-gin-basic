// cmd/migrate/main.go adalah perintah terpisah untuk menjalankan database migration.
//
// Perintah ini memungkinkan Anda menjalankan migration MySQL atau PostgreSQL
// secara eksplisit, terpisah dari proses menjalankan server aplikasi.
//
// Cara penggunaan:
//
//	# Jalankan migration (driver ditentukan oleh DB_DRIVER di .env)
//	go run cmd/migrate/main.go
//
//	# Jalankan migration dengan arah tertentu (up/down)
//	go run cmd/migrate/main.go -direction=up
//	go run cmd/migrate/main.go -direction=down
package main

import (
	// errors digunakan untuk membandingkan error sentinel dari golang-migrate.
	"errors"

	// flag digunakan untuk membaca argumen command-line.
	"flag"

	// fmt digunakan untuk mencetak pesan ke console.
	"fmt"

	// log digunakan untuk logging fatal error.
	"log"

	// os digunakan untuk keluar dari program dengan kode status.
	"os"

	// config berisi fungsi untuk memuat konfigurasi dari environment.
	"github.com/andreantama/go-gin-basic/config"

	// database berisi fungsi koneksi database dan runner migration.
	"github.com/andreantama/go-gin-basic/infrastructure/database"

	// migrations adalah package internal yang mengekspor embed.FS berisi file SQL.
	"github.com/andreantama/go-gin-basic/migrations"

	// migrate adalah package utama golang-migrate.
	migrate "github.com/golang-migrate/migrate/v4"

	// migmysql adalah driver golang-migrate untuk MySQL.
	migmysql "github.com/golang-migrate/migrate/v4/database/mysql"

	// migpostgres adalah driver golang-migrate untuk PostgreSQL.
	migpostgres "github.com/golang-migrate/migrate/v4/database/postgres"

	// iofs adalah source driver golang-migrate untuk membaca file dari embed.FS.
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	// ─── PARSE FLAG ──────────────────────────────────────────────────────────
	// Flag -direction menentukan arah migration: "up" (default) atau "down".
	// "up"   → jalankan semua migration yang belum dieksekusi.
	// "down" → rollback satu langkah migration terakhir.
	direction := flag.String("direction", "up", `Arah migration: "up" untuk menerapkan migration, "down" untuk rollback`)
	flag.Parse()

	// Validasi nilai direction yang diberikan.
	if *direction != "up" && *direction != "down" {
		log.Fatalf("❌ Nilai -direction tidak valid: %q. Gunakan \"up\" atau \"down\".", *direction)
	}

	// ─── MUAT KONFIGURASI ─────────────────────────────────────────────────────
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Gagal memuat konfigurasi: %v", err)
	}

	// ─── BUAT KONEKSI DATABASE ───────────────────────────────────────────────
	// Buat koneksi ke database sesuai DB_DRIVER di .env (mysql atau postgres).
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database %s: %v", cfg.DBDriver, err)
	}

	// Dapatkan koneksi *sql.DB yang mendasari GORM untuk digunakan golang-migrate.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Gagal mendapatkan SQL DB: %v", err)
	}

	// ─── INISIALISASI DRIVER MIGRATION ───────────────────────────────────────
	var m *migrate.Migrate

	switch cfg.DBDriver {
	case "mysql":
		// Buat database driver golang-migrate untuk MySQL.
		dbDriver, err := migmysql.WithInstance(sqlDB, &migmysql.Config{})
		if err != nil {
			log.Fatalf("❌ Gagal membuat driver migrasi MySQL: %v", err)
		}

		// Buat source driver dari embed.FS, mengarah ke sub-direktori "mysql".
		sourceDriver, err := iofs.New(migrations.MigrationFiles, "mysql")
		if err != nil {
			log.Fatalf("❌ Gagal membuat source migrasi MySQL: %v", err)
		}

		// Buat instance migrate.
		m, err = migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
		if err != nil {
			log.Fatalf("❌ Gagal menginisialisasi migrasi MySQL: %v", err)
		}

	case "postgres":
		// Buat database driver golang-migrate untuk PostgreSQL.
		// DatabaseName wajib diisi agar migration berjalan di database yang benar
		// sesuai dengan DB_NAME yang dikonfigurasi di .env.
		dbDriver, err := migpostgres.WithInstance(sqlDB, &migpostgres.Config{DatabaseName: cfg.DBName})
		if err != nil {
			log.Fatalf("❌ Gagal membuat driver migrasi PostgreSQL: %v", err)
		}

		// Buat source driver dari embed.FS, mengarah ke sub-direktori "postgres".
		sourceDriver, err := iofs.New(migrations.MigrationFiles, "postgres")
		if err != nil {
			log.Fatalf("❌ Gagal membuat source migrasi PostgreSQL: %v", err)
		}

		// Buat instance migrate.
		m, err = migrate.NewWithInstance("iofs", sourceDriver, cfg.DBName, dbDriver)
		if err != nil {
			log.Fatalf("❌ Gagal menginisialisasi migrasi PostgreSQL: %v", err)
		}

	default:
		log.Fatalf("❌ Driver '%s' tidak didukung. Gunakan \"mysql\" atau \"postgres\".", cfg.DBDriver)
	}

	// ─── JALANKAN MIGRATION ──────────────────────────────────────────────────
	fmt.Printf("🗄️  Menjalankan migration %s untuk database %s (%s)...\n",
		*direction, cfg.DBName, cfg.DBDriver)

	switch *direction {
	case "up":
		// Jalankan semua migration yang belum dieksekusi.
		err = m.Up()
	case "down":
		// Rollback satu langkah migration terakhir.
		err = m.Steps(-1)
	}

	// Tangani hasil migration.
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// Cetak pesan error dan keluar dengan kode status 1.
		fmt.Fprintf(os.Stderr, "❌ Migration %s gagal: %v\n", *direction, err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		// Tidak ada migration baru yang perlu dijalankan.
		fmt.Printf("✅ Tidak ada perubahan migration — database %s sudah up-to-date.\n", cfg.DBName)
	} else {
		// Migration berhasil dijalankan.
		fmt.Printf("✅ Migration %s berhasil dijalankan untuk database %s (%s).\n",
			*direction, cfg.DBName, cfg.DBDriver)
	}
}
