-- =============================================================
-- Migration 001: Skema awal BE User (layanan akun, keputusan 63)
-- Data akun dipindah dari Portal: user, role pengelola akun, aplikasi,
-- akses user ke aplikasi, dan kode SSO. Tidak ada tabel sesi: BE User
-- tidak membuat sesi browser, setiap aplikasi membuat sesinya sendiri.
-- Asal: migrasi Portal 001, 003, 004 (digabung).
-- =============================================================

-- -------------------------------------------------------------
-- Akun
-- -------------------------------------------------------------

-- Role pengelola akun hanya 2 (admin, user) — keputusan 65.
CREATE TABLE IF NOT EXISTS roles (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(50)  UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description  VARCHAR(255) NULL,
    is_system    TINYINT(1)   DEFAULT 0,
    is_active    TINYINT(1)   DEFAULT 1,
    created_at   DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- must_change_password = 1: password sementara (dibuat/di-reset admin),
-- login ke aplikasi mana pun ditolak sampai diganti di Portal (keputusan 70).
CREATE TABLE IF NOT EXISTS users (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    username             VARCHAR(50)  UNIQUE NOT NULL,
    password             VARCHAR(255) NOT NULL,
    full_name            VARCHAR(100) NOT NULL,
    role_id              INT          NULL,
    is_active            TINYINT(1)   DEFAULT 1,
    must_change_password TINYINT(1)   NOT NULL DEFAULT 0,
    created_at           DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------------
-- Aplikasi, akses, dan SSO
-- -------------------------------------------------------------

-- Aplikasi yang boleh berbicara dengan BE User (termasuk Portal sendiri).
--   code           : pengenal aplikasi (app_id), contoh 'toko'
--   path           : alamat aplikasi, path di domain yang sama ('/toko/') atau URL penuh (dev)
--   app_key_hash   : SHA-256 (hex) dari kunci aplikasi. Kunci asli tidak disimpan di sini.
--   show_in_portal : 1 = tampil sebagai kartu di halaman Aplikasi Portal, 0 = tidak (mis. 'portal')
CREATE TABLE IF NOT EXISTS apps (
    id             INT AUTO_INCREMENT PRIMARY KEY,
    code           VARCHAR(50)  UNIQUE NOT NULL,
    name           VARCHAR(100) NOT NULL,
    description    VARCHAR(255) NULL,
    path           VARCHAR(255) NOT NULL,
    icon           VARCHAR(50)  NOT NULL DEFAULT 'app',
    sort_order     INT          NOT NULL DEFAULT 0,
    is_active      TINYINT(1)   NOT NULL DEFAULT 1,
    show_in_portal TINYINT(1)   NOT NULL DEFAULT 1,
    app_key_hash   CHAR(64)     NULL,
    created_at     DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
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
-- Logging
-- -------------------------------------------------------------

CREATE TABLE IF NOT EXISTS log_schedulers (
    id             VARCHAR(36)              NOT NULL PRIMARY KEY,
    scheduler_name VARCHAR(100)             NOT NULL,
    status         ENUM('success','failed') NOT NULL,
    message        TEXT                     NULL,
    duration_ms    INT                      NULL,
    executed_at    DATETIME                 NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS log_requests (
    id            VARCHAR(36)  NOT NULL PRIMARY KEY,
    method        VARCHAR(10)  NOT NULL,
    endpoint      VARCHAR(255) NOT NULL,
    status_code   SMALLINT     NULL,
    request_body  TEXT         NULL,
    response_body TEXT         NULL,
    user_id       INT          NULL,
    duration_ms   INT          NULL,
    ip_address    VARCHAR(45)  NULL,
    error_message TEXT         NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
