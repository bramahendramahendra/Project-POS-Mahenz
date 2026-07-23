-- -------------------------------------------------------------
-- Produk Expired — batch tanggal expired opsional per baris item Pembelian.
-- Lihat catatan desain: warning dihitung MURNI dari expired_date + status
-- (bukan dari estimasi FEFO sisa stok), karena di toko self-service sistem
-- tidak bisa tahu unit fisik mana yang benar-benar terjual duluan. Warning
-- baru hilang kalau staf konfirmasi manual setelah cek fisik rak.
-- -------------------------------------------------------------

CREATE TABLE IF NOT EXISTS product_expiry_batches (
    id               INT AUTO_INCREMENT PRIMARY KEY,
    product_id       INT           NULL,
    purchase_item_id INT           NULL,
    qty              DECIMAL(15,3) NOT NULL,
    expired_date     DATE          NOT NULL,
    status           ENUM('active', 'cleared', 'written_off') NOT NULL DEFAULT 'active',
    resolved_by      INT           NULL,
    resolved_at      DATETIME      NULL,
    notes            TEXT          NULL,
    created_at       DATETIME      DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id)       REFERENCES products(id)       ON DELETE SET NULL,
    FOREIGN KEY (purchase_item_id) REFERENCES purchase_items(id) ON DELETE CASCADE,
    FOREIGN KEY (resolved_by)      REFERENCES users(id)          ON DELETE SET NULL,
    INDEX idx_expiry_batches_product_status (product_id, status, expired_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tambah jenis mutasi "expired" untuk write-off stok yang dimusnahkan karena
-- sudah lewat tanggal expired (dikonfirmasi manual oleh staf setelah cek fisik).
ALTER TABLE stock_mutations
    MODIFY COLUMN mutation_type ENUM('in','out','adjustment','void','return','void_purchase','expired') NOT NULL;
