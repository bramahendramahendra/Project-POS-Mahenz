# Rencana Perbaikan: Presisi Stok Multi-Satuan

> Dokumen rancangan (bukan kode). Status saat ini: **tahap desain, belum implementasi**.
> Terkait tapi berbeda topik dari breakdown tampilan stok (`docs/RENCANA_STOK_BREAKDOWN.md`, jika masih ada) — dokumen itu soal *cara menampilkan* stok yang benar, dokumen ini soal *cara menyimpan* stok supaya angkanya sendiri tidak salah dari akarnya.

> **Catatan rencana eksekusi**: implementasi di lokal akan dijalankan memakai **salinan data production** (bukan data dev/dummy biasa), supaya kondisi yang diuji saat pengembangan (termasuk skenario migrasi/backfill) sudah representatif dengan kondisi nyata sejak awal. Tujuannya: begitu siap rilis ke production, skenarionya sudah sama persis dengan yang sudah diuji di lokal — bukan baru ketahuan bedanya pas rilis.

## Masalah

Kolom `products.stock` bertipe `DECIMAL(15,3)` dan menyimpan stok dalam **satu satuan tunggal** (satuan anchor/`is_default` di `product_packages`, mis. Kardus/Slop/Krak).

Produk yang punya rasio konversi satuan tidak habis bagi rapi di basis 10 — mis. 24 Botol = 1 Kardus (faktor 1/24 = 0.041666...), 1/12, 1/3 — mengalami **truncation presisi** setiap kali dikonversi:

- BE menghitung faktor konversi dengan benar sebagai float presisi penuh (`model.ResolvePackageFactor`).
- Tapi begitu ditulis ke kolom `transaction_items.conversion_qty` (juga `DECIMAL(15,3)`), angka itu **terpotong** (mis. 0.041666... → 0.042).
- Nilai yang sudah terpotong inilah yang dipakai mengurangi `products.stock` — bukan nilai aslinya.
- **Satu kali transaksi saja sudah cukup** membuat stok "kotor" secara permanen — bukan soal akumulasi banyak transaksi.

### Bukti konkret (data nyata & direproduksi via test)

**Kasus nyata** — "Air Nestle Purelife" (product_id 104), rasio 24 Botol = 1 Kardus:

| Kejadian | `products.stock` |
|---|---|
| Beli 2 Kardus | 0.000 → 2.000 (bersih) |
| Jual 1 Botol (`conversion_qty` tersimpan 0.042, bukan 0.041666...) | 2.000 → **1.958** (kotor) |

Breakdown hasilnya: **"1 Kardus 22.992 Botol"** — padahal jawaban benar seharusnya bersih **"1 Kardus 23 Botol"** (2 Kardus − 1 Botol = 47/24 = tepat 1 + 23/24).

**Direproduksi ulang** dari produk baru (id 199, "TEST DRIFT Kardus Botol", rasio sama 1/24) lewat API asli (create produk → pembelian → penjualan) — hasil identik, mengonfirmasi ini bukan kebetulan data lama, tapi bug struktural yang akan selalu muncul untuk rasio non-terminating.

**Scope di data dev**: perkiraan awal 8 produk **sudah diverifikasi ulang (verifikasi ke-5) dan ternyata kurang tepat — sebenarnya 47 produk** dari ~140+ punya rasio rantai penuh yang tidak bulat. Rasio 1/12 sendiri sudah sangat umum (lusinan/karton), dan rantai berjenjang saling mengalikan memperbanyak yang kena (mis. 1/5 × 1/12 = 1/60, non-terminating meski 1/5 sendiri bulat). "Coffee Candy Kapal Api" (id 58) sudah menunjukkan drift nyata di riwayat `stock_mutations` (`stock = 9.680`, desimal ganjil untuk produk hitungan). Implikasi: scope perbaikan ini **jauh lebih luas** dari perkiraan awal — bukan cuma segelintir produk, tapi hampir sepertiga dari total produk yang punya lebih dari 1 satuan.

### Kenapa ini bukan cuma soal tampilan

`products.stock` dipakai untuk validasi penjualan, laporan, dan nilai persediaan (rupiah) — bukan cuma ditampilkan di Detail Produk. Kalau angka dasarnya menyimpang, semua yang dihitung dari situ ikut menyimpang.

## Solusi yang disepakati

### Prinsip dasar

Stok **tidak lagi disimpan sebagai satu angka desimal di satuan anchor**. Sebaliknya, **setiap level satuan produk punya angka stoknya sendiri**, tersimpan sebagai bilangan bulat (integer) — kecuali level paling bawah yang satuannya memang kontinu (kg/gram/liter), yang boleh tetap desimal karena itu wajar secara fisik.

Karena satuan transaksi (mis. jual 1 Botol) langsung match dengan satuan penyimpanan, **tidak pernah ada pembagian desimal** yang perlu dilakukan — sumber bug hilang total, bukan ditutupi.

### Perubahan skema

**Pindahkan kolom stok dari `products` ke `product_packages`** (bukan bikin tabel baru terpisah) — karena `product_packages` sudah berpola "banyak baris per produk, satu baris per level satuan", jadi tinggal ditambah kolom:

| Kolom | Lokasi | Sifat |
|---|---|---|
| `stock` | `product_packages` (baru) | **per-level** (satu angka per baris/satuan), integer kecuali level kontinu |
| `reserved_qty` | `product_packages` (baru) | **per-level**, sejajar dengan `stock` — retur supplier dibuat di satuan tertentu, jadi ditahan di level yang sama persis (lihat celah #14 — `supplier_return_items` perlu kolom `package_id` baru supaya ini bisa diimplementasikan) |
| `min_stock` | `products` (tetap) | 1 angka umum, diinput di satuan anchor — **tapi perbandingan "stok menipis" dilakukan di satuan terkecil** (lihat revisi di bawah), bukan langsung anchor vs anchor |
| `is_active` | `product_packages` (baru) | **ganti fungsi delete** — satuan tidak pernah dihapus fisik dari DB |

`product_packages.is_default` (sudah ada) tetap dipakai sebagai penanda satuan anchor — dipakai untuk cek stok menipis/`min_stock` dan sebagai referensi baris utama produk.

Catatan penting: **peran `product_packages` berubah** — awalnya cuma untuk data "grosiran"/satuan alternatif harga, sekarang jadi representasi **setiap level satuan produk**, termasuk anchor.

### Aturan operasional

1. **Setiap produk wajib punya minimal 1 baris `product_packages`** (`is_default=true`). Produk lama yang belum punya (satuan tunggal, belum pernah diisi "satuan alternatif") perlu dimigrasi dulu — dibuatkan baris anchor otomatis sebelum fitur ini aktif.

2. **Hapus satuan = nonaktifkan (`is_active=false`), bukan `DELETE` fisik.** Alasan: stok yang masih ada di baris itu tidak hilang/perlu dipindah paksa; riwayat transaksi lama (`purchase_items.package_id`, `transaction_items.unit_id`) yang merujuk baris itu tetap valid selamanya tanpa perlu jadi strict FK.

3. **Locking saat transaksi**: kunci **semua baris `product_packages` milik produk itu sekaligus** (`WHERE product_id = ? FOR UPDATE`), bukan cuma baris level yang sedang ditransaksikan — supaya breakdown antar level tetap konsisten kalau ada 2 transaksi jalan bersamaan pada produk yang sama.

4. **Jual/kurangi lintas-level** (mis. sisa lepasan Botol tidak cukup, perlu "buka" 1 Kardus): jangan kurangi langsung di kolom level yang dijual (bisa minus). Alurnya:
   - Untuk tiap level, hitung **stok yang boleh dipakai = `stock − reserved_qty`** level itu dulu (bukan `stock` mentah) — supaya bagian yang ditahan retur tidak ikut kepakai/kepinjam.
   - Turunkan semua level ke total **stok-boleh-dipakai** itu di satuan terkecil (integer).
   - Kurangi di situ. Kalau tidak cukup (termasuk karena sebagian stok sedang ditahan), **tolak transaksi** — jangan sampai "buka" baris yang stoknya sedang ditahan buat retur.
   - Breakdown ulang total baru itu naik ke semua level (`div`/`mod` berjenjang), lalu update semua baris sekaligus (`reserved_qty` tiap baris tidak ikut berubah, cuma `stock`-nya).
   - **Ditemukan di verifikasi ke-5**: tanpa langkah pengurangan `reserved_qty` di awal ini, algoritma cascade bisa "membobol" baris yang sedang ditahan (mis. Kardus stock=1 reserved=1, Botol stock=23 total 47 setara Botol — jual 30 Botol lolos cek `47≥30` padahal yang bebas cuma 23, dan Kardus yang ditahan retur ikut "dibuka"). Sudah diperbaiki dengan langkah pengurangan `reserved_qty` per-level di atas sebelum masuk ke pool.
   - **Prasyarat**: alur ini cuma valid kalau rantai `ref_package_id` produk itu benar-benar linear (1 induk → 1 anak di tiap level). Ditemukan saat verifikasi akhir: produk **1, 83, 112, 173** punya struktur **bercabang** (beberapa baris `product_packages` sama-sama menunjuk `ref_package_id` yang sama — mirip kasus "2 baris Gram ambigu" yang sudah diketahui sebelumnya). Untuk produk begini, fungsi terpusat **tidak boleh menebak otomatis** mau "pinjam" ke baris anak yang mana — **tolak diproses, tandai untuk ditinjau manual**. Kemungkinan besar ini data yang salah input dari awal, bukan struktur yang disengaja; baiknya dibetulkan strukturnya (jadi rantai linear) daripada dipaksa dapat aturan otomatis yang belum tentu benar.

5. **Satu fungsi terpusat untuk semua perubahan stok** — ini titik risiko tertinggi. Bug awal muncul karena tiap jalur (pembelian, penjualan, retur, void, sync offline, edit manual) punya logika sendiri-sendiri yang tidak konsisten. Jalur yang harus dialihkan ke fungsi terpusat ini:
   - Pembelian (create/update/void)
   - Penjualan (create/void)
   - **Jalur sync/offline** (`ApplySyncTransaction`, `ReturnStockForRejectSync`) — implementasi terpisah yang saat ini duplikat logikanya dari jalur online, risiko tinggi kalau tidak ikut diubah. **Dikonfirmasi verifikasi ke-14: ini fitur AKTIF dipakai** (terhubung route `/sync/push` & `/sync/conflicts/:id/resolve`, ada halaman FE `SyncCenterPage.tsx`) — bukan dead code, WAJIB ikut dipindah ke fungsi terpusat, bukan dilewati.
   - Retur ke supplier (reserve/release/reduce)
   - Write-off stok kadaluarsa
   - Edit produk manual — saat ini bisa override `stock` langsung **tanpa** tercatat di `stock_mutations` sama sekali (celah audit yang ditemukan saat investigasi)
   - Stok opname/adjustment — fitur ini **belum ada** di kodingan sekarang; kalau dibuat nanti, wajib lewat fungsi ini juga

6. **Data lama yang sudah kotor direkonstruksi dari riwayat, bukan dikoreksi manual dari awal.** (Direvisi — lihat detail lengkap di bagian "Migrasi data lama & skema DB" di bawah.) Ringkasnya: `quantity` asli di `purchase_items`/`transaction_items`/`supplier_return_items` tidak rusak — yang rusak cuma `conversion_qty`/`products.stock` hasil hitung yang kepotong. Jadi stok yang benar bisa dihitung ulang otomatis dari riwayat transaksi (pakai faktor konversi presisi penuh, bukan yang tersimpan kepotong), tanpa perlu opname manual untuk memperbaiki **kesalahan sistem**. Opname fisik tetap dijadwalkan terpisah, tapi fungsinya jadi verifikasi akhir (nangkep selisih fisik riil — barang hilang/rusak), bukan alat utama buat memperbaiki bug ini.

## Celah yang sudah diidentifikasi & keputusannya

| # | Celah | Keputusan |
|---|---|---|
| 1 | Produk 1 satuan mungkin belum punya baris `product_packages` | Wajib migrasi — setiap produk minimal 1 baris anchor |
| 2 | Locking lintas-baris saat transaksi barengan | Kunci semua baris produk sekaligus |
| 3 | Semua jalur tulis stok harus konsisten | Satu fungsi terpusat, semua jalur (termasuk sync) wajib lewat situ |
| 4 | Rasio satuan diedit admin belakangan | Hitung ulang breakdown saat rasio berubah |
| 5 | Jual lintas-level bisa jadi minus kalau naif | Turun ke total satuan terkecil dulu → kurangi → breakdown ulang naik |
| 6 | Hapus baris satuan yang masih ada stok | Nonaktifkan (`is_active=false`), tidak pernah `DELETE` fisik |
| 7 | `reserved_qty` di level mana | Per-level, sejajar `stock` |
| 8 | `min_stock` di level mana | Tetap 1 angka (input di satuan anchor) — **direvisi**: perbandingan "stok menipis" TIDAK dilakukan anchor-vs-anchor (lihat celah #12) |
| 12 | Cek anchor-vs-anchor untuk `min_stock` salah alarm — ditemukan di verifikasi ke-4: produk dengan sisa stok < 1 unit anchor penuh (mis. 0 Kardus + 20 Botol dari kapasitas 48 Botol, ~42% masih ada) akan selalu ke-flag "stok menipis" walau stoknya masih banyak dalam satuan kecil. Ini bukan kasus langka — kejadian rutin di setiap produk begitu sisa turun di bawah 1 unit anchor. | **Opsi A dipilih**: perbandingan dikonversi dulu ke satuan terkecil (total stok semua level, dan `min_stock` dikonversi ke satuan terkecil pakai faktor konversi anchor→terkecil) sebelum dibandingkan — bukan anchor vs anchor langsung. Alasan: Opsi B (terima false-positive) berisiko bikin admin mengabaikan alert karena sering salah, termasuk pas beneran menipis. |
| 13 | `reserved_qty` bisa "dibobol" lewat jual lintas-level — ditemukan di verifikasi ke-5: algoritma cascade Aturan #4 aslinya beroperasi di `stock` mentah, bukan `stock - reserved_qty`, jadi bisa "buka" baris yang stoknya sedang ditahan retur supplier kalau totalnya (termasuk yang ditahan) masih cukup secara matematis. | **Opsi A dipilih**: kurangi `reserved_qty` tiap level dulu sebelum masuk ke pool satuan terkecil yang boleh dijual (lihat revisi Aturan Operasional #4). Alasan: lebih presisi dari Opsi B (blokir total kalau ada reserved di baris manapun) yang terlalu longgar sekaligus terlalu ketat — bisa menolak transaksi valid padahal cuma sebagian kecil stok baris itu yang ditahan. |
| 9 | Satuan kontinu (kg/gram/liter) dipaksa integer? | Tidak — cuma level terbawah yang kontinu boleh desimal |
| 10 | `ref_package_id` bercabang (bukan rantai linear) — ditemukan di produk 1, 83, 112, 173 | Tolak diproses otomatis, tandai untuk ditinjau manual (kemungkinan data salah input) |
| 11 | Unique constraint "1 anchor per produk" — constraint biasa gagal di data yang sudah ada (banyak baris `is_default=0`) | Pakai generated column + unique index (`is_default_flag`), bukan constraint komposit biasa |
| 14 | `supplier_return_items` tidak punya kolom `package_id`/`unit_id` — cuma `product_id`, `quantity`, `unit` (teks bebas). Ditemukan di verifikasi ke-6: Aturan #7 ("retur ditahan di level yang sama persis") tidak bisa diimplementasikan langsung karena tidak ada data yang nunjuk ke baris `product_packages` yang mana. Satu-satunya jalan sekarang: join tidak langsung lewat `purchase_item_id → purchase_items.package_id`, rapuh kalau data pembelian aslinya berubah/dihapus. | **Opsi A dipilih**: tambah kolom `package_id` langsung di `supplier_return_items`, masuk sekalian ke migrasi 004 — lebih eksplisit & tidak bergantung join tidak langsung. |
| 14b | **Pola yang sama persis ketemu lagi di jalur write-off** (verifikasi ke-8): `product_expiry_batches` juga cuma punya `product_id`, `purchase_item_id`, `qty` — tidak ada `package_id`. Write-off tidak akan tahu level mana yang harus dikurangi, persis kasus celah #14. Efek samping: `ExpiryWarningModal.tsx` (FE) saat ini menampilkan placeholder harfiah kata **"unit"**, bukan nama satuan asli — bug kosmetik lama yang jadi kelihatan pentingnya begitu `package_id` ada buat lookup nama satuan yang benar. | Perbaikan sama seperti #14: tambah `package_id` di `product_expiry_batches`, masuk migrasi 004, isi lewat join `purchase_item_id → purchase_items.package_id`. Sekalian perbaiki placeholder "unit" di `ExpiryWarningModal.tsx` jadi nama satuan asli. |
| 15 | **Sumber join yang dipakai celah #14/#14b ternyata kosong total** — ditemukan di verifikasi ke-9: `purchase_items.package_id` **100% NULL** di data live (202/202 baris). FE sudah mengirim `PackageID`, skema sudah menyediakan kolomnya, tapi `createPurchaseItemQuery` (BE) tidak pernah mengisinya. Ini bukan "sebagian data lama belum lengkap" — rencana backfill #14/#14b akan gagal total kalau dijalankan sekarang karena joinnya selalu mentok NULL. **Diperdalam di verifikasi ke-10**: bahkan setelah BE diperbaiki supaya pakai `PackageID` dari FE, **~44% baris pembelian tetap akan NULL** — karena `PurchaseFormModal.tsx` (FE) cuma mengisi `itemSelectedPackageId` kalau produk punya **lebih dari 1** satuan (dropdown dipaksa muncul); untuk produk 1-satuan (87 dari 196 produk di data dev, 44%), FE sengaja tidak pernah set nilainya, jadi dikirim `undefined`. Ini bukan lagi cuma masalah BE tidak pakai data yang dikirim — FE sendiri memang tidak selalu mengirim. | Perbaiki 3 sisi: **(a)** ke depan (BE) — `createPurchaseItemQuery` diisi `package_id` dari `PackageID` yang dikirim FE, masuk Fase 4 poin 1. **(b)** ke depan (BE, tambahan) — **Opsi B dipilih** untuk kasus FE tidak kirim nilai: kalau `PackageID` yang diterima BE kosong/nil, BE cari sendiri package anchor (`is_default=true`) produk itu sebelum insert — jangan bergantung FE selalu benar mengirim data. Lebih robust daripada Opsi A (perbaiki FE saja), konsisten dengan prinsip validasi/fallback ada di server. **(c)** data lama (202 baris) — backfill manual dengan cocokkan `unit` (teks) + `conversion_qty` tersimpan ke `product_packages` produk itu, **wajib selesai sebelum Fase 3**. |
| 16 | `stock_mutations` (log audit stok) tidak punya kolom `package_id` — cuma `product_id`. Di desain per-level yang baru, log ini tidak bisa menelusuri riwayat "level satuan yang mana" berubah, cuma total produk. Belum pernah diputuskan eksplisit di Aturan #5. | **Ditambahkan**: kolom `package_id` (nullable) di `stock_mutations`, masuk migrasi 004, diisi otomatis oleh fungsi terpusat (yang sudah pasti tahu level mana yang berubah). Alasan: tanpa ini, audit trail kehilangan presisi per level — balik ke masalah yang sama dengan alasan awal seluruh perbaikan ini dilakukan. |
| 17 | Nonaktifkan **produk** (bukan satuan) — ditemukan di verifikasi ke-11: `toggleProductStatusQuery` cuma `UPDATE products SET is_active = NOT is_active`, tanpa cek stok sama sekali dan tanpa cascade ke `product_packages`. Produk yang masih ada stoknya bisa dinonaktifkan bebas, dan dokumen belum punya aturan eksplisit soal ini (beda dari celah #6 yang soal nonaktifkan satuan, ini soal nonaktifkan produk secara keseluruhan). **Diperdalam di verifikasi ke-12**: dampaknya lebih serius dari sekadar "hilang dari tampilan aktif" — dikonfirmasi langsung ke kode, `getLowStockQuery`, semua query nilai persediaan di `report_repo.go`, dan `lowStockCountQuery` di `business_summary_repo.go` semuanya filter `is_active=1`. Jadi produk nonaktif yang masih ada stok **juga hilang dari laporan nilai persediaan (rupiah) dan dashboard**, selamanya, tanpa ada cara melihatnya lagi — ini bukan cuma isu UX, tapi **integritas laporan keuangan** (nilai persediaan understated). | **Diizinkan tetap** (tidak diblokir) — ada alasan bisnis wajar (produk discontinued, sisa stok masih mau dihabiskan/diretur). `stock`/`reserved_qty` di semua baris `product_packages` produk itu TIDAK ikut diubah/dihapus saat nonaktifkan. **Mitigasi tambahan (bukan blokir)**: tambah 1 laporan/query audit baru "produk nonaktif yang masih ada stok" (join `products.is_active=0` + `product_packages.stock > 0`), supaya nilai persediaan yang "tersembunyi" ini tetap bisa ditelusuri kapan pun dibutuhkan, bukan hilang sama sekali dari radar. |
| 21 | **Void tidak mengembalikan angka yang persis sama — drift permanen** — ditemukan lewat testing browser menyeluruh: jual 3 Botol (rasio 1/24) dihitung saat jual pakai pecahan asli `3×1/24=0.125` (benar). Tapi **void**-nya menghitung ulang dari `transaction_items.conversion_qty` yang **sudah tersimpan kepotong** (`0.042`, bukan `0.041666...`), jadi yang dikembalikan `0.126` bukan `0.125` — stok akhir jadi `11.501`, harusnya `11.500`. Selisih +0.001 permanen tiap kali void produk rasio non-bulat. Beda dari celah #15/#18 — ini soal rumus jual vs void yang tidak konsisten, bukan soal `package_id`/atomicity. | Di desain baru, `ApplyStockDelta` untuk void **HARUS** selalu hitung ulang fresh dari `package_id` + `quantity` asli (integer) yang tersimpan di `transaction_items`, **BUKAN** membaca balik `conversion_qty` yang sudah kepotong. Ini otomatis teratasi kalau void dipaksa lewat fungsi terpusat yang sama seperti jual (arah delta dibalik, tapi sumber datanya sama-sama fresh) — masuk syarat wajib Fase 4 poin 2. |
| 22 | **Bug UI: rasio satuan baru bisa kesimpen terbalik** — ditemukan lewat testing browser: form "Tambah Paket" (edit produk, `GrosirRowForm`) saat diisi "1 Botol = 24 Kardus", label preview-nya kebalik, dan data yang tersimpan **rasionya terbalik** (`qty=1, ref_qty=24`) dibanding cara paket dibuat saat produk pertama kali dibuat (`qty=24, ref_qty=1`). Dampak nyata: jual 3 unit itu **salah ditolak** ("stok tidak cukup") karena sistem hitung kebutuhan `3×24=72` bukan `3×1/24=0.125`. Rusak diam-diam, baru ketahuan pas transaksi, bukan pas disimpan. | Perbaiki `GrosirRowForm` (FE) supaya konsisten dengan konvensi penyimpanan yang dipakai saat pembuatan produk pertama kali — satu cara mapping `qty`/`ref_qty`, bukan dua cara berbeda tergantung form mana yang dipakai. Tambahkan validasi BE: setelah simpan paket baru, resolve faktornya dan cross-check masuk akal (mis. bandingkan dengan harga per unit yang diinput, kalau rasionya kebalik biasanya harga per unit jadi tidak masuk akal) sebagai jaring pengaman tambahan. |
| 23 | **Hapus paket/satuan tidak ada konfirmasi khusus & tidak ada penjagaan server** — ditemukan lewat testing browser: paket yang masih ada stok & riwayat transaksinya bisa dihapus langsung cuma lewat dialog umum "Update Produk?", tidak ada peringatan spesifik satuan itu masih dipakai. Memperkuat urgensi celah #6 (yang sudah rencanakan `is_active` bukan `DELETE` fisik) — sekarang terbukti nyata lewat UI, bukan cuma risiko teoretis. | Tidak perlu keputusan baru — celah #6 sudah menjawab ini (`is_active=false`, tidak pernah `DELETE` fisik). Ditambahkan sebagai bukti konkret kenapa celah #6 penting, dan sebagai pengingat: FE juga perlu tombol "hapus" satuan diubah jadi "nonaktifkan" dengan pesan yang jelas, bukan cuma backend-nya yang diubah. |
| 24 | **Edit qty pembelian mengubah stok tanpa tercatat di `stock_mutations`** — ditemukan lewat testing browser: `purchase_repo.go` Update (nambah/kurangi qty item pembelian yang sudah ada) mengubah `products.stock` dengan benar, tapi **tidak menulis baris `stock_mutations`** sama sekali — beda dari create & void yang keduanya sudah benar. Sekelas dengan celah #1 (edit produk manual bypass audit), tapi di jalur berbeda yang belum pernah disebut eksplisit. | Tambahkan insert `stock_mutations` di jalur Update pembelian (`purchase_repo.go` `Update`), masuk cakupan Fase 4 poin 1 (Pembelian) — pastikan SEMUA perubahan `stock` (bukan cuma create/void) selalu lewat fungsi terpusat yang otomatis mencatat mutasi, tanpa terkecuali. |
| 20 | **Produk yang dilewati backfill Fase 3 tidak punya penanda permanen** — ditemukan di verifikasi ke-16: skrip backfill menandai produk gagal/dilewati (mis. celah #10, rantai gagal di-resolve) sebagai "perlu ditinjau manual", tapi cuma di **laporan hasil skrip saat itu dijalankan** — begitu laporan itu tidak dilihat lagi, tidak ada jejak di data produknya sendiri. Begitu Fase 5 (jalur baca) aktif, produk ini `product_packages.stock` default 0, **kelihatan sama persis seperti stok memang habis** — admin tidak bisa bedakan "stok habis beneran" vs "backfill dilewati, data belum benar". **Diperdalam verifikasi ke-17**: 3 lubang tambahan ketemu — (a) belum ada cara admin **mematikan** flag ini setelah ditinjau (cuma dijelaskan cara dipasang & ditampilkan); (b) kolom baru ini belum masuk kontrak API/DTO respons produk, jadi badge FE tidak akan punya data buat ditampilkan; (c) belum ada aturan eksplisit "blokir operasi stok kalau `needs_stock_review=true`" di fungsi terpusat — sejauh ini cuma "kebetulan aman" karena stok defaultnya 0, tapi produk yang ditandai karena percabangan (celah #10) bisa saja punya stok non-zero di sebagian level kalau ada perbaikan manual parsial sebelum ditinjau tuntas. | Tambahkan kolom `products.needs_stock_review` + `stock_review_note`, diisi otomatis skrip Fase 3. **(a)** Tambah checkbox/tombol "Tandai sudah ditinjau" di `ProductFormModal.tsx` (Fase 6) yang mematikan flag + kosongkan catatan. **(b)** Masukkan kedua kolom ke response DTO produk (Fase 5), sebut eksplisit di bagian "Dampak ke kode > BE — DTO". **(c)** Tambah aturan eksplisit di Fase 2: `ApplyStockDelta` **blokir semua operasi stok** (bukan cuma andalkan stok=0 kebetulan) kalau `needs_stock_review=true`. |
| 19 | **Klaim "jalur sync tidak dipakai" ternyata salah** — ditemukan di verifikasi ke-14: dokumen berulang kali (sejak awal) menyebut `ApplySyncTransaction`/`ReturnStockForRejectSync` "diabaikan, tidak dipakai" tanpa pernah benar-benar diverifikasi. Faktanya terhubung ke route HTTP aktif + ada halaman FE (`SyncCenterPage.tsx`) yang memakainya. **Diperdalam di verifikasi ke-15** (audit langsung ke jalur sync yang baru masuk scope): (a) bug konversi ternyata **dobel** — `ReturnStockForRejectSync` (pembalik retur sync) punya bug simetris yang sama persis (pakai qty mentah, bukan qty × conversion_qty), bukan cuma `ApplySyncTransaction`. (b) payload sync (`SyncTransactionItemPayload` di `dto_sync.go`) **sama sekali tidak punya field `package_id`** — bukan cuma NULL seperti celah #15, tapi konsepnya memang tidak ada di kontrak sync, indikasi fitur ini dibangun sebelum sistem multi-satuan ada. Kabar baik: atomik & locking-nya sudah benar (tidak kena pola celah #18), dan arsitekturnya memang cocok dipanggilkan ke fungsi terpusat (bukan skenario terdistribusi rumit). | **Dibatalkan, dimasukkan ke scope**: jalur sync WAJIB ikut dipindah ke fungsi terpusat di Fase 4 poin 6 — **KEDUA** fungsi (`ApplySyncTransaction` DAN `ReturnStockForRejectSync`), bukan cuma satu. Tambahkan field `package_id` baru ke `SyncTransactionItemPayload`/`dto_sync.go` supaya jalur ini bisa disamakan dengan jalur online. |
| 18 | **Beberapa jalur multi-item TIDAK atomik hari ini** — ditemukan di verifikasi ke-12, **diperluas & dikoreksi di verifikasi ke-13**: `transaction_repo.go` `Create()` (checkout kasir) menjalankan loop item pakai `Exec`/`Raw` biasa TANPA transaksi DB. Dokumen sempat salah sebut `Void` sebagai "sudah benar" untuk pembanding — **setelah dicek ulang baris per baris, `transaction_repo.go` `Void()` DAN `purchase_repo.go` `Void()` ternyata SAMA-SAMA kena bug ini** (bukan cuma `Create`). **KOREKSI KEDUA (saat eksekusi Fase 4 poin 2)**: klaim untuk `transaction_repo.go` `Create`/`Void` ternyata **false positive** — keduanya SUDAH dibungkus transaksi di **service layer** (`transaction_service.go`, `s.repo.GetDB().Transaction(...)` + `WithTx(tx)`), cuma tidak kelihatan kalau cuma baca file repo-nya sendiri terpisah dari service. Dibuktikan lewat test API nyata: transaksi 2-item, item ke-2 sengaja dibuat gagal (stok tidak cukup) → item ke-1 **tidak ikut ter-commit**, atomik terbukti. `purchase_repo.go` `Void()` juga ternyata sudah dibungkus service layer serupa — perbaikan yang sempat ditambahkan di Fase 4 poin 1 (wrap transaksi sendiri di repo) jadi nested transaction/savepoint yang tidak berbahaya tapi redundan. Kabar baik: `expiry_batch_repo.go` `WriteOff` tetap contoh pola locking `FOR UPDATE` yang valid untuk dicontoh di `ApplyStockDelta`. | Fase 4 poin 1 & 2: cukup ganti hitungan stok inline jadi panggil `ApplyStockDelta` (yang sudah punya `FOR UPDATE` sendiri) — TIDAK perlu tambah pembungkus transaksi baru di `transaction_repo.go` Create/Void karena sudah ada dari service layer. |

## Temuan audit tambahan (di luar bug utama)

Hasil audit lanjutan BE/FE/DB mencari gap lain yang belum terbahas. ~~Temuan terkait jalur sync/offline diabaikan (fitur itu tidak dipakai)~~ — **koreksi verifikasi ke-14 (celah #19): klaim ini salah, jalur sync ternyata aktif dipakai, sudah dimasukkan ke scope perbaikan.** Sisanya:

| # | Temuan | Rekomendasi |
|---|---|---|
| A | Write-off stok kadaluarsa (`expiry_batch_repo.go:55`) diam-diam clamp ke 0 (`GREATEST(stock - qty, 0)`) kalau qty write-off melebihi stok — `stock_mutations` tetap mencatat qty penuh yang diminta, bukan yang benar-benar terpotong, jadi `stock_before - quantity ≠ stock_after` (audit trail tidak balance). | Cek dulu `qty write-off ≤ stok tersedia` sebelum eksekusi. Kalau tidak cukup, **tolak dengan error eksplisit** (bukan clamp senyap) — kelebihan qty write-off itu sendiri sinyal ada masalah data lain yang perlu diperiksa. `stock_mutations` harus selalu mencatat qty yang benar-benar dikurangi. |
| B | Tidak ada penjamin di level DB bahwa "1 produk = 1 baris anchor (`is_default=true`)" — cuma dijamin di kode aplikasi (`insertAnchorPackageOnCreate`). Desain baru bergantung penuh pada asumsi ini. | Tambahkan **unique partial index** `UNIQUE (product_id) WHERE is_default = true` di `product_packages` (cek dukungan versi MySQL yang dipakai) sebagai backstop kalau ada bug aplikasi yang lolos. |
| C | Tidak ada `CHECK` constraint apa pun di skema (`001_init_schema.sql`) — stok/harga negatif murni dijaga di kode aplikasi. | Tambahkan `CHECK (stock >= 0)`, `CHECK (reserved_qty >= 0)`, `CHECK (min_stock >= 0)`, dst pada kolom stok baru — **dengan catatan**: MySQL baru benar-benar *enforce* `CHECK` sejak 8.0.16+, versi lama cuma parse tanpa berlaku. Perlu cek versi MySQL dulu; kalau tidak didukung, tetap andalkan validasi aplikasi + row locking sebagai baris pertahanan utama, `CHECK` cuma lapisan tambahan. |
| D | Sistem ini berasumsi single-location (tidak ada tabel cabang/gudang) — `products`, `stock_mutations`, `product_packages` semua asumsi satu pool stok global. | Tidak perlu diubah sekarang (di luar scope, menambah kompleksitas berisiko melebar). Cukup **dicatat sebagai asumsi eksplisit** — kalau nanti ada kebutuhan multi-cabang, skema `product_packages` hasil redesign ini perlu ditinjau ulang, jangan diasumsikan otomatis mendukung multi-lokasi. |

**KOREKSI PENTING (verifikasi ke-14)**: sepanjang dokumen ini sebelumnya, jalur sync/offline (`ApplySyncTransaction`, `ReturnStockForRejectSync`) berulang kali disebut "diabaikan, tidak dipakai" — **klaim ini SALAH**, tidak pernah benar-benar diverifikasi. Faktanya: kedua fungsi ini terhubung ke route HTTP aktif (`POST /sync/push`, `POST /sync/conflicts/:id/resolve`, terdaftar di `BE/routes/segment/sync_routes.go` & `protected_routes.go`, permission `operasional.sync` ada di seed data), dan **ada halaman FE yang aktif memakainya** (`FE/src/features/operational/sync/SyncCenterPage.tsx`). Ini bukan dead code.

**Keputusan direvisi**: jalur sync **dimasukkan ke scope perbaikan**, bukan dilewati. Alasan: bug di jalur ini sebenarnya lebih parah dari sekadar truncation — `ApplySyncTransaction` **sama sekali tidak menerapkan faktor konversi satuan** saat deduct/restore stok (bukan cuma presisi kurang, salah total kalau satuan yang dipakai bukan anchor). Kalau dibiarkan di luar scope, begitu jalur lain selesai dipindah ke fungsi terpusat, jalur sync yang aktif dipakai ini justru jadi satu-satunya titik yang masih pakai logika lama yang rusak.

## Dampak ke kode (inventarisasi awal, belum final)

Hasil audit baca kode BE & FE — daftar lengkap file/fungsi yang perlu berubah:

**BE — jalur tulis stok** (semua perlu switch ke fungsi terpusat):
- `BE/domain/transaction/repo/transaction_repo.go` — `Create` (jual), `Void`, `ApplySyncTransaction`, `ReturnStockForRejectSync`
- `BE/domain/supplier_purchase/repo/purchase_repo.go` — create/update/void pembelian
- `BE/domain/expiry_batch/repo/expiry_batch_repo.go` — write-off kadaluarsa
- `BE/domain/supplier_return/repo/supplier_return_repo.go` — reserve/release/reduce retur
- `BE/domain/product/repo/product_repo.go` + `product_service.go` — edit produk manual (saat ini bypass `stock_mutations` sama sekali)

**BE — jalur baca stok** (perlu join ke `product_packages WHERE is_default=true` atau agregasi):
- `product_repo.go` — list/get produk, `GetLowStock`
- `BE/domain/report/repo/report_repo.go` — laporan stok, nilai persediaan
- `BE/domain/business_summary/repo/business_summary_repo.go` — dashboard low-stock count

**BE — skema**:
- Tabel baru/kolom baru di `product_packages` (`stock`, `reserved_qty`, `is_active`)
- Perlu flag `is_continuous` (di `units` atau `product_packages`) untuk menandai satuan kontinu vs diskrit — belum ada field ini sekarang
- `ResolvePackageFactor` (`BE/domain/product/model/product_package_resolve.go`) perlu fungsi pendamping yang mengembalikan faktor integer relatif ke satuan terkecil (bukan float relatif ke anchor)

**BE — DTO** (dikonfirmasi via verifikasi lanjutan):
- `BE/domain/product/dto/dto_product.go` — **3 struct terpisah** (create/update/response) masing-masing punya field `Stock float64`, ketiganya perlu diubah, bukan cuma 1 tempat
- `BE/domain/product/dto/dto_product.go` (response struct saja, celah #20) — tambahkan `NeedsStockReview bool` + `StockReviewNote string` supaya badge peringatan di FE (Fase 6) punya data buat ditampilkan
- `BE/domain/business_summary/dto/dto_business_summary.go` — `LowStockCount int64` (agregat, bukan field stok flat) — kemungkinan aman as-is, cuma query sumbernya yang berubah (sudah tercakup di `business_summary_repo.go`)

**FE**:
- `FE/src/features/products/products/products.utils.ts` — `breakdownStock()`/`formatStockBreakdown()` saat ini menghitung ulang breakdown di client dari float `resolved_factor`; setelah BE menyimpan angka per-level asli, fungsi ini jadi tidak perlu lagi — breakdown harusnya datang langsung dari server, sudah benar per level, tidak perlu direkonstruksi dari desimal.
- `FE/src/features/products/products/components/ProductDetailModal.tsx` — konsumsi data per-level baru
- `FE/src/features/sales/cashier/components/ProductSearch.tsx` — validasi stok & tampilan sisa
- `FE/src/features/products/products/components/ProductFormModal.tsx` — **(ditemukan di verifikasi lanjutan)** baca/tulis `product.stock` langsung + ada field form "Stok" dengan validasinya sendiri, perlu disesuaikan ke model per-level
- `FE/src/features/products/products/components/ProductTableColumns.tsx` — **(ditemukan di verifikasi lanjutan)** baca `row.stock`/`row.min_stock` buat kolom tabel daftar produk, perlu disesuaikan
- `FE/src/features/reporting/stock/components/StockReportTableColumns.tsx` — dikonfirmasi pakai nama field beda (`current_stock`/`stock_value`, bukan `product.stock` langsung) — bentuknya beda dari dugaan awal, tetap perlu ditinjau ulang
- `FE/src/features/procurement/purchases/components/PurchaseFormModal.tsx` — perlu dicek ulang setelah bentuk data `product.stock` berubah
- `FE/src/features/products/products/components/ExpiryWarningModal.tsx` — dikonfirmasi **tidak** menyentuh field stok sama sekali, aman/tidak perlu diubah (dicantumkan sebelumnya secara defensif)

## Migrasi data lama & skema DB

Migrasi dipecah 2 bagian karena sifatnya beda — skema lewat migration runner otomatis (`BE/database/migrate.go`, baca `BE/database/migrations/*.sql` urut nomor file), backfill data lewat skrip terpisah karena butuh logika penelusuran rantai satuan yang tidak bisa jadi SQL murni.

### Bagian 1 — Migrasi skema (`BE/database/migrations/004_stock_per_package_level.sql`)

File berikutnya di urutan migrasi yang sudah ada (`001_init_schema.sql`, `002_seed_data.sql`, `003_fix_receivables_status_enum.sql`):

```sql
ALTER TABLE product_packages
  ADD COLUMN stock DECIMAL(15,3) NOT NULL DEFAULT 0,
  ADD COLUMN reserved_qty DECIMAL(15,3) NOT NULL DEFAULT 0,
  ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1;

-- backstop DB-level: cegah stok/reserved negatif (Temuan #C — sebelumnya
-- cuma tercatat di tabel keputusan, belum pernah dieksekusi ke SQL, sudah
-- dikoreksi di verifikasi ke-18). MySQL 8.4.7 (dikonfirmasi Fase 0) enforce
-- CHECK penuh sejak 8.0.16+, jadi aman dipakai di sini.
ALTER TABLE product_packages
  ADD CONSTRAINT chk_product_packages_stock_nonneg CHECK (stock >= 0),
  ADD CONSTRAINT chk_product_packages_reserved_nonneg CHECK (reserved_qty >= 0);

ALTER TABLE products
  ADD CONSTRAINT chk_products_min_stock_nonneg CHECK (min_stock >= 0);

-- backstop DB-level: jamin cuma 1 anchor per produk (celah #B)
-- CATATAN (diperbaiki setelah verifikasi akhir): unique constraint biasa
-- (product_id, is_default) TIDAK BISA dipakai — mayoritas produk sudah
-- punya 2+ baris is_default=0, constraint biasa akan mengira semua baris
-- is_default=0 itu saling bentrok dan ALTER TABLE langsung gagal di data
-- yang sudah ada. Wajib pakai generated column (MySQL abaikan NULL di
-- unique index, jadi cuma baris is_default=1 yang benar-benar dijaga unik):
ALTER TABLE product_packages
  ADD COLUMN is_default_flag INT GENERATED ALWAYS AS (IF(is_default = 1, product_id, NULL)) STORED;

ALTER TABLE product_packages
  ADD UNIQUE KEY uq_product_default (is_default_flag);

-- penanda satuan kontinu (kg/gram/liter) vs diskrit (celah #9)
ALTER TABLE units
  ADD COLUMN is_continuous TINYINT(1) NOT NULL DEFAULT 0;

-- penanda produk yang dilewati/gagal backfill Fase 3, butuh tinjauan manual (celah #20)
ALTER TABLE products
  ADD COLUMN needs_stock_review TINYINT(1) NOT NULL DEFAULT 0,
  ADD COLUMN stock_review_note TEXT NULL;

-- retur supplier perlu tahu level satuan yang mana buat reserve/release (celah #14)
ALTER TABLE supplier_return_items
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_supplier_return_items_package FOREIGN KEY (package_id) REFERENCES product_packages(id);
  -- backfill data lama: isi dari purchase_items.package_id lewat purchase_item_id
  -- kalau kolom itu ada; kalau tidak ada relasinya, biarkan NULL dan tangani
  -- sebagai kasus "level tidak diketahui" di Fase 3/4 (bukan diasumsikan anchor)

-- write-off kadaluarsa punya masalah sama persis (celah #14b)
ALTER TABLE product_expiry_batches
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_expiry_batches_package FOREIGN KEY (package_id) REFERENCES product_packages(id);
  -- backfill sama: isi dari purchase_items.package_id lewat purchase_item_id

-- audit trail per-level (celah #16)
ALTER TABLE stock_mutations
  ADD COLUMN package_id INT NULL,
  ADD CONSTRAINT fk_stock_mutations_package FOREIGN KEY (package_id) REFERENCES product_packages(id);
  -- diisi otomatis oleh fungsi terpusat (Fase 2) mulai berlaku, data lama biarkan NULL

-- CATATAN PENTING (celah #15): purchase_items.package_id sendiri 100% NULL
-- di data live sekarang (verifikasi ke-9) — join di atas (celah #14/#14b)
-- akan mentok NULL kalau kolom sumber ini belum diperbaiki dulu. Backfill-nya
-- BUKAN sekadar ALTER TABLE (kolomnya sudah ada dari skema awal), tapi perlu
-- skrip terpisah yang mencocokkan `unit` (teks) + `conversion_qty` tersimpan
-- ke `product_packages` yang sesuai — lihat Fase 3 & Fase 4 poin 1.
```

`products.stock` dan `products.reserved_qty` **tidak dihapus** di migrasi ini — tetap dipertahankan sebagai cache turunan sampai semua jalur baca dipindah, biar rollout bisa bertahap tanpa merusak fitur yang belum sempat diubah.

### Bagian 2 — Backfill data (skrip Go one-off, bukan migration runner)

Tidak bisa pakai `UPDATE` SQL polos karena perlu telusuri rantai `product_packages` tiap produk (sampai 5 level, berjenjang lewat `ref_package_id`) — pakai ulang logika yang sudah ada di `ResolvePackageFactor` supaya konsisten dengan cara sistem menghitung faktor konversi selama ini.

**Pendekatan (revisi): rekonstruksi dari riwayat, bukan konversi langsung dari `products.stock` yang sudah kepotong.**

Insight kunci: yang rusak cuma kolom hasil hitung (`transaction_items.conversion_qty`, `products.stock`) — data mentahnya (`purchase_items.quantity`, `transaction_items.quantity`, `supplier_return_items.quantity`, dst, semuanya dalam satuan asli/integer) **tidak rusak**. Jadi stok yang benar bisa dihitung ulang dari nol, bukan diperbaiki dari angka yang sudah salah.

1. Pastikan **setiap produk punya minimal 1 baris `product_packages` (`is_default=true`)** — yang belum punya, dibuatkan otomatis dari `products.unit_id`.
2. Untuk tiap produk, ambil seluruh riwayatnya **berurutan waktu** lewat `stock_mutations` (untuk urutan & jenis kejadian: in/out/adjustment/void/return/expired) dicocokkan ke tabel sumbernya (`purchase_items`, `transaction_items`, `supplier_return_items`, `expiry_batch`).
3. Untuk tiap kejadian, hitung ulang delta-nya pakai `quantity` **asli** (integer, satuan aslinya) × faktor konversi **presisi penuh** (dihitung fresh lewat `ResolvePackageFactor`, bukan ambil dari `conversion_qty` yang tersimpan kepotong).
4. Jumlahkan semua delta secara berurutan di **satuan terkecil** (integer/exact fraction) → hasilnya stok yang seharusnya, bebas dari error akumulasi truncation.
5. Bandingkan hasil rekonstruksi dengan `products.stock` yang tersimpan sekarang:
   - **Cocok/dekat** (dalam toleransi kecil) → pakai hasil rekonstruksi, breakdown ke `product_packages.stock` per level (cascading integer div/mod).
   - **Selisih signifikan yang tidak bisa dijelaskan cuma dari truncation** → tandai di laporan skrip sebagai kandidat "kemungkinan ada masalah data lain" (bukan otomatis dianggap opname fisik, tapi perlu ditinjau — bisa jadi ada mutasi yang tidak konsisten, chain rasio pernah diubah di tengah jalan, dll).
6. Diproses **per produk dalam transaksi terpisah** (bukan satu transaksi raksasa untuk semua ~140+ produk) — produk yang chain-nya gagal di-resolve (circular reference, anchor hilang) dicatat & di-skip untuk ditinjau manual, tidak menggagalkan seluruh proses.
7. Setelah backfill diverifikasi, jalur baca BE dipindah bertahap dari `products.stock` ke `product_packages.stock`.
8. `products.stock`/`reserved_qty` baru dihapus di **migrasi terpisah** nanti (bukan bagian dari 004), setelah semua jalur baca & tulis sudah pindah total dan terbukti stabil.
9. **Stok opname fisik tetap dijadwalkan terpisah setelah backfill** — bukan buat memperbaiki bug ini (itu sudah selesai lewat rekonstruksi otomatis di atas), tapi buat menangkap selisih **fisik** riil (barang hilang/rusak/salah hitung) yang tidak akan pernah terlihat dari riwayat sistem manapun.

## Fase Pengembangan

Dipecah supaya tiap fase bisa diverifikasi sendiri sebelum lanjut — bukan dikerjakan sekaligus. Tiap fase **wajib** lewat gerbang standar project sebelum dianggap selesai: type-check → lint → build (FE), `go build`/`go vet` (BE) → baru test browser.

**Fase 0 — Persiapan**
- Backup penuh DB dev (dan nanti production saat rilis) sebelum migrasi apa pun dijalankan.
- ~~Cek versi MySQL~~ — **sudah dikonfirmasi**: MySQL 8.4.7, mendukung `CHECK` constraint & generated column. Draft SQL Fase 1 (versi generated-column) **sudah diuji jalan bersih** di tabel salinan data real (323 baris), tidak perlu dicek ulang.
- Tentukan daftar satuan yang ditandai `is_continuous = true` — dasarnya nama satuan yang sudah ada di tabel `units`, ditinjau manual satu-satu (bukan ditebak dari pola nama). Data `units` sekarang (22 baris): Pieces, Pack, Kilogram, Kardus, Karton, Batang, Slop, Pres, Bungkus, Renteng, Botol, Tabung, Sak, Pouch, Sachet, Kaleng, Galon, Gelas, Ikat, Krak, Gram, Cup. Hasil cross-check awal: **Kilogram, Gram, Galon** jelas kontinu; **Gelas** ambigu (perlu dipastikan maksudnya — takar cair lepas atau kemasan tetap?); sisanya diskrit. **Catatan**: "Liter"/"ml" yang jadi contoh ilustrasi di dokumen ini **tidak ada** di data sekarang — keputusan final daftar `is_continuous` tetap harus dikonfirmasi user, bukan diasumsikan dari cross-check ini.

**Fase 1 — Migrasi skema**
- Jalankan `004_stock_per_package_level.sql` (kolom baru di `product_packages`, `units.is_continuous`, unique constraint anchor).
- Verifikasi migrasi jalan bersih di dev DB, tidak mengganggu data/fitur yang sudah ada (karena `products.stock` lama belum disentuh).

**Fase 2 — Fungsi terpusat stok (BE, inti logika)**
- Bangun fungsi `ApplyStockDelta` (atau nama serupa) di `domain/product` — locking semua baris produk, konversi ke satuan terkecil, validasi cukup/tidak, cascading breakdown balik ke semua level, update atomik.
- Unit test khusus kasus-kasus yang sudah dibahas (daftar lengkap, jangan cuma "kasus yang sudah dibahas" secara umum):
  1. Rasio non-terminating (1/24, 1/3) — pastikan tidak ada presisi hilang.
  2. Jual lintas-level (borrow dari level atas) menghasilkan breakdown yang benar.
  3. Penolakan kalau stok tidak cukup (total, bukan cuma 1 level).
  4. Locking konsisten saat 2 operasi bersamaan pada produk sama.
  5. **`reserved_qty` tidak boleh ikut kepakai/kepinjam** (celah #13) — kasus spesifik: baris punya `reserved_qty > 0`, transaksi yang seharusnya cukup kalau `reserved_qty` diabaikan tapi tidak cukup kalau dihitung, harus ditolak.
  6. **Produk dengan `ref_package_id` bercabang** (celah #10) ditolak diproses otomatis, bukan menebak jawaban.
  7. **Produk `needs_stock_review=true`** (celah #20) — semua operasi stok (jual/beli/retur/write-off) ditolak eksplisit, TIDAK cuma andalkan stok=0 yang kebetulan menolak. Ini eksplisit karena produk bercabang (celah #10) bisa saja punya stok non-zero di sebagian level kalau ada perbaikan manual parsial sebelum ditinjau tuntas.
  8. **Void selalu hitung ulang fresh dari `package_id`+`quantity` asli** (celah #21) — TIDAK boleh baca balik `conversion_qty` yang tersimpan (berpotensi sudah kepotong). Test: void penjualan rasio non-bulat (1/24), pastikan stok kembali PERSIS ke angka semula, tidak ada drift +0.001 dst.
- Belum dipakai jalur manapun di fase ini — murni bangun & uji fungsinya sendiri dulu.

**Fase 3 — Backfill data lama**
> **Wajib sebelum Fase 4/5** — kalau jalur baca/tulis dipindah duluan sementara `product_packages.stock` masih default 0 (belum diisi), semua produk akan tampil stok 0 di periode transisi. Backfill harus selesai & terverifikasi dulu sebelum ada jalur lain yang bergantung pada kolom ini.
- **Prasyarat #1 — perbaiki `purchase_items.package_id` DULU (celah #15)**: kolom ini 100% NULL di data live sekarang. Backfill Fase 3 butuh kolom ini terisi untuk resolve faktor konversi tiap kejadian pembelian (bukan cuma dipakai celah #14/#14b) — kalau belum diperbaiki, seluruh algoritma rekonstruksi di bawah ini tidak bisa jalan sama sekali, bukan cuma retur/write-off yang kena. Jalankan skrip pencocokan `unit` (teks) + `conversion_qty` tersimpan ke `product_packages` yang sesuai, verifikasi hasilnya ke user, baru lanjut ke langkah berikutnya.
- **Prasyarat #2 — ✅ SELESAI (data retur)**: sebelumnya `stock_mutations.reference_type` di data dev cuma ada `purchase`/`transaction`. Sudah dibuat data uji nyata lewat UI asli (produk id 200 "TEST MIGRASI Kardus Botol": beli 3 Kardus → retur 1 Kardus disetujui → jual 5 Botol via Kasir) — sekarang `reference_type` sudah ada `supplier_return` juga. Sekaligus mengonfirmasi ulang bug precision-nya direproduksi di data segar (jual 5 Botol → `conversion_qty` 0.208333... kepotong jadi 0.208, stok jadi 1.792 bukan 1.79166...). **Masih perlu**: data uji jalur write-off/expired (`expired`, `void_purchase`) — belum dibuat.
- Jalankan skrip rekonstruksi dari riwayat (`purchase_items`/`transaction_items`/`supplier_return_items`/`expiry_batch` + `stock_mutations`) sesuai algoritma yang sudah dirancang di atas.
- Produk dengan struktur `ref_package_id` bercabang (celah #10 — produk 1, 83, 112, 173 di data dev) **dilewati otomatis**, ditandai untuk ditinjau manual, tidak dipaksa direkonstruksi.
- **(celah #20)** Setiap produk yang dilewati/gagal (bercabang, chain tidak resolve, selisih signifikan) WAJIB di-set `products.needs_stock_review = true` + `stock_review_note` diisi alasannya — jangan cuma dicatat di laporan skrip yang hilang begitu sesi berakhir. Ini penanda permanen di data, bukan cuma output sekali jalan.
- Tinjau manual daftar produk yang gagal/selisih signifikan dari hasil laporan skrip.
- Disarankan jalan di jam sepi/maintenance window kalau nanti dieksekusi di production — bukan sambil toko masih transaksi aktif, supaya tidak ada transaksi baru yang nyelip di tengah proses rekonstruksi riwayat.

**Fase 4 — Migrasi jalur tulis satu per satu**
Pindahkan tiap jalur ke fungsi terpusat, **satu jalur, verifikasi, baru lanjut** (bukan sekaligus semua) — sekarang bertumpu di atas baseline `product_packages.stock` yang sudah benar dari Fase 3:
1. Pembelian (create/update/void) — sekalian perbaiki celah #15: `createPurchaseItemQuery` diisi `package_id` dari `PackageID` yang dikirim FE, **DENGAN fallback BE** — kalau `PackageID` kosong/nil (FE tidak selalu kirim, khususnya utuk produk 1-satuan, ~44% kasus), BE cari sendiri package anchor (`is_default=true`) produk itu sebelum insert. Jangan cuma andalkan FE selalu mengirim (data lama sudah dibereskan di Fase 3 prasyarat). **Sekalian perbaiki celah #18 untuk `purchase_repo.go` `Void()`** — saat ini TANPA transaksi DB pembungkus, bungkus dalam `r.db.Transaction(...)`. **Sekalian perbaiki celah #24**: jalur `Update` (edit qty item pembelian) mengubah stok tapi TIDAK menulis `stock_mutations` — pastikan lewat fungsi terpusat yang otomatis mencatat mutasi.
2. Penjualan (create/void) — **wajib sekalian perbaiki celah #18 untuk KEDUA fungsi**: `Create()` maupun `Void()` di `transaction_repo.go` sama-sama menjalankan loop item TANPA dibungkus 1 transaksi DB (dokumen sebelumnya salah kutip `Void` sebagai "sudah benar" — sudah dikoreksi). Kalau item ke-N gagal, item sebelumnya sudah permanen ter-commit (sukses sebagian diam-diam). Bungkus seluruh loop di kedua fungsi dalam `r.db.Transaction(...)` + `FOR UPDATE` — contoh pola yang benar sudah ada di `expiry_batch_repo.go` `WriteOff()`, tinggal dicontoh.
3. Retur ke supplier (reserve/release/reduce)
4. Write-off kadaluarsa (sekalian perbaiki gap #A — tolak eksplisit kalau qty melebihi stok, bukan clamp senyap)
5. Edit produk manual (sekalian tutup celah bypass `stock_mutations`)
6. **Jalur sync/offline** (`ApplySyncTransaction`, `ReturnStockForRejectSync`) — **direvisi verifikasi ke-14: TIDAK lagi dilewati**, ternyata fitur aktif dipakai (route `/sync/push`, `/sync/conflicts/:id/resolve`, halaman FE `SyncCenterPage.tsx`). Sekalian perbaiki bug lama: `ApplySyncTransaction` sama sekali tidak menerapkan faktor konversi satuan saat deduct/restore stok (lebih parah dari truncation biasa). Verifikasi lewat halaman Sync Center di browser, bukan cuma go build.

**Fase 5 — Migrasi jalur baca**
- `product_repo.go` (list/get, `GetLowStock`), `report_repo.go` (laporan stok, nilai persediaan), `business_summary_repo.go` (dashboard low-stock) — pindah dari baca `products.stock` langsung ke join/agregasi `product_packages`.
- **Perbaikan sekalian (celah #12)**: `GetLowStock` dan semua cek "stok menipis" lain (`business_summary_repo.go` low-stock count, dsb) diubah dari perbandingan anchor-vs-anchor jadi perbandingan di satuan terkecil (total stok semua level vs `min_stock` yang dikonversi ke satuan terkecil) — supaya tidak salah alarm untuk produk yang sisanya < 1 unit anchor penuh.
- Definisikan bentuk response API baru untuk data produk (array per-level, bukan angka tunggal) — **termasuk `needs_stock_review`/`stock_review_note`** (celah #20) — draft kontraknya sebelum sentuh FE di Fase 6.

**Fase 6 — FE**
- `products.types.ts` — bentuk tipe data baru sesuai kontrak API Fase 4.
- `products.utils.ts` — pensiunkan `breakdownStock()`/`formatStockBreakdown()` versi rekonstruksi client-side, ganti jadi formatter murni dari data server.
- `ProductDetailModal.tsx`, `ProductSearch.tsx` (kasir), `ProductFormModal.tsx`, `ProductTableColumns.tsx`, `StockReportTableColumns.tsx`, `PurchaseFormModal.tsx` — sesuaikan konsumsi data baru. (`ExpiryWarningModal.tsx` dikonfirmasi tidak perlu diubah — tidak menyentuh field stok.)
- **(celah #20)** Tambahkan badge peringatan di `ProductDetailModal.tsx` & `ProductTableColumns.tsx` untuk produk dengan `needs_stock_review = true` — supaya admin bisa bedakan "stok habis beneran" vs "backfill dilewati, perlu ditinjau", bukan diam-diam terlihat sama (stok 0). **Tambah juga checkbox/tombol "Tandai sudah ditinjau" di `ProductFormModal.tsx`** yang mematikan flag + kosongkan `stock_review_note` — tanpa ini flag permanen menyala walau sudah diperbaiki manual.
- **(celah #22)** Perbaiki `GrosirRowForm` (bagian dari `ProductFormModal.tsx`, form "Tambah Paket") — rasio yang diinput bisa kesimpen terbalik (`qty`/`ref_qty` tertukar dibanding konvensi saat produk pertama dibuat), rusak diam-diam sampai ketahuan pas transaksi. Satukan cara mapping-nya, jangan beda antara form tambah-satuan-baru dengan form buat-produk-baru.

**Fase 7 — Regresi & verifikasi menyeluruh**
- Skenario end-to-end lewat browser: pembelian → penjualan lintas satuan → retur → write-off → cek breakdown di semua halaman (produk, kasir, laporan) konsisten.
- Bandingkan angka sebelum/sesudah untuk beberapa produk yang sudah diketahui datanya (mis. Air Nestle Purelife, Coffee Candy Kapal Api) sebagai kasus regresi wajib.

**Fase 8 — Bersih-bersih (setelah stabil di production)**
- Hapus `products.stock`/`reserved_qty` (kolom lama) lewat migrasi terpisah, setelah dipastikan tidak ada jalur baca/tulis yang masih bergantung padanya.

## Prompt Eksekusi per Fase

Tiap fase adalah sesi kerja implementasi (bukan diskusi lagi — desainnya sudah selesai di dokumen ini). Salin blok yang sesuai ke sesi AI baru (atau lanjutkan di sesi yang sama). Jangan lompat ke fase berikutnya sebelum fase saat ini disepakati/diverifikasi selesai.

---

#### PROMPT — Fase 0 (Persiapan)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Kerjakan FASE 0 (Persiapan) dari bagian "Fase Pengembangan".

TUGAS:
1. Backup penuh DB dev (pos_retail_db) sebelum melakukan apa pun — simpan dump-nya, konfirmasi ke saya lokasi filenya.
2. Cek versi MySQL yang dipakai (SELECT VERSION()) — laporkan apakah mendukung CHECK constraint (8.0.16+) dan unique constraint kombinasi boolean seperti yang didraft di bagian "Migrasi data lama & skema DB" > Bagian 1. Kalau tidak didukung, usulkan penyesuaian draft SQL-nya.
3. Baca isi tabel `units` yang ada sekarang, buat daftar satuan mana yang seharusnya ditandai `is_continuous = true` (kg, gram, liter, ml, dst) — tinjau manual satu-satu, jangan ditebak dari pola nama. Tampilkan daftarnya ke saya untuk saya konfirmasi sebelum dipakai di fase berikutnya.
4. JANGAN jalankan migrasi skema atau ubah kode apa pun di fase ini — murni persiapan & pengumpulan info.
```

---

#### PROMPT — Fase 1 (Migrasi skema)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 0 (Persiapan) harus sudah selesai — pastikan backup DB sudah ada dan daftar `is_continuous` sudah dikonfirmasi user sebelum lanjut.

TUGAS — Kerjakan FASE 1 (Migrasi skema):
1. Buat file migrasi baru `BE/database/migrations/004_stock_per_package_level.sql`, lanjutan urutan migrasi yang sudah ada (001, 002, 003) — isi sesuai draft di bagian "Migrasi data lama & skema DB" > Bagian 1 dokumen ini, sesuaikan dengan hasil cek versi MySQL dari Fase 0.
2. Jalankan BE (migration runner otomatis) di dev, pastikan migrasi jalan bersih tanpa error.
3. Verifikasi: struktur tabel `product_packages` dan `units` sudah punya kolom baru, data lama (`products.stock`, dll) tidak berubah/tidak rusak, dan tidak ada fitur existing yang rusak (jalankan aplikasi, cek beberapa halaman produk sekilas).
4. JANGAN mulai bangun fungsi terpusat atau ubah kode domain lain di fase ini.
```

---

#### PROMPT — Fase 2 (Fungsi terpusat stok)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 1 (Migrasi skema) harus sudah selesai.

TUGAS — Kerjakan FASE 2 (Fungsi terpusat stok, BE):
1. Bangun fungsi terpusat (mis. `ApplyStockDelta`) di `BE/domain/product` sesuai aturan operasional di dokumen ini (poin 3, 4, 5 di bagian "Aturan operasional"): locking semua baris `product_packages` milik produk, konversi ke satuan terkecil, validasi cukup/tidak (**tolak kalau tidak cukup, TANPA kecuali — termasuk write-off, sesuai gap #A: tolak eksplisit, jangan clamp/izinkan sebagian**. Koreksi verifikasi ke-11: kalimat sebelumnya di prompt ini sempat kontradiksi dengan keputusan gap #A, sudah diluruskan), cascading breakdown balik ke semua level, update atomik dalam 1 DB transaction.
2. Tulis unit test untuk fungsi ini mencakup MINIMAL 8 skenario ini (lihat detail lengkap di dokumen bagian "Fase 2"): rasio non-terminating (1/24, 1/3), jual lintas-level yang perlu "buka" level atas, penolakan saat stok tidak cukup, locking konsisten saat 2 operasi bersamaan pada produk sama, `reserved_qty` tidak boleh ikut kepakai/kepinjam (celah #13), produk `ref_package_id` bercabang ditolak diproses otomatis (celah #10), produk `needs_stock_review=true` menolak SEMUA operasi stok eksplisit (celah #20), dan **void selalu hitung ulang fresh dari `package_id`+`quantity` asli, TIDAK baca balik `conversion_qty` tersimpan — pastikan stok kembali PERSIS ke angka semula tanpa drift (celah #21)**.
3. Fungsi ini BELUM dipakai jalur manapun (purchase/transaction/dst) di fase ini — murni dibangun & diuji sendiri dulu, isolated.
4. Jalankan go build, go vet, dan seluruh unit test sampai bersih sebelum lapor selesai.
```

---

#### PROMPT — Fase 3 (Backfill data lama)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 1 & 2 harus sudah selesai. PENTING: fase ini wajib selesai & terverifikasi SEBELUM Fase 4/5 — jangan pindah jalur baca/tulis manapun sebelum backfill ini beres, supaya tidak ada periode "stok tampil 0".

TUGAS — Kerjakan FASE 3 (Backfill data lama):
1. Bangun skrip Go one-off yang merekonstruksi stok tiap produk dari riwayat — ikuti algoritma persis di bagian "Migrasi data lama & skema DB" > Bagian 2 dokumen ini (bukan konversi langsung dari products.stock yang sudah kepotong, tapi hitung ulang dari quantity asli di purchase_items/transaction_items/supplier_return_items/expiry_batch, dicocokkan urutannya lewat stock_mutations, pakai faktor konversi presisi penuh).
2. Pastikan dulu tiap produk punya minimal 1 baris product_packages is_default=true (buatkan otomatis untuk yang belum punya).
3. Proses per produk dalam transaksi terpisah (bukan 1 transaksi raksasa) — catat & skip produk yang chain-nya gagal di-resolve, jangan gagalkan seluruh proses.
4. Untuk tiap produk, bandingkan hasil rekonstruksi vs products.stock sekarang — kalau selisihnya di luar wajar (bukan sekadar efek truncation), tandai di laporan sebagai "perlu ditinjau", jangan langsung dipaksa masuk.
5. Jalankan dulu di DB DEV, tunjukkan laporan hasilnya ke saya (jumlah produk berhasil, jumlah yang ditandai perlu ditinjau, contoh sebelum/sesudah untuk beberapa produk termasuk Air Nestle Purelife & Coffee Candy Kapal Api) sebelum saya putuskan lanjut atau tidak. JANGAN jalankan ke production di fase ini.
```

---

#### PROMPT — Fase 4 (Migrasi jalur tulis)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 1, 2, 3 harus sudah selesai dan backfill sudah terverifikasi oleh user.

TUGAS — Kerjakan FASE 4 (Migrasi jalur tulis), SATU JALUR DULU lalu berhenti untuk verifikasi sebelum lanjut ke jalur berikutnya (urutan sesuai dokumen):
1. Pembelian (create/update/void) — BE/domain/supplier_purchase/repo/purchase_repo.go. PENTING (celah #15): `createPurchaseItemQuery` saat ini TIDAK PERNAH mengisi `package_id` walau FE mengirim `PackageID` — tambahkan itu di INSERT-nya. TAPI FE (PurchaseFormModal.tsx) sendiri cuma mengisi nilai itu untuk produk dengan >1 satuan (dropdown) — untuk produk 1-satuan (~44% dari total produk) FE kirim kosong. Jadi BE WAJIB punya fallback: kalau `PackageID` yang diterima nil, cari sendiri package `is_default=true` produk itu sebelum insert. Ini prasyarat celah #14/#14b (retur & write-off) bisa di-backfill dengan benar. PENTING JUGA (celah #18): `Void()` di file yang sama TANPA transaksi DB pembungkus — bungkus dalam `r.db.Transaction(...)`. PENTING JUGA (celah #24): jalur `Update` (edit qty item pembelian yang sudah ada) mengubah `products.stock` dengan benar tapi TIDAK menulis `stock_mutations` sama sekali — pastikan lewat fungsi terpusat yang otomatis mencatat mutasi.
2. Penjualan (create/void) — BE/domain/transaction/repo/transaction_repo.go. PENTING (celah #18, DIKOREKSI): BUKAN cuma `Create()` yang tidak atomik — `Void()` di file yang sama JUGA sama-sama pakai `Exec`/`Raw` lepas-lepas tanpa transaksi DB (dokumen sebelumnya keliru menyebut `Void` sudah benar). Kalau item ke-N di keranjang/void gagal, item sebelumnya sudah permanen ter-commit. Bungkus SELURUH loop di KEDUA fungsi (`Create` dan `Void`) dalam `r.db.Transaction(func(tx *gorm.DB) error {...})` + `FOR UPDATE` pada baris yang dikunci — contoh pola yang sudah benar ada di `expiry_batch_repo.go` `WriteOff()`, pakai itu sebagai acuan, jangan cuma ganti hitungan stok inline jadi panggil `ApplyStockDelta` di dalam loop yang masih tidak atomik.
3. Retur ke supplier (reserve/release/reduce) — BE/domain/supplier_return/repo/supplier_return_repo.go. PENTING: saat insert baris baru ke `supplier_return_items` (`createReturnItemQuery`), isi kolom `package_id` (ditambahkan di migrasi 004) dari `purchase_items.package_id` lewat join `purchase_item_id` — FE tidak perlu diubah (`purchase_item_id` sudah dikirim), tapi kalau langkah ini terlewat, `package_id` akan tetap NULL untuk SEMUA retur baru setelah migrasi, bukan cuma data lama.
4. Write-off kadaluarsa — BE/domain/expiry_batch/repo/expiry_batch_repo.go (sekalian perbaiki gap #A: tolak eksplisit kalau qty write-off melebihi stok, jangan clamp senyap ke 0. PENTING, sama seperti retur di poin 3: isi kolom `package_id` (migrasi 004) pada `product_expiry_batches` lewat join `purchase_item_id → purchase_items.package_id` saat insert batch baru — celah #14b. Sekalian perbaiki placeholder "unit" di ExpiryWarningModal.tsx jadi nama satuan asli.)
5. Edit produk manual — BE/domain/product (sekalian tutup celah: pastikan lewat fungsi terpusat supaya tercatat di stock_mutations, tidak lagi bisa bypass)

6. Jalur sync/offline (ApplySyncTransaction, ReturnStockForRejectSync) — BE/domain/transaction/repo/transaction_repo.go. DIREVISI (celah #19): sebelumnya dianggap tidak dipakai, ternyata FITUR AKTIF (route /sync/push, /sync/conflicts/:id/resolve, halaman FE SyncCenterPage.tsx). WAJIB diperbaiki, bukan dilewati — perbaiki KEDUA fungsi (bukan cuma ApplySyncTransaction): ReturnStockForRejectSync punya bug simetris yang sama (qty mentah, bukan qty × conversion_qty). Tambahkan field `package_id` baru ke SyncTransactionItemPayload/dto_sync.go — payload sync saat ini sama sekali tidak punya konsep package_id (bukan cuma NULL, field-nya memang tidak ada). Atomik & locking jalur ini sudah benar (tidak perlu diperbaiki di sisi itu).

Untuk TIAP jalur: ganti logika inline lama supaya memanggil fungsi terpusat dari Fase 2, jalankan go build/go vet, verifikasi lewat browser (alur asli: beli/jual/retur/write-off/edit produk/sync beneran), baru lapor selesai untuk jalur itu sebelum lanjut ke jalur berikutnya.
```

---

#### PROMPT — Fase 5 (Migrasi jalur baca)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 1-4 harus sudah selesai.

TUGAS — Kerjakan FASE 5 (Migrasi jalur baca):
1. Pindahkan product_repo.go (list/get produk, GetLowStock), report_repo.go (laporan stok, nilai persediaan), business_summary_repo.go (dashboard low-stock) dari baca products.stock langsung ke join/agregasi product_packages.
2. Definisikan & dokumentasikan bentuk response API baru untuk data stok produk (array per-level satuan, bukan angka tunggal) — tunjukkan draftnya ke saya sebelum lanjut ke Fase 6 (FE), karena FE akan bergantung pada bentuk ini.
3. go build/go vet sampai bersih, verifikasi lewat API (curl/Postman) bentuk response-nya sesuai draft.
4. JANGAN ubah kode FE di fase ini — itu Fase 6.
```

---

#### PROMPT — Fase 6 (FE)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 5 harus sudah selesai, kontrak API baru untuk data stok sudah didefinisikan dan disepakati.

TUGAS — Kerjakan FASE 6 (FE):
1. Update FE/src/features/products/products/products.types.ts sesuai kontrak API baru dari Fase 5.
2. Di FE/src/features/products/products/products.utils.ts, pensiunkan breakdownStock()/formatStockBreakdown() versi lama (yang merekonstruksi breakdown dari float resolved_factor di client) — ganti jadi formatter murni yang menampilkan data per-level yang sudah benar dari server, tanpa hitung ulang.
3. Sesuaikan konsumsi data baru di: ProductDetailModal.tsx, ProductSearch.tsx (kasir), ProductFormModal.tsx, ProductTableColumns.tsx, StockReportTableColumns.tsx, PurchaseFormModal.tsx. (ExpiryWarningModal.tsx sudah dikonfirmasi tidak perlu diubah.) PENTING (celah #20): tambahkan badge peringatan di ProductDetailModal.tsx & ProductTableColumns.tsx untuk produk `needs_stock_review=true`, plus checkbox/tombol "Tandai sudah ditinjau" di ProductFormModal.tsx yang mematikan flag + kosongkan stock_review_note. PENTING JUGA (celah #22): perbaiki GrosirRowForm (bagian dari ProductFormModal.tsx, form "Tambah Paket") — rasio yang diinput bisa kesimpen terbalik (qty/ref_qty tertukar dibanding konvensi saat produk pertama dibuat), satukan cara mapping-nya.
4. Jalankan urutan wajib: type-check → lint → build, baru verifikasi lewat browser (buka tiap halaman yang disentuh, pastikan breakdown stok tampil benar, tidak ada desimal aneh di satuan diskrit).
5. JANGAN lanjut ke Fase 7 sebelum semua halaman yang disentuh sudah dicek visual di browser.
```

---

#### PROMPT — Fase 7 (Regresi & verifikasi menyeluruh)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase 1-6 harus sudah selesai — ini fase verifikasi akhir sebelum dianggap siap production.

TUGAS — Kerjakan FASE 7 (Regresi & verifikasi menyeluruh) lewat browser sungguhan:
1. Jalankan skenario end-to-end: pembelian produk multi-satuan → penjualan lintas satuan (termasuk satuan yang dulu bermasalah, rasio 1/24 & 1/3, mis. Air Nestle Purelife, Djarum Super Kretek 12) → retur ke supplier → write-off kadaluarsa. Screenshot tiap langkah.
2. Bandingkan tampilan stok yang sama di 4 tempat sekaligus (tabel produk, Detail Produk, kartu Kasir, Laporan Stok) untuk beberapa produk yang sama — pastikan semuanya konsisten satu sama lain.
3. Cek khusus kasus regresi wajib: Air Nestle Purelife dan Coffee Candy Kapal Api (data yang sudah diketahui rusak sebelum perbaikan) — pastikan sekarang breakdown-nya bulat/benar, bandingkan dengan angka sebelum perbaikan yang sudah dicatat di dokumen ini.
4. Cek console browser (devtools) sepanjang seluruh skenario, laporkan kalau ada error/warning baru yang muncul.
5. Susun laporan akhir: apa yang sudah diverifikasi benar, apa yang masih meragukan, rekomendasi lanjut atau tidak ke Fase 8.
```

---

#### PROMPT — Fase 8 (Bersih-bersih, setelah stabil di production)
```
Baca docs/RENCANA_PERBAIKAN_STOK_PRESISI.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase ini HANYA dikerjakan setelah user eksplisit konfirmasi sistem sudah stabil di production dengan skema baru (bukan otomatis lanjut dari Fase 7).

TUGAS — Kerjakan FASE 8 (Bersih-bersih):
1. Pastikan ulang tidak ada satu pun jalur baca/tulis di BE atau FE yang masih bergantung pada products.stock/reserved_qty (grep menyeluruh untuk pastikan).
2. Buat migrasi baru (005_...) untuk menghapus kolom products.stock dan products.reserved_qty.
3. Backup DB sebelum menjalankan migrasi ini (kolom yang dihapus tidak bisa dikembalikan tanpa restore).
4. Verifikasi aplikasi tetap berjalan normal setelah kolom dihapus.
```

## Status

Tahap desain sudah disepakati, sudah dipecah jadi fase pengembangan.

**Fase 0 (Persiapan) — ✅ SELESAI** (16 Agu 2026): backup DB (`backups/pos_retail_db_backup_20260816_120158.sql`), MySQL 8.4.7 dikonfirmasi mendukung CHECK constraint & generated column, daftar `is_continuous` dikonfirmasi user (Kilogram, Gram, Galon, Gelas = kontinu).

**Fase 1 (Migrasi skema) — ✅ SELESAI** (16 Agu 2026): `BE/database/migrations/004_stock_per_package_level.sql` berhasil dijalankan bersih. Semua kolom/constraint baru terverifikasi ada, data lama (`products.stock`) tidak rusak.

Dua bug ditemukan & diperbaiki saat eksekusi (di luar rencana desain, murni soal cara migrasinya dijalankan):
1. **`STORED` generated column gagal** (error 1215 "Cannot add foreign key constraint") — `product_packages` punya FK ke dirinya sendiri (`ref_package_id`), MySQL perlu rebuild tabel penuh untuk `STORED` yang konflik dengan self-referencing FK. **Diperbaiki: pakai `VIRTUAL`** (fungsinya sama untuk kebutuhan unique index `is_default_flag`, tidak perlu rebuild tabel).
2. **Migration runner project ini (`BE/database/migrate.go`) memecah SQL murni berdasarkan tanda `;` literal** (regex/string split sederhana, bukan SQL parser sungguhan) — komentar `--` yang mengandung `;` di tengah kalimat (mis. "...ada; kalau...") merusak pemotongan statement berikutnya, jadi fragmen komentar disangka SQL dan menyebabkan syntax error. **Pelajaran untuk migrasi berikutnya**: jangan taruh tanda titik koma di dalam teks komentar SQL sama sekali, dan hindari komentar panjang trailing setelah `;` sebuah statement — taruh komentar SEBELUM statement yang dijelaskan, bukan sesudahnya.

**Fase 2 (Fungsi terpusat stok) — ✅ SELESAI** (16 Agu 2026): `ComputeStockDelta` (logika murni, `BE/domain/product/model/stock_delta.go`) + `ApplyStockDelta` (pembungkus DB, `BE/domain/product/repo/stock_delta_repo.go`). 8 skenario wajib semua PASS (7 unit test murni + 1 integration test locking `FOR UPDATE` pakai DB dev nyata, 20 goroutine konkuren, hasil akhir tepat 0 — tidak ada lost update). Belum dipakai jalur manapun di fase ini (isolated, sesuai rencana).

**Fase 3 (Backfill data lama) — ✅ SELESAI** (16 Agu 2026): 184/195 produk berhasil direkonstruksi dari riwayat, 11 ditandai `needs_stock_review` (4 bercabang/celah #10, 3 "stok tidak cukup"/"paket tidak ketemu" saat replay riwayat, 4 selisih signifikan di luar wajar) — bukan dipaksa masuk. Air Nestle Purelife (104): 1.9580 → **1.9583** (matematis benar, bukan yang kepotong). Coffee Candy Kapal Api (58) & Djarum Super Kretek 12 (197): cocok persis, tidak ada drift.

Sekalian diselesaikan sebagai prasyarat wajib sebelum backfill bisa jalan:
- **Prasyarat #1 (celah #15)**: `purchase_items.package_id` diperbaiki — jalur tulis (`purchase_repo.go`, dengan fallback BE) + backfill 203 baris lama (skrip `BE/cmd/backfill_purchase_package_id`), sekarang 0 NULL.
- **Prasyarat #2**: data uji retur & write-off dibuat lewat API asli (bukan insert manual) — `stock_mutations.reference_type` sekarang lengkap ada `purchase`/`transaction`/`supplier_return`/`expiry_batch`.
- Skrip backfill utama: `BE/cmd/backfill_stock_per_level` — mereplay riwayat lewat `ComputeStockDelta` yang sama dipakai jalur online (bukan reimplementasi terpisah), idempoten (aman dijalankan ulang).

**Fase 4 (Migrasi jalur tulis) — poin 1 (Pembelian) ✅ SELESAI** (16 Agu 2026), poin 2-6 belum: `purchase_repo.go` (Create/Update/AddItems/Void) sudah pakai `ApplyStockDelta`. Sekalian diperbaiki:
- Celah #18: `Void()` dibungkus `r.db.Transaction(...)` (sebelumnya lepas-lepas tanpa transaksi).
- Celah #24: `Update()` sekarang menulis `stock_mutations` (sebelumnya tidak sama sekali) — **menyingkap bug baru**: handler `Update` ternyata tidak pernah mengisi `req.UserID` dari sesi login (beda dari `Create`/`AddItems`), laten karena dulu tidak pernah exercised. Sudah diperbaiki di `purchase_handler.go`.
- `ApplyStockDelta` diperluas sekalian sinkronkan `products.stock` sebagai cache turunan (bukan cuma `product_packages.stock`) — supaya jalur baca yang belum dipindah (Fase 5) tidak baca data basi selama masa transisi.
- Diverifikasi lewat API asli (bukan cuma unit test): pembelian, edit qty (net delta lintas-level benar), tambah item (breakdown cascading benar), void (kembali persis tanpa drift), void ganda (error bisnis rapi, bukan 500) — semua PASS.

**Tinjauan tambahan setelah selesai**: sempat salah diagnosis `Delete()` PO sebagai "tidak membalikkan stok" — ternyata ada guard di `purchase_service.go` yang mewajibkan PO di-void dulu sebelum boleh dihapus, dan `Void()` sudah membalikkan stok saat itu. Perbaikan yang sempat ditambahkan (Delete ikut membalikkan stok lagi) akan menyebabkan **pengurangan dobel** — langsung dibatalkan setelah diverifikasi lewat API (void 8→6, lalu delete tetap 6, tidak berkurang lagi). Edge case minor tersisa (item PO yang sudah kena write-off parsial sebelum di-edit) — direverse dengan qty kotor, bukan net-of-writeoff; dicatat sebagai catatan, tidak blocking, prioritas rendah.

**Fase 4 poin 2 (Penjualan) — ✅ SELESAI** (16 Agu 2026): `transaction_repo.go` `Create`/`Void` pakai `ApplyStockDelta`. **Koreksi**: celah #18 false positive untuk file ini (sudah atomik via service layer, lihat detail di celah #18). Diverifikasi API nyata: jual, void (kembali persis), tolak stok kurang (pesan ramah), **multi-item item ke-2 gagal → item ke-1 tidak ikut ter-commit** (atomik terbukti).

**Fase 4 poin 3 (Retur ke supplier) — ✅ SELESAI** (16 Agu 2026): `supplier_return_repo.go` — celah #14 selesai (`supplier_return_items.package_id` terisi dari `purchase_items.package_id`), reservasi pindah ke level package (`product_packages.reserved_qty`, bukan lagi `products.reserved_qty`). Urutan penting di approve: **lepas reservasi dulu, baru kurangi stok** lewat `ApplyStockDelta` — kalau dibalik, `ApplyStockDelta` akan salah tolak (menganggap qty yang mau dikonsumsi masih "ditahan"). Diverifikasi API nyata: create (reserve di package benar), approve (reduce+release benar, stock_mutations punya package_id), reject (cuma lepas reservasi, stok tidak berkurang) — semua PASS.

**Fase 4 poin 4 (Write-off kadaluarsa) — ✅ SELESAI** (16 Agu 2026): `expiry_batch_repo.go` `WriteOff()` pakai `ApplyStockDelta`. Gap #A selesai (tolak eksplisit qty melebihi stok, bukan clamp senyap `GREATEST(...,0)`). Celah #14b selesai (`product_expiry_batches.package_id` diisi otomatis lewat `purchase_repo.go`'s `insertExpiryBatches`, dari `resolvedPackageID` yang sudah dihitung di titik yang sama). Placeholder "unit" di FE (`ExpiryWarningModal.tsx`) diganti nama satuan asli — perlu tambah field `unit_name` di response BE (`WarningResponse`/`BatchHistoryResponse`, join lewat `package_id → product_packages → units`) dan type FE (`expiry-batches.types.ts`). Diverifikasi API nyata: write-off normal (stok berkurang benar, `unit_name` tampil "Pack"), tolak qty melebihi stok (pesan jelas, batch tetap `active`, stok sama sekali tidak tersentuh).

**Fase 4 poin 5 (Edit produk manual) — ✅ SELESAI** (16 Agu 2026): celah #1 selesai — `product_repo.go` `Update()` tidak lagi menimpa `products.stock` langsung. Selisih `req.Stock` vs stok anchor sekarang diterjemahkan jadi delta lewat `ApplyStockDelta` (mutation_type `adjustment`, reference_type `product_edit`). Sekalian ditambah `UserID` di `dto.UpdateRequest` (sebelumnya tidak ada, jadi mutasi tidak bisa diatribusikan ke siapa yang edit). Diverifikasi API nyata: edit stok 8→10 Pack tercatat lengkap di `stock_mutations` (`user_id` terisi benar), edit lain tanpa ubah stok tidak menghasilkan mutasi baru (tidak ada noise).

**Fase 4 poin 6 (Sync/offline) — ✅ SELESAI** (16 Agu 2026), **Fase 4 tuntas semua 6 poin**: `transaction_repo.go` `ApplySyncTransaction()` dan `ReturnStockForRejectSync()` pakai `ApplyStockDelta`. Field `unit_id` (=package_id) di `SyncTransactionItemPayload` ternyata sudah ada sejak awal — koreksi temuan verifikasi sebelumnya (celah #19) yang bilang field-nya sama sekali tidak ada. Bug yang diperbaiki: (1) pre-check manual lama membandingkan `(stock - reserved_qty)` di satuan anchor langsung dengan qty di satuan paket yang dipilih tanpa faktor konversi — salah besaran total, bukan cuma kurang presisi; dihapus, validasi cukup/tidak sekarang ditangani `ApplyStockDelta`. (2) `ReturnStockForRejectSync` punya bug simetris (qty mentah tanpa faktor konversi) plus `mutation_type='REJECT_SYNC'` yang bukan bagian enum `stock_mutations.mutation_type` — diganti `'void'`. Diverifikasi lewat API asli `/api/sync/push` (jual 1 Pieces dari produk yang anchornya Pack, faktor 1/12): `stock_mutations.quantity` tercatat presisi `0.083` (bukan `1.000` mentah) dengan `package_id` terisi benar (sebelumnya selalu NULL); lalu diverifikasi jalur reject via `/api/sync/conflicts/:id/resolve` (`action=reject`): stok balik presis ke nilai semula, `mutation_type='void'` tersimpan valid. **Catatan proses**: run pertama sempat memakai binary BE lama (belum di-rebuild) sehingga sempat menulis data salah (potong `products.stock` cache 1.000 penuh, `package_id` NULL) — ketahuan dari hasil tidak sesuai ekspektasi, dibersihkan (hapus transaksi test only, bukan data asli), BE di-rebuild+restart, lalu diulang dengan hasil benar seperti di atas.

**Verifikasi ulang Fase 4 (semua 6 poin) — ✅ SELESAI** (16 Agu 2026): dicek ulang independen (agent terpisah, bukan cuma percaya laporan sebelumnya) — semua 6 jalur tulis stok dikonfirmasi PASS: memanggil `ApplyStockDelta` (bukan raw SQL lama), dibungkus transaksi (baik lewat `r.db.Transaction` sendiri atau lewat service layer yang sudah membungkus), `package_id` terisi benar (bukan NULL) untuk poin 3/4, `go build`/`go vet` bersih. Satu temuan kosmetik: `product_repo.UpdateStock(id, delta)` (dead code, tidak ada pemanggil di manapun, masih pakai raw `UPDATE products SET stock = stock + ?`) — dihapus beserta const query dan deklarasi interface-nya, supaya tidak ada lagi jalur raw-SQL stock tersisa di codebase sama sekali (walau sudah tidak pernah dieksekusi).

**Fase 5 (Migrasi jalur baca) — ✅ SELESAI** (16 Agu 2026): `product_repo.go` (`GetAll`/`GetByID`/`GetByBarcode`/`Search`/`GetLowStock`), `report_repo.go` (laporan stok, nilai persediaan), `business_summary_repo.go` (dashboard low-stock count) semua pindah dari baca `products.stock` langsung ke agregasi `product_packages`.

**Cara kerja**: `ComputeStockDelta` (Fase 2) di-refactor — logika "cari anchor/leaf/faktor konversi tiap baris" yang tadinya nempel di dalamnya diekstrak jadi `analyzePackageChain` (dipakai bersama), lalu ditambah `ComputeStockSummary` (murni baca, tidak pernah menulis) yang menghasilkan total stok per produk (dikonversi ke satuan anchor untuk tampilan) DAN status `is_low_stock` — dibandingkan di **satuan terkecil** (celah #12), bukan anchor-vs-anchor, supaya produk dengan sisa < 1 unit anchor (mis. 0 Kardus + 20 Botol) tidak salah alarm "stok menipis". Repo baru `stock_read_repo.go` (`domain/product/repo`) expose `BuildStockSummaries(db, minStockByProduct)` sebagai fungsi lintas-domain (pola sama seperti `ApplyStockDelta`), dipakai `report_repo.go` & `business_summary_repo.go` tanpa duplikasi logika.

**Produk yang gagal dianalisis** (rantai bercabang/celah #10, dll) sengaja TIDAK dimasukkan hasil agregasi — fallback ke kolom cache lama `products.stock`/`reserved_qty` (masih ada & disinkronkan `ApplyStockDelta`), `is_low_stock` dibiarkan `false` (mengandalkan badge `needs_stock_review` di FE Fase 6, bukan alert yang datanya sendiri meragukan).

**Filter "stok menipis" di list produk** (`GetAll` dengan `low_stock=true`) tidak bisa lagi difilter di SQL (perbandingan sekarang di satuan terkecil, bukan kolom SQL biasa) — solusinya muat kandidat yang cocok filter lain tanpa paginasi, hitung `is_low_stock` di Go, filter+paginasi di memori. Diterima karena dataset kecil (~200 produk).

**Diverifikasi lewat API asli** (bukan cuma go build/vet): `/api/products/detail/197` (stock sekarang `10.416666666666666`, presisi penuh, bukan cache `10.417` yang dibulatkan), `/api/products/list` (termasuk `low_stock=true`), `/api/products/search`, `/api/products/197/packages/list` (breakdown per-level: 10 Pack + 2 Pieces + 1 Sachet, semua bulat), `/api/reports/stock/list`, `/api/reports/stock/summary`, `/api/reports/business-summary/stats` (`low_stock_count`) — semua PASS, angka konsisten dengan logika satuan-terkecil yang baru.

### Draft kontrak API baru (untuk Fase 6)

**`ProductResponse`** (`GET/POST /products/detail/:id`, `/products/list`, dll) — field baru ditambahkan, field lama TIDAK dihapus (supaya tidak breaking untuk FE yang belum diupdate):
```jsonc
{
  "id": 197,
  // ...field lama tidak berubah (barcode, name, category, harga, dst)...
  "stock": 10.416666666666666,   // TOTAL semua level, satuan anchor, presisi penuh (dulu cache dibulatkan) -- masih angka tunggal untuk kompatibilitas tampilan lama
  "reserved_qty": 0,
  "min_stock": 5,
  "unit_id": 2,
  "unit_name": "Pack",
  "is_low_stock": false,          // BARU -- dibandingkan di satuan terkecil (celah #12), pakai ini, JANGAN hitung ulang stock<=min_stock di FE
  "needs_stock_review": false,    // BARU -- celah #20, kalau true tampilkan badge "perlu ditinjau", stock/is_low_stock tidak bisa dipercaya penuh
  "stock_review_note": ""         // BARU -- alasan singkat dari backfill Fase 3, ditampilkan di badge/tooltip
}
```

**`PackageResponse`** (`POST /products/:id/packages/list`) — array per-level satuan, field baru `stock`/`reserved_qty`/`is_active` ditambahkan ke struktur yang sudah ada (endpoint ini SUDAH ada sejak awal, cuma belum expose kolom stok baru):
```jsonc
[
  { "id": 323, "unit_id": 2, "unit_name": "Pack", "is_default": true,  "resolved_factor": 1,        "stock": 10, "reserved_qty": 0, "is_active": true },
  { "id": 324, "unit_id": 1, "unit_name": "Pieces", "is_default": false, "resolved_factor": 0.0833..., "stock": 2,  "reserved_qty": 0, "is_active": true },
  { "id": 325, "unit_id": 15, "unit_name": "Sachet", "is_default": false, "resolved_factor": 0.25,     "stock": 1,  "reserved_qty": 0, "is_active": true }
]
```
Ini yang dipakai FE Fase 6 untuk breakdown "1 Kardus 23 Botol" dkk — **langsung dari `stock` tiap baris, TIDAK direkonstruksi dari `resolved_factor` di client** (`breakdownStock()`/`formatStockBreakdown()` lama di `products.utils.ts` dipensiunkan di Fase 6, bukan diperbaiki).

**`StockItem`** (laporan stok) & **`StockSummary`** (ringkasan stok) — bentuk/nama field TIDAK berubah, cuma sumber datanya (`current_stock`, `is_low_stock`, dll sekarang dari agregasi `product_packages`, bukan `products.stock` mentah).

**Fase 6 (FE) — ✅ SELESAI** (16 Agu 2026): `products.types.ts` dapat field baru (`ProductPackage.stock/reserved_qty/is_active`, `Product.is_low_stock/needs_stock_review/stock_review_note`). `products.utils.ts`: `breakdownStock()`/`formatStockBreakdown()` (rekonstruksi dari `resolved_factor` di client) dihapus total, diganti `formatPackageBreakdown()` (formatter murni, breakdown datang langsung dari `stock` tiap baris server) + `formatStockNumber()` (pembulatan tampilan 3 desimal, dipakai di semua tempat yang masih menampilkan satu angka total — server Fase 5 sekarang kirim presisi penuh, mis. `10.416666666666666`, bukan cache lama yang sudah dibulatkan).

Konsumsi data baru disesuaikan di: `ProductDetailModal.tsx` (breakdown per-level + badge `needs_stock_review`), `ProductSearch.tsx`/kasir (badge `is_low_stock` dari server, bukan `stock < min_stock` naif lagi, "Sisa..." pakai breakdown), `ProductTableColumns.tsx` (badge + kolom Stok pakai `is_low_stock`+`formatStockNumber`), `StockReportTableColumns.tsx` (pakai `is_low_stock` dari API yang sudah ada di DTO tapi belum kepakai — sebelumnya masih hitung ulang `current_stock < min_stock` sendiri, celah #12 versi FE), `PurchaseFormModal.tsx` (format angka di warning banner). `ProductFormModal.tsx`: badge peringatan + tombol "Tandai Sudah Ditinjau" (celah #20) — **butuh endpoint BE baru** `POST /products/mark-reviewed/:id` (tidak ada sebelumnya, ditambahkan sekalian: dto/repo/service/handler/route) supaya bisa mematikan `needs_stock_review`+kosongkan `stock_review_note` sebagai aksi terpisah dari Update biasa (tidak ke-clear tidak sengaja saat edit lain).

**Celah #22 (GrosirRowForm) — dicek, TIDAK ditemukan/sudah fixed sebelumnya**: kode saat ini (`qty`/`ref_qty` convention, label preview, contoh di form) sudah konsisten dengan cara BE menyimpan & cara produk pertama kali dibuat — dikonfirmasi baik dari baca kode maupun screenshot browser nyata (tabel paket produk 197: "12 Pieces = 1 Pack", "1 Sachet = 3 Pieces", cocok dengan data DB `qty=12,ref_qty=1` dan `qty=1,ref_qty=3`). Kemungkinan bug ini sudah diperbaiki di commit terpisah (`72a9bb5`, 27 Jul 2026) sebelum sesi perbaikan stok ini dimulai. Tidak ada perubahan kode untuk celah ini.

**Verifikasi visual browser** — dilakukan lewat Playwright (headless Chromium, `npx playwright` diinstall on-demand karena tidak ada MCP Playwright terpasang di sesi ini), auth di-inject langsung ke `localStorage` (`auth-session`, sama seperti yang dipakai `zustand persist`) pakai token admin yang sudah ada, bukan lewat form login manual. Semua halaman yang disentuh dicek dengan screenshot + console error check (0 error di semua):
- Daftar Produk (angka bersih, badge "Perlu Ditinjau" utk produk id 1 tampil benar)
- Detail Produk (breakdown "10 Pack 1 Sachet 2 Pieces" utk produk 197 — dari data server langsung, bukan hitung ulang; produk id 1 tampil banner kuning + catatan review)
- Edit Produk (banner + tombol "Tandai Sudah Ditinjau" — diklik sungguhan, toast sukses muncul, banner hilang, dicek ulang ke DB `needs_stock_review` benar jadi 0, lalu dikembalikan manual ke kondisi semula karena ini data hasil backfill Fase 3 yang sah, bukan data uji)
- Laporan Stok (badge "Stok Rendah" konsisten dengan `is_low_stock` API)
- Kasir/`ProductSearch` (kartu produk stok menipis: "Sisa 1 Slop — stok menipis", breakdown bukan angka mentah)
- Form "Tambah Paket" (`GrosirRowForm`) — tabel paket existing tampil benar, konfirmasi celah #22 di atas

**Bug ditemukan & diperbaiki selama verifikasi visual** (bukan cuma dari baca kode statis):
1. Field input "Stok" di form Edit Produk menampilkan angka mentah 15-digit (`10,416666666666666`) apa adanya di `<input type="number">` — jelek & berisiko submit ulang angka presisi-berlebih kalau user tidak sengaja menyentuhnya. Diperbaiki: `mapProductToForm` membulatkan `product.stock` ke 3 desimal sebelum diisi ke form.
2. Tombol "Tandai Sudah Ditinjau" sempat gagal dengan toast "Route not found" — BE process yang jalan masih binary lama dari sebelum endpoint `mark-reviewed` ditambahkan (bukan bug kode, cuma belum di-rebuild+restart). Setelah restart BE, berfungsi normal.

**Fase 7 (Regresi & verifikasi menyeluruh) — ✅ SELESAI** (16 Agu 2026): dilakukan lewat browser sungguhan (Playwright headless Chromium, sama seperti Fase 6), bukan cuma baca kode/API. Skenario end-to-end penuh untuk **Air Nestle Purelife (104)**, produk dengan rasio 1/24 yang jadi kasus regresi wajib:

1. **Pembelian**: 24 Botol (Rp2.000/Botol) via `/suppliers/purchases` — form asli, supplier dipilih lewat combobox, submit lewat dialog konfirmasi sungguhan. Hasil: `product_packages` Kardus 1.000→2.000 (24 Botol persis = 1 Kardus, tidak nyisa), `stock_mutations.quantity` tercatat `1.000` exact.
2. **Penjualan**: 3 Botol via Kasir (`/kasir`), bayar tunai Rp15.000, kembalian Rp4.500 — struk transaksi WEB-20260816-002 muncul. `stock_mutations.quantity = 0.125` (**3/24 persis**, bukan `0.042` yang kepotong seperti bug lama).
3. **Retur ke supplier**: 5 Botol dari PO yang sama via `/suppliers/returns`, status Pending → disetujui (`Setujui`) via detail modal. `reserved_qty` naik ke 5 saat create, balik ke 0 + stock berkurang saat approve — urutan release-then-reduce (Aturan Operasional #4/#13) terbukti benar di UI nyata. `stock_mutations.quantity = 0.208` (**5/24 persis**).
4. **Write-off kadaluarsa**: PO baru dengan tanggal expired diisi saat pembelian (10 Botol, expired 1 Agu 2026) → badge "Expired" muncul di tabel produk → modal `ExpiryWarningModal` menampilkan **"10 Botol"** (bukan lagi placeholder "unit", konfirmasi ulang fix Fase 4 poin 4) → "Musnahkan" → batch jadi `written_off`. `stock_mutations.quantity = 0.417` (**10/24**, presisi utuh dijaga sepanjang cascading breakdown meski kolom legacy `purchase_items.conversion_qty`/tampilan tetap dibulatkan 3 desimal untuk display).

Stok akhir produk 104 setelah 4 operasi: **2 Kardus + 15 Botol** (`2.625` di satuan anchor) — dihitung manual dari kombinasi delta di atas dan cocok persis dengan yang ditampilkan sistem, tidak ada drift sedikit pun.

**Konsistensi 4 tempat tampilan** (poin 2 tugas) — dicek bersamaan untuk produk 104 pada state akhir yang sama:
| Tempat | Tampilan |
|---|---|
| Tabel Produk (`/products`) | `2.625` |
| Detail Produk (modal) | `2 Kardus 15 Botol` |
| Kartu Kasir (`/kasir`) | `Sisa 2 Kardus 15 Botol — stok menipis` |
| Laporan Stok (`/reports/stock`) | `2.625`, badge "Stok Rendah", Nilai Stok **Rp 126.000** (=2.625 × Rp48.000 harga beli, dihitung benar) |

Semua konsisten satu sama lain — angka gross sama (`2.625`), breakdown yang menampilkan per-level sama (`2 Kardus 15 Botol`), status "stok rendah"/"stok menipis" sama-sama aktif (min_stock=5 Kardus, dikonversi ke satuan terkecil = 120 Botol, sisa 63 Botol < 120 → benar ditandai menipis).

**Kasus regresi wajib** (poin 3 tugas) — dicek breakdown-nya sekarang lewat Detail Produk (browser nyata, bukan cuma query DB):
- **Air Nestle Purelife (104)**: sebelum perbaikan breakdown salah "1 Kardus 22.992 Botol" (dicatat di bagian "Masalah" dokumen ini) — **sekarang bersih**, dan setelah 4 operasi test di atas berakhir di "2 Kardus 15 Botol" (integer semua, tidak ada sisa desimal ganjil).
- **Coffee Candy Kapal Api (58)**: sebelum perbaikan `stock = 9.680` (desimal ganjil untuk produk hitungan, dicatat di bagian "Masalah") — **sekarang "9 Bungkus 11 Sachet 1 Pieces"**, seluruhnya integer, dikonfirmasi lewat screenshot Detail Produk browser nyata.

**Console browser** (poin 4 tugas): dipantau di SETIAP langkah skenario (page + pageerror listener aktif sepanjang seluruh sesi Playwright, bukan cuma dicek sekali di akhir) — **0 error/warning baru** di semua 22 script/langkah yang dijalankan (pembelian, penjualan, retur, write-off, cross-check 4 tempat, kedua kasus regresi).

**Laporan akhir Fase 7**:
- **Sudah diverifikasi benar**: precision fix bekerja end-to-end lewat UI asli (bukan cuma API) untuk keempat jalur tulis (beli/jual/retur/write-off) pada kasus rasio non-terminating (1/24); breakdown per-level konsisten di 4 tempat tampilan berbeda; kedua kasus regresi wajib (Air Nestle Purelife, Coffee Candy Kapal Api) sudah bersih; tidak ada console error baru di sepanjang pengujian.
- **Masih meragukan / di luar cakupan pengujian ini**: (a) skenario ini cuma menguji 1 produk secara mendalam (104) plus baca-saja untuk produk lain — belum menguji ulang seluruh ~195 produk satu-satu (tidak praktis manual, tapi Fase 3 backfill + Fase 4 verifikasi API sebelumnya sudah mencakup itu secara terpisah); (b) jalur sync/offline (`SyncCenterPage`) tidak diuji ulang lewat UI di Fase 7 ini (sudah diverifikasi API asli di Fase 4 poin 6, tidak diulang di sini karena butuh device kedua/simulasi offline yang di luar cakupan verifikasi browser biasa); (c) 11 produk yang ditandai `needs_stock_review` dari Fase 3 (termasuk produk id 1 yang dipakai uji badge di Fase 6) sengaja TIDAK diikutkan skenario Fase 7 ini — datanya memang belum valid sampai ditinjau manual, di luar cakupan "regresi" (bukan sesuatu yang rusak oleh perubahan Fase 1-6, tapi kondisi awal yang sudah diketahui & didokumentasikan).
- **Rekomendasi**: **lanjut ke Fase 8** kapan pun user siap merilis ke production dan sudah memantau stabilitas beberapa waktu — tidak ditemukan regresi apa pun dari seluruh rangkaian Fase 1-7. Sebelum Fase 8 (hapus kolom `products.stock`/`reserved_qty`), disarankan sekali lagi: (1) tinjau manual 11 produk `needs_stock_review` (di luar cakupan otomatis manapun, perlu keputusan bisnis), (2) pastikan tidak ada proses lain (laporan pihak ketiga, export terjadwal, dll) yang masih baca kolom lama itu langsung dari DB di luar aplikasi ini.

**Fase 8 (Bersih-bersih) — ⚠️ DIKERJAKAN DI DB DEV SAJA, BUKAN FINAL/PRODUCTION-READY** (16 Agu 2026): user secara eksplisit BELUM mengonfirmasi rilis production stabil (syarat wajib prompt Fase 8) — sistem ini masih di tahap dev. Atas keputusan user, Fase 8 tetap disiapkan & diuji penuh di DB dev sekarang sebagai latihan/persiapan, TAPI status ini harus diulang lagi (backup baru, migrasi baru kalau perlu) saat production sungguhan sudah stabil nanti — bukan dianggap "Fase 8 sudah beres selamanya".

**1. Grep menyeluruh** — ditemukan 1 gap nyata yang belum pernah terdokumentasi sebelumnya: `product_repo.go` `Create()` (tambah produk baru) menulis stok awal (`req.Stock`) HANYA ke kolom lama `products.stock`, tidak pernah disinkronkan ke `product_packages.stock` (anchor package selalu mulai dari 0 lewat `insertAnchorPackageOnCreate`, tidak ada langkah lanjutan). Dibuktikan lewat data nyata: produk id 55 "Susu Dancow Sachet Putih" (`products.stock=1`, `product_packages.stock=0`, tidak ada `stock_mutations` sama sekali) — untungnya sudah ketahuan & ditandai `needs_stock_review=1` oleh Fase 3 backfill, tapi kalau tidak diperbaiki, PRODUK BARU yang dibuat lewat form "Tambah Produk" dengan stok awal setelah Fase 3 akan diam-diam kehilangan stok awalnya begitu Fase 8 menghapus `products.stock` (datanya tidak akan tercatat di mana pun sama sekali). **Diperbaiki**: `Create()` sekarang memanggil `ApplyStockDelta` untuk stok awal (kalau > 0), tercatat di `stock_mutations` (`reference_type='product_create'`), atribusi ke `UserID` (field baru di `CreateRequest`, diisi handler dari sesi login sama seperti `UpdateRequest`). Diverifikasi lewat API asli: produk baru id 200 dengan stok awal 7 → `product_packages.stock=7` benar, `stock_mutations` tercatat lengkap.

Selain itu, grep memastikan TIDAK ADA lagi query SQL yang membaca/menulis `products.stock`/`products.reserved_qty` (dibedakan dari `product_packages.stock`/`reserved_qty` yang memang sumber kebenaran, alias `pp.` bukan `p.`) — sisa hasil grep cuma komentar penjelas & referensi ke skrip backfill Fase 3 yang sudah tidak jalan lagi otomatis. Perubahan tambahan yang diperlukan supaya bisa drop kolom dengan aman:
- `stock_delta_repo.go`: `ApplyStockDelta` berhenti menulis `products.stock` (cache sync dihapus, sudah tidak perlu sejak Fase 5 semua jalur baca pindah ke `product_packages`).
- `product_repo.go`: `p.stock`/`p.reserved_qty` dihapus dari semua SELECT (`GetAll`/`GetByID`/`GetByBarcode`/`Search`/`GetLowStock`); sortir "stock" di `GetAll` (dulu `ORDER BY p.stock`) dipindah ke sortir in-memory (gabung dengan jalur `low_stock` yang sudah in-memory sejak Fase 5).
- `report_repo.go`: sortir "current_stock"/"stock_value" (dulu `ORDER BY p.stock`) juga dipindah ke in-memory dengan pola yang sama.
- `stock_read_repo.go`: `BuildStockSummaries` fallback untuk produk yang gagal dianalisis (rantai bercabang dkk, celah #10) TIDAK LAGI baca `products.stock` (kolom sudah hilang) — diganti baca langsung baris anchor `product_packages` yang sudah dimuat di query yang sama, tetap akurat untuk baris yang trusted itu.

**2. Migrasi `005_drop_legacy_product_stock_columns.sql`** dibuat — `ALTER TABLE products DROP COLUMN stock, DROP COLUMN reserved_qty`. Dikonfirmasi dulu lewat `SHOW CREATE TABLE products` tidak ada CHECK constraint atau FOREIGN KEY yang menempel di kedua kolom itu, jadi `DROP COLUMN` polos aman.

**3. Backup** — `mysqldump` penuh DB dev dijalankan SEBELUM migrasi: `backups/pos_retail_db_backup_before_fase8_20260816_152603.sql` (~14MB).

**4. Migrasi dijalankan** (migration runner otomatis, BE direstart) — `005_drop_legacy_product_stock_columns.sql` sukses, `SHOW COLUMNS FROM products` mengonfirmasi `stock`/`reserved_qty` sudah hilang, `min_stock`/`needs_stock_review`/`stock_review_note` tetap ada (benar, tidak ikut terhapus).

**Verifikasi aplikasi tetap normal setelah kolom dihapus** (lewat API asli & browser, bukan cuma go build/vet):
- `go build`/`go vet` bersih setelah semua perubahan.
- Semua endpoint baca: detail produk, list (termasuk sort by stock & filter low_stock), search, laporan stok (termasuk sort by current_stock), ringkasan stok, dashboard business-summary — semua PASS, angka identik dengan sebelum kolom dihapus.
- Jalur tulis: pembelian baru (PO-20260816-010, 1 Pack produk 197) → stok naik benar (10.417→11.417) → di-void → stok balik persis (11.417→10.417, tidak drift).
- Jalur tulis BARU yang diperbaiki: produk baru dengan stok awal 7 → `product_packages.stock=7` + `stock_mutations` tercatat benar (celah yang ditemukan di langkah 1).
- Guard `Delete()` (`exists.Stock > 0`) tetap berfungsi benar lewat data agregasi baru (dicoba hapus produk test yang masih ada stok → ditolak dengan pesan yang benar).
- Browser (Playwright): halaman Produk & Laporan Stok tampil identik dengan sebelum migrasi, 0 console error.
- **Temuan sampingan (bukan dari Fase 8, pre-existing)**: `Delete()` produk yang sudah punya riwayat `stock_mutations` gagal dengan FK error (`stock_mutations.package_id` tidak punya `ON DELETE CASCADE`) — bukan regresi dari Fase 8 (constraint ini sudah ada sejak migrasi 004), tapi baru kena saat mencoba hapus produk test yang sudah pernah dapat stok awal. Tidak diperbaiki di sesi ini (di luar cakupan Fase 8, dan `Delete()` memang dirancang untuk produk yang belum pernah dipakai — pesan errornya sendiri sudah menyarankan "Nonaktifkan produk sebagai gantinya" untuk kasus lain yang mirip).

**PENTING — batasan status ini**: seluruh Fase 8 di atas dijalankan di **DB dev** (yang isinya salinan data production, per catatan Fase 0), BUKAN di production sungguhan. Migrasi 005 sudah tersedia dan teruji di dev, tapi baru boleh dijalankan di production setelah: (1) production sudah rilis dengan skema Fase 1-7, (2) dipantau stabil beberapa waktu, (3) backup production diambil terpisah sesaat sebelum migrasi 005 dijalankan di sana (backup dev di atas TIDAK bisa dipakai untuk production). 11 produk `needs_stock_review` dari Fase 3 juga masih menunggu tinjauan manual — tidak terpengaruh/tidak diperbaiki oleh Fase 8 (memang di luar cakupannya).
