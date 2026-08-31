# Rencana Perbaikan: Mutasi Stok Salah Tipe saat Edit Pembelian

> Status: **LAPISAN 1 SELESAI & TERVERIFIKASI (via UI browser).** Lapisan 2 (data lama)
> masih menunggu keputusan (Opsi 2a vs 2b, lihat Bagian 4).
> Ditemukan saat investigasi kenapa banyak produk `needs_stock_review` setelah migrasi stok.

---

## 1. Ringkasan Masalah

Saat sebuah Purchase Order (PO) **diedit** (lewat fungsi `Update`), SEMUA perubahan stok
akibat edit dicatat di `stock_mutations` dengan tipe **`adjustment`**, bukan tipe pembelian
yang benar (`in` untuk penambahan / `void_purchase` untuk pengurangan). Ini berlaku untuk
**semua jenis edit**: menambah produk baru, menaikkan qty, menurunkan qty, maupun menghapus
produk (lihat baseline lengkap di Bagian 7). Angka stok tetap benar — yang salah hanya **tipe**
mutasinya.

Akibatnya:
1. **Kartu stok jadi salah makna** — pembelian tampak seperti "koreksi manual", bukan pembelian dari supplier.
2. **Rekonstruksi stok (backfill migrasi) gagal** untuk produk itu — karena backfill tidak bisa memproses tipe `adjustment`, produk jadi minus saat replay dan ditandai `needs_stock_review`.

> **Inkonsistensi di dalam kode sendiri:** menambah produk lewat tombol **"Tambah Item"**
> (fungsi `AddItems`, untuk PO status partial/paid) SUDAH mencatat `in` dengan benar. Hanya
> jalur **Edit PO** (`Update`) yang salah pakai `adjustment`. Jadi target perbaikan (`in`/
> `void_purchase`) sudah punya preseden benar di kode yang sama — lihat Bagian 3.

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
| **Edit PO — semua perubahan item (`purchase_repo.go` `Update`)** | `purchase` | Penyesuaian stok karena PO diedit | ⚠️ **Salah — semua jenis edit** |
| Koreksi manual menu Rekonsiliasi Stok | `stock_reconciliation` | Koreksi hasil tinjauan | ✅ |

Sebagai pembanding, jalur pembelian yang SUDAH benar (tidak pakai `adjustment`):

| Lokasi | `mutation_type` | Keterangan |
|--------|-----------------|-----------|
| Buat PO baru (`purchase_repo.go` `Create`) | `in` | ✅ benar |
| Tambah Item ke PO (`purchase_repo.go` `AddItems`) | `in` | ✅ benar — jadi acuan perbaikan |
| Void PO (`purchase_repo.go` `Void`) | `void_purchase` | ✅ benar |

Kesimpulan: `adjustment` **sah untuk koreksi manual**, tapi **salah dipakai untuk pembelian**. Perubahan item lewat Edit PO adalah **pembelian nyata** (barang masuk/berkurang dari supplier) → harus `in` (nambah) / `void_purchase` (kurang), sama seperti `AddItems`/`Void`.

---

## 3. Lokasi Kode Sumber Masalah

File: `BE/domain/supplier_purchase/repo/purchase_repo.go`, fungsi `Update()` — blok "Terapkan hanya SELISIH bersih per (produk, paket)" (sekitar baris 517-555).

Logika menghitung `delta = newQty - oldQty` per (produk, paket) tetap dipertahankan
(pendekatan "delta bersih"). Yang diperbaiki hanya **pemilihan `MutationType`**.

**SEBELUM (bug):** semua delta dicatat `MutationType: "adjustment"`, termasuk item yang
benar-benar baru.

**SESUDAH (implementasi Lapisan 1):**

```go
// (ringkas) — versi SESUDAH perbaikan
for k := range keys {
    delta := newQty[k] - oldQty[k]
    if delta == 0 { continue }
    direction := StockIn; qty := delta
    mutationType := "in"                                    // delta positif = pembelian tambah
    notes := "Edit PO ID ... -- tambah pembelian (selisih item baru vs lama)"
    if delta < 0 {
        direction = StockOut; qty = -delta
        mutationType = "void_purchase"                      // delta negatif = pembelian dikurangi
        notes = "Edit PO ID ... -- kurangi pembelian (selisih item baru vs lama)"
    }
    ApplyStockDelta(... MutationType: mutationType, ReferenceType: "purchase" ...)
}
```

Pola ini konsisten dengan `AddItems()` (`in`) dan `Void()` (`void_purchase`) di file yang
sama, sehingga rekonstruksi stok (backfill/kartu stok) bisa menelusuri perubahan edit-PO
sebagai pembelian, bukan koreksi manual.

---

## 4. Rencana Perbaikan

### Lapisan 1 — Perbaiki `Update()` (root cause, untuk data BARU) — ✅ SELESAI & TERVERIFIKASI

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

**Opsi 2a — Perbaiki skrip backfill lalu jalankan ulang (otomatis)**
Ajari `backfill_stock_restore` supaya mutasi `adjustment` dengan `reference_type='purchase'` diperlakukan sebagai pembelian:
- `adjustment` + `reference_type='purchase'` + arah menambah → perlakukan seperti `in`
- `adjustment` + `reference_type='purchase'` + arah mengurangi → perlakukan seperti `void_purchase`
- Arah (+/−) ditentukan dari `stock_after` vs `stock_before` pada baris mutasi itu (JANGAN diasumsikan).

> ⚠️ **KOMPLIKASI TEKNIS PENTING (harus ditangani, kalau tidak justru salah hitung).**
> Backfill sekarang memproses `in`/`void_purchase` lewat fungsi `purchaseItemEvents`, yang
> membaca **qty FINAL dari `purchase_items`** (`SELECT quantity FROM purchase_items WHERE
> purchase_id=? AND product_id=?`) — BUKAN qty delta yang tercatat di baris mutasi. Untuk
> mutasi `in` asli itu benar. TAPI untuk kasus **edit menaikkan qty** (mis. 5→9), ledger
> berisi `in`=5 (create) + `adjustment`=4 (delta). Kalau `adjustment` diproses lewat
> `purchaseItemEvents` yang sama, ia akan membaca qty final purchase_items (=9) dan menambah
> 9 LAGI → total jadi 5+9=14, padahal benar 9 (dobel hitung).
>
> Jadi Opsi 2a TIDAK cukup hanya "mapping adjustment→in". Skrip harus memproses mutasi
> `adjustment`-purchase memakai **qty DELTA dari baris `stock_mutations` itu sendiri**
> (`stock_after − stock_before`, sudah dalam satuan anchor), BUKAN lewat `purchaseItemEvents`.
> Ini menambah kompleksitas: perlu jalur replay khusus untuk adjustment-purchase yang tidak
> membaca ulang purchase_items.

Kelebihan: mayoritas 32 produk beres otomatis tanpa input manual, DAN backfill jadi tahan
terhadap data campuran (`in` + `adjustment` di PO yang sama). Kekurangan: **mengubah skrip
backfill** (perlu izin — user sebelumnya menekankan jangan ubah skrip backfill yang sudah ada),
dan komplikasi qty-delta di atas menaikkan risiko bug kalau tidak hati-hati.

> CATATAN penting soal urutan: perbaikan Lapisan 1 (Update) hanya mengubah perilaku data BARU.
> Data lama tetap `adjustment` sampai Lapisan 2 dijalankan. Kedua lapisan saling melengkapi,
> bukan menggantikan.

**Opsi 2b — Koreksi manual lewat menu Rekonsiliasi Stok**
Input satu-satu pakai angka di `docs/ANALISIS_REKONSILIASI_STOK_32_PRODUK.md`. Kelebihan: tidak menyentuh skrip, tidak ada risiko dobel-hitung. Kekurangan: manual, lebih lama.

---

## 5. Dampak & Cakupan Perubahan

**Lapisan 1 (Update):**
- File: `BE/domain/supplier_purchase/repo/purchase_repo.go` — hanya blok pemilihan `MutationType` di `Update()` (sekitar baris 545, di dalam loop `for k := range keys`).
- Tidak mengubah alur `Create()` (sudah benar: pakai `in`).
- Tidak mengubah `AddItems()` (sudah benar: pakai `in` — bisa jadi acuan pola).
- Tidak mengubah `Void()` (sudah benar: pakai `void_purchase`).
- Perlu test ulang skenario: tambah produk, naikkan qty, turunkan qty, hapus produk saat edit — plus **regresi**: pastikan Create, AddItems, dan Void TIDAK berubah perilakunya.

**Lapisan 2 (kalau Opsi 2a dipilih):**
- File: `BE/cmd/backfill_stock_restore/main.go` — tambah handling `adjustment` bersumber `purchase`.
- Setelah itu jalankan ulang backfill (butuh reset/kondisi yang sesuai).

---

## 6. Pertanyaan yang Perlu Diputuskan

1. ~~**Lapisan 1 — pembedaan tipe**~~ — **SELESAI.** Disepakati & diimplementasi: `in`
   (delta naik) / `void_purchase` (delta turun), bukan `adjustment`. Terverifikasi via UI
   browser (Bagian 7b).
2. **Lapisan 2 — data lama:** pilih **2a** (ubah backfill, otomatis) atau **2b** (manual lewat menu)? Ini mengubah skrip backfill, jadi butuh izin eksplisit. **← MASIH MENUNGGU KEPUTUSAN.**
3. ~~Test tambahan~~ — **SUDAH SELESAI.** Baseline 5 skenario (tambah/naik/turun/hapus/void)
   sudah diuji via UI browser, lihat Bagian 7. Tidak perlu test tambahan sebelum implementasi.

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

## 7b. Hasil Uji SESUDAH Perbaikan (Lapisan 1) — ✅ SEMUA LULUS

Setelah Lapisan 1 diimplementasi, `go build ./...` bersih (exit 0), lalu BE (kode baru) + FE
dijalankan dan kelima skenario diuji ULANG via UI browser (Playwright, login admin, sesi
sama). Tipe mutasi diperiksa langsung di `stock_mutations`. Semua data test dibersihkan &
DB dikembalikan ke kondisi semula (baseline `max_mut=816`, `max_po=154`, stok pkg4=7, pkg5=3
terverifikasi pulih 100%).

Produk uji: id 2 (Toppas Merah 16 Kretek, pkg 4) & id 3 (Geo Mild 16, pkg 5). Semua PO
berstatus `unpaid` (Hutang).

| Skenario | Aksi edit | Tipe mutasi edit SESUDAH | Target | Hasil |
|----------|-----------|--------------------------|--------|-------|
| SCN1 | Tambah produk baru (Geo, +3) | **`in`** | `in` | ✅ LULUS |
| SCN2 | Naikkan qty (2→5, +3) | **`in`** | `in` | ✅ LULUS |
| SCN3 | Turunkan qty (5→2, −3) | **`void_purchase`** | `void_purchase` | ✅ LULUS |
| SCN4 | Hapus produk (Geo, −3) | **`void_purchase`** | `void_purchase` | ✅ LULUS |
| SCN5 | Void PO | **`void_purchase`** | `void_purchase` | ✅ LULUS (tak berubah) |

**Kesimpulan sesudah perbaikan:**
1. **Tidak ada lagi `adjustment`** di seluruh mutasi hasil edit-PO. Semua terlabeli `in`
   (delta positif) atau `void_purchase` (delta negatif) — persis target.
2. Notes mutasi sudah deskriptif & benar: *"...tambah pembelian (selisih item baru vs lama)"*
   untuk `in`, *"...kurangi pembelian (selisih item baru vs lama)"* untuk `void_purchase`.
3. Angka stok tetap konsisten (stock_before/after berurutan tanpa lompatan janggal) di semua
   skenario — perbaikan hanya mengubah TIPE, tidak merusak nilai stok.
4. **Regresi aman:** Create (`in`) dan Void (`void_purchase`) TIDAK berubah perilakunya
   (SCN5 mengonfirmasi keduanya masih benar).

---

## 8. Kaitan dengan Dokumen Lain

- `docs/ANALISIS_REKONSILIASI_STOK_32_PRODUK.md` — daftar 32 produk & rekomendasi angka. **Perlu dikoreksi narasinya**: penyebab sebenarnya adalah bug edit-PO (`adjustment`), bukan "pembelian lama pra-ledger".
- `docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md` — panduan migrasi (tidak terpengaruh langsung, tapi jumlah `needs_stock_review` akan berkurang kalau Lapisan 2 dijalankan).
