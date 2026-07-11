-- Menu Backup & Restore, sebelumnya hanya ada endpoint backend tanpa representasi
-- menu/permission apa pun (akses dihardcode lewat RoleMiddleware). Sekarang dipindah
-- ke pola permission-menu standar seperti modul lain (role_menu_access).
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.backup', 'Backup & Restore', 'DatabaseBackup', '/settings/backup', 7
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/settings/backup', 'Backup & Restore');

-- OWNER: akses penuh (view, create, delete/restore) — mengikuti pola blanket
-- "owner full access" yang sudah berjalan untuk semua menu lain.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r, menus m
WHERE r.name = 'owner' AND m.key_name = 'sistem.backup';

-- ADMIN: boleh lihat & buat backup rutin, TIDAK boleh restore (can_delete=0) —
-- mengikuti pola default admin (can_delete owner-only) yang sudah berlaku di
-- seluruh menu lain di aplikasi ini.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 0
FROM roles r, menus m
WHERE r.name = 'admin' AND m.key_name = 'sistem.backup';
