-- Migration: 000001_create_users_table (UP)
-- Membuat tabel users beserta index yang diperlukan.
-- Migration ini dieksekusi otomatis oleh golang-migrate saat aplikasi pertama kali dijalankan.

CREATE TABLE IF NOT EXISTS `users` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(255)    NOT NULL,
    `email`      VARCHAR(255)    NOT NULL,
    `password`   VARCHAR(255)    NOT NULL,
    `role`       VARCHAR(50)     NOT NULL DEFAULT 'user',
    `created_at` DATETIME(3)     DEFAULT NULL,
    `updated_at` DATETIME(3)     DEFAULT NULL,
    `deleted_at` DATETIME(3)     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_email`      (`email`),
    KEY         `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
