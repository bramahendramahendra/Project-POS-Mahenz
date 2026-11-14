-- =============================================================
-- Migration 005: Pembelian dua jenis dan tautan barang jadi (tahap 4.5, keputusan 91-93)
-- =============================================================

-- Baris pembelian: item_type product = barang jadi (produk + kemasan, seperti Toko),
-- ingredient = bahan baku (jumlah dalam satuan stok bahan, desimal, keputusan 92).
ALTER TABLE purchase_items ADD COLUMN item_type varchar(12) NOT NULL DEFAULT 'product' AFTER purchase_id;
ALTER TABLE purchase_items ADD COLUMN ingredient_id int DEFAULT NULL AFTER product_id;
ALTER TABLE purchase_items ADD KEY idx_purchase_items_ingredient (ingredient_id);
ALTER TABLE purchase_items ADD CONSTRAINT fk_purchase_items_ingredient FOREIGN KEY (ingredient_id) REFERENCES ingredients (id) ON DELETE RESTRICT;

-- Menu bermode barang jadi ditautkan ke tepat satu produk (keputusan 91).
-- Satu produk hanya untuk satu menu (kunci unik). Menu racikan tidak memakai produk.
ALTER TABLE menu_items ADD COLUMN product_id int DEFAULT NULL AFTER stock_mode;
ALTER TABLE menu_items ADD UNIQUE KEY uq_menu_items_product (product_id);
ALTER TABLE menu_items ADD CONSTRAINT fk_menu_items_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT;

-- Grup menu navigasi "Produk" (salinan Toko) menjadi katalog barang jadi.
UPDATE menus SET label = 'Barang Jadi' WHERE key_name = 'produk';
