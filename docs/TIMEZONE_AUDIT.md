# Audit Timezone — POS Mahenz (BE + FE)

> **STATUS: ✅ SELESAI DITANGANI (9 Agustus 2026)** — Seluruh temuan di audit ini sudah diperbaiki lewat 9 fase migrasi di `TIMEZONE_MIGRATION_PLAN.md` (Fase 0–8), diverifikasi build/lint/test + browser testing. Dokumen ini disimpan sebagai referensi historis kondisi sebelum perbaikan.

> Dibuat: 9 Agustus 2026
> Tujuan: Menyelaraskan seluruh sumber tanggal/jam di aplikasi (backend & frontend) supaya konsisten penuh menggunakan **Asia/Jakarta (WIB)**.
> Dokumen ini murni **audit/inventaris fakta** — rencana eksekusi ada di `TIMEZONE_MIGRATION_PLAN.md`.

---

## 1. Ringkasan Masalah

Aplikasi ini punya **4 sumber "jam sekarang" yang berbeda**, dipakai bercampur tanpa aturan baku:

| # | Sumber | Lokasi | Ikut WIB? |
|---|---|---|---|
| 1 | Jam MySQL (`NOW()`, `CURDATE()`, `CURRENT_TIMESTAMP`) | Backend — banyak query & default kolom | Kebetulan (ikut setting OS server DB), tidak dijamin |
| 2 | Jam Go mentah (`time.Now()` tanpa konversi) | Backend — beberapa service/handler | Tidak dijamin (ikut OS server aplikasi) |
| 3 | Jam WIB eksplisit (`config.Location` / `GetTimeNow()`) | Backend — helper `helper/time/time.go` | Ya, pasti |
| 4 | Jam device/browser user (`new Date()`, `Date.now()`) | Frontend — seluruh kode | Tidak — ikut jam laptop/HP user apa adanya |

**Bukti nyata bug** (direproduksi via API & browser, lihat §5): kas yang `status='open'` tapi `open_time` bukan tanggal hari ini (mis. kebawa nginap dari kemarin) membuat endpoint `/cash-drawer/current` (Dashboard) dan `/cash-drawer/my-cash` (Kas Saya) memberi jawaban BERBEDA untuk kondisi yang seharusnya identik — satu bilang "Kas Buka", satu bilang "Tutup".

---

## 2. Inventaris Backend (Go)

### 2.1 File yang pakai `time.Now()` mentah (17 file)

Tidak ada konversi eksplisit ke `config.Location` — ikut jam OS server aplikasi Go.

| File | Dipakai untuk |
|---|---|
| `domain/backup/service/backup_service.go` | Nama file backup (timestamp) |
| `domain/supplier_purchase/repo/purchase_repo.go` | Kode PO, tanggal pembayaran default |
| `domain/supplier_return/service/supplier_return_service.go` | Validasi tanggal retur (`returnDate.After(time.Now()...)`) |
| `domain/supplier_return/repo/supplier_return_repo.go` | Kode retur |
| `domain/dashboard/service/dashboard_service.go` | "Hari ini" untuk ringkasan dashboard |
| `domain/business_summary/handler/business_summary_handler.go` | Batas periode ringkasan bisnis |
| `domain/business_summary/service/business_summary_service.go` | Batas periode ringkasan bisnis |
| `domain/transaction/repo/transaction_repo.go` | Suffix tanggal kode transaksi, `transaction_date` saat insert |
| `domain/sync/service/sync_service.go` | Timestamp job sinkronisasi |
| `domain/report/handler/sales_report_handler.go` | Default `date_from`/`date_to` laporan penjualan |
| `helper/time/time.go` | (basis `GetTimeNow()` — lihat §2.3) |
| `middleware/login_rate_limit_middleware.go` | Window rate-limit login |
| `scheduler/cash_drawer_scheduler.go` | Pencatatan durasi & log eksekusi (bukan jadwal pemicunya — itu sudah WIB) |
| `repository/log_scheduler_repository.go` | `ExecutedAt` log scheduler |
| `helper/global.go` | Generator ID/token/tahun |
| `middleware/log_request_middleware.go` | Durasi & `createdAt` request log |
| `pkg/janitor/janitor.go` | TTL cache (relatif, risiko rendah) |

### 2.2 File dengan SQL yang pakai `NOW()`/`CURDATE()`/`CURRENT_TIMESTAMP` (24 file)

| File | Query/kolom |
|---|---|
| `domain/cash_drawer/repo/cash_drawer_repo.go` | `open_time`, `close_time`, `updated_at` (`NOW()`); **`getOpenYesterdayQuery`** (`DATE(open_time) < CURDATE()`); **`getMyCashQuery`** (`DATE(open_time) = CURDATE()` ← **sumber bug utama**) |
| `domain/transaction/repo/transaction_repo.go` | `generateTransactionCodeQuery` (`DATE(transaction_date) = CURDATE()`), `updated_at = NOW()` |
| `domain/business_summary/repo/business_summary_repo.go` | `MONTH(NOW())`/`YEAR(NOW())` untuk rekap bulan berjalan |
| `domain/expiry_batch/repo/expiry_batch_repo.go` | `CURDATE()` untuk cek produk akan expired |
| `domain/supplier_purchase/repo/purchase_repo.go` | `updated_at`/`payment_date` |
| `domain/supplier_return/repo/supplier_return_repo.go` | `updated_at` |
| `domain/supplier/repo/supplier_repo.go` | `updated_at` |
| `domain/product_unit/repo/unit_repo.go` | `updated_at` |
| `domain/product/repo/product_repo.go` | `updated_at` |
| `domain/product/repo/product_package_repo.go` | `updated_at` |
| `domain/sync/repo/sync_repo.go` | `desktop_time`, `online_time`, `updated_at` |
| `domain/receivable/repo/receivable_repo.go` | `updated_at` |
| `domain/menu/repo/menu_repo.go` | `updated_at` |
| `domain/user/repo/user_repo.go` | `updated_at` |
| `domain/shift/repo/shift_repo.go` | `updated_at` |
| `domain/customer/repo/customer_repo.go` | `updated_at` |
| `domain/expense/repo/expense_repo.go` | `updated_at` |
| `domain/role/repo/role_repo.go` | `updated_at` |
| `domain/product_category/repo/product_category_repo.go` | `updated_at` |
| `domain/access/repo/access_repo.go` | `updated_at` |
| `domain/setting/repo/setting_repo.go` | `updated_at` |
| `domain/auth/repo/auth_repo.go` | `created_at = NOW()` pada upsert sesi (⚠️ di baris yang sama, `expires_at` dihitung dari Go pakai WIB — **2 sumber jam beda dalam 1 row**) |
| `domain/_sample_reference/repo/userIntegration_data.go` | contoh/referensi, dampak rendah |
| `database/migrate.go` | `executed_at DATETIME DEFAULT CURRENT_TIMESTAMP` (tabel migrasi) |

### 2.3 File yang SUDAH benar (WIB eksplisit via `GetTimeNow()` / `config.Location`) — 13 file

| File | Dipakai untuk |
|---|---|
| `helper/time/time.go` | Helper pusat: `GetTimeNow()`, `GetTimeWithFormat()`, `GetEndTime()`, `NormalizeDateRange()`, dll. |
| `config/config.go` | Load `Cfg.Timezone` → `config.Location` |
| `domain/auth/service/auth_service.go` | `ExpiresAt` token/sesi |
| `domain/expiry_batch/service/expiry_batch_service.go` | Cek "hari ini" untuk expiry |
| `domain/report/repo/report_repo.go` | `NormalizeDateRange` default laporan |
| `domain/report/service/report_service.go` | idem |
| `scheduler/cash_drawer_scheduler.go` | Jadwal pemicu "tengah malam WIB" (baris 23-24) |
| `helper/log/log.go` | Timestamp log |
| `helper/error/error.go` | Timestamp log error |
| `pkg/logger/zap.go` | Penamaan file log harian |
| `pkg/jwt/jwt.go` | Klaim `exp` JWT |
| `middleware/log_request_middleware.go` | `GetTimeWithFormat()` untuk start request (⚠️ file ini JUGA pakai `time.Now()` mentah untuk durasi — campur) |

### 2.4 Skema database — default kolom (60 occurrence, 1 file)

`database/migrations/001_init_schema.sql` — **30+ tabel** (users, products, transactions, purchases, expenses, receivables, cash_drawer, shifts, sync, log_scheduler, dst.) memakai:
```sql
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```
→ Untuk INSERT/UPDATE manapun yang **tidak** mengirim value eksplisit dari Go, MySQL yang menentukan waktunya — independen dari `config.Location`.

---

## 3. Inventaris Frontend (React + TypeScript)

### 3.1 Library tanggal

**Tidak ada.** `package.json` tidak memuat `dayjs`, `date-fns`, `moment`, atau `luxon`. Seluruh logika tanggal ditulis manual dengan `Date` bawaan JavaScript. **Tidak ada satupun kode yang timezone-aware.**

### 3.2 Helper pusat

`FE/src/shared/utils/date.ts` — dipakai di 45+ file lewat `shared/utils/index.ts`:
- `formatDate`, `formatDateShort`, `formatDateTime`, `formatRelative` — format tampilan
- `toISODate`, `todayStr`, `monthStart`, `weekStart` — hitung "hari ini" / rentang tanggal

Semua fungsi ini berbasis `new Date()` browser — **tidak ada normalisasi ke WIB sama sekali**.

### 3.3 File yang panggil `new Date()` / `Date.now()` langsung (9 file)

| File | Dipakai untuk |
|---|---|
| `shared/utils/date.ts` | Basis semua helper tanggal (lihat 3.2) |
| `features/dashboard/components/GreetingHeader.tsx` | Header "Sabtu, 8 Agustus 2026" — murni jam device |
| `features/dashboard/dashboard.utils.ts` | `getGreeting()` — Pagi/Siang/Malam berdasar jam device |
| `features/reporting/business-summary/business-summary.api.ts` | `periodToDateRange()` — pakai `toISOString()` (⚠️ konversi ke **UTC** dulu, risiko lompat tanggal ganda) |
| `features/finance/cash-drawer/components/OpenCashDrawerModal.tsx` | `detectShiftId()` — auto-pilih shift aktif berdasar jam device saat buka kas |
| `features/customers/receivables/components/ReceivableTableColumns.tsx` | `isOverdue()` — badge jatuh tempo piutang |
| `features/procurement/returns/returns.schema.ts` | Cap validasi `return_date <= today` (snapshot UTC saat modul dimuat) |
| `features/settings/printer/components/PrinterSettingsTab.tsx` | Preview struk (kosmetik, tidak dikirim ke server) |
| `features/products/products/products.utils.ts` | Nama file export (tidak berdampak bisnis) |

### 3.4 File yang memakai `todayStr()` / `monthStart()` / `weekStart()` (19 file — turunan dari 3.3)

Semua filter tanggal default & preset "Hari ini/Minggu ini/Bulan ini" yang **dikirim ke backend sebagai parameter query**, browser-local, tanpa penyesuaian WIB:

| File | Fitur |
|---|---|
| `features/finance/cash-drawer/components/CashDrawerTable.tsx` | Kas Harian — rentang default |
| `features/finance/cash-drawer/components/CashDrawerFilterBar.tsx` | Kas Harian — preset "Bulan ini" |
| `features/finance/overview/components/FinanceTable.tsx` | Ringkasan keuangan — rentang default |
| `features/finance/overview/components/FinanceFilterBar.tsx` | Ringkasan keuangan — preset Hari/Minggu/Bulan ini |
| `features/finance/expenses/components/ExpenseTable.tsx` | Pengeluaran — rentang default |
| `features/finance/expenses/components/ExpenseFormModal.tsx` | Default `expense_date` |
| `features/reporting/profit-loss/components/ProfitLossTab.tsx` | Laba rugi — rentang default |
| `features/reporting/profit-loss/components/ProfitLossFilterBar.tsx` | Laba rugi — preset tanggal |
| `features/reporting/sales/components/SalesReportTab.tsx` | Laporan penjualan — rentang default |
| `features/reporting/cashier-performance/components/CashierPerformanceTab.tsx` | Performa kasir — rentang default |
| `features/reporting/cashier-performance/components/CashierPerformanceFilterBar.tsx` | Performa kasir — preset tanggal |
| `features/procurement/purchases/components/PurchaseTable.tsx` | Pembelian — rentang default |
| `features/procurement/purchases/components/PurchaseFormModal.tsx` | Default `purchase_date` |
| `features/procurement/purchases/purchases.schema.ts` | Cap validasi `purchase_date <= today` |
| `features/procurement/purchases/components/PaymentModal.tsx` | Default `payment_date` |
| `features/procurement/returns/components/ReturnTable.tsx` | Retur — rentang default |
| `features/procurement/returns/components/ReturnFormModal.tsx` | Default `return_date` + cap max |
| `features/customers/receivables/components/PaymentRecordModal.tsx` | Default `payment_date` |
| `shared/utils/date.ts` | (sumber) |

### 3.5 Date-picker native (`<input type="date">`) — 14 file

Semua terisi default dari fungsi di §3.4, submit string `YYYY-MM-DD` browser-local apa adanya ke API:
`PurchaseFormModal.tsx`, `ReturnFormModal.tsx`, `ProfitLossFilterBar.tsx`, `CashierPerformanceFilterBar.tsx`, `FinanceFilterBar.tsx`, `ExpenseFormModal.tsx`, `PaymentRecordModal.tsx`, `PaymentModal.tsx`, `SalesReportFilterBar.tsx`, `CashDrawerFilterBar.tsx`, `ReturnFilterBar.tsx`, `PurchaseFilterBar.tsx`, `TransactionFilterBar.tsx`, `ExpenseFilterBar.tsx`.

### 3.6 Duplikasi helper (tidak reuse `date.ts`)

- `features/customers/receivables/components/ReceivableTableColumns.tsx` — implementasi `formatDate()` sendiri
- `features/products/products/components/ExpiryWarningModal.tsx` — implementasi `formatDate()` sendiri

---

## 4. Peta Risiko (diurutkan dari paling kritikal)

| Prioritas | Domain | Kenapa kritikal |
|---|---|---|
| 🔴 P0 | **Kas** (cash_drawer) | Bug sudah **terbukti terjadi di produksi** (kasus Wahyuni) — status "buka/tutup" beda antar halaman, transaksi hilang dari tampilan |
| 🔴 P0 | **Transaksi** — generate kode | `generateTransactionCodeQuery` pakai `CURDATE()` MySQL, tapi `transaction_date` di-insert pakai `time.Now()` Go mentah — risiko kode transaksi duplikat/lompat di sekitar tengah malam |
| 🟠 P1 | **Dashboard** | "Hari ini" dihitung 2x beda sumber (Go raw di `dashboard_service.go`, browser device di FE `GreetingHeader`) |
| 🟠 P1 | **Business Summary** | 3 sumber jam berbeda untuk 1 fitur (SQL `NOW()`, Go raw, dan FE `toISOString()` yang double-convert ke UTC) |
| 🟠 P1 | **Laporan** (Sales, Profit-Loss, Cashier Performance) | Handler pakai `time.Now()` raw, repo pakai `NormalizeDateRange` WIB — dua default berbeda dalam satu alur request; ditambah FE kirim filter browser-local |
| 🟡 P2 | **Auth/Session** | Sebagian besar sudah WIB (`GetTimeNow()`), tapi `created_at` sesi masih `NOW()` MySQL — bisa beda dari `expires_at` |
| 🟡 P2 | **Pembelian/Retur/Piutang** | Validasi tanggal (`<=today`) pakai jam device FE — risiko reject/accept salah di edge case dekat tengah malam |
| 🟢 P3 | **Semua tabel lain (`updated_at`/`created_at` default)** | Risiko rendah selama server DB & aplikasi tetap di WIB, tapi tidak eksplisit/tidak future-proof jika pindah hosting |

---

## 5. Bukti Reproduksi Bug (Testing)

### 5.1 Via API (curl langsung ke backend, tanpa lewat FE)

**Kondisi awal** — kas baru dibuka (`open_time` = hari ini), kedua endpoint konsisten:
```json
// POST /api/cash-drawer/current
{"data":{"id":2,"status":"open","open_time":"2026-08-09T11:10:51+07:00", ...}}

// POST /api/cash-drawer/my-cash
{"data":{"id":2,"status":"open","open_time":"2026-08-09T11:10:51+07:00", ...}}
```

**Setelah `open_time` dimundurkan 1 hari secara manual** (mensimulasikan kas yang kebawa nginap):
```sql
UPDATE cash_drawer SET open_time = DATE_SUB(open_time, INTERVAL 1 DAY) WHERE id = 2;
-- Hasil: open_time = 2026-08-08 11:10:51, tgl_buka = 2026-08-08, tgl_sekarang (CURDATE()) = 2026-08-09
```

```json
// POST /api/cash-drawer/current  →  MASIH "open" (benar, status kolom memang open)
{"data":{"id":2,"status":"open","open_time":"2026-08-08T11:10:51+07:00", ...}}

// POST /api/cash-drawer/my-cash  →  JADI "closed" + id null (BUG)
{"data":{"id":null,"status":"closed","open_time":null,"transactions":[], ...}}
```

### 5.2 Via Browser (screenshot, kondisi sama seperti di atas)

- **Dashboard** → badge hijau **"Kas Buka"**, shift & jam tampil benar, tombol "Mulai Transaksi"
- **Kas Saya** (`/finance/my-cash`) → badge abu-abu **"Tutup"**, 0 transaksi, tombol kembali jadi "Buka Kas"

→ **Persis mereplikasi laporan bug asli dari user "Wahyuni Putri Romadhona".**

### 5.3 Root cause pasti (bukan dugaan)

`BE/domain/cash_drawer/repo/cash_drawer_repo.go`:
```go
// baris 27 — dipakai /cash-drawer/current (Dashboard)
getCurrentCashDrawerQuery = `... WHERE cd.user_id = ? AND cd.status = 'open' LIMIT 1`

// baris 125-133 — dipakai /cash-drawer/my-cash (Kas Saya)
getMyCashQuery = `... WHERE cd.user_id = ? AND DATE(cd.open_time) = CURDATE() AND cd.status = 'open' LIMIT 1`
```
Syarat tambahan `DATE(cd.open_time) = CURDATE()` di `getMyCashQuery` itulah yang menyebabkan mismatch — bukan soal selisih timezone server (sudah dicek: MySQL `SYSTEM` timezone = `SE Asia Standard Time`, sama dengan `config.Location` app = `Asia/Jakarta`, tidak ada offset). Murni karena `open_time` (ditulis via `NOW()` MySQL) sudah lewat batas tanggal kalender sebelum baris ini dibaca ulang.

---

## 6. Kesimpulan Audit

1. **Config `Timezone: Asia/Jakarta` di `config_dev.json`/`config_prod.json` sudah BENAR dan BERFUNGSI** — tapi hanya dipakai di ~13 file backend (auth, JWT, report, scheduler-trigger, logging). Sisanya (mayoritas kode) melewatinya sama sekali.
2. **Frontend 100% tidak punya kesadaran timezone** — semua "hari ini" dihitung dari jam device user, dikirim mentah-mentah ke backend sebagai parameter.
3. Ini bukan bug tunggal — ini **pola arsitektur yang tidak konsisten**, berpotensi muncul lagi di domain lain (transaksi, laporan, business summary) dengan gejala serupa (data beda antar layar, kode duplikat, validasi salah).
4. Scheduler auto-close kas (`AutoCloseYesterday`) **berfungsi normal** dan bukan bagian dari masalah ini — dia hanya menangani kasus "kas dari hari sebelumnya", bukan mismatch dalam hari yang sama.

Rencana eksekusi penyelarasan lengkap: lihat **`TIMEZONE_MIGRATION_PLAN.md`**.
