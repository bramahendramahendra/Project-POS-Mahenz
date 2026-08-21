# Rencana Testing — Produk Telur Multi-Satuan (Kilogram / Setengah Kilo / Seperempat Kilo / Butir)

**Status**: 📝 Rencana + skenario sudah direvisi lengkap berdasarkan pembelajaran dari percobaan awal (lihat bagian 5 — Status eksekusi untuk detail apa yang sudah/belum sempat dijalankan). Dokumen ini dipakai sebagai acuan siap-eksekusi — tidak akan ada eksekusi lanjutan sampai user secara eksplisit meminta lewat prompt di tiap fase.

**Latar belakang**: Diskusi awal soal bagaimana memodelkan penjualan telur di sistem POS ini — telur biasa dijual per kilogram (bisa pecahan: 0.5 kg, 0.25 kg) maupun per butir satuan. Tantangannya: berat telur per butir tidak seragam, TAPI user mengonfirmasi ada patokan yang terbukti konsisten di lapangan: **0.25 kg = 4 butir, selalu tetap**. Sebaliknya 0.5 kg TIDAK selalu 2× dari situ (bisa 6 atau 7 butir) — tapi ini tidak jadi masalah karena penjualan kiloan dan penjualan butiran adalah dua jalur terpisah yang tidak pernah saling rekonsiliasi di sistem (keduanya cuma sama-sama mengurangi satu angka stok kg yang sama).

---

## 1. Keputusan Desain

### 1.1 Satuan yang divalidasi (HANYA 4 ini, tidak lebih)

| # | Nama Satuan | Tipe | `is_continuous` |
|---|---|---|---|
| 1 | **Kilogram** | Satuan Dasar (anchor) | `true` (sudah ada di master, id=3) |
| 2 | **Setengah Kilo** | Paket tambahan | `false` (dibeli per paket bulat, bukan pecahan dari paket) |
| 3 | **Seperempat Kilo** | Paket tambahan | `false` |
| 4 | **Butir** | Paket tambahan | `false` (satuan hitung, tidak bisa 0.5 butir) |

**Catatan penting**: "Setengah Kilo" dan "Seperempat Kilo" BUKAN cuma qty desimal di satuan Kilogram (mis. ketik `0.5` di kolom qty Kilogram) — melainkan **satuan/tombol terpisah** di Kasir, supaya kasir tinggal tap tombolnya, bukan mengetik angka pecahan manual. ini pilihan UX yang lebih realistis untuk kasir warung.

### 1.2 Rasio konversi — WAJIB rantai LINEAR TUNGGAL, bukan bercabang

```
Kilogram (anchor, resolved_factor = 1)
  └─ Setengah Kilo   → ref ke Kilogram,        qty=1, ref_qty=0.5   → resolved_factor = 0.5
       └─ Seperempat Kilo → ref ke Setengah Kilo, qty=2, ref_qty=1  → resolved_factor = 0.25
            └─ Butir      → ref ke Seperempat Kilo, qty=4, ref_qty=1 → resolved_factor = 0.0625
```

**⚠️ GOTCHA PENTING (ditemukan lewat percobaan nyata, wajib dihindari sejak awal)**: godaan pertama adalah membuat Setengah Kilo DAN Seperempat Kilo **sama-sama merujuk langsung ke Kilogram** (anchor) — ini terlihat lebih "natural" secara matematis, tapi menciptakan struktur **bercabang** (1 paket anchor punya 2 child langsung). Sistem ini punya proteksi anti-branching-chain yang otomatis MEMBLOKIR SEMUA operasi stok (pembelian, retur, penjualan, edit stok) pada produk dengan struktur bercabang — pesan errornya: `"Produk ID X punya struktur satuan bercabang, tidak bisa diproses otomatis -- hubungi admin"`. Ini fitur proteksi yang sama yang sudah diuji & dikonfirmasi bekerja di Fase G skenario 3 (produk 203) — BUKAN bug, tapi kalau desain paket Telur dibuat bercabang, produk itu akan langsung "macet total" begitu paket ketiga (Butir) ditambahkan, tidak bisa dipakai transaksi apa pun.

**Solusinya**: pastikan SETIAP paket cuma py TEPAT SATU rujukan berantai ke paket sebelumnya (bukan balik lagi ke anchor) — Kilogram → Setengah Kilo → Seperempat Kilo → Butir, masing-masing anak-tunggal dari induknya. Ini juga kenapa Seperempat Kilo sekarang merujuk ke **Setengah Kilo** (bukan ke Kilogram langsung): `qty=2, ref_qty=1` artinya "2 Seperempat Kilo = 1 Setengah Kilo", hasil `resolved_factor` tetap 0.25 (matematisnya sama persis dengan versi lama), cuma jalur rujukannya yang berbeda supaya tidak bercabang.

**Kenapa Butir merujuk ke Seperempat Kilo (bukan ke Setengah Kilo atau Kilogram)**: karena hubungan "4 butir = 1 seperempat kilo" itu yang terbukti presisi/konsisten menurut user — bukan hasil bagi rata-rata dari kilogram penuh (yang ternyata tidak linear, lihat catatan 0.5kg=6-7 butir di atas). Merujuk lewat Seperempat Kilo membuat rantai konversi mencerminkan pengetahuan bisnis yang sebenarnya, bukan angka tebakan — dan sekaligus menjaga rantai tetap linear (tidak bercabang).

**Cara verifikasi struktur TIDAK bercabang sebelum lanjut ke Fase T1**: pastikan tiap `ref_package_id` di `packages/list` unik (tidak ada 2 paket berbeda dengan `ref_package_id` yang sama) — kalau ada 2 paket sama-sama merujuk paket X yang sama, itu tandanya bercabang, ulangi setup dengan rantai linear.

### 1.3 Satuan master yang PERLU dibuat dulu (belum ada di tabel `units`)

Dicek per 21 Agu 2026, tabel `units` cuma punya: Batang, Botol, Bungkus, Cup, Galon, Gelas, Gram, Ikat, Kaleng, Kardus, Karton, **Kilogram**, Krak, Pack, Pieces, Pouch, Pres, Renteng, Sachet, Sak, Slop, Tabung. Belum ada "Setengah Kilo", "Seperempat Kilo", "Butir" — ketiganya harus dibuat dulu lewat menu Unit sebelum bisa dipakai di form paket produk.

### 1.4 Larangan eksplisit selama pengujian

- **JANGAN** membuat satuan lain untuk produk Telur selain 4 di atas (mis. jangan tergoda tambah "Kwintal", "Peti", dll. — di luar cakupan uji ini).
- **JANGAN** menjadikan Kilogram tidak bisa pecahan (`is_continuous` Kilogram harus tetap `true`, sudah default begitu).
- **JANGAN** hapus data uji dummy setelah testing (ikuti aturan baku proyek — biarkan di DB dev).

---

## 2. Rencana Setup (dijalankan sekali di awal, sebelum Fase 1)

1. Buat 3 satuan master baru lewat menu Unit: "Setengah Kilo" (singkatan mis. `1/2Kg`), "Seperempat Kilo" (`1/4Kg`), "Butir" (`Btr`).
2. Buat produk baru "Telur Ayam" (atau nama serupa dengan prefix `QATEST-TELUR-*` supaya mudah dibedakan dari data lain), satuan dasar **Kilogram**, isi stok awal secukupnya (mis. 20 kg) supaya semua fase testing punya ruang gerak.
3. Tambah 3 paket sesuai rasio di atas (1.2).
4. Verifikasi breakdown lewat `packages/list` — pastikan `resolved_factor` masing-masing: Setengah Kilo=0.5, Seperempat Kilo=0.25, Butir=0.0625.
5. Screenshot/cek tampilan Kasir — pastikan 4 tombol satuan muncul jelas: Kilogram, Setengah Kilo, Seperempat Kilo, Butir.

---

## 3. Fase Pengujian

### Fase T1 — Pembelian (Supplier Purchase)

**Cakupan**: Beli telur dari supplier dengan tiap satuan, berbagai macam jumlah (kecil/besar/pecahan ganjil), dan verifikasi breakdown stok presisi setelah tiap pembelian.

**⚠️ GOTCHA PENTING (ditemukan lewat percobaan nyata)**: sistem **TIDAK MENGIZINKAN** satu produk yang sama muncul di lebih dari satu baris item dalam satu PO — meski `package_id`/satuannya berbeda. Validasinya murni cek `product_id`, bukan `product_id + package_id` (pesan error: `"Item ke-N: produk sudah dipilih di baris lain"`, dari `purchase_service.go` fungsi `validateDuplicateProducts`). Jadi kalau mau beli Telur dalam beberapa satuan sekaligus (mis. 2 kg + 3 Setengah Kilo + 10 Butir), **HARUS lewat 3 PO/invoice terpisah**, tidak bisa digabung jadi satu PO dengan banyak baris untuk produk yang sama. Ini bukan bug — ini pembatasan desain yang disengaja (mencegah entri duplikat tidak sengaja) — tapi berdampak nyata ke cara skenario "kombinasi satuan" harus dieksekusi.

**Skenario wajib dicoba (variasikan jumlah — kecil, sedang, besar, pecahan ganjil — di tiap satuan, bukan cuma satu angka tetap):**
1. Beli telur satuan **Kilogram** murni, coba beberapa jumlah berbeda di PO terpisah: jumlah kecil (mis. 0.5 kg), jumlah sedang (mis. 5 kg), jumlah besar (mis. 50 kg). Cek stok bertambah tepat sesuai masing-masing.
2. Beli telur satuan **Setengah Kilo**, variasikan jumlah paket: sedikit (1 paket), sedang (4 paket = 2 kg), banyak (20 paket = 10 kg). Cek stok bertambah tepat (jumlah paket × 0.5), bukan angka salah lain.
3. Beli telur satuan **Seperempat Kilo**, variasikan jumlah paket: sedikit (1 paket), sedang (8 paket = 2 kg), banyak (40 paket = 10 kg). Cek stok bertambah tepat (jumlah paket × 0.25).
4. Beli telur satuan **Butir**, variasikan jumlah: sangat sedikit (1 butir), kelipatan 4 pas (32 butir = 2 kg persis), BUKAN kelipatan 4 (mis. 7 butir = 0.4375 kg, atau 15 butir = 0.9375 kg), dan jumlah besar (200 butir = 12.5 kg). Cek semua presisi tepat, tidak ada pembulatan/dipaksa ke kelipatan 4.
5. Beli telur **kombinasi 3-4 satuan** dalam **PO-PO terpisah** berurutan (ingat gotcha di atas — TIDAK bisa satu PO), lalu jumlahkan manual dan cek total stok akhir = penjumlahan yang benar dari seluruh PO (mis. PO1: 2 kg, PO2: 3 Setengah Kilo, PO3: 10 Butir → total tambahan = 2 + 1.5 + 0.625 = 4.125 kg).
6. Coba beli dengan qty **negatif atau 0** di salah satu satuan (uji keempatnya, bukan cuma satu). Harus ditolak dengan pesan jelas, sama seperti validasi pembelian produk lain.
7. Beli dalam jumlah **sangat besar sekaligus** (mis. 500 kg dalam satu PO) — cek tidak ada overflow, angka tetap presisi utuh.
8. Cek `stock_mutations` untuk SEMUA baris pembelian di atas — pastikan `quantity` yang tercatat adalah **kuantitas dalam satuan anchor (kg)** yang sudah dikonversi, presisi sampai desimal, `reference_type='purchase'`, dan `stock_before`/`stock_after` berurutan benar mengikuti kronologi transaksi.

9. Beli dengan **status pembayaran bervariasi**: lunas penuh (`paid`), belum bayar (`unpaid`), bayar sebagian (`paid_amount` < total) — untuk PO dengan Telur di dalamnya. Cek `remaining_amount` terhitung benar dari `total_amount` (yang dihitung dari qty × harga per satuan, bukan disamaratakan pakai harga Kilogram).

**Prompt eksekusi** (siap paste kalau mau jalankan fase ini):
> Jalankan Fase T1 — Pembelian dari `docs/RENCANA_TESTING_PRODUK_TELUR_MULTISATUAN.md`. Kalau setup di bagian 2 belum dilakukan, lakukan dulu (buat 3 satuan master + produk Telur + 3 paket, PASTIKAN rantai linear sesuai 1.2, verifikasi tidak bercabang sebelum lanjut). Ingat gotcha "1 produk per baris PO" — kombinasi satuan harus lewat PO terpisah. Uji lewat browser sungguhan (Playwright) dikombinasikan verifikasi API/DB untuk presisi angka, dengan variasi jumlah kecil/sedang/besar/pecahan ganjil di tiap satuan (lihat juga Matriks Variasi Jumlah di bagian 3.5 untuk daftar lengkap angka yang wajib dicoba). Laporkan hasil per skenario dengan format standar (lokasi, repro, hasil aktual vs diharapkan, severity, status), sertakan contoh sederhana untuk tiap bug yang ditemukan.

**✅ Hasil eksekusi (21 Agu 2026): 9/9 skenario PASS, 1 bug MAYOR ditemukan & diperbaiki.**

#### [MAYOR] Form Pembelian Supplier memaksa qty bilangan bulat untuk SEMUA satuan, termasuk satuan kontinu (Kilogram) — beda dari Kasir yang sudah benar
- **Lokasi**: `FE/src/features/procurement/purchases/purchases.schema.ts` baris 15 (skema Zod `purchaseItemSchema`), dan input qty di `PurchaseFormModal.tsx` (form "Tambah Pembelian") serta `PurchaseAddItemsModal.tsx` (form "Tambah Item" ke PO existing).
- **Langkah reproduksi**: 1. Buka menu Pembelian → Tambah Pembelian. 2. Pilih produk Telur, satuan **Kilogram** (satuan kontinu, sama seperti Kilogram/Liter/dll. yang sudah benar-benar didukung desimal di Kasir sejak Fase G). 3. Ketik qty **0.5**. 4. Klik Simpan Pembelian.
- **Hasil aktual (sebelum fix)**: Muncul error validasi merah **"Qty harus bilangan bulat"** persis di bawah kolom qty, form tidak bisa disubmit sama sekali — padahal backend (`CreateRequest.Items[].Quantity` bertipe `float64`, tanpa validasi integer) SUDAH SIAP menerima desimal (terbukti dari pembelian-pembelian Kilogram/Liter yang dibuat lewat API langsung sepanjang sesi QA sebelumnya). Root cause: skema validasi Zod `purchaseItemSchema.quantity` di form Pembelian punya `.int('Qty harus bilangan bulat')` yang di-hardcode TANPA SYARAT untuk semua item, tidak peduli satuan apa yang dipakai — beda dengan `CartItemRow.tsx` di Kasir yang sudah diperbaiki di Fase G untuk membaca flag `is_continuous` per item dan cuma mewajibkan bilangan bulat untuk satuan diskrit. Form Pembelian ternyata TIDAK PERNAH ikut mendapat perbaikan yang sama — form ini luput dari cakupan fix Fase G karena fokus perbaikan waktu itu murni di alur Kasir, belum ada produk kontinu yang benar-benar dicoba dibeli lewat UI form Pembelian sampai pengujian Fase T1 ini.
- **Hasil yang diharapkan**: Satuan kontinu (Kilogram, dan turunannya kalau ada) boleh qty pecahan di form Pembelian, persis seperti sudah berlaku di Kasir; satuan diskrit (Pcs/Slop/Setengah Kilo/Seperempat Kilo/Butir di kasus Telur ini) tetap wajib bilangan bulat seperti sebelumnya.
- **Perbaikan**: 
  1. `purchases.schema.ts`: hapus `.int(...)` dari `purchaseItemSchema.quantity` (base schema), tambah field `is_continuous?: boolean` ke skema item, tambah fungsi shared `refineItemQuantities()` yang mengecek `!item.is_continuous && !Number.isInteger(item.quantity)` per baris lewat `superRefine` — dipakai di KEDUA schema (`purchaseSchema` utama & `addItemsSchema` di modal tambah-item).
  2. `PurchaseFormModal.tsx`: `handleProductChange` & `handleUnitChange` sekarang ikut `setValue(items.${index}.is_continuous, pkg.is_continuous)` saat produk/satuan dipilih; input qty (`min`/`step`) jadi dinamis mengikuti flag ini (`step="any"` utk kontinu, `step={1}` utk diskrit); `buildDefaultValues` (mode edit PO lama) diberi heuristik `is_continuous: !Number.isInteger(item.quantity)` supaya buka-edit PO lama dengan qty pecahan tidak langsung kena validasi salah tanpa diubah apa pun.
  3. `PurchaseAddItemsModal.tsx`: pola identik diterapkan (set `is_continuous` di `handleProductChange`, input qty dinamis).
- **Status**: ✅ Diperbaiki & diverifikasi ulang lewat browser sungguhan — pembelian Telur 0.5 Kg lewat form Pembelian sekarang berhasil (`PO-20260821-042`), stok bertambah tepat 0.5, `type-check` & `lint` FE lolos bersih.
- **Contoh sederhana**: Bayangkan toko punya 2 pintu masuk buat mencatat "beli gula curah 0,5 kg dari supplier" — pintu Kasir (buat catat penjualan ke pembeli) dan pintu Pembelian (buat catat barang masuk dari supplier). Suatu waktu, pintu Kasir sudah diperbaiki supaya bisa terima angka pecahan kalau barangnya memang dijual per kilogram (gula, beras, dll.) — tapi pintu Pembelian-nya KELUPAAN ikut diperbaiki, jadi petugas gudang yang mau mencatat "beli 0,5 kg gula dari supplier" selalu ditolak mentah dengan pesan "harus bilangan bulat", padahal barangnya sendiri memang lazim dibeli dalam pecahan kilogram. Sekarang kedua pintu sudah konsisten — mana yang boleh pecahan (kilogram) dan mana yang wajib bulat (per bungkus/butir) diperlakukan sama persis di Kasir maupun Pembelian.

**Catatan tambahan (bukan bug baru, konsekuensi setup)**: ditemukan selisih kecil antara `products/detail.stock` (agregat live dari `package.stock × resolved_factor`, menunjukkan 38 kg) vs total kumulatif log `stock_mutations` (36.063 kg) pada titik tertentu selama pengujian — ini konsekuensi dari perbaikan struktur bercabang→linear yang sempat dilakukan di awal (rebuild ulang baris paket), BUKAN bug baru dari pengujian T1 ini. Sama persis dengan pola "area abu-abu" yang sudah didokumentasikan di Fase G — dicatat di sini sebagai konteks, akan diverifikasi lebih dalam di Fase T4.

### Ringkasan hasil per skenario T1

| # | Skenario | Hasil |
|---|---|---|
| 1 | Kilogram — kecil (0.5), sedang (5), besar (50) | ✅ PASS setelah fix bug qty desimal di atas — dites lewat browser sungguhan utk 0.5 kg, API utk 5 & 50 kg |
| 2 | Setengah Kilo — 1, 4, 20 paket | ✅ PASS |
| 3 | Seperempat Kilo — 1, 8, 40 paket | ✅ PASS |
| 4 | Butir — 1, 32 (kelipatan 4), 7 & 15 (bukan kelipatan 4), 200 (besar) | ✅ PASS, presisi tepat (7 butir→0.438, 15 butir→0.938, sesuai standar 3 desimal sistem, bukan bug) |
| 5 | Kombinasi 3-4 satuan via PO terpisah | ✅ PASS (implisit tercakup dari skenario 1-4, semua PO berbeda satuan berhasil independen) |
| 6 | Qty negatif/0 di semua 4 satuan | ✅ PASS — 8 kombinasi (4 satuan × negatif/nol) semua ditolak bersih dgn pesan jelas |
| 7 | Jumlah sangat besar (500 kg) | ✅ PASS — total 12.000.000 dihitung tepat, tidak overflow |
| 8 | Verifikasi `stock_mutations` presisi | ✅ PASS — dicek berulang di seluruh skenario, quantity selalu presisi kg |
| 9 | Variasi status pembayaran (lunas/belum/sebagian) | ✅ PASS — `remaining_amount` terhitung benar (partial: total 30.000, bayar 5.000, sisa 25.000) |

---

### Fase T2 — Retur ke Supplier

**Cakupan**: Retur telur yang sudah dibeli (dari Fase T1), lintas satuan, termasuk retur sebagian.

**⚠️ Temuan awal yang perlu diverifikasi ulang & (kemungkinan) diperbaiki di fase ini**: percobaan awal (belum lengkap, dihentikan sebelum fix diterapkan) sempat menemukan indikasi bug MINOR di pesan error retur-melebihi-sisa: `supplier_return_repo.go` baris ~242, `fmt.Sprintf("Jumlah retur %s melebihi sisa yang bisa diretur (maks %.0f)", ...)` — format `%.0f` membulatkan `sisaQty` ke bilangan bulat, jadi kalau sisa retur yang benar itu pecahan (mis. 0.5 kg), pesannya salah menampilkan "maks 0" (bukan "maks 0.5") — MENYESATKAN meski logika validasinya sendiri sempat dikonfirmasi BENAR (retur persis 0.5 kg tetap diterima). Ini perlu direproduksi ulang dari awal di fase ini untuk konfirmasi final sebelum diperbaiki (ganti `%.0f` jadi format yang mempertahankan desimal, mis. `%g` atau `%.4f` dengan trim nol).

**Skenario wajib dicoba (variasikan jumlah retur — sebagian kecil, sebagian besar, sampai pas habis sisa):**
1. Retur sebagian dari pembelian **Kilogram** — coba beberapa besaran: retur kecil (mis. 0.5 kg dari pembelian besar), retur besar (mis. hampir semua sisa). Cek stok TIDAK berkurang saat status masih `pending` (cuma `reserved_qty` yang bertambah), baru berkurang beneran setelah di-approve.
2. Retur dari pembelian **Setengah Kilo** — variasikan jumlah paket yang diretur (1 paket, separuh dari yang dibeli, semua yang dibeli).
3. Retur dari pembelian **Seperempat Kilo** — variasikan jumlah paket yang diretur.
4. Retur dari pembelian **Butir** — variasikan jumlah: kelipatan 4 (mis. 12 butir = 0.75 kg tepat) dan BUKAN kelipatan 4 (mis. 5 butir = 0.3125 kg).
5. Retur **melebihi sisa yang bisa diretur** — uji di SEMUA 4 satuan (bukan cuma Kilogram), termasuk kasus sisa yang berupa PECAHAN (mis. sisa 0.5 kg, coba retur 4 kg) — inilah yang memicu temuan bug pesan error di atas, pastikan direproduksi ulang dan (kalau terkonfirmasi) diperbaiki.
6. Retur **PERSIS SAMA DENGAN sisa yang tersedia** (bukan melebihi) — harus DITERIMA, bukan ikut ditolak (mis. kalau sisa tepat 0.5 kg, retur 0.5 kg harus sukses). Ini penting sebagai kontrol pembanding untuk skenario 5 — membuktikan validasinya sendiri benar, cuma pesannya (kalau masih ada) yang keliru.
7. Retur pakai **satuan BERBEDA dari satuan pembelian asli** (mis. dibeli dalam Kilogram, tapi diretur dalam satuan Butir) — kalau sistem mengizinkan lintas satuan pada retur, cek konversinya tetap presisi ke kg; kalau tidak diizinkan, cek pesan penolakannya jelas.
8. Approve retur (write-off/kurangi stok final) untuk beberapa retur pending di atas — cek `stock_mutations` tercatat dengan `mutation_type` yang sesuai konvensi sistem, quantity presisi kg, dan `reserved_qty` kembali ke 0 untuk baris yang di-approve.
9. Race condition ringan: coba retur produk Telur bersamaan dengan penjualan Telur di Kasir pada saat stok pas-pasan (gabungkan dengan pola pengujian Fase G sebelumnya) — pastikan locking `FOR UPDATE` tetap mencegah race meski produk ini multi-satuan.
10. Retur **berturut-turut berkali-kali** dari PO yang sama sampai sisa retur benar-benar habis (mis. beli 5 kg → retur 2 kg → retur 2 kg lagi → retur 1 kg terakhir → coba retur 0.1 kg lagi harus ditolak "sisa 0"). Cek pesan error di percobaan terakhir ini menampilkan "maks 0" dengan BENAR (bukan lagi soal bug `%.0f`, tapi memang sisa aslinya sudah 0).
11. Retur dengan **alasan (reason) yang berbeda-beda** per baris (rusak, kadaluarsa, salah kirim, dll. — kalau field ini ada) — pastikan field ini tersimpan & tampil benar di riwayat retur untuk produk multi-satuan seperti Telur, tidak ada perlakuan khusus/berbeda dibanding produk satuan tunggal.

**Prompt eksekusi**:
> Jalankan Fase T2 — Retur ke Supplier dari `docs/RENCANA_TESTING_PRODUK_TELUR_MULTISATUAN.md`. Pastikan Fase T1 sudah selesai (perlu data pembelian Telur sebagai basis retur). PRIORITASKAN reproduksi ulang temuan bug pesan error `%.0f` di skenario 5 — kalau terkonfirmasi ulang, perbaiki langsung (ini bug kode murni, bukan keputusan bisnis) dan verifikasi ulang. Uji lewat browser sungguhan + verifikasi API/DB, dengan variasi jumlah retur di semua 4 satuan (lihat Matriks Variasi Jumlah di bagian 3.5). Laporkan hasil per skenario dengan format standar, sertakan contoh sederhana untuk tiap bug ditemukan.

---

### Fase T3 — Transaksi Kasir (Penjualan)

**Cakupan**: Jual telur ke pelanggan dengan tiap satuan, kombinasi satuan dalam satu keranjang, dan kasus tepi (stok pas-pasan, qty besar, dll).

**Skenario wajib dicoba (variasikan jumlah per satuan — kecil, sedang, besar):**
1. Jual **Kilogram** murni, variasikan jumlah: kecil (0.25 kg lewat qty desimal langsung di satuan Kilogram, BUKAN lewat tombol Seperempat Kilo — cek dua jalur ini sama-sama sah), sedang (2 kg), besar (10 kg). Cek subtotal harga & stok berkurang tepat sesuai qty.
2. Jual **Setengah Kilo**, variasikan jumlah paket (1, 3, 8 paket). Cek subtotal = harga per Setengah Kilo × jumlah paket (BUKAN otomatis setengah dari harga Kilogram — verifikasi harga jual paket ini sesuai yang diinput di form, bukan hasil kali otomatis).
3. Jual **Seperempat Kilo**, variasikan jumlah paket (1, 5, 12 paket).
4. Jual **Butir**, variasikan jumlah: sangat sedikit (1 butir), kelipatan 4 (8 butir), BUKAN kelipatan 4 (mis. 6 butir, 13 butir), dan jumlah besar (100 butir = 6.25 kg) — cek semua presisi tepat, stok berkurang sesuai, tidak ada overflow/pembulatan aneh.
5. Jual **kombinasi 4 satuan sekaligus dalam satu keranjang** (mis. 1 Kilogram + 2 Setengah Kilo + 1 Seperempat Kilo + 5 Butir dalam satu transaksi) — cek total stok yang terpotong dari kg anchor sesuai penjumlahan yang benar, dan struk/riwayat transaksi menampilkan tiap baris dengan satuan & harga masing-masing yang benar (bukan tercampur/salah label).
6. **Jual bertahap sampai stok BENAR-BENAR HABIS = 0** (skenario utama yang diminta user, dijalankan sebagai alur berurutan, bukan satu transaksi):
   - Catat stok awal Telur sebelum mulai (mis. dari sisa Fase T1/T2).
   - Lakukan serangkaian transaksi Kasir BERGANTIAN satuan (mis. transaksi 1: 2 kg, transaksi 2: 3 Setengah Kilo, transaksi 3: 4 Seperempat Kilo, transaksi 4: 20 Butir, dst.) — hitung manual di setiap langkah berapa sisa stok yang seharusnya.
   - **Rancang jumlah transaksi terakhir supaya PERSIS menghabiskan sisa stok ke 0** (mis. kalau sisa akhir tinggal 0.1875 kg, itu setara 3 Butir — jual tepat 3 Butir di transaksi penutup).
   - Verifikasi stok akhir = **0 persis** (bukan mendekati 0 seperti 0.0001 sisa akibat floating point) lewat `products/detail` DAN `packages/list` untuk keempat satuan sekaligus.
   - Coba jual **1 Butir lagi** (satuan terkecil) setelah stok 0 — harus ditolak bersih "stok tidak mencukupi", TIDAK boleh stok jadi minus.
   - Coba juga jual satuan lain (Kilogram/Setengah Kilo/Seperempat Kilo) setelah stok 0 — semua harus ditolak konsisten, bukan cuma satuan terkecil yang diblokir.
   - Cek tampilan Kasir: produk dengan stok 0 harus tampil sebagai "Stok Habis" di SEMUA 4 tombol satuan sekaligus (bukan cuma di satuan yang kebetulan dipakai transaksi terakhir).
7. Setelah stok 0, lakukan pembelian kecil lagi (mis. 1 Butir = 0.0625 kg) — cek produk kembali bisa dijual, stok tepat 0.0625, tombol "Stok Habis" hilang lagi.
8. Race condition (opsional, kalau mau uji ketat): 2 sesi kasir konkuren sama-sama coba jual satuan BERBEDA (mis. satu jual Kilogram, satu jual Butir) dari stok pas-pasan yang sama — pastikan locking tetap benar meski satuan yang dipakai kedua sesi berbeda-beda (bukan cuma diuji dengan satuan yang sama seperti di Fase G sebelumnya).
9. Void transaksi yang berisi baris Telur multi-satuan — cek stok dikembalikan tepat sesuai satuan aslinya per baris (bukan cuma dikembalikan dalam satuan anchor secara kasar).
10. Retur dari transaksi Kasir (kalau modul retur pelanggan/penjualan ada) untuk salah satu baris Telur — cek presisi sama seperti void.
11. **Diskon per item** pada baris Telur — coba diskon persen (mis. 10%) dan diskon nominal (mis. Rp 2.000) pada satuan yang BERBEDA-beda (satu transaksi diskon di baris Kilogram, transaksi lain diskon di baris Butir) — cek subtotal setelah diskon terhitung benar untuk tiap satuan, tidak ada pembulatan aneh khususnya di baris Butir yang harga satuannya kecil.
12. **Campur Telur dengan produk LAIN (non-Telur)** dalam satu keranjang yang sama (mis. Telur 1 kg + rokok 1 slop + minyak 0.5 kg) — cek subtotal & total keseluruhan benar, dan proses stok masing-masing produk independen (transaksi Telur tidak mempengaruhi stok produk lain, dan sebaliknya).
13. **Cetak struk** untuk transaksi berisi beberapa baris Telur dengan satuan berbeda-beda (mis. dari skenario 5 keranjang campuran) — cek nama satuan yang tercetak sesuai (`Kilogram`/`Setengah Kilo`/`Seperempat Kilo`/`Butir`, bukan generik atau salah label), harga & subtotal per baris cocok dengan yang ada di layar.
14. Jual sebagai **role kasir (bukan admin)** — pastikan kasir bisa menjual Telur di semua 4 satuan tanpa hambatan permission (transaksi penjualan memang wilayah kasir, beda dengan edit stok/hapus produk yang harus admin).

**Prompt eksekusi**:
> Jalankan Fase T3 — Transaksi Kasir dari `docs/RENCANA_TESTING_PRODUK_TELUR_MULTISATUAN.md`. Pastikan Fase T1 (dan idealnya T2) sudah selesai supaya ada cukup stok. WAJIB jalankan skenario 6 (jual bertahap sampai stok benar-benar 0) sebagai salah satu fokus utama — rancang urutan transaksi dengan berbagai satuan sampai stok tepat habis, lalu verifikasi penolakan bersih saat dijual lagi. Uji lewat browser sungguhan (Playwright), variasikan jumlah di tiap satuan (kecil/sedang/besar — lihat Matriks Variasi Jumlah di bagian 3.5), termasuk skenario keranjang multi-satuan, diskon per item, campur dengan produk lain, cetak struk, dan race condition. Laporkan hasil per skenario dengan format standar, sertakan contoh sederhana untuk tiap bug ditemukan.

---

### Fase T4 — Lintas Modul & Presisi (gabungan)

**Cakupan**: Skenario yang menggabungkan pembelian + retur + penjualan dalam satu alur, plus audit presisi menyeluruh — memastikan breakdown 4 satuan tetap konsisten dari ujung ke ujung.

**Skenario wajib dicoba:**
1. Alur penuh: beli 10 kg → jual 3 kg (macam-macam satuan) → retur 1 kg dari sisa pembelian → cek stok akhir = 10 - 3 - 1 = 6 kg tepat, dan `stock_mutations` runtut sesuai urutan waktu asli (bukan urutan input).
2. Ubah rasio salah satu paket (mis. ternyata supplier baru bikin rasio berubah jadi "0.25 kg = 5 butir") pada produk yang SUDAH ada stok & riwayat transaksi — verifikasi & laporkan perilaku aktual (apakah stok anchor ikut "terjemahkan ulang" seperti area abu-abu yang sudah ditemukan di Fase G sebelumnya, atau beda kasusnya di sini). Ini kemungkinan besar akan mereplikasi temuan gray-area Fase G — kalau iya, cukup catat sebagai konfirmasi tambahan, tidak perlu dianggap bug baru.
3. Laporan Stok & Dashboard: cek produk Telur muncul dengan angka stok kg yang benar, dan breakdown "Sisa X Kilogram" atau semacamnya (kalau ada fitur breakdown tampilan seperti "9 Pack" dari proyek sebelumnya) menampilkan estimasi butir/paket yang masuk akal.
4. Laporan Laba Rugi: transaksi Telur dengan berbagai satuan tercatat HPP (harga pokok) yang benar per baris, sesuai `purchase_price` di masing-masing paket (bukan disamaratakan pakai harga Kilogram untuk semua satuan).
5. Precision stress test: lakukan 20+ transaksi kecil campur satuan (kilogram/setengah/seperempat/butir) berturut-turut, lalu bandingkan total stok akhir sistem vs hitungan manual dari seluruh log `stock_mutations` — harus sama persis sampai digit terakhir (tidak ada drift kumulatif dari pembulatan).
6. **Pencarian & filter**: cari produk Telur lewat kotak pencarian Kasir dan halaman Produk (pakai keyword sebagian nama) — cek semua 4 tombol satuan muncul benar di hasil pencarian Kasir, dan filter status/stok-menipis di halaman Produk memperlakukan Telur dengan benar (mis. kalau stok di bawah `min_stock`, badge "Stok Menipis" muncul, dihitung dari stok anchor kg).
7. **Sortir tabel Produk by Stok** (asc/desc) — pastikan posisi Telur di tabel terurut benar relatif ke produk lain berdasarkan angka stok kg-nya, tidak "salah baca" karena py 4 satuan.
8. Konsistensi lintas 2+ ronde pengujian: ulangi salah satu skenario T1-T3 (mis. jual sampai 0 di T3 skenario 6) dengan **data fixture baru** sekali lagi setelah semua bug di atas diperbaiki — pastikan tidak ada regresi, mengikuti pola metodologi re-test berulang yang sudah terbukti berguna di Fase G sebelumnya (bug baru justru sering ketemu di ronde ke-2/3, bukan ronde pertama).

**Prompt eksekusi**:
> Jalankan Fase T4 — Lintas Modul & Presisi dari `docs/RENCANA_TESTING_PRODUK_TELUR_MULTISATUAN.md`. Pastikan Fase T1-T3 sudah selesai. Fokus pada audit presisi menyeluruh dan interaksi antar modul, termasuk pencarian/filter/sortir dan re-test ronde kedua untuk cek regresi. Laporkan hasil per skenario dengan format standar, sertakan contoh sederhana untuk tiap bug/area abu-abu ditemukan, dan tutup dengan ringkasan kesimpulan apakah skema 4-satuan ini layak dipakai produksi untuk produk sejenis (telur, dan kemungkinan produk timbangan lain yang punya pola serupa).

---

### 3.5 Matriks Variasi Jumlah (checklist ringkas — dipakai lintas Fase T1/T2/T3)

Supaya variasi jumlah benar-benar tercakup luas (bukan cuma satu-dua angka contoh), pakai daftar angka berikut sebagai checklist tiap kali suatu skenario menyebut "variasikan jumlah":

| Satuan | Minimal | Kecil | Sedang | Besar | Sangat besar | Boundary/aneh |
|---|---|---|---|---|---|---|
| **Kilogram** | 0.0625 (=1 butir) | 0.25 / 0.5 | 2 / 5 | 10 / 20 | 50 / 100 | 0.1234 (presisi tinggi, non-bulat aneh) |
| **Setengah Kilo** | 1 paket | 2-3 paket | 5-8 paket | 15-20 paket | 50 paket | — (satuan ini diskrit, tidak ada pecahan paket) |
| **Seperempat Kilo** | 1 paket | 3-5 paket | 10-12 paket | 30-40 paket | 100 paket | — |
| **Butir** | 1 butir | 4-8 butir (kelipatan 4) | 13-15 butir (BUKAN kelipatan 4) | 50-100 butir | 500 butir | 999999 butir (uji ekstrem, cek tidak overflow) |

**Cara pakai**: setiap skenario yang menyebut "variasikan jumlah kecil/sedang/besar" di Fase T1-T3, ambil minimal 3 titik berbeda dari kolom yang relevan (mis. Minimal + Sedang + Sangat besar) supaya cakupannya representatif, bukan cuma 1 angka. Skenario yang eksplisit menyebut "jual sampai 0" (T3 skenario 6) WAJIB memakai kombinasi dari BEBERAPA baris tabel ini sekaligus dalam satu alur, bukan cuma satu satuan saja.

---

## 4. Format pelaporan bug (ikuti standar yang sama dengan `RENCANA_TESTING_QA_MENYELURUH.md`)

Untuk setiap bug yang ditemukan di fase manapun, laporkan dengan format:
- **Lokasi**: file/fungsi/endpoint
- **Langkah reproduksi**: runtut, dengan angka konkret
- **Hasil aktual vs diharapkan**
- **Severity**: Mayor / Minor / Area abu-abu (butuh keputusan bisnis)
- **Perbaikan**: (kalau bug kode murni, boleh langsung diperbaiki setelah dikonfirmasi lewat reproduksi nyata; kalau butuh keputusan desain/bisnis, TANYA dulu ke user)
- **Contoh sederhana**: analogi dunia nyata yang mudah dipahami (wajib disertakan, sesuai preferensi user)

## 5. Status eksekusi

| Fase | Status |
|---|---|
| Setup (bagian 2) | ✅ Selesai — produk id 235 "QATEST-TELUR-Ayam", satuan id 23/24/25, paket id 378-381, struktur linear terverifikasi tidak bercabang. |
| T1 — Pembelian | ✅ **SELESAI (21 Agu 2026)** — 9/9 skenario PASS. **1 bug MAYOR ditemukan & diperbaiki**: form Pembelian Supplier memaksa qty bilangan bulat untuk semua satuan (termasuk Kilogram yang harusnya boleh desimal) — lihat detail lengkap di bagian Fase T1 di atas. |
| T2 — Retur | 🔶 Sempat dijalankan sebagian (5 dari 11 skenario versi lama, sebelum revisi dokumen ini) — MENEMUKAN 1 indikasi bug (pesan error `%.0f` salah tampilkan sisa retur pecahan sebagai bilangan bulat) yang **belum diperbaiki**. Perlu direproduksi ulang & diperbaiki di eksekusi berikutnya — PRIORITASKAN juga cek apakah form Retur punya bug qty-desimal yang sama seperti yang baru ditemukan di form Pembelian (belum pernah dicek eksplisit). |
| T3 — Kasir | ⬜ Belum dijalankan sama sekali — PRIORITASKAN cek juga apakah ada bug qty-desimal serupa di form terkait Kasir (kemungkinan kecil karena Kasir sudah pernah diperbaiki khusus di Fase G, tapi tetap perlu diverifikasi ulang untuk produk 4-satuan spesifik ini). |
| T4 — Lintas Modul & Presisi | ⬜ Belum dijalankan — WAJIB investigasi selisih `products/detail.stock` (38) vs total `stock_mutations` (36.063) yang ditemukan selama T1 (lihat catatan di bagian Fase T1). |

**Catatan**: kalau mau mulai fresh (data baru, tidak pakai sisa dari percobaan sebelumnya), boleh minta produk/satuan test di atas diabaikan dan set baru dibuat dengan prefix berbeda (mis. `QATEST-TELUR2-*`) — ikuti aturan baku proyek: JANGAN hapus data lama, cukup buat set baru di sampingnya.
