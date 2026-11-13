-- =============================================================
-- Migration 003: Bahan baku (tahap 4.3, keputusan 24, 80-84)
--
-- Stok bahan = pembelian (tahap 4.5) dikurangi pemakaian yang diinput manual dan
-- barang terbuang. Tanpa opname (keputusan 84). Resep tidak mengurangi stok.
-- Salah input dibetulkan dengan membatalkan entri (keputusan 81).
-- =============================================================

-- Katalog bahan. Satu satuan stok per bahan dari master Unit (keputusan 82),
-- jumlah boleh desimal. Stok bahan baru selalu 0 (keputusan 83).
CREATE TABLE IF NOT EXISTS `ingredients` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `unit_id` int NOT NULL,
  `stock` decimal(15,3) NOT NULL DEFAULT '0.000',
  `total_used` decimal(15,3) NOT NULL DEFAULT '0.000',
  `total_wasted` decimal(15,3) NOT NULL DEFAULT '0.000',
  `min_stock` decimal(15,3) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_by` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_ingredients_name` (`name`),
  KEY `idx_ingredients_unit` (`unit_id`),
  CONSTRAINT `fk_ingredients_unit` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Riwayat stok bahan. type:
--   usage  = terpakai (qty_change negatif)
--   waste  = terbuang, alasan wajib (qty_change negatif)
--   cancel = pembatalan entri usage/waste (qty_change positif, reference_id = entri asal)
--   purchase dan purchase_void menyusul di tahap 4.5
-- Entri tidak pernah diubah atau dihapus. Entri yang dibatalkan diberi tanda cancelled_*.
CREATE TABLE IF NOT EXISTS `ingredient_mutations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ingredient_id` int NOT NULL,
  `type` varchar(20) NOT NULL,
  `qty_change` decimal(15,3) NOT NULL,
  `stock_before` decimal(15,3) NOT NULL,
  `stock_after` decimal(15,3) NOT NULL,
  `notes` varchar(255) DEFAULT NULL,
  `reference_type` varchar(30) DEFAULT NULL,
  `reference_id` int DEFAULT NULL,
  `cancelled_at` datetime DEFAULT NULL,
  `cancelled_by` int DEFAULT NULL,
  `cancel_reason` varchar(255) DEFAULT NULL,
  `created_by` int NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ingredient_mutations_ingredient` (`ingredient_id`, `created_at`),
  KEY `idx_ingredient_mutations_created` (`created_at`),
  CONSTRAINT `fk_ingredient_mutations_ingredient` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- -------------------------------------------------------------
-- Menu navigasi: grup "Bahan Baku" setelah Produk (grup berikutnya digeser satu).
-- -------------------------------------------------------------
UPDATE menus SET order_index = order_index + 1 WHERE parent_id IS NULL AND order_index >= 4;

INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('bahan', 'Bahan Baku', 'Wheat', NULL, 4);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'bahan.daftar', 'Daftar Bahan', 'Wheat', '/ingredients', 1
FROM menus m WHERE m.key_name = 'bahan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'bahan.riwayat', 'Riwayat Stok Bahan', 'History', '/ingredients/history', 2
FROM menus m WHERE m.key_name = 'bahan';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/ingredients', 'Daftar Bahan'),
    ('/ingredients/history', 'Riwayat Stok Bahan');

-- -------------------------------------------------------------
-- Izin (keputusan 80):
--   bahan.daftar  : lihat (owner, admin, kasir), kelola bahan (owner, admin)
--   bahan.riwayat : lihat + tambah = catat terpakai/terbuang (owner, admin, kasir),
--                   hapus = batalkan entri (owner, admin)
-- -------------------------------------------------------------
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name IN ('bahan', 'bahan.daftar', 'bahan.riwayat')
WHERE r.name IN ('owner', 'admin');

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name IN ('bahan', 'bahan.daftar')
WHERE r.name = 'kasir';

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'bahan.riwayat'
WHERE r.name = 'kasir';
