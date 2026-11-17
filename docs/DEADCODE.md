# Daftar Dead Code

Hasil analisis. Item bertanda **[DIHAPUS]** sudah dieksekusi; sisanya belum. Tiap baris diberi saran: **HAPUS** (pasti tidak dipakai), **OPSIONAL**
(hanya merapikan, tidak mengubah perilaku), **PERTAHANKAN** (alat salah tebak atau sengaja disimpan), **PUTUSKAN** (butuh keputusan user).
Urutan aplikasi: Coffee → Portal → BE User → Toko.

Alat: Go `deadcode` (fungsi tak terjangkau dari main), `staticcheck -checks U1000` (const/var/field), analisis tipe sendiri
untuk method interface (`testing/deadcode/ifacecheck`), pencocokan route BE dengan alamat API di FE (`testing/deadcode/endpoints.js`),
FE `knip` (file, export, library). Setiap temuan dicek ulang di kode; salah tebak alat sudah dicoret.

---

## Coffee (dianalisis 15 Nov 2026)

**Cara memilah (disepakati user):** tiap temuan dibandingkan dengan kode yang berjalan.
1. Benar-benar tidak berguna → hapus langsung (dieksekusi 15 Nov 2026).
2. Fungsinya sudah ada tetapi kode yang berjalan menulis ulang logika yang sama → diskusi (mungkin dipakai, bukan dihapus).
3. Seharusnya dipakai tetapi terlewat → diskusi (mungkin disambungkan).
4. Mirip tetapi kode berjalan sudah lebih baik → diskusi (kemungkinan dihapus).

**Eksekusi kemungkinan 1 (15 Nov 2026):** BE `pkg/redis` + `RedisConfig` + komentar Redis di `main.go`, `pkg/encryptor` (+uji),
`helper/security` (semua repo sudah memakai daftar kolom urutan yang diizinkan), `middleware/auth/basic_auth_middleware.go` +
`BuildBasicAuthCreds` (aplikasi memakai JWT), paket model payment_method/payment_status/stock_mutation, `GetDB()` di 9 repo yang
tidak pernah dipanggil (6 repo lain tetap); FE `ui/dropdown-menu.tsx`, library `@radix-ui/react-dropdown-menu`,
`@radix-ui/react-toast`, `@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import`.
Uji: build/vet, `go test ./...` 12 paket ok; FE type-check/lint/build; regresi API Coffee 4.3-4.10 semua lolos
(33, 40, 38, 17, 49, 36, 17, 15); e2e browser 4.11 22/22. Temuan `npm audit` tetap 27 (tiga alat lint itu masih terpasang
sebagai turunan `typescript-eslint`).

### BE-A. Paket yang sama sekali tidak diimpor aplikasi — HAPUS seluruh folder

| Paket | Isi | Catatan |
|---|---|---|
| `domain/payment_method/model` **[DIHAPUS]** | model | handler/service/repo payment_method tetap dipakai |
| `domain/payment_status/model` **[DIHAPUS]** | model | idem |
| `domain/stock_mutation/model` **[DIHAPUS]** | model | repo stock_mutation tetap dipakai |
| `helper/regexp` **[DIHAPUS]** | 8 fungsi validasi teks | bukan duplikat: validasi berjalan memakai tag validator dan regex lain di `validation`; `IsAccountLoan` format bank proyek lain |
| `helper/security` **[DIHAPUS]** | sanitasi SQL, validasi header | repo sudah memakai whitelist kolom urutan |
| `pkg/encryptor` **[DIHAPUS]** | enkripsi AES | ikut file uji |
| `pkg/redis` **[DIHAPUS]** | pengelola Redis | `RedisConfig` dan komentar Redis di `main.go` ikut dibuang |

Semua item BE-A sudah dihapus.

### BE-B. Fungsi tak terpanggil di paket yang masih dipakai — **[SEMUA DIHAPUS]** (kelompok 3)

**Kelompok 3 (fungsi bantu BE), diputuskan AI atas serahan user, 16 Nov 2026:** tidak ada yang layak dipakai; semua
peninggalan template proyek lain (faktur pajak, BRI/Brifaktur, ESB/Brigate, header PN, terbilang) atau tergantikan cara
yang lebih tepat di kode berjalan (kemungkinan 4): format rupiah struk `rupiah()` angka bulat; kode transaksi/PO dari urutan
database; `IN ?` GORM; middleware membaca sesi dan header Bearer sendiri; batas hari sebagai teks SQL; handler laporan
memakai `DefaultQuery` (nilai bawaan, `BindQuery` akan menolak parameter kosong); nama satuan dasar dari `packagesQuery`.
Dihapus 48 fungsi + 16 fungsi uji yang hanya menguji fungsi itu (`TestBindQuery_*`, `TestBindMultipartForm_*`,
`TestHashPasswordWithParams`) + paket `helper/regexp` + const `baseUnitNameQuery`; sisa variabel/konstanta global yang hanya
dipakai fungsi tersebut dibersihkan berdasarkan `staticcheck` (penghitung ID, regex desimal, format invoice BRI, daftar angka
terbilang, konstanta log ESB, peta klaim `rulesVerficationMap`). Diff: 13 file, 1.537 baris dihapus, 0 baris ditambah.
Catatan: middleware menyimpan `device_info` ke context tetapi tidak ada yang membacanya; dibiarkan (murah, tidak berbahaya).
Uji: `go build`/`go vet` (termasuk file uji), `deadcode` kosong, `staticcheck` hanya menyisakan `field private` di uji binder
(sudah ada sejak awal), `go test ./...` 11 paket ok; Coffee BE dijalankan ulang dengan kode baru: regresi API 4.3-4.10 lolos semua
(251 kasus), e2e 4.11 22/22, SSO Portal ke Coffee + token palsu/setelah logout ditolak; log BE tanpa 500/panic.

Daftar asal (semua kini dihapus):

| File | Fungsi |
|---|---|
| `helper/global.go` (29) | GenerateNomorInvoice, GetClaims, GetClaimsId, ConvertDataToJson, ConvertStringToMap, ReadhttpHeader, ParseUserHeader, GenerateCustomID, GeneratePlaceholderForIds, WrapOtherArgumentsWithIds, GenerateTimeBasedID, generateSecureTimeBasedIDFallback, PrettyPrintJSON, MapResponseBody, ~~BuildBasicAuthCreds~~ [DIHAPUS], GenerateExternalId, GenerateReferenceNumber, GetMaxBackdateBrifaktur, FormatRupiah, FormatRupiahString, ParseCommaSeparatedToMap, GenerateNomorFaktur, TerbilangBigRat, terbilangBigInt, FormatDecimalRupiah, DecimalToPercentage, IndoToMySQLNumber, NormalizeTwoDecimal, WriteCSVWithGzipResponse |
| `helper/time/time.go` (9) | StartOfDay, EndOfDay, ConvertToDateFormatString, GenerateMasaPajak, GetIndonesianDayNameAndDate, getIndonesianDayName, formatIndonesianDate, formatTwoDigits, formatFourDigits |
| `helper/log/log.go` (3) | LogExternalCall, LogExternalCallWithResponse, IsDebugMode |
| `helper/request/request.go` (3) | ConvertParamStringToInt, ConvertQueryValueToInt, GetTokenParts |
| `helper/context_helper.go` (1) | GetDeviceInfo |
| `middleware/auth/basic_auth_middleware.go` **[DIHAPUS]** | BasicAuthMiddleware dan turunannya |
| `pkg/binder/binder.go` (2) | BindQuery, BindMultipartForm (ikut ±32 kasus uji di `binder_test.go`) |
| `pkg/jwt/jwt.go` (1) | FillResultMapFromClaims |
| `pkg/bcrypt/bcrypt.go` (1) | HashPasswordWithParams (ikut 5 kasus uji) |
| `domain/stock_reconciliation/repo/stock_reconciliation_repo.go` | const `baseUnitNameQuery` |

### BE-C. Method interface yang dideklarasikan + diimplementasi tetapi tidak pernah dipanggil — **[SEMUA DIHAPUS]** (kelompok 4)

**Kelompok 4, diputuskan AI atas serahan user, 16 Nov 2026:** `GetTodayStats`/`GetTodayExpenses` + 2 query (kemungkinan 4:
identik dengan versi rentang yang dipakai layanan untuk "hari ini"); `GetLowStock` service + repo + query + DTO
`GetLowStockResponse` + model `LowStockProduct` (kemungkinan 4: filter "Stok Menipis" dan kartu Ringkasan Bisnis sudah memakai
`BuildStockSummaries` yang sama, satuan terkecil; `GetLowStock` tidak pernah punya endpoint); **seluruh fitur PIN**
(`domain/pin`, `routes/segment/pin_routes.go`, pendaftaran route) karena tidak pernah dipakai FE Coffee maupun FE Toko (asal
salinan), tidak ada di data izin menu, 0 user punya PIN; kolom `users.pin_hash` di database DIBIARKAN (tanpa migrasi).
Efek berantai: `pkg/bcrypt` (HashPassword/VerifyPassword) hanya dipakai PIN, jadi seluruh paket + ujinya ikut dihapus
(verifikasi password ada di BE User). Uji: build/vet, `ifacecheck` 0, `deadcode` kosong, `staticcheck` hanya temuan lama
`field private`, `go test` 10 paket ok; BE dijalankan ulang: kartu Ringkasan Bisnis (habis+menipis) = filter "Stok Menipis"
(42 = 42, semua tersaring berstatus menipis); alamat PIN kini "Route not found" seperti alamat tak dikenal lainnya;
regresi API 4.3-4.10 (251 kasus) lolos, `b_4_10` 6/6, `b_4_11` 22/22, log BE tanpa 500/panic.

Daftar asal:

| Interface | Method |
|---|---|
| `business_summary/repo` | GetTodayStats, GetTodayExpenses |
| `pin/repo` | ClearPinHash |
| `product/service` | GetLowStock (beserta `repo.GetLowStock` yang hanya dipakai olehnya) |
| `cash_drawer`, `customer`, `finance`, `product_category`, `product_unit`, `shift`, `stock_mutation`, `stock_reconciliation`, `supplier` (repo) | GetDB (9 interface) **[DIHAPUS]** |

### BE-D. Endpoint API yang tidak dipanggil FE Coffee — **[SELESAI]** (kelompok 5)

**Kelompok 5, diputuskan AI atas serahan user, 16 Nov 2026.** Prinsip: kemungkinan 1/4 dihapus; kemungkinan 3 yang terkait
keputusan resmi tidak dihapus dan tidak dibuatkan fitur (dilaporkan ke user). Riwayat git dipakai untuk memastikan asal-usul:
sebagian besar endpoint tidak pernah dipanggil FE mana pun sejak dibuat (template API awal); `summary-extra` sengaja dihapus
dari FE oleh user (commit "standart", 28 Jun 2026); `by-barcode` dipakai kasir Toko lalu ditinggalkan saat kasir Coffee ditulis
ulang (4.7b).

- **Dipertahankan:** `GET /health` (cek server SOP); `GET /reports/business-summary/payment-methods` (**kemungkinan 3**:
  keputusan 117 menulis "metode bayar tetap" di Ringkasan Bisnis, tetapi tidak pernah ditampilkan FE mana pun; kini hanya ada di
  Laporan Penjualan tab Metode Bayar); `POST /cash-drawer/auto-close-yesterday` (**kemungkinan 3**: dibuat sengaja di 4.9 untuk
  admin, belum ada tombolnya; penutupan otomatis tengah malam tetap berjalan tanpa endpoint ini).
- **Dihapus 23 endpoint** (route + handler + method service/repo yang hanya dipakai olehnya + query/DTO/tipe sisa):
  verify-token, auth/me, summary-extra, sales/chart, cash-drawer update-sales/update-expenses (berbahaya: mengubah total kas
  manual padahal total dihitung otomatis), customers/active, detail expenses/ingredients/modifier-groups/categories/units/users,
  receivables summary/payments, settings/reset, shifts summary/detail, stock-mutations list/product (seluruh domain
  `stock_mutation`; tabelnya tetap ditulis domain lain), menus/reorder, products by-barcode dan :id/prices/list.
- **FE:** hook `useGetMeQuery`, `useReorderMenuMutation`, `useShiftDetailQuery`, `useProductBarcodeQuery`,
  `useProductPricesQuery` + `fetchProductPrices`, tipe `ReorderMenuPayload`, `SummaryExtraResponse`.
- **Alat:** `cascade.js` (ulang deadcode + ifacecheck sampai kosong), `remove_consts.js`, `unused_types.js` +
  `remove_types.js` (tipe DTO/model ekspor yang tidak terdeteksi staticcheck, termasuk 3 yang sudah mati sejak sebelumnya:
  SalesTrendFilter, ReceivablePayment, StockMutation di model transaksi), `fe_to_be.js` (setiap alamat API FE harus punya route).
- **Uji:** build/vet, deadcode kosong, ifacecheck 0, tipe tak terpakai 0, staticcheck hanya temuan lama; FE type-check/lint/build,
  knip tanpa export mati; endpoint mati tersisa tepat 3 yang dipertahankan; semua alamat API FE punya route; regresi API Coffee
  4.3-4.10 lolos (251 kasus; 4 bacaan ulang di skrip uji api_4_3/api_4_4 dialihkan dari endpoint detail ke endpoint daftar);
  36 halaman dibuka tanpa error dan tanpa panggilan "Route not found"; e2e 22/22 dan kas tertunda 6/6 di laptop/tablet/HP,
  open bill 11/11, laporan 6/6, kelompok 2 7/7; log BE tanpa 500/panic.

Daftar asal:

Setelah salah tebak dicoret (Kategori Menu/Ukuran/Meja dipanggil lewat `createMasterApi`; `GET /settings/:key` dipakai
`/settings/pagination_sizes`). Menghapus endpoint = route + handler + method service/repo yang hanya dipakai olehnya.

| Endpoint | Saran |
|---|---|
| `GET /health` | PERTAHANKAN (dipakai cek server di SOP/cutover) |
| `POST /auth/verify-token` | HAPUS |
| `POST /auth/me` (hanya dipakai hook FE yang mati) | HAPUS |
| `GET /reports/business-summary/payment-methods`, `GET /reports/business-summary/summary-extra` | HAPUS |
| `GET /reports/sales/chart` | HAPUS |
| `POST /cash-drawer/update-sales/:id`, `POST /cash-drawer/update-expenses/:id` | HAPUS (sisa Toko; total kas dihitung otomatis) |
| `POST /cash-drawer/auto-close-yesterday` | **PUTUSKAN** (dibuat 4.9 untuk tutup kas kemarin oleh admin, tetapi belum ada tombolnya di FE) |
| `POST /customers/active` | HAPUS |
| `GET /pin/check`, `POST /pin/set`, `/pin/verify`, `/pin/change` | **[DIHAPUS]** bersama seluruh fitur PIN (kelompok 4) |
| `POST /receivables/summary`, `POST /receivables/payments/:id` | HAPUS |
| `POST /settings/reset` | HAPUS |
| `POST /shifts/summary`, `POST /shifts/detail/:id` | HAPUS |
| `POST /stock-mutations/list`, `POST /stock-mutations/product/:product_id` | HAPUS (rute saja; repo stock_mutation tetap dipakai) |
| `POST /menus/reorder` | HAPUS |
| `POST /products/by-barcode/:barcode`, `POST /products/:id/prices/list` | HAPUS |
| `POST /expenses/detail/:id`, `/ingredients/detail/:id`, `/modifier-groups/detail/:id`, `/categories/detail/:id`, `/units/detail/:id`, `/users/detail/:id` | **PUTUSKAN** (FE memakai data dari daftar; endpoint detail standar, bisa dipertahankan untuk kelengkapan API) |

### FE-E. File tidak terpakai

| File | Saran |
|---|---|
| `shared/components/ui/dropdown-menu.tsx` | **[DIHAPUS]** (tidak ada menu tarik-turun di aplikasi) |
| `shared/components/ui/password-input.tsx` | **[DIPAKAI]** kemungkinan 2: isinya identik dengan tombol mata yang ditulis ulang di `LoginForm`; kini `LoginForm` memakai `PasswordInput` (sama seperti Portal) |
| `shared/components/ui/radio-group.tsx` | **[DIHAPUS]** kemungkinan 4: bulatan radio 16 px, sedangkan kode berjalan sengaja memakai tombol segmen besar untuk tablet |
| `shared/types/webBluetooth.d.ts` | PERTAHANKAN (tipe `navigator.bluetooth` untuk printer Bluetooth; alat salah tebak) |
| 15 barrel `index.ts` tak diimpor (customers/receivables, dashboard, finance/overview, menu-items, procurement/purchases, reporting/{business-summary, cashier-performance, profit-loss, sales, stock}, sales/transactions, settings, settings/{backup, menus, users}) | OPSIONAL (hanya pengumpul export, tidak ikut ke hasil build) |

**Kelompok 1 (komponen UI), diputuskan AI atas serahan user, 15 Nov 2026:** password-input dipakai, radio-group dihapus,
dan temuan sampingan disatukan: tombol [Kembali ke stok | Dibuang] yang ditulis dua kali (`BillDialogs` dan
`VoidTransactionDialog`) kini satu komponen `features/sales/cashier/components/StockChoice.tsx`. Uji: type-check/lint/build;
uji login baru `testing/deadcode/b_login_password.js` 5/5 (password tersembunyi, tombol mata, Enter pindah fokus, password salah
ditolak, login berhasil) di laptop/tablet/HP; `b_4_7b` 7/7 (void dibuang), `b_4_8` 11/11 (hapus item dibuang, batal) dan
`b_4_11` 22/22 (void kembali ke stok) di laptop/tablet/HP.

### FE-F. Library npm tidak terpakai

| Library | Saran |
|---|---|
| `@radix-ui/react-dropdown-menu`, `@radix-ui/react-toast` | **[DIHAPUS]** |
| `@radix-ui/react-radio-group` | **[DIHAPUS]** |
| devDependencies `@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import` | **[DIHAPUS]** (lint memakai paket `typescript-eslint`) |
| `lint-staged` | PERTAHANKAN (dipakai hook Husky `.husky/pre-commit`) |

### FE-G. Fungsi/hook diekspor yang tidak dipakai di mana pun

**Kelompok 2 (format, hitung, utilitas FE), diputuskan AI atas serahan user, 15 Nov 2026:**

| Fungsi | Kemungkinan | Hasil |
|---|---|---|
| `formatProductPackage` | 2 (rumus "Satuan (kemasan)" ditulis ulang di 3 tempat) | **[DIPAKAI]** di `PurchaseFormModal`, `ProductDetailModal`, `ProductFormModal` |
| `isPaymentSufficient` | 4 (modal bayar lebih lengkap: utang, total 0) | **[DIHAPUS]** |
| `percentage`, `roundTo` | 4 (laporan sengaja membulatkan ke angka bulat) | **[DIHAPUS]** |
| `parseRupiah` | 4 (`RupiahInput` sengaja angka bulat) | **[DIHAPUS]** |
| `PAYMENT_STATUS_LABEL` | 4 (label sama terpusat di `StatusBadge`) | **[DIHAPUS]** |
| `usePermission` | 4 (`useAuth` sudah menyediakan isOwner/isAdmin/isKasir) | **[DIHAPUS]** |
| `EmptyState` | 4 (hampir identik dengan `DataTableEmpty` yang dipakai) | **[DIHAPUS]** |
| `getPlatform` + tipe `DevicePlatform` | 4 (`isAndroid` sudah cukup) | **[DIHAPUS]** |
| `clamp`, `truncate`, `capitalize`, `slugify`, `formatRelative`, `useLocalStorage`, `getApplicablePrice` | 1 | **[DIHAPUS]** (`string.ts` dihapus) |
| `openReceiptPdf` | 3 (pratinjau PDF di tab baru; sudah ada Simpan/Bagikan PDF, tab baru sering diblokir di WebView APK) | **[DIHAPUS]** |

Temuan sampingan disatukan: `pctChange` (persen naik/turun) yang identik di `SummaryCards` dan `SalesReportSummaryCard`
kini satu di `shared/utils/number.ts`; `CashDrawerTable` memakai `useAuth` (isOwner/isAdmin) untuk tombol Tutup Paksa.
Uji: type-check/lint/build; knip ulang hanya menyisakan 5 hook API (kelompok 5); uji baru `testing/deadcode/b_kelompok2.js`
7/7 di laptop/tablet/HP (label "Dus (Dus)" di detail produk, form produk, pilihan satuan pembelian "Dus (Dus) (x24)";
Ringkasan Bisnis dan Laporan Penjualan tampil; tombol Tutup Paksa untuk admin; Simpan PDF struk menghasilkan file .pdf);
regresi `b_4_10` 6/6, `b_4_11` 22/22 di 3 perangkat, `b_4_5` bagian Coffee 5/5 (kata cari di skrip uji diganti nama lengkap
karena data uji "Croissant" menumpuk; bagian Toko butuh FE Toko yang tidak dinyalakan).

Sisa kelompok ini (hook API, dibahas bersama endpoint di kelompok 5):

`auth.api.ts` useGetMeQuery; `menu.api.ts` useReorderMenuMutation; `shifts.api.ts` useShiftDetailQuery;
`products.api.ts` useProductBarcodeQuery, useProductPricesQuery (+ `fetchProductPrices` yang hanya dipakai olehnya).

### FE-H. Hanya kata `export` yang berlebih (dipakai di file sendiri) — OPSIONAL

formatTime (shifts.utils), fetchProductPackages, printViaBle, receiptToPlainText, isAndroid, calcGrandTotal, OPEN_BILLS_KEY,
buildReceiptPdf, buildReceiptText, authClient, useMediaQuery, roundTo.

### FE-I. Lainnya — PERTAHANKAN / OPSIONAL

- 43 tipe TypeScript diekspor tetapi tidak dipakai: OPSIONAL (tidak ikut ke hasil build).
- Sub-export komponen shadcn (`DialogTrigger`, `SelectGroup`, `CardFooter`, dll.): PERTAHANKAN (pustaka komponen standar).


**Susulan 17 Nov 2026 — 6 tipe sisa [DIHAPUS]** (ditemukan setelah pemeriksa tipe diperbaiki saat analisis Toko; pemeriksa lama tidak
membaca tipe di blok `type ( ... )`): `GetByIDRequest` di dto expense, product_category, product_unit, shift, user (milik endpoint
detail yang dihapus di kelompok 5) dan model `transaction.Transaction` (kode memakai DTO transaksi; migrasi berupa file SQL, tanpa
AutoMigrate/refleksi). Diff: 6 file, -45 baris. Uji: build/vet, `go test` + uji DB produk dan nomor antrean, `deadcode` kosong,
`ifacecheck` 0, tipe tak terpakai 0 (`unused_types_pkg.js`), `staticcheck` hanya 2 temuan lama, gofmt hanya file yang rapi di HEAD;
BE dijalankan ulang: regresi API Coffee 4.3-4.10 lolos 251/251 (33, 40, 38, 17, 49, 6, 36, 17, 15); log BE tanpa 500/panic.

---

## Portal (dianalisis 16 Nov 2026) — **SELESAI** (kelompok 1 BE, 2 FE, 3 kemungkinan 3)

**Keputusan user 16 Nov 2026:** semua saran HAPUS/PAKAI disetujui; `/auth/me` dipakai (pilihan b); validasi Login/Refresh
dijalankan (pilihan b). Eksekusi 3 kelompok: 1 BE, 2 FE, 3 dua perubahan kemungkinan 3.

**Kelompok 1 (BE Portal) [SELESAI], 16 Nov 2026:** seluruh BE-A, BE-B, BE-C, dan 2 endpoint BE-D dihapus; berantai ikut:
handler `VerifyToken` + method service + DTO `VerifyTokenRequest/Response`, handler `GetByID`, variabel/konstanta sisa
(penghitung ID, regex desimal, format invoice BRI, daftar angka terbilang, konstanta log ESB/Brigate, `rulesVerficationMap`),
36 fungsi uji binder + 3 pembantunya, uji `validator` untuk 12 aturan yang dihapus (`AddCustomErrorMessage` ikut terbuang),
library Go `go-redis` + 3 turunannya (`go mod tidy`). Dipertahankan: `NotFoundError.Error`, `InternalServerError.Error`, `/health`,
`/auth/me`, pesan `nefield` (catatan: DTO Ganti Password juga tidak memakai `nefield`; dibiarkan karena kecil dan bawaan).
Diff: 38 file, 5.184 baris dihapus, 0 ditambah. Uji: `go build`/`go vet`, `go test` 6 paket ok, `deadcode` hanya 2 method error
yang dipertahankan, `ifacecheck` 0, tipe tak terpakai 0, `staticcheck` hanya `field private` (sudah ada sejak awal), gofmt bersih;
endpoint tak dipanggil tinggal `/health`, semua alamat API FE punya route; BE dijalankan ulang: uji API `auth_3b3/api_3b3.js` 55/55
(P11 dialihkan ke `/users/update/abc`, P13 ke `/users/list`, tambah P13b/P13c: dua endpoint terhapus menjawab "Route not found"),
browser laptop `v_portal.js` 40/40 dan SSO ke Toko `v_toko.js` 32/32; log BE Portal tanpa 500/panic.

**Kelompok 2 (FE Portal) [SELESAI], 16 Nov 2026:** seluruh FE-E (13 file ui), FE-F (10 library + `dayjs` + 3 alat lint;
192 paket turunan ikut terbuang, tidak ada library lain yang berubah versi) dan FE-G (komponen ActionModal, ExtendedConfirmDialog,
EmptyState, DetailField, ScrollableTable, SummaryCard; hook usePagination, useDisclosure, useTableSelection, useLocalStorage;
seluruh `date.ts` dan `string.ts`; barrel dirapikan) dihapus. `LoginForm` memakai `PasswordInput` (versi Portal, area sentuh
40 px, sama dengan halaman Ganti Password); ikon gembok diberi `z-10` karena pembungkus `PasswordInput` ber-position relative
dan bisa menutupinya. Perbedaan kecil yang disengaja: tombol mata di login kini sama dengan Ganti Password (warna abu-abu, tidak
ikut urutan Tab). `useGetMeQuery` dibiarkan untuk kelompok 3. Diff `src`: 35 file, 5 baris ditambah, 1.313 dihapus.
Uji: type-check, lint, build; knip tanpa file/library tak terpakai; uji baru `testing/deadcode/b_portal_login_password.js` 9/9 di
laptop/tablet/HP (ikon gembok dan user tidak tertutup, tombol mata 40x40 bekerja, password salah ditolak dengan pesan, Enter
mengirim login, tombol mata di Ganti Password bekerja, tanpa error JavaScript); regresi `v_portal.js` 40/40 di laptop/tablet/HP.
**Temuan tambahan (belum dihapus, menunggu keputusan user):** `ui/tooltip.tsx` hanya dipakai lewat `TooltipProvider` di
`app/providers.tsx`; `Tooltip`/`TooltipContent`/`TooltipTrigger` tidak dipakai di halaman mana pun (kemungkinan 1:
provider + file + library `@radix-ui/react-tooltip` bisa dihapus). Terlewat di analisis awal karena dikira sub-export shadcn.

**Kelompok 3 (kemungkinan 3 + Tooltip) [SELESAI], 16 Nov 2026:**
- Keputusan 129: handler Login dan Refresh menjalankan `validator.Validate.Struct` (pola Ganti Password). Kiriman tidak sah
  (`device_info` selain desktop/web/android, username/password kosong, refresh token kosong) ditolak 422 dengan rincian; kiriman
  ini tidak dihitung batas percobaan login (hanya password salah yang dihitung).
- Keputusan 128: `useGetMeQuery` dipakai di `ProtectedRoute`; hasil `/auth/me` langsung menyegarkan nama, role, dan status wajib
  ganti password di sesi tersimpan (di dalam `queryFn`, tanpa effect). Dipanggil saat Portal dibuka, paling sering sekali per
  5 menit saat pindah bagian (staleTime bawaan); bila gagal, sesi tersimpan tetap dipakai.
- `Tooltip` (temuan tambahan kelompok 2, disetujui user): `TooltipProvider` di `providers.tsx`, `ui/tooltip.tsx`, library
  `@radix-ui/react-tooltip` dihapus (tanpa perubahan versi library lain).
Uji: BE build/vet/test, gofmt; FE type-check/lint/build, knip tanpa file/library tak terpakai dan `useGetMeQuery` tidak lagi
tercatat mati; uji baru `testing/deadcode/b_portal_k3.js` 18/18 di laptop/tablet/HP (V1-V8 validasi Login/Refresh termasuk login
android APK; S1-S5: admin menurunkan role dan mengganti nama, Portal dibuka ulang, `/auth/me` dipanggil, sesi berubah, langsung
keluar dari Kelola User dan menunya hilang tanpa login ulang; role dinaikkan kembali; `/auth/me` diputus, Portal tetap jalan; tanpa
error JavaScript); regresi API Portal 55/55, `v_portal.js` 40/40 dan login 9/9 di laptop/tablet, HP: 40/40 dan login 9/9; SSO Toko 32/32 (laptop);
SSO Coffee `b_4_11` 22/22 (laptop). Catatan uji: Portal hanya mengizinkan satu sesi aktif per akun, jadi skrip uji tidak boleh
login ulang akun yang tokennya masih dipakai.

Cara memilah sama dengan Coffee (kemungkinan 1-4). Latar belakang penting: sejak 3B.3 (keputusan 59) BE Portal hanya
**meneruskan** permintaan ke BE User (kelola user/role/aplikasi, SSO launch) dan hanya mengurus sesi Portal sendiri. Karena itu
banyak alat bantu template (pagination, ORDER BY, binder query/form, pemeriksaan role) tidak lagi punya tempat dipakai.
Hasil alat: `deadcode` 123 fungsi, `ifacecheck` 1 method, tipe DTO/model tak terpakai 0, endpoint tak dipanggil 3, alamat API FE
tanpa route 0; FE `knip` 13 file, 10 library + 3 alat lint, 29 export mati, 8 kata `export` berlebih. Setiap temuan dicek di kode
Portal sendiri (bukan disalin dari Coffee): kode hidup dicari yang menulis ulang logika serupa (JSON, pembuat ID/token, cache,
regex, tanggal, badge/kartu, penyimpanan lokal). Hasil khas Portal yang tidak ada di Coffee: aturan validasi khusus tidak dipakai
DTO mana pun, validasi Login/Refresh tidak dijalankan, dan seluruh `date.ts` + library `dayjs` mati.

### BE-A. Paket/file yang sama sekali tidak diimpor aplikasi

| Paket | Isi | Kemungkinan | Saran |
|---|---|---|---|
| `pkg/redis` + `RedisConfig` + komentar Redis di `main.go` | pengelola Redis | 1 (Portal tidak memakai Redis) | HAPUS (sama Coffee) |
| `pkg/encryptor` (+uji) | enkripsi AES | 1 | HAPUS (sama Coffee) |
| `pkg/janitor` (+uji) | cache LRU memori bersharding | 1 (satu-satunya penyimpanan memori hidup adalah pencatat percobaan login di `login_rate_limit_middleware.go`, tujuannya lain: daftar waktu gagal per IP+username, bukan cache berukuran) | HAPUS |
| `helper/security` | sanitasi SQL, validasi header | 1 (Portal tidak membuat query daftar; header dibaca middleware) | HAPUS (sama Coffee) |
| `helper/regexp` (+uji) | 8 validasi teks | 1 (validasi password/username di BE User) | HAPUS (sama Coffee) |
| `helper/codegen` | kode unik berprefiks | 1 (Portal tidak membuat kode transaksi) | HAPUS |
| `middleware/auth/basic_auth_middleware.go` | Basic Auth | 1 (Portal memakai JWT + sesi) | HAPUS (sama Coffee) |

### BE-B. Fungsi tak terpanggil di paket yang masih dipakai

| File | Fungsi | Kemungkinan | Saran |
|---|---|---|---|
| `helper/global.go` (29) | GenerateNomorInvoice, GetClaims, GetClaimsId, ConvertDataToJson, ConvertStringToMap, ReadhttpHeader, ParseUserHeader, GenerateCustomID, GeneratePlaceholderForIds, WrapOtherArgumentsWithIds, GenerateTimeBasedID (+fallback), PrettyPrintJSON, MapResponseBody, BuildBasicAuthCreds, GenerateExternalId, GenerateReferenceNumber, GetMaxBackdateBrifaktur, FormatRupiah, FormatRupiahString, ParseCommaSeparatedToMap, GenerateNomorFaktur, TerbilangBigRat (+BigInt), FormatDecimalRupiah, DecimalToPercentage, IndoToMySQLNumber, NormalizeTwoDecimal, WriteCSVWithGzipResponse | 1 (peninggalan template proyek lain: faktur pajak, BRI, terbilang); pembuat ID/kode acak kemungkinan 4: refresh token di `auth_service` memakai 32 byte acak, lebih kuat dari GenerateExternalId (6 byte) / GenerateReferenceNumber (jam); ConvertDataToJson tidak dipakai `helper/error` yang menulis `json.Marshal` sendiri karena butuh penanganan galat berbeda | HAPUS |
| `helper/time/time.go` (11) | StartOfDay, EndOfDay, ToSQLDate, ConvertToDateFormatString, NormalizeDateRange, GenerateMasaPajak, GetIndonesianDayNameAndDate + 4 fungsi kecil di dalamnya | 1 (Portal tidak punya laporan/rentang tanggal) | HAPUS |
| `helper/log/log.go` (3) | LogExternalCall, LogExternalCallWithResponse, IsDebugMode | 1 (menulis ke tabel log ESB/Brigate proyek lain; panggilan ke BE User sudah tercatat di log request) | HAPUS (sama Coffee) |
| `helper/request` (4) | BuildOrderClause, NormalizePagination, ConvertQueryValueToInt, GetTokenParts | 1 untuk 3 pertama (urut/halaman diurus BE User); GetTokenParts kemungkinan 4 (middleware dan `extractBearerToken` sudah membaca token, tidak peka huruf besar/kecil) | HAPUS |
| `helper/response` | SetPagination | 1 (data halaman diteruskan apa adanya dari BE User) | HAPUS |
| `helper/context_helper.go` (2) | GetUserRole, GetDeviceInfo | 1 (pemeriksaan role pindah ke BE User, keputusan 59) | HAPUS (sama Coffee) |
| `pkg/binder` (4) | BindURI, BindQuery, BindMultipartForm, setFieldValue (+uji yang hanya menguji fungsi ini) | 1 (Portal hanya menerima JSON; id URL dibaca `helper/forward`) | HAPUS (sama Coffee) |
| `pkg/jwt` | FillResultMapFromClaims (+peta `rulesVerficationMap`) | 1 (klaim pegawai BRI) | HAPUS (sama Coffee) |
| `errors/error_types.go` | NotFoundError.Error, InternalServerError.Error | - | PERTAHANKAN (kelengkapan tipe error yang ditangani `errorHandler`) |
| `validation/validator.go` (tidak terdeteksi alat, didaftarkan lewat nama tag) | 12 aturan khusus + pesannya: maxfilesize, fileextension, alphanumspace, dateformat, date_range_end, date_range_start, required_with, fraction, decimal, alphanumdash, notblank, timeformat (+kasus uji di `validator_test.go`) | 1 (DTO Portal hanya memakai `required`, `oneof`, `min`, `max`, `nefield`; Portal tidak menerima file, tanggal, atau angka desimal) | HAPUS (pesan `nefield` tetap) |

**Temuan khas Portal, kemungkinan 3 (PUTUSKAN): validasi Login dan Refresh tidak dijalankan.** `LoginRequest` (`username`,
`password`, `device_info` harus `desktop`/`web`/`android`) dan `RefreshRequest` punya tag `validate`, tetapi handler hanya
menjalankan `validator.Validate.Struct` di Ganti Password. Dampak sekarang kecil: username/password kosong ditolak BE User,
refresh kosong ditolak karena sesi tidak ditemukan; tetapi `device_info` berisi apa saja (atau kosong) tetap disimpan ke sesi.
FE Portal selalu mengirim `web`. Pilihan: (a) biarkan, (b) jalankan validasi di Login dan Refresh (2 baris per handler, pesan
galat 422 seperti Ganti Password).

### BE-C. Method interface tak terpanggil

| Interface | Method | Kemungkinan | Saran |
|---|---|---|---|
| `repository.LogSchedulerRepository` | InsertLogScheduler (+model `LogSchedulerModel`, baris pendaftaran di `routes/router.go`) | 1 (Portal tidak punya penjadwal) | HAPUS kode; tabel `log_schedulers` di DB dibiarkan (tanpa migrasi) |

### BE-D. Endpoint API

| Endpoint | Saran |
|---|---|
| `GET /health` | PERTAHANKAN (cek server di SOP/cutover) |
| `POST /auth/verify-token` | HAPUS (tidak dipanggil siapa pun; Toko/Coffee memeriksa sesi lewat BE User sejak 3B.4; Coffee sudah menghapus yang sama) |
| `POST /users/detail/:id` | HAPUS (form ubah user memakai data baris daftar; Coffee kelompok 5 menghapus detail yang sama) |
| `POST /auth/me` + hook FE `useGetMeQuery` | **PUTUSKAN** (kemungkinan 3, lihat bawah) |

**`/auth/me` (kemungkinan 3, terkait keputusan 31 dan 59).** Endpoint ini sengaja dibuat: mengambil profil terbaru dari BE User
dan termasuk 3 alamat yang boleh dibuka saat wajib ganti password. Tetapi FE tidak pernah memanggilnya (hook `useGetMeQuery`
tidak dipakai); FE memakai data user yang disimpan saat login. Akibatnya, bila admin mengubah role/nama seseorang, menu di
Portal orang itu baru berubah setelah ia login ulang (hak akses tetap aman karena diperiksa BE User tiap permintaan).
Pilihan: (a) biarkan seperti sekarang, (b) pakai: FE memanggil `/auth/me` saat Portal dibuka untuk menyegarkan nama/role,
(c) hapus endpoint + hook + entri di daftar wajib ganti password.

### FE-E. File tidak terpakai (13)

| File | Kemungkinan | Saran |
|---|---|---|
| `ui/badge.tsx` | 4 (status memakai `StatusBadge`; label role di Kelola User ditulis langsung dengan garis tepi dan tanpa efek hover, sedangkan `Badge` varian `secondary` tanpa garis tepi dan punya efek hover/fokus untuk elemen yang bisa diklik) | HAPUS |
| `ui/card.tsx` | 4 (kartu aplikasi Beranda `AppCard` sengaja berdesain sendiri: sudut besar, ikon, kerangka muat) | HAPUS |
| `ui/combobox.tsx`, `ui/async-combobox.tsx`, `ui/command.tsx`, `ui/popover.tsx` | 1 (pilihan role memakai `Select`; tidak ada pencarian di dalam pilihan) | HAPUS |
| `ui/collapsible.tsx`, `ui/radio-group.tsx`, `ui/separator.tsx`, `ui/switch.tsx`, `ui/tabs.tsx`, `ui/textarea.tsx`, `ui/toggle-group.tsx` | 1 (tidak ada bagian lipat, radio, tab, atau isian panjang; status aktif diubah lewat tombol + konfirmasi) | HAPUS |

### FE-F. Library npm tidak terpakai

`@radix-ui/react-collapsible`, `-popover`, `-radio-group`, `-separator`, `-switch`, `-tabs`, `-toast`, `-toggle-group`, `cmdk`,
`react-is`; **`dayjs`** (berantai: hanya dipakai `date.ts` yang seluruhnya mati, lihat FE-G); alat lint `@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import` (lint memakai
paket `typescript-eslint`). Saran: HAPUS (sama Coffee).

### FE-G. Fungsi/komponen/hook diekspor yang tidak dipakai

| Item | Kemungkinan | Saran |
|---|---|---|
| `LoginForm` menulis ulang tombol mata password | 2 (`ui/password-input.tsx` sudah dipakai halaman Ganti Password) | PAKAI `PasswordInput` di `LoginForm` (sama Coffee) |
| `ActionModal`, `ExtendedConfirmDialog` | 4 (semua modal memakai `FormModal` dan `ConfirmDialog`) | HAPUS |
| `EmptyState` | 4 (hampir identik `DataTableEmpty` yang dipakai di dalam `DataTable`) | HAPUS (sama Coffee) |
| `DetailField`, `ScrollableTable`, `SummaryCard` | 1 (Portal tidak punya halaman detail/ringkasan) | HAPUS |
| `usePagination` | 4 (halaman Kelola User punya paging sendiri yang otomatis kembali ke halaman 1 saat pencarian/ukuran berubah; hook ini tidak) | HAPUS |
| `useDisclosure` | 4 (modal menyimpan user yang dipilih; hook ini hanya buka/tutup) | HAPUS |
| `useTableSelection` | 1 (tidak ada aksi massal) | HAPUS |
| `useLocalStorage` | 1 (sesi disimpan zustand persist) | HAPUS (sama Coffee) |
| **seluruh `date.ts`** (formatDate, formatDateShort, formatDateTime, formatRelative, toISODate, todayStr, monthStart, weekStart, getWIBNow, APP_TIMEZONE) | 1 (tidak ada satu pun halaman Portal yang menampilkan tanggal; formatDate/getWIBNow hanya dipakai fungsi mati di file yang sama; dari `shared/utils` hanya `cn` yang diimpor) | HAPUS file + export di `utils/index.ts` + library `dayjs` |
| `string.ts` (truncate, capitalize, slugify) | 1 | HAPUS file (sama Coffee) |
| `useGetMeQuery` | 3 | **PUTUSKAN** bersama `/auth/me` |

### FE-H. Hanya kata `export` yang berlebih — OPSIONAL

updateAppSchema, apiClient, authClient, useMediaQuery, ApiResponse. (APP_TIMEZONE, getWIBNow, formatDate ikut terhapus bersama `date.ts`.)

### FE-I. Lainnya — PERTAHANKAN / OPSIONAL

- 11 tipe TypeScript diekspor tetapi tidak dipakai: OPSIONAL.
- Sub-export komponen shadcn (`DialogTrigger`, `SelectGroup`, dll.) dan barrel `index.ts`: PERTAHANKAN (sama Coffee).

### Catatan sampingan Portal (bukan dead code, tidak diubah tanpa izin)

- Kotak pesan galat merah yang sama ditulis ulang 6 kali (LoginForm, Ganti Password, form user, reset password, akses user,
  form aplikasi). Bisa disatukan jadi satu komponen; murni merapikan.
- Pencatat percobaan login gagal tidak pernah membuang kunci IP+username yang tidak datang lagi (memori bertambah sangat
  pelan sampai BE dimulai ulang). Risiko kecil; bisa dibersihkan berkala bila diinginkan.
- Skema DB lama (`roles`, `apps`, `user_apps`, `sso_codes`, kolom `password`/`role_id`) sudah dibuang migrasi 006; tidak ada kode
  yang masih memakainya.

---

## BE User (dianalisis 16 Nov 2026) — **SELESAI**

Dianalisis dari kode BE User sendiri. Latar belakang: BE User hanya dipanggil server ke server (Portal, Toko, Coffee) lewat
`127.0.0.1:8080` dengan kunci aplikasi `X-App-Key` (keputusan 63); tidak menerbitkan token/sesi, tidak punya FE, tidak di nginx.
Hasil alat: `deadcode` 119 fungsi, `ifacecheck` 1 method, tipe DTO/model tak terpakai 0 dari 45, `staticcheck` hanya `field private`
(uji binder, sudah ada sejak awal). Pemanggil tiap route dicocokkan dengan `client.Post(...)` di Portal, Toko, Coffee: Toko dan
Coffee hanya memanggil `/login/verify` dan `/sso/exchange`; sisanya dipanggil Portal. Alat ops `cmd/set_app_key` dan
`cmd/internal/migrationdb` dipakai SOP/panduan cutover: PERTAHANKAN.

Temuan khas BE User (tidak ada di Portal/Coffee): seluruh CRUD Role (handler, service, repo, DTO) tidak punya route; middleware
`Cors` + config `CorsAllowOrigins` tidak dipasang; config token (`SecretKey`, `TokenExpire`, `RefreshTokenExpire`) tidak dipakai
karena BE User memang tidak menerbitkan token; `user_repo` menulis ulang logika `BuildOrderClause`; `/users/detail/:id` tercantum
di daftar API resmi rancangan (baris "Kelola akun dan aplikasi") tetapi kini tidak punya pemanggil.

### BE-A. Paket/file yang tidak diimpor aplikasi

| Paket | Kemungkinan | Saran |
|---|---|---|
| `pkg/redis` + `RedisConfig` (+ field `PoolSize`, `MinIdleConns`) + komentar Redis di `main.go` | 1 (tidak memakai Redis) | HAPUS |
| `pkg/encryptor` (+uji) | 1 | HAPUS |
| `pkg/janitor` (+uji) | 1 (pencatat percobaan login di `login_rate_limit_middleware.go` berdiri sendiri, tujuannya lain) | HAPUS |
| `helper/security` | 1 (urutan daftar user memakai daftar kolom yang diizinkan; header dibaca `app_key_middleware`) | HAPUS |
| `helper/regexp` (+uji) | 1 (validasi memakai tag bawaan validator; kode aplikasi diperiksa `codePattern` sendiri di `app_service`) | HAPUS |
| `helper/codegen` | 1 (kode aplikasi diketik admin dan diperiksa pola + `GetByCode`, bukan dibuat otomatis) | HAPUS |
| `middleware/cors_middleware.go` + field config `CorsAllowOrigins` | 1 (tidak pernah dipasang di router; BE User tidak dibuka browser, hanya `127.0.0.1`) | HAPUS (nilai di file config dibiarkan, tidak dibaca lagi) |

### BE-B. Fungsi tak terpanggil di paket yang masih dipakai

| File | Fungsi | Kemungkinan | Saran |
|---|---|---|---|
| `helper/global.go` (30) | 29 fungsi template yang sama dengan Portal (faktur, BRI, terbilang, ID acak, JSON) + **GetSecretKey** | 1; GetSecretKey: BE User tidak menandatangani token | HAPUS |
| `config` | field `SecretKey`, `TokenExpire`, `RefreshTokenExpire` | 1 (keputusan 63: BE User tanpa JWT/sesi; komentar di `config.go` sudah menyebut SecretKey tidak wajib) | HAPUS field (nilai di file config dibiarkan) |
| `helper/time/time.go` (11) | StartOfDay, EndOfDay, ToSQLDate, ConvertToDateFormatString, NormalizeDateRange, GenerateMasaPajak, GetIndonesianDayNameAndDate + 4 fungsi kecil | 1 | HAPUS |
| `helper/log/log.go` (3) | LogExternalCall(+WithResponse), IsDebugMode | 1 (tabel ESB/Brigate proyek lain) | HAPUS |
| `helper/request` (3) | ConvertParamStringToInt, ConvertQueryValueToInt, GetTokenParts | 4 untuk ConvertParamStringToInt (id URL dibaca `binder.BindURI` dengan validasi `gt=0`); 1 untuk lainnya (tanpa query string, tanpa Bearer) | HAPUS |
| `helper/request/filter_request_param.go` | **BuildOrderClause** | **2**: `user_repo.GetAll` menulis ulang logika yang sama (daftar kolom diizinkan, arah ASC/DESC) | **PAKAI** di `user_repo` (perilaku sama; bonus: "DESC" huruf besar ikut dikenali) |
| `helper/context_helper.go` | GetUserRole | 1 (role admin diperiksa `ActingUserMiddleware`) | HAPUS |
| `pkg/bcrypt` | HashPasswordWithParams (+uji) | 1 (`HashPassword` dengan parameter tetap dipakai buat user dan ganti password). Catatan: isi paket memakai Argon2id walau bernama bcrypt; tidak diubah | HAPUS |
| `pkg/binder` | BindQuery, BindMultipartForm (+uji) | 1 (hanya JSON dan URI) | HAPUS |
| `validation/validator.go` | 12 aturan khusus (maxfilesize, fileextension, alphanumspace, dateformat, date_range_end/start, required_with, fraction, decimal, alphanumdash, notblank, timeformat) + `AddCustomErrorMessage` (+uji) | 1 (DTO hanya memakai `required`, `min`, `max`, `gt`, `alphanum`, `dive`) | HAPUS (pesan `nefield` dibiarkan seperti Portal) |
| `errors/error_types.go` | ServiceUnavailableError.Error | - | PERTAHANKAN (kelengkapan tipe error, sama seperti Portal) |

### BE-C. Method interface tak terpanggil

| Interface | Method | Saran |
|---|---|---|
| `repository.LogSchedulerRepository` | InsertLogScheduler (+model, baris di `router.go`) | HAPUS (tidak ada penjadwal; tabel `log_schedulers` dibiarkan) |

### BE-D. Endpoint dan handler

| Item | Kemungkinan | Saran |
|---|---|---|
| CRUD Role: `RoleHandler.GetAll/GetByID/Create/Update/Delete/ToggleStatus` tanpa route, + method service/repo yang hanya dipakainya + DTO-nya | 1 (rancangan hanya mencantumkan `/roles/options`; role tetap `admin`/`user` dari data awal; salinan template Portal sebelum keputusan 31) | HAPUS (tetap: `GetActiveOptions` untuk `/roles/options`, `repo.GetByID` dipakai user service saat mengganti role, model `Role`; `repo.GetByName` ikut terhapus karena hanya dipakai `Create`) |
| `GET /health` | - | PERTAHANKAN (SOP dan panduan cutover: `curl 127.0.0.1:8080/api/health`) |
| `POST /users/detail/:id` | **PUTUSKAN**: tercantum di daftar API resmi rancangan, tetapi hanya Portal yang boleh memanggil grup ini (`PortalOnlyMiddleware` + admin) dan Portal sudah tidak meneruskannya (dead code Portal kelompok 1). Jadi praktis tidak bisa dijangkau siapa pun | Saran AI: HAPUS dan perbarui baris daftar API di rancangan |

### Lainnya — OPSIONAL (tidak diubah tanpa izin)

- Field config `AppName`, `AppAuthor`, `AppVersion`, `AppHost` (dibaca dari `.env`) dan `Timezone` tidak dipakai kode mana pun. Bawaan
  template yang juga ada di Portal, Toko, Coffee. Tidak berbahaya. Hal yang sama berlaku di Portal (`App*`, `Timezone`,
  `FormatDate`), terlewat di analisis Portal karena alat tidak memeriksa field config.

**Eksekusi [SELESAI], 16-17 Nov 2026** (user mengikuti semua saran AI: HAPUS/PAKAI disetujui, `/users/detail/:id` dihapus
dan baris daftar API di rancangan diperbarui, field config OPSIONAL dibiarkan):
- BE-A, BE-B, BE-C dihapus seperti daftar; berantai ikut: seluruh CRUD Role (handler 6, service 6, repo 6 termasuk `GetByName`,
  `toRoleResponse`), `UserService.GetByID` + handler, 9 DTO tak terpakai (8 DTO Role + `user/dto.GetByIDRequest`; alat tipe
  mencocokkan nama saja sehingga terlewat, ditemukan dengan pemeriksaan per paket; pemeriksaan yang sama di Portal bersih),
  variabel/konstanta sisa `global.go`/`log.go`, 16 fungsi uji binder/bcrypt + 2 pembantu, library `golang-jwt` dan `go-redis`
  (`go mod tidy`). Field config `SecretKey`, `TokenExpire`, `RefreshTokenExpire`, `CorsAllowOrigins`, `RedisConfig` dihapus dari
  struct (nilai di file config tidak disentuh, kini diabaikan). `validator.go`/`validator_test.go` disamakan dengan Portal.
- Kemungkinan 2: `user_repo.GetAll` memakai `BuildOrderClause`. Perilaku lama dijaga: kolom tak dikenal/kosong diarahkan ke `id`
  sehingga `sort_order` tetap berlaku untuk kolom bawaan (cascade sempat menghapus `BuildOrderClause` karena belum dipakai; file
  dipulihkan dari git, identik dengan HEAD).
- Diff: 41 file, +8/-4.894. gofmt hanya pada file yang rapi di HEAD (6 file sudah tidak rapi sejak awal, dibiarkan).
Uji: build/vet/test 7 paket, `deadcode` hanya `ServiceUnavailableError.Error` (dipertahankan), `ifacecheck` 0, tipe tak terpakai 0,
`staticcheck` hanya `field private` lama; rekaman urutan `/users/list` sebelum/sesudah (`testing/deadcode/beuser_sort_snapshot.js`)
20/22 kombinasi identik, 2 berbeda tepat yang diharapkan (`sort_order: "DESC"` huruf besar kini sama dengan `desc`; injeksi lewat
`sort_by` tetap aman); API BE User `api_3b2` 33/33, API Portal 55/55, API Toko `api_3b4` 31/31, endpoint terhapus menjawab
"Route not found", `/roles/options` dan `/health` 200; browser laptop: `b_portal_k3` 18/18 (ganti role lewat `roleRepo.GetByID`),
`v_portal` 40/40, SSO Toko 32/32, SSO Coffee `b_4_11` 22/22; log BE User tanpa 500/panic.
Perbaikan skrip uji (bukan aplikasi): `api_3b2` S1 menerima path toko dev berupa alamat lengkap (sudah gagal sebelum perubahan);
`b_4_11` mencari produk dengan nama lengkap (21 produk uji "Croissant Kemasan" di DB dev membuat produk baru tidak tampil).
Catatan data dev: `api_3b2` A3 memutar kunci Coffee di BE User; kunci Coffee di `pos_coffee_db.auth_credentials` diperbarui
dari `testing/auth_3b2/keys.txt` agar Coffee dev bisa login lagi (setiap menjalankan `api_3b2`, lakukan hal yang sama).

---

## Toko (dianalisis 17 Nov 2026) — **SELESAI** (kelompok 1 BE, 2 endpoint, 3 FE)

**Keputusan user 17 Nov 2026:** mengikuti semua saran AI (HAPUS/PERTAHANKAN sesuai daftar; Saldo manual dibiarkan;
`formatProductPackage` dipakai; `password-input` dihapus dan `LoginForm` Toko dibiarkan; endpoint dihapus setelah cek log nginx).

**Kelompok 1 (BE tanpa endpoint) [SELESAI], 17 Nov 2026:** BE-A, BE-B, BE-C seperti daftar (kecuali Saldo `Adjust`/`GetBalance` dan
`ClearPinHash` yang ikut fitur PIN di kelompok 2; dijaga lewat `KEEP_IFACE` baru di `cascade.js`). Berantai: query
`todayStatsQuery`/`todayExpensesQuery`/`getAllActiveProductsForLowStockQuery`, `GetLowStockResponse`, variabel sisa template,
16 uji binder/bcrypt + 2 pembantu, library `go-redis`. Validator: 9 aturan + 18 fungsi uji + 18 struct uji dihapus,
`dateformat`/`notblank`/`timeformat` tetap; `TestAddCustomErrorMessage` lama (memakai aturan yang dihapus) diganti uji baru
`TestNotBlankAndTimeformat` (aturan yang dipakai tetapi belum pernah diuji). Komentar penjelasan "stok baru per produk" dipindah
ke `packagesQuery` (sama dengan Coffee) dan judul `// RESPONSE` di DTO produk dikembalikan. Akhir baris CRLF dipertahankan;
gofmt hanya pada 15 file yang rapi di HEAD (file lain banyak yang CRLF/tidak rapi sejak awal, dibiarkan). Diff: 62 file, +30/-4.495.
**Celah alat yang ditemukan:** pemeriksa tipe per paket memakai `\t` di pola grep (tidak dibaca sebagai tab), sehingga tipe di blok
`type ( ... )` tidak pernah diperiksa. Diganti `testing/deadcode/unused_types_pkg.js`; diulang di semua BE: BE User dan Portal tetap 0,
Toko menemukan `GetLowStockResponse` (ikut dihapus), **Coffee masih punya 6 sisa** (5 `GetByIDRequest` milik endpoint detail yang
dihapus di kelompok 5, model `transaction.Transaction`) — belum disentuh, menunggu keputusan user.
Uji: build/vet, `go test` 10 paket + uji DB produk, `deadcode` kosong, `ifacecheck` tepat 3 yang dipertahankan, tipe tak terpakai 0,
`staticcheck` hanya 2 temuan lama; BE Toko dijalankan ulang: uji baru `testing/deadcode/api_toko_k1.js` 15/15 (notblank/timeformat/
dateformat menolak dengan pesan Indonesia di Shift, Pelanggan, Pengeluaran, Retur; aturan bawaan tetap; Ringkasan Bisnis
today/week/month, tren, produk terlaris; endpoint PIN masih ada), `api_3b4` 31/31, `v_portal` 40/40 + `v_toko` 32/32 di laptop/
tablet/HP, `r_toko` 86 layar: 4 layar bermasalah sama persis dengan sebelum perubahan (tabel lebar di dialog Pembelian dan Laporan
Penjualan; lebar isi berubah karena data uji hari ini); log BE tanpa 500/panic. Usang, tidak dijalankan: `portal_toko_e2e/b_flow.js`
(port dan tabel sebelum Fase 3B) dan `test-edit-pembelian.js` (port sebelum restrukturisasi).

**Kelompok 2 (endpoint BE) [SELESAI], 17 Nov 2026.** Cek log nginx produksi oleh user (access.log, .1, dan semua .gz): 0, kecuali
1 baris `GET /api/auth/me` 14 Nov 2026 11:26 UTC (sebelum cutover), jawaban 400, tanpa referer, Firefox Linux: pemindai otomatis
(FE Toko selalu memakai POST; APK memakai WebView Android), bukan pengguna. Dihapus 26 endpoint (daftar BE-D, termasuk
`/reports/sales/chart` yang sempat terlewat skrip karena ditulis `"/sales/chart"` dan ketahuan saat pencocokan ulang BE-FE) beserta
handler, method service, method repo dan query yang hanya dipakainya (`GetTopCategories`, `GetPaymentMethods`, `GetSummaryExtra` +
`GetHighestTransaction`/`GetPeakHour`/`GetAvgTransaction`, `GetSalesChart`, `Reorder`, `ResetAll`, dst.), domain `pin` dan
`stock_mutation` (tabel tetap), 25 DTO, dan paket `pkg/bcrypt` (Argon2; ternyata hanya dipakai PIN, sama seperti Coffee;
`golang.org/x/crypto` kini dependensi tidak langsung). **Repo `UpdateSales`/`UpdateExpenses` tetap** (dipakai transaksi dan
pengeluaran untuk menambah total kas otomatis); yang dihapus hanya service/handler endpoint manualnya. Saldo manual tetap.
Komentar route kas disesuaikan ("open/close"). Route yang tidak dipanggil FE kini tinggal `/health` dan `/settings/:key`;
4 alamat FE tanpa route berasal dari hook FE mati (dihapus di kelompok 3).
Uji: build/vet, `go test` + uji DB produk, `deadcode` kosong, `ifacecheck` hanya Saldo, tipe tak terpakai 0, `staticcheck` hanya
2 temuan lama; uji baru `testing/deadcode/api_toko_k2.js` 50/50 (26 endpoint terhapus "Route not found"; 21 endpoint yang
dipertahankan/bernama mirip tetap 200 termasuk `/menus/options`, `/shifts/active`, `/settings/pagination_sizes`, riwayat Saldo;
pengeluaran tunai menambah total pengeluaran kas otomatis 0 -> 1.234); `api_3b4` 31/31 (T30/T30a dialihkan dari `/auth/me` ke
`/menus/my`, uji satu sesi per user tetap), `api_toko_k1` 15/15 (K1 kini memastikan PIN terhapus), `v_portal` 40/40 + `v_toko`
32/32 di laptop/tablet/HP, `r_toko` 86 layar tanpa masalah baru; `log_requests` Toko: tidak ada "Route not found" dari FE (hanya
dari skrip uji, termasuk T21/T22 `api_3b4` yang sengaja memeriksa endpoint lama 3B.4); log BE tanpa 500/panic.
Diff gabungan kelompok 1+2 (belum di-commit): 154 file, +33/-6.571.

**Kelompok 3 (FE Toko) [SELESAI], 17 Nov 2026.** Dihapus: `ui/dropdown-menu`, `ui/radio-group`, `ui/password-input` (LoginForm Toko
TIDAK diubah), `EmptyState`, `usePermission`, `useLocalStorage`, `string.ts`, `number.ts` (clamp, percentage, lalu `roundTo` yang
hanya dipakai percentage), `capital-reconciliation.constants.ts`, `getApplicablePrice` salinan di products.utils (kasir memakai salinan
identik di cashier.utils), `isPaymentSufficient`, `PAYMENT_STATUS_LABEL`, `PAYMENT_METHOD_VALUES`, `getPlatform` + tipe
`DevicePlatform`, `openReceiptPdf`, `parseRupiah`, `formatRelative`, hook mati (`useGetMeQuery`, `useReorderMenuMutation`,
`useShiftDetailQuery`, `useProductOptionsQuery`, `useProductBarcodeQuery`, `useProductPricesQuery`, `cashier.useProductBarcodeSearchQuery`,
re-export di `backdate-cashier.api.ts`), tipe milik endpoint yang dihapus (`ReorderMenuPayload`, `ProductOption`,
`SummaryExtraResponse`), barrel dirapikan; library `@radix-ui/react-dropdown-menu`, `-radio-group`, `-toast`, alat lint
`@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import` (109 paket turunan, tanpa perubahan versi
library lain). Kemungkinan 2: `formatProductPackage` dipakai di Pembelian, Detail Produk, Form Produk, Kasir Backdate (ekspresi identik,
tipe `ProductPackage` sama). Dipertahankan: hook Saldo manual, `fetchProductPrices`, `webBluetooth.d.ts`, `lint-staged`. CRLF dijaga.
Diff `src`: 38 file, +10/-543. knip: export mati tinggal 3 hook Saldo (sengaja); sisanya OPSIONAL (barrel, kata export, tipe).
Uji: type-check/lint/build; uji baru `testing/deadcode/b_toko_k3.js` 8/8 di laptop/tablet/HP (label "Bungkus (Bungkus (Isi 10 Buah))" di
detail produk, form produk, pilihan satuan pembelian; login Toko dan tombol mata; 7 halaman yang filenya diubah terbuka di alamatnya
sendiri tanpa error, termasuk Kasir Backdate; Simpan PDF struk menghasilkan .pdf; tanpa error JavaScript); `v_portal` + `v_toko` 32/32
di laptop/tablet/HP; `r_toko` 86 layar tanpa masalah baru; tidak ada panggilan "Route not found" baru dari FE; log BE tanpa 500/panic.
Belum diuji langsung: label kemasan di keranjang Kasir Backdate (butuh kas backdate terbuka); dijamin oleh ekspresi identik + type-check.
Catatan skrip uji: pemeriksaan halaman awalnya terlalu longgar (route tak dikenal dialihkan ke "/" tanpa pesan error, dan alamat
Shift salah ditulis); diperbaiki dengan memastikan alamat akhir sama dengan halaman yang dibuka.
**Dead code Toko SELESAI** (kelompok 1-3).

Dianalisis dari kode Toko sendiri. Alat: `deadcode` 100 fungsi, `ifacecheck` 18 method, `staticcheck` 5 temuan, tipe DTO/model tak
terpakai (pemeriksaan per paket) 8, endpoint BE tanpa pemanggil FE 24, FE `knip` 21 file, 3 library + 4 alat dev, 31 export mati.

**Pemeriksaan keamanan produksi (khas Toko):**
- APK hanya cangkang yang membuka web produksi (`mobile/capacitor.config.ts`, `server.url`), FE tanpa service worker: semua
  perangkat selalu memakai FE terbaru, tidak ada FE lama di tablet yang bisa memanggil endpoint lama.
- Tiap endpoint dicek riwayat git FE Toko (termasuk struktur lama `FE/src`) dengan potongan alamat pendek, karena FE bisa menulis
  alamat terpisah (`base + '/summary-extra'`; pencarian alamat lengkap sempat salah menyimpulkan "tidak pernah ada").
- Endpoint yang alamatnya hanya tertulis di hook FE yang mati dihitung TIDAK dipanggil (`/auth/me`, `/menus/reorder`,
  `/shifts/detail/:id`, `/products/options`, saldo manual); alat endpoint menghitungnya "dipakai".
- Folder `domain/_sample_reference` (contoh integrasi ESB BRI) tidak pernah dikompilasi: Go mengabaikan folder berawalan `_`, dan
  paket `pkg/transport` yang diimpornya tidak ada. Alat tidak bisa melihatnya.
- Sebelum endpoint dihapus, disarankan cek log nginx produksi (perintah di bawah) untuk memastikan tidak ada pemanggil di luar FE.

### BE-A. Paket/folder yang tidak diimpor aplikasi — HAPUS

`pkg/redis` (+`RedisConfig` dan komentar di `main.go`), `pkg/encryptor` (+uji), `helper/security`, `helper/regexp` (+uji),
`middleware/auth/basic_auth_middleware.go`, paket model `payment_method`, `payment_status`, `stock_mutation` (domain tetap), folder
`domain/_sample_reference`. Berbeda dari Portal/BE User: `pkg/janitor` dan `helper/codegen` DIPAKAI Toko (izin menu dan kode
produk/kategori/supplier) — PERTAHANKAN.

### BE-B. Fungsi, konstanta, tipe tak terpakai — HAPUS

| Item | Kemungkinan |
|---|---|
| `helper/global.go` 29 fungsi template (faktur, BRI, terbilang, ID, JSON) | 1 |
| `helper/time` EndOfDay, ConvertToDateFormatString, GenerateMasaPajak, GetIndonesianDayNameAndDate + 4 fungsi kecil | 1 |
| `helper/log` LogExternalCall(+WithResponse), IsDebugMode | 1 |
| `helper/request` ConvertParamStringToInt, ConvertQueryValueToInt, GetTokenParts | 4/1 (id URL lewat binder; middleware membaca Bearer sendiri) |
| `helper/context_helper.go` GetDeviceInfo | 1 |
| `pkg/bcrypt` HashPasswordWithParams (+uji), `pkg/binder` BindQuery, BindMultipartForm (+uji), `pkg/jwt` FillResultMapFromClaims | 1 |
| const `updateBackdateSalesQuery` (backdate), `baseUnitNameQuery` (rekonsiliasi stok) | 1 |
| Tipe: `SalesTrendFilter`, `ReceivablePayment`, model `payment_method.PaymentMethod`, `payment_status.PaymentStatus`, `stock_mutation.StockMutation`, `supplier_return.SupplierReturn` (yang dipakai `supplier/model.SupplierReturn`), `transaction/model.Transaction` dan `.StockMutation` (yang dipakai hanya `TransactionItem`) | 1 |
| `validation/validator.go` 9 aturan khusus tak dipakai DTO: maxfilesize, fileextension, alphanumspace, date_range_end, date_range_start, required_with, fraction, decimal, alphanumdash (+uji) | 1. **Toko memakai `dateformat`, `notblank`, `timeformat`: PERTAHANKAN** (beda dari Portal/BE User) |
| `errors` tipe error | PERTAHANKAN |

### BE-C. Method interface tak terpanggil

| Interface | Method | Saran |
|---|---|---|
| 12 repo | `GetDB` | HAPUS (1) |
| `business_summary/repo` | GetTodayStats, GetTodayExpenses (+query) | HAPUS (4: layanan memakai `GetStatsByRange`/`GetExpensesByRange` untuk hari ini) |
| `product/service` | GetLowStock (+repo, DTO, model `LowStockProduct`) | HAPUS (1: tidak pernah punya route) |
| `pin/repo` | ClearPinHash | HAPUS bersama seluruh fitur PIN (BE-D) |
| `customer_balance/repo` | **Adjust**, GetBalance | **PUTUSKAN** (lihat kemungkinan 3 Saldo) |

### BE-D. Endpoint tanpa pemanggil — HAPUS (route + handler + service/repo yang hanya dipakainya)

Tidak pernah dipanggil FE Toko sepanjang riwayat git: `/auth/verify-token`, `/reports/business-summary/{top-categories,
payment-methods}`, `/reports/sales/chart`, `/cash-drawer/update-sales/:id` dan `/update-expenses/:id` (berbahaya: mengubah total kas
manual, padahal total dihitung otomatis), `/customers/active`, detail `/expenses|categories|units|users/detail/:id`, seluruh `/pin/*`
(domain `pin`; tidak ada di data menu; kolom `users.pin_hash` dibiarkan), `/receivables/summary`, `/receivables/payments/:id`,
`/settings/reset` (berbahaya: mengembalikan pengaturan ke bawaan), `/shifts/summary`, `/stock-mutations/list` dan
`/product/:product_id` (seluruh domain `stock_mutation`; tabelnya tetap ditulis domain lain).
Pernah dipanggil lalu **sengaja dihapus user** dari FE: `/reports/business-summary/summary-extra` (commit "standart", 28 Jun 2026).
Hanya dipanggil hook FE yang mati: `/auth/me` (uji `portal_toko_e2e/b_flow.js` B1b ikut disesuaikan), `/menus/reorder`,
`/shifts/detail/:id`, `/products/options`.
**PERTAHANKAN:** `GET /health` (SOP), `GET /settings/:key` (salah tebak alat: dipakai `/settings/pagination_sizes`),
`/products/by-barcode/:code` (dipakai pemindai barcode kasir).

Cek log produksi sebelum eksekusi (copas di VPS; harus: semua 0):
```bash
P='/(toko/)?api/(auth/verify-token|auth/me|reports/business-summary/(top-categories|payment-methods|summary-extra)|reports/sales/chart|cash-drawer/update-(sales|expenses)|customers/active|(expenses|categories|units|users|shifts)/detail|pin/|receivables/(summary|payments)|settings/reset|shifts/summary|stock-mutations|menus/reorder|products/options)'
sudo grep -hcE "\"[A-Z]+ $P" /var/log/nginx/access.log /var/log/nginx/access.log.1
sudo zgrep -hcE "\"[A-Z]+ $P" /var/log/nginx/access.log.*.gz
```

### Kemungkinan 3 — PUTUSKAN (tidak dihapus, tidak dibuatkan fitur)

**Saldo pelanggan manual.** Saldo yang berjalan: bertambah lewat "simpan kembalian ke saldo" di struk, dipakai bayar di kasir,
dikembalikan saat void. Yang tidak punya layar: top up manual, refund manual, riwayat saldo (route `/customers/:id/balance/{topup,
refund,history}` + hook FE dari commit "saldo" 24 Agu 2026 yang tidak pernah dipasang di layar), dan `Adjust`/`GetBalance` di repo.
Pilihan: (a) biarkan (saran AI: fitur produk, butuh keputusan pemilik toko), (b) buatkan layar, (c) hapus.

### FE Toko

| Item | Kemungkinan | Saran |
|---|---|---|
| `ui/dropdown-menu.tsx`, `ui/radio-group.tsx` + library `@radix-ui/react-dropdown-menu`, `-radio-group`, `-toast`; dev `@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import` | 1 | HAPUS |
| `ui/password-input.tsx` | 2 terhadap `LoginForm` (tombol mata ditulis ulang) | **PUTUSKAN**: saran AI hapus file dan **biarkan LoginForm** (tampilan login produksi tidak berubah) |
| `formatProductPackage` | 2: rumus "Satuan (kemasan)" ditulis ulang di 4 tempat (Pembelian, Detail Produk, Form Produk, Kasir Backdate) | **PUTUSKAN**: saran AI pakai (teks hasil identik) |
| `getApplicablePrice` (products.utils) | 4: salinan identik di `cashier.utils` yang dipakai kasir | HAPUS salinan mati |
| `usePermission` | 4: peran tetap (owner/admin) bertentangan dengan izin menu; layar memakai `useAuth` + `RoleGuard` | HAPUS |
| `isPaymentSufficient` | 4: modal bayar sudah menangani kredit dan total 0 | HAPUS |
| `PAYMENT_STATUS_LABEL` | 4: label sama di `StatusBadge` | HAPUS |
| `PAYMENT_METHOD_VALUES` | 4: tipe `PaymentMethod` lebih lengkap (`balance`) | HAPUS |
| `CAPITAL_RECON_MENU_KEY` | 1: router menulis 34 kunci menu sebagai teks | HAPUS |
| `EmptyState` | 4: hampir identik `DataTableEmpty` | HAPUS |
| `getPlatform` | 4: `isAndroid` cukup | HAPUS |
| `openReceiptPdf` | 1: Simpan PDF dipakai, buka di tab baru tidak (sering diblokir WebView APK) | HAPUS |
| `parseRupiah`, `percentage`, `clamp`, `truncate`, `capitalize`, `slugify` (`string.ts`), `formatRelative`, `useLocalStorage` | 1 | HAPUS |
| Hook mati: `useGetMeQuery`, `useReorderMenuMutation`, `useShiftDetailQuery`, `useProductOptionsQuery`, `useProductBarcodeQuery`, `useProductPricesQuery` (`fetchProductPrices` tetap dipakai kasir), `cashier.useProductBarcodeSearchQuery` (pemindai memakai `useBarcodeScan`), re-export di `backdate-cashier.api.ts` | 1 | HAPUS (bersama endpointnya bila ada) |
| Hook saldo manual (`useCustomerBalanceHistoryQuery`, `useTopupBalanceMutation`, `useRefundBalanceMutation`) | 3 | PUTUSKAN bersama Saldo |
| `shared/types/webBluetooth.d.ts`, `lint-staged` | salah tebak alat | PERTAHANKAN |

### PERTAHANKAN / OPSIONAL

- Alat `cmd/`: 5 skrip `backfill_*` (sekali jalan untuk data stok/HPP lama; `backfill_stock_restore` untuk memulihkan backup
  produksi lama) dan `migrate_users_to_auth` (panduan cutover): PERTAHANKAN.
- 17 barrel `index.ts` tak diimpor, kata `export` berlebih, 44 tipe TS tak terpakai, sub-export shadcn: OPSIONAL.
- Field config `AppName`, `AppAuthor`, `AppHost`, `Timezone`: OPSIONAL (sama dengan aplikasi lain).
- `main.go` mengimpor `pkg/database` dua kali (gaya, ST1019): OPSIONAL.

### Usulan eksekusi (setelah disetujui)

1. **Kelompok 1, BE tanpa endpoint** (BE-A, BE-B, BE-C kecuali Saldo): tidak mengubah perilaku API.
2. **Kelompok 2, endpoint BE** (BE-D) setelah user menjalankan cek log nginx produksi di atas dan hasilnya 0.
3. **Kelompok 3, FE** (hapus + kemungkinan 2 sesuai keputusan).
Tiap kelompok: build/vet/test, alat dead code, regresi API dan browser Toko di laptop/tablet/HP (`v_toko`, `portal_toko_e2e`,
uji kasir/backdate/pembelian/laporan yang ada), log BE tanpa 500/panic; deploy Toko oleh user.
