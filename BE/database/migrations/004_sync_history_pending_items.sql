-- Kolom pending_items menyimpan jumlah item non-transaksi yang diterima & diantrekan
-- ke sync_queue tapi belum benar-benar diterapkan ke tabel target (lihat komentar
-- di sync_service.go PushSync). Sebelumnya jumlah ini dihitung di kode tapi tidak
-- pernah disimpan, sehingga riwayat sync bisa berstatus "partial" tanpa penjelasan
-- angka apa pun (total_items > synced+conflict+failed, selisihnya tidak terlihat).
ALTER TABLE sync_history
    ADD COLUMN pending_items INT DEFAULT 0 AFTER failed_items;
