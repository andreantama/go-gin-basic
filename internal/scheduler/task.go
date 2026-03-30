// Package scheduler menyediakan fitur penjadwalan tugas otomatis (task scheduling)
// yang terinspirasi dari Laravel Task Scheduler.
//
// Package ini memungkinkan Anda mendefinisikan tugas-tugas yang akan dijalankan
// secara otomatis berdasarkan jadwal tertentu (cron expression, interval, dll).
package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// TaskFunc adalah tipe fungsi yang dijalankan oleh scheduler.
// Setiap tugas terjadwal harus berupa fungsi dengan signature ini.
type TaskFunc func() error

// Task merepresentasikan satu tugas terjadwal.
// Setiap task memiliki nama, jadwal, fungsi yang dijalankan, dan metadata eksekusi.
type Task struct {
	// Name adalah nama unik untuk mengidentifikasi tugas ini dalam log.
	Name string

	// Schedule adalah cron expression yang menentukan kapan tugas dijalankan.
	// Format: "detik menit jam hari-bulan bulan hari-minggu"
	// Contoh: "0 * * * *" berarti setiap jam pada menit ke-0.
	Schedule string

	// Fn adalah fungsi yang akan dijalankan sesuai jadwal.
	Fn TaskFunc

	// WithoutOverlapping jika true, mencegah tugas berjalan bersamaan
	// jika eksekusi sebelumnya belum selesai.
	WithoutOverlapping bool

	// running menandakan apakah tugas sedang berjalan (untuk overlap prevention).
	running bool

	// mu adalah mutex untuk melindungi field running dari race condition.
	mu sync.Mutex

	// lastRunAt mencatat waktu terakhir tugas dijalankan.
	lastRunAt *time.Time

	// lastError mencatat error terakhir yang terjadi saat tugas dijalankan.
	lastError error
}

// Run menjalankan fungsi tugas dengan proteksi overlap dan pencatatan error.
// Method ini dipanggil secara internal oleh scheduler sesuai jadwal.
func (t *Task) Run() {
	// Jika WithoutOverlapping aktif, cek apakah tugas masih berjalan.
	if t.WithoutOverlapping {
		t.mu.Lock()
		if t.running {
			t.mu.Unlock()
			log.Printf("⏭️  [Scheduler] Task '%s' dilewati — masih berjalan dari eksekusi sebelumnya", t.Name)
			return
		}
		t.running = true
		t.mu.Unlock()

		// Pastikan flag running di-reset setelah tugas selesai.
		defer func() {
			t.mu.Lock()
			t.running = false
			t.mu.Unlock()
		}()
	}

	// Catat waktu mulai eksekusi.
	startTime := time.Now()
	log.Printf("▶️  [Scheduler] Menjalankan task '%s'...", t.Name)

	// Jalankan fungsi tugas dan tangkap panic jika terjadi.
	err := t.safeRun()

	// Catat waktu eksekusi dan hasilnya.
	duration := time.Since(startTime)
	now := time.Now()
	t.lastRunAt = &now
	t.lastError = err

	if err != nil {
		log.Printf("❌ [Scheduler] Task '%s' gagal setelah %v: %v", t.Name, duration, err)
	} else {
		log.Printf("✅ [Scheduler] Task '%s' selesai dalam %v", t.Name, duration)
	}
}

// safeRun menjalankan fungsi tugas dengan recovery dari panic.
// Jika fungsi panic, error akan dikembalikan sebagai nilai return.
func (t *Task) safeRun() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return t.Fn()
}

// LastRunAt mengembalikan waktu terakhir tugas dijalankan.
// Mengembalikan nil jika tugas belum pernah dijalankan.
func (t *Task) LastRunAt() *time.Time {
	return t.lastRunAt
}

// LastError mengembalikan error terakhir yang terjadi saat tugas dijalankan.
// Mengembalikan nil jika eksekusi terakhir berhasil atau belum pernah dijalankan.
func (t *Task) LastError() error {
	return t.lastError
}
