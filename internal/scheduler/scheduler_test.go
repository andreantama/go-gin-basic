package scheduler

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewScheduler memastikan scheduler dibuat dengan benar.
func TestNewScheduler(t *testing.T) {
	s := NewScheduler()
	if s == nil {
		t.Fatal("NewScheduler() mengembalikan nil")
	}
	if s.cron == nil {
		t.Fatal("cron instance tidak diinisialisasi")
	}
	if len(s.tasks) != 0 {
		t.Fatalf("tasks seharusnya kosong, mendapat %d", len(s.tasks))
	}
}

// TestRegisterValidTask memastikan task valid bisa didaftarkan.
func TestRegisterValidTask(t *testing.T) {
	s := NewScheduler()

	err := s.Register(&Task{
		Name:     "test-task",
		Schedule: "@every 1s",
		Fn: func() error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Register() gagal: %v", err)
	}

	if len(s.GetTasks()) != 1 {
		t.Fatalf("seharusnya ada 1 task, mendapat %d", len(s.GetTasks()))
	}
	if s.GetTasks()[0].Name != "test-task" {
		t.Fatalf("nama task seharusnya 'test-task', mendapat '%s'", s.GetTasks()[0].Name)
	}
}

// TestRegisterInvalidTask memastikan validasi task berjalan.
func TestRegisterInvalidTask(t *testing.T) {
	s := NewScheduler()

	// Task tanpa nama.
	err := s.Register(&Task{
		Schedule: "@every 1s",
		Fn:       func() error { return nil },
	})
	if err == nil {
		t.Fatal("seharusnya error untuk task tanpa nama")
	}

	// Task tanpa fungsi.
	err = s.Register(&Task{
		Name:     "no-fn-task",
		Schedule: "@every 1s",
	})
	if err == nil {
		t.Fatal("seharusnya error untuk task tanpa fungsi")
	}

	// Task tanpa schedule.
	err = s.Register(&Task{
		Name: "no-schedule-task",
		Fn:   func() error { return nil },
	})
	if err == nil {
		t.Fatal("seharusnya error untuk task tanpa schedule")
	}

	// Task dengan cron expression tidak valid.
	err = s.Register(&Task{
		Name:     "invalid-cron-task",
		Schedule: "invalid expression",
		Fn:       func() error { return nil },
	})
	if err == nil {
		t.Fatal("seharusnya error untuk cron expression tidak valid")
	}
}

// TestTaskExecution memastikan task dijalankan sesuai jadwal.
func TestTaskExecution(t *testing.T) {
	s := NewScheduler()

	var counter int64
	err := s.Register(&Task{
		Name:     "counter-task",
		Schedule: "@every 1s",
		Fn: func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Register() gagal: %v", err)
	}

	s.Start()
	// Tunggu cukup lama agar task berjalan minimal 1 kali.
	time.Sleep(2500 * time.Millisecond)
	s.Stop()

	count := atomic.LoadInt64(&counter)
	if count < 1 {
		t.Fatalf("task seharusnya berjalan minimal 1 kali, berjalan %d kali", count)
	}
}

// TestTaskPanicRecovery memastikan scheduler tidak crash ketika task panic.
func TestTaskPanicRecovery(t *testing.T) {
	task := &Task{
		Name:     "panic-task",
		Schedule: "@every 1s",
		Fn: func() error {
			panic("test panic")
		},
	}

	// Jalankan task secara langsung — tidak boleh panic.
	task.Run()

	if task.LastError() == nil {
		t.Fatal("seharusnya ada error setelah panic")
	}
	if task.LastError().Error() != "panic: test panic" {
		t.Fatalf("error seharusnya 'panic: test panic', mendapat '%s'", task.LastError().Error())
	}
}

// TestTaskErrorHandling memastikan error task dicatat dengan benar.
func TestTaskErrorHandling(t *testing.T) {
	expectedErr := errors.New("test error")
	task := &Task{
		Name:     "error-task",
		Schedule: "@every 1s",
		Fn: func() error {
			return expectedErr
		},
	}

	task.Run()

	if task.LastError() == nil {
		t.Fatal("seharusnya ada error")
	}
	if task.LastError().Error() != expectedErr.Error() {
		t.Fatalf("error seharusnya '%v', mendapat '%v'", expectedErr, task.LastError())
	}
	if task.LastRunAt() == nil {
		t.Fatal("lastRunAt seharusnya ter-set setelah eksekusi")
	}
}

// TestWithoutOverlapping memastikan overlap prevention berfungsi.
func TestWithoutOverlapping(t *testing.T) {
	var running int64
	var maxRunning int64

	task := &Task{
		Name:               "overlap-task",
		Schedule:            "@every 1s",
		WithoutOverlapping: true,
		Fn: func() error {
			current := atomic.AddInt64(&running, 1)
			// Catat jumlah eksekusi bersamaan tertinggi.
			for {
				old := atomic.LoadInt64(&maxRunning)
				if current <= old || atomic.CompareAndSwapInt64(&maxRunning, old, current) {
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt64(&running, -1)
			return nil
		},
	}

	// Jalankan beberapa eksekusi secara bersamaan.
	done := make(chan struct{}, 5)
	for i := 0; i < 5; i++ {
		go func() {
			task.Run()
			done <- struct{}{}
		}()
	}

	// Tunggu semua selesai.
	for i := 0; i < 5; i++ {
		<-done
	}

	if atomic.LoadInt64(&maxRunning) > 1 {
		t.Fatalf("seharusnya tidak ada overlap, maxRunning = %d", atomic.LoadInt64(&maxRunning))
	}
}

// TestHelperMethods memastikan semua helper method mendaftarkan task dengan benar.
func TestHelperMethods(t *testing.T) {
	fn := func() error { return nil }

	tests := []struct {
		name   string
		register func(s *Scheduler) error
	}{
		{"EverySecond", func(s *Scheduler) error { return s.EverySecond("t1", fn) }},
		{"EveryFiveSeconds", func(s *Scheduler) error { return s.EveryFiveSeconds("t2", fn) }},
		{"EveryTenSeconds", func(s *Scheduler) error { return s.EveryTenSeconds("t3", fn) }},
		{"EveryFifteenSeconds", func(s *Scheduler) error { return s.EveryFifteenSeconds("t4", fn) }},
		{"EveryThirtySeconds", func(s *Scheduler) error { return s.EveryThirtySeconds("t5", fn) }},
		{"EveryMinute", func(s *Scheduler) error { return s.EveryMinute("t6", fn) }},
		{"EveryTwoMinutes", func(s *Scheduler) error { return s.EveryTwoMinutes("t7", fn) }},
		{"EveryFiveMinutes", func(s *Scheduler) error { return s.EveryFiveMinutes("t8", fn) }},
		{"EveryTenMinutes", func(s *Scheduler) error { return s.EveryTenMinutes("t9", fn) }},
		{"EveryFifteenMinutes", func(s *Scheduler) error { return s.EveryFifteenMinutes("t10", fn) }},
		{"EveryThirtyMinutes", func(s *Scheduler) error { return s.EveryThirtyMinutes("t11", fn) }},
		{"Hourly", func(s *Scheduler) error { return s.Hourly("t12", fn) }},
		{"Daily", func(s *Scheduler) error { return s.Daily("t13", fn) }},
		{"DailyAt", func(s *Scheduler) error { return s.DailyAt("t14", 9, 30, fn) }},
		{"Weekly", func(s *Scheduler) error { return s.Weekly("t15", fn) }},
		{"Monthly", func(s *Scheduler) error { return s.Monthly("t16", fn) }},
		{"Cron", func(s *Scheduler) error { return s.Cron("t17", "0 0 * * * *", fn) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScheduler()
			err := tt.register(s)
			if err != nil {
				t.Fatalf("%s gagal: %v", tt.name, err)
			}
			if len(s.GetTasks()) != 1 {
				t.Fatalf("%s: seharusnya ada 1 task, mendapat %d", tt.name, len(s.GetTasks()))
			}
		})
	}
}

// TestStartStop memastikan scheduler bisa start dan stop tanpa error.
func TestStartStop(t *testing.T) {
	s := NewScheduler()

	err := s.EveryMinute("test-start-stop", func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("Register gagal: %v", err)
	}

	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()
}

// TestTaskLastRunAt memastikan timestamp eksekusi ter-set.
func TestTaskLastRunAt(t *testing.T) {
	task := &Task{
		Name:     "timestamp-task",
		Schedule: "@every 1s",
		Fn: func() error {
			return nil
		},
	}

	// Sebelum dijalankan, LastRunAt harus nil.
	if task.LastRunAt() != nil {
		t.Fatal("LastRunAt seharusnya nil sebelum eksekusi")
	}

	// Sebelum dijalankan, LastError harus nil.
	if task.LastError() != nil {
		t.Fatal("LastError seharusnya nil sebelum eksekusi")
	}

	before := time.Now()
	task.Run()
	after := time.Now()

	if task.LastRunAt() == nil {
		t.Fatal("LastRunAt seharusnya ter-set setelah eksekusi")
	}
	if task.LastRunAt().Before(before) || task.LastRunAt().After(after) {
		t.Fatalf("LastRunAt (%v) seharusnya antara %v dan %v", task.LastRunAt(), before, after)
	}
}
