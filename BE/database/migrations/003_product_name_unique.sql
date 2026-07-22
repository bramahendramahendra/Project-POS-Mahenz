-- =============================================================
-- Migration 003: Cegah nama produk duplikat
-- Nama produk dibandingkan case-insensitive & trimmed lewat kolom
-- generated name_normalized, di-unique-kan agar aman dari race
-- condition (bukan cuma dicek di service layer).
-- =============================================================

ALTER TABLE products
    ADD COLUMN name_normalized VARCHAR(200)
        GENERATED ALWAYS AS (LOWER(TRIM(name))) STORED,
    ADD UNIQUE KEY unique_product_name (name_normalized);
