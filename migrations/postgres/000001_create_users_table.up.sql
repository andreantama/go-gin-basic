-- Migration: 000001_create_users_table (UP)
-- Membuat tabel users beserta index yang diperlukan.
-- Migration ini dieksekusi otomatis oleh golang-migrate saat aplikasi pertama kali dijalankan.

CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL    NOT NULL,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(50)  NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ  DEFAULT NULL,
    updated_at TIMESTAMPTZ  DEFAULT NULL,
    deleted_at TIMESTAMPTZ  DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT idx_users_email UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
