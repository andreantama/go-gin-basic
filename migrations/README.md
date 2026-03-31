# 🗄️ Database Migrations

Folder ini berisi file-file SQL migration yang digunakan oleh **golang-migrate** untuk mengelola skema database secara terstruktur dan terkontrol.

---

## 📋 Daftar Isi

- [Mengapa golang-migrate?](#mengapa-golang-migrate)
- [Struktur Folder](#struktur-folder)
- [Konvensi Penamaan File](#konvensi-penamaan-file)
- [Cara Kerja](#cara-kerja)
- [Menambahkan Migration Baru](#menambahkan-migration-baru)
- [Cara Rollback](#cara-rollback)
- [Tabel schema_migrations](#tabel-schema_migrations)
- [Perbedaan GORM AutoMigrate vs golang-migrate](#perbedaan-gorm-automigrate-vs-golang-migrate)

---

## 🤔 Mengapa golang-migrate?

Proyek ini sebelumnya menggunakan **GORM AutoMigrate** yang berjalan otomatis saat aplikasi start. Fitur tersebut digantikan dengan **golang-migrate** karena:

| Fitur | GORM AutoMigrate | golang-migrate |
|-------|-----------------|----------------|
| File SQL eksplisit | ❌ (berbasis struct) | ✅ (file .sql terpisah) |
| Rollback (DOWN) | ❌ | ✅ |
| Audit trail perubahan skema | ❌ | ✅ (tabel `schema_migrations`) |
| Review SQL sebelum dieksekusi | ❌ | ✅ |
| Idempotent (hanya jalankan sekali) | ⚠️ (partial) | ✅ |
| Cocok untuk production | ⚠️ | ✅ |

---

## 📁 Struktur Folder

```
migrations/
├── embed.go                                    # Mengekspor embed.FS ke package lain
│
├── mysql/                                      # Migration khusus MySQL
│   ├── 000001_create_users_table.up.sql        # Membuat tabel users
│   └── 000001_create_users_table.down.sql      # Menghapus tabel users
│
└── postgres/                                   # Migration khusus PostgreSQL
    ├── 000001_create_users_table.up.sql        # Membuat tabel users
    └── 000001_create_users_table.down.sql      # Menghapus tabel users
```

> **Catatan:** File-file SQL di-embed ke dalam binary menggunakan Go `embed` directive, sehingga tidak perlu mendistribusikan file SQL secara terpisah saat deployment.

---

## 📝 Konvensi Penamaan File

Format nama file migration adalah:

```
{versi}_{deskripsi}.{arah}.sql
```

| Bagian | Keterangan | Contoh |
|--------|------------|--------|
| `{versi}` | Nomor urut 6 digit, bertambah secara sequential | `000001`, `000002` |
| `{deskripsi}` | Deskripsi singkat perubahan menggunakan snake_case | `create_users_table`, `add_index_to_orders` |
| `{arah}` | `up` untuk menerapkan, `down` untuk rollback | `up`, `down` |

### Contoh nama file yang benar:

```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_phone_to_users.up.sql
000002_add_phone_to_users.down.sql
000003_create_orders_table.up.sql
000003_create_orders_table.down.sql
```

---

## ⚙️ Cara Kerja

1. **Saat aplikasi start**, fungsi `RunMigrations()` di `infrastructure/database/migrate.go` dipanggil otomatis.
2. golang-migrate **memeriksa tabel `schema_migrations`** di database untuk melihat migration mana yang sudah dieksekusi.
3. Hanya migration yang **belum dieksekusi** yang akan dijalankan (idempotent).
4. Setiap migration yang berhasil dicatat di tabel `schema_migrations` dengan nomor versinya.

```
Aplikasi Start
     │
     ▼
RunMigrations()
     │
     ▼
Cek tabel schema_migrations
     │
     ├─── Sudah ada versi 000001? ──► Skip
     │
     └─── Belum ada? ──────────────► Jalankan 000001_*.up.sql
                                              │
                                              ▼
                                     Catat di schema_migrations
```

---

## ➕ Menambahkan Migration Baru

Untuk menambahkan migrasi baru (misalnya menambahkan kolom `phone` ke tabel `users`):

### 1. Buat file SQL untuk MySQL

**`migrations/mysql/000002_add_phone_to_users.up.sql`**
```sql
ALTER TABLE `users`
    ADD COLUMN `phone` VARCHAR(20) DEFAULT NULL AFTER `role`;
```

**`migrations/mysql/000002_add_phone_to_users.down.sql`**
```sql
ALTER TABLE `users`
    DROP COLUMN `phone`;
```

### 2. Buat file SQL untuk PostgreSQL

**`migrations/postgres/000002_add_phone_to_users.up.sql`**
```sql
ALTER TABLE users
    ADD COLUMN phone VARCHAR(20) DEFAULT NULL;
```

**`migrations/postgres/000002_add_phone_to_users.down.sql`**
```sql
ALTER TABLE users
    DROP COLUMN IF EXISTS phone;
```

### 3. Update domain struct (opsional, jika menggunakan GORM untuk query)

```go
// internal/domain/user.go
type User struct {
    // ... field yang sudah ada ...
    Phone string `json:"phone,omitempty" gorm:"default:null"`
}
```

### 4. Jalankan aplikasi

```bash
go run main.go
```

Migration baru akan otomatis dieksekusi saat aplikasi start.

---

## ↩️ Cara Rollback

golang-migrate mendukung rollback migration menggunakan file `.down.sql`. Untuk saat ini, rollback hanya dapat dilakukan secara programatik. Berikut contoh membuat CLI tool untuk rollback:

```go
// Contoh: rollback 1 migration terakhir
m.Steps(-1)

// Contoh: rollback ke versi tertentu
m.Migrate(1) // rollback ke versi 000001
```

> Untuk keperluan development, Anda juga bisa menghapus record dari tabel `schema_migrations` dan menjalankan ulang aplikasi.

---

## 🗃️ Tabel schema_migrations

golang-migrate secara otomatis membuat tabel `schema_migrations` di database Anda:

```sql
-- MySQL
CREATE TABLE IF NOT EXISTS `schema_migrations` (
    `version`  bigint      NOT NULL,
    `dirty`    tinyint(1)  NOT NULL,
    PRIMARY KEY (`version`)
);

-- PostgreSQL
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint  NOT NULL,
    dirty   boolean NOT NULL,
    PRIMARY KEY (version)
);
```

| Kolom | Keterangan |
|-------|------------|
| `version` | Nomor versi migration terakhir yang dieksekusi |
| `dirty` | `true` jika migration terakhir gagal di tengah jalan |

> **Penting:** Jangan pernah menghapus tabel `schema_migrations` di production karena golang-migrate menggunakannya untuk melacak status migration.

---

## 🔄 Perbedaan GORM AutoMigrate vs golang-migrate

### GORM AutoMigrate (sebelumnya)

```go
// Cara lama — otomatis, tapi tidak terkontrol
db.AutoMigrate(&domain.User{})
```

**Kekurangan:**
- Tidak ada file SQL yang bisa di-review
- Tidak mendukung rollback
- Perilaku tidak terprediksi di production (menambah kolom tanpa konfirmasi)
- Tidak ada audit trail

### golang-migrate (sekarang)

```sql
-- 000001_create_users_table.up.sql — eksplisit dan bisa di-review
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL,
    ...
);
```

**Keunggulan:**
- ✅ File SQL eksplisit dan dapat di-review sebelum dieksekusi
- ✅ Mendukung rollback dengan file `.down.sql`
- ✅ Setiap migration hanya dieksekusi sekali
- ✅ Audit trail di tabel `schema_migrations`
- ✅ Cocok untuk tim dengan code review workflow

---

## 📖 Referensi

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [golang-migrate source/iofs](https://pkg.go.dev/github.com/golang-migrate/migrate/v4/source/iofs)
- [Go embed Package](https://pkg.go.dev/embed)
