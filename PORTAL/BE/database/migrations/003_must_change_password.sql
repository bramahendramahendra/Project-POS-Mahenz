-- =============================================================
-- Migration 003: Wajib ganti password
-- must_change_password = 1 artinya password saat ini bersifat sementara
-- (hasil reset oleh admin) dan user wajib menggantinya sebelum memakai
-- fitur lain di Portal.
-- =============================================================

ALTER TABLE users ADD COLUMN must_change_password TINYINT(1) NOT NULL DEFAULT 0 AFTER is_active;
