# Audit: Modul Procurement (purchases / returns / suppliers)

Scope: `FE/src/features/procurement/**`, `BE/domain/supplier_purchase/**`, `BE/domain/supplier_return/**`, `BE/domain/supplier/**`

## Temuan

### 1. [KRITIS] Update purchase order menghitung total salah & menghilangkan data pembayaran
- **File**: `BE/domain/supplier_purchase/repo/purchase_repo.go` — `Update` (baris 321-331)
- **Masalah**: Menghitung ulang `total_amount` dari subtotal item mentah **tanpa mengurangi diskon**, dan diam-diam menghilangkan `discount_amount`, `payment_status`, `paid_amount`, `payment_method` dari statement UPDATE — padahal `dto.UpdateRequest` dan FE (`PurchaseFormModal.tsx` baris 206-232) mengirimkannya.
- **Pembanding**: `Create` (baris 211-231) benar: `subtotal - req.DiscountAmount`.
- **Dampak**: Edit PO menghasilkan `total_amount` terlalu tinggi dan kehilangan perubahan diskon/pembayaran.

### 2. [Tinggi] Purchase create/update baca stok tanpa row locking
- **File**: `BE/domain/supplier_purchase/repo/purchase_repo.go` (sekitar baris 273 untuk Create, dan di Update)
- **Pembanding**: `BE/domain/supplier_return/repo/supplier_return_repo.go` (baris 30, 34) pakai `FOR UPDATE` untuk perubahan stok.
- **Dampak**: Transaksi purchase konkuren pada produk yang sama bisa race pada baca/update stok (lost update).

### 3. [Medium] Aturan hapus return berbeda antara FE dan BE
- **File BE**: `supplier_return_service.go` (baris 165-178) — boleh hapus return status apa pun kecuali `approved` (termasuk `rejected`)
- **File FE**: `FE/src/features/procurement/returns/components/ReturnTableColumns.tsx` (baris 85) — hanya tampilkan tombol delete jika status `pending`
- **Dampak**: Return berstatus `rejected` bisa dihapus lewat API langsung tapi UI tidak pernah menyediakan aksi itu.

### 4. [Medium] Input uang tidak konsisten di PaymentModal
- **File**: `FE/src/features/procurement/purchases/components/PaymentModal.tsx` (baris 168-176) — pakai `<Input type="number">` polos dengan `valueAsNumber: true`
- **Pembanding**: `PurchaseFormModal.tsx` (baris 386-397, 458-464) pakai `RupiahInput` via `Controller`.
- **Dampak**: Mengosongkan field amount menghasilkan `NaN` tanpa guard/format.

### 5. [Medium] Validasi tanggal return hanya di BE
- **File FE**: `FE/src/features/procurement/returns/returns.schema.ts` (baris 7-10) — hanya cek `return_date <= today`
- **File BE**: `supplier_return_service.go` (baris 76-98) — tambahan cek `return_date >= purchase_date`, tidak ada padanan di FE.
- **Dampak**: User baru tahu error setelah submit (BE error), bukan validasi inline seperti field lain.

### 6. [Low] `Update...Payload` di-duplikasi manual
- **File**: `purchases.types.ts` (baris 85-96), `suppliers.types.ts` (baris 68-75)
- **Catatan**: Kemungkinan menjadi akar penyebab bug #1 di atas.

### 7. [Low] Kolom `sortable` tidak konsisten & default pagination berbeda
- **File**: `ReturnTableColumns.tsx` (baris 42-47 vs 59-63); `PurchaseTable.tsx` baris 26 set `{initialPageSize: 10}` eksplisit vs `ReturnTable.tsx` baris 24 pakai default implisit.

### 8. [Low] Duplikasi aturan "tidak bisa edit setelah dibayar" tanpa sumber tunggal
- **File**: `PurchaseTableColumns.tsx` (baris 99) & `purchase_service.go` (baris 147-149) — keduanya cek `paid_amount === 0` secara independen, tanpa constant/kontrak bersama.

## Yang Sudah Baik
Modul Suppliers (CRUD service/repo/DTO) solid dan konsisten dengan validasi FE — tidak ada temuan berarti.
