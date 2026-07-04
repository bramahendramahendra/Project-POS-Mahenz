# Audit: Modul Customers (customers / receivables)

Scope: `FE/src/features/customers/**`, `BE/domain/customer/**`, `BE/domain/receivable/**`

## Temuan

### 1. [KRITIS] Fitur pembayaran piutang rusak end-to-end
- **File**: `BE/domain/receivable/dto/dto_receivable.go` (`PayRequest`, sekitar baris 60-66), `FE/src/features/customers/receivables/receivables.types.ts` (`CreatePaymentPayload`, baris 32-36), `FE/src/features/customers/receivables/components/PaymentRecordModal.tsx` (baris 60-72), `BE/domain/receivable/repo/receivable_repo.go` (`CreatePayment`, baris ~99)
- **Masalah**:
  - FE mengirim `{ amount, payment_date, notes }`.
  - BE `PayRequest` **mewajibkan** `payment_method` (`validate:"required,oneof=cash transfer card qris kredit"`) yang **tidak pernah dikirim FE**.
  - BE tidak punya field `payment_date` sama sekali.
  - Bahkan jika field itu diperbaiki, `CreatePayment` di repo meng-hardcode `time.Now()`, mengabaikan tanggal yang dikirim.
- **Dampak**: Setiap submit pembayaran piutang **gagal validasi BE**. Ini alur uang utama fitur piutang dan saat ini tidak berfungsi sama sekali.
- **Kategori**: correctness — **perbaiki paling prioritas**

### 2. [Medium] Tombol "Bayar" di Receivables tidak ada RoleGuard
- **File**: `FE/src/features/customers/receivables/components/ReceivableTableColumns.tsx` (baris 93-103)
- **Pembanding**: `CustomerTableColumns.tsx` (baris 71, 89) membungkus edit/delete dengan `<RoleGuard menuKey="pelanggan.pelanggan" .../>`
- **Dampak**: User tanpa izin `pelanggan.piutang` `can_edit` tetap melihat tombol Bayar yang aktif, baru gagal saat submit (403 dari server) — tidak konsisten dengan pola "sembunyikan di client" pada modul Customers.

### 3. [Medium] Generate kode customer rawan duplikat (race condition)
- **File**: `BE/domain/customer/model/customer.go` (baris 95-102), dipakai di `customer_service.go` (baris 78-82)
- **Masalah**: Kode `CUS-xxx` dihasilkan dari `COUNT(*)` tanpa locking/sequence.
- **Dampak**: Duplicate code berpotensi muncul setelah hard-delete atau pembuatan customer secara bersamaan (TOCTOU).

### 4. [Low] `UpdateCustomerPayload` di-duplikasi manual
- **File**: `FE/src/features/customers/customers/customers.types.ts` (baris 37-43)
- **Dampak**: Sama seperti temuan pola berulang lain — field baru di Create tidak otomatis ikut ke Update.

### 5. [Low] Arsitektur state tabel Customers vs Receivables berbeda
- **File**: `CustomerTable.tsx` (baris 21-141, smart component, punya state sendiri) vs `ReceivableTable.tsx` (baris 14-27, dumb component, state di `ReceivablesPage.tsx`)
- **Dampak**: Dua pola berbeda untuk domain yang sama menyulitkan maintenance jangka panjang; juga menyebabkan reset filter ditangani di layer berbeda (`ReceivableFilterBar.tsx` baris 28-42 vs `CustomerFilterBar.tsx` baris 17-62).

## Yang Sudah Baik
Handler layer BE konsisten (bind → validate → service → WrapResponse). Transisi status pembayaran (`receivable_service.go` baris 40-77) sudah benar menjaga terhadap double-payment dan overpayment.
