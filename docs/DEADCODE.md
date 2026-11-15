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
| `helper/regexp` — DISKUSI | 8 fungsi validasi teks | bandingkan dengan validasi yang ditulis ulang; ikut uji `regexp_test.go` |
| `helper/security` **[DIHAPUS]** | sanitasi SQL, validasi header | repo sudah memakai whitelist kolom urutan |
| `pkg/encryptor` **[DIHAPUS]** | enkripsi AES | ikut file uji |
| `pkg/redis` **[DIHAPUS]** | pengelola Redis | `RedisConfig` dan komentar Redis di `main.go` ikut dibuang |

Sisa untuk diskusi: `helper/regexp`.

### BE-B. Fungsi tak terpanggil di paket yang masih dipakai — DISKUSI (kecuali yang ditandai)

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
| `shared/components/ui/password-input.tsx` | DISKUSI, kemungkinan 2: `LoginForm` menulis ulang tampil/sembunyi password sendiri |
| `shared/components/ui/radio-group.tsx` | DISKUSI, kemungkinan 2: 4 tempat membuat tombol pilihan sendiri (`role="radio"`) |
| `shared/types/webBluetooth.d.ts` | PERTAHANKAN (tipe `navigator.bluetooth` untuk printer Bluetooth; alat salah tebak) |
| 15 barrel `index.ts` tak diimpor (customers/receivables, dashboard, finance/overview, menu-items, procurement/purchases, reporting/{business-summary, cashier-performance, profit-loss, sales, stock}, sales/transactions, settings, settings/{backup, menus, users}) | OPSIONAL (hanya pengumpul export, tidak ikut ke hasil build) |

### FE-F. Library npm tidak terpakai

| Library | Saran |
|---|---|
| `@radix-ui/react-dropdown-menu`, `@radix-ui/react-toast` | **[DIHAPUS]** |
| `@radix-ui/react-radio-group` | ikut keputusan `radio-group.tsx` |
| devDependencies `@typescript-eslint/eslint-plugin`, `@typescript-eslint/parser`, `eslint-plugin-import` | **[DIHAPUS]** (lint memakai paket `typescript-eslint`) |
| `lint-staged` | PERTAHANKAN (dipakai hook Husky `.husky/pre-commit`) |

### FE-G. Fungsi/hook diekspor yang tidak dipakai di mana pun — DISKUSI

`auth.api.ts` useGetMeQuery; `menu.api.ts` useReorderMenuMutation; `shifts.api.ts` useShiftDetailQuery;
`purchases.types.ts` PAYMENT_STATUS_LABEL; `products.api.ts` useProductBarcodeQuery, useProductPricesQuery (+ `fetchProductPrices`
yang hanya dipakai olehnya); `products.utils.ts` getApplicablePrice, formatProductPackage; `blePrinter.ts` getPlatform;
`cashier.utils.ts` isPaymentSufficient; `receiptPdf.ts` openReceiptPdf; komponen `EmptyState`; hook `useLocalStorage`, `usePermission`;
`currency.ts` parseRupiah; `date.ts` formatRelative; `number.ts` clamp, percentage; `string.ts` truncate, capitalize, slugify.

### FE-H. Hanya kata `export` yang berlebih (dipakai di file sendiri) — OPSIONAL

formatTime (shifts.utils), fetchProductPackages, printViaBle, receiptToPlainText, isAndroid, calcGrandTotal, OPEN_BILLS_KEY,
buildReceiptPdf, buildReceiptText, authClient, useMediaQuery, roundTo.

### FE-I. Lainnya — PERTAHANKAN / OPSIONAL

- 43 tipe TypeScript diekspor tetapi tidak dipakai: OPSIONAL (tidak ikut ke hasil build).
- Sub-export komponen shadcn (`DialogTrigger`, `SelectGroup`, `CardFooter`, dll.): PERTAHANKAN (pustaka komponen standar).
