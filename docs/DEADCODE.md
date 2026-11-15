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

### BE-C. Method interface yang dideklarasikan + diimplementasi tetapi tidak pernah dipanggil — DISKUSI (kecuali yang ditandai)

| Interface | Method |
|---|---|
| `business_summary/repo` | GetTodayStats, GetTodayExpenses |
| `pin/repo` | ClearPinHash |
| `product/service` | GetLowStock (beserta `repo.GetLowStock` yang hanya dipakai olehnya) |
| `cash_drawer`, `customer`, `finance`, `product_category`, `product_unit`, `shift`, `stock_mutation`, `stock_reconciliation`, `supplier` (repo) | GetDB (9 interface) **[DIHAPUS]** |

### BE-D. Endpoint API yang tidak dipanggil FE Coffee

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
| `POST /pin/check`, `/pin/set`, `/pin/verify`, `/pin/change` | HAPUS (fitur PIN tidak dipakai di FE Coffee; seluruh domain `pin`) |
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
