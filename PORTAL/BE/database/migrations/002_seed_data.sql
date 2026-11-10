-- =============================================================
-- Migration 002: Seed Data — POS Portal
-- Gunakan INSERT IGNORE agar aman jika dijalankan ulang.
-- =============================================================

-- -------------------------------------------------------------
-- Roles (2 role sistem, tidak bisa dihapus)
-- -------------------------------------------------------------
INSERT IGNORE INTO roles (name, display_name, description, is_system, is_active) VALUES
    ('admin', 'Admin', 'Mengelola user, aplikasi, dan akses Portal', 1, 1),
    ('user',  'User',  'Pengguna biasa, memilih aplikasi yang diizinkan', 1, 1);

-- -------------------------------------------------------------
-- Akun awal: satu admin agar Portal bisa dipakai pertama kali.
-- password: admin -> 'admin123' (sama dengan seed POS Toko)
-- WAJIB diganti sebelum production. Skrip cmd/migrate_users_to_portal (di TOKO/BE)
-- mengambil alih akun ini bila ada user "admin" di Toko.
-- Aturan role: hanya "admin" Toko yang menjadi admin Portal, sehingga
-- akun "owner" TIDAK diseed di sini.
-- -------------------------------------------------------------
INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'admin', '$argon2id$v=19$m=65536,t=3,p=4$TJxWkwEYoSuvSZ0JtwrfJw$ne9pID5QjltGeu0cNxiKYGCOXkUCRRjunQUYvttxCsM', 'Administrator', id
FROM roles WHERE name = 'admin';
