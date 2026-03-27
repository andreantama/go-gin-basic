// Package usecase berisi implementasi business logic (logika bisnis) aplikasi.
//
// Dalam Clean Architecture, lapisan Use Case (juga disebut "Application Business Rules")
// berisi semua aturan dan logika yang spesifik untuk aplikasi ini.
// Use Case TIDAK boleh bergantung pada:
//   - Framework web (GIN)
//   - Database (MySQL, GORM)
//   - Detail delivery (HTTP, gRPC)
//
// Use Case HANYA bergantung pada lapisan Domain (interface dan entitas).
// Ini memungkinkan use case untuk di-test secara independen tanpa database atau web server.
package usecase

import (
	// errors digunakan untuk membuat dan menginspeksi error.
	"errors"

	// fmt digunakan untuk memformat pesan error.
	"fmt"

	// time digunakan untuk membuat waktu kadaluarsa token JWT.
	"time"

	// domain berisi interface dan entitas yang digunakan oleh use case.
	"github.com/andreantama/go-gin-basic/internal/domain"

	// config berisi konfigurasi aplikasi (JWT secret, expire time, dll).
	"github.com/andreantama/go-gin-basic/config"

	// jwt adalah library untuk membuat dan memvalidasi JSON Web Token.
	"github.com/golang-jwt/jwt/v5"

	// bcrypt digunakan untuk hashing dan verifikasi password.
	"golang.org/x/crypto/bcrypt"

	// gorm digunakan untuk memeriksa tipe error (ErrRecordNotFound).
	"gorm.io/gorm"
)

// userUsecase adalah struct yang mengimplementasikan interface domain.UserUsecase.
// Struct ini menyimpan dependency yang dibutuhkan: repository dan konfigurasi.
type userUsecase struct {
	// userRepo adalah interface repository untuk mengakses data pengguna.
	// Use case bergantung pada INTERFACE, bukan implementasi konkret.
	// Ini memungkinkan dependency injection dan unit testing yang mudah.
	userRepo domain.UserRepository

	// cfg adalah konfigurasi aplikasi yang berisi JWT secret dan expire time.
	cfg *config.Config
}

// NewUserUsecase adalah constructor function untuk membuat instance userUsecase baru.
// Fungsi ini menerima dependency melalui parameter (dependency injection).
// Mengembalikan interface domain.UserUsecase untuk menyembunyikan implementasi konkret.
//
// Parameter:
//   - userRepo: implementasi repository pengguna (bisa real DB atau mock untuk testing).
//   - cfg: konfigurasi aplikasi.
//
// Contoh penggunaan:
//
//	userUC := usecase.NewUserUsecase(userRepo, cfg)
func NewUserUsecase(userRepo domain.UserRepository, cfg *config.Config) domain.UserUsecase {
	// Buat dan kembalikan instance baru dengan dependency yang diinjeksikan.
	return &userUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// GetAllUsers mengambil daftar semua pengguna dari repository.
// Tidak ada business logic kompleks di sini — hanya meneruskan ke repository.
// Method ini mengimplementasikan domain.UserUsecase.GetAllUsers().
func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	// Delegasikan ke repository untuk mengambil semua data pengguna.
	return u.userRepo.FindAll()
}

// GetUserByID mengambil detail pengguna berdasarkan ID.
// Method ini mengimplementasikan domain.UserUsecase.GetUserByID().
func (u *userUsecase) GetUserByID(id uint) (*domain.User, error) {
	// Delegasikan ke repository untuk mencari pengguna berdasarkan ID.
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		// Periksa apakah error adalah "record not found" dari GORM.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Kembalikan error yang lebih deskriptif untuk client.
			return nil, fmt.Errorf("pengguna dengan ID %d tidak ditemukan", id)
		}
		// Kembalikan error lain apa adanya.
		return nil, err
	}

	// Kembalikan data pengguna jika berhasil ditemukan.
	return user, nil
}

// RegisterUser mendaftarkan pengguna baru setelah melakukan validasi bisnis.
// Business logic yang diterapkan:
//  1. Validasi email tidak boleh sudah terdaftar.
//  2. Hash password menggunakan bcrypt sebelum disimpan ke database.
//
// Method ini mengimplementasikan domain.UserUsecase.RegisterUser().
func (u *userUsecase) RegisterUser(user *domain.User) (*domain.User, error) {
	// Langkah 1: Periksa apakah email sudah terdaftar di database.
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Jika ada error selain "record not found", kembalikan error.
		return nil, fmt.Errorf("gagal memeriksa email: %w", err)
	}
	if existingUser != nil {
		// Jika email sudah ada, tolak registrasi dengan pesan error yang jelas.
		return nil, fmt.Errorf("email %s sudah terdaftar", user.Email)
	}

	// Langkah 2: Hash password menggunakan bcrypt sebelum menyimpan ke database.
	// bcrypt.GenerateFromPassword menghasilkan hash yang aman dengan cost factor 12.
	// Semakin tinggi cost factor, semakin aman tapi semakin lambat prosesnya.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// Kembalikan error jika proses hashing gagal.
		return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}
	// Ganti password plain text dengan hash yang dihasilkan bcrypt.
	user.Password = string(hashedPassword)

	// Langkah 3: Set role default ke "user" jika tidak disediakan.
	if user.Role == "" {
		// Berikan role default "user" untuk pengguna baru.
		user.Role = "user"
	}

	// Langkah 4: Simpan pengguna baru ke database melalui repository.
	if err := u.userRepo.Create(user); err != nil {
		// Kembalikan error jika penyimpanan ke database gagal.
		return nil, fmt.Errorf("gagal membuat pengguna: %w", err)
	}

	// Kembalikan data pengguna yang baru dibuat (dengan ID yang sudah diisi oleh database).
	return user, nil
}

// LoginUser memvalidasi kredensial pengguna dan mengembalikan token JWT jika valid.
// Business logic yang diterapkan:
//  1. Cari pengguna berdasarkan email.
//  2. Bandingkan password dengan hash yang tersimpan di database.
//  3. Generate token JWT jika kredensial valid.
//
// Method ini mengimplementasikan domain.UserUsecase.LoginUser().
func (u *userUsecase) LoginUser(email, password string) (string, error) {
	// Langkah 1: Cari pengguna berdasarkan email di database.
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		// Jika email tidak ditemukan, kembalikan pesan error generik.
		// Kita sengaja tidak membedakan antara "email tidak ada" dan "password salah"
		// untuk mencegah "user enumeration attack" (attacker tidak bisa tahu email mana yang terdaftar).
		return "", fmt.Errorf("email atau password salah")
	}

	// Langkah 2: Bandingkan password yang diberikan dengan hash di database.
	// bcrypt.CompareHashAndPassword akan mengembalikan nil jika password cocok.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password tidak cocok — kembalikan pesan error generik yang sama.
		return "", fmt.Errorf("email atau password salah")
	}

	// Langkah 3: Generate token JWT untuk pengguna yang berhasil login.
	token, err := u.generateJWT(user)
	if err != nil {
		// Kembalikan error jika pembuatan token gagal.
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}

	// Kembalikan token JWT yang berhasil dibuat.
	return token, nil
}

// UpdateUser memperbarui informasi pengguna yang sudah ada.
// Business logic yang diterapkan:
//  1. Pastikan pengguna yang akan diupdate ada di database.
//  2. Jika password baru diberikan, hash sebelum disimpan.
//
// Method ini mengimplementasikan domain.UserUsecase.UpdateUser().
func (u *userUsecase) UpdateUser(user *domain.User) (*domain.User, error) {
	// Langkah 1: Pastikan pengguna dengan ID tersebut ada di database.
	existingUser, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		// Kembalikan error jika pengguna tidak ditemukan.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengguna dengan ID %d tidak ditemukan", user.ID)
		}
		return nil, err
	}

	// Langkah 2: Update field yang diberikan, pertahankan field yang tidak diubah.
	// Hanya update nama jika disediakan dalam request.
	if user.Name != "" {
		existingUser.Name = user.Name
	}

	// Langkah 3: Jika password baru diberikan, hash sebelum disimpan.
	if user.Password != "" {
		// Hash password baru menggunakan bcrypt.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("gagal mengenkripsi password baru: %w", err)
		}
		// Simpan hash password baru ke existing user.
		existingUser.Password = string(hashedPassword)
	}

	// Langkah 4: Simpan perubahan ke database melalui repository.
	if err := u.userRepo.Update(existingUser); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengguna: %w", err)
	}

	// Kembalikan data pengguna yang sudah diperbarui.
	return existingUser, nil
}

// DeleteUser menghapus pengguna berdasarkan ID.
// Method ini mengimplementasikan domain.UserUsecase.DeleteUser().
func (u *userUsecase) DeleteUser(id uint) error {
	// Langkah 1: Pastikan pengguna dengan ID tersebut ada sebelum menghapus.
	_, err := u.userRepo.FindByID(id)
	if err != nil {
		// Kembalikan error yang sesuai jika pengguna tidak ditemukan.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pengguna dengan ID %d tidak ditemukan", id)
		}
		return err
	}

	// Langkah 2: Hapus pengguna melalui repository (soft delete).
	if err := u.userRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus pengguna: %w", err)
	}

	// Kembalikan nil jika penghapusan berhasil.
	return nil
}

// generateJWT adalah helper method untuk membuat token JWT.
// Token berisi informasi (claims) tentang pengguna yang digunakan untuk autentikasi.
// Method ini bersifat private — hanya bisa dipanggil dari dalam package usecase.
func (u *userUsecase) generateJWT(user *domain.User) (string, error) {
	// Tentukan waktu kadaluarsa token berdasarkan konfigurasi (default: 24 jam).
	expirationTime := time.Now().Add(time.Duration(u.cfg.JWTExpireHours) * time.Hour)

	// Buat claims (payload) JWT yang berisi informasi pengguna.
	// Claims ini akan dienkripsi dan disertakan dalam token.
	claims := jwt.MapClaims{
		// "sub" (subject) adalah ID pengguna yang sedang login.
		"sub": user.ID,

		// "email" adalah email pengguna, disertakan untuk kemudahan.
		"email": user.Email,

		// "role" adalah peran pengguna, digunakan untuk otorisasi.
		"role": user.Role,

		// "exp" (expiration) adalah waktu kadaluarsa token dalam format Unix timestamp.
		"exp": expirationTime.Unix(),

		// "iat" (issued at) adalah waktu token dibuat dalam format Unix timestamp.
		"iat": time.Now().Unix(),
	}

	// Buat token JWT menggunakan algoritma HMAC-SHA256 (HS256) dan claims di atas.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token menggunakan secret key dari konfigurasi.
	// Tanda tangan ini memastikan token tidak bisa dipalsukan tanpa mengetahui secret key.
	tokenString, err := token.SignedString([]byte(u.cfg.JWTSecret))
	if err != nil {
		// Kembalikan string kosong dan error jika signing gagal.
		return "", err
	}

	// Kembalikan token JWT dalam bentuk string.
	return tokenString, nil
}
