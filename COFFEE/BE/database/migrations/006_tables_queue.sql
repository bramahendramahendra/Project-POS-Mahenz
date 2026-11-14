-- =============================================================
-- Migration 006: Meja dan nomor antrean harian (tahap 4.6, keputusan 56, 95-98)
-- =============================================================

-- Daftar meja sederhana: nama/nomor unik, urutan tampil di kasir, aktif (keputusan 97).
CREATE TABLE IF NOT EXISTS `dining_tables` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `sort_order` int NOT NULL DEFAULT '0',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_dining_tables_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Penghitung nomor antrean: satu baris per tanggal (Asia/Jakarta), nomor terakhir yang dipakai.
-- Tanggal baru otomatis mulai dari 1 (reset 00:00, keputusan 56). Nomor tidak pernah dipakai ulang.
CREATE TABLE IF NOT EXISTS `queue_counters` (
  `queue_date` date NOT NULL,
  `last_number` int NOT NULL DEFAULT '0',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`queue_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Menu navigasi "Meja" di grup Menu setelah Ukuran (keputusan 98).
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'menu_jual.meja', 'Meja', 'Armchair', '/menu-items/tables', 5
FROM menus m WHERE m.key_name = 'menu_jual';

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/menu-items/tables', 'Meja');

-- Izin (keputusan 98): owner dan admin penuh, kasir hanya melihat.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name = 'menu_jual.meja'
WHERE r.name IN ('owner', 'admin');

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'menu_jual.meja'
WHERE r.name = 'kasir';
