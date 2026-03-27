# 📁 pkg

Package `pkg` berisi **utility packages** yang bersifat umum dan bisa digunakan di seluruh aplikasi, bahkan di proyek lain.

---

## 🎯 Perbedaan `pkg` vs `internal`

| Aspek | `pkg/` | `internal/` |
|-------|--------|-------------|
| **Akses** | Bisa diimpor oleh proyek lain | Hanya untuk proyek ini |
| **Konten** | Utility umum yang reusable | Logika spesifik aplikasi ini |
| **Contoh** | response formatter, error types | domain entities, use cases |

---

## 📄 Struktur File

```
pkg/
├── response/
│   └── response.go   # Format respons JSON yang konsisten
└── errors/
    └── errors.go     # Tipe error kustom dan HTTP status mapper
```

---

## 📦 Package `response`

Menyediakan format respons JSON yang **seragam** untuk seluruh endpoint API.

```go
// Respons sukses
c.JSON(http.StatusOK, response.Success("Data diambil", data))
// → {"success": true, "message": "Data diambil", "data": {...}}

// Respons error
c.JSON(http.StatusBadRequest, response.Error("Input tidak valid"))
// → {"success": false, "message": "Input tidak valid"}
```

---

## 📦 Package `errors`

Menyediakan tipe error kustom yang bisa di-map ke HTTP status code yang tepat.

```go
// Di use case — kembalikan error dengan tipe yang tepat
return nil, errors.NewNotFoundError("pengguna tidak ditemukan")
return nil, errors.NewConflictError("email sudah terdaftar")

// Di handler — konversi ke HTTP status code
statusCode := pkgerrors.HTTPStatus(err)
c.JSON(statusCode, response.Error(err.Error()))
```

**Mapping Error → HTTP Status:**

| Error Type | HTTP Status |
|------------|-------------|
| `NotFoundError` | 404 Not Found |
| `ValidationError` | 400 Bad Request |
| `ConflictError` | 409 Conflict |
| `UnauthorizedError` | 401 Unauthorized |
| Error lainnya | 500 Internal Server Error |

---

## 🏛️ Posisi dalam Clean Architecture

Package `pkg` dapat digunakan oleh lapisan manapun karena bersifat umum dan tidak mengandung business logic atau framework dependency.
