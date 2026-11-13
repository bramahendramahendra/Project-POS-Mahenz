-- =============================================================
-- cleanup_transaksi_cangkang_20260913.sql
--
-- Pembersihan data akibat BUG TRANSAKSI CANGKANG (header tanpa item).
-- Lihat: docs/BUG_TRANSAKSI_CANGKANG_HEADER_TANPA_ITEM.md
--
-- Yang dibersihkan:
--   1) Void 3 transaksi cangkang: id 545 (WEB-20260908-004),
--      635 (WEB-20260913-005), 636 (WEB-20260913-006).
--      Ketiganya: status 'completed', 0 item, 0 mutasi stok, 0 receivable.
--   2) Koreksi cash_drawer id 21 (8 Sep, sudah closed): closing_balance &
--      expected_balance dari 1.120.501 -> 1.083.501 (buang kontaminasi 37.000
--      dari cangkang 545 yang ikut terhitung saat closing).
--
-- Cash drawer id 26 (13 Sep, open) TIDAK diubah: kolom tersimpannya sudah benar
-- (total_cash_sales 613.500). Tampilan "kas saat ini" otomatis benar setelah
-- cangkang di-void, karena laporan mengecualikan status = 'void'.
--
-- CARA PAKAI (JALANKAN MANUAL, BUKAN MIGRASI OTOMATIS):
--   1. WAJIB backup dulu (lihat perintah di bawah, di luar file ini).
--   2. Jalankan file ini. Ia berjalan dalam SATU transaksi dan menampilkan
--      hasil verifikasi SEBELUM commit. Kalau ada baris "VERIFIKASI GAGAL",
--      JANGAN commit -- ganti COMMIT di akhir dengan ROLLBACK.
--   3. Kalau semua verifikasi OK, biarkan COMMIT berjalan.
--
-- CATATAN: Guard di setiap UPDATE memastikan hanya baris yang benar-benar
-- cangkang (0 item) yang tersentuh. Kalau kondisi tak terpenuhi, 0 baris
-- ter-update (aman, tidak merusak data lain).
-- =============================================================

START TRANSACTION;

-- ---------------------------------------------------------------
-- 1) VOID transaksi cangkang.
--    Guard: hanya yang status='completed' DAN tidak punya item sama sekali.
-- ---------------------------------------------------------------
UPDATE transactions t
SET t.status = 'void', t.updated_at = NOW()
WHERE t.id IN (545, 635, 636)
  AND t.status = 'completed'
  AND NOT EXISTS (SELECT 1 FROM transaction_items ti WHERE ti.transaction_id = t.id);

-- ---------------------------------------------------------------
-- 2) Koreksi cash_drawer 21 (closed) -- buang kontaminasi 37.000.
--    Guard: hanya kalau nilainya masih persis angka tercemar (1.120.501),
--    supaya idempotent (kalau sudah dikoreksi, 0 baris ter-update).
-- ---------------------------------------------------------------
UPDATE cash_drawer
SET closing_balance  = 1083501.00,
    expected_balance = 1083501.00,
    updated_at       = NOW()
WHERE id = 21
  AND status = 'closed'
  AND closing_balance = 1120501.00
  AND expected_balance = 1120501.00;

-- ---------------------------------------------------------------
-- VERIFIKASI (baca sebelum commit)
-- ---------------------------------------------------------------

-- 2a. Ketiga transaksi harus 'void' sekarang.
SELECT CASE WHEN COUNT(*) = 3 THEN 'OK: 3 transaksi cangkang sudah void'
            ELSE CONCAT('VERIFIKASI GAGAL: transaksi void = ', COUNT(*), ' (harusnya 3)') END AS cek_void
FROM transactions WHERE id IN (545,635,636) AND status = 'void';

-- 2b. Tidak boleh ada transaksi cangkang 'completed' tersisa di seluruh tabel.
SELECT CASE WHEN COUNT(*) = 0 THEN 'OK: tidak ada cangkang completed tersisa'
            ELSE CONCAT('VERIFIKASI GAGAL: masih ada ', COUNT(*), ' cangkang completed') END AS cek_sisa_cangkang
FROM transactions t
WHERE t.status = 'completed'
  AND NOT EXISTS (SELECT 1 FROM transaction_items ti WHERE ti.transaction_id = t.id);

-- 2c. Drawer 21 konsisten: closing_balance = opening_balance + total_cash_sales - total_expenses.
SELECT CASE WHEN closing_balance = (opening_balance + total_cash_sales - total_expenses)
            THEN 'OK: drawer 21 konsisten (closing = opening + cash_sales - expenses)'
            ELSE CONCAT('VERIFIKASI GAGAL: drawer 21 closing=', closing_balance,
                        ' != ', (opening_balance + total_cash_sales - total_expenses)) END AS cek_drawer21
FROM cash_drawer WHERE id = 21;

-- 2d. Total penjualan cash per hari (harus: 8 Sep=1.083.500, 13 Sep=613.500).
SELECT DATE(transaction_date) AS tgl, SUM(total_amount) AS total_cash_completed
FROM transactions
WHERE payment_method = 'cash' AND status = 'completed'
  AND DATE(transaction_date) IN ('2026-09-08','2026-09-13')
GROUP BY DATE(transaction_date);

-- ---------------------------------------------------------------
-- Kalau SEMUA verifikasi di atas OK -> biarkan COMMIT.
-- Kalau ADA "VERIFIKASI GAGAL"       -> ganti baris berikut jadi: ROLLBACK;
-- ---------------------------------------------------------------
COMMIT;
