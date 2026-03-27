# 📁 infrastructure

Package `infrastructure` berisi **implementasi teknis** dari lapisan paling luar dalam Clean Architecture — koneksi ke database, external services, dan konfigurasi frameworks.

---

## 🎯 Tanggung Jawab

- Mengkonfigurasi dan membuat koneksi **database** (MySQL via GORM).
- Menjalankan **database migrations** (auto-create/update tabel).
- Mengkonfigurasi **connection pool** untuk performa optimal.

---

## 📄 Struktur File

```
infrastructure/
└── database/
    └── mysql.go   # Koneksi MySQL, connection pool, dan auto-migration
```

---

## 🏛️ Posisi dalam Clean Architecture

```
┌─────────────────────────────────────┐
│ ► Frameworks & Drivers ◄            │  ← Infrastructure ada di sini
│   (GIN, GORM, MySQL)                │
├─────────────────────────────────────┤
│         Interface Adapters          │
├─────────────────────────────────────┤
│           Application               │
├─────────────────────────────────────┤
│         Domain / Entities           │
└─────────────────────────────────────┘
```

---

## 🗄️ Konfigurasi Database

### Connection Pool

| Parameter | Nilai | Deskripsi |
|-----------|-------|-----------|
| `MaxIdleConns` | 10 | Koneksi idle dalam pool |
| `MaxOpenConns` | 100 | Maksimum koneksi aktif |
| `ConnMaxLifetime` | 1 jam | Umur maksimum koneksi |

### Auto Migration

GORM `AutoMigrate` secara otomatis:
- ✅ Membuat tabel baru jika belum ada.
- ✅ Menambahkan kolom baru yang belum ada.
- ❌ **Tidak** menghapus kolom yang sudah tidak ada di struct (untuk keamanan data).

---

## 🚀 Cara Penggunaan

```go
// Di main.go
cfg, _ := config.LoadConfig()
db, err := database.NewMySQLConnection(cfg)
if err != nil {
    log.Fatal(err)
}

// Gunakan db untuk membuat repository
userRepo := repository.NewUserRepository(db)
```

---

## 🔄 Menambah Database Baru

Untuk mendukung database lain (misalnya PostgreSQL), buat file baru:

```
infrastructure/
└── database/
    ├── mysql.go       # Koneksi MySQL
    └── postgres.go    # Koneksi PostgreSQL (baru)
```

```go
// postgres.go
import "gorm.io/driver/postgres"

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

Kemudian di `main.go`, cukup ganti panggilan dari `NewMySQLConnection` ke `NewPostgresConnection` — **use case dan domain tidak perlu diubah sama sekali**.
