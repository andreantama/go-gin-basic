# 📁 internal/middleware

Package `middleware` berisi fungsi-fungsi middleware HTTP untuk GIN. Middleware dieksekusi **sebelum** handler utama dan dapat menghentikan atau melanjutkan proses request.

---

## 🎯 Tanggung Jawab

- **Autentikasi** — Validasi JWT token di header Authorization.
- **Otorisasi** — Periksa role pengguna (admin/user).
- **CORS** — Konfigurasi header untuk cross-origin request.

---

## 📄 Struktur File

```
internal/middleware/
└── auth_middleware.go   # JWT auth, admin-only, dan CORS middleware
```

---

## 🔐 AuthMiddleware

Memvalidasi JWT token dari header `Authorization: Bearer <token>`.

**Alur Validasi:**
```
Request masuk
    ↓
Periksa header Authorization
    ↓ (tidak ada)         ↓ (ada)
401 Unauthorized    Parse & validasi token JWT
                        ↓ (invalid)     ↓ (valid)
                    401 Unauthorized  Set userID, email, role ke context
                                          ↓
                                    Lanjut ke handler
```

**Data yang disimpan di GIN Context:**

| Key | Tipe | Deskripsi |
|-----|------|-----------|
| `userID` | `uint` | ID pengguna dari claim `sub` |
| `email` | `string` | Email pengguna |
| `role` | `string` | Role pengguna (`user` atau `admin`) |

**Cara membaca dari handler:**
```go
userID, _ := c.Get("userID")
email, _ := c.Get("email")
role, _ := c.Get("role")
```

---

## 👑 AdminOnly

Middleware otorisasi yang hanya mengizinkan pengguna dengan role `admin`.

> ⚠️ Harus digunakan **setelah** `AuthMiddleware`.

```go
adminRoutes := router.Group("/admin")
adminRoutes.Use(middleware.AuthMiddleware(cfg), middleware.AdminOnly())
```

---

## 🌍 CORSMiddleware

Mengkonfigurasi header CORS untuk mengizinkan request dari browser.

> ⚠️ Di production, ganti `Access-Control-Allow-Origin: *` dengan domain spesifik.

---

## 🚀 Cara Penggunaan

```go
// Di main.go
router := gin.Default()

// Tambah CORS middleware ke semua route
router.Use(middleware.CORSMiddleware())

// Grup route yang dilindungi JWT
protected := router.Group("/api")
protected.Use(middleware.AuthMiddleware(cfg))
{
    protected.GET("/profile", handler.GetProfile)
}

// Grup route khusus admin
admin := router.Group("/admin")
admin.Use(middleware.AuthMiddleware(cfg), middleware.AdminOnly())
{
    admin.DELETE("/users/:id", handler.DeleteUser)
}
```

---

## 🛡️ Keamanan

- **Algorithm Confusion Attack Prevention**: Middleware memvalidasi bahwa algoritma yang digunakan adalah HMAC (HS256), mencegah attacker mengganti ke algoritma `"none"`.
- **Generic Error Messages**: Pesan error tidak membedakan antara "token tidak ada" dan "token invalid" untuk mencegah information leakage.
- **Context Isolation**: Data pengguna disimpan di GIN context, bukan variabel global, sehingga aman untuk concurrent request.
