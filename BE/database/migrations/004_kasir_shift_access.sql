-- -------------------------------------------------------------
-- Migrasi: grant akses view-only Kasir ke operasional.shift.
--
-- Kasir butuh baca daftar shift aktif (GET /shifts/active) untuk memilih
-- shift saat Buka Kas (modal OpenCashDrawerModal, dipakai baik dari
-- halaman Kas Saya maupun landing page Dashboard baru). Sebelum migrasi
-- ini, Kasir tidak punya akses menu operasional.shift sama sekali, jadi
-- endpoint itu selalu 403 dan dropdown pilih shift selalu kosong.
-- -------------------------------------------------------------

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'operasional.shift'
WHERE r.name = 'kasir';
