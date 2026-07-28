# Rencana Perbaikan UI/UX — POS Mahenz

Dokumen ini berisi rencana perbaikan desain UI/UX seluruh menu aplikasi (FE), dipecah menjadi fase-fase kecil — 1 fase = 1 submenu, dipecah lagi jadi beberapa sub-fase kalau submenunya punya banyak fitur/state. Setiap fase punya prompt siap-pakai — cukup rujuk fase & file ini ke asisten (mis. "Jalankan Fase 1.1 dari docs\UIUX_REDESIGN_PLAN.md"), tidak harus berurutan dalam satu sesi.

Fokus dokumen ini adalah **desain UI/UX** (konsistensi visual, alur interaksi, responsivitas), bukan debug fungsional — untuk itu sudah ada [DEBUG_TESTING_PLAN.md](DEBUG_TESTING_PLAN.md) yang jadi rekan sejalan (boleh dijalankan terpisah, tidak saling gantung). Kalau saat mengerjakan fase di sini ditemukan **bug** (bukan cuma soal tampilan — misal data salah, validasi bolong, alur rusak) yang berhubungan dengan FE maupun BE, langsung diperbaiki juga di fase yang sama, jangan ditunda.

---

## Semua dikerjakan lewat AI — tidak ada langkah manual di luar AI

Prinsip dokumen ini: **tidak ada kondisi yang mengharuskan user membuka browser sendiri, mengklik-klik manual, atau menilai tampilan dengan matanya sendiri di tengah proses.** Semua observasi (screenshot, log console, log network) diambil otomatis oleh AI lewat script, dibaca ulang oleh AI (lewat tool baca gambar), dan diperbaiki oleh AI. User hanya perlu: (1) menyetujui rencana/prompt tiap fase, (2) meninjau hasil akhir (ringkasan + screenshot before/after) yang dilaporkan AI, (3) approve kalau ada perubahan yang berisiko (mis. sentuh BE, migration, atau restore data).

**Tools yang dipakai AI dalam tiap fase:**
- **Playwright (Node.js, headless Chromium)** — dijalankan lewat script di folder scratchpad (`node nama-script.js`), untuk: membuka aplikasi, login, navigasi ke menu, isi form/klik tombol, dan mengambil screenshot di titik-titik penting serta di 3 breakpoint (`375px` mobile, `768px` tablet, `1280px` desktop).
- **Tool baca file/gambar bawaan AI** — untuk membuka & menilai screenshot yang dihasilkan Playwright (bukan cuma percaya teks log), supaya penilaian visual (spacing, overflow, kontras, keterpotongan) benar-benar diperiksa, bukan diasumsikan.
- **Dev server BE (`go run main.go`) & FE (`npm run dev`)** — dijalankan AI lewat tool shell, dibiarkan hidup selama fase berjalan, dimatikan AI di akhir fase.
- **Baca kode langsung (komponen, style, schema)** — untuk menelusuri akar masalah desain (mis. komponen shared yang dipakai banyak halaman) sebelum mengubah apa pun, bukan tempel-tempel CSS di titik gejala.
- **type-check / lint / build (FE)**, **go build / go vet (BE kalau disentuh)** — dijalankan AI di akhir tiap fase sebagai gerbang wajib sebelum fase dianggap selesai.

Tidak ada instruksi di dokumen ini yang meminta user melakukan sesuatu secara manual di luar AI (misal "coba buka browser dan lihat sendiri") — kalau nanti muncul kebutuhan seperti itu (misal verifikasi ke printer fisik/BLE sungguhan), itu akan disebut eksplisit sebagai pengecualian, bukan aturan umum.

---

## Prasyarat & Cara Menjalankan (wajib dibaca sebelum fase manapun)

Sama seperti [DEBUG_TESTING_PLAN.md](DEBUG_TESTING_PLAN.md) — ringkasnya:

1. **Jalankan BE & FE** (AI yang menjalankan lewat tool shell, bukan user):
   ```
   # Terminal 1 — Backend
   cd BE
   go run main.go
   # → API di http://localhost:8080/api, health check di http://localhost:8080/health

   # Terminal 2 — Frontend
   cd FE
   npm run dev
   # → aplikasi di http://localhost:3000 (port aktual ikuti output npm run dev)
   ```
2. **Database & migration**: otomatis jalan sendiri saat BE start (lihat detail di DEBUG_TESTING_PLAN.md poin 2). Data dummy hasil eksplorasi tidak perlu dihapus.
3. **Kredensial**: Owner `owner`/`owner123`, Admin `admin`/`admin123`, Kasir belum ada user seeded — buat dulu lewat Manajemen User kalau fase butuh role Kasir.
4. **Playwright**: kalau belum pernah dipakai di lingkungan kerja saat itu, install dulu sekali (`npm init -y && npm install playwright && npx playwright install chromium`) di folder scratchpad — bukan per fase.

---

## Aturan Umum (berlaku di SETIAP fase)

**Mindset:** bertindak sebagai UI/UX reviewer + front-end engineer, bukan cuma "menjalankan fitur sekali lalu dianggap bagus". Untuk tiap halaman/komponen di submenu yang sedang dikerjakan, periksa sistematis poin-poin berikut (catat juga yang sudah oke, jangan cuma catat yang bermasalah):

**1. Konsistensi visual (dibanding komponen shared & halaman lain yang sudah baik):**
- Spacing, ukuran font, warna, radius, shadow konsisten dengan komponen shared (`FE/src/shared/components/ui/*`, `DataTable`, `FormModal`/`ActionModal`, `AppLayout`) — bukan style ad-hoc yang menyimpang tanpa alasan.
- Ikon, label tombol, dan istilah (bahasa Indonesia) konsisten dengan menu lain yang sejenis (mis. semua tombol hapus pakai label & warna yang sama).
- Empty state, loading state, dan error state ada dan konsisten gaya-nya dengan halaman lain (bukan blank putih atau teks default browser).

**2. Alur interaksi (UX):**
- Alur logis: langkah wajar user bisa diselesaikan tanpa jalan buntu, tanpa bolak-balik yang tidak perlu.
- Feedback jelas untuk tiap aksi: loading indicator saat proses, toast/notifikasi sukses-gagal, disabled state yang masuk akal (bukan tombol yang kelihatan aktif tapi tidak ngapa-ngapain).
- Form: label & placeholder jelas, pesan validasi muncul dekat field yang salah (bukan cuma toast generik), urutan tab/fokus wajar.
- Konfirmasi untuk aksi destruktif (hapus, void, dsb) — konsisten dengan pola di halaman lain.

**3. Responsivitas (WAJIB di 375px, 768px, dan 1280px untuk setiap halaman di fase ini):**
- Tidak ada horizontal overflow di level dokumen (`document.documentElement.scrollWidth - clientWidth`, cek lewat script, bukan sekilas).
- `DataTable` beralih ke card-view yang benar di <768px (field judul & yang disembunyikan sesuai `mobileLabel`/`mobileHidden` di `*TableColumns.tsx`).
- Modal/dialog tidak overflow ke luar layar, tombol aksi tetap terjangkau (tidak ketutup keyboard virtual/footer).
- Sidebar drawer & elemen navigasi berfungsi normal di semua breakpoint.

**4. Aksesibilitas dasar:**
- Kontras warna teks-background cukup (terutama teks abu-abu di atas latar terang/gelap).
- Elemen interaktif (tombol, link) punya target tap yang cukup besar di mobile, tidak berdempetan.
- Ikon-only button punya `title`/`aria-label` supaya jelas fungsinya.

**5. Bug non-visual yang kebetulan ketemu saat eksplorasi UI:**
- Kalau saat mengetes alur ditemukan data salah, kalkulasi keliru, validasi bolong, atau error console/network — ini **bukan** cakupan utama fase ini (itu domain DEBUG_TESTING_PLAN.md), TAPI tetap wajib diperbaiki di fase yang sama sesuai instruksi user, baik di FE maupun BE.

**Cara kerja tiap fase (semua lewat AI):**
1. AI baca kode halaman/komponen terkait dulu (struktur, style yang dipakai, komponen shared yang direuse) sebelum menilai — supaya rekomendasi ubah selaras dengan pola yang sudah ada, bukan bikin pola baru.
2. AI tulis script Playwright: login sesuai role yang relevan untuk submenu itu, navigasi ke halaman, jalankan skenario UI utama (termasuk buka form/modal, isi data, state loading/error kalau bisa disimulasikan), ambil screenshot di 375px/768px/1280px di titik-titik penting, rekam console error & response API yang gagal.
3. AI jalankan script, baca semua screenshot (tool baca gambar) untuk menilai poin 1-4 di atas, dan baca log untuk poin 5.
4. AI susun daftar temuan (per kategori: visual/UX/responsive/aksesibilitas/bug), lalu **perbaiki langsung** — perbaikan seminimal & setepat mungkin menyasar akar masalah, prioritaskan reuse komponen shared yang sudah ada daripada bikin style baru.
5. Setelah tiap perbaikan (atau kumpulan perbaikan kecil yang berkaitan), AI ulangi screenshot untuk verifikasi before/after — jangan cuma percaya bahwa perbaikan kode "pasti benar" tanpa verifikasi visual ulang.
6. Kalau ada perbaikan di komponen shared yang dipakai banyak halaman (`DataTable`, `AppLayout`, `dialog.tsx`, dst), AI wajib spot-check minimal 1-2 halaman lain di luar submenu fase ini untuk pastikan tidak regresi.
7. Setelah semua perbaikan di fase itu selesai, jalankan gerbang wajib:
   ```
   npm run type-check
   npm run lint
   npm run build
   ```
   Ketiganya harus **0 error dan 0 warning**. Kalau ada file BE yang ikut diubah (karena ditemukan bug terkait), tambahkan juga `go build ./...` dan `go vet ./...` di folder `BE` sampai bersih.
8. Laporan akhir tiap fase (ditulis AI, bukan diminta dari user): daftar perubahan UI/UX yang dilakukan (sebelum → sesudah, alasan), daftar bug non-visual yang ditemukan & diperbaiki (kalau ada), dan daftar hal yang sudah diperiksa tapi ternyata sudah baik (supaya kelihatan cakupannya, bukan cuma yang diubah).
9. Data uji yang dibuat selama eksplorasi (kalau ada) **tidak perlu dihapus**, biarkan di database development.

---

## Daftar Fase

Urutan fase sama dengan pengelompokan menu di aplikasi (9 grup, 28 submenu). Fase tidak harus berurutan, tapi disarankan ikut nomor karena beberapa halaman saling reuse komponen (kalau Fase 1 memperbaiki `DataTable`, fase-fase berikutnya otomatis ikut lebih rapi).

### Fase 1 — Beranda: Dashboard
```
Jalankan Fase 1 (Beranda — Dashboard) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 1: review & perbaiki UI/UX halaman Dashboard (FE/src/features/dashboard/DashboardPage.tsx) untuk role Kasir, Admin, dan Owner (tampilannya berbeda per role — cek ketiganya). Ikuti Aturan Umum & Cara kerja tiap fase di docs\UIUX_REDESIGN_PLAN.md secara penuh (baca kode dulu, Playwright screenshot 375/768/1280px, nilai 5 kategori, perbaiki, verifikasi ulang, gerbang type-check+lint+build). Ini fase pertama yang menyentuh komponen shared (AppLayout/Sidebar/DataTable) — kalau ditemukan perbaikan di level shared, catat dampaknya ke fase-fase berikutnya di laporan akhir.
```

### Fase 2 — Penjualan: Kasir
```
Jalankan Fase 2 (Penjualan — Kasir) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 2: review & perbaiki UI/UX halaman Kasir (pencarian produk, keranjang, checkout, modal pembayaran per metode, struk). Ini halaman paling sering dipakai (kritis untuk kecepatan kerja kasir) — perhatikan khusus: kepadatan informasi vs kecepatan tap di mobile/tablet (perangkat kasir riil sering tablet), kontras angka total/kembalian (harus paling menonjol di layar), dan transisi antar tab Produk/Keranjang di mobile. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 3 — Penjualan: Transaksi
```
Jalankan Fase 3 (Penjualan — Transaksi) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 3: review & perbaiki UI/UX halaman Transaksi (list riwayat, filter tanggal/kasir/metode, detail transaksi, aksi void, cetak ulang). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 4 — Produk: Produk
```
Jalankan Fase 4 (Produk — Produk) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 4: review & perbaiki UI/UX halaman Produk — list (search, filter, DataTable/card-view), form Tambah/Edit (termasuk bagian satuan/paket dinamis, batch expired, gating stok per role), modal Import, modal Cetak Label, badge warning expired. Ini halaman dengan form paling kompleks di aplikasi — perhatikan khusus keterbacaan form panjang di mobile (grouping section, jangan satu scroll panjang tanpa struktur) dan kejelasan indikator readonly (field stok untuk Owner). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 5 — Produk: Kategori
```
Jalankan Fase 5 (Produk — Kategori) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 5: review & perbaiki UI/UX halaman Kategori Produk (list, CRUD modal). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 6 — Produk: Unit
```
Jalankan Fase 6 (Produk — Unit) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 6: review & perbaiki UI/UX halaman Unit/Satuan (list, CRUD modal). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 7 — Pengadaan: Supplier
```
Jalankan Fase 7 (Pengadaan — Supplier) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 7: review & perbaiki UI/UX halaman Supplier (list, search/filter, form tambah/edit, toggle aktif). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 8 — Pengadaan: Pembelian

**8.1 List, filter & create**
```
Jalankan Fase 8.1 (Pengadaan — Pembelian: List, filter & create) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 8.1: review & perbaiki UI/UX list Pembelian Supplier (filter tanggal/supplier/status, sorting, pagination, empty state) dan form Create PO (multi-item, kalkulasi total real-time, bagian batch expired). Perhatikan khusus kejelasan status pembayaran (badge warna) dan kepadatan tabel item di mobile. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

**8.2 Detail, edit, bayar, void, tambah item**
```
Jalankan Fase 8.2 (Pengadaan — Pembelian: Detail, edit, bayar, void, tambah item) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 8.2: review & perbaiki UI/UX halaman Detail PO beserta modal Edit, Bayar, Void, dan Tambah Item. Perhatikan khusus kejelasan banner info tier pembayaran (unpaid/partial/paid — field mana yang terkunci harus jelas terlihat, bukan cuma disabled diam-diam) dan histori badge ganda (status void + payment_status). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 9 — Pengadaan: Retur
```
Jalankan Fase 9 (Pengadaan — Retur) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 9: review & perbaiki UI/UX halaman Retur Pembelian (list, form create dengan pilih PO & checklist item, alur approve/reject). Perhatikan khusus kejelasan state setelah ganti PO yang dipilih (checklist harus jelas ter-reset). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 10 — Pelanggan: Pelanggan
```
Jalankan Fase 10 (Pelanggan — Pelanggan) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 10: review & perbaiki UI/UX halaman Pelanggan (list, search, form tambah/edit termasuk field Limit Kredit). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 11 — Pelanggan: Piutang
```
Jalankan Fase 11 (Pelanggan — Piutang) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 11: review & perbaiki UI/UX halaman Piutang (list, filter status, modal pelunasan sebagian/lunas). Perhatikan khusus kejelasan sisa piutang vs sudah dibayar. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 12 — Keuangan: Dashboard Keuangan
```
Jalankan Fase 12 (Keuangan — Dashboard Keuangan) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 12: review & perbaiki UI/UX halaman Dashboard Keuangan (filter tanggal manual & preset, 4 kartu ringkasan, tabel arus kas). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 13 — Keuangan: Kas Harian & Kas Saya
```
Jalankan Fase 13 (Keuangan — Kas Harian & Kas Saya) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 13: review & perbaiki UI/UX halaman Kas Harian (rekap semua kasir) dan Kas Saya (termasuk blok guard "Belum bisa membuka kas" & form Buka/Tutup Kas). Pastikan konsisten dengan CashDrawerStatusCard di Dashboard (Fase 1) karena reuse komponen yang sama. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 14 — Keuangan: Pengeluaran
```
Jalankan Fase 14 (Keuangan — Pengeluaran) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 14: review & perbaiki UI/UX halaman Pengeluaran (list, form tambah/edit dengan kategori, hapus). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 15 — Pelaporan: Penjualan
```
Jalankan Fase 15 (Pelaporan — Penjualan) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 15: review & perbaiki UI/UX Laporan Penjualan (filter, tabel/grafik, export kalau ada). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 16 — Pelaporan: Laba Rugi
```
Jalankan Fase 16 (Pelaporan — Laba Rugi) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 16: review & perbaiki UI/UX Laporan Laba Rugi (breakdown pendapatan/HPP/pengeluaran/laba). Perhatikan khusus keterbacaan angka besar & struktur breakdown bertingkat di mobile. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 17 — Pelaporan: Stok
```
Jalankan Fase 17 (Pelaporan — Stok) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 17: review & perbaiki UI/UX Laporan Stok (tabel mutasi stok, filter). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 18 — Pelaporan: Kinerja Kasir
```
Jalankan Fase 18 (Pelaporan — Kinerja Kasir) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 18: review & perbaiki UI/UX Laporan Kinerja Kasir (filter per kasir/tanggal, tabel/grafik perbandingan). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 19 — Pelaporan: Ringkasan Bisnis
```
Jalankan Fase 19 (Pelaporan — Ringkasan Bisnis) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 19: review & perbaiki UI/UX halaman Ringkasan Bisnis (SummaryCards, SalesChart, TopProductsTable, filter periode). Perhatikan khusus keterbacaan grafik tren di layar kecil (label sumbu, tooltip tidak ketutup). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 20 — Operasional: Shift
```
Jalankan Fase 20 (Operasional — Shift) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 20: review & perbaiki UI/UX halaman Manajemen Shift (list, CRUD jam mulai/selesai & nama shift, assign ke user). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 21 — Operasional: Sync Center
```
Jalankan Fase 21 (Operasional — Sync Center) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 21: review & perbaiki UI/UX halaman Sync Center (status sinkronisasi tiap modul, indikator progress, penanganan konflik kalau ada UI-nya). Perhatikan khusus kejelasan status (berhasil/gagal/pending) — ini modul teknis yang gampang membingungkan kalau statusnya ambigu. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 22 — Sistem: Profil Toko
```
Jalankan Fase 22 (Sistem — Profil Toko) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 22: review & perbaiki UI/UX halaman Profil Toko (form edit nama/alamat/kontak/logo, preview logo). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 23 — Sistem: Manajemen User
```
Jalankan Fase 23 (Sistem — Manajemen User) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 23: review & perbaiki UI/UX halaman Manajemen User (list, CRUD, assign role, toggle aktif/nonaktif). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 24 — Sistem: Manajemen Role & Role Access
```
Jalankan Fase 24 (Sistem — Manajemen Role & Role Access) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 24: review & perbaiki UI/UX halaman Manajemen Role (CRUD role) dan Role Access (matrix permission can_view/can_create/can_edit/can_delete per menu). Perhatikan khusus keterbacaan matrix permission yang biasanya lebar (banyak kolom) — cek strategi tampilan di mobile (scroll horizontal terkontrol vs restrukturisasi jadi per-menu). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 25 — Sistem: Manajemen Menu
```
Jalankan Fase 25 (Sistem — Manajemen Menu) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 25: review & perbaiki UI/UX halaman Manajemen Menu (tambah/edit menu, pilih parent & path dari route registry, icon picker, drag/ubah urutan). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 26 — Sistem: Pengaturan Printer
```
Jalankan Fase 26 (Sistem — Pengaturan Printer) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 26: review & perbaiki UI/UX halaman Pengaturan Printer (form ukuran kertas/header/footer/logo/auto-print, preview struk live). Preview struk & tombol "Test Print" cukup diverifikasi dialog print browser muncul dengan konten benar (verifikasi cetak fisik nyata di luar cakupan fase ini). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 27 — Sistem: Versi Aplikasi
```
Jalankan Fase 27 (Sistem — Versi Aplikasi) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 27: review & perbaiki UI/UX halaman Versi Aplikasi (list versi, form tambah, toggle wajib update). Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 28 — Sistem: Backup & Restore
```
Jalankan Fase 28 (Sistem — Backup & Restore) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 28: review & perbaiki UI/UX halaman Backup & Restore (tombol backup manual, list riwayat backup, form restore). HATI-HATI: jangan benar-benar menjalankan aksi restore saat eksplorasi UI (cukup review tampilan form/dialog konfirmasinya) karena berpotensi menimpa data development — kalau perlu verifikasi alur restore sungguhan, konfirmasi dulu ke user secara eksplisit sebelum menjalankan. Ikuti Aturan Umum & Cara kerja di docs\UIUX_REDESIGN_PLAN.md.
```

### Fase 29 — Regresi Visual Akhir
```
Jalankan Fase 29 (Regresi Visual Akhir) dari docs\UIUX_REDESIGN_PLAN.md :

Jalankan Fase 29: regresi visual menyeluruh. Ambil screenshot smoke-test 375px/768px/1280px untuk seluruh 28 submenu (Fase 1-28) dalam satu sesi, bandingkan konsistensi lintas halaman (spacing, warna, komponen shared) sekarang setelah semua fase selesai — pastikan tidak ada halaman yang "ketinggalan gaya" dibanding yang lain. Jalankan type-check, lint, build (FE) dan go build+go vet (BE kalau ada perubahan BE sepanjang seluruh fase) final — harus 0 error 0 warning. Susun laporan ringkasan akhir: daftar seluruh perubahan UI/UX per fase, daftar bug non-visual yang ikut diperbaiki, dan rekomendasi (kalau ada) untuk perbaikan lanjutan di luar cakupan 28 fase ini.
```

---

## Catatan Penggunaan

- Fase tidak harus dijalankan berurutan, tapi urutan di atas disarankan karena Fase 1 (Dashboard) adalah fase pertama yang menyentuh komponen shared (`AppLayout`, `Sidebar`, `DataTable`) — perbaikan di situ otomatis berdampak positif ke fase-fase berikutnya.
- Kalau satu fase menemukan banyak isu dan terasa kepanjangan, boleh dihentikan di tengah dan dilanjutkan nanti — sebutkan fase & bagian mana yang terakhir dikerjakan.
- Dokumen ini saudara dari [DEBUG_TESTING_PLAN.md](DEBUG_TESTING_PLAN.md) (fokus fungsional/bug) — keduanya boleh dicicil bergantian per submenu yang sama tanpa saling mengganggu, karena aturan "data uji tidak perlu dihapus" berlaku di keduanya.
- Screenshot yang dihasilkan Playwright selama proses boleh disimpan sementara di folder scratchpad — tidak perlu dirapikan/disimpan permanen di repo kecuali user memintanya secara eksplisit untuk dokumentasi.
