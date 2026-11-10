-- =============================================================
-- Migration 004: Aplikasi, akses user ke aplikasi, dan kode SSO
-- =============================================================

-- Daftar aplikasi yang bisa dibuka dari Portal.
--   code         : pengenal aplikasi (dipakai sebagai app_id saat SSO), contoh 'toko'
--   path         : alamat aplikasi, path di domain yang sama ('/toko/') atau URL penuh
--   app_key_hash : SHA-256 (hex) dari kunci aplikasi. Kunci asli tidak disimpan di Portal.
CREATE TABLE IF NOT EXISTS apps (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    code         VARCHAR(50)  UNIQUE NOT NULL,
    name         VARCHAR(100) NOT NULL,
    description  VARCHAR(255) NULL,
    path         VARCHAR(255) NOT NULL,
    icon         VARCHAR(50)  NOT NULL DEFAULT 'app',
    sort_order   INT          NOT NULL DEFAULT 0,
    is_active    TINYINT(1)   NOT NULL DEFAULT 1,
    app_key_hash CHAR(64)     NULL,
    created_at   DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- User boleh membuka aplikasi mana saja.
CREATE TABLE IF NOT EXISTS user_apps (
    user_id    INT      NOT NULL,
    app_id     INT      NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, app_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (app_id)  REFERENCES apps(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Kode SSO sekali pakai (berlaku singkat). Yang disimpan hanya hash SHA-256 dari kode.
CREATE TABLE IF NOT EXISTS sso_codes (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    code_hash  CHAR(64) UNIQUE NOT NULL,
    user_id    INT      NOT NULL,
    app_id     INT      NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at    DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (app_id)  REFERENCES apps(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------------
-- Seed: aplikasi POS Toko dan akses akun awal.
-- Kunci aplikasi (app_key_hash) sengaja dikosongkan: isi lewat SQL setelah kunci dibuat.
-- -------------------------------------------------------------
INSERT IGNORE INTO apps (code, name, description, path, icon, sort_order) VALUES
    ('toko', 'POS Toko', 'Kasir, stok, pembelian, dan laporan toko', '/toko/', 'store', 1);

INSERT IGNORE INTO user_apps (user_id, app_id)
SELECT u.id, a.id FROM users u CROSS JOIN apps a
WHERE a.code = 'toko' AND u.username IN ('owner', 'admin');
