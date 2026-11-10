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
-- Users default (disamakan dengan seed POS Toko)
-- password: admin -> 'admin123' | owner -> 'owner123'
-- Keduanya berperan admin di Portal.
-- Ganti password sebelum deploy ke production.
-- -------------------------------------------------------------
INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'owner', '$argon2id$v=19$m=65536,t=3,p=4$gcWNYu+y/FCRuCT6WS/+eg$1qAjwCU3HoK87s/CbDPHYo06W6B4mhoHF1SQmkigukY', 'Owner', id
FROM roles WHERE name = 'admin';

INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'admin', '$argon2id$v=19$m=65536,t=3,p=4$TJxWkwEYoSuvSZ0JtwrfJw$ne9pID5QjltGeu0cNxiKYGCOXkUCRRjunQUYvttxCsM', 'Administrator', id
FROM roles WHERE name = 'admin';
