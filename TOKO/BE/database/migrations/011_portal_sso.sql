-- =============================================================
-- Migration 011: Integrasi login terpusat (Portal SSO)
-- Aditif: login lama (username + password) tetap berfungsi.
-- =============================================================

-- Tautan ke akun di Portal. NULL = akun lokal yang belum ditautkan.
ALTER TABLE users ADD COLUMN portal_user_id INT NULL;

ALTER TABLE users ADD UNIQUE KEY uq_users_portal_user_id (portal_user_id);

-- Role untuk profil lokal yang dibuat otomatis saat SSO pertama.
-- Sengaja tanpa baris role_menu_access, jadi tidak punya izin menu apa pun
-- sampai admin mengatur role yang sebenarnya.
INSERT IGNORE INTO roles (name, display_name, description, is_system, is_active) VALUES
    ('menunggu_akses', 'Menunggu Akses', 'Akun baru dari Portal, belum diberi izin oleh admin', 1, 1);

-- Kredensial aplikasi ini untuk menukar kode SSO ke Portal. Satu baris saja (id = 1).
--   app_id  : kode aplikasi ini di Portal (kolom apps.code), mis. 'toko'
--   app_key : kunci aplikasi yang diberikan Portal saat aplikasi didaftarkan (tampil sekali)
-- Isi lewat SQL manual, contoh:
--   INSERT INTO portal_credentials (id, app_id, app_key) VALUES (1, 'toko', 'pak_xxx')
--   ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key);
-- Catatan: kunci tersimpan apa adanya, jadi file backup SQL harus diperlakukan rahasia.
CREATE TABLE IF NOT EXISTS portal_credentials (
    id         INT          PRIMARY KEY DEFAULT 1,
    app_id     VARCHAR(50)  NOT NULL,
    app_key    VARCHAR(255) NOT NULL,
    updated_at DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Login langsung (username + password di Toko): true = form login tampil dan endpoint login aktif.
INSERT IGNORE INTO settings (setting_key, setting_value) VALUES ('direct_login_enabled', 'true');
