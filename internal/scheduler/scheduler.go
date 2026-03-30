// scheduler.go berisi mesin utama scheduler yang mengelola pendaftaran,
// penjadwalan, dan eksekusi tugas-tugas otomatis.
//
// Scheduler ini terinspirasi dari Laravel Task Scheduler dan menyediakan
// API yang familiar seperti EveryMinute(), Hourly(), Daily(), Weekly(), dll.
//
// Contoh penggunaan:
//
//	s := scheduler.NewScheduler()
//
//	s.Register(&scheduler.Task{
//	    Name:     "cleanup-logs",
//	    Schedule: "@daily",
//	    Fn: func() error {
//	        log.Println("Membersihkan log lama...")
//	        return nil
//	    },
//	})
//
//	s.Start()
//	defer s.Stop()
package scheduler

import (
	"fmt"
	"log"

	// robfig/cron adalah library cron scheduler yang populer di Go.
	// Mendukung cron expression standar dan predefined schedule.
	"github.com/robfig/cron/v3"
)

// Scheduler adalah engine utama yang mengelola tugas-tugas terjadwal.
// Struct ini membungkus robfig/cron dan menyediakan API yang lebih ekspresif.
type Scheduler struct {
	// cron adalah instance cron scheduler dari library robfig/cron.
	cron *cron.Cron

	// tasks menyimpan semua tugas yang sudah terdaftar.
	tasks []*Task
}

// NewScheduler membuat instance Scheduler baru.
// Scheduler dibuat dengan opsi cron.WithSeconds() agar mendukung presisi detik.
func NewScheduler() *Scheduler {
	return &Scheduler{
		// cron.WithSeconds() mengaktifkan format cron 6-field (detik menit jam hari bulan hariMinggu).
		// Tanpa opsi ini, cron hanya mendukung format 5-field standar.
		cron: cron.New(cron.WithSeconds()),
		tasks: make([]*Task, 0),
	}
}

// Register mendaftarkan tugas baru ke scheduler menggunakan cron expression.
// Cron expression mendukung format 6-field (dengan detik) dan predefined schedule.
//
// Format cron 6-field:
//
//	┌──────── detik (0-59)
//	│ ┌────── menit (0-59)
//	│ │ ┌──── jam (0-23)
//	│ │ │ ┌── hari dalam bulan (1-31)
//	│ │ │ │ ┌ bulan (1-12)
//	│ │ │ │ │ ┌ hari dalam minggu (0-6, 0=Minggu)
//	│ │ │ │ │ │
//	* * * * * *
//
// Predefined schedule:
//
//	@yearly    → Sekali setahun (1 Jan 00:00)
//	@monthly   → Sekali sebulan (tanggal 1, 00:00)
//	@weekly    → Sekali seminggu (Minggu, 00:00)
//	@daily     → Sekali sehari (00:00)
//	@hourly    → Sekali sejam (menit ke-0)
//	@every 5m  → Setiap 5 menit
//	@every 30s → Setiap 30 detik
func (s *Scheduler) Register(task *Task) error {
	if task.Name == "" {
		return fmt.Errorf("nama task tidak boleh kosong")
	}
	if task.Fn == nil {
		return fmt.Errorf("fungsi task '%s' tidak boleh nil", task.Name)
	}
	if task.Schedule == "" {
		return fmt.Errorf("schedule task '%s' tidak boleh kosong", task.Name)
	}

	// Daftarkan task ke cron scheduler.
	// task mengimplementasikan cron.Job melalui method Run().
	_, err := s.cron.AddJob(task.Schedule, task)
	if err != nil {
		return fmt.Errorf("gagal mendaftarkan task '%s': %w", task.Name, err)
	}

	s.tasks = append(s.tasks, task)
	log.Printf("📋 [Scheduler] Task '%s' terdaftar dengan jadwal: %s", task.Name, task.Schedule)
	return nil
}

// ─── HELPER METHODS (Laravel-style API) ─────────────────────────────────────

// EverySecond mendaftarkan tugas yang berjalan setiap detik.
// Setara dengan Laravel: $schedule->call(...)->everySecond()
func (s *Scheduler) EverySecond(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "* * * * * *",
		Fn:       fn,
	})
}

// EveryFiveSeconds mendaftarkan tugas yang berjalan setiap 5 detik.
func (s *Scheduler) EveryFiveSeconds(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "*/5 * * * * *",
		Fn:       fn,
	})
}

// EveryTenSeconds mendaftarkan tugas yang berjalan setiap 10 detik.
func (s *Scheduler) EveryTenSeconds(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "*/10 * * * * *",
		Fn:       fn,
	})
}

// EveryFifteenSeconds mendaftarkan tugas yang berjalan setiap 15 detik.
func (s *Scheduler) EveryFifteenSeconds(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "*/15 * * * * *",
		Fn:       fn,
	})
}

// EveryThirtySeconds mendaftarkan tugas yang berjalan setiap 30 detik.
func (s *Scheduler) EveryThirtySeconds(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "*/30 * * * * *",
		Fn:       fn,
	})
}

// EveryMinute mendaftarkan tugas yang berjalan setiap menit.
// Setara dengan Laravel: $schedule->call(...)->everyMinute()
func (s *Scheduler) EveryMinute(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 * * * * *",
		Fn:       fn,
	})
}

// EveryTwoMinutes mendaftarkan tugas yang berjalan setiap 2 menit.
func (s *Scheduler) EveryTwoMinutes(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 */2 * * * *",
		Fn:       fn,
	})
}

// EveryFiveMinutes mendaftarkan tugas yang berjalan setiap 5 menit.
// Setara dengan Laravel: $schedule->call(...)->everyFiveMinutes()
func (s *Scheduler) EveryFiveMinutes(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 */5 * * * *",
		Fn:       fn,
	})
}

// EveryTenMinutes mendaftarkan tugas yang berjalan setiap 10 menit.
// Setara dengan Laravel: $schedule->call(...)->everyTenMinutes()
func (s *Scheduler) EveryTenMinutes(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 */10 * * * *",
		Fn:       fn,
	})
}

// EveryFifteenMinutes mendaftarkan tugas yang berjalan setiap 15 menit.
// Setara dengan Laravel: $schedule->call(...)->everyFifteenMinutes()
func (s *Scheduler) EveryFifteenMinutes(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 */15 * * * *",
		Fn:       fn,
	})
}

// EveryThirtyMinutes mendaftarkan tugas yang berjalan setiap 30 menit.
// Setara dengan Laravel: $schedule->call(...)->everyThirtyMinutes()
func (s *Scheduler) EveryThirtyMinutes(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 */30 * * * *",
		Fn:       fn,
	})
}

// Hourly mendaftarkan tugas yang berjalan setiap jam (pada menit ke-0).
// Setara dengan Laravel: $schedule->call(...)->hourly()
func (s *Scheduler) Hourly(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 0 * * * *",
		Fn:       fn,
	})
}

// Daily mendaftarkan tugas yang berjalan setiap hari pada pukul 00:00.
// Setara dengan Laravel: $schedule->call(...)->daily()
func (s *Scheduler) Daily(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 0 0 * * *",
		Fn:       fn,
	})
}

// DailyAt mendaftarkan tugas yang berjalan setiap hari pada jam tertentu.
// Parameter hour dalam format 24 jam (0-23), minute (0-59).
// Setara dengan Laravel: $schedule->call(...)->dailyAt('13:00')
func (s *Scheduler) DailyAt(name string, hour, minute int, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: fmt.Sprintf("0 %d %d * * *", minute, hour),
		Fn:       fn,
	})
}

// Weekly mendaftarkan tugas yang berjalan setiap hari Minggu pukul 00:00.
// Setara dengan Laravel: $schedule->call(...)->weekly()
func (s *Scheduler) Weekly(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 0 0 * * 0",
		Fn:       fn,
	})
}

// Monthly mendaftarkan tugas yang berjalan setiap tanggal 1 pukul 00:00.
// Setara dengan Laravel: $schedule->call(...)->monthly()
func (s *Scheduler) Monthly(name string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: "0 0 0 1 * *",
		Fn:       fn,
	})
}

// Cron mendaftarkan tugas dengan cron expression kustom (6-field dengan detik).
// Ini adalah method yang paling fleksibel untuk jadwal yang tidak standar.
// Setara dengan Laravel: $schedule->call(...)->cron('* * * * *')
func (s *Scheduler) Cron(name, expression string, fn TaskFunc) error {
	return s.Register(&Task{
		Name:     name,
		Schedule: expression,
		Fn:       fn,
	})
}

// Start memulai scheduler dan mulai menjalankan tugas sesuai jadwal.
// Method ini non-blocking (berjalan di goroutine terpisah).
func (s *Scheduler) Start() {
	log.Printf("🚀 [Scheduler] Dimulai dengan %d task terdaftar", len(s.tasks))
	s.cron.Start()
}

// Stop menghentikan scheduler secara graceful.
// Menunggu semua tugas yang sedang berjalan selesai sebelum berhenti.
func (s *Scheduler) Stop() {
	log.Println("🛑 [Scheduler] Menghentikan scheduler...")
	ctx := s.cron.Stop()
	// Tunggu semua tugas yang sedang berjalan selesai.
	<-ctx.Done()
	log.Println("✅ [Scheduler] Scheduler berhenti dengan sukses")
}

// GetTasks mengembalikan daftar semua tugas yang terdaftar.
// Berguna untuk monitoring atau endpoint status.
func (s *Scheduler) GetTasks() []*Task {
	return s.tasks
}
