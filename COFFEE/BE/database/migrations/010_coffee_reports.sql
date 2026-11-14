-- =============================================================
-- Migration 010: Laporan Coffee (tahap 4.10, keputusan 114-117)
-- =============================================================

-- Dua halaman laporan baru di grup Pelaporan (keputusan 116).
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.bahan_baku', 'Bahan Baku', 'Wheat', '/reports/ingredients', 7
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.terbuang', 'Barang Terbuang', 'Trash2', '/reports/waste', 8
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/reports/ingredients', 'Laporan Bahan Baku'),
    ('/reports/waste', 'Laporan Barang Terbuang');

-- Izin sama dengan laporan lain: owner dan admin.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name IN ('pelaporan.bahan_baku', 'pelaporan.terbuang')
WHERE r.name IN ('owner', 'admin');
