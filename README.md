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
│   └── middleware/                  # ← Middleware HTTP
│       ├── README.md
│       └── auth_middleware.go       # JWT auth, admin-only, CORS middleware
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
        └── mysql.go                 # Koneksi MySQL, connection pool, auto-migration
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
| **MySQL** | 8.0+ | Database relasional |
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

### 5. Jalankan Aplikasi

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

### 6. 📦 Package Utilitas (`pkg/`)
Helper packages: format respons dan custom error types. → [Baca selengkapnya](pkg/README.md)

### 7. 🏭 Infrastructure (`infrastructure/`)
Konfigurasi **database** dan external services. → [Baca selengkapnya](infrastructure/README.md)

### 8. ⚙️ Config (`config/`)
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

---

## 📖 Referensi Belajar

- [Clean Architecture — Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [GIN Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [JWT Introduction](https://jwt.io/introduction)
- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
