-- =============================================================
-- Migration 012: Toko menjadi klien BE User (Fase 3B, tahap 3B.4)
--
-- Akun (username, password, aktif, wajib ganti password, akses aplikasi) dikelola
-- BE User (pos_auth_db). Toko menyimpan profil lokal: role, PIN, aktif/nonaktif,
-- dan riwayat transaksi. Password lokal tidak dipakai lagi (keputusan 66).
-- =============================================================

-- Tautan profil lokal ke akun BE User (keputusan 71).
ALTER TABLE users RENAME COLUMN portal_user_id TO auth_user_id;
ALTER TABLE users RENAME INDEX uq_users_portal_user_id TO uq_users_auth_user_id;

-- Isi lama adalah id user Portal (sebelum BE User), tidak sama dengan id di BE User.
-- Dikosongkan, lalu ditautkan ulang otomatis saat login pertama yang terverifikasi
-- BE User ke profil dengan username sama (keputusan 67 dan 77). Production belum
-- pernah memakai kolom ini (Toko baru belum di-deploy).
UPDATE users SET auth_user_id = NULL WHERE auth_user_id IS NOT NULL;

-- Profil baru dari BE User dibuat tanpa password lokal.
-- Hash lama dibiarkan: dipakai skrip penyalinan user Toko ke BE User (tahap 3B.5).
ALTER TABLE users MODIFY password VARCHAR(255) NULL;

-- Kunci aplikasi Toko untuk memanggil BE User (header X-App-Key).
-- Kunci lama (tabel portal_credentials, untuk Portal) tidak berlaku di BE User: tabel lama
-- tidak dipakai kode lagi tetapi belum dihapus, dan dibuat tabel baru yang kosong.
-- Isi setelah kunci dibuat di BE User
-- (cd AUTH/BE lalu go run ./cmd/set_app_key --app toko --confirm)
--   INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'toko', 'pak_xxx')
--   ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key)
-- Catatan: kunci tersimpan apa adanya, jadi file backup SQL harus diperlakukan rahasia.
CREATE TABLE IF NOT EXISTS auth_credentials (
    id         INT          PRIMARY KEY DEFAULT 1,
    app_id     VARCHAR(50)  NOT NULL,
    app_key    VARCHAR(255) NOT NULL,
    updated_at DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

