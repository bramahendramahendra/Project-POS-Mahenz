-- Fase 8 (docs/RENCANA_PERBAIKAN_STOK_PRESISI.md) -- Bersih-bersih.
-- Menghapus products.stock dan products.reserved_qty: kolom cache turunan
-- lama dari sebelum stok dipindah ke product_packages (Fase 1-5). Semua
-- jalur baca/tulis BE sudah dipindah sepenuhnya ke product_packages sejak
-- Fase 5, dan ApplyStockDelta tidak lagi menulis kolom ini sejak Fase 8 --
-- lihat grep menyeluruh yang dilakukan sebelum migrasi ini ditulis.
--
-- IRREVERSIBLE tanpa restore dari backup -- wajib backup DB dulu sebelum
-- menjalankan ini (lihat prompt Fase 8 di dokumen).
--
-- Tidak ada CHECK constraint atau FOREIGN KEY yang menempel di kolom ini
-- (dikonfirmasi lewat SHOW CREATE TABLE products sebelum migrasi ditulis),
-- jadi DROP COLUMN polos aman dijalankan.
ALTER TABLE products
  DROP COLUMN stock,
  DROP COLUMN reserved_qty;
