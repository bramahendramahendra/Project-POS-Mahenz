# Analisis Rekonsiliasi Stok — 32 Produk `needs_stock_review`

> Laporan analisis (Opsi B). **Tidak ada data yang diubah** — ini murni rekomendasi.
> Keputusan akhir & eksekusi koreksi ada di tangan Anda (lewat menu Rekonsiliasi Stok).
> Dibuat: hasil migrasi stok prod → skema baru.

---

## 1. Kesimpulan Utama

> **KOREKSI NARASI (update final, berdasarkan investigasi langsung pada DB prod).**
> Versi-versi sebelumnya sempat menuduh penyebabnya "pembelian pra-ledger" lalu
> "bug edit-PO `adjustment`". Keduanya TIDAK tepat untuk data prod ini: setelah dicek,
> **DB prod TIDAK punya satu pun baris `adjustment`** — isinya hanya `in/purchase`,
> `out/transaction`, `void/transaction`. Penyebab sebenarnya di bawah ini.
> **Angka rekomendasi di Bagian 3–4 tetap valid** (dasarnya membandingkan stok lama vs
> ledger); yang dikoreksi hanya penjelasan *kenapa* dan *bagaimana* ditangani.

Akar masalah **bukan** status pembayaran (hutang/lunas) secara langsung — pembelian hutang
tetap menambah stok.

**Penyebab utama (terverifikasi): sebagian `purchase_items` tidak punya baris mutasi `in`
di ledger.** Ada 24 baris item pembelian (dari 4 PO active) yang barangnya benar-benar
dibeli dan tercatat di `purchase_items`, TAPI jejak stok masuknya (`in`) tidak pernah masuk
ke `stock_mutations` (kemungkinan besar akibat alur lama yang gagal menuliskan sebagian
mutasi saat PO dibuat/diedit). Akibatnya saat backfill me-replay riwayat, barang itu "tak
terlihat" (hanya penjualan `out` yang tercatat) → stok minus → produk ditandai
`needs_stock_review`. Ini menjelaskan korelasi dengan PO hutang/partial (PO besar/kompleks
yang mutasinya paling banyak bolong, mis. PO-20260822-003 yang 17 item-nya kehilangan `in`).

**Penyebab lain (tidak bisa ditambal otomatis):**
- **Rantai satuan bercabang (celah #10)** — struktur `ref_package_id` tidak linear; produk
  tidak bisa direkonstruksi otomatis sama sekali.
- **Paket satuan tidak ditemukan** — ada penjualan lama memakai `package_id` yang kini sudah
  tidak ada pada produk (satuan berubah/dihapus sejak transaksi).
- **Selisih di luar wajar** — hasil hitung ulang masih beda jauh dari stok lama (penjualan
  benar-benar melebihi total pembelian tercatat).

**Sudah ditangani otomatis.** Penyebab utama (missing `in`) kini ditambal oleh skrip
`BE/cmd/backfill_missing_purchase_in` (menyisipkan baris `in` yang hilang dari data
`purchase_items`, idempotent, `created_at` di-backdate ke tanggal PO). Setelah skrip ini +
`backfill_stock_restore` dijalankan, **`needs_stock_review` turun dari 32 → 12**. Sisa 12
produk adalah kasus "tidak bisa otomatis" di atas — inilah yang masih perlu Rekonsiliasi
Stok manual. Detail alur: `docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md` (langkah 5b).

**Bukti kuat (tetap berlaku):** kolom `stock_after` pada mutasi terakhir di ledger
**hampir selalu sama persis dengan stok lama** di backup. Artinya **stok lama umumnya akurat**
dan bisa dipakai sebagai angka koreksi untuk sisa produk yang masih perlu ditinjau.

---

## 2. Cara Membaca Rekomendasi

Untuk tiap produk saya bandingkan 3 sumber:
- **Stok Lama** = `products_stock_backup.stock` (angka sistem lama sebelum migrasi).
- **Ledger terakhir** = `stock_after` mutasi paling akhir (kondisi stok terakhir yang tercatat sistem).
- **Beli − Jual** = total pembelian aktif dikurangi total penjualan (kalau bisa dihitung).

Tingkat keyakinan:
- 🟢 **TINGGI** — stok lama = ledger terakhir (dua sumber cocok). Rekomendasi: pakai angka itu.
- 🟡 **SEDANG** — stok lama & ledger beda tipis, atau ledger lengkap tapi beda dari backup. Perlu Anda putuskan mana yang benar.
- 🔴 **CEK FISIK** — data tidak cukup / bertentangan. Sebaiknya hitung fisik di gudang.

> Semua angka dalam **satuan dasar** produk.

---

## 3. Rekomendasi per Produk

### Grup A — Keyakinan TINGGI (stok lama = ledger terakhir) 🟢
Rekomendasi: **isi stok = Stok Lama**. Dua sumber independen cocok.

| ID | Produk | Stok Lama | Ledger Terakhir | Rekomendasi |
|----|--------|-----------|-----------------|-------------|
| 1 | Toppas Merah 12 Kretek | 4.917 | 4.917 | **4.917** |
| 22 | Aqua 1600ml | 11.000 | 11.000 | **11** |
| 24 | Surya 12 | 7.166 | 7.166 | **7.166** |
| 53 | Kerupuk | 49.000 | 49.000 | **49** |
| 63 | Gula | 9.000 | 9.000 | **9** |
| 64 | Crispy Crackers | 2.000 | 2.000 | **2** |
| 83 | Telur | 0.087 | 0.087 | **0.087** |
| 112 | Kapal Api Classic | 0.950 | 0.950 | **0.95** |
| 139 | Dunhill Hitam 16 | 0.900 | 0.900 | **0.9** |
| 172 | Kopi ABC Botol | 17.910 | 17.910 | **17.91** |
| 173 | Top Kopi Susu | 1.000 | 1.000 | **1** |
| 217 | Hansaplast | 226.000 | 226.000 | **226** |
| 235 | Head & Shoulders 160mL | 0.000 | 0.000 | **0** |
| 241 | Indomie Hype Abis | 10.000 | 10.000 | **10** |
| 242 | LA Putih 12 | 12.000 | 12.000 | **12** |
| 244 | 234 Magnum Kretek 12 | 3.833 | 3.833 | **3.833** |
| 245 | Box Kotakan Kertas Nasi | 21.000 | 21.000 | **21** |
| 246 | Mika Ikan Lux | 2.990 | 2.990 | **2.99** |
| 248 | Stapler EP-10 | 1.000 | 1.000 | **1** |
| 251 | Kertas Alas Makanan Samir | 7.900 | 7.900 | **7.9** |

> Catatan: angka pecahan seperti 4.917 / 7.166 / 0.087 / 0.95 adalah efek satuan turunan (mis. rokok per batang, telur per kg). Itu wajar dari sistem lama — kalau Anda mau dibulatkan sesuai fisik, silakan sesuaikan saat input.

### Grup B — Tidak ada ledger sama sekali (murni pembelian lama) 🟢/🟡
Produk ini **tidak punya mutasi apa pun** di ledger (jml_mutasi=0). Berarti belum pernah ada penjualan tercatat di sistem baru, cuma pembelian lama. Rekomendasi: **isi stok = Stok Lama** (yang praktis = total beli, karena belum ada penjualan tercatat).

| ID | Produk | Stok Lama | Total Beli | Rekomendasi |
|----|--------|-----------|-----------|-------------|
| 47 | Marlboro Black Motion 12 | 6.000 | 6 | **6** |
| 55 | Susu Dancow Sachet Putih | 1.000 | 1 | **1** |
| 62 | Tepung Segitiga | 13.000 | 13 | **13** |
| 237 | Kurma First Dates | 1.000 | 1 | **1** |
| 243 | LA Putih 16 | 4.000 | 4 | **4** |
| 247 | Stapler HD-10-Mini | 2.000 | 2 | **2** |
| 249 | Stapler HD-50-C | 1.000 | 1 | **1** |
| 250 | Kertas Alas Makanan Padi | 2.000 | 2 | **2** |

### Grup C — Perlu keputusan Anda (ledger ≠ stok lama) 🟡🔴

| ID | Produk | Stok Lama | Ledger Terakhir | Catatan & Rekomendasi |
|----|--------|-----------|-----------------|-----------------------|
| 70 | LA Bold 20 | **-0.100** | 1.000 | Stok lama MINUS (efek bug presisi). Ledger terakhir `out` 0.1 dari 1.1 → **1.0**. 🟡 Rekomendasi: **1** (ikuti ledger, bukan angka minus). |
| 140 | Dji Samsoe Hitam 16 | 2.000 | **1.000** | Ledger LENGKAP & berurutan (in 3, out 1, out 1 → 1.0). Backup 2.0 kemungkinan belum ter-update. 🟡 Rekomendasi: **1** (ikuti ledger yang lengkap). |
| 236 | Marlboro Black Motion 20 | 4.990 | 0.490 | Ledger cuma 1 `out` 0.01 dari `stock_before` 0.5 → 0.49. Tapi stok lama 4.99 & total beli 5. Ada penjualan pecahan aneh (0.01). 🔴 **CEK FISIK** — kemungkinan besar sekitar 4.99, tapi angka 0.49 di ledger mencurigakan. |
| 238 | Basreng | 5.900 | 5.000 | Ledger 1 `out` 0.1 dari 5.1 → 5.0. Stok lama 5.9, total beli 6 (2 PO). 🟡 Rekomendasi: **5.9** (ikuti stok lama; ledger cuma catat 1 penjualan kecil). |

> Grup C sengaja dipisah karena angkanya bertentangan antar sumber. Untuk 236 khususnya, saya sarankan cek fisik karena pola penjualan 0.01 (satuan pecahan sangat kecil) mencurigakan dan bisa jadi salah input di masa lalu.

---

## 4. Ringkasan Angka Rekomendasi (siap input)

Kalau Anda setuju dengan rekomendasi di atas, berikut daftar angka final untuk di-input di menu Rekonsiliasi Stok (satuan dasar):

```
[1]   Toppas Merah 12 Kretek        = 4.917
[22]  Aqua 1600ml                   = 11
[24]  Surya 12                      = 7.166
[47]  Marlboro Black Motion 12      = 6
[53]  Kerupuk                       = 49
[55]  Susu Dancow Sachet Putih      = 1
[62]  Tepung Segitiga               = 13
[63]  Gula                          = 9
[64]  Crispy Crackers               = 2
[70]  LA Bold 20                    = 1        (bukan -0.1)
[83]  Telur                         = 0.087
[112] Kapal Api Classic             = 0.95
[139] Dunhill Hitam 16              = 0.9
[140] Dji Samsoe Hitam 16           = 1        (ikuti ledger, bukan 2)
[172] Kopi ABC Botol                = 17.91
[173] Top Kopi Susu                 = 1
[217] Hansaplast                    = 226
[235] Head & Shoulders 160mL        = 0
[236] Marlboro Black Motion 20      = CEK FISIK (kandidat 4.99)
[237] Kurma First Dates             = 1
[238] Basreng                       = 5.9
[241] Indomie Hype Abis             = 10
[242] LA Putih 12                   = 12
[243] LA Putih 16                   = 4
[244] 234 Magnum Kretek 12          = 3.833
[245] Box Kotakan Kertas Nasi       = 21
[246] Mika Ikan Lux                 = 2.99
[247] Stapler HD-10-Mini            = 2
[248] Stapler EP-10                 = 1
[249] Stapler HD-50-C               = 1
[250] Kertas Alas Makanan Padi      = 2
[251] Kertas Alas Makanan Samir     = 7.9
```

- **30 produk**: rekomendasi jelas (Grup A + B + kasus C yang sudah ada rekomendasi).
- **1 produk (236)**: disarankan cek fisik.
- **1 produk (140)**: perhatikan — rekomendasi ikut ledger (1), berbeda dari stok lama (2).

---

## 5. Cara Menerapkan (di menu Rekonsiliasi Stok)

Per produk:
1. Buka menu **Pelaporan → Rekonsiliasi Stok** (login admin).
2. Klik **Tinjau** pada produk.
3. Di bagian **Koreksi Manual**, isi kolom "Stok Benar" (satuan dasar) sesuai angka rekomendasi.
4. Klik **Simpan Koreksi Stok**.

Setiap koreksi otomatis:
- Update stok produk.
- Menghapus flag `needs_stock_review`.
- Mencatat mutasi `adjustment` (ada jejak audit: siapa, kapan, dari berapa ke berapa).

---

## 6. Catatan Penting

- Angka pecahan (4.917 dst) berasal dari **satuan turunan** (produk dijual per batang/gram sementara satuan dasarnya lebih besar). Itu bukan error — tapi kalau menurut Anda tidak masuk akal secara fisik, koreksi manual saat input.
- Produk Grup C (70, 140, 236, 238) adalah satu-satunya yang datanya bertentangan. Untuk ini pertimbangkan cek fisik, terutama **236**.
- Penyebab utama (item pembelian yang mutasi `in`-nya hilang dari ledger) **sudah ditangani
  otomatis** oleh `backfill_missing_purchase_in` (lihat `docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md`
  langkah 5b), sehingga `needs_stock_review` sudah turun dari 32 → 12. Sisa 12 inilah yang
  masih perlu Rekonsiliasi Stok manual (branching chain, paket tak ditemukan, atau selisih riil).
- Terpisah: bug tipe mutasi saat **Edit PO** juga sudah diperbaiki (edit PO kini mencatat
  `in`/`void_purchase`, bukan `adjustment` — lihat `docs/PERBAIKAN_MUTASI_STOK_EDIT_PEMBELIAN.md`,
  Lapisan 1). Ini pencegahan agar transaksi ke depan tidak menghasilkan tipe mutasi keliru;
  bukan penyebab 32 produk di data prod ini (yang memang tidak punya baris `adjustment`).
