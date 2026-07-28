-- Fix: role_menu_access Kasir untuk menu keuangan.kas_saya kurang can_edit.
--
-- keuangan.kas_saya adalah menu self-service (buka/tutup kas milik sendiri) —
-- endpoint /cash-drawer/open|close|update-sales|update-expenses digembok ke
-- permission menu ini (lihat cash_drawer_routes.go), kepemilikan kas tetap
-- divalidasi di service layer. Tanpa can_edit, Kasir bisa buka kas tapi tidak
-- pernah bisa menutupnya sendiri.
UPDATE role_menu_access rma
JOIN roles r ON r.id = rma.role_id
JOIN menus m ON m.id = rma.menu_id
SET rma.can_edit = 1
WHERE r.name = 'kasir' AND m.key_name = 'keuangan.kas_saya';
