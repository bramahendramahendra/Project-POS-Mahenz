-- =============================================================
-- Migration 006: Buang data akun lama Portal (Fase 3B, tahap 3B.5, keputusan 79)
--
-- Sejak 3B.3 akun, role, aplikasi, akses, dan kode SSO dikelola BE User (pos_auth_db).
-- Kolom/tabel di bawah tidak dipakai kode lagi. Data dev Portal tidak dipindah ke
-- BE User (keputusan 78). Portal belum pernah di-deploy, jadi production tidak terdampak.
-- Backup DB dev diambil sebelum migrasi ini dijalankan.
-- =============================================================

-- Tabel anak dulu (foreign key ke users dan apps).
DROP TABLE IF EXISTS sso_codes;
DROP TABLE IF EXISTS user_apps;
DROP TABLE IF EXISTS apps;

-- users.role_id tidak memakai foreign key, jadi tabel roles bisa langsung dibuang.
DROP TABLE IF EXISTS roles;

ALTER TABLE users DROP COLUMN password;
ALTER TABLE users DROP COLUMN role_id;
ALTER TABLE users DROP COLUMN must_change_password;
