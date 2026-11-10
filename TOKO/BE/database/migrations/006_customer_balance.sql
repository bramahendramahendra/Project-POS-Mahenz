-- =============================================================
-- Migration 006: Customer Balance (Saldo Pelanggan)
-- Fitur deposit/titipan uang pelanggan untuk belanja berikutnya.
-- =============================================================

-- 1. Tambah kolom balance di customers
ALTER TABLE customers ADD COLUMN balance DECIMAL(15,2) NOT NULL DEFAULT 0;

-- 2. Tambah kolom balance_used di transactions
ALTER TABLE transactions ADD COLUMN balance_used DECIMAL(15,2) NOT NULL DEFAULT 0;

-- 3. Tambah payment method 'balance'
INSERT IGNORE INTO payment_methods (code, label, is_active, sort_order)
VALUES ('balance', 'Saldo', 1, 6);

-- 4. Tabel mutasi saldo pelanggan
CREATE TABLE IF NOT EXISTS customer_balance_mutations (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    customer_id     INT           NOT NULL,
    amount          DECIMAL(15,2) NOT NULL,
    balance_after   DECIMAL(15,2) NOT NULL,
    type            ENUM('topup', 'usage', 'refund', 'adjustment') NOT NULL,
    reference_type  VARCHAR(50)   NULL,
    reference_id    INT           NULL,
    notes           TEXT          NULL,
    user_id         INT           NOT NULL,
    created_at      DATETIME      DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    INDEX idx_cbm_customer (customer_id),
    INDEX idx_cbm_type (type),
    INDEX idx_cbm_ref (reference_type, reference_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
