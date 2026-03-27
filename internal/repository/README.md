# 📁 internal/repository

Package `repository` berisi **implementasi konkret** dari interface repository yang didefinisikan di lapisan `domain`. Repository bertanggung jawab atas semua operasi akses data ke database.

---

## 🎯 Tanggung Jawab

- Mengimplementasikan interface `domain.UserRepository`.
- Melakukan operasi CRUD ke database menggunakan **GORM**.
- Menerjemahkan query database menjadi objek domain.
- Menyembunyikan detail implementasi database dari lapisan use case.

---

## 📄 Struktur File

```
internal/repository/
└── user_repository.go   # Implementasi UserRepository menggunakan GORM + MySQL
```

---

## 🏛️ Posisi dalam Clean Architecture

```
┌─────────────────────────────────────┐
│         Frameworks & Drivers        │  (GIN, GORM, MySQL)
├─────────────────────────────────────┤
│ ► Interface Adapters ◄              │  ← Repository ada di sini
│   (HTTP Handlers, Repositories)     │
├─────────────────────────────────────┤
│           Application               │  (Use Cases)
├─────────────────────────────────────┤
│         Domain / Entities           │
└─────────────────────────────────────┘
```

---

## 📝 Cara Penggunaan

```go
// 1. Inisialisasi koneksi database GORM
db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})

// 2. Buat instance repository dengan dependency injection
userRepo := repository.NewUserRepository(db)

// 3. Gunakan repository di use case
userUC := usecase.NewUserUsecase(userRepo, cfg)
```

---

## 🗃️ Operasi Database

| Method | SQL Query | Deskripsi |
|--------|-----------|-----------|
| `FindAll()` | `SELECT * FROM users WHERE deleted_at IS NULL` | Ambil semua user |
| `FindByID(id)` | `SELECT * FROM users WHERE id = ? LIMIT 1` | Ambil user by ID |
| `FindByEmail(email)` | `SELECT * FROM users WHERE email = ? LIMIT 1` | Ambil user by email |
| `Create(user)` | `INSERT INTO users (...)` | Simpan user baru |
| `Update(user)` | `UPDATE users SET ... WHERE id = ?` | Update data user |
| `Delete(id)` | `UPDATE users SET deleted_at = NOW() WHERE id = ?` | Soft delete user |

---

## 🔄 Dependency Injection Pattern

Repository menggunakan pola **Constructor Injection**:

```go
// Constructor menerima *gorm.DB dan mengembalikan INTERFACE, bukan struct konkret.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
    return &userRepository{db: db}
}
```

Keuntungan mengembalikan interface:
- **Testability**: Bisa di-mock untuk unit testing.
- **Fleksibilitas**: Bisa diganti dengan implementasi lain tanpa mengubah kode use case.
- **Encapsulation**: Detail implementasi (`userRepository`) tersembunyi dari luar.

---

## 🛡️ Soft Delete

Repository menggunakan fitur **Soft Delete** dari GORM. Data tidak dihapus secara permanen, melainkan field `deleted_at` diisi dengan timestamp.

```
HARD DELETE: DELETE FROM users WHERE id = 1  ← Data hilang permanen
SOFT DELETE: UPDATE users SET deleted_at = NOW() WHERE id = 1  ← Data tetap ada
```

Keuntungan soft delete: Data bisa di-restore, audit trail tetap terjaga.
