-- =============================================================
-- Migration 011: Label metode bayar 'kredit' menjadi 'Hutang' (keputusan 135)
-- Kode tetap 'kredit' (dipakai transaksi, piutang, validasi), hanya teks tampilan yang berubah,
-- agar rincian Kas Harian/Kas Saya dan pilihan Pembelian sama dengan kasir dan struk.
-- =============================================================

UPDATE payment_methods SET label = 'Hutang' WHERE code = 'kredit';
