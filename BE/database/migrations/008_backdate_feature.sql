-- =============================================================
-- Migration 008: Backdate (Transaksi Lampau) Feature
-- Menambah kolom & menu untuk fitur input transaksi historis.
-- =============================================================

-- -------------------------------------------------------------
-- Schema Changes
-- -------------------------------------------------------------

-- Pembeda kas historis vs kas biasa
ALTER TABLE cash_drawer ADD COLUMN is_backdate TINYINT(1) NOT NULL DEFAULT 0;

-- Audit: siapa admin yang buka kas historis
ALTER TABLE cash_drawer ADD COLUMN created_by INT NULL;

-- Audit: siapa admin yang input transaksi historis
ALTER TABLE transactions ADD COLUMN created_by INT NULL;

-- Relasi transaksi → kas (untuk void rollback yang akurat)
ALTER TABLE transactions ADD COLUMN cash_drawer_id INT NULL;

-- Index untuk performa query backdate
ALTER TABLE cash_drawer ADD INDEX idx_cd_backdate (is_backdate, status);
ALTER TABLE transactions ADD INDEX idx_tx_cash_drawer (cash_drawer_id);

-- -------------------------------------------------------------
-- Menu: Grup "Historis" (setelah Penjualan, order_index = 3)
-- Geser group lain ke bawah agar Historis masuk di posisi 3
-- -------------------------------------------------------------
UPDATE menus SET order_index = order_index + 1 WHERE parent_id IS NULL AND order_index >= 3;

INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('historis', 'Historis', 'CalendarClock', NULL, 3);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'historis.kas_historis', 'Kas Historis', 'CalendarClock', '/backdate/cash', 1
FROM menus m WHERE m.key_name = 'historis';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'historis.kasir_historis', 'Kasir Historis', 'ShoppingCart', '/backdate/cashier', 2
FROM menus m WHERE m.key_name = 'historis';

-- -------------------------------------------------------------
-- Route Registry — daftarkan path FE baru
-- -------------------------------------------------------------
INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/backdate/cash', 'Kas Historis'),
    ('/backdate/cashier', 'Kasir Historis');

-- -------------------------------------------------------------
-- Role Menu Access — hanya Admin & Owner yang bisa akses
-- (Kasir TIDAK mendapat akses menu Historis)
-- -------------------------------------------------------------

-- Admin: full CRUD
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
CROSS JOIN menus m
WHERE r.name = 'admin'
  AND m.key_name IN ('historis', 'historis.kas_historis', 'historis.kasir_historis');

-- Owner: full CRUD
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
CROSS JOIN menus m
WHERE r.name = 'owner'
  AND m.key_name IN ('historis', 'historis.kas_historis', 'historis.kasir_historis');
