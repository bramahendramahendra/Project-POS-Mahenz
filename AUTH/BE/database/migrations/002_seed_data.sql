-- =============================================================
-- Migration 002: Seed Data — BE User
-- Gunakan INSERT IGNORE agar aman jika dijalankan ulang.
-- =============================================================

-- -------------------------------------------------------------
-- Role pengelola akun (keputusan 65)
-- -------------------------------------------------------------
INSERT IGNORE INTO roles (name, display_name, description, is_system, is_active) VALUES
    ('admin', 'Admin', 'Mengelola user, aplikasi, dan akses',          1, 1),
    ('user',  'User',  'Pengguna biasa, memakai aplikasi yang diizinkan', 1, 1);

-- -------------------------------------------------------------
-- Akun awal: admin / admin123 (hash Argon2id sama dengan seed Toko dan Portal).
-- Ganti password sebelum production. User Toko lama disalin ke sini oleh
-- skrip migrasi saat cutover (tahap 3B.5).
-- -------------------------------------------------------------
INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'admin', '$argon2id$v=19$m=65536,t=3,p=4$TJxWkwEYoSuvSZ0JtwrfJw$ne9pID5QjltGeu0cNxiKYGCOXkUCRRjunQUYvttxCsM', 'Administrator', id
FROM roles WHERE name = 'admin';

-- -------------------------------------------------------------
-- Aplikasi yang berbicara dengan BE User.
-- Kunci aplikasi (app_key_hash) sengaja kosong: dibuat dengan perintah
-- di BE User lalu dipasang di database aplikasinya (tahap 3B.2-3B.4).
-- 'portal' tidak tampil sebagai kartu (show_in_portal = 0).
-- -------------------------------------------------------------
INSERT IGNORE INTO apps (code, name, description, path, icon, sort_order, show_in_portal) VALUES
    ('portal', 'Portal',     'Login terpusat dan pengelolaan akun',       '/',        'app',    0, 0),
    ('toko',   'POS Toko',   'Kasir, stok, pembelian, dan laporan toko',  '/toko/',   'store',  1, 1),
    ('coffee', 'POS Coffee', 'Kasir kedai kopi',                          '/coffee/', 'coffee', 2, 1);

-- Akses akun admin ke semua aplikasi.
INSERT IGNORE INTO user_apps (user_id, app_id)
SELECT u.id, a.id FROM users u CROSS JOIN apps a
WHERE u.username = 'admin';
