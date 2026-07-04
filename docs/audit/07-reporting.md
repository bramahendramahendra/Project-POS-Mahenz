# Audit: Modul Reporting (cashier-performance / profit-loss / sales / stock)

Scope: `FE/src/features/reporting/**`, `BE/domain/report/**`, `BE/routes/segment/report_routes.go`

## Temuan

### 1. [Tinggi] Endpoint laporan sales & stock tanpa permission gate
- **File**: `BE/routes/segment/report_routes.go` (baris 34-38, 44-46)
- **Masalah**: `GET /reports/sales`, `/reports/sales/chart`, `POST /reports/sales/list`, `/reports/sales/summary`, `GET /reports/stock`, `POST /reports/stock/list`, `/reports/stock/summary` tidak ada middleware permission sama sekali.
- **Pembanding**: Profit-loss men-gate semua route dengan `permPL("can_view")` (baris 40-42); cashier-performance dengan `permCashier("can_view")` (baris 49-51); bahkan endpoint export sales/stock sendiri sudah di-gate (`permSales`/`permStock`, baris 36/47).
- **Dampak**: Data penjualan (nominal transaksi, nama customer) dan data stok (harga modal, nilai stok) bisa dibaca siapa pun yang login, tanpa cek izin — hanya endpoint utama sales/stock yang bolong, sementara sibling-nya sudah aman.

### 2. [Tinggi] Revenue profit-loss tidak konsisten dengan sales report
- **File**: `BE/domain/report/repo/report_repo.go` — profit-loss query (baris 56-66) menjumlahkan `ti.subtotal` (level item, sebelum diskon); sales query (baris 21-49) menjumlahkan `t.total_amount` (level transaksi, sudah termasuk diskon).
- **Dampak**: Untuk periode dengan transaksi berdiskon, `total_revenue` profit-loss tidak akan sama dengan sales report, dan gross/net profit di P&L **overstated** sebesar nilai diskon. Perlu konfirmasi ke bisnis: profit seharusnya net dari diskon atau tidak.

### 3. [Medium] Sales & stock report tidak punya kolom sortable/sort wiring
- **File**: `FE/src/features/reporting/sales/components/SalesReportTableColumns.tsx`, `StockReportTableColumns.tsx` — tidak ada `sortable`; `SalesReportTab.tsx` (baris 55-62), `StockReportTab.tsx` (baris 41-48) tidak passing `currentSort`/`onSort`.
- **Pembanding**: `CashierPerformanceTableColumns.tsx` + `CashierPerformanceTab.tsx` (baris 39-43) sudah wiring sort ke BE.
- **Catatan**: DTO `SalesListRequest`/`StockListRequest` juga tidak punya `SortBy`/`SortOrder`, jadi FE memang tidak bisa wiring meski mau.

### 4. [Medium] Export tidak konsisten antar 4 laporan
- **File**: BE expose `GET /reports/{sales,profit-loss,stock,cashier}/export` secara seragam (baris 36, 42, 47, 51), tapi FE hanya `SalesReportFilterBar.tsx` (baris 76-85) punya tombol export — dan itu pun **export CSV client-side dari halaman saat ini**, bukan memanggil endpoint BE xlsx yang sudah ada.
- **Dampak**: Cashier-performance, profit-loss, stock tidak punya tombol export sama sekali walau endpoint BE-nya sudah siap & di-gate; sales report sendiri exportnya juga tidak memakai endpoint yang benar.

### 5. [Medium] Halaman profit-loss tidak menangani error state secara eksplisit
- **File**: `FE/src/features/reporting/profit-loss/components/ProfitLossTab.tsx` (baris 80-86) — hanya handle `isLoading` dan `!report` (empty state), tidak ada `isError`.
- **Dampak**: Request gagal akan tampil sama seperti "Belum ada data", menyesatkan user.

### 6. [Info] Tidak ada validasi date_from <= date_to di modul manapun (FE/BE)
- **Dampak**: Range tanggal terbalik diam-diam menghasilkan hasil kosong tanpa error yang jelas. Konsisten di semua 4 laporan (bukan inkonsistensi sibling, tapi tetap worth-fix untuk UX).

## Yang Sudah Baik
- Stock report membaca `p.stock` langsung dari tabel live (real-time, bukan snapshot basi).
- Void transaksi sudah benar dikeluarkan dari perhitungan sales figures, dan dihitung terpisah di `void_count` cashier report — tidak ada double counting.
- Username enumeration pada login (dicek ulang di modul auth) tidak relevan di sini, disebutkan hanya untuk lengkap.
