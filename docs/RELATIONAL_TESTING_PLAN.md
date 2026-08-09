# Rencana Testing Relasional Antar-Menu (via Browser)

> Dasar: hasil pemetaan relasi menu (11 grup) setelah ditemukannya bug Dashboard vs Kas Saya.
> Tujuan: **bukan** testing "buka menu satu-satu pastikan tidak error" (itu sudah dilakukan di smoke test Fase 8 migrasi timezone). Ini testing **lintas menu** — memastikan angka/status yang sama, yang ditampilkan di beberapa halaman berbeda, benar-benar konsisten satu sama lain pada state data yang sama.
> Semua testing **wajib lewat browser sungguhan** (bukan cuma curl API) — karena bug aslinya (Dashboard vs Kas Saya) hanya kelihatan kalau kedua halaman dibuka dan dibandingkan visualnya.

---

## Kenapa Pendekatan Ini Berbeda dari Smoke Test Biasa

Smoke test per-menu menjawab: *"apakah halaman ini rusak?"*
Testing relasional menjawab: *"apakah halaman ini **berbohong** dibanding halaman lain yang menampilkan data yang sama?"*

Sebuah halaman bisa tampil sempurna (tidak error, tidak crash, datanya "masuk akal") tapi tetap salah — kalau angkanya beda dari halaman lain yang seharusnya menampilkan hal yang sama. Bug Dashboard vs Kas Saya persis begini: kedua-duanya tampil normal, tidak ada error di console, tapi salah satu bilang "Kas Buka" dan satunya bilang "Tutup". Testing per-menu tidak akan pernah menangkap ini karena masing-masing "terlihat benar" saat dicek sendirian.

## Metodologi Setiap Fase

Setiap fase di bawah ini mengikuti pola yang sama:

1. **Siapkan state yang diketahui** — lewat browser (bukan cuma API), lakukan satu aksi konkret (buka kas, buat transaksi, void transaksi, dst.)
2. **Screenshot halaman A** segera setelah aksi
3. **Screenshot halaman B, C, D...** (semua halaman yang menampilkan data yang sama) **tanpa mengubah apapun di antaranya**
4. **Bandingkan angka/status secara eksplisit** — bukan "kelihatannya benar", tapi tulis angka dari tiap halaman berdampingan
5. **Uji juga skenario pembalikan** (void, batal, hapus) — apakah semua halaman ikut ter-update, atau ada yang "nyangkut" nilai lama
6. Kalau ketemu bug: dokumentasikan dengan bukti sama seperti bug Dashboard/Kas Saya kemarin (root cause di kode, bukan cuma gejala)

Testing ini **tidak mengubah kode** kecuali ditemukan bug dan diminta diperbaiki — fokus fase ini murni investigasi/pembuktian.

---

## Urutan Fase (berdasarkan risiko)

| Fase | Grup | Kenapa urutan ini |
|---|---|---|
| R1 | Grup 8 — Agregasi Finansial (Dashboard/Ringkasan Bisnis/Dashboard Keuangan/Laba Rugi) | Pola paling mirip bug asli, dan lebih parah (3-4 query independen) |
| R2 | Grup 1 — Kas (verifikasi ulang lebih dalam pasca-fix) | Sudah pernah ada bug di sini, perlu dibuktikan benar-benar tuntas |
| R3 | Grup 2 — Transaksi → Stok → Laporan | Efek samping tulis-lalu-baca, dampak operasional harian |
| R4 | Grup 3 — Transaksi Kredit → Piutang | Sama seperti R3, fokus ke piutang |
| R5 | Grup 4 — Shift ↔ Kas ↔ Kinerja Kasir | Menyambung dari R2, soal validasi shift |
| R6 | Grup 5 & 6 — Pembelian ↔ Retur ↔ Stok | Alur pengadaan, dampak stok |
| R7 | Grup 7 — Produk ↔ Kategori ↔ Unit | Master data, risiko rendah tapi cepat dicek |
| R8 | Grup 11 — Batch Expired ↔ Produk ↔ Pembelian | Terkait Pembelian, sekalian |
| R9 | Grup 9 — Sync Center | Paling kompleks (offline sync), dikerjakan setelah semua alur normal teruji |
| R10 | Grup 10 — Users/Roles/Access/Menus | Fungsional (hak akses), bukan soal angka — beda sifat, dikerjakan terakhir |

---

## FASE R1 — Agregasi Finansial (Dashboard, Ringkasan Bisnis, Dashboard Keuangan, Laba Rugi)

### Skenario yang harus diuji lewat browser
1. Catat kondisi awal: buka keempat halaman ini secara berurutan di browser yang sama (login sekali), screenshot masing-masing, catat angka "Pendapatan/Penjualan Hari Ini" dan "Laba" yang tampil di tiap halaman.
2. Buat 1 transaksi tunai baru lewat menu Kasir (jumlah yang gampang dilacak, misal Rp 77.000).
3. Reload keempat halaman satu per satu (Dashboard → Ringkasan Bisnis → Dashboard Keuangan → Laba Rugi), screenshot tiap halaman.
4. Bandingkan: apakah kenaikan Rp 77.000 itu muncul **sama persis** di keempat halaman, atau ada yang beda/telat/tidak berubah?
5. Ulangi dengan 1 pengeluaran (expense) baru — cek apakah semua halaman mengurangi laba dengan jumlah yang sama.
6. Ulangi dengan transaksi kredit (menambah piutang) — cek apakah "Laba Kotor"/pendapatan menghitungnya secara konsisten sebagai pendapatan atau tidak, di semua halaman.
7. Void transaksi yang baru dibuat di langkah 2 — cek apakah keempat halaman **sama-sama** kembali ke angka semula.
8. Ganti rentang tanggal filter (kalau ada) di Ringkasan Bisnis dan Laba Rugi ke "Bulan Ini" — bandingkan totalnya dengan cara jumlahkan manual dari Laporan Penjualan pada rentang yang sama.

### PROMPT EKSEKUSI — Fase R1
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R1 (Agregasi Finansial) — testing relasional lintas menu via browser, BUKAN testing per-menu biasa.

Konteks: sebelumnya ditemukan bug dimana Dashboard dan Kas Saya menampilkan status kas yang berbeda untuk data yang sama, karena masing-masing pakai query SQL terpisah. Grup ini (Dashboard, Ringkasan Bisnis / Business Summary, Dashboard Keuangan / Finance Overview, Laba Rugi / Profit-Loss) punya pola serupa: keempatnya independen menghitung pendapatan/laba dari query masing-masing (business_summary_repo.go, dashboard_repo.go, report_repo.go semuanya beda file, beda query).

TESTING WAJIB (semua lewat browser, gunakan playwright-core headless jika tersedia seperti sesi sebelumnya, atau tools browser testing yang ada):
1. Jalankan BE dan FE. Login sebagai admin/admin123.
2. Buka 4 halaman ini satu per satu, screenshot masing-masing, catat angka "pendapatan/penjualan hari ini" dan "laba" yang tertulis di layar:
   - Dashboard (/dashboard)
   - Ringkasan Bisnis (/reports/business-summary)
   - Dashboard Keuangan (/finance)
   - Laba Rugi (/reports/profit-loss)
3. Buat 1 transaksi tunai baru via Kasir dengan nominal yang gampang dilacak (misal Rp 77.000) — kalau perlu buka kas dulu.
4. Reload keempat halaman di langkah 2, screenshot ulang, bandingkan angka SEBELUM vs SESUDAH secara eksplisit dalam bentuk tabel di laporan akhir. WAJIB semua 4 halaman naik dengan jumlah yang sama.
5. Buat 1 pengeluaran (expense) baru dengan nominal yang jelas. Reload 4 halaman lagi, cek laba turun konsisten di semua halaman.
6. Buat 1 transaksi kredit (is_credit) ke pelanggan. Reload 4 halaman, cek bagaimana masing-masing memperlakukan pendapatan kredit ini (apakah dihitung sebagai pendapatan langsung atau tidak) — WAJIB konsisten cara hitungnya di semua halaman, meskipun cara hitungnya sendiri boleh "belum dibayar tidak dihitung", yang penting SAMA di keempat halaman.
7. Void transaksi tunai dari langkah 3. Reload 4 halaman, pastikan keempatnya kembali ke angka SEBELUM langkah 3 (dari tabel di langkah 2).
8. Kalau Ringkasan Bisnis dan Laba Rugi punya filter tanggal, ganti ke rentang yang sama (misal "Bulan Ini"), bandingkan totalnya satu sama lain, dan silangkan dengan total manual dari Laporan Penjualan (/reports/sales) pada rentang yang sama.
9. Untuk SETIAP temuan selisih/inkonsisten: baca kode repo yang relevan (business_summary_repo.go, dashboard_repo.go/dashboard_service.go, report_repo.go, finance-related repo) untuk cari root cause pastinya (bukan dugaan) — sebutkan file dan baris kode yang jadi biang keroknya, tapi JANGAN perbaiki dulu, cukup laporkan dengan bukti lengkap (screenshot + angka + root cause di kode).

Laporkan hasil lengkap: tabel perbandingan angka per langkah, screenshot, dan kalau ada bug, root cause di kode. Jangan lanjut ke fase lain.
```

---

## FASE R2 — Kas (Verifikasi Ulang Lebih Dalam)

### Skenario yang harus diuji lewat browser
1. Buka kas dari Dashboard dengan user A. Cek konsisten di Kas Saya dan Kas Harian.
2. Dari Kasir, buat 2-3 transaksi tunai dan 1 transaksi non-tunai (QRIS/transfer, kalau ada). Cek total "Total Masuk Tunai" di Kas Saya vs Kas Harian vs saldo ekspektasi di Dashboard — harus sama.
3. Void salah satu transaksi tunai tadi. Cek total kas di ketiga halaman turun dengan jumlah yang sama.
4. Tutup kas dari Kas Saya dengan saldo akhir yang sengaja beda dari ekspektasi (misal selisih Rp 5.000) — cek "Selisih" tampil konsisten di Kas Harian.
5. Coba buka kas lagi tanpa tutup kas sebelumnya (harus ditolak) — cek pesan errornya konsisten kalau dicoba dari Dashboard maupun Kas Saya.
6. **Multi-shift dalam sehari**: tutup kas, buka kas baru lagi dengan shift berbeda di hari yang sama — cek Kas Harian menampilkan 2 baris terpisah dengan benar, dan Dashboard/Kas Saya menunjuk ke sesi yang sedang aktif (bukan yang sudah ditutup).

### PROMPT EKSEKUSI — Fase R2
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R2 (Kas — verifikasi ulang lebih dalam) — testing relasional lintas menu via browser.

Konteks: bug asli (Dashboard vs Kas Saya beda status) sudah diperbaiki di getMyCashQuery (cash_drawer_repo.go). Fase ini untuk membuktikan perbaikannya benar-benar tuntas dengan skenario yang lebih dalam dari sekadar cek status buka/tutup, termasuk multi-transaksi dan multi-shift.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123.
2. Buka kas dari Dashboard. Screenshot Dashboard, Kas Saya (/finance/my-cash), Kas Harian (/finance/cash-drawer) — pastikan status dan saldo awal konsisten di ketiganya.
3. Buat 2-3 transaksi tunai dan (kalau tersedia metode pembayaran non-tunai) 1 transaksi non-tunai via Kasir. Screenshot ulang ketiga halaman, bandingkan "Total Masuk Tunai"/"Saldo Ekspektasi" — WAJIB sama di ketiganya.
4. Void salah satu transaksi tunai (dari halaman Transaksi). Screenshot ulang ketiga halaman, pastikan totalnya turun dengan jumlah yang SAMA di ketiganya.
5. Tutup kas dari Kas Saya, masukkan saldo akhir yang SENGAJA beda dari saldo ekspektasi (misal kurang Rp 5.000). Screenshot hasil "Selisih" di Kas Saya, lalu cek Kas Harian menampilkan selisih yang SAMA PERSIS.
6. Coba buka kas lagi TANPA menutup kas manapun yang mungkin masih aktif (kalau sudah ditutup di langkah 5, buka baru itu wajar; tapi coba juga trigger error "sudah ada kas terbuka" dengan membuka 2x berturutan) — cek pesan error yang tampil konsisten baik dicoba dari tombol di Dashboard maupun di Kas Saya.
7. Setelah kas dari langkah 5 ditutup, buka kas BARU dengan shift yang berbeda (hari yang sama). Buat 1 transaksi. Screenshot Kas Harian — pastikan tampil 2 baris riwayat terpisah (sesi pertama tertutup, sesi kedua terbuka) dengan data yang benar masing-masing, dan Dashboard/Kas Saya menunjuk ke sesi KEDUA (yang aktif), bukan tercampur dengan sesi pertama.

Laporkan hasil lengkap dengan screenshot tiap langkah dan tabel perbandingan angka. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu, laporkan dulu). Jangan lanjut ke fase lain.
```

---

## FASE R3 — Transaksi → Stok → Laporan

### Skenario yang harus diuji lewat browser
1. Catat stok produk tertentu di halaman Produk.
2. Jual produk itu (qty tertentu) lewat Kasir.
3. Cek stok di Produk berkurang sesuai qty, dan riwayat mutasi stok (kalau ada halaman/expand-nya) mencatat "out" dengan qty yang sama.
4. Cek transaksi itu muncul di Laporan Penjualan dan Laba Rugi dengan angka yang benar.
5. Void transaksi itu dari halaman Transaksi.
6. Cek stok di Produk **kembali** ke angka semula, ada mutasi "void" tercatat, dan transaksi itu hilang/terkoreksi dari Laporan Penjualan & Laba Rugi.
7. Ulangi dengan produk yang punya varian satuan/paket (kalau ada) — pastikan konversi qty ke stok tetap benar dan konsisten di semua halaman.

### PROMPT EKSEKUSI — Fase R3
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R3 (Transaksi → Stok → Laporan) — testing relasional lintas menu via browser.

Konteks: transaction_repo.go (Go/backend) langsung UPDATE products.stock dan INSERT ke stock_mutations dalam transaksi database yang sama saat sebuah penjualan dibuat, dan mengembalikannya saat transaksi di-void. Fase ini membuktikan efek tulis-lalu-baca ini benar-benar tercermin di semua halaman yang menampilkan stok/laporan.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123.
2. Buka halaman Produk (/products), catat stok produk tertentu (screenshot).
3. Jual produk itu qty tertentu (misal 3) lewat Kasir (pastikan kas sudah terbuka).
4. Screenshot Produk lagi — stok WAJIB berkurang tepat 3.
5. Buka Laporan Penjualan (/reports/sales) dan Laba Rugi (/reports/profit-loss) untuk hari ini — screenshot, catat apakah transaksi/nilai penjualan produk ini sudah tercermin dengan benar.
6. Void transaksi tadi dari halaman Transaksi (/transactions).
7. Screenshot Produk lagi — stok WAJIB kembali ke angka semula (langkah 2).
8. Screenshot ulang Laporan Penjualan dan Laba Rugi — transaksi yang di-void WAJIB tidak lagi dihitung sebagai penjualan aktif (baik hilang dari daftar, atau tetap muncul tapi berstatus void dan tidak masuk ke total).
9. Kalau produk yang dipakai punya satuan/paket alternatif (unit_id/package selain satuan dasar), ulangi langkah 3-8 dengan menjual pakai satuan alternatif itu — pastikan konversi qty ke stok dasar tetap benar dan konsisten antara Produk, Kasir (saat input), dan riwayat mutasi.

Laporkan hasil lengkap dengan screenshot tiap langkah, angka stok sebelum/sesudah/setelah-void dalam bentuk tabel. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R4 — Transaksi Kredit → Piutang

### Skenario yang harus diuji lewat browser
1. Buat transaksi kredit (is_credit) ke pelanggan tertentu via Kasir.
2. Cek muncul di Piutang dengan status "Hutang"/unpaid dan jumlah yang benar.
3. Cek juga muncul di halaman detail Pelanggan (kalau ada riwayat piutang per pelanggan).
4. Bayar sebagian piutang itu dari halaman Piutang — cek status berubah jadi "partial", sisa piutang berkurang benar.
5. Void transaksi kredit itu dari Transaksi (kalau sistem mengizinkan void transaksi yang sudah ada pembayaran piutang — kalau ditolak, itu perilaku benar, catat sebagai temuan positif bukan bug).
6. Kalau void berhasil (untuk transaksi kredit yang BELUM dibayar sama sekali): cek piutang di Piutang ikut ter-void/hilang, tidak nyangkut sebagai piutang "hantu".

### PROMPT EKSEKUSI — Fase R4
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R4 (Transaksi Kredit → Piutang) — testing relasional lintas menu via browser.

Konteks: transaction_repo.go INSERT ke tabel receivables saat transaksi is_credit=true, dan UPDATE status='void' pada receivable terkait saat transaksinya di-void. Fase ini membuktikan alur ini konsisten dari sisi Kasir/Transaksi ke Piutang dan Pelanggan.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123. Pastikan ada minimal 1 pelanggan terdaftar (buat dulu di /customers kalau belum ada).
2. Buat transaksi kredit (is_credit=true) ke pelanggan itu via Kasir, nominal yang gampang dilacak.
3. Screenshot Piutang (/receivables) — transaksi WAJIB muncul dengan status "Hutang"/unpaid dan jumlah yang sama persis dengan transaksi.
4. Kalau ada halaman/detail riwayat piutang per pelanggan di /customers, screenshot itu juga — pastikan datanya sama dengan yang di Piutang.
5. Bayar SEBAGIAN piutang itu dari halaman Piutang (misal bayar setengahnya). Screenshot — status WAJIB berubah jadi "partial"/sebagian, sisa piutang berkurang tepat sesuai pembayaran.
6. Coba void transaksi kredit itu dari halaman Transaksi. Screenshot hasilnya:
   - Kalau sistem MENOLAK void (karena sudah ada pembayaran) — itu perilaku yang masuk akal, catat sebagai temuan POSITIF (bukan bug), lanjut ke langkah 7 pakai transaksi kredit BARU yang belum dibayar sama sekali.
   - Kalau sistem MENGIZINKAN void — screenshot Piutang setelahnya, pastikan piutang itu ikut berubah status (void/dibatalkan), TIDAK boleh nyangkut sebagai piutang aktif yang tidak ada transaksi validnya.
7. Buat transaksi kredit BARU (belum dibayar sama sekali), langsung void dari Transaksi. Screenshot Piutang — piutang itu WAJIB ikut hilang/void, tidak nyangkut.

Laporkan hasil lengkap dengan screenshot tiap langkah. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R5 — Shift ↔ Kas ↔ Kinerja Kasir

### Skenario yang harus diuji lewat browser
1. Buat shift baru di halaman Shift dengan jam operasional tertentu.
2. Buka kas dengan shift itu, buat beberapa transaksi.
3. Cek performa kasir (kalau ada laporannya) menampilkan angka yang sama dengan yang tercatat di kas untuk shift itu.
4. Coba buat transaksi dengan shift_id yang TIDAK sesuai kas yang sedang terbuka — pastikan ditolak dengan pesan yang jelas (ini business rule yang sudah ada, verifikasi masih jalan).
5. Nonaktifkan shift itu dari halaman Shift — cek apakah kas yang masih terbuka dengan shift itu tetap normal (tidak rusak), dan shift yang nonaktif tidak muncul lagi sebagai pilihan saat buka kas baru.

### PROMPT EKSEKUSI — Fase R5
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R5 (Shift ↔ Kas ↔ Kinerja Kasir) — testing relasional lintas menu via browser.

Konteks: transaction_service.Create() memvalidasi req.ShiftID harus sama dengan shift_id milik kas yang sedang terbuka sebelum mengizinkan transaksi. shift_repo.go membaca dari cash_drawer dan transactions untuk laporan. Fase ini membuktikan keterkaitan Shift-Kas-Laporan Kinerja Kasir konsisten.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123.
2. Buat shift baru di /shifts dengan jam operasional yang jelas (misal 08:00-16:00), screenshot.
3. Buka kas dengan shift itu dari Dashboard. Buat 3 transaksi dengan nominal berbeda-beda yang mudah dijumlahkan (misal 10rb, 20rb, 30rb = total 60rb).
4. Buka Laporan Kinerja Kasir (/reports/cashier) untuk hari ini, screenshot — pastikan total penjualan kasir yang login sekarang untuk shift ini menunjukkan Rp 60.000 dan 3 transaksi, sama dengan yang tercatat di Kas Saya/Kas Harian.
5. Coba buat transaksi baru tapi dengan shift_id yang BEDA dari kas yang sedang terbuka (kalau UI mengizinkan memilih shift saat transaksi — kalau tidak ada pilihan shift eksplisit di Kasir, catat itu sebagai temuan: berarti validasi ini tidak bisa dites dari UI, cukup laporkan). Pastikan sistem menolak dengan pesan error yang jelas kalau memang ada jalur untuk trigger ini.
6. Nonaktifkan shift dari langkah 2 (toggle status di /shifts) SEMENTARA kas dari langkah 3 masih terbuka. Screenshot Kas Saya/Kas Harian — pastikan kas yang sudah terbuka itu TETAP normal (tidak rusak/hilang), karena menonaktifkan shift seharusnya cuma mencegah PEMILIHAN BARU, bukan merusak data yang sudah ada.
7. Coba buka kas BARU (dengan user/sesi berbeda kalau perlu, atau setelah tutup kas yang lama) — pastikan shift yang baru dinonaktifkan TIDAK muncul lagi di pilihan dropdown shift.

Laporkan hasil lengkap dengan screenshot tiap langkah. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R6 — Pembelian ↔ Retur ↔ Stok

### Skenario yang harus diuji lewat browser
1. Catat stok produk tertentu.
2. Buat Pembelian dari Supplier dengan qty tertentu.
3. Cek stok di Produk bertambah sesuai qty, mutasi stok "in" tercatat.
4. Buat Retur dari Pembelian itu dengan qty sebagian.
5. Cek stok di Produk berkurang sesuai qty retur, dan sisa qty "net received" di halaman Pembelian (kalau ditampilkan) ter-update benar.
6. Coba retur melebihi qty yang tersisa (harus ditolak) — verifikasi validasi masih jalan.
7. Kalau pembelian punya expiry batch: cek batch itu muncul konsisten di Produk dan tersambung ke Pembelian asalnya.

### PROMPT EKSEKUSI — Fase R6
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R6 (Pembelian ↔ Retur ↔ Stok, termasuk Grup 11 Batch Expired) — testing relasional lintas menu via browser.

Konteks: purchase_repo.go dan supplier_return_repo.go sama-sama langsung UPDATE products.stock dan mencatat stock_mutations dalam raw SQL. Retur menghitung ulang "net received" terhadap purchase aslinya. Fase ini membuktikan rantai Pembelian → Stok → Retur → Stok konsisten, termasuk data expiry batch kalau ada.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123. Pastikan ada minimal 1 supplier terdaftar.
2. Buka Produk (/products), catat stok produk tertentu (screenshot).
3. Buat Pembelian (/suppliers/purchases) untuk produk itu, qty tertentu (misal 20), isi expiry date kalau formnya menyediakan.
4. Screenshot Produk lagi — stok WAJIB bertambah tepat 20.
5. Kalau ada halaman/expand untuk expiry batch di Produk, screenshot itu juga — pastikan batch baru muncul dengan qty dan tanggal expired yang sama dengan yang diinput di Pembelian.
6. Buat Retur (/suppliers/returns) dari purchase order tadi, qty sebagian (misal 5).
7. Screenshot Produk lagi — stok WAJIB berkurang tepat 5 (jadi net +15 dari awal).
8. Screenshot detail Pembelian asalnya — kalau ada info "sudah diretur X dari Y", pastikan menunjukkan 5 dari 20 dengan benar.
9. Coba buat retur lagi dengan qty yang melebihi sisa (misal coba retur 20 lagi padahal cuma sisa 15 net) — pastikan sistem menolak dengan pesan jelas.

Laporkan hasil lengkap dengan screenshot tiap langkah dan tabel angka stok di tiap tahap. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R7 — Produk ↔ Kategori ↔ Unit

### Skenario yang harus diuji lewat browser
1. Buat kategori baru, buat produk dengan kategori itu.
2. Ganti nama kategori — cek nama yang tampil di listing Produk ikut berubah (bukan cache lama).
3. Coba hapus kategori yang masih dipakai produk — harus ditolak/diberi peringatan jelas, bukan silent fail atau produk jadi orphan.
4. Sama untuk Unit — ganti nama satuan, cek konsisten di Produk; coba hapus satuan yang dipakai, harus ditolak.

### PROMPT EKSEKUSI — Fase R7
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R7 (Produk ↔ Kategori ↔ Unit) — testing relasional lintas menu via browser.

Konteks: product_service menginjeksi CategoryRepo dan UnitRepo langsung untuk validasi. Fase ini membuktikan perubahan di master data (Kategori/Unit) tercermin benar di Produk, dan tidak bisa menciptakan data yatim (orphan).

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123.
2. Buat kategori baru "Kategori Relasi Test" di /categories.
3. Buat produk baru dengan kategori itu di /products.
4. Ganti nama kategori itu jadi "Kategori Relasi Test - Updated" di /categories.
5. Reload /products — nama kategori yang tampil di baris produk tadi WAJIB ikut berubah jadi nama baru (bukan nama lama yang ke-cache).
6. Coba nonaktifkan/hapus kategori itu SEMENTARA masih dipakai produk aktif — screenshot hasilnya, pastikan sistem menolak dengan pesan jelas (bukan berhasil diam-diam meninggalkan produk dengan category_id yatim).
7. Ulangi langkah 2-6 dengan Unit (/units) sebagai gantinya Kategori — buat unit baru, pasang ke produk, ganti nama, cek konsisten di Produk, coba hapus yang masih dipakai.

Laporkan hasil lengkap dengan screenshot tiap langkah. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R8 — Sync Center

### Skenario yang harus diuji lewat browser
1. Buka halaman Sync Center, cek state awal (antrian kosong / ada riwayat).
2. Kalau ada cara mensimulasikan transaksi "offline" dari UI (device lain / mode simulasi), lakukan itu.
3. Jalankan proses sync, cek transaksi itu muncul konsisten di Transaksi, Produk (stok berkurang), Kas (kalau tunai).
4. Kalau memungkinkan memicu konflik sync (dua perubahan pada data yang sama), cek halaman resolve-conflict menampilkan data yang benar dari kedua sisi, dan hasil resolusinya konsisten di semua halaman terkait.

### PROMPT EKSEKUSI — Fase R8
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R8 (Sync Center) — testing relasional lintas menu via browser.

Konteks: sync_service menginjeksi transactionRepo, expenseRepo, productRepo, cashDrawerRepo langsung -- artinya proses sync menyentuh 4 domain sekaligus dalam satu alur. Ini fase paling kompleks, kerjakan setelah R1-R7 selesai.

TESTING WAJIB (lewat browser):
1. Jalankan BE dan FE, login admin/admin123.
2. Buka halaman Sync/Pusat Sinkronisasi, screenshot state awal (antrian, riwayat konflik).
3. Cari tahu dari kode (BE/domain/sync) apakah ada endpoint test/manual untuk push sync item secara manual (misal lewat request builder atau form di UI). Kalau UI menyediakan cara mensimulasikan push sync (device_id, local_id, payload transaksi/expense), gunakan itu.
4. Kalau tidak ada jalur UI untuk mensimulasikan sync push, laporkan itu sebagai keterbatasan (tidak bisa ditest murni dari browser tanpa device kedua/mode offline sungguhan), dan sebagai gantinya: baca kode sync_service.go untuk PushSync dan applySyncCashDrawer, verifikasi secara code-review bahwa update yang terjadi (products.stock, cash_drawer, transactions, expenses) memakai fungsi yang SAMA dengan yang sudah diperbaiki di migrasi timezone sebelumnya (time_helper.GetTimeNow() bukan time.Now() raw), dan bahwa tidak ada jalur yang melewati validasi konsistensi yang sudah diuji di R2-R4.
5. Kalau ada data konflik sync yang sudah ada (dari testing manual sebelumnya atau seed data), buka halaman resolve conflict, screenshot data "versi desktop" vs "versi online" yang ditampilkan, approve salah satu, lalu cek hasilnya konsisten di Transaksi/Produk/Kas sesuai domain yang terlibat.

Laporkan hasil lengkap dengan screenshot dan/atau catatan code-review. Kalau ada temuan bug, cari root cause di kode (jangan diperbaiki dulu). Jangan lanjut ke fase lain.
```

---

## FASE R9 — Users/Roles/Access/Menus

### Skenario yang harus diuji lewat browser
1. Buat role baru dengan akses terbatas (misal cuma bisa lihat Produk, tidak bisa Kasir).
2. Buat user baru dengan role itu.
3. Login sebagai user itu di browser terpisah (context baru), pastikan menu yang muncul di sidebar sesuai — menu yang tidak diizinkan TIDAK muncul.
4. Coba akses langsung via URL menu yang tidak diizinkan — harus diblokir (redirect/403), bukan cuma disembunyikan dari sidebar tapi tetap bisa diakses lewat URL.
5. Ganti akses role itu (kasih izin Kasir) tanpa logout user — cek apakah butuh re-login untuk berlaku, atau langsung update (dokumentasikan perilaku sebenarnya, ini bukan berarti bug kalau butuh re-login, asal jelas & konsisten).

### PROMPT EKSEKUSI — Fase R9
```
Baca docs/RELATIONAL_TESTING_PLAN.md di project d:\Develop\Project_pos_mahenz. Kerjakan FASE R9 (Users/Roles/Access/Menus) — testing relasional lintas menu via browser. Ini fase terakhir dari rencana testing relasional.

Konteks: access_service dan user_service sama-sama menginjeksi roleRepo. role_menu_access menentukan menu apa yang boleh diakses role tertentu, digerbang lewat ProtectedRoute menuKey di FE router.

TESTING WAJIB (lewat browser, gunakan 2 context/session terpisah -- satu admin, satu user baru):
1. Jalankan BE dan FE. Login sebagai admin/admin123 di context pertama.
2. Buat role baru "Role Terbatas" di /roles.
3. Di halaman Role Access, kasih role itu akses HANYA ke menu Produk (can_view saja), pastikan menu lain (Kasir, Kas Harian, dll) tidak dicentang.
4. Buat user baru dengan role "Role Terbatas" di /users, screenshot.
5. Buka context/session browser BARU (jangan pakai context admin), login sebagai user baru itu.
6. Screenshot sidebar -- WAJIB cuma menu Produk yang muncul, menu lain (Kasir, Kas Harian, dst) TIDAK ADA di sidebar.
7. Coba akses langsung via URL menu yang tidak diizinkan (misal ketik langsung /finance/cash-drawer di address bar). Screenshot hasilnya -- WAJIB diblokir (redirect ke halaman lain / pesan tidak berwenang), BUKAN berhasil menampilkan halaman itu.
8. Kembali ke context admin, kasih tambahan akses Kasir ke "Role Terbatas" di Role Access.
9. Di context user (TANPA logout/reload dulu), coba akses /pos (Kasir) -- catat apakah langsung bisa akses atau masih diblokir sampai reload/re-login. Lalu reload halaman di context user, cek lagi -- dokumentasikan perilaku sebenarnya (butuh reload atau tidak), ini untuk didokumentasikan bukan divonis bug kecuali perilakunya membingungkan/tidak konsisten.

Laporkan hasil lengkap dengan screenshot tiap langkah dari kedua context. Kalau ada temuan bug (terutama soal poin 7 -- akses URL langsung harus diblokir, ini celah keamanan kalau tidak), cari root cause di kode (jangan diperbaiki dulu, laporkan dengan detail). Ini fase terakhir rencana testing relasional.
```

---

## Ringkasan Setelah Semua Fase Selesai

Setelah R1-R9 selesai, susun satu laporan gabungan berisi:
- Daftar semua bug yang ditemukan (kalau ada), diurutkan dari paling kritikal (soal uang: Grup 8, 2, 3) ke paling ringan (Grup 7, 10)
- Untuk tiap bug: halaman mana vs halaman mana yang tidak konsisten, root cause di kode, dan rekomendasi perbaikan (belum dieksekusi — tunggu keputusan lanjut)
- Grup yang sudah terbukti konsisten (tidak ada temuan) juga dicatat, supaya jelas cakupan yang sudah teruji
