-- Fase 8 (docs/RENCANA_PERBAIKAN_STOK_PRESISI.md) -- Bersih-bersih.
-- Menghapus products.stock dan products.reserved_qty (kolom cache turunan
-- lama dari sebelum stok dipindah ke product_packages, Fase 1-5).
--
-- Sebelum drop, stok lama diselamatkan ke products_stock_backup supaya angka
-- stok prod lama tidak hilang dan bisa dipakai backfill_stock_restore untuk
-- validasi silang. Dulu langkah ini manual sehingga rawan lupa atau kebalik
-- urutan di production, kini jadi bagian migrasi supaya otomatis dan aman.
-- CREATE TABLE IF NOT EXISTS membuatnya idempotent bila migrasi ter-retry.
CREATE TABLE IF NOT EXISTS products_stock_backup AS
SELECT id, stock, reserved_qty FROM products;

ALTER TABLE products
  DROP COLUMN stock,
  DROP COLUMN reserved_qty;
