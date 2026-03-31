-- Migration: 000001_create_users_table (DOWN)
-- Menghapus tabel users (rollback dari migration UP).
-- Hati-hati: perintah ini akan menghapus SEMUA data di tabel users.

DROP TABLE IF EXISTS users;
