// Package migrations menyediakan file-file SQL migration yang di-embed ke dalam binary aplikasi.
//
// Menggunakan Go embed (//go:embed), semua file .sql akan dikompilasi langsung
// ke dalam binary sehingga tidak perlu mendistribusikan file SQL terpisah saat deployment.
//
// Struktur direktori:
//
//	migrations/
//	├── mysql/
//	│   ├── 000001_create_users_table.up.sql    # Membuat tabel users (MySQL)
//	│   └── 000001_create_users_table.down.sql  # Menghapus tabel users (MySQL)
//	└── postgres/
//	    ├── 000001_create_users_table.up.sql    # Membuat tabel users (PostgreSQL)
//	    └── 000001_create_users_table.down.sql  # Menghapus tabel users (PostgreSQL)
package migrations

import "embed"

// MigrationFiles adalah embed.FS yang berisi semua file SQL migration untuk
// MySQL (direktori "mysql/") dan PostgreSQL (direktori "postgres/").
//
// Di-embed menggunakan direktif //go:embed sehingga file-file ini menjadi bagian
// dari binary dan tidak perlu di-copy secara terpisah ke server produksi.
//
//go:embed mysql postgres
var MigrationFiles embed.FS
