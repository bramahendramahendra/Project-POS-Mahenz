# Rencana Migrasi Timezone → Asia/Jakarta (WIB) Penuh

> **STATUS: ✅ SELESAI (9 Agustus 2026)** — Fase 0–8 seluruhnya tuntas dan diverifikasi (build, lint, test, browser). Lihat ringkasan penutup di paling bawah dokumen ini.

> Dasar: `TIMEZONE_AUDIT.md`
> Tujuan akhir: **satu aturan tunggal** — semua tanggal/jam di BE dan FE dihitung/ditampilkan dalam Asia/Jakarta (WIB), tidak ada lagi campuran jam MySQL / jam OS server Go / jam device browser.
> Sifat pekerjaan: **bukan bug fix darurat** — ini pekerjaan konsistensi menyeluruh, dikerjakan bertahap per fase supaya bisa diverifikasi satu-satu.

---

## Aturan Baku yang Akan Diberlakukan

**Backend (Go):**
- Semua kebutuhan "jam sekarang" WAJIB lewat `time_helper.GetTimeNow()` (sudah ada di `helper/time/time.go`) — tidak boleh ada `time.Now()` mentah lagi.
- Semua SQL yang menulis timestamp WAJIB terima value dari Go (parameter `?`), bukan `NOW()`/`CURDATE()`/`CURRENT_TIMESTAMP` bawaan MySQL.
- Semua SQL yang membandingkan tanggal ("hari ini", "kemarin", dst.) WAJIB terima batas tanggal dari Go (WIB), bukan `CURDATE()` MySQL.
- Kolom `DEFAULT CURRENT_TIMESTAMP` di skema tabel — dipertahankan sebagai *fallback* saja (bukan sumber utama), karena mengubah default kolom butuh migration terpisah dan berisiko ke data lama.

**Frontend (React):**
- Install `dayjs` + plugin `utc` & `timezone`.
- Buat instance ter-konfigurasi WIB di satu tempat (`shared/utils/date.ts`), semua fungsi tanggal existing (`todayStr`, `monthStart`, `weekStart`, `formatDate`, dst.) direfaktor untuk pakai instance ini.
- Tidak ada lagi `new Date()` / `Date.now()` langsung di kode fitur — semua wajib lewat helper pusat.

---

## Fase Pengerjaan

Urutan berdasarkan prioritas risiko di `TIMEZONE_AUDIT.md` §4. Setiap fase = satu unit kerja yang bisa dites & di-commit terpisah.

| Fase | Domain | Scope |
|---|---|---|
| 0 | Persiapan | FE: install dayjs, refactor `date.ts`. BE: pastikan `GetTimeNow()` helper lengkap (tambah varian jika perlu) |
| 1 | **Kas (cash_drawer)** | Perbaikan bug asli + selaraskan semua query kas ke WIB |
| 2 | **Transaksi** | Kode transaksi + `transaction_date` |
| 3 | **Dashboard** | `dashboard_service.go` + `GreetingHeader.tsx`/`dashboard.utils.ts` |
| 4 | **Business Summary** | Satukan 3 sumber jam jadi 1 |
| 5 | **Laporan** (Sales, Profit-Loss, Cashier Performance) | Handler + repo + FE filter bar |
| 6 | **Auth/Session** | Selaraskan `created_at` sesi ke WIB |
| 7 | **Pembelian/Retur/Piutang** | Validasi tanggal FE + kode PO/retur BE |
| 8 | **Sisanya** (`updated_at` generic di 20+ repo) | Sapu bersih file yang tersisa dari audit §2.1/2.2 |

Setiap fase FE dan BE untuk domain yang sama dikerjakan **berpasangan** (tidak ada gunanya benerin BE doang kalau FE masih kirim tanggal browser-local, begitu juga sebaliknya).

---

## Cara Pakai Dokumen Ini

Setiap fase di bawah punya blok **"PROMPT EKSEKUSI"** — itu teks siap-tempel untuk memulai sesi baru mengerjakan fase tersebut. Jalankan satu fase, review, test, commit, baru lanjut fase berikutnya. Jangan gabung banyak fase dalam satu sesi supaya mudah di-review dan di-rollback kalau ada masalah.

---

## FASE 0 — Persiapan

### Scope
- FE: `npm install dayjs`, buat konfigurasi timezone di `shared/utils/date.ts`
- BE: cek kelengkapan `helper/time/time.go` — pastikan ada fungsi untuk: jam sekarang (`GetTimeNow`), awal hari (`StartOfDay`), akhir hari (`GetEndTime` — sudah ada), format tanggal SQL-ready (`YYYY-MM-DD`)

### Testing
- BE: unit test helper waktu memastikan konsisten menghasilkan WIB
- FE: cek `dayjs().tz('Asia/Jakarta').format()` menghasilkan waktu yang benar di browser dengan timezone device berbeda-beda (bisa disimulasikan lewat DevTools Sensors)

### PROMPT EKSEKUSI — Fase 0
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz.

Kerjakan FASE 0 (Persiapan) dari rencana migrasi timezone:

BACKEND:
1. Buka BE/helper/time/time.go, pastikan ada fungsi-fungsi berikut (tambahkan jika belum ada):
   - GetTimeNow() time.Time — sudah ada, jangan diubah
   - StartOfDay(t time.Time) time.Time — jam 00:00:00 WIB dari tanggal t
   - EndOfDay(t time.Time) time.Time — jam 23:59:59 WIB dari tanggal t (cek dulu apakah GetEndTime() sudah memenuhi ini)
   - ToSQLDate(t time.Time) string — format "2006-01-02" untuk dipakai sebagai parameter query MySQL
2. Jangan ubah behavior GetTimeNow() yang sudah ada, cuma tambahkan fungsi baru kalau memang belum ada.

FRONTEND:
1. Install dayjs: npm install dayjs (di folder FE)
2. Di FE/src/shared/utils/date.ts, tambahkan setup dayjs dengan plugin utc dan timezone, set default timezone ke 'Asia/Jakarta'. Buat helper baru getWIBNow() yang mengembalikan dayjs object WIB.
3. JANGAN ubah dulu implementasi todayStr/monthStart/weekStart/formatDate yang sudah ada — itu tugas fase berikutnya. Fase ini cuma nyiapkan fondasinya.

Setelah selesai, jalankan:
- BE: go build ./... untuk pastikan compile sukses
- FE: npm run type-check && npm run lint

Laporkan hasil, jangan lanjut ke fase lain.
```

---

## FASE 1 — Kas (cash_drawer) 🔴 P0

### Scope Backend
File: `BE/domain/cash_drawer/repo/cash_drawer_repo.go`
1. **Fix bug utama**: `getMyCashQuery` (baris ~132) — hapus `AND DATE(cd.open_time) = CURDATE()`, samakan logic dengan `getCurrentCashDrawerQuery` (cukup `status = 'open'`)
2. `openCashDrawerQuery` (baris 30) — ganti `NOW()` jadi parameter `?`, kirim `time_helper.GetTimeNow()` dari service/handler
3. `closeCashDrawerQuery` (baris 32), `updateSalesQuery`/`updateExpensesQuery` (baris 33-34) — sama, ganti `NOW()` jadi parameter
4. `getOpenYesterdayQuery` (baris 81) — ganti `CURDATE()` jadi parameter tanggal WIB dari Go
5. Baris 61, 70, 76, 106, 117 — audit ulang satu-satu, ganti `NOW()` fallback jadi parameter WIB
6. Sesuaikan service layer (`cash_drawer_service.go`) dan repo interface untuk terima parameter waktu tambahan ini

### Scope Frontend
1. `FE/src/features/finance/cash-drawer/components/OpenCashDrawerModal.tsx` — `detectShiftId()` pakai `getWIBNow()` bukan `new Date()`
2. `FE/src/features/finance/cash-drawer/cash-drawer.api.ts`, `CashDrawerTable.tsx`, `CashDrawerFilterBar.tsx` — default rentang tanggal pakai helper WIB baru
3. `FE/src/features/finance/my-cash/` — cek ulang tidak ada pemakaian `new Date()` tersembunyi

### Testing (WAJIB, browser + API)
1. **Regresi dasar**: buka kas → cek Dashboard, Kas Saya, Kas Harian tampil konsisten "Buka" — via browser
2. **Reproduksi kasus bug**: buka kas, mundurkan `open_time` manual 1 hari via SQL (`UPDATE cash_drawer SET open_time = DATE_SUB(open_time, INTERVAL 1 DAY) WHERE id = ?`), reload ketiga halaman → SEKARANG HARUS TETAP KONSISTEN "Buka" di ketiganya (sebelum fix: Kas Saya salah jadi "Tutup")
3. **Via API langsung** (curl): panggil `/cash-drawer/current` dan `/cash-drawer/my-cash` dengan kondisi sama seperti di atas, bandingkan response JSON — harus identik status-nya
4. Tutup kas → cek scheduler (`AutoCloseYesterday`) masih jalan benar untuk kas yang genuinely ketinggalan dari hari sebelumnya
5. Test dekat tengah malam WIB (23:55-00:05) kalau memungkinkan — atau simulasikan dengan mengubah jam sistem test/mock waktu

### PROMPT EKSEKUSI — Fase 1
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0 (persiapan helper waktu WIB di BE dan FE) sudah selesai dikerjakan sebelumnya.

Kerjakan FASE 1 (Kas / cash_drawer) dari rencana migrasi timezone. Ini memperbaiki bug nyata: endpoint /cash-drawer/current (Dashboard) dan /cash-drawer/my-cash (Kas Saya) bisa menampilkan status berbeda untuk kas yang sama, karena getMyCashQuery di BE/domain/cash_drawer/repo/cash_drawer_repo.go punya syarat tambahan "DATE(open_time) = CURDATE()" yang tidak ada di getCurrentCashDrawerQuery.

BACKEND — di BE/domain/cash_drawer/repo/cash_drawer_repo.go, service/cash_drawer_service.go, handler terkait:
1. getMyCashQuery: hapus syarat "AND DATE(cd.open_time) = CURDATE()", samakan dengan getCurrentCashDrawerQuery (cukup status='open')
2. openCashDrawerQuery, closeCashDrawerQuery, updateSalesQuery, updateExpensesQuery: ganti semua NOW() jadi parameter waktu yang dikirim dari Go, pakai time_helper.GetTimeNow() dari service layer
3. getOpenYesterdayQuery: ganti CURDATE() jadi parameter tanggal WIB dikirim dari Go
4. Baris-baris lain di file ini yang masih pakai NOW()/CURDATE() (cek baris 61,70,76,106,117 sebagai referensi awal, tapi baca ulang file karena nomor baris bisa berubah): audit dan ganti jadi parameter WIB
5. Sesuaikan signature function di repo interface (BE/domain/cash_drawer/repo/cash_drawer_repo_interface.go jika ada) dan service layer supaya waktu WIB mengalir dari service ke repo

FRONTEND — di FE/src/features/finance/cash-drawer/ dan FE/src/features/finance/my-cash/:
1. OpenCashDrawerModal.tsx: fungsi detectShiftId() ganti new Date() jadi getWIBNow() dari shared/utils/date.ts (fungsi ini sudah dibuat di Fase 0)
2. Cek semua file di kedua folder itu untuk pemakaian new Date()/Date.now() lain yang belum ketahuan, ganti ke helper WIB

TESTING WAJIB (lakukan semua, jangan lewati):
1. Jalankan BE (go run main.go) dan FE (npm run dev) di project ini
2. Login sebagai admin/admin123 via browser (gunakan playwright-core headless jika tersedia, atau tools browser testing lain yang ada)
3. Buka kas via Dashboard, screenshot Dashboard + Kas Saya + Kas Harian — pastikan ketiganya konsisten "Buka"
4. Mundurkan open_time kas itu 1 hari via query SQL langsung ke MySQL (database pos_retail_db, tabel cash_drawer)
5. Reload ketiga halaman, screenshot lagi — WAJIB ketiganya tetap konsisten "Buka" (ini yang membuktikan fix berhasil, sebelumnya Kas Saya akan salah jadi "Tutup")
6. Test juga langsung via curl ke /api/cash-drawer/current dan /api/cash-drawer/my-cash dengan kondisi yang sama, bandingkan response JSON-nya harus konsisten
7. Tutup kas, verifikasi close berjalan normal
8. go build ./... dan npm run type-check && npm run lint harus lulus

Laporkan hasil lengkap dengan bukti screenshot dan response API. Jangan lanjut ke fase lain.
```

---

## FASE 2 — Transaksi 🔴 P0

### Scope Backend
File: `BE/domain/transaction/repo/transaction_repo.go`
1. `generateTransactionCodeQuery` (baris 21) — `DATE(transaction_date) = CURDATE()` → kirim tanggal WIB sebagai parameter
2. Baris 193, 197, 294, 444, 448 — `time.Now()` mentah (suffix kode, `TransactionDate`) → ganti `GetTimeNow()`
3. Pastikan **satu titik waktu WIB** dipakai konsisten untuk: generate kode transaksi DAN nilai `transaction_date` yang di-insert (saat ini berpotensi 2 pembacaan jam berbeda dalam 1 request)

### Testing
1. Buat transaksi mendekati tengah malam WIB (atau simulasi dengan mundurkan/majukan jam sistem test) — pastikan kode transaksi tidak duplikat/skip
2. Cek beberapa transaksi berurutan dalam 1 hari — nomor urut kode harus tetap sekuensial
3. Via API: `POST` beberapa transaksi berturutan, cek field `transaction_date` di response konsisten dengan kode yang di-generate

### PROMPT EKSEKUSI — Fase 2
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0 dan Fase 1 (kas) sudah selesai.

Kerjakan FASE 2 (Transaksi) dari rencana migrasi timezone.

BACKEND — BE/domain/transaction/repo/transaction_repo.go:
1. generateTransactionCodeQuery: ganti CURDATE() jadi parameter tanggal WIB dikirim dari Go (pakai time_helper.GetTimeNow())
2. Cari semua pemakaian time.Now() mentah di file ini (referensi awal: sekitar baris 193,197,294,444,448, tapi baca ulang karena nomor baris bisa berubah) — ganti ke time_helper.GetTimeNow()
3. PENTING: pastikan generate kode transaksi dan insert transaction_date memakai NILAI WAKTU YANG SAMA (hitung sekali di awal function, jangan panggil GetTimeNow() berkali-kali di function yang sama untuk hal yang seharusnya konsisten)
4. Cek juga service/handler transaction untuk time.Now() lain yang terlewat

TESTING WAJIB:
1. Jalankan BE, login via API dapatkan token
2. Buat beberapa transaksi berturutan via API (POST /api/transactions atau endpoint yang sesuai), cek kode transaksi sekuensial dan tidak duplikat
3. Cek response transaction_date konsisten dengan tanggal WIB saat ini
4. go build ./... harus lulus

Laporkan hasil dengan bukti response API. Jangan lanjut ke fase lain.
```

---

## FASE 3 — Dashboard 🟠 P1

### Scope
- BE: `domain/dashboard/service/dashboard_service.go` — ganti `time.Now()` → `GetTimeNow()`
- FE: `features/dashboard/components/GreetingHeader.tsx`, `features/dashboard/dashboard.utils.ts` (`getGreeting`) — ganti `new Date()` → helper WIB

### Testing
- Browser: cek header dashboard tampil tanggal & sapaan (Pagi/Siang/Malam) sesuai jam WIB, bukan jam device — bisa divalidasi dengan ubah timezone device di DevTools lalu reload, harus tetap tampil WIB
- Cek angka "Transaksi Saya Hari Ini"/"Penjualan Saya Hari Ini" konsisten dengan data aktual (query API `/dashboard/today-summary`)

### PROMPT EKSEKUSI — Fase 3
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-2 sudah selesai.

Kerjakan FASE 3 (Dashboard) dari rencana migrasi timezone.

BACKEND — BE/domain/dashboard/service/dashboard_service.go:
1. Ganti time.Now() jadi time_helper.GetTimeNow() untuk perhitungan "hari ini" di ringkasan dashboard

FRONTEND — FE/src/features/dashboard/:
1. dashboard.utils.ts: getGreeting() ganti new Date() jadi getWIBNow() dari shared/utils/date.ts
2. components/GreetingHeader.tsx: ganti new Date().toLocaleDateString(...) jadi format dari getWIBNow(), pastikan hasil tampilan formatnya tetap sama ("Sabtu, 8 Agustus 2026" style, locale id-ID)

TESTING WAJIB (via browser):
1. Jalankan BE dan FE
2. Login, buka Dashboard, screenshot header (tanggal + sapaan)
3. Ganti timezone browser di Chrome DevTools (Sensors > Location/Timezone) ke timezone lain (misal America/New_York), reload halaman
4. Screenshot lagi — header WAJIB tetap tampil tanggal & sapaan WIB, TIDAK berubah ikut timezone device
5. npm run type-check && npm run lint harus lulus, go build ./... harus lulus

Laporkan hasil dengan bukti screenshot sebelum/sesudah ganti timezone device. Jangan lanjut ke fase lain.
```

---

## FASE 4 — Business Summary 🟠 P1

### Scope
- BE: `domain/business_summary/repo/business_summary_repo.go` (SQL `NOW()`), `handler/business_summary_handler.go` + `service/business_summary_service.go` (Go `time.Now()`) — satukan ke `GetTimeNow()` + parameter
- FE: `features/reporting/business-summary/business-summary.api.ts` — `periodToDateRange()` ganti dari `toISOString()` (UTC) ke helper WIB

### Testing
- Browser: cek halaman business summary top-products, ganti timezone device, pastikan data yang tampil tidak berubah
- API: bandingkan hasil `MONTH(NOW())`/`YEAR(NOW())` lama vs parameter baru — pastikan angka sama persis untuk data yang sama

### PROMPT EKSEKUSI — Fase 4
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-3 sudah selesai.

Kerjakan FASE 4 (Business Summary) dari rencana migrasi timezone. Domain ini sebelumnya punya 3 sumber jam berbeda: SQL NOW() di repo, Go time.Now() di service/handler, dan FE toISOString() yang konversi ke UTC dulu (double bug potensial).

BACKEND:
1. BE/domain/business_summary/repo/business_summary_repo.go: ganti MONTH(NOW())/YEAR(NOW()) dan pattern serupa jadi terima parameter bulan/tahun WIB dari Go
2. BE/domain/business_summary/handler/business_summary_handler.go dan service/business_summary_service.go: ganti semua time.Now() jadi time_helper.GetTimeNow()

FRONTEND — FE/src/features/reporting/business-summary/business-summary.api.ts:
1. periodToDateRange(): ganti implementasi dari new Date()/toISOString() jadi pakai getWIBNow() dari shared/utils/date.ts, JANGAN pakai toISOString() sama sekali (itu convert ke UTC, sumber bug ganda)

TESTING WAJIB (browser + API):
1. Jalankan BE dan FE, login, buka halaman Business Summary / Top Products
2. Screenshot data yang tampil untuk periode "bulan ini"
3. Ganti timezone browser (DevTools Sensors) ke timezone lain, reload, screenshot lagi — data periode WAJIB tidak berubah
4. Via curl, panggil endpoint business-summary langsung, cek date range yang dihasilkan konsisten dengan WIB
5. go build ./... dan npm run type-check && npm run lint harus lulus

Laporkan hasil dengan bukti. Jangan lanjut ke fase lain.
```

---

## FASE 5 — Laporan (Sales, Profit-Loss, Cashier Performance) 🟠 P1

### Scope
- BE: `domain/report/handler/sales_report_handler.go` (default `time.Now()`) — selaraskan dengan `report_repo.go`/`report_service.go` yang sudah pakai `NormalizeDateRange` (WIB)
- FE: `features/reporting/profit-loss/`, `features/reporting/sales/`, `features/reporting/cashier-performance/` — semua filter bar & default range ganti dari `todayStr()`/`monthStart()`/`weekStart()` versi lama ke versi WIB baru (setelah `date.ts` direfaktor di sini)

### Refactor `date.ts` (bagian dari fase ini)
Ini titik pentingnya: refactor `FE/src/shared/utils/date.ts` supaya `todayStr()`, `monthStart()`, `weekStart()`, `formatDate`, `formatDateTime`, `formatRelative` semua pakai `dayjs().tz('Asia/Jakarta')` di baliknya — begitu direfactor di sini, **otomatis berlaku** untuk 19 file yang sudah listing di audit §3.4 (tidak perlu ubah manual satu-satu di file konsumennya, karena mereka cuma import fungsi ini).

**Temuan tambahan (ditemukan saat testing Fase 4, dikonfirmasi masuk scope di sini):** `FE/src/features/reporting/business-summary/components/SalesChart.tsx` punya fungsi lokal `formatDateLabel()` yang mem-parse label tanggal dari API pakai `new Date(label)` miliknya sendiri **sebelum** mengoper hasilnya ke `formatDateShort()`. Karena `date.ts`'s `toDate()` cuma re-wrap value yang sudah berupa `Date` (tidak parse ulang), refactor `date.ts` saja tidak akan otomatis membenarkan file ini — parsing manual di `formatDateLabel()` perlu dihapus, biarkan string mentah dari API (`label`) langsung dioper ke `formatDateShort()` supaya lewat jalur parsing WIB yang sama. Pola bug-nya sama persis dengan `ReceivableTableColumns.tsx`/`ExpiryWarningModal.tsx` di atas (re-implementasi parsing sendiri, bypass helper pusat) — tambahkan file ini ke daftar yang diperbaiki di langkah 4 pada prompt eksekusi.

### Testing
- Browser: buka masing-masing halaman laporan, cek preset "Hari ini/Minggu ini/Bulan ini" menghasilkan rentang tanggal WIB yang benar meski timezone device diganti
- Cek 19 file konsumen `date.ts` lainnya (Kas Harian, Pengeluaran, Pembelian, Retur, Piutang, dll) tidak regresi — spot check beberapa

### PROMPT EKSEKUSI — Fase 5
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-4 sudah selesai.

Kerjakan FASE 5 (Laporan + refactor pusat date.ts) dari rencana migrasi timezone. Ini fase paling luas dampaknya karena me-refactor helper pusat yang dipakai 19+ file.

BACKEND — BE/domain/report/handler/sales_report_handler.go:
1. Default date_from/date_to yang sekarang pakai time.Now() raw: ganti jadi time_helper.GetTimeNow(), selaraskan dengan NormalizeDateRange yang sudah dipakai di report_repo.go/report_service.go supaya tidak ada 2 default berbeda dalam satu alur request

FRONTEND — REFACTOR PUSAT FE/src/shared/utils/date.ts:
1. Refactor todayStr(), monthStart(), weekStart() supaya berbasis dayjs().tz('Asia/Jakarta') (setup dayjs sudah ada dari Fase 0), BUKAN new Date() browser lagi
2. Refactor formatDate, formatDateShort, formatDateTime, formatRelative supaya konsisten pakai instance dayjs WIB yang sama untuk parsing DAN formatting
3. JANGAN ubah signature/return type function-function ini (tetap return string sama seperti sebelumnya) supaya 19 file consumer TIDAK PERLU diubah satu-satu -- mereka otomatis ikut benar begitu date.ts diperbaiki
4. Setelah refactor, baca ulang FE/src/features/customers/receivables/components/ReceivableTableColumns.tsx dan FE/src/features/products/products/components/ExpiryWarningModal.tsx -- kedua file ini punya implementasi formatDate() sendiri (duplikat, tidak reuse date.ts), ganti supaya reuse helper pusat yang sudah diperbaiki
5. Juga perbaiki FE/src/features/reporting/business-summary/components/SalesChart.tsx -- fungsi lokal formatDateLabel() di file ini melakukan new Date(label) sendiri sebelum manggil formatDateShort(), sehingga bypass parsing WIB yang baru diperbaiki di date.ts. Hapus parsing manual itu, biarkan string label mentah dari API langsung diproses formatDateShort(). (Ditemukan saat testing Fase 4 -- gejalanya: label sumbu-X grafik "Grafik Penjualan" berubah tanggal kalau timezone device diganti, padahal data di baliknya sudah benar.)

TESTING WAJIB (browser, cakupan luas karena dampaknya lebar):
1. Jalankan BE dan FE, login
2. Buka SATU PER SATU halaman berikut, screenshot rentang tanggal default yang tampil: Kas Harian, Pengeluaran, Laba Rugi, Laporan Penjualan, Performa Kasir, Pembelian, Retur, Piutang, Ringkasan Keuangan
3. Ganti timezone browser (DevTools Sensors) ke timezone lain (contoh: America/New_York atau Pacific/Auckland), reload SETIAP halaman di atas, screenshot ulang -- rentang tanggal default WAJIB tetap sama (WIB), tidak ikut geser
4. Test preset "Hari ini/Minggu ini/Bulan ini" di minimal 2 halaman laporan, verifikasi rentang yang dihasilkan benar
5. Buka halaman Ringkasan Bisnis (Business Summary), cek label sumbu-X grafik "Grafik Penjualan" tidak berubah tanggal setelah ganti timezone device (ini yang sebelumnya jadi temuan di Fase 4)
6. npm run type-check && npm run lint && npm run build harus lulus semua
7. go build ./... harus lulus

Laporkan hasil lengkap per halaman dengan bukti screenshot sebelum/sesudah ganti timezone device. Jangan lanjut ke fase lain.
```

---

## FASE 6 — Auth/Session 🟡 P2

### Scope
- BE: `domain/auth/repo/auth_repo.go` baris 8 — `created_at = NOW()` pada upsert sesi → ganti jadi parameter WIB dari Go, selaraskan dengan `expires_at` yang sudah WIB via `auth_service.go`

### Testing
- Login, cek row `sessions` di DB — `created_at` dan `expires_at` harus konsisten satu clock (selisihnya = `TokenExpire` config, bukan ganjil karena beda clock)

### PROMPT EKSEKUSI — Fase 6
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-5 sudah selesai.

Kerjakan FASE 6 (Auth/Session) dari rencana migrasi timezone.

BACKEND — BE/domain/auth/repo/auth_repo.go dan auth_service.go:
1. Baris created_at=NOW() pada query upsert sesi (referensi awal baris 8, baca ulang karena nomor bisa berubah): ganti jadi terima parameter waktu WIB dari Go
2. Pastikan created_at dan expires_at pada row yang sama dihitung dari SATU pemanggilan GetTimeNow() yang sama di service layer, bukan dipanggil terpisah

TESTING WAJIB:
1. Jalankan BE, login via API
2. Cek langsung ke database tabel sessions (atau tabel yang relevan) untuk user yang baru login -- verifikasi created_at dan expires_at selisihnya sama persis dengan TokenExpire di config (dalam hitungan detik/menit yang wajar, bukan ganjil karena beda timezone)
3. go build ./... harus lulus

Laporkan hasil dengan bukti query database. Jangan lanjut ke fase lain.
```

---

## FASE 7 — Pembelian/Retur/Piutang 🟡 P2

### Scope
- BE: `supplier_purchase/repo/purchase_repo.go`, `supplier_return/repo/supplier_return_repo.go` + `service/supplier_return_service.go` — kode PO/retur & validasi tanggal
- FE: `purchases.schema.ts`, `returns.schema.ts`, `PurchaseFormModal.tsx`, `ReturnFormModal.tsx`, `ReceivableTableColumns.tsx` (`isOverdue`)

### Testing
- Browser: input pembelian/retur dengan tanggal mendekati batas validasi (hari ini), ganti timezone device, pastikan validasi tidak salah tolak/terima
- Cek badge "jatuh tempo" piutang konsisten meski timezone device diganti

### PROMPT EKSEKUSI — Fase 7
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-6 sudah selesai.

Kerjakan FASE 7 (Pembelian/Retur/Piutang) dari rencana migrasi timezone.

BACKEND:
1. BE/domain/supplier_purchase/repo/purchase_repo.go: ganti time.Now() (kode PO, payment_date default) jadi time_helper.GetTimeNow()
2. BE/domain/supplier_return/repo/supplier_return_repo.go dan service/supplier_return_service.go: ganti time.Now() (kode retur, validasi returnDate.After(...)) jadi time_helper.GetTimeNow()

FRONTEND:
1. FE/src/features/procurement/purchases/purchases.schema.ts: cap validasi todayStr() -- ini otomatis sudah benar kalau Fase 5 sudah selesai (date.ts sudah WIB-aware), tapi verifikasi ulang
2. FE/src/features/procurement/returns/returns.schema.ts: module-level `today = new Date().toISOString().slice(0,10)` -- ganti ke helper WIB dari shared/utils/date.ts, JANGAN pakai toISOString()
3. FE/src/features/customers/receivables/components/ReceivableTableColumns.tsx: isOverdue() ganti new Date() jadi getWIBNow()

TESTING WAJIB (browser):
1. Jalankan BE dan FE, login
2. Buka form Pembelian dan Retur, coba input tanggal = hari ini (WIB), pastikan tidak ditolak validasi
3. Ganti timezone browser ke yang lebih lambat dari WIB (misal America/Los_Angeles, yang jauh di belakang), reload form, ulangi input tanggal hari ini WIB -- pastikan validasi TETAP menerima (sebelumnya berisiko ditolak karena "besok" menurut device)
4. Buka halaman Piutang, cek badge jatuh tempo tidak berubah setelah ganti timezone device
5. npm run type-check && npm run lint harus lulus, go build ./... harus lulus

Laporkan hasil dengan bukti screenshot. Jangan lanjut ke fase lain.
```

---

## FASE 8 — Sisanya (sapu bersih) 🟢 P3

### Scope
Semua file `updated_at = NOW()` generic yang belum tersentuh di fase 1-7 (lihat daftar lengkap di `TIMEZONE_AUDIT.md` §2.2): `access_repo.go`, `customer_repo.go`, `menu_repo.go`, `product_repo.go`, `product_category_repo.go`, `product_unit/unit_repo.go`, `product_package_repo.go`, `receivable_repo.go`, `role_repo.go`, `setting_repo.go`, `shift_repo.go`, `supplier_repo.go`, `sync_repo.go`, `user_repo.go`, plus sisa file `time.Now()` mentah di §2.1 yang belum kesentuh (`backup_service.go`, `sync_service.go`, `login_rate_limit_middleware.go`, `log_scheduler_repository.go`, `helper/global.go`, `log_request_middleware.go`, `janitor.go`).

### Testing
- Regression test menyeluruh: jalankan seluruh test suite BE (`go test ./...`) dan FE (`npm run type-check && npm run lint && npm run build`)
- Smoke test browser di semua menu utama aplikasi (checklist manual singkat)

### PROMPT EKSEKUSI — Fase 8
```
Baca docs/TIMEZONE_AUDIT.md dan docs/TIMEZONE_MIGRATION_PLAN.md di project d:\Develop\Project_pos_mahenz. Fase 0-7 sudah selesai, ini fase terakhir untuk menyapu bersih sisa file yang belum tersentuh.

Kerjakan FASE 8 (Sapu Bersih) dari rencana migrasi timezone.

Daftar file yang HARUS dicek dan diperbaiki (ganti NOW() SQL jadi parameter WIB dari Go, ganti time.Now() mentah jadi time_helper.GetTimeNow()):
- BE/domain/access/repo/access_repo.go
- BE/domain/customer/repo/customer_repo.go
- BE/domain/menu/repo/menu_repo.go
- BE/domain/product/repo/product_repo.go
- BE/domain/product_category/repo/product_category_repo.go
- BE/domain/product_unit/repo/unit_repo.go
- BE/domain/product/repo/product_package_repo.go
- BE/domain/receivable/repo/receivable_repo.go
- BE/domain/role/repo/role_repo.go
- BE/domain/setting/repo/setting_repo.go
- BE/domain/shift/repo/shift_repo.go
- BE/domain/supplier/repo/supplier_repo.go
- BE/domain/sync/repo/sync_repo.go dan service/sync_service.go
- BE/domain/user/repo/user_repo.go
- BE/domain/backup/service/backup_service.go
- BE/middleware/login_rate_limit_middleware.go
- BE/repository/log_scheduler_repository.go
- BE/helper/global.go
- BE/middleware/log_request_middleware.go (masih ada time.Now() raw untuk durasi, cek apakah perlu diubah atau memang boleh tetap relatif)
- BE/pkg/janitor/janitor.go (cek apakah ini relative-duration only, kalau iya boleh tetap time.Now())

Untuk setiap file, gunakan judgment: kalau time.Now()/NOW() dipakai untuk MENGUKUR DURASI/SELISIH RELATIF (misal TTL cache, request duration), boleh dibiarkan karena tidak sensitif terhadap timezone absolut. Kalau dipakai untuk MENYIMPAN/MEMBANDINGKAN TANGGAL KALENDER (created_at, updated_at, tanggal transaksi/dokumen), WAJIB diganti ke WIB eksplisit.

TESTING WAJIB:
1. go build ./... dan go test ./... harus lulus semua
2. npm run type-check && npm run lint && npm run build harus lulus di FE
3. Jalankan BE+FE, lakukan smoke test manual lewat browser: login, buka minimal 10 menu berbeda (Dashboard, Kasir, Produk, Supplier, Pembelian, Retur, Pelanggan, Piutang, Kas Harian, Laporan), screenshot masing-masing pastikan tidak ada error/crash
4. Cek BE/logs/ tidak ada error baru muncul setelah perubahan

Laporkan hasil lengkap dengan checklist semua menu yang sudah dites. Ini fase terakhir -- setelah ini seluruh audit di TIMEZONE_AUDIT.md dianggap selesai ditangani.
```

---

## Checklist Akhir (setelah Fase 8 selesai)

- [x] Tidak ada lagi `time.Now()` mentah untuk tanggal kalender di BE (kecuali kasus relative-duration yang sudah di-review sengaja dipertahankan)
- [x] Tidak ada lagi `NOW()`/`CURDATE()` di SQL untuk perbandingan/logic bisnis (boleh tetap di `DEFAULT CURRENT_TIMESTAMP` skema sebagai fallback)
- [x] Tidak ada lagi `new Date()`/`Date.now()` mentah di FE untuk kebutuhan bisnis (kosmetik seperti nama file export boleh dikecualikan)
- [x] Semua 19 file konsumen `date.ts` otomatis WIB-consistent karena helper pusatnya sudah diperbaiki
- [x] `go build ./...` lulus. `go test ./...` — 2 kegagalan pre-existing tidak terkait (dikonfirmasi tidak ada perubahan di file terkait): `pkg/binder` (bug overflow angka besar) dan `pos_api/helper` (masalah cwd `.env` saat test, bukan regresi)
- [x] `npm run type-check`, `npm run lint`, `npm run build` lulus
- [x] Smoke test manual browser di seluruh menu utama tidak menemukan regresi (23 halaman dicek di Fase 8)
- [x] Reproduksi kasus bug asli (kas "kebawa nginap") sudah tidak terjadi lagi — dibuktikan via API dan browser

---

## ✅ Ringkasan Penutup — Migrasi Selesai (9 Agustus 2026)

Semua 9 fase (Fase 0–8) sudah dikerjakan dan diverifikasi, masing-masing dengan testing browser nyata (bukan cuma build/lint):

| Fase | Domain | Testing browser |
|---|---|---|
| 0 | Persiapan (helper WIB BE + dayjs FE) | — (fondasi, tidak ada UI) |
| 1 | Kas (cash_drawer) — **bug asli** | ✅ Reproduksi bug asli via API + screenshot Dashboard/Kas Saya/Kas Harian sebelum-sesudah fix |
| 2 | Transaksi | ✅ Buat transaksi berturutan via API, cek kode & tanggal konsisten |
| 3 | Dashboard | ✅ Screenshot header lintas timezone (Jakarta vs New York) |
| 4 | Business Summary | ✅ Screenshot + response API lintas timezone, termasuk fix grafik |
| 5 | Laporan + refactor pusat `date.ts` | ✅ 9 halaman discreenshot lintas timezone (Jakarta vs Auckland), termasuk seed data manual untuk halaman yang butuh prasyarat |
| 6 | Auth/Session | ✅ Query database langsung, verifikasi selisih `created_at`/`expires_at` |
| 7 | Pembelian/Retur/Piutang | ✅ Validasi tanggal via API + badge overdue lintas 3 timezone |
| 8 | Sapu bersih seluruh sisa file | ✅ Smoke test 23 halaman + spot-check update produk/customer/unit |

**Root cause asli** (`getMyCashQuery` punya syarat `DATE(open_time) = CURDATE()` yang tidak ada di `getCurrentCashDrawerQuery`) — sudah diperbaiki di Fase 1, dan seluruh pola serupa di aplikasi (campuran jam MySQL/Go raw/browser) sudah disapu bersih di fase-fase berikutnya.

Dokumen ini dan `TIMEZONE_AUDIT.md` bisa dijadikan referensi historis — tidak perlu tindakan lebih lanjut kecuali ditemukan regresi baru.
