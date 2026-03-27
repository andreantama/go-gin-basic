# 📁 internal/delivery

Package `delivery` berisi **interface adapters** yang bertanggung jawab untuk menerima input dari dunia luar (HTTP request, gRPC, CLI) dan mengembalikan output dalam format yang sesuai.

---

## 🎯 Tanggung Jawab

- Menerima **HTTP request** dari client.
- Melakukan **parsing dan validasi** input (request body, URL params, query string).
- Memanggil **use case** yang sesuai.
- Mengkonversi hasil use case ke **HTTP response** (JSON).
- **TIDAK** mengandung business logic.

---

## 📄 Struktur File

```
internal/delivery/
└── http/
    └── user_handler.go   # HTTP Handler untuk endpoint /users dan /auth
```

---

## 🌐 REST API Endpoints

### 🔓 Public Endpoints (Tanpa JWT)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/auth/register` | Daftarkan pengguna baru |
| `POST` | `/auth/login` | Login dan dapatkan JWT token |

### 🔒 Protected Endpoints (Memerlukan JWT Token)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/users` | Ambil semua pengguna |
| `GET` | `/users/:id` | Ambil pengguna berdasarkan ID |
| `PUT` | `/users/:id` | Perbarui data pengguna |
| `DELETE` | `/users/:id` | Hapus pengguna |

---

## 📝 Contoh Request & Response

### POST /auth/register

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
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

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
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

### GET /users/:id

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Data pengguna berhasil diambil",
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

---

## 🔄 Alur Request-Response

```
HTTP Request
    ↓
[GIN Router]        → routing ke handler yang tepat
    ↓
[Middleware]        → JWT validation, logging, CORS
    ↓
[HTTP Handler]      → parse input, panggil use case, format output
    ↓
[Use Case]          → business logic
    ↓
[Repository]        → akses database
    ↓
HTTP Response
```

---

## 🏛️ Posisi dalam Clean Architecture

```
┌─────────────────────────────────────┐
│         Frameworks & Drivers        │
├─────────────────────────────────────┤
│ ► Interface Adapters ◄              │  ← Delivery ada di sini
│   (HTTP Handlers)                   │
├─────────────────────────────────────┤
│           Application               │
├─────────────────────────────────────┤
│         Domain / Entities           │
└─────────────────────────────────────┘
```
