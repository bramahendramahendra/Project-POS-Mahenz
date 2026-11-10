-- =============================================================
-- Migration 007: Ensure 'balance' payment method exists
-- Fix: migration 006 mungkin gagal INSERT karena ALTER error sebelumnya
-- =============================================================

INSERT IGNORE INTO payment_methods (code, label, is_active, sort_order)
VALUES ('balance', 'Saldo', 1, 6);
