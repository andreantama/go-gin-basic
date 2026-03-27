// Package repository berisi implementasi konkret dari interface repository yang
// didefinisikan di lapisan domain. Repository bertanggung jawab untuk semua
// operasi baca/tulis data ke database menggunakan GORM sebagai ORM.
//
// Dalam Clean Architecture, lapisan ini termasuk dalam "Interface Adapters" —
// ia mengadaptasi interface domain ke implementasi database yang nyata.
// Use Case tidak tahu bahwa implementasi ini menggunakan MySQL/GORM.
package repository

import (
	// domain berisi interface UserRepository dan struct User yang akan kita gunakan.
	"github.com/andreantama/go-gin-basic/internal/domain"

	// gorm adalah ORM yang digunakan untuk berinteraksi dengan database.
	"gorm.io/gorm"
)

// userRepository adalah struct yang mengimplementasikan interface domain.UserRepository.
// Field `db` adalah koneksi database GORM yang diinjeksikan saat pembuatan instance.
// Struct ini bersifat private (huruf kecil) — hanya bisa dibuat melalui fungsi NewUserRepository.
type userRepository struct {
	// db adalah instance koneksi database GORM.
	// Menggunakan *gorm.DB karena GORM bekerja dengan pointer ke struct DB.
	db *gorm.DB
}

// NewUserRepository adalah constructor function (factory function) untuk membuat
// instance baru dari userRepository.
// Fungsi ini mengembalikan interface domain.UserRepository, bukan struct konkret.
// Ini adalah pola "return interface, accept interface" yang umum di Go.
//
// Parameter:
//   - db: instance koneksi database GORM yang sudah dikonfigurasi.
//
// Contoh penggunaan:
//
//	userRepo := repository.NewUserRepository(db)
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	// Buat dan kembalikan instance baru userRepository dengan koneksi database yang diberikan.
	return &userRepository{db: db}
}

// FindAll mengambil semua data pengguna dari database.
// Method ini mengimplementasikan domain.UserRepository.FindAll().
// GORM akan mengeksekusi query: SELECT * FROM users WHERE deleted_at IS NULL
func (r *userRepository) FindAll() ([]domain.User, error) {
	// Deklarasi slice untuk menampung hasil query.
	var users []domain.User

	// Jalankan query menggunakan GORM.
	// `.Find(&users)` mengisi slice users dengan semua record dari tabel users.
	// GORM otomatis menambahkan kondisi WHERE deleted_at IS NULL untuk soft delete.
	result := r.db.Find(&users)

	// Periksa apakah ada error dari query.
	if result.Error != nil {
		// Kembalikan slice kosong dan error jika query gagal.
		return nil, result.Error
	}

	// Kembalikan data pengguna dan nil (tidak ada error).
	return users, nil
}

// FindByID mengambil satu data pengguna berdasarkan ID-nya.
// Method ini mengimplementasikan domain.UserRepository.FindByID().
// GORM akan mengeksekusi query: SELECT * FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	// Deklarasi variabel untuk menampung satu record pengguna.
	var user domain.User

	// Jalankan query dengan kondisi WHERE id = id menggunakan GORM.
	// `.First(&user, id)` mengambil record pertama yang cocok dengan ID.
	// Jika tidak ditemukan, GORM mengembalikan error `gorm.ErrRecordNotFound`.
	result := r.db.First(&user, id)

	// Periksa apakah ada error (termasuk record not found).
	if result.Error != nil {
		// Kembalikan nil dan error jika record tidak ditemukan atau query gagal.
		return nil, result.Error
	}

	// Kembalikan pointer ke user dan nil (tidak ada error).
	return &user, nil
}

// FindByEmail mengambil satu data pengguna berdasarkan alamat email.
// Method ini mengimplementasikan domain.UserRepository.FindByEmail().
// GORM akan mengeksekusi query: SELECT * FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	// Deklarasi variabel untuk menampung satu record pengguna.
	var user domain.User

	// Jalankan query dengan kondisi WHERE email = email.
	// `.Where("email = ?", email)` adalah cara yang aman untuk menghindari SQL injection.
	// Tanda `?` adalah placeholder yang akan diisi oleh GORM dengan nilai yang di-escape.
	result := r.db.Where("email = ?", email).First(&user)

	// Periksa apakah ada error.
	if result.Error != nil {
		// Kembalikan nil dan error jika record tidak ditemukan atau query gagal.
		return nil, result.Error
	}

	// Kembalikan pointer ke user dan nil (tidak ada error).
	return &user, nil
}

// Create menyimpan data pengguna baru ke database.
// Method ini mengimplementasikan domain.UserRepository.Create().
// GORM akan mengeksekusi query: INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)
func (r *userRepository) Create(user *domain.User) error {
	// Jalankan operasi INSERT menggunakan GORM.
	// `.Create(user)` menyimpan semua field dari struct user ke database.
	// GORM juga otomatis mengisi field ID, CreatedAt, dan UpdatedAt.
	result := r.db.Create(user)

	// Kembalikan error jika operasi INSERT gagal, atau nil jika berhasil.
	return result.Error
}

// Update memperbarui data pengguna yang sudah ada di database.
// Method ini mengimplementasikan domain.UserRepository.Update().
// GORM akan mengeksekusi query: UPDATE users SET ... WHERE id = ?
func (r *userRepository) Update(user *domain.User) error {
	// Jalankan operasi UPDATE menggunakan GORM.
	// `.Save(user)` akan mengupdate semua field (termasuk field zero-value).
	// Ini berbeda dengan `.Updates(user)` yang hanya mengupdate field non-zero.
	result := r.db.Save(user)

	// Kembalikan error jika operasi UPDATE gagal, atau nil jika berhasil.
	return result.Error
}

// Delete menghapus data pengguna dari database berdasarkan ID menggunakan soft delete.
// Method ini mengimplementasikan domain.UserRepository.Delete().
// Soft delete berarti GORM mengisi field deleted_at, bukan benar-benar menghapus record.
// GORM akan mengeksekusi query: UPDATE users SET deleted_at = NOW() WHERE id = ?
func (r *userRepository) Delete(id uint) error {
	// Jalankan operasi soft delete menggunakan GORM.
	// `.Delete(&domain.User{}, id)` mengisi field DeletedAt dengan waktu saat ini.
	// Record tidak benar-benar dihapus dari database, hanya "ditandai" sebagai dihapus.
	result := r.db.Delete(&domain.User{}, id)

	// Kembalikan error jika operasi DELETE gagal, atau nil jika berhasil.
	return result.Error
}
