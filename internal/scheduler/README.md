# 📁 internal/scheduler

Package `scheduler` menyediakan fitur **penjadwalan tugas otomatis** (task scheduling) yang terinspirasi dari **Laravel Task Scheduler**. Package ini memungkinkan Anda mendaftarkan tugas-tugas yang akan dijalankan secara otomatis berdasarkan jadwal tertentu.

---

## 🎯 Tanggung Jawab

- Menjalankan **tugas-tugas otomatis** sesuai jadwal (cron expression atau interval).
- Menyediakan **API ekspresif** mirip Laravel (`EveryMinute()`, `Hourly()`, `Daily()`, dll).
- **Mencegah overlap** eksekusi tugas yang sama (opsional).
- **Recovery dari panic** — satu tugas yang gagal tidak menghentikan scheduler.
- **Graceful shutdown** — menunggu tugas yang sedang berjalan selesai sebelum berhenti.

---

## 📄 Struktur File

```
internal/scheduler/
├── scheduler.go   # Engine utama scheduler (registrasi, start, stop)
├── task.go        # Definisi Task (nama, jadwal, fungsi, metadata eksekusi)
└── README.md      # Dokumentasi ini
```

---

## 🏛️ Posisi dalam Clean Architecture

```
┌─────────────────────────────────────┐
│         Frameworks & Drivers        │
├─────────────────────────────────────┤
│         Interface Adapters          │
├─────────────────────────────────────┤
│  ► Application / Use Cases ◄        │  ← Scheduler ada di sini
├─────────────────────────────────────┤
│         Domain / Entities           │
└─────────────────────────────────────┘
```

Scheduler berada di lapisan **Application** karena ia mengorkestrasi eksekusi tugas-tugas bisnis tanpa bergantung pada framework HTTP atau database tertentu.

---

## 🚀 Cara Penggunaan

### 1. Buat Instance Scheduler

```go
s := scheduler.NewScheduler()
```

### 2. Daftarkan Tugas

**Menggunakan helper method (Laravel-style):**

```go
// Jalankan setiap menit
s.EveryMinute("kirim-notifikasi", func() error {
    log.Println("Mengirim notifikasi...")
    return nil
})

// Jalankan setiap hari pukul 02:00
s.DailyAt("cleanup-database", 2, 0, func() error {
    log.Println("Membersihkan data lama...")
    return nil
})

// Jalankan setiap jam
s.Hourly("sync-data", func() error {
    log.Println("Sinkronisasi data...")
    return nil
})
```

**Menggunakan cron expression kustom:**

```go
s.Cron("backup-database", "0 0 2 * * *", func() error {
    log.Println("Backup database...")
    return nil
})
```

**Menggunakan Register() langsung (dengan opsi lanjutan):**

```go
s.Register(&scheduler.Task{
    Name:               "proses-antrian",
    Schedule:           "*/30 * * * * *",  // Setiap 30 detik
    WithoutOverlapping: true,               // Cegah overlap
    Fn: func() error {
        // Proses antrian yang mungkin memakan waktu lama
        return nil
    },
})
```

### 3. Jalankan Scheduler

```go
s.Start()        // Non-blocking, berjalan di goroutine terpisah
defer s.Stop()   // Hentikan saat aplikasi shutdown
```

---

## 📋 Helper Methods (Laravel-style API)

| Method | Jadwal | Setara Laravel |
|--------|--------|----------------|
| `EverySecond()` | Setiap detik | `everySecond()` |
| `EveryFiveSeconds()` | Setiap 5 detik | `everyFiveSeconds()` |
| `EveryTenSeconds()` | Setiap 10 detik | `everyTenSeconds()` |
| `EveryFifteenSeconds()` | Setiap 15 detik | `everyFifteenSeconds()` |
| `EveryThirtySeconds()` | Setiap 30 detik | `everyThirtySeconds()` |
| `EveryMinute()` | Setiap menit | `everyMinute()` |
| `EveryTwoMinutes()` | Setiap 2 menit | `everyTwoMinutes()` |
| `EveryFiveMinutes()` | Setiap 5 menit | `everyFiveMinutes()` |
| `EveryTenMinutes()` | Setiap 10 menit | `everyTenMinutes()` |
| `EveryFifteenMinutes()` | Setiap 15 menit | `everyFifteenMinutes()` |
| `EveryThirtyMinutes()` | Setiap 30 menit | `everyThirtyMinutes()` |
| `Hourly()` | Setiap jam | `hourly()` |
| `Daily()` | Setiap hari (00:00) | `daily()` |
| `DailyAt(hour, minute)` | Setiap hari pada jam tertentu | `dailyAt('13:00')` |
| `Weekly()` | Setiap Minggu (00:00) | `weekly()` |
| `Monthly()` | Setiap tanggal 1 (00:00) | `monthly()` |
| `Cron(expression)` | Custom cron expression | `cron('* * * * *')` |

---

## 🔄 Format Cron Expression

Scheduler ini menggunakan format **6-field** (dengan detik):

```
┌──────── detik (0-59)
│ ┌────── menit (0-59)
│ │ ┌──── jam (0-23)
│ │ │ ┌── hari dalam bulan (1-31)
│ │ │ │ ┌ bulan (1-12)
│ │ │ │ │ ┌ hari dalam minggu (0-6, 0=Minggu)
│ │ │ │ │ │
* * * * * *
```

**Predefined schedule** juga didukung:

| Ekspresi | Deskripsi |
|----------|-----------|
| `@yearly` | Sekali setahun (1 Jan, 00:00:00) |
| `@monthly` | Sekali sebulan (tanggal 1, 00:00:00) |
| `@weekly` | Sekali seminggu (Minggu, 00:00:00) |
| `@daily` | Sekali sehari (00:00:00) |
| `@hourly` | Sekali sejam (menit ke-0) |
| `@every <durasi>` | Setiap durasi tertentu (contoh: `@every 5m`, `@every 30s`) |

---

## 🛡️ Fitur Keamanan

### Pencegahan Overlap

Aktifkan `WithoutOverlapping: true` untuk mencegah tugas yang sama berjalan bersamaan:

```go
s.Register(&scheduler.Task{
    Name:               "long-running-task",
    Schedule:           "0 * * * * *",    // Setiap menit
    WithoutOverlapping: true,              // Jika masih berjalan, lewati
    Fn: func() error {
        // Proses yang mungkin memakan waktu > 1 menit
        time.Sleep(2 * time.Minute)
        return nil
    },
})
```

### Recovery dari Panic

Jika sebuah tugas mengalami panic, scheduler akan menangkapnya dan melanjutkan tugas-tugas lain:

```go
s.EveryMinute("risky-task", func() error {
    // Jika panic terjadi di sini, scheduler tetap berjalan
    panic("something went wrong")
})
```

---

## 📊 Monitoring

### Mendapatkan Daftar Tugas

```go
tasks := s.GetTasks()
for _, task := range tasks {
    fmt.Printf("Task: %s\n", task.Name)
    if task.LastRunAt() != nil {
        fmt.Printf("  Terakhir dijalankan: %v\n", task.LastRunAt())
    }
    if task.LastError() != nil {
        fmt.Printf("  Error terakhir: %v\n", task.LastError())
    }
}
```

---

## 🔗 Integrasi dengan main.go

Scheduler diintegrasikan di `main.go` dan dijalankan sebagai goroutine terpisah bersama server HTTP:

```go
func main() {
    // ... inisialisasi config, database, router ...

    // Buat dan konfigurasi scheduler
    taskScheduler := scheduler.NewScheduler()

    // Daftarkan tugas-tugas
    taskScheduler.EveryMinute("health-check", func() error {
        log.Println("Server berjalan normal")
        return nil
    })

    // Jalankan scheduler
    taskScheduler.Start()
    defer taskScheduler.Stop()

    // Jalankan server HTTP
    router.Run(":8080")
}
```

---

## 💡 Perbandingan dengan Laravel

| Fitur | Laravel | Go Scheduler |
|-------|---------|-------------|
| Definisi jadwal | `$schedule->call(...)` | `scheduler.Register(...)` |
| Helper method | `->everyMinute()` | `.EveryMinute()` |
| Cron expression | `->cron('...')` | `.Cron('...')` |
| Cegah overlap | `->withoutOverlapping()` | `WithoutOverlapping: true` |
| Recovery dari error | Built-in | Built-in (panic recovery) |
| Logging | Built-in | Built-in (log package) |
| Graceful shutdown | Signal handling | `.Stop()` method |

---

## 🧪 Testability

Scheduler mudah di-test karena `TaskFunc` adalah fungsi sederhana:

```go
// Test fungsi tugas secara langsung
fn := func() error {
    // logika tugas
    return nil
}

err := fn()
assert.NoError(t, err)
```
