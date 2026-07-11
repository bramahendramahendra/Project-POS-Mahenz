# Audit: Modul Sales (cashier / transactions) — Checkout Flow

Scope: `FE/src/features/sales/**`, `BE/domain/transaction/**`, `BE/routes/segment/transaction_routes.go`

Modul paling kritis karena ini alur checkout inti POS.

> **Update 2026-07-11**: Temuan #1, #2, #3, #4, #5, #6 diverifikasi ulang lewat testing keamanan
> manual (curl dengan payload yang di-tamper) dan pembacaan kode terbaru — **semua sudah tidak
> valid**, sudah diperbaiki di kode saat ini. Dalam proses verifikasi ditemukan 2 bug baru
> (gap `payment_amount`/`change_amount`, dan bug double-counting `expected_balance` di
> cash-drawer akibat pola self-reference MySQL `UPDATE ... SET a = a + ?, b = CASE WHEN a ...`)
> — keduanya sudah diperbaiki sekaligus. Detail di catatan masing-masing di bawah.

## Temuan

### 1. [SUDAH DIPERBAIKI] ~~BE tidak pernah menghitung ulang uang — percaya total dari client~~
- **Status**: Tidak valid lagi. `BE/pkg/pricing/pricing.go` (`Recalculate`) dipanggil dari
  `transaction_service.go` (`recalculateTotals`) sebelum insert — subtotal, harga per-item, dan
  total_amount dihitung ulang dari master data produk, mengabaikan nilai dari payload client.
  Diverifikasi: kirim `price:1` untuk produk seharga Rp3.000 → tersimpan `price:3000` (harga
  asli), bukan nilai yang di-tamper.
- **Catatan turunan yang DITEMUKAN & DIPERBAIKI (2026-07-11)**: meski subtotal/total_amount
  sudah benar, `payment_amount`/`change_amount` **tidak** ikut divalidasi/direkalkulasi — bisa
  submit "dibayar Rp1" untuk transaksi Rp3.000 tanpa ditolak, dan `change_amount` bisa
  dipalsukan bebas. Diperbaiki di `transaction_service.go` (`Create`): tolak kalau
  `payment_amount < total_amount` untuk transaksi non-kredit, dan `change_amount` selalu
  dihitung ulang server (`payment_amount - total_amount`, floor di 0).

### 2. [SUDAH DIPERBAIKI] ~~Endpoint create transaksi tanpa permission middleware~~
- **Status**: Tidak valid lagi. `BE/routes/segment/transaction_routes.go:31` — `POST /create`
  sudah di-gate `perm("can_create")` pada key `penjualan.transaksi`, pola sama dengan
  `supplier_purchase_routes.go`. Diverifikasi empiris: user tanpa izin ini mendapat 403.

### 3. [SUDAH DIPERBAIKI] ~~`shift_id` tidak divalidasi terhadap user yang login~~
- **Status**: Tidak valid lagi. `transaction_service.go` (`Create`, baris ~34-43) mencari kas
  terbuka milik `userID` dari token (bukan dari payload), dan menolak kalau `shift_id` di
  payload tidak cocok dengan shift dari kas milik user tersebut. Diverifikasi: dua user dengan
  kas terbuka pada `shift_id` yang sama — transaksi user A hanya mengkredit kas milik A, kas
  user B tidak berubah sama sekali.

### 4. [SUDAH DIPERBAIKI] ~~Void transaksi tidak membalikkan total penjualan cash-drawer~~
- **Status**: Tidak valid lagi. `transaction_service.go` (`Void`, baris ~132-142) memanggil
  `cashDrawerRepo.UpdateSales(drawer.ID, -t.TotalAmount, -t.TotalAmount)` untuk transaksi cash,
  membalikkan `total_sales`/`total_cash_sales` kas milik pemilik transaksi.
- **Bug turunan yang DITEMUKAN & DIPERBAIKI (2026-07-11)**: `updateSalesQuery` di
  `cash_drawer_repo.go` punya bug self-reference MySQL — kolom `expected_balance` dihitung dari
  `total_cash_sales` yang **sudah** ter-update oleh assignment sebelumnya di statement yang
  sama, lalu menambahkan parameter yang sama sekali lagi, sehingga `expected_balance` selalu
  lebih tinggi dari seharusnya sebesar nilai transaksi (create) atau lebih rendah (void, bug
  yang sama juga ada di `updateExpensesQuery`). Diverifikasi dengan data nyata: sebelum fix,
  `expected_balance` kas kasir1 tersimpan Rp225.500 padahal seharusnya Rp222.500 (selisih
  persis sebesar transaksi terakhir). Fix: hapus referensi ganda ke kolom yang sudah
  ter-update, biarkan MySQL membaca nilai barunya secara alami (satu kali saja).

### 5. [SUDAH DIPERBAIKI] ~~Update cash-drawer tidak atomik dengan create/void~~
- **Status**: Tidak valid lagi. `transaction_service.go` (`Create` & `Void`) memanggil
  `cashDrawerRepo.UpdateSales` **di dalam** `s.repo.GetDB().Transaction(func(tx *gorm.DB) ...)`
  yang sama dengan insert/void transaksi (pakai `cashDrawerRepo.WithTx(tx)`) — satu DB
  transaction, rollback otomatis kalau salah satu langkah gagal.

### 6. [SUDAH DIPERBAIKI] ~~Field diskon tidak ada batas atas di DTO~~
- **Status**: Tidak valid lagi. `dto_transaction.go` — `Discount` dan `DiscountItem` punya
  validasi `ltefield=Subtotal`. Diverifikasi 2 lapis: (1) diskon > subtotal yang dikirim →
  ditolak validator (`ltefield`); (2) subtotal dipalsukan tinggi supaya diskon "lolos" relatif
  → tetap ditolak di `pricing.Recalculate` karena dibandingkan terhadap subtotal hasil hitung
  ulang server, bukan subtotal dari payload.

### 7. [Medium] Model rounding tidak konsisten antara diskon cart-level dan item-level
- **File**: `FE/src/features/sales/cashier/cashier.utils.ts` — diskon cart-level (baris 19, 31) membulatkan nilai akhir; diskon persen item-level (baris 79-84) membulatkan harga satuan dulu baru dikalikan quantity.
- **Dampak**: Kesalahan pembulatan yang berbeda antara dua jalur, tanpa rekonsiliasi.

### 8. [Low] Cart di localStorage tidak di-refresh sebelum checkout
- **File**: `FE/src/features/sales/cashier/cashier.store.ts` (baris 218-227) — cart/diskon/pajak dipersist ke localStorage dengan harga yang diambil saat add-to-cart (`ProductSearch.tsx` baris 95-104), tidak pernah di-refresh.
- **Dampak (direvisi)**: Sejak temuan #1 tidak valid lagi (BE selalu merekalkulasi harga dari master data), ini murni isu UX (kasir bisa melihat harga lama sesaat di keranjang sebelum ditolak/dikoreksi saat submit kalau harga produk berubah di tengah shift), bukan lagi celah keamanan.

### 9. [Low] Perlu verifikasi lanjutan: void tidak cek pembayaran parsial pada receivable
- **File**: `BE/domain/transaction/repo/transaction_repo.go` (baris 25, 350) — void set receivable ke `'void'` tanpa cek apakah sudah ada pembayaran parsial.

### 10. [Low] Tabel transaksi tanpa kolom sortable
- **File**: `FE/src/features/sales/transactions/components/TransactionTableColumns.tsx`
- **Pembanding**: `PurchaseTableColumns.tsx` (baris 32, 45, 54, 63) punya sortable.

## Prioritas Perbaikan
~~#1 dan #2/#3 paling parah...#4/#5 berikutnya...~~ — **semua sudah diperbaiki per 2026-07-11**,
lihat catatan status di masing-masing temuan di atas. Sisa yang masih relevan: #7 (rounding
tidak konsisten), #8 (cart stale, kini murni UX), #9 (void vs receivable parsial, perlu
verifikasi terpisah), #10 (kolom sortable).
