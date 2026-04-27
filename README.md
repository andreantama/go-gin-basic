# 🚀 go-gin-basic — Clean Architecture dengan GIN Framework

Proyek ini adalah contoh implementasi **Clean Architecture** menggunakan **Go** dan **GIN Framework** yang sering dipakai di industri enterprise backend. Setiap baris kode dilengkapi dengan komentar penjelasan untuk membantu pemula memahami konsep yang diterapkan.

---

## 📋 Daftar Isi

- [Apa itu Clean Architecture?](#apa-itu-clean-architecture)
- [Struktur Folder](#struktur-folder)
- [Diagram Arsitektur](#diagram-arsitektur)
- [Teknologi yang Digunakan](#teknologi-yang-digunakan)
- [Cara Menjalankan](#cara-menjalankan)
- [REST API Endpoints](#rest-api-endpoints)
- [Contoh Request & Response](#contoh-request--response)
- [Task Scheduler](#-task-scheduler)
- [Database Migrations](#-database-migrations)
- [Penjelasan Setiap Layer](#penjelasan-setiap-layer)

---

## 🏛️ Apa itu Clean Architecture?

Clean Architecture adalah sebuah pendekatan desain software yang dipopulerkan oleh **Robert C. Martin (Uncle Bob)**. Tujuan utamanya adalah membuat kode yang:

- ✅ **Independen dari framework** — Logika bisnis tidak bergantung pada GIN, GORM, atau library apapun.
- ✅ **Mudah di-test** — Setiap lapisan bisa di-test secara terisolasi menggunakan mock.
- ✅ **Independen dari database** — Bisa ganti MySQL ke PostgreSQL tanpa mengubah logika bisnis.
- ✅ **Independen dari UI** — Bisa menambah gRPC/CLI tanpa mengubah use case.
- ✅ **Mudah dipelihara** — Perubahan di satu lapisan tidak mempengaruhi lapisan lain.

### Aturan Utama: Dependency Rule

```
Lapisan luar BOLEH bergantung pada lapisan dalam.
Lapisan dalam TIDAK BOLEH bergantung pada lapisan luar.
```

---

## 📁 Struktur Folder

```
go-gin-basic/
│
├── main.go                          # Entry point — inisialisasi dan dependency injection
├── go.mod                           # Definisi modul dan dependensi
├── .env.example                     # Template konfigurasi environment
│
├── cmd/                             # Perintah CLI yang bisa dijalankan secara terpisah
│   └── migrate/
│       └── main.go                  # Perintah untuk menjalankan database migration
│
├── migrations/                      # File SQL migration (golang-migrate)
│   ├── README.md
│   ├── embed.go                     # Ekspor embed.FS berisi file SQL
│   ├── mysql/                       # Migration khusus MySQL
│   │   ├── 000001_create_users_table.up.sql
│   │   └── 000001_create_users_table.down.sql
│   └── postgres/                    # Migration khusus PostgreSQL
│       ├── 000001_create_users_table.up.sql
│       └── 000001_create_users_table.down.sql
│
├── config/                          # Konfigurasi aplikasi
│   ├── README.md
│   └── config.go                    # LoadConfig(), struct Config, DSN()
│
├── internal/                        # Kode privat aplikasi ini
│   ├── domain/                      # ← Lapisan paling dalam (Entities)
│   │   ├── README.md
│   │   └── user.go                  # Struct User, interface UserRepository, UserUsecase
│   │
│   ├── usecase/                     # ← Lapisan Use Cases (Application Business Rules)
│   │   ├── README.md
│   │   └── user_usecase.go          # Business logic: register, login, CRUD user
│   │
│   ├── repository/                  # ← Interface Adapters (data access)
│   │   ├── README.md
│   │   └── user_repository.go       # Implementasi GORM untuk UserRepository
│   │
│   ├── delivery/                    # ← Interface Adapters (HTTP)
│   │   ├── README.md
│   │   └── http/
│   │       └── user_handler.go      # GIN HTTP handlers untuk /users dan /auth
│   │
│   ├── middleware/                  # ← Middleware HTTP
│   │   ├── README.md
│   │   └── auth_middleware.go       # JWT auth, admin-only, CORS middleware
│   │
│   └── scheduler/                  # ← Task Scheduler (mirip Laravel)
│       ├── README.md
│       ├── scheduler.go             # Engine scheduler: registrasi, start, stop
│       └── task.go                  # Definisi Task: nama, jadwal, fungsi, metadata
│
├── pkg/                             # Utility packages yang bisa dipakai ulang
│   ├── README.md
│   ├── response/
│   │   └── response.go              # Format respons JSON yang konsisten
│   └── errors/
│       └── errors.go                # Custom error types dan HTTP status mapper
│
└── infrastructure/                  # ← Lapisan paling luar (Frameworks & Drivers)
    ├── README.md
    └── database/
        ├── database.go              # Factory function NewDatabaseConnection
        ├── migrate.go               # Runner migration menggunakan golang-migrate
        ├── mysql.go                 # Koneksi MySQL dan connection pool
        └── postgres.go              # Koneksi PostgreSQL dan connection pool
```

---

## 🔄 Diagram Arsitektur

```
┌──────────────────────────────────────────────────────────────────┐
│                    FRAMEWORKS & DRIVERS                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   GIN (HTTP)    │  │  GORM (ORM)     │  │  MySQL (DB)     │  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  │
└───────────┼────────────────────┼────────────────────┼───────────┘
            │                    │                    │
┌───────────▼────────────────────▼────────────────────▼───────────┐
│                    INTERFACE ADAPTERS                             │
│  ┌─────────────────────────┐  ┌──────────────────────────────┐   │
│  │    HTTP Handlers        │  │    Repository (GORM impl)    │   │
│  │  (delivery/http/)       │  │    (repository/)             │   │
│  └─────────────┬───────────┘  └───────────────┬──────────────┘   │
└────────────────┼──────────────────────────────┼─────────────────┘
                 │                              │
┌────────────────▼──────────────────────────────▼─────────────────┐
│                    APPLICATION / USE CASES                        │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    User Use Case                            │ │
│  │               (usecase/user_usecase.go)                     │ │
│  └─────────────────────────────┬───────────────────────────────┘ │
└────────────────────────────────┼────────────────────────────────┘
                                 │
┌────────────────────────────────▼────────────────────────────────┐
│                    DOMAIN / ENTITIES                              │
│  ┌──────────────────────┐  ┌───────────────────────────────────┐ │
│  │    User Entity       │  │  Repository & UseCase Interfaces  │ │
│  │   (domain/user.go)   │  │       (domain/user.go)            │ │
│  └──────────────────────┘  └───────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Teknologi yang Digunakan

| Teknologi | Versi | Kegunaan |
|-----------|-------|---------|
| **Go** | 1.21+ | Bahasa pemrograman utama |
| **GIN** | v1.9.1 | HTTP Web Framework |
| **GORM** | v1.25.7 | ORM untuk akses database |
| **golang-migrate** | v4.19.1 | Manajemen SQL migration yang terkontrol |
| **MySQL** | 8.0+ | Database relasional |
| **PostgreSQL** | 14+ | Database relasional alternatif |
| **golang-jwt** | v5.2.1 | Pembuatan dan validasi JWT token |
| **bcrypt** | - | Hashing password yang aman |
| **godotenv** | v1.5.1 | Membaca file `.env` |

---

## 🚀 Cara Menjalankan

### Prasyarat

- Go 1.21 atau lebih baru: [Download Go](https://go.dev/dl/)
- MySQL 8.0 atau lebih baru
- Git

### 1. Clone Repository

```bash
git clone https://github.com/andreantama/go-gin-basic.git
cd go-gin-basic
```

### 2. Install Dependensi

```bash
go mod tidy
```

### 3. Konfigurasi Environment

```bash
# Salin template konfigurasi
cp .env.example .env

# Edit file .env sesuai konfigurasi lokal Anda
nano .env
```

Isi file `.env`:

```env
APP_NAME=go-gin-clean-arch
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go_gin_db

JWT_SECRET=your-very-strong-secret-key
JWT_EXPIRE_HOURS=24
```

### 4. Buat Database MySQL

```sql
CREATE DATABASE go_gin_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. Jalankan Migration Database

Migration **tidak** berjalan otomatis saat `go run main.go`. Anda perlu menjalankannya secara terpisah menggunakan perintah `cmd/migrate`.

#### MySQL

Pastikan `.env` sudah dikonfigurasi dengan `DB_DRIVER=mysql`, lalu:

```bash
go run cmd/migrate/main.go
```

#### PostgreSQL

Pastikan `.env` sudah dikonfigurasi dengan `DB_DRIVER=postgres`, lalu:

```bash
go run cmd/migrate/main.go
```

#### Rollback Migration (DOWN)

```bash
go run cmd/migrate/main.go -direction=down
```

Output yang diharapkan (contoh MySQL):

```
🗄️  Menjalankan migration up untuk database go_gin_db (mysql)...
✅ Migration up berhasil dijalankan untuk database go_gin_db (mysql).
```

### 6. Jalankan Aplikasi

```bash
go run main.go
```

Output yang diharapkan:

```
🚀 Server go-gin-clean-arch berjalan di http://localhost:8080
📌 Environment: development
[GIN-debug] POST   /auth/register
[GIN-debug] POST   /auth/login
[GIN-debug] GET    /users
[GIN-debug] GET    /users/:id
[GIN-debug] PUT    /users/:id
[GIN-debug] DELETE /users/:id
```

---

## 🌐 REST API Endpoints

### 🔓 Public Endpoints (Tanpa Autentikasi)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/auth/register` | Daftarkan pengguna baru |
| `POST` | `/auth/login` | Login dan dapatkan JWT token |

### 🔒 Protected Endpoints (Memerlukan `Authorization: Bearer <token>`)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/users` | Ambil daftar semua pengguna |
| `GET` | `/users/:id` | Ambil pengguna berdasarkan ID |
| `PUT` | `/users/:id` | Perbarui data pengguna |
| `DELETE` | `/users/:id` | Hapus pengguna (soft delete) |

---

## 📝 Contoh Request & Response

### POST /auth/register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### POST /auth/login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### GET /users (dengan JWT Token)

```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Daftar pengguna berhasil diambil",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### PUT /users/:id

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'
```

### DELETE /users/:id

```bash
curl -X DELETE http://localhost:8080/users/1 \
  -H "Authorization: Bearer <token>"
```

---

## ⏰ Task Scheduler

Proyek ini dilengkapi dengan fitur **Task Scheduler** yang terinspirasi dari **Laravel Task Scheduler**. Fitur ini memungkinkan Anda mendaftarkan tugas-tugas yang berjalan secara otomatis di background sesuai jadwal.

### Fitur Utama

- 🕐 **API ekspresif** mirip Laravel: `EveryMinute()`, `Hourly()`, `Daily()`, `Weekly()`, dll.
- 📝 **Cron expression** kustom dengan dukungan format 6-field (detik).
- 🛡️ **Pencegahan overlap** — tugas yang sama tidak berjalan bersamaan.
- 🔄 **Recovery dari panic** — satu tugas gagal tidak menghentikan scheduler.
- 🛑 **Graceful shutdown** — menunggu tugas selesai saat aplikasi berhenti.

### Contoh Penggunaan

```go
// Di main.go
taskScheduler := scheduler.NewScheduler()

// Jalankan setiap menit
taskScheduler.EveryMinute("health-check", func() error {
    log.Println("Server berjalan normal")
    return nil
})

// Jalankan setiap hari pukul 02:00
taskScheduler.DailyAt("cleanup-data", 2, 0, func() error {
    log.Println("Membersihkan data lama...")
    return nil
})

// Jalankan setiap jam
taskScheduler.Hourly("sync-data", func() error {
    log.Println("Sinkronisasi data...")
    return nil
})

// Cron expression kustom (setiap hari Senin pukul 09:00)
taskScheduler.Cron("weekly-report", "0 0 9 * * 1", func() error {
    log.Println("Membuat laporan mingguan...")
    return nil
})

taskScheduler.Start()
defer taskScheduler.Stop()
```

### Daftar Helper Method

| Method | Jadwal | Setara Laravel |
|--------|--------|----------------|
| `EverySecond()` | Setiap detik | `everySecond()` |
| `EveryFiveSeconds()` | Setiap 5 detik | `everyFiveSeconds()` |
| `EveryMinute()` | Setiap menit | `everyMinute()` |
| `EveryFiveMinutes()` | Setiap 5 menit | `everyFiveMinutes()` |
| `EveryFifteenMinutes()` | Setiap 15 menit | `everyFifteenMinutes()` |
| `EveryThirtyMinutes()` | Setiap 30 menit | `everyThirtyMinutes()` |
| `Hourly()` | Setiap jam | `hourly()` |
| `Daily()` | Setiap hari (00:00) | `daily()` |
| `DailyAt(hour, minute)` | Setiap hari pada jam tertentu | `dailyAt('13:00')` |
| `Weekly()` | Setiap Minggu | `weekly()` |
| `Monthly()` | Setiap tanggal 1 | `monthly()` |
| `Cron(expression)` | Custom cron expression | `cron('* * * * *')` |

→ [Baca dokumentasi lengkap Scheduler](internal/scheduler/README.md)

---

## 🗄️ Database Migrations

Proyek ini menggunakan **golang-migrate** untuk manajemen skema database yang terkontrol. Migration sebelumnya menggunakan GORM AutoMigrate telah digantikan untuk memberikan kontrol penuh atas perubahan skema database.

### Mengapa golang-migrate?

| Fitur | GORM AutoMigrate (lama) | golang-migrate (sekarang) |
|-------|------------------------|--------------------------|
| File SQL eksplisit | ❌ | ✅ |
| Rollback (DOWN) | ❌ | ✅ |
| Audit trail | ❌ | ✅ |
| Review SQL sebelum deploy | ❌ | ✅ |
| Idempotent (hanya sekali) | ⚠️ | ✅ |

### Cara Kerja

Migration **tidak** dijalankan otomatis saat aplikasi start. Migration dijalankan secara eksplisit menggunakan perintah `cmd/migrate`. File SQL di-embed ke dalam binary menggunakan Go `embed`, sehingga tidak perlu file terpisah saat deployment.

```
go run cmd/migrate/main.go
     │
     ▼
cmd/migrate/main.go
     │
     ├── Baca DB_DRIVER dari .env (mysql atau postgres)
     ├── Buat koneksi ke database sesuai DB_NAME
     ├── Cek tabel schema_migrations
     ├── Jalankan migration yang belum dieksekusi
     └── Catat versi yang sudah dijalankan
```

### Menjalankan Migration MySQL

1. Set `DB_DRIVER=mysql` di file `.env`
2. Pastikan database sudah dibuat:
   ```sql
   CREATE DATABASE go_gin_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
3. Jalankan migration:
   ```bash
   go run cmd/migrate/main.go
   ```

### Menjalankan Migration PostgreSQL

1. Set `DB_DRIVER=postgres` di file `.env`
2. Pastikan database sudah dibuat:
   ```sql
   CREATE DATABASE go_gin_db;
   ```
3. Jalankan migration:
   ```bash
   go run cmd/migrate/main.go
   ```

### Rollback Migration

```bash
# Rollback satu langkah terakhir
go run cmd/migrate/main.go -direction=down
```

### Flag yang Tersedia

| Flag | Nilai | Default | Deskripsi |
|------|-------|---------|-----------|
| `-direction` | `up` / `down` | `up` | Arah migration: terapkan atau rollback |

### Struktur Migration

```
migrations/
├── mysql/
│   ├── 000001_create_users_table.up.sql    # Membuat tabel
│   └── 000001_create_users_table.down.sql  # Rollback
└── postgres/
    ├── 000001_create_users_table.up.sql
    └── 000001_create_users_table.down.sql
```

### Menambahkan Migration Baru

Buat pasangan file UP dan DOWN dengan nomor versi berikutnya:

```bash
# MySQL
migrations/mysql/000002_add_phone_to_users.up.sql
migrations/mysql/000002_add_phone_to_users.down.sql

# PostgreSQL
migrations/postgres/000002_add_phone_to_users.up.sql
migrations/postgres/000002_add_phone_to_users.down.sql
```

→ [Baca dokumentasi lengkap Migration](migrations/README.md)

---

## 📚 Penjelasan Setiap Layer

### 1. 🏗️ Domain (`internal/domain/`)
Lapisan **paling dalam** yang berisi entitas bisnis dan kontrak (interface). Tidak bergantung pada apapun. → [Baca selengkapnya](internal/domain/README.md)

### 2. 🧠 Use Case (`internal/usecase/`)
Berisi **logika bisnis** aplikasi. Hanya bergantung pada interface dari domain. → [Baca selengkapnya](internal/usecase/README.md)

### 3. 🗄️ Repository (`internal/repository/`)
Implementasi **akses data** ke database menggunakan GORM. → [Baca selengkapnya](internal/repository/README.md)

### 4. 🌐 Delivery (`internal/delivery/`)
**HTTP handlers** yang menangani request dan response menggunakan GIN. → [Baca selengkapnya](internal/delivery/README.md)

### 5. 🔐 Middleware (`internal/middleware/`)
**Middleware** untuk autentikasi JWT, otorisasi role, dan CORS. → [Baca selengkapnya](internal/middleware/README.md)

### 6. ⏰ Scheduler (`internal/scheduler/`)
**Task Scheduler** untuk penjadwalan tugas otomatis (mirip Laravel Scheduler). → [Baca selengkapnya](internal/scheduler/README.md)

### 7. 📦 Package Utilitas (`pkg/`)
Helper packages: format respons dan custom error types. → [Baca selengkapnya](pkg/README.md)

### 8. 🏭 Infrastructure (`infrastructure/`)
Konfigurasi **database** dan external services. → [Baca selengkapnya](infrastructure/README.md)

### 9. 🗄️ Migrations (`migrations/`)
File SQL migration menggunakan **golang-migrate** untuk manajemen skema database. → [Baca selengkapnya](migrations/README.md)

### 10. ⚙️ Config (`config/`)
Pembacaan **konfigurasi** dari environment variable. → [Baca selengkapnya](config/README.md)

---

## 🔑 Konsep Penting yang Dipelajari

| Konsep | Implementasi |
|--------|-------------|
| **Dependency Injection** | Constructor functions menerima dependency sebagai parameter |
| **Dependency Inversion** | Use case bergantung pada interface, bukan implementasi konkret |
| **Interface Segregation** | Interface kecil dan spesifik (`UserRepository`, `UserUsecase`) |
| **Soft Delete** | Data tidak dihapus permanen, `deleted_at` diisi timestamp |
| **JWT Authentication** | Token berisi claims user (ID, email, role) yang diverifikasi |
| **Password Hashing** | `bcrypt` untuk hashing yang aman, tidak bisa di-reverse |
| **Connection Pool** | GORM connection pool untuk performa database yang optimal |
| **Middleware Chain** | GIN middleware dieksekusi berurutan sebelum handler |
| **SQL Migration** | golang-migrate dengan file SQL bernomor versi dan dukungan rollback |
| **Task Scheduling** | Penjadwalan tugas otomatis mirip Laravel Scheduler |

---

## 📖 Referensi Belajar

- [Clean Architecture — Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [GIN Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [JWT Introduction](https://jwt.io/introduction)
- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [Go embed Package](https://pkg.go.dev/embed)
