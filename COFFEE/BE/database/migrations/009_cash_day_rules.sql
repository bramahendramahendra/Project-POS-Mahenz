-- =============================================================
-- Migration 009: Kas dan batas hari untuk pesanan terbuka (tahap 4.9, keputusan 34, 42, 45, 110-113)
-- =============================================================

-- Pelaku yang menyelesaikan pembayaran pesanan terbuka (kasir penerima, atau admin/owner yang menyelesaikan
-- pesanan tertunda atas nama kasir, keputusan 113). Untuk pesanan tertunda, user_id = pemilik kas tertunda.
ALTER TABLE transactions ADD COLUMN settled_by int DEFAULT NULL AFTER cancel_reason;
