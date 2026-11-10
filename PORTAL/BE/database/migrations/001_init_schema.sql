-- =============================================================
-- Migration 001: Initial Schema — POS Portal
-- Skema awal Portal (login terpusat). Tabel aplikasi, akses user ke
-- aplikasi, dan kode SSO ditambahkan pada migrasi berikutnya (tahap 2.5).
-- =============================================================

-- -------------------------------------------------------------
-- Auth
-- -------------------------------------------------------------

-- Role Portal hanya 2 (admin, user). Tabel dipertahankan agar struktur
-- query sama dengan aplikasi Toko.
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

CREATE TABLE IF NOT EXISTS users (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    username   VARCHAR(50)   UNIQUE NOT NULL,
    password   VARCHAR(255)  NOT NULL,
    full_name  VARCHAR(100)  NOT NULL,
    role_id    INT           NULL,
    is_active  TINYINT(1)    DEFAULT 1,
    created_at DATETIME      DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Satu sesi aktif per user (UNIQUE unique_user), sama seperti Toko.
-- user_role disimpan sebagai VARCHAR agar token tidak perlu decode untuk cek role.
CREATE TABLE IF NOT EXISTS sessions (
    id            INT AUTO_INCREMENT PRIMARY KEY,
    user_id       INT          NOT NULL,
    user_role     VARCHAR(20)  NOT NULL DEFAULT 'user',
    token         TEXT         NOT NULL,
    refresh_token TEXT         NOT NULL,
    device_info   VARCHAR(100) NULL,
    ip_address    VARCHAR(50)  NULL,
    created_at    DATETIME     DEFAULT CURRENT_TIMESTAMP,
    expires_at    DATETIME     NOT NULL,
    UNIQUE KEY unique_user (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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
