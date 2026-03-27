// Package domain berisi definisi entitas bisnis inti (core business entities).
//
// Dalam Clean Architecture, lapisan Domain (juga disebut "Entities") adalah lapisan
// paling dalam yang TIDAK boleh bergantung pada lapisan lain.
// Domain hanya berisi:
//   - Struct entitas (representasi objek bisnis)
//   - Interface repository (kontrak data access, diimplementasikan di lapisan luar)
//   - Domain errors (error yang spesifik untuk aturan bisnis)
//
// Lapisan ini murni berisi logika bisnis dan tidak mengenal framework apapun.
package domain

import (
	// time digunakan untuk tipe data waktu pada field CreatedAt, UpdatedAt, DeletedAt.
	"time"
)

// User adalah entitas utama yang merepresentasikan pengguna dalam sistem.
// Struct ini adalah "pusat kebenaran" (source of truth) untuk data pengguna.
// GORM menggunakan field-field ini untuk memetakan ke tabel database.
type User struct {
	// ID adalah primary key unik untuk setiap pengguna.
	// `gorm:"primaryKey"` menandai field ini sebagai primary key di database.
	ID uint `json:"id" gorm:"primaryKey"`

	// Name adalah nama lengkap pengguna.
	// `not null` memastikan field ini wajib diisi di database.
	Name string `json:"name" gorm:"not null"`

	// Email adalah alamat email pengguna, harus unik di seluruh sistem.
	// `uniqueIndex` membuat index unik di database untuk mencegah duplikasi email.
	Email string `json:"email" gorm:"uniqueIndex;not null"`

	// Password adalah password pengguna yang sudah di-hash menggunakan bcrypt.
	// `json:"-"` berarti field ini TIDAK akan disertakan dalam respons JSON.
	// Ini adalah praktik keamanan penting agar password tidak terekspos ke client.
	Password string `json:"-" gorm:"not null"`

	// Role adalah peran pengguna dalam sistem, contoh: "admin" atau "user".
	// `default:'user'` memberikan nilai default "user" jika tidak diisi.
	Role string `json:"role" gorm:"default:'user'"`

	// CreatedAt adalah waktu saat data pengguna pertama kali dibuat.
	// GORM otomatis mengisi field ini saat record dibuat.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt adalah waktu terakhir data pengguna diperbarui.
	// GORM otomatis mengupdate field ini setiap kali record diupdate.
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt digunakan untuk implementasi "soft delete".
	// Daripada menghapus data dari database, GORM hanya mengisi field ini.
	// `gorm:"index"` membuat index untuk performa query yang lebih baik.
	// `json:"-"` menyembunyikan field ini dari respons JSON.
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName menentukan nama tabel database untuk entitas User.
// Dengan mendefinisikan method ini, kita bisa menggunakan nama tabel custom
// daripada nama default yang dihasilkan GORM (biasanya nama struct dalam bentuk plural).
func (User) TableName() string {
	// Kembalikan nama tabel "users" sebagai nama tabel di database.
	return "users"
}

// UserRepository mendefinisikan kontrak (interface) untuk operasi data pengguna.
// Interface ini dideklarasikan di lapisan Domain agar Use Case dapat menggunakannya
// tanpa mengetahui detail implementasi database yang sebenarnya.
// Prinsip ini disebut "Dependency Inversion" dalam Clean Architecture.
type UserRepository interface {
	// FindAll mengambil semua data pengguna dari database.
	// Mengembalikan slice of User dan error jika terjadi masalah.
	FindAll() ([]User, error)

	// FindByID mengambil satu data pengguna berdasarkan ID-nya.
	// Mengembalikan pointer to User (nil jika tidak ditemukan) dan error.
	FindByID(id uint) (*User, error)

	// FindByEmail mengambil satu data pengguna berdasarkan alamat email.
	// Digunakan untuk login dan validasi duplikasi email saat registrasi.
	FindByEmail(email string) (*User, error)

	// Create menyimpan data pengguna baru ke database.
	// Menerima pointer to User dan mengembalikan error jika gagal.
	Create(user *User) error

	// Update memperbarui data pengguna yang sudah ada di database.
	// Menerima pointer to User (harus memiliki ID) dan mengembalikan error jika gagal.
	Update(user *User) error

	// Delete menghapus data pengguna dari database berdasarkan ID.
	// Implementasi menggunakan soft delete (mengisi DeletedAt, bukan benar-benar dihapus).
	Delete(id uint) error
}

// UserUsecase mendefinisikan kontrak (interface) untuk business logic pengguna.
// Interface ini memudahkan testing karena handler HTTP hanya bergantung pada interface,
// bukan implementasi konkret — sehingga kita bisa mock usecase saat unit testing.
type UserUsecase interface {
	// GetAllUsers mengambil daftar semua pengguna.
	GetAllUsers() ([]User, error)

	// GetUserByID mengambil detail pengguna berdasarkan ID.
	GetUserByID(id uint) (*User, error)

	// RegisterUser mendaftarkan pengguna baru dengan validasi bisnis.
	// Mengembalikan pengguna yang sudah dibuat dan error jika gagal.
	RegisterUser(user *User) (*User, error)

	// LoginUser memvalidasi kredensial dan mengembalikan token JWT jika valid.
	// Mengembalikan token string dan error jika kredensial tidak valid.
	LoginUser(email, password string) (string, error)

	// UpdateUser memperbarui informasi pengguna yang sudah ada.
	UpdateUser(user *User) (*User, error)

	// DeleteUser menghapus pengguna berdasarkan ID.
	DeleteUser(id uint) error
}
