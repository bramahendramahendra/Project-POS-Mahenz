-- =============================================================
-- Migration 005: Remove Rekap Kas menu & update order_index
-- Menghapus menu keuangan.rekap_kas karena fiturnya digabung
-- ke dalam Kas Harian (keuangan.kas_harian).
-- =============================================================

-- Hapus akses role untuk menu rekap_kas
DELETE rma FROM role_menu_access rma
JOIN menus m ON rma.menu_id = m.id
WHERE m.key_name = 'keuangan.rekap_kas';

-- Hapus menu rekap_kas
DELETE FROM menus WHERE key_name = 'keuangan.rekap_kas';

-- Rapikan order_index: pengeluaran 4→3, kas_saya 5→4
UPDATE menus SET order_index = 3 WHERE key_name = 'keuangan.pengeluaran';
UPDATE menus SET order_index = 4 WHERE key_name = 'keuangan.kas_saya';
