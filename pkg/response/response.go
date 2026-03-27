// Package response menyediakan format respons JSON yang konsisten untuk seluruh API.
//
// Dengan menggunakan format respons yang seragam, client (frontend, mobile app)
// dapat mengharapkan struktur yang sama dari semua endpoint, memudahkan parsing.
//
// Format respons standar:
//
//	{
//	    "success": true/false,
//	    "message": "Pesan deskriptif",
//	    "data": { ... } atau null
//	}
package response

// JSONResponse adalah struct yang mendefinisikan format respons JSON standar.
// Semua endpoint API menggunakan struct ini untuk respons yang konsisten.
type JSONResponse struct {
	// Success menunjukkan apakah request berhasil (true) atau gagal (false).
	// Client dapat memeriksa field ini untuk menentukan keberhasilan request.
	Success bool `json:"success"`

	// Message adalah pesan deskriptif yang menjelaskan hasil request.
	// Contoh: "Data berhasil diambil", "Email sudah terdaftar", dll.
	Message string `json:"message"`

	// Data berisi payload respons (data yang diminta).
	// Menggunakan `interface{}` agar bisa menampung tipe data apapun (struct, slice, map, nil).
	// `omitempty` berarti field ini tidak akan disertakan dalam JSON jika nilainya nil.
	Data interface{} `json:"data,omitempty"`
}

// Success membuat respons sukses dengan data payload.
// Digunakan ketika request berhasil diproses dan ada data untuk dikembalikan.
//
// Parameter:
//   - message: pesan sukses yang deskriptif.
//   - data: payload data yang akan dikirim ke client.
//
// Contoh penggunaan:
//
//	c.JSON(http.StatusOK, response.Success("Data berhasil diambil", users))
//
// Menghasilkan JSON:
//
//	{"success": true, "message": "Data berhasil diambil", "data": [...]}
func Success(message string, data interface{}) JSONResponse {
	// Buat dan kembalikan respons sukses dengan data yang diberikan.
	return JSONResponse{
		// Set success ke true karena ini adalah respons berhasil.
		Success: true,

		// Set pesan yang diberikan.
		Message: message,

		// Set data payload yang diberikan.
		Data: data,
	}
}

// Error membuat respons error tanpa data payload.
// Digunakan ketika request gagal karena validasi, not found, unauthorized, dll.
//
// Parameter:
//   - message: pesan error yang menjelaskan kenapa request gagal.
//
// Contoh penggunaan:
//
//	c.JSON(http.StatusBadRequest, response.Error("Email sudah terdaftar"))
//
// Menghasilkan JSON:
//
//	{"success": false, "message": "Email sudah terdaftar"}
func Error(message string) JSONResponse {
	// Buat dan kembalikan respons error tanpa data.
	return JSONResponse{
		// Set success ke false karena ini adalah respons gagal.
		Success: false,

		// Set pesan error yang diberikan.
		Message: message,

		// Data tidak diisi (nil) untuk respons error — client tidak perlu data.
		Data: nil,
	}
}
