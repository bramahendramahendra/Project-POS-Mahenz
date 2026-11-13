-- =============================================================
-- Migration 013: Buang tabel portal_credentials (Fase 3B, tahap 3B.5, keputusan 79)
--
-- Sejak migrasi 012 kunci Toko untuk BE User disimpan di auth_credentials. Kunci lama
-- (untuk Portal) tidak berlaku di BE User. Tabel ini belum pernah ada di production
-- (migrasi 011 belum di-deploy). Backup DB dev diambil sebelum migrasi ini dijalankan.
--
-- Kolom users.password SENGAJA tetap ada: dibaca skrip cmd/migrate_users_to_auth saat
-- cutover production (hash disalin ke BE User agar kasir tetap login dengan password lama).
-- =============================================================

DROP TABLE IF EXISTS portal_credentials;
