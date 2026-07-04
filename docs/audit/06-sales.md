# Audit: Modul Sales (cashier / transactions) — Checkout Flow

Scope: `FE/src/features/sales/**`, `BE/domain/transaction/**`, `BE/routes/segment/transaction_routes.go`

Modul paling kritis karena ini alur checkout inti POS.

## Temuan

### 1. [KRITIS] BE tidak pernah menghitung ulang uang — percaya total dari client
- **File**: `BE/domain/transaction/repo/transaction_repo.go` (baris 189-196)
- **Masalah**: `req.Subtotal, req.Discount, req.Tax, req.TotalAmount` dan per-item `Price/Subtotal/DiscountItem` di-insert langsung dari payload client. `dto_transaction.go` (baris 38-51) hanya validasi `min=0`/`required`, tidak ada cross-check aritmatika.
- **Pembanding**: `BE/domain/supplier_purchase/repo/purchase_repo.go` (baris 211-218) menghitung ulang subtotal dari `PurchasePrice * Quantity` server-side.
- **Dampak**: Request yang di-crafting bisa submit total transaksi sembarangan sementara stok tetap dikurangi sesuai quantity asli — **celah manipulasi harga jual**.

### 2. [KRITIS] Endpoint create transaksi tanpa permission middleware
- **File**: `BE/routes/segment/transaction_routes.go` (baris 24-27) — create tanpa gate sama sekali (void hanya di-gate `can_delete` generik)
- **Pembanding**: `supplier_purchase_routes.go` (baris 29-32) men-gate create dengan `can_create`.
- **Dampak**: Siapa pun yang login bisa membuat transaksi penjualan tanpa izin apa pun.

### 3. [Tinggi] `shift_id` tidak divalidasi terhadap user yang login
- **File**: `BE/domain/transaction/service/transaction_service.go` (baris 30-48)
- **Masalah**: `req.ShiftID` di-insert apa adanya, sementara cash drawer yang di-kredit adalah milik `userID` pemanggil — tidak ada pengecekan shift itu benar milik user tersebut.
- **Dampak**: User bisa submit transaksi dengan `shift_id` milik user lain, mengkredit/mempengaruhi cash drawer yang salah.

### 4. [Tinggi] Void transaksi tidak membalikkan total penjualan cash-drawer
- **File**: `BE/domain/transaction/repo/transaction_repo.go` — `Void` (baris 309-356)
- **Masalah**: Void mengembalikan stok dan membatalkan receivable, tapi tidak pernah memanggil balik `cashDrawerRepo.UpdateSales` yang dipanggil saat create (`transaction_service.go` baris 40-45).
- **Dampak**: Void transaksi cash secara permanen menaikkan `total_cash_sales` cash drawer (tidak pernah dikurangi kembali).

### 5. [Tinggi] Update cash-drawer tidak atomik dengan create/void
- **File**: `BE/domain/transaction/service/transaction_service.go` (baris 40-44)
- **Masalah**: `UpdateSales` dipanggil setelah transaksi DB commit, bukan dalam transaction yang sama.
- **Dampak**: Jika `UpdateSales` gagal setelah create commit, stok & record transaksi tetap tersimpan tapi total drawer diam-diam tidak sinkron, tanpa rollback/log error.

### 6. [Medium] Field diskon tidak ada batas atas di DTO
- **File**: `BE/domain/transaction/dto/dto_transaction.go` (baris 33, `DiscountItem`; baris 41-42, `Discount`/`Tax`)
- **Masalah**: Tidak ada validasi upper-bound, sementara FE membatasi diskon maks 100%/subtotal (`cashier.utils.ts` baris 19-21, 80).
- **Dampak**: Dikombinasikan dengan temuan #1, request yang di-crafting bisa membuat revenue negatif tanpa ditolak server.

### 7. [Medium] Model rounding tidak konsisten antara diskon cart-level dan item-level
- **File**: `FE/src/features/sales/cashier/cashier.utils.ts` — diskon cart-level (baris 19, 31) membulatkan nilai akhir; diskon persen item-level (baris 79-84) membulatkan harga satuan dulu baru dikalikan quantity.
- **Dampak**: Kesalahan pembulatan yang berbeda antara dua jalur, tanpa rekonsiliasi.

### 8. [Low] Cart di localStorage tidak di-refresh sebelum checkout
- **File**: `FE/src/features/sales/cashier/cashier.store.ts` (baris 218-227) — cart/diskon/pajak dipersist ke localStorage dengan harga yang diambil saat add-to-cart (`ProductSearch.tsx` baris 95-104), tidak pernah di-refresh.
- **Dampak**: Memperparah temuan #1 — harga bisa stale saat akhirnya di-submit.

### 9. [Low] Perlu verifikasi lanjutan: void tidak cek pembayaran parsial pada receivable
- **File**: `BE/domain/transaction/repo/transaction_repo.go` (baris 25, 350) — void set receivable ke `'void'` tanpa cek apakah sudah ada pembayaran parsial.

### 10. [Low] Tabel transaksi tanpa kolom sortable
- **File**: `FE/src/features/sales/transactions/components/TransactionTableColumns.tsx`
- **Pembanding**: `PurchaseTableColumns.tsx` (baris 32, 45, 54, 63) punya sortable.

## Prioritas Perbaikan
#1 dan #2/#3 paling parah — memungkinkan client memalsukan total penjualan dan bypass role check pada alur checkout inti. #4/#5 berikutnya karena langsung merusak akuntansi cash-drawer saat void.
