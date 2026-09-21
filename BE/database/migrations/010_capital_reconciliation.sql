-- =============================================================
-- 010_capital_reconciliation.sql
-- Menu "Rekonsiliasi Modal" (Pelaporan) — koreksi manual modal (HPP) per baris
-- penjualan yang TIDAK bisa dibetulkan otomatis oleh skrip backfill_transaction_cogs.
-- HANYA untuk role admin.
--
-- Dijalankan otomatis oleh RunMigrations (database/migrate.go) sekali, tercatat
-- di migrations_history. Semua INSERT pakai IGNORE agar aman bila ter-retry.
-- =============================================================

-- -------------------------------------------------------------
-- 1) Tabel daftar tinjauan: 1 baris = 1 transaction_item yang modalnya
--    (purchase_price) meragukan dan perlu dikoreksi manual admin.
--    Diisi oleh skrip backfill_transaction_cogs (baris berstatus MANUAL),
--    dan diselesaikan lewat menu Rekonsiliasi Modal.
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS transaction_cost_review (
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    transaction_item_id INT           NOT NULL,
    transaction_id      INT           NOT NULL,
    product_id          INT           NULL,
    -- snapshot info saat baris ditandai (biar menu informatif tanpa join berat)
    reason              VARCHAR(255)  NOT NULL DEFAULT '',
    old_purchase_price  DECIMAL(15,2) NOT NULL DEFAULT 0,  -- modal tersimpan saat ditandai
    suggested_price     DECIMAL(15,2) NOT NULL DEFAULT 0,  -- hasil hitung otomatis (meragukan)
    sell_price          DECIMAL(15,2) NOT NULL DEFAULT 0,  -- harga jual per satuan (acuan)
    -- status tinjauan
    status              ENUM('pending','resolved','skipped') NOT NULL DEFAULT 'pending',
    resolved_price      DECIMAL(15,2) NULL,                -- modal benar yang di-input admin
    note                VARCHAR(500)  NULL,
    resolved_by         INT           NULL,
    resolved_at         DATETIME      NULL,
    created_at          DATETIME      DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_cost_review_item (transaction_item_id),
    INDEX idx_cost_review_status (status),
    INDEX idx_cost_review_product (product_id),
    FOREIGN KEY (transaction_item_id) REFERENCES transaction_items(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id)      REFERENCES transactions(id)      ON DELETE CASCADE,
    FOREIGN KEY (product_id)          REFERENCES products(id)          ON DELETE SET NULL,
    FOREIGN KEY (resolved_by)         REFERENCES users(id)             ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------------
-- 2) Menu (child dari grup 'pelaporan'), tepat di bawah Rekonsiliasi Stok.
--    order_index 7 (Rekonsiliasi Stok = 6).
-- -------------------------------------------------------------
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.rekonsiliasi_modal', 'Rekonsiliasi Modal', 'Coins', '/reports/capital-reconciliation', 7
FROM menus m WHERE m.key_name = 'pelaporan';

-- -------------------------------------------------------------
-- 3) Whitelist path untuk dropdown di Manajemen Menu
-- -------------------------------------------------------------
INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/reports/capital-reconciliation', 'Rekonsiliasi Modal');

-- -------------------------------------------------------------
-- 4) Hak akses HANYA admin (view + edit untuk koreksi manual).
--    Owner & kasir TIDAK diberi akses => tidak melihat menu ini.
-- -------------------------------------------------------------
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 1, 0
FROM roles r
JOIN menus m ON m.key_name = 'pelaporan.rekonsiliasi_modal'
WHERE r.name = 'admin';
