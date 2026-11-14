-- =============================================================
-- Migration 007: Pesanan Coffee (tahap 4.7a, keputusan 25, 40, 41, 56, 95-101)
-- Pesanan = tabel transactions salinan Toko yang diperluas, agar kas, utang, metode bayar,
-- laporan, dan cetak struk tetap dipakai.
-- =============================================================

-- Kepala pesanan: jenis pesanan, meja (dine-in), nomor antrean harian, nama pemesan (opsional).
-- Nama meja disimpan juga agar struk lama tetap benar walau meja diganti nama.
ALTER TABLE transactions ADD COLUMN order_type varchar(12) DEFAULT NULL AFTER status;
ALTER TABLE transactions ADD COLUMN table_id int DEFAULT NULL AFTER order_type;
ALTER TABLE transactions ADD COLUMN table_name varchar(50) DEFAULT NULL AFTER table_id;
ALTER TABLE transactions ADD COLUMN queue_date date DEFAULT NULL AFTER table_name;
ALTER TABLE transactions ADD COLUMN queue_number int DEFAULT NULL AFTER queue_date;
ALTER TABLE transactions ADD COLUMN order_name varchar(50) DEFAULT NULL AFTER queue_number;
ALTER TABLE transactions ADD KEY idx_transactions_queue (queue_date, queue_number);
ALTER TABLE transactions ADD KEY idx_transactions_table (table_id);
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_table FOREIGN KEY (table_id) REFERENCES dining_tables (id) ON DELETE RESTRICT;

-- Item pesanan: menu dan kombinasi varian yang dijual, nama kombinasi saat dijual, catatan.
-- product_id hanya terisi untuk menu barang jadi (stok dipotong). Kombinasi yang kelak dihapus
-- dari form menu tidak menghapus riwayat (SET NULL), nama tetap tersimpan.
ALTER TABLE transaction_items ADD COLUMN menu_item_id int DEFAULT NULL AFTER transaction_id;
ALTER TABLE transaction_items ADD COLUMN variant_id int DEFAULT NULL AFTER menu_item_id;
ALTER TABLE transaction_items ADD COLUMN variant_name varchar(100) DEFAULT NULL AFTER product_name;
ALTER TABLE transaction_items ADD COLUMN notes varchar(200) DEFAULT NULL AFTER variant_name;
ALTER TABLE transaction_items ADD KEY idx_transaction_items_menu (menu_item_id);
ALTER TABLE transaction_items ADD CONSTRAINT fk_transaction_items_menu FOREIGN KEY (menu_item_id) REFERENCES menu_items (id) ON DELETE RESTRICT;
ALTER TABLE transaction_items ADD CONSTRAINT fk_transaction_items_variant FOREIGN KEY (variant_id) REFERENCES menu_item_variants (id) ON DELETE SET NULL;

-- Tambahan yang dipilih per item, dengan nama dan harga saat dijual.
CREATE TABLE IF NOT EXISTS `transaction_item_modifiers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `transaction_item_id` int NOT NULL,
  `option_id` int DEFAULT NULL,
  `group_name` varchar(100) NOT NULL,
  `option_name` varchar(100) NOT NULL,
  `price` decimal(15,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `idx_tim_item` (`transaction_item_id`),
  KEY `idx_tim_option` (`option_id`),
  CONSTRAINT `fk_tim_item` FOREIGN KEY (`transaction_item_id`) REFERENCES `transaction_items` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tim_option` FOREIGN KEY (`option_id`) REFERENCES `modifier_options` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Barang jadi yang dibuang (void atau batal pesanan, keputusan 44, 100): jenis mutasi stok baru.
ALTER TABLE stock_mutations MODIFY COLUMN mutation_type enum('in','out','adjustment','void','return','void_purchase','expired','waste') NOT NULL;
