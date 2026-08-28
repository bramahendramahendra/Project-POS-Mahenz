# Desain Menu: Rekonsiliasi Stok (Stok Lama vs Stok Baru)

> Status: **DRAFT DESAIN — belum implementasi.** Dokumen ini untuk didiskusikan sampai final dulu.
> Tujuan: menu untuk membandingkan stok lama (sebelum migrasi) vs stok baru (hasil rekonstruksi), menampilkan **rincian dari mana stok baru dihitung**, lalu memungkinkan **input manual per level satuan** untuk menetapkan stok yang benar.

---

## 1. Latar Belakang

Setelah migrasi stok (lihat `docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md`), ada **32 produk** yang ditandai `needs_stock_review` karena hasil rekonstruksi tidak cocok dengan stok lama, atau replay riwayatnya gagal. Stok produk-produk ini sengaja **tidak diisi otomatis**.

Selama ini untuk mengecek kebenaran stok, harus buka data satu per satu di banyak menu (pembelian, transaksi, retur). Menu ini menyatukan semua itu dalam satu layar: **stok lama, stok baru, rincian perhitungannya, dan form koreksi manual**.

Menu ini juga berguna terus-menerus (bukan sekali pakai): kapan pun ada produk yang stoknya dicurigai salah, admin bisa buka menu ini untuk audit + koreksi.

---

## 2. Konsep Data (Bagaimana "Stok Baru" Dihitung)

### 2.1 Sumber kebenaran: tabel `stock_mutations` (ledger / kartu stok)

Setiap pergerakan stok tercatat sebagai satu baris di `stock_mutations`. Stok baru **bukan angka manual** — dia hasil menjumlahkan seluruh mutasi ini. Kolom penting:

| Kolom | Arti |
|-------|------|
| `product_id` | Produk terkait |
| `mutation_type` | Jenis pergerakan (lihat tabel di bawah) |
| `quantity` | Jumlah (dalam **satuan dasar / anchor**) |
| `stock_before` / `stock_after` | Stok sebelum & sesudah mutasi (satuan dasar) |
| `reference_type` / `reference_id` | Menunjuk balik ke nota sumbernya (pembelian/transaksi/retur/expiry) |
| `package_id` | Level satuan (diisi untuk mutasi baru, **NULL untuk data lama**) |
| `notes`, `user_id`, `created_at` | Catatan, pelaku, waktu |

### 2.2 Peta jenis mutasi → efek → sumber

Ini adalah "dari mana stok baru dihitung", diambil persis dari logika skrip `backfill_stock_restore` (`resolveEventsForMutation`):

| `mutation_type` | Arti | Efek stok | Tabel sumber |
|-----------------|------|-----------|--------------|
| `in` | Pembelian diterima | **+ tambah** | `purchase_items` |
| `void_purchase` | Pembelian dibatalkan | **− kurang** | `purchase_items` |
| `out` | Penjualan kasir | **− kurang** | `transaction_items` |
| `void` | Penjualan dibatalkan | **+ tambah** | `transaction_items` |
| `return` | Retur ke supplier | **− kurang** | `supplier_return_items` |
| `expired` | Kadaluarsa / rusak | **− kurang** | `product_expiry_batches` |
| `adjustment` | Koreksi manual | (lihat catatan) | — (tidak ada tabel sumber) |

**Rumus stok baru (satuan dasar):**
```
stok_baru = Σ(in) − Σ(void_purchase) − Σ(out) + Σ(void) − Σ(return) − Σ(expired) ± Σ(adjustment)
```

### 2.3 Batasan yang HARUS diakui di UI (biar tidak menyesatkan)

1. **`quantity` dalam satuan dasar (anchor).** Semua angka breakdown ditampilkan dalam satuan dasar agar konsisten. Kalau ingin per level satuan, hanya bisa untuk mutasi yang `package_id`-nya terisi.
2. **Data lama `package_id` = NULL.** Breakdown **per level satuan** tidak tersedia untuk mutasi lama. Breakdown **per jenis** (beli/jual/retur/dst) tetap akurat.
3. **`adjustment` tidak bisa ditelusuri ke sumber.** Kalau ada mutasi `adjustment`, jumlahnya tetap ditampilkan (dari ledger), tapi tanpa link nota. Produk yang gagal migrasi kemarin bukan karena adjustment, jadi ini jarang.
4. **Stok baru resmi** yang dipakai aplikasi = jumlah `product_packages.stock` dikonversi ke anchor (via `ResolvePackageFactor`). Nilai dari ledger dipakai untuk **menjelaskan** angka itu, bukan menggantikannya.

---

## 3. Skema Data Terkait

### 3.1 Stok lama — `products_stock_backup`
Tabel cadangan (dibuat manual saat migrasi, bukan migrasi resmi):
```
id INT, stock DECIMAL(15,3), reserved_qty DECIMAL(15,3)
```
> UI harus menangani kasus **tabel ini tidak ada** (di environment yang tidak lewat proses migrasi) → tampilkan "stok lama: tidak tersedia".

### 3.2 Stok baru per level — `product_packages`
```
id, product_id, unit_id, unit_name, package_name, ref_package_id, qty, ref_qty,
stock DECIMAL(15,3), reserved_qty, is_default (anchor), is_active, is_continuous
```
- Tiap produk punya beberapa baris paket membentuk **rantai linear** via `ref_package_id` sampai ke **anchor** (`is_default=1`, `ref_package_id=NULL`).
- Konversi antar level pakai `ResolvePackageFactor(packages, packageID)`.
- **Stok fisik disimpan di kolom `stock` tiap baris paket.** Input manual per level = update `stock` di baris-baris paket ini.

### 3.3 Flag review — `products`
```
needs_stock_review TINYINT(1) DEFAULT 0
stock_review_note  TEXT NULL
```

---

## 4. Rancangan UI/UX

Menu baru sub dari **Pelaporan** (atau grup Operasional — lihat bagian Pertanyaan). Mengikuti pola halaman `reporting/stock` yang sudah ada (PageHeader + FilterBar + DataTable).

### 4.1 Layar 1 — Daftar Rekonsiliasi

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Rekonsiliasi Stok                                                        │
│  Pelaporan / Rekonsiliasi Stok                                            │
├─────────────────────────────────────────────────────────────────────────┤
│  [ Cari produk...]  [Kategori ▾]  [Filter: ◉ Perlu Ditinjau ○ Semua ]     │
├─────────────────────────────────────────────────────────────────────────┤
│  Kode   Produk              Stok Lama   Stok Baru   Selisih   Status  Aksi│
│  P-001  Toppas Merah 12       4.917       0.000     −4.917    ⚠ Tinjau [›]│
│  P-053  Kerupuk              49.000      43.000     −6.000    ⚠ Tinjau [›]│
│  P-172  Kopi ABC Botol       17.910       0.667    −17.243    ⚠ Tinjau [›]│
│  P-010  Indomie Goreng       24.000      24.000      0.000    ✓ Cocok  [›]│
│  ...                                                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                            « ‹  1 2 3 … ›  »               │
└─────────────────────────────────────────────────────────────────────────┘
```

Kolom:
- **Stok Lama**: `products_stock_backup.stock` (satuan dasar). "—" bila tabel tidak ada.
- **Stok Baru**: Σ `product_packages.stock` → anchor.
- **Selisih**: `stok_baru − stok_lama` (warna merah bila ≠ 0).
- **Status**: `⚠ Perlu Ditinjau` (needs_stock_review=1) / `✓ Cocok` / `● Sudah Dikoreksi`.
- **Aksi [›]**: buka panel detail (Layar 2).

Default filter: **Perlu Ditinjau** (32 produk), bisa diubah ke Semua.

### 4.2 Layar 2 — Detail Rekonsiliasi Produk (drawer/panel samping atau halaman detail)

```
┌───────────────────────────────────────────────────────────────────────┐
│  ← Kerupuk (P-053)                                    [Tandai Selesai]  │
├───────────────────────────────────────────────────────────────────────┤
│  Stok Lama : 49.000 Bungkus     Stok Baru : 43.000 Bungkus              │
│  Selisih   : −6.000             Catatan tinjau: "selisih di luar wajar" │
├───────────────────────────────────────────────────────────────────────┤
│  RINCIAN PERHITUNGAN STOK BARU (satuan dasar: Bungkus)                   │
│                                                                         │
│    + Pembelian (in)            :   +120.000                             │
│    − Pembelian dibatalkan      :     −0.000                             │
│    − Penjualan (out)           :    −79.000                             │
│    + Penjualan dibatalkan      :     +2.000                             │
│    − Retur supplier            :     −0.000                             │
│    − Kadaluarsa/rusak          :     −0.000                             │
│    ± Koreksi manual            :     −0.000                             │
│    ─────────────────────────────────────────                           │
│    = Stok Baru                 :     43.000                             │
├───────────────────────────────────────────────────────────────────────┤
│  KARTU STOK (riwayat mutasi, terbaru dulu)          [semua ▾][filter]   │
│  Tgl         Jenis      Qty    Sebelum  Sesudah  Ref            Oleh    │
│  2026-08-10  Penjualan  −2.0   45.0     43.0     TRX-0912       kasir1  │
│  2026-08-09  Pembelian +20.0   25.0     45.0     PB-0231        admin   │
│  ...                                                                    │
├───────────────────────────────────────────────────────────────────────┤
│  KOREKSI MANUAL (input per level satuan)                                │
│                                                                         │
│    Level          Satuan     Stok Sekarang    Stok Benar                │
│    Anchor (dasar) Bungkus         43.000       [  49.000 ]              │
│    Karton         Karton (×20)     2.000       [   2.000 ]  (= 40 dasar)│
│                                                                         │
│    Total setara satuan dasar (baru): 49.000                             │
│    Catatan koreksi: [ ................................................ ]│
│                                                                         │
│                                   [ Batal ]   [ Simpan Koreksi Stok ]   │
└───────────────────────────────────────────────────────────────────────┘
```

Perilaku form koreksi manual:
- Menampilkan **semua level satuan aktif** produk (`product_packages` `is_active=1`).
- User isi "Stok Benar" **per level**. Di sampingnya ditampilkan setara satuan dasarnya (real-time, pakai faktor konversi) supaya user tahu totalnya.
- Validasi: tiap nilai `≥ 0` (sesuai CHECK constraint DB), desimal diizinkan.
- Tombol **Simpan Koreksi Stok**:
  - Update `product_packages.stock` untuk tiap level.
  - Tulis baris `stock_mutations` bertipe `adjustment` per level yang berubah (audit trail: siapa, kapan, dari berapa ke berapa, catatan). **Direkomendasikan** agar koreksi tetap tercatat di kartu stok.
  - Set `needs_stock_review = 0` dan simpan `stock_review_note` = catatan koreksi + timestamp.
- Tombol **Tandai Selesai** (tanpa ubah angka): untuk produk yang setelah ditinjau ternyata stok barunya sudah benar → cukup clear flag `needs_stock_review`.

---

## 5. Rancangan Backend (API)

Domain baru: `stock_reconciliation` (mengikuti struktur domain lain: `dto/ model/ repo/ service/ handler/`). Registrasi via segment route baru + 1 baris di `protected_routes.go`. Permission pakai menu key **`pelaporan.stok`** (atau menu baru khusus — lihat Pertanyaan).

### 5.1 `POST /stock-reconciliation/list`
Daftar produk + stok lama vs baru + selisih + status. Filter: `search`, `category_id`, `only_review` (bool), pagination.

Response item:
```json
{
  "product_id": 53, "product_code": "P-053", "product_name": "Kerupuk",
  "base_unit": "Bungkus",
  "old_stock": 49.000, "old_stock_available": true,
  "new_stock": 43.000,
  "diff": -6.000,
  "needs_stock_review": true,
  "stock_review_note": "selisih rekonstruksi vs stok lama di luar wajar..."
}
```

### 5.2 `GET /stock-reconciliation/detail/:product_id`
Detail 1 produk: header (lama/baru/selisih/note), breakdown per jenis mutasi, daftar level satuan + stok sekarang, dan kartu stok (bisa reuse `stock_mutations` GetByProduct yang sudah ada).

Response:
```json
{
  "product_id": 53, "product_name": "Kerupuk", "base_unit": "Bungkus",
  "old_stock": 49.000, "old_stock_available": true,
  "new_stock": 43.000, "diff": -6.000,
  "needs_stock_review": true, "stock_review_note": "...",
  "breakdown": {
    "purchase_in": 120.000, "purchase_void": 0, "sale_out": 79.000,
    "sale_void": 2.000, "supplier_return": 0, "expired": 0, "adjustment": 0,
    "computed_new_stock": 43.000
  },
  "packages": [
    { "package_id": 10, "unit_name": "Bungkus", "is_default": true,
      "factor_to_base": 1, "current_stock": 43.000 },
    { "package_id": 11, "unit_name": "Karton", "is_default": false,
      "factor_to_base": 20, "current_stock": 2.000 }
  ]
}
```
> Breakdown dihitung dengan `SUM(quantity)` di `stock_mutations` `GROUP BY mutation_type WHERE product_id=?`. Angka dalam satuan dasar.

### 5.3 `POST /stock-reconciliation/adjust` (perlu `can_edit`)
Koreksi manual per level. Payload:
```json
{
  "product_id": 53,
  "levels": [
    { "package_id": 10, "new_stock": 49.000 },
    { "package_id": 11, "new_stock": 2.000 }
  ],
  "note": "Hasil hitung fisik gudang 2026-08-28"
}
```
Proses (dalam 1 transaksi DB):
1. Validasi produk ada, tiap `package_id` milik produk itu & aktif, `new_stock ≥ 0`.
2. Untuk tiap level yang berubah: catat `stock_mutations` tipe `adjustment` (`stock_before`, `stock_after`, `package_id`, `notes`, `user_id`), lalu update `product_packages.stock`.
3. Set `products.needs_stock_review = 0`, `stock_review_note = note + " (dikoreksi manual <user> <tgl>)"`.
4. Kembalikan detail terbaru.

### 5.4 `POST /stock-reconciliation/mark-reviewed` (perlu `can_edit`)
Clear flag tanpa ubah stok. Payload `{ "product_id": 53, "note": "sudah dicek, stok baru benar" }`.

> Catatan: `adjustment` sebaiknya mulai didukung juga di logika replay ke depan supaya konsisten. Untuk fitur ini, penulisan `adjustment` cukup untuk audit trail; tidak wajib mengubah skrip backfill (skrip itu one-off).

---

## 6. Rancangan Frontend

Ikuti struktur `FE/src/features/reporting/stock/`:
```
FE/src/features/reporting/stock-reconciliation/
  StockReconciliationPage.tsx
  stock-reconciliation.api.ts       (useReconListQuery, useReconDetailQuery, useAdjustMutation, useMarkReviewedMutation)
  stock-reconciliation.types.ts
  components/
    ReconFilterBar.tsx              (search, kategori, toggle "Perlu Ditinjau")
    ReconTableColumns.tsx
    ReconDetailDrawer.tsx           (breakdown + kartu stok + form koreksi)
    ReconBreakdownCard.tsx
    ReconLevelInputTable.tsx        (input per level + total setara dasar real-time)
```
Registrasi:
- Route baru di `FE/src/shared/constants/routes.ts` + `FE/src/app/router.tsx` dengan `menuKey`.
- Reuse komponen `DataTable`, `PageHeader`, `formatStockNumber` yang sudah ada.

---

## 7. Rencana Perubahan (Ringkasan File)

**Backend (baru):**
- `BE/domain/stock_reconciliation/{dto,model,repo,service,handler}/...`
- `BE/routes/segment/stock_reconciliation_routes.go`
- 1 baris di `BE/routes/protected_routes.go`
- (opsional) migrasi menu baru di seed bila pakai menu terpisah

**Frontend (baru):**
- `FE/src/features/reporting/stock-reconciliation/...`
- Tambahan route + konstanta

**Tidak diubah:** skrip backfill (`backfill_stock_restore`, `backfill_stock_per_level`), logika stok live, migrasi yang sudah jalan.

---

## 8. Pertanyaan untuk Difinalkan (mohon keputusan)

1. **Penempatan menu**: taruh sebagai menu terpisah **"Rekonsiliasi Stok"** di grup Operasional/Pelaporan (butuh seed menu + permission baru), atau **jadikan tab di dalam halaman Laporan Stok** yang sudah ada (pakai permission `pelaporan.stok`, tanpa seed baru)? Rekomendasi saya: **menu terpisah** karena ini fitur audit+edit, beda sifat dari laporan read-only.

2. **Catat koreksi manual sebagai mutasi `adjustment`?** Rekomendasi saya: **ya**, supaya setiap koreksi masuk kartu stok (audit lengkap, bisa ditelusuri siapa yang mengubah). Alternatif: update `stock` diam-diam tanpa jejak (tidak disarankan).

3. **Stok lama saat tabel `products_stock_backup` tidak ada**: tampilkan "—/tidak tersedia" dan sembunyikan kolom selisih? (Rekomendasi: ya, biar tetap jalan di environment tanpa migrasi.)

4. **Hak akses**: siapa yang boleh mengoreksi stok di menu ini — hanya Owner/Admin, atau juga role lain? (memengaruhi seed permission)

5. **Reserved qty**: apakah form koreksi cukup mengatur `stock` saja, atau juga perlu bisa mengatur `reserved_qty` per level? (Rekomendasi: cukup `stock` dulu, reserved diurus alur normal.)

---

## 9. Lampiran — Contoh Query Breakdown (referensi implementasi)

```sql
-- Breakdown per jenis mutasi (satuan dasar) untuk 1 produk
SELECT mutation_type, SUM(quantity) AS total
FROM stock_mutations
WHERE product_id = ?
GROUP BY mutation_type;

-- Stok lama
SELECT stock FROM products_stock_backup WHERE id = ?;

-- Level satuan + stok sekarang
SELECT pp.id, u.name AS unit_name, pp.is_default, pp.ref_package_id,
       pp.qty, pp.ref_qty, pp.stock
FROM product_packages pp
JOIN units u ON u.id = pp.unit_id
WHERE pp.product_id = ? AND pp.is_active = 1;
```
