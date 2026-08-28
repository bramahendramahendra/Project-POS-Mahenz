-- =============================================================
-- 009_stock_reconciliation_menu.sql
-- Menambahkan menu "Rekonsiliasi Stok" (stok lama vs baru + koreksi manual)
-- di grup Pelaporan. HANYA untuk role admin.
--
-- Catatan: file ini dijalankan otomatis oleh RunMigrations (database/migrate.go)
-- pada boot BE berikutnya, dan hanya sekali (tercatat di migrations_history).
-- Semua INSERT pakai IGNORE supaya aman kalau dijalankan ulang di DB yang
-- sudah pernah punya baris menu-nya.
-- =============================================================

-- 1) Menu (child dari grup 'pelaporan')
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.rekonsiliasi_stok', 'Rekonsiliasi Stok', 'ClipboardCheck', '/reports/stock-reconciliation', 6
FROM menus m WHERE m.key_name = 'pelaporan';

-- 2) Whitelist path untuk dropdown di Manajemen Menu
INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/reports/stock-reconciliation', 'Rekonsiliasi Stok');

-- 3) Hak akses HANYA untuk admin (view + edit untuk koreksi manual).
--    Tidak ada grant untuk owner/kasir => role lain tidak melihat menu ini.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 1, 0
FROM roles r
JOIN menus m ON m.key_name = 'pelaporan.rekonsiliasi_stok'
WHERE r.name = 'admin';
