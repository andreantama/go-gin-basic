# 📁 config

Package `config` bertanggung jawab untuk **membaca dan menyediakan konfigurasi** aplikasi dari environment variable atau file `.env`.

---

## 🎯 Tanggung Jawab

- Memuat konfigurasi dari file `.env` menggunakan `godotenv`.
- Menyediakan struct `Config` yang berisi semua nilai konfigurasi.
- Menyediakan helper `DSN()` untuk menghasilkan connection string database.
- Memvalidasi konfigurasi kritis (misalnya JWT secret di production).

---

## 📄 Struktur File

```
config/
└── config.go   # Definisi struct Config dan fungsi LoadConfig()
```

---

## ⚙️ Environment Variables

Buat file `.env` di root proyek berdasarkan `.env.example`:

| Variable         | Default                               | Keterangan                          |
|------------------|---------------------------------------|-------------------------------------|
| `APP_NAME`       | `go-gin-clean-arch`                   | Nama aplikasi                       |
| `APP_ENV`        | `development`                         | Environment: dev/staging/production |
| `APP_PORT`       | `8080`                                | Port server HTTP                    |
| `DB_HOST`        | `localhost`                           | Host database MySQL                 |
| `DB_PORT`        | `3306`                                | Port database MySQL                 |
| `DB_USER`        | `root`                                | Username database                   |
| `DB_PASSWORD`    | *(kosong)*                            | Password database                   |
| `DB_NAME`        | `go_gin_db`                           | Nama database                       |
| `JWT_SECRET`     | `default-secret-change-in-production` | Secret untuk JWT (wajib di prod)    |
| `JWT_EXPIRE_HOURS`| `24`                                 | Masa berlaku JWT dalam jam          |

---

## 🚀 Cara Penggunaan

```go
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

fmt.Println(cfg.AppPort)  // "8080"
fmt.Println(cfg.DSN())    // "root:pass@tcp(localhost:3306)/go_gin_db?..."
```

---

## 🏛️ Posisi dalam Clean Architecture

```
[ Frameworks & Drivers ] ← config berada di sini
        ↓
[ Interface Adapters ]
        ↓
[ Use Cases ]
        ↓
[ Domain / Entities ]
```

Package `config` termasuk dalam lapisan **Frameworks & Drivers** karena ia berinteraksi langsung dengan sistem operasi (environment variable) dan library eksternal (`godotenv`).
