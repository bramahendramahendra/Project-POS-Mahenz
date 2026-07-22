# Rencana: Satuan Produk Dinamis

> Dokumen rancangan (bukan kode). Diupdate seiring diskusi & kasus baru.

## Masalah

Satuan Dasar produk saat ini tetap sejak awal dibuat. Semua satuan grosir (`conversion_qty`) wajib dihitung relatif ke basis itu dan harus ≥ 1 — sehingga tidak bisa menambah satuan yang lebih kecil dari basis di kemudian hari. Mengubah dropdown Satuan Dasar di form edit juga tidak memicu penyesuaian apa pun ke satuan lain atau ke stok (bug).

## Solusi FINAL: Jaringan Konversi Antar-Satuan (bebas pasangan, per produk)

> Menggantikan ide "Ubah Satuan Dasar (rebase)" di bawah — lebih sederhana, tidak perlu rebase/migrasi stok sama sekali.

Tiap produk punya jaringan konversi sendiri: satuan-satuan yang dipakai produk itu (Slop, Pack, Batang, dst) saling terhubung lewat pasangan (mis. "1 Slop = 10 Pack", "1 Pack = 12 Batang"), **bukan** selalu relatif ke satu Satuan Dasar tetap.

Aturan:
1. **Tambah satuan baru** — cukup isi rasio ke **satu** satuan yang sudah ada (bebas pilih ke satuan mana). Hubungan ke satuan lain yang belum diisi langsung dihitung otomatis oleh sistem (mis. Batang-Slop = 12×10 = 120), **tidak perlu diisi manual** — supaya tidak pernah kontradiksi antar pasangan.
2. **Ubah rasio** — edit angka di satu pasangan, satuan lain ikut menyesuaikan otomatis (dihitung ulang saat ditampilkan, bukan disimpan ulang).
3. **Stok tetap disimpan di satu satuan internal tetap** (satuan pertama yang dipakai produk itu, tidak pernah perlu dikonversi ulang) — tapi ini murni detail teknis, tidak terlihat user.
4. **"Satuan Dasar" jadi cuma pilihan tampilan** — bebas pilih satuan mana saja untuk ditampilkan sebagai acuan utama di layar produk, bisa diganti kapan saja, karena semua satuan sudah saling terhubung lewat jaringan yang sama. Ganti pilihan ini **tidak perlu hitung ulang atau migrasi stok apa pun**.
5. Riwayat transaksi/pembelian lama tetap tersimpan dengan rasio yang berlaku saat itu (snapshot per transaksi, seperti yang sudah berjalan sekarang) — tidak berubah walau jaringan konversi produk diedit belakangan.

### Kenapa ini lebih baik dari ide rebase sebelumnya
- Tidak ada operasi "berat" (rebase) yang perlu migrasi stok.
- Tidak ada pembatasan urutan (harus dari besar ke kecil atau sebaliknya) — satuan bisa ditambah dari arah mana saja, kapan saja.
- "Satuan Dasar" tidak lagi jadi konsep yang mengunci apa pun — murni pilihan tampilan yang bebas diganti.

<details>
<summary>Ide sebelumnya: "Ubah Satuan Dasar" (rebase) — disimpan sebagai riwayat diskusi, sudah digantikan di atas</summary>

Tambah kemampuan baru: Satuan Dasar produk bisa diganti kapan saja ke satuan yang lebih kecil. Saat diganti, sistem otomatis:
1. Menghitung ulang semua satuan lain yang sudah ada supaya tetap benar relatif ke Satuan Dasar baru (angka lama dikali rasio yang diinput).
2. Mengonversi stok yang sudah tercatat ke satuan baru (dikali rasio yang sama).
3. Riwayat transaksi/pembelian lama tidak berubah (rasio yang berlaku saat itu sudah tercatat permanen di baris transaksinya sendiri).

Satuan yang **lebih besar** dari yang sudah ada (mis. tambah "Karton" di atas Slop) tetap pakai alur "Tambah Paket Grosir" yang sudah ada — tidak berubah.

Mengubah **rasio** antar dua satuan yang sudah ada (bukan basis) — cukup edit angka satuan itu langsung, tidak perlu rebase, stok tidak ikut berubah.

</details>

## Contoh alur (rokok: Slop → Pack → rasio berubah → Batang)

| Tahun | Aksi | Basis | Stok | Satuan lain |
|---|---|---|---|---|
| 1 | Mulai jual Slop | Slop | 20 Slop | – |
| 2 | Ubah Satuan Dasar → Pack, "1 Slop = 10 Pack" | Pack | 200 Pack | Slop = 10 |
| 3 | Edit rasio Slop: 10 → 5 | Pack | 200 Pack (tetap) | Slop = 5 |
| 4 | Ubah Satuan Dasar → Batang, "1 Pack = 12 Batang" | Batang | 2.400 Batang | Pack = 12, Slop = 60 |

## Kasus yang bisa dihandle

- Tambah satuan lebih kecil dari basis, kapan saja, berkali-kali, urutan bebas.
- Tambah satuan lebih besar dari yang sudah ada (alur existing, tidak berubah).
- Rasio kemasan berubah, tanpa memengaruhi stok atau riwayat lama.
- Riwayat transaksi/pembelian lama tetap akurat meski basis produk berubah belakangan.

## Kasus yang di luar cakupan solusi ini

- Harga tingkat/grosir (tier price) tidak tersambung ke satuan tertentu di backend — masalah terpisah.
- Form pembelian saat ini percaya rasio konversi dari layar tanpa verifikasi ulang ke server — masalah terpisah.
- Form "tambah item ke PO yang sudah ada" tidak punya pilihan satuan — masalah terpisah.

## Kasus dari user (diisi & dievaluasi bertahap)

### Kasus 1: Arah terbalik — eceran ke grosir (Gula Pasir)

Tahun 1: jual per kg. Tahun 2: tambah Karung (50 kg). Tahun 3: rasio karung berubah jadi 25 kg. Tahun 4: tambah Pallet (20 Karung).

**Hasil evaluasi: sudah bisa dihandle tanpa fitur baru.** Basis (kg) dari awal sudah yang terkecil, semua tambahan berikutnya (Karung, Pallet) selalu lebih besar — ini alur "Tambah Paket Grosir" yang sudah ada sekarang, tidak butuh "Ubah Satuan Dasar".

| Tahun | Basis | Satuan lain |
|---|---|---|
| 1 | kg | – |
| 2 | kg | Karung = 50 |
| 3 | kg | Karung = 25 |
| 4 | kg | Karung = 25, Pallet = 500 (20×25, relatif ke kg bukan ke Karung) |

### Kasus 2: Rasio berubah naik-turun, arah konversi tetap (Air Mineral Botol)

Tahun 1: jual per Dus (isi 24 Botol). Tahun 2: mulai eceran per Botol. Tahun 3: isi dus jadi 12. Tahun 4: isi dus naik lagi jadi 15.

**Hasil evaluasi: bisa dihandle, cuma butuh 1× rebase (Tahun 2), sisanya edit angka.** Begitu Botol (terkecil) jadi basis, perubahan rasio Dus berikutnya — naik maupun turun — sama-sama cuma edit satu angka, tidak perlu rebase ulang.

| Tahun | Aksi | Basis | Satuan lain |
|---|---|---|---|
| 1 | Jual per Dus | Dus | – |
| 2 | Mulai eceran per Botol → rebase ke Botol, "1 Dus = 24 Botol" | Botol | Dus = 24 |
| 3 | Isi dus jadi 12 → edit angka | Botol | Dus = 12 |
| 4 | Isi dus naik jadi 15 → edit angka | Botol | Dus = 15 |

### Kasus 3: Dua batch beda rasio beredar bersamaan (Rokok, kemasan lama vs baru)

Stok lama (1 Slop = 10 Pack) dan stok baru (1 Slop = 5 Pack) ada bersamaan di rak.

**Hasil evaluasi: TIDAK sepenuhnya bisa dihandle — limitasi nyata.** Solusi rebase cuma simpan satu rasio "saat ini" per produk, tidak ada konsep rasio berbeda per batch.
- Aman selama stok lama sudah dipecah jadi Pack saat diterima dulu (rasio waktu itu sudah kepakai, jadi angka Pack).
- Bermasalah kalau Slop lama dijual **utuh** setelah rasio berubah — kasir potong stok pakai rasio sekarang (5), padahal fisik Slop itu isinya 10 Pack → stok kelebihan tercatat diam-diam.
- Solusi sebenarnya: stok per-batch/lot (tiap kulakan simpan rasio sendiri) — fitur terpisah, jauh lebih besar dari rebase satuan. Di luar cakupan rencana ini.

### Kasus Mie Instan: multi-level (Dus = 4 Renceng = 20 Bungkus)

Tahun 1: jual per Dus (40 Bungkus). Tahun 2: eceran per Bungkus. Tahun 3: dus jadi isi 20 Bungkus. Tahun 4: tambah Renceng (5 Bungkus).

**Hasil evaluasi: bisa dihandle penuh**, tanpa fitur tambahan.

| Tahun | Aksi | Basis | Satuan lain |
|---|---|---|---|
| 1 | Jual per Dus | Dus | – |
| 2 | Rebase ke Bungkus, "1 Dus=40 Bungkus" | Bungkus | Dus=40 |
| 3 | Edit Dus: 40→20 | Bungkus | Dus=20 |
| 4 | Tambah Renceng (5 Bungkus) — alur grosir biasa, karena lebih besar dari basis | Bungkus | Dus=20, Renceng=5 |

"1 Dus = 4 Renceng" otomatis konsisten (20÷5=4) karena semua relatif ke basis yang sama, tidak perlu disimpan terpisah.

### Kasus Telur: rasio berat↔jumlah yang tidak eksak

Tahun 1: per Peti (±15 kg). Tahun 2: per kg. Tahun 3: peti jadi 10 kg. Tahun 4: per butir (1 kg ≈ 16 butir, bisa berubah-ubah).

**Hasil evaluasi: bisa dihandle secara mekanisme, dengan catatan akurasi.** Mekanismenya sama seperti kasus lain (rebase 2×: Peti→kg→butir). Bedanya: rasio kg↔butir itu **perkiraan**, bukan hitungan fisik eksak seperti Dus/Pack/Slop — bisa beda tiap batch (telur kecil vs jumbo). Solusi tetap jalan (rasio diedit tiap kali beda), tapi stok dalam satuan "butir" jadi estimasi, bukan pasti. Ini batasan yang perlu disadari, bukan gap fitur.

## Ringkasan & Rekomendasi

**Solusi inti** — 3 mekanisme, semua sudah ada di aplikasi kecuali #1:
1. **Ubah Satuan Dasar (rebase)** — fitur baru. Dipakai saat mau jual satuan yang lebih kecil dari basis saat ini. Semua satuan lain & stok ikut terkonversi otomatis, riwayat lama tidak berubah.
2. **Tambah Paket Grosir** — sudah ada. Dipakai saat mau tambah satuan yang lebih besar dari yang sudah ada.
3. **Edit angka rasio** — sudah ada. Dipakai saat rasio kemasan berubah tapi arah/urutan satuan tetap sama.

**Skor dari 5 kasus yang diuji:** 4 dari 5 kasus (Gula Pasir, Air Mineral, Mie Instan, Telur) tertangani penuh oleh 3 mekanisme di atas — termasuk kasus multi-level (Dus=Renceng=Bungkus) dan rasio naik-turun berkali-kali. Cuma **Kasus 3 (dua batch rasio beda beredar bersamaan)** yang tidak tertangani, karena itu butuh pelacakan stok per-batch/lot — kategori fitur yang berbeda sama sekali dari urusan konversi satuan.

**Rekomendasi:**
1. **Kerjakan sekarang** (masih tahap dev, aman dilakukan): fitur Ubah Satuan Dasar + kunci dropdown Satuan Dasar di form edit + perbaiki pembelian supaya rasio konversi selalu diverifikasi dari server (bukan percaya input layar). Tiga ini saling terkait dan sama-sama menyentuh titik yang sama (konversi satuan & stok), jadi masuk akal dikerjakan sekaligus.
2. **Putuskan dulu, jangan langsung dikerjakan:** apakah kasus "dua batch rasio beda" (Kasus 3) benar-benar akan terjadi di operasional toko kamu? Kalau toko selalu memecah kemasan besar jadi satuan kecil segera saat barang diterima (praktik umum), limitasi ini tidak pernah kerasa. Kalau toko memang sengaja menyimpan & menjual kemasan besar utuh dalam jangka panjang sementara rasio sering berubah, ini perlu direncanakan sebagai fitur terpisah (stok per-lot) — lebih besar scope-nya, sebaiknya dipikirkan sebelum ada data produksi supaya tidak perlu migrasi data nanti.
3. **Boleh ditunda:** penanda "estimasi vs eksak" untuk kasus seperti Telur (kg↔butir), dan perbaikan form "tambah item ke PO existing" yang belum bisa pilih satuan. Ini nice-to-have, tidak menghalangi kasus inti.

## Status

Draf diskusi. Belum ada kode yang diubah.
