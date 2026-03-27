# 📁 internal/usecase

Package `usecase` berisi **logika bisnis aplikasi** (Application Business Rules). Lapisan ini adalah "otak" dari aplikasi yang mengorkestrasi alur data antara domain dan interface.

---

## 🎯 Tanggung Jawab

- Mengimplementasikan interface `domain.UserUsecase`.
- Menerapkan aturan bisnis seperti validasi email unik, hashing password.
- Menggunakan `domain.UserRepository` (interface) untuk akses data — **bukan** implementasi konkret.
- **Tidak** mengenal HTTP, JSON, atau framework web apapun.

---

## 📄 Struktur File

```
internal/usecase/
└── user_usecase.go   # Implementasi business logic untuk User
```

---

## 🏛️ Posisi dalam Clean Architecture

```
┌─────────────────────────────────────┐
│         Frameworks & Drivers        │
├─────────────────────────────────────┤
│         Interface Adapters          │
├─────────────────────────────────────┤
│ ► Application / Use Cases ◄         │  ← Kita ada di sini
├─────────────────────────────────────┤
│         Domain / Entities           │
└─────────────────────────────────────┘
```

---

## 📋 Business Logic yang Diterapkan

### `RegisterUser`
1. ✅ Validasi email belum terdaftar (cek duplikasi).
2. 🔒 Hash password menggunakan `bcrypt` sebelum disimpan.
3. 👤 Set role default `"user"` jika tidak disediakan.

### `LoginUser`
1. 🔍 Cari user berdasarkan email.
2. 🔑 Bandingkan password dengan hash di database menggunakan `bcrypt`.
3. 🎫 Generate JWT token jika kredensial valid.
4. 🛡️ Gunakan pesan error generik untuk mencegah *user enumeration attack*.

### `UpdateUser`
1. ✅ Pastikan user dengan ID tersebut ada.
2. 📝 Update hanya field yang diberikan (partial update).
3. 🔒 Hash password baru jika disediakan.

### `DeleteUser`
1. ✅ Pastikan user ada sebelum menghapus.
2. 🗑️ Delegasikan soft delete ke repository.

---

## 🔐 JWT Token

Token JWT yang dibuat berisi claims berikut:

```json
{
  "sub": 1,
  "email": "user@example.com",
  "role": "user",
  "exp": 1234567890,
  "iat": 1234567890
}
```

Token ditandatangani dengan algoritma **HS256** menggunakan secret key dari konfigurasi.

---

## 🧪 Testability

Use case mudah di-test karena hanya bergantung pada interface:

```go
// Buat mock repository untuk testing
type mockUserRepo struct { ... }
func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) { ... }

// Inject mock ke use case
uc := usecase.NewUserUsecase(&mockUserRepo{}, cfg)

// Test logika bisnis tanpa database nyata
user, err := uc.RegisterUser(&domain.User{...})
```

---

## 💡 Kenapa Use Case Tidak Bergantung pada Framework?

Jika use case bergantung langsung pada GIN atau database, maka:
- Sulit untuk di-test (harus setup server dan database sungguhan).
- Tidak bisa mengganti framework tanpa mengubah logika bisnis.
- Logika bisnis "terkontaminasi" dengan detail implementasi.

Dengan Clean Architecture, use case hanya bergantung pada abstraksi (interface) sehingga bisa di-test secara terisolasi.
