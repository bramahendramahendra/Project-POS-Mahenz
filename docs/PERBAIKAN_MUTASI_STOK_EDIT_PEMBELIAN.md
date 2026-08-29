# Rencana Perbaikan: Mutasi Stok Salah Tipe saat Edit Pembelian

> Status: **CATATAN / RENCANA — belum diimplementasi.** Dokumen ini untuk disepakati dulu.
> Ditemukan saat investigasi kenapa banyak produk `needs_stock_review` setelah migrasi stok.

---

## 1. Ringkasan Masalah

Saat sebuah Purchase Order (PO) **diedit** untuk **menambah produk baru**, stok produk itu bertambah dengan benar, TAPI dicatat di `stock_mutations` dengan tipe **`adjustment`**, bukan **`in`** (pembelian).

Akibatnya:
1. **Kartu stok jadi salah makna** — pembelian tampak seperti "koreksi manual", bukan pembelian dari supplier.
2. **Rekonstruksi stok (backfill migrasi) gagal** untuk produk itu — karena backfill tidak bisa memproses tipe `adjustment`, produk jadi minus saat replay dan ditandai `needs_stock_review`.

### Bukti (hasil test terkontrol)
Test: buat PO status `unpaid` dengan 1 produk (id 2, qty 5) → edit, tambah produk kedua (id 3, qty 7). Hasil di `stock_mutations`:

| Produk | Cara masuk | Tipe tercatat | Seharusnya |
|--------|-----------|---------------|------------|
| 2 (saat create) | Pembelian | `in` ✅ | `in` |
| 3 (ditambah saat edit) | Pembelian | **`adjustment`** ❌ | **`in`** |

Catatan mutasi produk 3: *"Edit PO ID ... — penyesuaian stok (selisih item baru vs lama)"*.

### Kenapa berkorelasi dengan status "hutang"
Dari data prod: pembelian **lunas** 0% bermasalah, tapi **hutang/bayar-sebagian** banyak yang bolong mutasi `in`-nya. Penyebabnya bukan status pembayaran langsung, melainkan **kebiasaan**: PO hutang/partial lebih sering diedit belakangan (menambah barang), dan setiap produk yang ditambah lewat edit tercatat `adjustment` → tidak terhitung saat migrasi.

---

## 2. Apa Itu `adjustment` (dan Kenapa Salah di Sini)

`adjustment` = penyesuaian stok yang **bukan** dari transaksi bisnis normal (bukan beli/jual/retur/expired). Dipakai untuk koreksi administratif/manual.

Pemakaian `adjustment` saat ini di kode:

| Lokasi | `reference_type` | Tujuan | Tepat? |
|--------|-----------------|--------|--------|
| Buat produk dgn stok awal (`product_repo.go`) | `product_create` | Set stok pembukaan | ✅ |
| Edit produk ubah stok (`product_repo.go`) | `product_edit` | Koreksi stok manual dari form produk | ✅ |
| **Edit PO tambah/ubah item (`purchase_repo.go`)** | `purchase` | Penyesuaian stok karena PO diedit | ⚠️ **Salah untuk item baru** |
| Koreksi manual menu Rekonsiliasi Stok | `stock_reconciliation` | Koreksi hasil tinjauan | ✅ |

Kesimpulan: `adjustment` **sah untuk koreksi manual**, tapi **salah dipakai untuk pembelian**. Produk yang benar-benar baru ditambahkan ke PO adalah **pembelian nyata** (barang masuk dari supplier) → harus `in`.

---

## 3. Lokasi Kode Sumber Masalah

File: `BE/domain/supplier_purchase/repo/purchase_repo.go`, fungsi `Update()` — blok "Terapkan hanya SELISIH bersih per (produk, paket)" (sekitar baris 517-555).

Logika sekarang: menghitung `delta = newQty - oldQty` per (produk, paket), lalu menerapkan delta itu dengan `MutationType: "adjustment"` untuk **SEMUA** perubahan — termasuk produk yang sebelumnya tidak ada di PO (item benar-benar baru).

```go
// (ringkas) — versi SEKARANG
for k := range keys {
    delta := newQty[k] - oldQty[k]
    if delta == 0 { continue }
    direction := StockIn; qty := delta
    if delta < 0 { direction = StockOut; qty = -delta }
    ApplyStockDelta(... MutationType: "adjustment", ReferenceType: "purchase" ...)
}
```

Masalah: tidak membedakan **item baru** (belum pernah ada di PO → seharusnya `in`) dari **perubahan qty item lama** (koreksi terhadap pembelian yang sudah tercatat).

---

## 4. Rencana Perbaikan

### Lapisan 1 — Perbaiki `Update()` (root cause, untuk data BARU) — DISETUJUI

Bedakan jenis perubahan saat edit PO, dan pilih tipe mutasi yang benar:

| Kondisi | Arti bisnis | Tipe mutasi yang benar |
|---------|-------------|------------------------|
| Produk **baru** ditambahkan (tidak ada di item lama) | Pembelian tambahan | **`in`** |
| Qty produk lama **dinaikkan** | Pembelian bertambah | **`in`** (sebesar selisih) |
| Qty produk lama **diturunkan** | Pembelian dikurangi | **`void_purchase`** (sebesar selisih) |
| Produk lama **dihapus** dari PO | Pembelian dibatalkan | **`void_purchase`** (sebesar qty lama) |

> Rasional: semua perubahan di dalam PO pada hakikatnya adalah **pembelian** (atau pembatalan sebagian pembelian). Jadi tipe yang tepat adalah `in` (nambah) / `void_purchase` (kurang), BUKAN `adjustment`. Dengan begitu backfill & kartu stok konsisten memperlakukannya sebagai pembelian.

> Catatan desain: pendekatan "delta bersih" yang sudah ada tetap dipertahankan (bagus, mencegah kegagalan kalau stok riil sudah berubah). Yang berubah hanya **pemilihan `MutationType`**: `in` untuk delta positif, `void_purchase` untuk delta negatif — menggantikan `adjustment`.

Alternatif yang DITOLAK: memakai `adjustment` untuk perubahan qty item lama. Ditolak karena tetap membuat backfill/analisis tidak bisa menelusuri sebagai pembelian.

### Lapisan 2 — Sesuaikan data LAMA yang terlanjur `adjustment` — PERLU KEPUTUSAN

Data 32 produk yang sudah terlanjur punya mutasi `adjustment` (dari edit-PO lama). Dua opsi:

**Opsi 2a — Perbaiki skrip backfill lalu jalankan ulang (otomatis, disarankan)**
Ajari `backfill_stock_restore` supaya mutasi `adjustment` dengan `reference_type='purchase'` diperlakukan sebagai pembelian:
- `adjustment` + `reference_type='purchase'` + arah menambah → perlakukan seperti `in`
- `adjustment` + `reference_type='purchase'` + arah mengurangi → perlakukan seperti `void_purchase`
- Arah (+/−) ditentukan dari `stock_after` vs `stock_before` pada baris mutasi itu (JANGAN diasumsikan).
- Data bisa direkonstruksi karena `purchase_items` (qty + package_id) masih lengkap dan `reference_id` menunjuk ke PO.

Kelebihan: mayoritas 32 produk beres otomatis tanpa input manual. Kekurangan: **mengubah skrip backfill** (perlu izin — user sebelumnya menekankan jangan ubah skrip backfill yang sudah ada).

**Opsi 2b — Koreksi manual lewat menu Rekonsiliasi Stok**
Input satu-satu pakai angka di `docs/ANALISIS_REKONSILIASI_STOK_32_PRODUK.md`. Kelebihan: tidak menyentuh skrip. Kekurangan: manual, lebih lama.

> CATATAN untuk Opsi 2a: arah `adjustment` dari edit-PO bisa positif (nambah barang) ATAU negatif (kurangi/hapus). Skrip harus menghormati arah dari `stock_before`/`stock_after`, bukan menyamaratakan jadi `in`. Ini syarat wajib supaya tidak salah hitung.

---

## 5. Dampak & Cakupan Perubahan

**Lapisan 1 (Update):**
- File: `BE/domain/supplier_purchase/repo/purchase_repo.go` — hanya blok pemilihan `MutationType` di `Update()`.
- Tidak mengubah alur `Create()` (sudah benar: pakai `in`).
- Tidak mengubah `Void()` (sudah benar: pakai `void_purchase`).
- Perlu test ulang skenario: tambah produk, naikkan qty, turunkan qty, hapus produk saat edit.

**Lapisan 2 (kalau Opsi 2a dipilih):**
- File: `BE/cmd/backfill_stock_restore/main.go` — tambah handling `adjustment` bersumber `purchase`.
- Setelah itu jalankan ulang backfill (butuh reset/kondisi yang sesuai).

---

## 6. Pertanyaan yang Perlu Diputuskan

1. **Lapisan 1 — pembedaan tipe:** setuju pakai `in` (delta naik) / `void_purchase` (delta turun), bukan `adjustment`? Atau cukup item baru saja yang `in`, sisanya biarkan?
2. **Lapisan 2 — data lama:** pilih **2a** (ubah backfill, otomatis) atau **2b** (manual lewat menu)? Ini mengubah skrip backfill, jadi butuh izin eksplisit.
3. **Test tambahan:** skenario edit mana lagi yang mau diuji sebelum implementasi (turunkan qty / hapus produk saat edit)?

---

## 7. Catatan Uji / Baseline (SEBELUM perbaikan)

Test lengkap dijalankan via browser (login UI Playwright + request sesi sama), 5 skenario,
lalu tipe mutasi diperiksa langsung di `stock_mutations`. Semua data test sudah dibersihkan,
DB dikembalikan ke kondisi semula. Hasil baseline (perilaku SEKARANG):

| Skenario | Aksi edit | Tipe mutasi tercatat SEKARANG | Seharusnya (target perbaikan) | Bug? |
|----------|-----------|-------------------------------|-------------------------------|------|
| SCN1 | Tambah produk baru | `in` + **`adjustment`** | `in` + **`in`** | ❌ |
| SCN2 | Naikkan qty (5→9) | `in` + **`adjustment`** (+4) | `in` + **`in`** (+4) | ❌ |
| SCN3 | Turunkan qty (8→3) | `in` + **`adjustment`** (−5) | `in` + **`void_purchase`** (5) | ❌ |
| SCN4 | Hapus produk (−6) | `in` + **`adjustment`** (−6) | `in` + **`void_purchase`** (6) | ❌ |
| SCN5 | Void PO | `in` + **`void_purchase`** | `in` + `void_purchase` | ✅ (sudah benar) |

**Kesimpulan baseline:**
1. Bug lebih luas dari dugaan awal — BUKAN cuma "tambah produk", tapi SEMUA jenis edit
   (tambah/naik/turun/hapus) tercatat `adjustment`. Ini menguatkan Tabel di Bagian 4.
2. Angka stok SECARA NILAI sudah benar di semua skenario (stock_before/after berurutan) —
   bug ini TIDAK merusak angka stok operasional, hanya salah melabeli TIPE mutasi.
   Masalah baru muncul saat migrasi/rekonstruksi yang mengandalkan tipe.
3. `Void` (SCN5) sudah benar (`void_purchase`) → perbaikan `Update` TIDAK boleh menyentuh
   `Void`, dan target tipe delta-negatif memang `void_purchase` (konsisten dengan Void).
4. Baseline ini jadi acuan pembanding: setelah perbaikan, SCN1-4 harus berubah ke kolom
   "Seharusnya", dan SCN5 harus TETAP sama (tidak boleh berubah/rusak).

Catatan: perbaikan `Update` harus dites ulang untuk kelima skenario ini + memastikan
angka stok akhir tetap identik dengan baseline (hanya tipe mutasi yang berubah).

---

## 8. Kaitan dengan Dokumen Lain

- `docs/ANALISIS_REKONSILIASI_STOK_32_PRODUK.md` — daftar 32 produk & rekomendasi angka. **Perlu dikoreksi narasinya**: penyebab sebenarnya adalah bug edit-PO (`adjustment`), bukan "pembelian lama pra-ledger".
- `docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md` — panduan migrasi (tidak terpengaruh langsung, tapi jumlah `needs_stock_review` akan berkurang kalau Lapisan 2 dijalankan).
