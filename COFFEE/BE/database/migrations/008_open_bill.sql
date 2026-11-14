-- =============================================================
-- Migration 008: Open bill Coffee (tahap 4.8, keputusan 27-29, 40, 44, 59, 96, 103-109)
-- =============================================================

-- Status pesanan: open = pesanan terbuka (belum dibayar), cancelled = pesanan terbuka dibatalkan,
-- merged = pesanan terbuka yang digabung ke pesanan lain. Laporan dan kas tetap hanya menghitung completed.
ALTER TABLE transactions MODIFY COLUMN status enum('pending','completed','void','open','cancelled','merged') DEFAULT 'completed';

-- Gabung bill: pesanan sumber menunjuk pesanan tujuan (keputusan 107).
ALTER TABLE transactions ADD COLUMN merged_into_id int DEFAULT NULL AFTER order_name;
-- Batal pesanan terbuka: siapa, kapan, alasan (keputusan 105).
ALTER TABLE transactions ADD COLUMN cancelled_by int DEFAULT NULL AFTER merged_into_id;
ALTER TABLE transactions ADD COLUMN cancelled_at datetime DEFAULT NULL AFTER cancelled_by;
ALTER TABLE transactions ADD COLUMN cancel_reason varchar(255) DEFAULT NULL AFTER cancelled_at;
ALTER TABLE transactions ADD KEY idx_transactions_status_table (status, table_id);

-- Catatan item yang dihapus dari pesanan terbuka atau ikut batal (keputusan 105): nama dan harga saat itu,
-- pilihan stok barang jadi (return = kembali ke stok, waste = dibuang, kosong = racikan), alasan, pelaku.
CREATE TABLE IF NOT EXISTS `transaction_item_cancellations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `transaction_id` int NOT NULL,
  `item_name` varchar(200) NOT NULL,
  `variant_name` varchar(100) DEFAULT NULL,
  `quantity` decimal(15,3) NOT NULL,
  `price` decimal(15,2) NOT NULL,
  `stock_action` varchar(10) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `user_id` int DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tic_transaction` (`transaction_id`),
  CONSTRAINT `fk_tic_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
