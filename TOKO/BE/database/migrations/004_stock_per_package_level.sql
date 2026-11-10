-- 004_stock_per_package_level.sql
-- Fase 1 dari docs/RENCANA_PERBAIKAN_STOK_PRESISI.md
-- Memindahkan penyimpanan stok dari products.stock (1 angka desimal di
-- satuan anchor) ke product_packages.stock (per-level, satu baris per
-- satuan). products.stock/reserved_qty TIDAK dihapus di migrasi ini —
-- tetap dipertahankan sebagai cache turunan sampai semua jalur baca/tulis
-- dipindah (Fase 4-5), supaya rollout bisa bertahap.

ALTER TABLE product_packages
  ADD COLUMN stock DECIMAL(15,3) NOT NULL DEFAULT 0,
  ADD COLUMN reserved_qty DECIMAL(15,3) NOT NULL DEFAULT 0,
  ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1;

-- backstop DB-level: cegah stok/reserved negatif (Temuan #C)
ALTER TABLE product_packages
  ADD CONSTRAINT chk_product_packages_stock_nonneg CHECK (stock >= 0),
  ADD CONSTRAINT chk_product_packages_reserved_nonneg CHECK (reserved_qty >= 0);

ALTER TABLE products
  ADD CONSTRAINT chk_products_min_stock_nonneg CHECK (min_stock >= 0);

-- backstop DB-level: jamin cuma 1 anchor (is_default=1) per produk (celah #B)
-- Unique constraint biasa (product_id, is_default) TIDAK BISA dipakai --
-- mayoritas produk sudah punya 2+ baris is_default=0, constraint biasa akan
-- mengira semua baris is_default=0 itu saling bentrok dan ALTER TABLE
-- langsung gagal di data yang sudah ada. Wajib pakai generated column
-- (MySQL abaikan NULL di unique index, jadi cuma baris is_default=1 yang
-- benar-benar dijaga unik).
-- CATATAN (ditemukan saat eksekusi Fase 1): STORED gagal dengan error 1215
-- "Cannot add foreign key constraint" karena product_packages punya FK ke
-- dirinya sendiri (ref_package_id -> product_packages.id) -- MySQL perlu
-- rebuild tabel penuh (ALGORITHM=COPY) untuk STORED, dan itu konflik dengan
-- self-referencing FK. Diuji & terbukti VIRTUAL tidak kena masalah ini
-- (dihitung on-the-fly, bukan disimpan fisik) -- fungsinya sama untuk
-- kebutuhan unique index ini, jadi dipakai VIRTUAL:
ALTER TABLE product_packages
  ADD COLUMN is_default_flag INT GENERATED ALWAYS AS (IF(is_default = 1, product_id, NULL)) VIRTUAL;

ALTER TABLE product_packages
  ADD UNIQUE KEY uq_product_default (is_default_flag);

-- penanda satuan kontinu (kg/gram/liter) vs diskrit (celah #9)
ALTER TABLE units
  ADD COLUMN is_continuous TINYINT(1) NOT NULL DEFAULT 0;

-- penanda produk yang dilewati/gagal backfill Fase 3, butuh tinjauan manual (celah #20)
ALTER TABLE products
  ADD COLUMN needs_stock_review TINYINT(1) NOT NULL DEFAULT 0,
  ADD COLUMN stock_review_note TEXT NULL;

-- retur supplier perlu tahu level satuan yang mana buat reserve/release (celah #14).
-- Backfill data lama: isi dari purchase_items.package_id lewat purchase_item_id
-- kalau relasinya ada -- kalau tidak ada, biarkan NULL dan tangani sebagai kasus
-- "level tidak diketahui" di Fase 3/4 (bukan diasumsikan anchor).
ALTER TABLE supplier_return_items
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_supplier_return_items_package FOREIGN KEY (package_id) REFERENCES product_packages(id);

-- write-off kadaluarsa punya masalah sama persis (celah #14b).
-- Backfill sama: isi dari purchase_items.package_id lewat purchase_item_id.
ALTER TABLE product_expiry_batches
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_expiry_batches_package FOREIGN KEY (package_id) REFERENCES product_packages(id);

-- audit trail per-level (celah #16).
-- Diisi otomatis oleh fungsi terpusat (Fase 2) mulai berlaku, data lama biarkan NULL.
ALTER TABLE stock_mutations
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_stock_mutations_package FOREIGN KEY (package_id) REFERENCES product_packages(id);

-- CATATAN PENTING (celah #15): purchase_items.package_id sendiri 100% NULL
-- di data live sekarang -- join di atas (celah #14/#14b) akan mentok NULL
-- kalau kolom sumber ini belum diperbaiki dulu. Backfill-nya BUKAN sekadar
-- ALTER TABLE (kolomnya sudah ada dari skema awal), tapi perlu skrip
-- terpisah yang mencocokkan unit (teks) dan conversion_qty tersimpan ke
-- product_packages yang sesuai -- lihat Fase 3 dan Fase 4 poin 1.

-- Isi is_continuous untuk satuan yang dikonfirmasi kontinu (Fase 0):
-- Kilogram, Gram, Galon, Gelas. Sisanya (Pieces, Pack, Kardus, Karton,
-- Batang, Slop, Pres, Bungkus, Renteng, Botol, Tabung, Sak, Pouch, Sachet,
-- Kaleng, Ikat, Krak, Cup) tetap diskrit (default 0).
UPDATE units SET is_continuous = 1 WHERE name IN ('Kilogram', 'Gram', 'Galon', 'Gelas');
