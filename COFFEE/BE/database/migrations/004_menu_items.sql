-- =============================================================
-- Migration 004: Menu jualan (tahap 4.4, keputusan 33, 43, 53-55, 85-89)
--
-- Nama tabel memakai awalan menu_ / modifier_ agar tidak bentrok dengan tabel menus
-- (menu navigasi dan izin, milik salinan Toko).
-- =============================================================

CREATE TABLE IF NOT EXISTS `menu_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `sort_order` int NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_menu_categories_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Daftar ukuran varian, dikelola admin (keputusan 85).
CREATE TABLE IF NOT EXISTS `menu_sizes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `sort_order` int NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_menu_sizes_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Menu jualan. stock_mode: recipe = racikan (stok bahan lewat pemakaian manual),
-- direct = langsung / barang jadi (stok di tahap 4.5, tanpa varian, keputusan 43).
-- has_variants = 0: satu harga, disimpan sebagai satu baris varian tanpa ukuran/suhu.
CREATE TABLE IF NOT EXISTS `menu_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `stock_mode` varchar(10) NOT NULL DEFAULT 'recipe',
  `has_variants` tinyint(1) NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `is_sold_out` tinyint(1) NOT NULL DEFAULT '0',
  `created_by` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_menu_items_name` (`name`),
  KEY `idx_menu_items_category` (`category_id`),
  CONSTRAINT `fk_menu_items_category` FOREIGN KEY (`category_id`) REFERENCES `menu_categories` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Kombinasi varian dan harganya (keputusan 53). size_id dan temperature boleh kosong:
-- menu tanpa varian punya tepat satu baris dengan keduanya kosong. temperature: hot / iced.
CREATE TABLE IF NOT EXISTS `menu_item_variants` (
  `id` int NOT NULL AUTO_INCREMENT,
  `menu_item_id` int NOT NULL,
  `size_id` int DEFAULT NULL,
  `temperature` varchar(10) DEFAULT NULL,
  `price` decimal(15,2) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_menu_item_variants_item` (`menu_item_id`),
  KEY `idx_menu_item_variants_size` (`size_id`),
  CONSTRAINT `fk_menu_item_variants_item` FOREIGN KEY (`menu_item_id`) REFERENCES `menu_items` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_menu_item_variants_size` FOREIGN KEY (`size_id`) REFERENCES `menu_sizes` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Resep per kombinasi varian: baris bahan + takaran teks (keputusan 86). Tidak mengurangi stok.
CREATE TABLE IF NOT EXISTS `menu_recipe_lines` (
  `id` int NOT NULL AUTO_INCREMENT,
  `variant_id` int NOT NULL,
  `ingredient_id` int NOT NULL,
  `amount` varchar(50) NOT NULL,
  `sort_order` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_menu_recipe_lines` (`variant_id`, `ingredient_id`),
  KEY `idx_menu_recipe_lines_ingredient` (`ingredient_id`),
  CONSTRAINT `fk_menu_recipe_lines_variant` FOREIGN KEY (`variant_id`) REFERENCES `menu_item_variants` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_menu_recipe_lines_ingredient` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Kelompok tambahan (keputusan 54, 87). selection_type: single = pilih satu, multiple = boleh banyak.
CREATE TABLE IF NOT EXISTS `modifier_groups` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `selection_type` varchar(10) NOT NULL DEFAULT 'single',
  `is_required` tinyint(1) NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_modifier_groups_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Pilihan dalam kelompok, harga 0 = gratis, paling banyak satu pilihan bawaan per kelompok.
CREATE TABLE IF NOT EXISTS `modifier_options` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `price` decimal(15,2) NOT NULL DEFAULT '0.00',
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `sort_order` int NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_modifier_options_group` (`group_id`),
  CONSTRAINT `fk_modifier_options_group` FOREIGN KEY (`group_id`) REFERENCES `modifier_groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Kelompok tambahan yang dipasang ke menu.
CREATE TABLE IF NOT EXISTS `menu_item_modifier_groups` (
  `menu_item_id` int NOT NULL,
  `group_id` int NOT NULL,
  `sort_order` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`menu_item_id`, `group_id`),
  KEY `idx_menu_item_modifier_groups_group` (`group_id`),
  CONSTRAINT `fk_mimg_item` FOREIGN KEY (`menu_item_id`) REFERENCES `menu_items` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_mimg_group` FOREIGN KEY (`group_id`) REFERENCES `modifier_groups` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- -------------------------------------------------------------
-- Menu navigasi: grup "Menu" setelah Penjualan (grup berikutnya digeser satu).
-- Kunci memakai awalan menu_jual agar tidak tertukar dengan menu navigasi.
-- -------------------------------------------------------------
UPDATE menus SET order_index = order_index + 1 WHERE parent_id IS NULL AND order_index >= 3;

INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('menu_jual', 'Menu', 'Coffee', NULL, 3);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'menu_jual.daftar', 'Daftar Menu', 'Coffee', '/menu-items', 1
FROM menus m WHERE m.key_name = 'menu_jual';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'menu_jual.kategori', 'Kategori Menu', 'Tags', '/menu-items/categories', 2
FROM menus m WHERE m.key_name = 'menu_jual';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'menu_jual.tambahan', 'Kelompok Tambahan', 'ListPlus', '/menu-items/modifiers', 3
FROM menus m WHERE m.key_name = 'menu_jual';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'menu_jual.ukuran', 'Ukuran', 'Ruler', '/menu-items/sizes', 4
FROM menus m WHERE m.key_name = 'menu_jual';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/menu-items', 'Daftar Menu'),
    ('/menu-items/categories', 'Kategori Menu'),
    ('/menu-items/modifiers', 'Kelompok Tambahan'),
    ('/menu-items/sizes', 'Ukuran');

-- -------------------------------------------------------------
-- Izin (keputusan 88): owner dan admin penuh. Kasir hanya melihat Daftar Menu
-- (lihat = boleh menandai habis/tersedia).
-- -------------------------------------------------------------
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name IN ('menu_jual', 'menu_jual.daftar', 'menu_jual.kategori', 'menu_jual.tambahan', 'menu_jual.ukuran')
WHERE r.name IN ('owner', 'admin');

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name IN ('menu_jual', 'menu_jual.daftar')
WHERE r.name = 'kasir';
