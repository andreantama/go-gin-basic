# 📁 internal/domain

Package `domain` berisi **inti bisnis aplikasi** — entitas dan kontrak interface. Ini adalah lapisan paling dalam dalam Clean Architecture.

---

## 🎯 Tanggung Jawab

- Mendefinisikan **struct entitas** (objek bisnis seperti `User`).
- Mendefinisikan **interface repository** (kontrak akses data).
- Mendefinisikan **interface use case** (kontrak logika bisnis).
- **TIDAK** bergantung pada framework, database, atau lapisan lain.

---

## 📄 Struktur File

```
internal/domain/
└── user.go   # Entitas User, UserRepository interface, UserUsecase interface
```

---

## 🏛️ Prinsip Clean Architecture

```
┌─────────────────────────────────────┐
│         Frameworks & Drivers        │  (GIN, GORM, MySQL)
├─────────────────────────────────────┤
│         Interface Adapters          │  (HTTP Handlers, Repositories)
├─────────────────────────────────────┤
│           Application               │  (Use Cases)
├─────────────────────────────────────┤
│    ► DOMAIN / ENTITIES ◄            │  ← Kita ada di sini
└─────────────────────────────────────┘
```

**Aturan Ketergantungan (Dependency Rule):**
- Panah ketergantungan selalu mengarah KE DALAM.
- Domain tidak tahu apa pun tentang lapisan di atasnya.
- Lapisan luar boleh menggunakan domain, tapi tidak sebaliknya.

---

## 📝 Entitas: `User`

```go
type User struct {
    ID        uint       `json:"id"`
    Name      string     `json:"name"`
    Email     string     `json:"email"`
    Password  string     `json:"-"`     // Disembunyikan dari JSON response
    Role      string     `json:"role"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

---

## 📋 Interface: `UserRepository`

Kontrak yang harus diimplementasikan oleh lapisan **infrastructure**:

| Method | Deskripsi |
|--------|-----------|
| `FindAll()` | Ambil semua pengguna |
| `FindByID(id)` | Ambil pengguna berdasarkan ID |
| `FindByEmail(email)` | Ambil pengguna berdasarkan email |
| `Create(user)` | Simpan pengguna baru |
| `Update(user)` | Perbarui data pengguna |
| `Delete(id)` | Hapus pengguna (soft delete) |

---

## 📋 Interface: `UserUsecase`

Kontrak yang harus diimplementasikan oleh lapisan **usecase**:

| Method | Deskripsi |
|--------|-----------|
| `GetAllUsers()` | Ambil semua pengguna |
| `GetUserByID(id)` | Ambil pengguna berdasarkan ID |
| `RegisterUser(user)` | Daftarkan pengguna baru |
| `LoginUser(email, pass)` | Login dan dapatkan JWT token |
| `UpdateUser(user)` | Perbarui data pengguna |
| `DeleteUser(id)` | Hapus pengguna |

---

## 💡 Kenapa Interface di Domain?

Dengan mendefinisikan interface di lapisan Domain, kita menerapkan **Dependency Inversion Principle (DIP)**:

- Use case bergantung pada **interface** `UserRepository`, bukan implementasi konkret.
- Ini memungkinkan kita mengganti implementasi database (MySQL → PostgreSQL → MongoDB) **tanpa mengubah** kode use case atau domain.
- Memudahkan **unit testing** karena kita bisa membuat mock repository.
