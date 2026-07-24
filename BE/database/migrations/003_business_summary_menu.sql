-- -------------------------------------------------------------
-- Migrasi: menu "Ringkasan Bisnis" (isi Dashboard lama dipindah ke sini,
-- khusus Owner/Admin) + grant akses Kasir ke beranda.dashboard (landing
-- page operasional baru).
-- -------------------------------------------------------------

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.ringkasan_bisnis', 'Ringkasan Bisnis', 'BarChart3', '/reports/business-summary', 5
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/reports/business-summary', 'Ringkasan Bisnis');

-- CATATAN: role_menu_access adalah tabel materialized (diisi sekali saat seed),
-- BUKAN view dinamis — wildcard query Owner/Admin di 002_seed_data.sql cuma berlaku
-- untuk menu yang sudah ada SAAT seed dijalankan. Menu baru tetap wajib digrant eksplisit.

-- OWNER: full akses ke menu baru (konsisten dengan pola Owner di 002_seed_data.sql).
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name = 'pelaporan.ringkasan_bisnis'
WHERE r.name = 'owner';

-- ADMIN: full akses ke menu baru (konsisten dengan pola Admin di 002_seed_data.sql).
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name = 'pelaporan.ringkasan_bisnis'
WHERE r.name = 'admin';

-- KASIR: tambah akses view-only ke beranda.dashboard (landing page operasional baru),
-- pola sama seperti grant sistem.profil_toko ke Kasir di 002_seed_data.sql.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'beranda.dashboard'
WHERE r.name = 'kasir';
