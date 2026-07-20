-- Menambahkan status void pada purchases (terpisah dari payment_status)
-- untuk mendukung fitur Void PO dan Tambah Item ke PO yang sudah dibayar.
ALTER TABLE purchases
    ADD COLUMN status ENUM('active', 'void') NOT NULL DEFAULT 'active' AFTER payment_method;

-- mutation_type sebelumnya tidak menampung nilai 'return' yang sudah dipakai
-- oleh fitur retur supplier (supplier_return_repo.go) — dilebarkan sekalian
-- agar konsisten, plus menambah 'void_purchase' untuk rollback stok saat void PO.
ALTER TABLE stock_mutations
    MODIFY COLUMN mutation_type ENUM('in', 'out', 'adjustment', 'void', 'return', 'void_purchase') NOT NULL;
