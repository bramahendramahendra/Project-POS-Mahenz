-- =============================================================
-- Migration 005: Portal menjadi klien BE User (Fase 3B, tahap 3B.3)
--
-- Akun, password, role, aplikasi, akses, dan kode SSO kini dikelola BE User
-- (pos_auth_db). Portal hanya menyimpan profil lokal (users + auth_user_id)
-- dan sesinya sendiri.
--
-- Kolom/tabel lama (users.password, users.role_id, users.must_change_password,
-- roles, apps, user_apps, sso_codes) TIDAK lagi dipakai kode, tetapi sengaja
-- belum dihapus: datanya dipindah ke BE User pada tahap 3B.5, baru dibuang
-- dengan migrasi berikutnya.
-- =============================================================

-- Profil lokal ditautkan ke akun BE User lewat auth_user_id.
ALTER TABLE users ADD COLUMN auth_user_id INT NULL AFTER id;
ALTER TABLE users ADD UNIQUE KEY uq_users_auth_user_id (auth_user_id);

-- Profil baru dibuat tanpa password (password ada di BE User).
ALTER TABLE users MODIFY password VARCHAR(255) NULL;

-- Username bukan lagi kunci akun (sumber kebenarannya BE User), hanya salinan tampilan.
ALTER TABLE users DROP INDEX username;
ALTER TABLE users ADD INDEX idx_users_username (username);

-- Status wajib ganti password disalin ke sesi saat login (dari jawaban BE User),
-- dipakai middleware Portal untuk membatasi endpoint sampai password diganti.
ALTER TABLE sessions ADD COLUMN must_change_password TINYINT(1) NOT NULL DEFAULT 0 AFTER user_role;

-- Kunci aplikasi Portal untuk memanggil BE User (header X-App-Key).
-- Isi setelah kunci dibuat di BE User (cd AUTH/BE lalu go run ./cmd/set_app_key --app portal --confirm)
--   INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'portal', 'pak_xxx')
--   ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key)
-- Catatan: kunci tersimpan apa adanya, jadi file backup SQL harus diperlakukan rahasia.
CREATE TABLE IF NOT EXISTS auth_credentials (
    id         INT          PRIMARY KEY DEFAULT 1,
    app_id     VARCHAR(50)  NOT NULL,
    app_key    VARCHAR(255) NOT NULL,
    updated_at DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
