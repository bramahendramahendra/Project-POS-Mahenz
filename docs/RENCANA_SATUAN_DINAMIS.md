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

**Update evaluasi: BISA dihandle** — dengan syarat desain tidak memaksa satu unit cuma boleh punya satu rasio. Solusinya: izinkan beberapa "paket" untuk unit yang sama, dibedakan label:
- Paket "Slop (Lama)" → 1 Slop = 10 Pack
- Paket "Slop (Baru)" → 1 Slop = 5 Pack

Keduanya aktif bersamaan, kasir/pembelian pilih paket yang sesuai fisik barangnya saat transaksi — stok terpotong benar sesuai paket yang dipilih.

**Yang tetap jadi tanggung jawab manusia, bukan sistem:** kasir harus tahu box yang dipegang itu dari paket "Lama" atau "Baru" (dari label fisik/tanggal terima), baru pilih yang benar di layar. Software cuma menyediakan pilihannya, tidak bisa otomatis tahu isi fisik box yang belum dibuka — ini berlaku untuk solusi apapun (termasuk sistem lot sekalipun).

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

### Kasus 4: Konversi tidak presisi, variasi per item (Jeruk)

Tahun 1: per kg. Tahun 2: per buah (berat 150–250 gram, tidak seragam). Tahun 3: per peti grosir (±20 kg, tidak pasti).

**Hasil evaluasi: bisa dihandle secara mekanisme (sama seperti Telur), tapi beda kategori — perlu pendekatan beda.** Bedanya dari Telur: rasio Telur cukup stabil (sesekali direvisi), sedangkan berat jeruk **beda di setiap transaksi**, bukan cuma sesekali antar-batch. Kalau dipaksa satu rasio tetap, stok "buah" akan meleset terus-menerus.

**Rekomendasi:** produk jenis ini sebaiknya dijual **per kg dengan berat ditimbang langsung saat transaksi** (kasir input berat aktual), bukan qty × rasio tetap. "Buah" cuma keterangan tambahan opsional, bukan dasar hitungan stok. Peti (±20 kg) aman dipakai untuk estimasi pemesanan borongan, bukan potong stok presisi per transaksi retail.

### Kasus 5-12: batch kedua (Minyak Goreng, Cukai, BBM, Beras, Sabun, Kue Kering, Semangka, Semen)

| Kasus | Status | Catatan |
|---|---|---|
| 5. Minyak Goreng — sisip Pack di tengah Dus-Jerigen | ✅ Handled | Graph model memang untuk ini — Pack cukup diisi 1 pasangan, pasangan lain otomatis konsisten (6×2=12) |
| 6. Cukai Rokok — rasio berubah karena regulasi, variasi per tanggal produksi | ✅ Handled | Mekanisme sama Kasus 3 — multi-paket berlabel (per periode aturan), kasir pilih sesuai tanggal produksi |
| 7. BBM/Oli — isi ulang custom (1.7 liter bebas) | ⚠️ Bukan soal konversi satuan — bug terpisah | Bukan "pilih paket", tapi input qty desimal bebas. Kasir saat ini cuma terima qty bilangan bulat (`parseInt`, min 1) — perbaikan terpisah di form kasir, di luar redesain satuan |
| 8. Beras — karung beda ukuran per supplier | ✅ Handled | Multi-paket per unit (Karung A=25kg, Karung B=20kg), cuma relevan di pembelian; begitu jadi stok kg sudah tercampur rata |
| 9. Sabun — promo bonus isi (48→50 pcs) | ✅ Handled | Multi-paket lagi — "Dus Normal"=48, "Dus Promo"=50, pilih sesuai batch beli |
| 10. Kue Kering — naik kelas Toples→Dus→Kontainer | ✅ Handled | Paling sederhana, sama pola Kasus 1 — basis dari awal sudah terkecil, tinggal tambah lebih besar |
| 11. Semangka — dipotong, harga per estimasi berat potongan | ❌ Di luar kategori — fitur beda total | Bukan "satuan produk sama", tapi produk utuh diubah bentuk (potongan ditimbang) — kategori "produk olahan/turunan", tidak terkait sistem konversi satuan |
| 12. Semen — migrasi "sak" lama ke kg + sak per merek | ✅ Handled ke depan, tidak untuk data lama | Ke depan: multi-paket ("Sak Merek A"=40kg, dst). Data transaksi lama yang cuma catat "sak" tanpa keterangan tidak bisa diperjelas otomatis — keterbatasan data lama |

**Temuan penting:** "Multi-paket per unit berlabel" (dari Kasus 3) jadi mekanisme berulang — muncul di Kasus 3, 6, 8, 9, 12. Ini pola umum untuk rasio beda karena sumber (supplier), waktu (regulasi/promo), atau batch produksi — bukan solusi khusus rokok saja.

## Evaluasi ulang 5 kasus terhadap desain FINAL (graph/pasangan bebas)

| Kasus | Status | Catatan |
|---|---|---|
| 1. Gula Pasir | ✅ Handled | Sama seperti sebelumnya |
| 2. Air Mineral | ✅ Handled | Lebih simpel — tidak perlu rebase/hitung ulang stok sama sekali lagi, cukup tambah 1 pasangan |
| 3. Dua batch rasio beda (Rokok) | ✅ Handled (update) | Selama desain izinkan >1 paket untuk unit yang sama (label pembeda "Lama"/"Baru") — kasir pilih manual sesuai fisik barang |
| Mie Instan (multi-level) | ✅ Handled | Lebih simpel — tanpa rebase, hubungan antar level otomatis konsisten |
| Telur (rasio perkiraan) | ✅ Handled, catatan sama | Mekanisme lebih simpel, tapi catatan "estimasi bukan pasti" tetap berlaku (soal sifat data, bukan desain) |

Skor jadi **5/5** — desain final ini **menghilangkan semua kebutuhan rebase/migrasi stok**, dan dengan izinkan multi-paket per unit, kasus batch pun tertangani (selama kasir pilih paket yang benar saat transaksi — bagian ini tetap tanggung jawab manusia, bukan software).

## Ringkasan & Rekomendasi

**Solusi final:** jaringan konversi antar-satuan per produk, pasangan bebas, boleh lebih dari satu paket untuk unit yang sama (beda label/rasio, untuk kasus batch), "Satuan Dasar" jadi cuma pilihan tampilan (lihat bagian "Solusi FINAL" di atas). Tidak ada lagi konsep rebase/migrasi stok.

**Rekomendasi:**
1. **Kerjakan sekarang** (masih tahap dev, aman dilakukan): bangun fitur konversi berbasis pasangan ini dari awal, izinkan multi-paket per unit berlabel, sekalian perbaiki pembelian supaya rasio konversi selalu diverifikasi dari server (bukan percaya input layar).
2. **Boleh ditunda:** penanda "estimasi vs eksak" (Telur/Jeruk), form "tambah item ke PO existing" yang belum bisa pilih satuan.
3. **Di luar cakupan redesain satuan — kategori masalah beda, evaluasi terpisah kalau dibutuhkan:**
   - Kasir cuma terima qty bilangan bulat (`parseInt`, min 1) — perlu diubah kalau butuh jual custom/desimal (Kasus 7, BBM isi ulang).
   - "Produk olahan/turunan" (Kasus 11, Semangka dipotong & ditimbang) — fitur baru sama sekali, bukan bagian dari sistem konversi satuan.
   - Cost/harga modal beda per batch/supplier (disinggung di Kasus 8) — itu topik akuntansi biaya (costing), beda dari konversi qty.

## Status

Draf diskusi. Belum ada kode yang diubah.
