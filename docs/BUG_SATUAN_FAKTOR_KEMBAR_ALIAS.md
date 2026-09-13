# Bug: Satuan Faktor Kembar (Alias) Ditolak — "struktur satuan bercabang"

Status: **Kode SUDAH diperbaiki & diverifikasi lewat unit test.**
Tanggal: 13 September 2026
Dilaporkan pada: Produk "Telur" (id 83, TEL-0001), menu Rekonsiliasi Stok.

---

## 1. Gejala

Saat membuka produk Telur di menu **Rekonsiliasi Stok**, muncul error:

> "Produk ID 83 punya struktur satuan bercabang, tidak bisa diproses otomatis
> -- hubungi admin"

Produk tidak bisa direkonsiliasi. Operasi stok lain untuk produk ini juga
berpotensi ikut terblokir (jalur tulis stok memakai guard yang sama).

## 2. Penyebab

Struktur `product_packages` Telur (id 83):

| id  | satuan   | definisi                | faktor ke anchor (Krak) |
|-----|----------|-------------------------|--------------------------|
| 139 | Krak     | anchor (is_default)     | 1                        |
| 140 | Kilogram | 1 Kilogram = 1 Krak     | **1** (kembar dgn Krak)  |
| 160 | Gram     | 2 Gram = 1 Kilogram     | 0,5                      |
| 161 | Gram     | 4 Gram = 1 Kilogram     | 0,25                     |
| 141 | Pieces   | 4 Pieces = 1 Gram(161)  | 0,0625                   |

Krak dan Kilogram punya **faktor konversi sama persis (= 1)** karena toko
memang menganggap keduanya identik ("1 krak = 1 kilogram"). Ini WAJAR secara
bisnis — dua nama untuk satuan yang sama (alias).

Namun `analyzePackageChain` di `BE/domain/product/model/stock_delta.go` punya
guard yang menolak (`ErrBranchingChain`) begitu ada dua paket berfaktor sama,
dengan alasan "pemilihan satuan terkecil jadi ambigu". Guard ini terlalu ketat:
kalau faktor dua satuan sama, hasil hitung stok (di satuan anchor) SELALU sama
apa pun yang dipilih sebagai "terkecil" — jadi tidak ada ambiguitas yang
berbahaya.

## 3. Perbaikan (kode, bukan data)

File: `BE/domain/product/model/stock_delta.go`.

- **Guard faktor-kembar dihapus.** Dua satuan berfaktor sama kini diperlakukan
  sebagai alias dan diproses normal (tidak lagi melempar `ErrBranchingChain`).
- **Tie-breaker deterministik ditambahkan** di dua tempat agar distribusi stok
  per-baris tidak bergantung urutan data / sort yang tidak stabil:
  - pemilihan `smallest` (satuan terkecil): kalau faktor sama, pilih paket
    dengan **ID terkecil**.
  - `sort.Slice` breakdown di `ComputeStockDelta`: kalau faktor sama, urut
    stabil berdasarkan **ID** paket.

Kenapa aman: total stok di satuan anchor selalu benar (penjumlahan komutatif).
Tie-breaker hanya menjamin baris tujuan stok konsisten antar-operasi (mis. jual
lalu void tidak memindah stok antar-alias).

`ErrBranchingChain` TIDAK dihapus dari kode (masih dipakai wrapper error di
`stock_delta_repo.go` & `purchase_repo.go`), hanya tidak pernah lagi dilempar
untuk kasus faktor kembar.

## 4. Verifikasi

- `go build ./...` sukses (EXIT=0).
- Unit test `domain/product/model` semua PASS, termasuk:
  - `TestComputeStockDelta_EggRealStructure_NotRejected` — struktur Telur nyata
    (Krak=Kilogram) kini diproses, total anchor benar.
  - `TestComputeStockSummary_EggRealStructure_NotRejected` — jalur baca juga OK.
  - `TestComputeStockDelta_EqualFactorAlias_Allowed_Deterministic` — faktor
    kembar diproses dengan distribusi deterministik (tie-break ID).
  - Semua test lama tetap hijau (void round-trip, cross-level borrow, reserved
    qty, non-terminating ratio, star topology, needs_stock_review, dll).

## 5. Catatan penting untuk admin/toko (belum ditindaklanjuti)

Perbaikan ini membuat sistem MENERIMA struktur satuan Telur apa adanya, tapi
struktur itu sendiri masih janggal secara fisik dan sebaiknya dirapikan lewat
UI master produk agar konversi sesuai kenyataan:

- Dua satuan sama-sama bernama "Gram" (id 160 & 161) padahal artinya "Setengah
  Kilo" dan "Seperempat" — membingungkan.
- "Pieces (butir telur)" menunjuk ke "Gram Seperempat", bukan ke Kilogram/Krak.
- Apakah 1 Krak benar-benar = 1 Kilogram secara berat? (umumnya 1 krak telur
  ~1,5-1,7 kg / ~30 butir). Kalau tidak, konversi ke Pieces/Gram akan meleset.

Ini keputusan data master (bukan bug kode) dan perlu konfirmasi pemilik toko.
