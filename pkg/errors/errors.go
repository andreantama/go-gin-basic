// Package errors menyediakan tipe error kustom dan helper untuk penanganan error.
//
// Dengan mendefinisikan tipe error kustom, handler dapat menentukan HTTP status code
// yang tepat berdasarkan tipe error yang dikembalikan oleh use case.
// Ini memisahkan logika HTTP dari logika bisnis — use case tidak perlu tahu
// tentang HTTP status code, cukup kembalikan tipe error yang sesuai.
package errors

import (
	// errors digunakan untuk memeriksa tipe error menggunakan errors.As().
	"errors"

	// net/http menyediakan konstanta HTTP status code.
	"net/http"
)

// NotFoundError merepresentasikan error ketika resource yang diminta tidak ditemukan.
// Akan dikonversi ke HTTP 404 Not Found oleh fungsi HTTPStatus().
type NotFoundError struct {
	// Message adalah pesan error yang deskriptif.
	Message string
}

// Error mengimplementasikan interface error untuk NotFoundError.
// Setiap tipe error di Go harus mengimplementasikan method Error() string.
func (e *NotFoundError) Error() string {
	// Kembalikan pesan error dari struct.
	return e.Message
}

// ValidationError merepresentasikan error ketika validasi input gagal.
// Akan dikonversi ke HTTP 400 Bad Request oleh fungsi HTTPStatus().
type ValidationError struct {
	// Message adalah pesan validasi yang menjelaskan apa yang salah.
	Message string
}

// Error mengimplementasikan interface error untuk ValidationError.
func (e *ValidationError) Error() string {
	// Kembalikan pesan validasi dari struct.
	return e.Message
}

// ConflictError merepresentasikan error ketika terjadi konflik data (contoh: email duplikat).
// Akan dikonversi ke HTTP 409 Conflict oleh fungsi HTTPStatus().
type ConflictError struct {
	// Message adalah pesan error yang menjelaskan konflik yang terjadi.
	Message string
}

// Error mengimplementasikan interface error untuk ConflictError.
func (e *ConflictError) Error() string {
	// Kembalikan pesan error dari struct.
	return e.Message
}

// UnauthorizedError merepresentasikan error ketika akses tidak diizinkan.
// Akan dikonversi ke HTTP 401 Unauthorized oleh fungsi HTTPStatus().
type UnauthorizedError struct {
	// Message adalah pesan error yang menjelaskan mengapa akses ditolak.
	Message string
}

// Error mengimplementasikan interface error untuk UnauthorizedError.
func (e *UnauthorizedError) Error() string {
	// Kembalikan pesan error dari struct.
	return e.Message
}

// HTTPStatus mengkonversi tipe error kustom ke HTTP status code yang sesuai.
// Handler menggunakan fungsi ini untuk menentukan status code respons berdasarkan tipe error.
//
// Mapping error ke status code:
//   - NotFoundError → 404 Not Found
//   - ValidationError → 400 Bad Request
//   - ConflictError → 409 Conflict
//   - UnauthorizedError → 401 Unauthorized
//   - Error lainnya → 500 Internal Server Error
//
// Contoh penggunaan di handler:
//
//	if err != nil {
//	    statusCode := pkgerrors.HTTPStatus(err)
//	    c.JSON(statusCode, response.Error(err.Error()))
//	    return
//	}
func HTTPStatus(err error) int {
	// Periksa apakah error adalah NotFoundError menggunakan errors.As().
	// errors.As() lebih baik daripada type assertion karena mendukung error wrapping.
	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		// Kembalikan 404 Not Found.
		return http.StatusNotFound
	}

	// Periksa apakah error adalah ValidationError.
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		// Kembalikan 400 Bad Request.
		return http.StatusBadRequest
	}

	// Periksa apakah error adalah ConflictError.
	var conflictErr *ConflictError
	if errors.As(err, &conflictErr) {
		// Kembalikan 409 Conflict.
		return http.StatusConflict
	}

	// Periksa apakah error adalah UnauthorizedError.
	var unauthorizedErr *UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		// Kembalikan 401 Unauthorized.
		return http.StatusUnauthorized
	}

	// Untuk error lain yang tidak dikenal, kembalikan 500 Internal Server Error.
	// Ini menunjukkan ada masalah yang tidak terduga di sisi server.
	return http.StatusInternalServerError
}

// NewNotFoundError membuat instance NotFoundError baru dengan pesan yang diberikan.
// Helper function ini memudahkan pembuatan error di use case.
//
// Contoh:
//
//	return nil, errors.NewNotFoundError("pengguna tidak ditemukan")
func NewNotFoundError(message string) error {
	// Buat dan kembalikan NotFoundError baru.
	return &NotFoundError{Message: message}
}

// NewValidationError membuat instance ValidationError baru.
func NewValidationError(message string) error {
	// Buat dan kembalikan ValidationError baru.
	return &ValidationError{Message: message}
}

// NewConflictError membuat instance ConflictError baru.
func NewConflictError(message string) error {
	// Buat dan kembalikan ConflictError baru.
	return &ConflictError{Message: message}
}

// NewUnauthorizedError membuat instance UnauthorizedError baru.
func NewUnauthorizedError(message string) error {
	// Buat dan kembalikan UnauthorizedError baru.
	return &UnauthorizedError{Message: message}
}
