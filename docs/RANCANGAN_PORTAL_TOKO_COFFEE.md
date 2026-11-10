# Rancangan Portal + POS Toko + POS Coffee

> Status: **RANCANGAN / CATATAN DISKUSI - belum ada implementasi.**
> Dokumen ini adalah catatan hidup sebelum implementasi. Setiap kali ada diskusi atau
> keputusan baru, dokumen ini diperbarui (lihat "Riwayat Pembaruan" di bagian akhir).
>
> Label status yang dipakai:
> - **[SEPAKAT]** sudah diputuskan bersama.
> - **[USULAN]** usulan AI, belum dikonfirmasi user.
> - **[TERBUKA]** belum dibahas / belum diputuskan.

---

## 1. Tujuan

Menambah dua aplikasi baru di dalam satu repo (folder root `Project_pos_mahenz`):

1. **Portal** - pintu login terpusat. Setelah login, user memilih aplikasi yang dituju.
2. **POS Coffee** - aplikasi kasir kedai kopi. BE dan FE distrukturkan sama dengan POS Toko,
   yang berbeda hanya fiturnya.

Aplikasi yang sudah ada (POS retail) disebut **POS Toko** dan dipindah ke folder `TOKO/`.

---

## 2. Ringkasan Keputusan

| # | Topik | Keputusan | Status |
|---|-------|-----------|--------|
| 1 | Struktur folder | `PORTAL/`, `TOKO/`, `COFFEE/`, masing-masing berisi `BE/` dan `FE/` | SEPAKAT |
| 2 | BE/FE yang ada sekarang | Dipindah ke `TOKO/BE` dan `TOKO/FE` | SEPAKAT |
| 3 | `mobile/`, `nginx/`, `docs/`, `.github/` | Tetap di root | SEPAKAT |
| 4 | Database | Terpisah per aplikasi (lihat bagian 4) | SEPAKAT |
| 5 | Urutan kerja | Pindah ke TOKO, lalu Portal, lalu hubungkan Portal ke Toko, lalu Coffee | SEPAKAT |
| 6 | Layout URL | Berbasis **path**, Portal di root domain | SEPAKAT |
| 7 | Stack Portal | Sama dengan Toko (Go/Gin/GORM + React/Vite/Tailwind); struktur sama, beda desain UI/UX | SEPAKAT |
| 8 | Login | Terpusat di Portal (SSO kode sekali pakai) | SEPAKAT |
| 9 | Pencocokan aplikasi | Pakai `app_id` (tabel aplikasi di Portal) dicocokkan dengan `app_id` di tabel `portal_credentials` pada database aplikasi (sebelumnya `APP_ID` di `.env`, diubah, lihat 6.10) | SEPAKAT |
| 10 | Akses user per aplikasi | Diatur di Portal (user boleh masuk aplikasi mana saja) | SEPAKAT |
| 11 | Akun lama | Disalin ke Portal dan ditautkan, bukan dipindah | SEPAKAT |
| 12 | Tabel user di tiap aplikasi | Tetap ada (Toko dan Coffee punya `users` sendiri) | SEPAKAT |
| 13 | Role default user baru di aplikasi | Role baru "menunggu akses" tanpa izin menu | SEPAKAT |
| 14 | Login langsung ke aplikasi | Tetap ada, dengan toggle on/off; hanya role admin yang boleh mengubah | SEPAKAT |
| 15 | Sesi | Satu sesi per user; token 8 jam, sesi 7 hari (diperpanjang selama dipakai) | SEPAKAT |
| 16 | Role di Portal | Hanya `admin` dan `user` | SEPAKAT |
| 17 | Halaman pilih aplikasi | Selalu tampil, walau hanya ada satu pilihan | SEPAKAT |
| 18 | Aturan logout | Dari aplikasi: akhiri sesi aplikasi lalu kembali ke Portal. Dari Portal: keluar Portal saja | SEPAKAT |
| 19 | Fitur Portal versi awal | Login, daftar aplikasi sesuai akses, kelola user dan akses, kelola aplikasi (UI, khusus admin), ganti password, logout | SEPAKAT |
| 20 | Penamaan key penyimpanan browser | Format `<aplikasi>.<nama>` termasuk Toko (lihat bagian 9) | SEPAKAT |
| 21 | APK tablet | Satu APK untuk semua aplikasi, tidak diubah, alamat awal = Portal | SEPAKAT |
| 22 | Pendekatan Coffee | Salin BE dan FE Toko dulu, ubah bertahap; menu Kasir dan Pembelian diubah lebih dulu | SEPAKAT |
| 23 | Menu Pembelian di Coffee | Pembelian dari supplier untuk dua jenis: bahan baku, dan barang jadi (menu bermode potong stok langsung, mis. roti dan snack) | SEPAKAT |
| 24 | Stok bahan baku | Model Lentur dengan pencatatan pemakaian manual: stok sisa = pembelian dikurangi pemakaian yang diinput, dikoreksi lewat opname. Resep hanya catatan takaran standar dan tidak mengurangi stok sisa | SEPAKAT |
| 25 | Kasir Coffee versi pertama | Varian, tambahan, nomor meja, dan nomor antrean | SEPAKAT |
| 26 | Modul Toko yang dipangkas di Coffee | `supplier_return`, `expiry_batch`, `customer_balance`, `capital_reconciliation`, `backdate`, `version` (`receivable` dikeluarkan dari daftar pangkas karena sistem utang dibawa, lihat 14.10) | SEPAKAT |
| 27 | Waktu pembayaran pesanan | Dua opsi: bayar sekarang, atau pesan dulu bayar belakangan (open bill) | SEPAKAT |
| 28 | Pesanan terbuka saat shift ditutup | Boleh dilanjutkan ke shift berikutnya dalam hari yang sama; tidak boleh lintas hari | SEPAKAT |
| 29 | Open bill lanjutan | Gabung bill dan pindah meja masuk versi pertama; split bill dihapus dari rencana | SEPAKAT |
| 30 | Pesanan yang belum dibayar saat batas hari | Menjadi utang seperti sistem utang di Toko (modul `receivable` dan `customer` dibawa ke Coffee); aturan utang sama seperti Toko apa adanya | SEPAKAT |
| 31 | Pencatatan pemakaian stok | Kapan saja, satu bahan per input; dicatat sebagai "stok terpakai" saat stok mulai dipakai | SEPAKAT |
| 32 | Peran resep di versi pertama | Catatan takaran standar per menu (data resep sudah ada); perkiraan pemakaian dan HPP menyusul | SEPAKAT |
| 33 | Mode stok per menu | Dua pilihan per menu: "potong stok langsung" (barang jadi dari supplier, mis. roti dan snack tertentu) dan "manual / racikan" (dibuat sendiri; stok bahan dikurangi lewat pemakaian manual). Tidak semua makanan memakai mode langsung | SEPAKAT |
| 34 | Tutup kas dan pesanan terbuka | Opsi A: tutup kas manual diblokir sampai pesanan terbuka diselesaikan (bayar, jadikan utang dengan memilih pelanggan, atau batalkan). Penutupan otomatis tengah malam: opsi i, lihat keputusan 42 | SEPAKAT |
| 35 | Batas hari | Mengikuti Toko (tengah malam, zona Asia/Jakarta) | SEPAKAT |
| 36 | Aturan utang | Sama seperti Toko apa adanya | SEPAKAT |
| 37 | Pemakaian stok | Jumlah boleh desimal | SEPAKAT |
| 38 | Data menu dan resep | Tanpa fitur impor; dicatat manual, resep diisi dari form menu | SEPAKAT |
| 39 | Satuan barang jadi | Konsep satuan bertingkat seperti Toko (dus dan pcs, konversi tetap). Berbeda dari bahan baku (satu satuan stok, tanpa konversi tetap) | SEPAKAT |
| 40 | Waktu potong stok barang jadi | Saat item masuk ke pesanan; stok dikembalikan bila item dihapus atau pesanan dibatalkan (mengikuti rekomendasi AI) | SEPAKAT |
| 41 | Stok barang jadi habis | Penjualan diblokir; ketidaksesuaian stok dikoreksi lewat opname (mengikuti rekomendasi AI) | SEPAKAT |
| 42 | Penutupan otomatis tengah malam dengan pesanan terbuka | Opsi i: kas yang masih punya pesanan terbuka tidak ditutup otomatis dan ditandai untuk diselesaikan | SEPAKAT |
| 43 | Varian barang jadi | Satu menu = satu stok, tanpa varian; rasa berbeda = menu terpisah. Tambahan boleh, tidak memengaruhi stok | SEPAKAT |
| 44 | Pembatalan setelah barang diserahkan | Kasir memilih "kembali ke stok" atau "dibuang" | SEPAKAT |
| 45 | Aturan pesanan melewati batas hari (pendamping opsi i) | Hanya bisa diselesaikan, tanpa tambah item; buka kas baru diblokir sampai kas tertunda ditutup; admin bisa menyelesaikan atas nama; ada daftar pesanan tertunda | SEPAKAT |
| 49 | User baru buatan admin | Wajib ganti password saat login pertama (sama seperti hasil reset); diserahkan user ke AI | SEPAKAT |
| 50 | Urutan tahap 2.4/2.5 | Tahap 2.5 (BE aplikasi dan SSO) dikerjakan lebih dulu, baru 2.4 (halaman pilih aplikasi) agar memakai data asli; user menyerahkan urutan ke rekomendasi AI | SEPAKAT |
| 46 | Cara membangun Portal | Usul user: struktur dan konsistensi kode BE dan FE Portal wajib sama dengan Toko; caranya menyalin `TOKO/BE` dan `TOKO/FE` ke `PORTAL/`, menghapus fitur yang tidak relevan, dan desain UI/UX diubah total. Implementasi dibuat bertahap | SEPAKAT (usulan user, ditanggapi AI) |
| 47 | Desain Portal | Diserahkan ke AI (bebas); target tampilan modern dan canggih | SEPAKAT |
| 51 | Siapa boleh mengubah toggle login langsung | Hanya role `admin` (di Toko juga: `owner` tidak boleh) | SEPAKAT |
| 48 | Skrip salin user Toko ke Portal | Dijalankan sebagai perintah terpisah (bukan otomatis saat Portal hidup); user menyerahkan keputusan ke AI | SEPAKAT |

---

## 3. Struktur Folder

```
Project_pos_mahenz/
├── PORTAL/
│   ├── BE/        (baru)
│   └── FE/        (baru)
├── TOKO/
│   ├── BE/        (pindahan dari BE/)
│   └── FE/        (pindahan dari FE/)
├── COFFEE/
│   ├── BE/        (baru, salinan struktur Toko)
│   └── FE/        (baru, salinan struktur Toko)
├── mobile/        (APK Capacitor, tetap)
├── nginx/         (tetap)
├── docs/          (tetap)
└── .github/       (tetap)
```

Prinsip: tiap aplikasi berdiri sendiri (`go.mod` dan `package.json` masing-masing), tanpa
workspace di root.

**[USULAN]** Detail tambahan:
- Nama module Go aplikasi baru diganti: `portal_api`, `coffee_api` (Toko tetap `pos_api`).
- Nama package FE aplikasi baru diganti (mis. `portal-web`, `coffee-web`).
- Port BE (diperbarui, lihat 10.1): Portal 8080, Toko 8081, Coffee 8082. Fitur versi lama/baru (Mode B) dihapus.
- Port dev FE (diperbarui, sejajar dengan port BE): Portal 3000, Toko 3001, Coffee 3002 (agar bisa jalan bersamaan di laptop).

---

## 4. Database

| Aplikasi | Database | Catatan |
|----------|----------|---------|
| Toko | `pos_retail_db` | Nama tidak diubah (menghindari risiko di production) |
| Coffee | `pos_coffee_db` | Baru |
| Portal | `pos_portal_db` | Baru |

- **[USULAN]** Tiap aplikasi memakai user MySQL sendiri (least privilege, sesuai pola di
  `DEPLOYMENT_PROD.md`).
- BE di repo ini migrasi otomatis dari `database/migrations/` saat start, jadi database baru
  terbentuk sendiri begitu BE baru pertama kali dijalankan.

---

## 5. Urutan Kerja

| Fase | Pekerjaan | Catatan |
|------|-----------|---------|
| 1 | Pindahkan BE/FE ke `TOKO/` | **[KODE SELESAI, belum di-commit dan belum di-deploy]** Tanpa perubahan fungsi, tetapi port Toko menjadi 8081 dan service menjadi `pos-toko-backend`. Deploy ke VPS dan verifikasi POS Toko jalan sama seperti sebelumnya |
| 2 | Bangun Portal | BE + FE + database baru |
| 3 | Hubungkan Portal ke Toko | SSO, login langsung on/off, penyeragaman key, penyesuaian URL |
| 4 | Bangun Coffee | Disalin dari `TOKO/BE` dan `TOKO/FE` di lokasi barunya, lalu fitur kopi |

**[USULAN]** Titik ubah URL production (Toko pindah ke `/toko`, Portal di root) terjadi **sekali**
saat go-live Fase 3, bukan di Fase 1. Dengan begitu Fase 1 benar-benar tidak mengubah apa pun
yang dilihat kasir.

### 5.0 Tahapan Fase 2 (Portal), dikerjakan bertahap

Cara membangun: salin `TOKO/BE` dan `TOKO/FE` ke `PORTAL/BE` dan `PORTAL/FE`, pangkas fitur retail, ganti
identitas (module Go `portal_api`, package FE `portal-web`, database `pos_portal_db`, `SecretKey` JWT berbeda,
port BE 8080 dan FE 3000, key browser `portal.*`). Struktur dan gaya kode wajib sama dengan Toko; desain UI/UX
dibuat baru dan modern. Setiap tahap berhenti dan menunggu persetujuan user sebelum lanjut.

| Tahap | Pekerjaan | Hasil yang dicek user |
|-------|-----------|------------------------|
| 2.1 **[SELESAI, menunggu persetujuan]** | Kerangka BE Portal: salin `TOKO/BE`, pangkas, ganti identitas, migrasi awal ke `pos_portal_db` | BE hidup di 8080, `/api/health` membalas, tabel terbentuk |
| 2.2 **[SELESAI, menunggu persetujuan]** | Auth BE: login, sesi, user, role `admin`/`user`, seed akun awal | Login dan refresh/logout lewat API |
| 2.3 **[SELESAI, menunggu persetujuan]** | Kerangka FE Portal: salin `TOKO/FE`, pangkas, sistem desain baru, halaman login | Login berjalan di 3000 dengan tampilan baru |
| 2.4 **[SELESAI, menunggu persetujuan]** | Halaman pilih aplikasi: kartu aplikasi, menu pengguna | Daftar aplikasi tampil setelah login. **Dikerjakan setelah 2.5** (urutan ditukar agar memakai data asli, bukan daftar sementara) |
| 2.5 **[SELESAI, menunggu persetujuan]** | Aplikasi dan akses (BE): tabel `apps`, `user_apps`, `sso_codes`, kunci, `/sso/exchange` | Kode SSO dibuat dan ditukar lewat API |
| 2.6 | Admin Portal: kelola user, kelola aplikasi, atur akses, reset password | Admin mengatur semuanya lewat UI. **Dibagi: 2.6a BE [SELESAI], 2.6b FE [SELESAI, menunggu persetujuan]** |
| 2.7 | Ganti password sendiri dan toggle login langsung | Fitur akun berjalan |

Keputusan teknis tahap 2.1-2.2:
- Database `pos_portal_db` sudah disiapkan user; migrasi dibuat oleh AI.
- Akun awal Portal sama dengan seed Toko (`002_seed_data.sql`): `owner` dan `admin` dengan hash Argon2id yang
  sama, sehingga password awal sama dengan Toko. Role Portal hanya `admin` dan `user`; akun seed Toko dipetakan:
  `owner` dan `admin` menjadi role `admin` Portal [USULAN].

### 5.0b Tahapan Fase 3 (hubungkan Portal ke Toko), dikerjakan bertahap

Tahap 2.7 (toggle login langsung) digabung ke Fase 3 (disetujui user) karena penegakannya ada di BE Toko.
Perubahan di Toko bersifat **aditif**: login lama tetap berfungsi sampai toggle dimatikan.

| Tahap | Pekerjaan | Hasil yang dicek |
|-------|-----------|------------------|
| 3.1 **[SELESAI, menunggu persetujuan]** | BE Toko: migrasi (`users.portal_user_id`, role "menunggu akses", tabel `portal_credentials`, setting `direct_login_enabled`), config alamat Portal, `POST /api/auth/sso` (tukar kode ke Portal, cari/buat user lokal, buat sesi) | Kode asli dari Portal menghasilkan sesi Toko lewat API |
| 3.2 **[SELESAI, menunggu persetujuan]** | BE Toko: toggle login langsung (endpoint status publik, login ditolak saat off, endpoint ubah untuk admin, flag darurat di config) | Toggle bekerja, flag darurat mengalahkan setelan DB |
| 3.3 | FE Toko: migrasi key browser ke `toko.*`, halaman `/sso`, tombol "Masuk lewat Portal" dan tombol kembali ke Portal, form login mengikuti toggle, UI toggle di Pengaturan, base path `/toko/` (lewat env) | Alur Portal ke Toko berjalan di browser (dev) |
| 3.4 | Skrip sekali jalan: salin user Toko ke Portal dan tautkan `portal_user_id` (dry-run, backup) | Dry-run dan eksekusi pada DB dev |
| 3.5 | Nginx (path `/`, `/portal-api/`, `/toko/`, `/toko/api/`), SOP deploy, panduan cutover VPS | Konfigurasi dan panduan siap, divalidasi `nginx -t` bila tersedia |
| 3.6 | Uji ujung ke ujung di dev (Portal 3000/8080 ke Toko 3001/8081) | Skenario lengkap lolos |

### 5.1 Dampak Fase 1 (pemindahan folder)

Kode tidak perlu diubah: module Go `pos_api` diimpor lewat nama module (bukan path folder), dan
alias `@` di Vite relatif ke folder FE-nya sendiri. Yang berubah hanya hal di sekelilingnya:

1. **`.gitignore`**: semua pola berawalan `BE/` dan `FE/` tidak akan cocok lagi (risiko
   `node_modules`, `dist`, log, binary ikut ter-commit). **[USULAN]** ganti ke pola generik
   (`**/node_modules/`, `**/dist/`, dst.) agar tahan perubahan struktur.
2. **Dokumen** yang menyebut path lama: `README.md`, `docs/DEPLOYMENT_PROD.md`,
   `docs/SOP_DEPLOY.md`, `mobile/README.md`, dan dokumen lain yang menyebut `BE/` atau `FE/`.
3. **VPS** (berdasarkan `SOP_DEPLOY.md` dan `DEPLOYMENT_PROD.md`):
   - Repo di `/opt/pos-mahenz`; systemd `pos-backend` memakai `WorkingDirectory=/opt/pos-mahenz/BE`
     dan `ExecStart=/opt/pos-mahenz/BE/pos_api`.
   - FE di-build di server dari `FE/`, hasil `dist` disalin ke `/var/www/pos-web`.
   - Backup database di `BE/backups/`; skrip maintenance diambil dari `FE/maintenance/`;
     template nginx di `FE/nginx.conf`.
   - File yang tidak ikut git (mis. `.env` produksi, log, backup) tidak pindah sendiri.
4. **Git**: pindahkan dengan `git mv` dalam satu commit khusus agar riwayat file tetap terlacak.
   Sebaiknya tidak ada branch lain terbuka saat itu.
5. **Husky** (hook git di `FE/.husky`): perlu dicek masih berfungsi setelah pindah.

**[USULAN]** Cara cutover VPS: sesuai kebiasaan Anda (rename folder lama, clone ulang). Clone baru
sudah berstruktur baru, file non-git dipindahkan manual, unit systemd diubah ke `TOKO/BE`, folder
lama dibiarkan sebagai jalur rollback. Pakai maintenance mode yang sudah ada selama prosesnya.

---

## 6. Login Terpusat (SSO)

### 6.1 Pembagian tanggung jawab
- **Portal**: siapa user itu (username, nama, status aktif, password) dan boleh masuk aplikasi mana.
- **Tiap aplikasi**: role, izin menu, dan PIN tetap diatur di aplikasinya sendiri.

### 6.2 Alur login

```
1. User login di Portal (username + password).
2. Portal menampilkan aplikasi yang boleh diakses user (halaman pilih aplikasi).
3. User klik kartu aplikasi (mis. POS Toko).
4. Portal membuat kode acak sekali pakai (berlaku sekitar 60 detik), terikat ke
   user dan app_id tujuan, lalu mengarahkan browser ke: /toko/sso?code=XXXX
5. FE Toko mengirim kode itu ke BE Toko.
6. BE Toko bertanya ke BE Portal (server ke server, tidak lewat browser):
     POST /portal-api/sso/exchange
     Header: X-App-Key: <kunci rahasia milik Toko>
     Body  : { "app_id": 1, "code": "XXXX" }
7. BE Portal menjawab:
     { "valid": true, "app_id": 1,
       "user": { "portal_user_id": 12, "username": "budi",
                 "full_name": "Budi Santoso", "is_active": true } }
8. BE Toko mencari user lokal yang tertaut ke portal_user_id tersebut (membuat profil
   lokal bila belum ada), membuat sesi seperti login biasa, lalu mengirim token Toko ke browser.
```

Hal penting:
- Yang lewat URL hanya **kode pendek sekali pakai**; tanpa password, nama, atau role.
- Setelah langkah 8, request berikutnya memakai **token milik aplikasi sendiri**, seperti sekarang.
  Portal tidak dipanggil di tiap request, jadi aplikasi tetap jalan walau Portal sedang mati.
- Role **tidak** dikirim Portal; role diambil dari data lokal aplikasi.

### 6.3 Pengecekan `app_id` (berlapis)
1. Portal saat membuat kode: kode terikat ke `app_id` tujuan, dan Portal memastikan user punya
   akses ke aplikasi itu.
2. Saat penukaran kode: BE aplikasi mengirim `APP_ID` (dari `.env`-nya) dan kunci rahasia. Portal
   mencocokkan keduanya dengan kode tadi.
3. Jawaban Portal memuat `app_id`; BE aplikasi membandingkannya dengan `APP_ID` di `.env`. Tidak
   sama = ditolak.

Catatan:
- ID seperti 1 atau 2 bukan rahasia, jadi **kunci rahasia (`X-App-Key`) tetap wajib**.
- `app_id` yang dipercaya adalah yang datang dari jawaban Portal, bukan yang tertulis di URL.
- Di token Toko sekarang sudah ada klaim `apps`, tetapi isinya info perangkat (`DeviceInfo`),
  bukan ID aplikasi. Agar tidak rancu, konsep baru ini disebut `app_id`.

### 6.4 Tabel Portal (gambaran awal, belum final)
| Tabel | Isi |
|-------|-----|
| `users` | Akun Portal (username, nama, hash password Argon2id, aktif/tidak, role `admin`/`user`) |
| `sessions` | Sesi login Portal |
| `apps` | Daftar aplikasi: id, kode, nama, path (mis. `/toko`), ikon, aktif/tidak, kunci (disimpan sebagai hash) |
| `user_apps` | User mana boleh masuk aplikasi mana |
| `sso_codes` | Kode SSO (disimpan sebagai hash), terikat user dan app, masa berlaku, sudah dipakai atau belum |

### 6.5 Perubahan di tiap aplikasi (Toko dan Coffee)
- Kolom baru `portal_user_id` (nullable) di tabel `users` sebagai tautan ke akun Portal.
- Konfigurasi `.env`: `APP_ID`, alamat Portal, dan kunci rahasia aplikasi.
- Endpoint SSO baru dan halaman FE `/sso`. **Aditif**: login lama tetap berfungsi.
- Role baru **"menunggu akses"** (tanpa izin menu apa pun) dibuat lewat migrasi; dipakai sebagai role
  default saat profil lokal dibuat otomatis pada SSO pertama.
- **[USULAN]** Panggilan server ke server memakai alamat lokal (localhost) karena satu VPS.
- **[USULAN]** `SecretKey` JWT dan kunci aplikasi diberikan lewat environment variable, tidak di-commit.

### 6.6 Akses user per aplikasi
- Admin Portal mengatur, per user, aplikasi mana yang boleh diakses (centang aplikasi).
- Halaman pilih aplikasi hanya menampilkan aplikasi yang dicentang.
- SSO juga menolak di sisi server bila akses dicabut (menyembunyikan kartu saja tidak cukup).

### 6.7 Fitur Portal versi awal
Login, halaman pilih aplikasi (selalu tampil), kelola user dan akses aplikasi (admin), kelola aplikasi
(UI, admin), ganti password, logout. Role Portal: `admin` dan `user`.

### 6.8 Sesi dan logout
- Satu user satu sesi aktif (sudah berlaku di Toko: `UNIQUE KEY unique_user (user_id)` di tabel `sessions`).
- Portal memakai aturan yang sama dengan Toko: token akses 8 jam; sesi 7 hari dan diperpanjang
  tiap token diperbarui (terverifikasi dari konfigurasi dan kode login/refresh Toko).
- Logout dari aplikasi: akhiri sesi aplikasi itu lalu kembali ke Portal.
- Logout dari Portal: keluar dari Portal saja. Sesi aplikasi lain tidak otomatis berakhir
  (di perangkat bersama, orang berikutnya bisa membuka `/toko` langsung dan masih login; sama dengan
  perilaku sekarang). Logout serentak = tahap lanjutan.
- Sesi aplikasi habis atau tautan dibuka langsung: bila login langsung off, diarahkan ke Portal
  dengan membawa halaman tujuan (`next`), lalu kembali ke halaman itu setelah login. **[USULAN]**
  `next` hanya boleh berupa path internal (cegah pengalihan ke situs lain).
- **[USULAN]** User dinonaktifkan di Portal: Portal memanggil aplikasi untuk menghapus sesi user itu
  (tanpa ini, sesi yang sudah berjalan bisa hidup sampai 7 hari). Tahap lanjutan.

### 6.9 Aturan password Portal ([TERBUKA], sedang dibahas)

Fakta Toko saat ini (dari kode): username 3-50 karakter alfanumerik; password minimal 6 karakter, tanpa aturan
lain; password hanya diganti oleh pengelola user lewat `PUT` per user id (tidak ada ganti password sendiri dengan
verifikasi password lama, tidak ada lupa password). **Koreksi:** Toko ternyata PUNYA pembatas login
(`middleware/login_rate_limit_middleware.go`: 5 kali gagal per kombinasi IP+username dalam 15 menit, lalu
ditahan); catatan awal yang menyebut tidak ada pembatas keliru.

Keputusan:
1. **[SEPAKAT]** Password minimal 8 karakter, tanpa syarat huruf besar/simbol. Password lama Toko (minimal 6)
   tetap berlaku sampai diganti.
2. **[SEPAKAT]** User dapat mengganti password sendiri di Portal dengan memasukkan password lama.
3. **[SEPAKAT]** Admin Portal dapat menetapkan password sementara; user diminta menggantinya saat login berikutnya.
4. **[SEPAKAT]** Lupa password: lewat admin (tanpa email/SMS).
5. **[SEPAKAT, diputuskan ulang]** Portal **memakai pembatas login seperti Toko** (5 kali gagal per kombinasi
   IP+username dalam 15 menit, lalu ditahan). Keputusan semula "tanpa batas" dibatalkan karena dibuat atas info
   keliru (Toko ternyata punya pembatas). Keputusan semula (dibatalkan): tidak ada batas percobaan login salah. Pesan error tetap generik
   ("Username atau password salah"). Risiko yang diterima: Portal terbuka dari internet sehingga password bisa
   ditebak berulang kali; pertimbangkan kembali bila ada indikasi penyalahgunaan (alternatif ringan: pembatasan
   di nginx atau jeda antar percobaan).
6. **[SEPAKAT]** Ganti password di Portal tidak mengubah password Toko selama login langsung masih aktif.

Catatan: poin 1-4 dan 6 dianggap disetujui karena user hanya menolak poin 5; user dapat mengoreksi.

### 6.10 Penyimpanan kunci aplikasi (`X-App-Key`)

- **[SEPAKAT]** Kunci aplikasi **disimpan di database aplikasi masing-masing, bukan di `.env`**, lewat **tabel
  baru khusus kunci** (nama usulan: `portal_credentials`: `app_id`, `app_key`, `updated_at`).
- Portal tetap hanya menyimpan **hash** kunci di tabel `apps`; aplikasi menyimpan kunci aslinya karena harus
  mengirimkannya ke Portal.
- Konsekuensi yang dicatat AI: kunci asli tersimpan di database aplikasi, sehingga ikut ke file backup SQL
  (`backups/*.sql`); backup perlu diperlakukan rahasia. Koneksi database dan `SecretKey` JWT tetap di config/`.env`
  (tidak bisa pindah ke DB).
- **[SEPAKAT]** Kunci dimasukkan ke tabel lewat **SQL manual** (tanpa layar pengaturan di aplikasi).
- **[SEPAKAT]** `APP_ID` **ikut pindah ke tabel yang sama** (satu baris berisi `app_id` dan `app_key`), tidak lagi
  di `.env`. Ini memperbarui keputusan 9 (pencocokan `app_id` jawaban Portal kini dengan nilai di tabel, bukan `.env`).
  Catatan AI: pengaman "salah tempel kunci ke aplikasi lain" tetap ada karena `app_id` dan `app_key` ditempel
  bersama dalam satu perintah SQL.
- Catatan implementasi: migrasi tiap aplikasi membuat tabelnya kosong; BE menolak login SSO dengan pesan jelas bila
  barisnya belum diisi. Panduan SQL pengisian masuk ke SOP deploy.

---

## 7. Migrasi User (Toko ke Portal)

Disalin dan ditautkan, **bukan dipindah**: tabel `users` Toko dirujuk sesi (FK) dan transaksi/shift/laci kas
menyimpan `user_id`, sehingga barisnya tidak boleh hilang.

- Skrip sekali jalan, bisa diulang aman, ada mode dry-run, didahului backup database.
- Yang disalin: username, nama, status aktif, **hash Argon2id yang sama** (password tidak berubah).
- Dibuat tautan `portal_user_id` di `users` Toko, plus akses aplikasi Toko di Portal.
- Setelah itu identitas (nama, aktif/nonaktif, password SSO) dikelola di Portal; Toko menyimpan profil
  lokal, role, dan PIN.
- Selama login langsung masih on, password Toko dan password Portal berjalan terpisah (tidak
  disinkronkan). Setelah login langsung off, hanya password Portal yang berlaku.
- Staf baru dibuat di Portal; profil lokal di aplikasi dibuat otomatis saat SSO pertama dengan role
  "menunggu akses", lalu admin aplikasi mengatur role sebenarnya.
- Coffee tidak punya akun lama, jadi tidak ada migrasi; profil dibuat otomatis saat SSO pertama.

---

## 8. Login Langsung (Toggle On/Off)

- Berlaku di Toko dan Coffee. On = form login aplikasi masih tampil. Off = form tidak tampil.
- Hanya role **admin** yang boleh mengubah toggle. **[SEPAKAT]**
- **[USULAN]** Disimpan di pengaturan aplikasi (tabel `settings`), jadi bisa diubah dari menu tanpa deploy.
- **[USULAN]** Ditegakkan di BE: saat off, endpoint login langsung menolak; FE menyembunyikan form
  lewat endpoint publik kecil yang memberi status.
- **[USULAN]** Pengaman lockout: bila Portal mati saat toggle off, tidak ada yang bisa masuk. Disediakan
  flag darurat di config server (mengalahkan setelan database) yang dinyalakan lewat SSH.

---

## 9. Penyimpanan Browser (Key)

Karena semua aplikasi satu origin (hanya beda path), localStorage dipakai bersama. Aplikasi Coffee adalah
salinan Toko, sehingga tanpa penamaan berbeda key-nya bentrok (termasuk keranjang kasir).

**Format [SEPAKAT]:** `<aplikasi>.<nama>`, termasuk Toko.

Key Toko yang ada sekarang dan nama barunya:

| Key sekarang (Toko) | Isi | Key baru |
|---------------------|-----|----------|
| `auth-session` | Token login | `toko.auth-session` |
| `cashier-draft` | Keranjang kasir yang belum selesai | `toko.cashier-draft` |
| `backdate-cashier-draft` | Keranjang kasir backdate | `toko.backdate-cashier-draft` |
| `pos.native_bt_printer_address` | Pengingat printer Bluetooth APK | `toko.native_bt_printer_address` |

Aplikasi baru: `coffee.auth-session`, `coffee.*`, `portal.auth-session`, `portal.*`. Aturan: semua key baru
wajib diawali nama aplikasinya.

**Konsekuensi perubahan key Toko (perlu diperhatikan):** bila key hanya diganti nama, kasir yang sedang login
akan ter-logout satu kali, keranjang yang belum selesai hilang, dan printer Bluetooth tablet harus dipilih ulang.

**[USULAN]** Migrasi key sekali jalan di FE Toko: saat aplikasi dimuat, bila key lama ada dan key baru belum,
salin nilainya ke key baru lalu hapus key lama. Harus dijalankan **sebelum** store (zustand persist)
di-hydrate. Dengan begitu tidak ada yang ter-logout atau kehilangan keranjang.

Catatan keamanan: karena satu origin, XSS di satu aplikasi secara teori bisa membaca token aplikasi lain.
Untuk aplikasi internal risikonya kecil, tetapi diterima sebagai trade-off pemilihan path.

---

## 10. URL dan Nginx

**[SEPAKAT]** Berbasis path, Portal di root domain `mahenz.duckdns.org`. Production "dirapikan di nginx dan
menyesuaikan URL baru".

Kondisi nginx sekarang: satu server block `mahenz.duckdns.org`; FE dari `/var/www/pos-web/dist`; `/api`
diteruskan ke 127.0.0.1:8080 (atau 8081 untuk versi lama, akan dihapus); ada maintenance mode dan bypass versi.

**[USULAN]** Pembagian path:

| Path | Tujuan |
|------|--------|
| `/` | FE Portal |
| `/portal-api/` | BE Portal |
| `/toko/` | FE Toko (Vite `base: '/toko/'`, router `basename`) |
| `/toko/api/` | BE Toko |
| `/coffee/` | FE Coffee |
| `/coffee/api/` | BE Coffee |
| `/api` (lama) | Tetap diarahkan ke BE Toko selama masa peralihan agar FE ter-cache / APK lama tidak putus |

Keuntungan satu origin: tanpa CORS, satu sertifikat HTTPS, dan APK tetap satu origin saat berpindah
Portal ke aplikasi sehingga jembatan Bluetooth seharusnya tetap aktif (perlu diuji).

### 10.1 Detail nginx (hasil diskusi)

- **[SEPAKAT]** Acuan konfigurasi nginx adalah `nginx/pos-web.conf` di repo. Catatan AI: file itu hanya
  `listen 80`; blok HTTPS (certbot) mungkin ada di server dan belum terlihat di repo. `nginx/pos-web-old.conf` dan
  `FE/nginx.conf` adalah salinan/variasi dengan mekanisme sama. Perlu dicek di VPS sebelum cutover, supaya
  konfigurasi baru dibuat dari yang benar-benar berjalan. Ketiga file perlu disatukan menjadi satu sumber.
- **[SEPAKAT]** Tiga FE satu domain, tiap FE di-build dengan base path sendiri dan ditaruh di folder terpisah
  (mis. `/var/www/pos-web/portal`, `/toko`, `/coffee`).
- **[SEPAKAT]** Maintenance mode: untuk saat ini satu saklar untuk **semua** aplikasi, tetapi dirancang juga bisa
  **per aplikasi** (skrip di `FE/maintenance/` menerima nama aplikasi; tanpa argumen = semua).
- **[SEPAKAT]** Cutover VPS: rename folder lama, lalu clone ulang (kebiasaan yang sudah ada).
- **[SEPAKAT]** Versi lama/baru (Mode B: cookie `ver_bypass`, port 8081, folder `dist_old`, dokumen
  `docs/DEPLOYMENT_REDEPLOY_FULL.md`) **tidak dipakai lagi**; semuanya satu versi baru. Konsekuensi saat
  implementasi: hapus `map $cookie_ver_bypass`, parameter `version` di `/maintenance-bypass`, dan `dist_old` dari
  nginx; hapus/tandai usang bagian Mode B di docs; hapus service `pos-backend-old` di VPS bila ada. Port 8081
  dipakai ulang oleh Toko, jadi service lama di 8081 wajib sudah mati sebelum cutover.
- **[SEPAKAT]** Port dan nama service (3 service, satu per BE):

  | Aplikasi | Service systemd | Port |
  |----------|-----------------|------|
  | Portal | `pos-portal-backend` | 8080 |
  | Toko | `pos-toko-backend` | 8081 |
  | Coffee | `pos-coffee-backend` | 8082 |

  Catatan AI: Toko berubah dari `pos-backend`:8080 menjadi `pos-toko-backend`:8081, jadi Fase 1 (pindah folder)
  bukan lagi tanpa perubahan production: service diganti nama, port berubah, dan `/api` di nginx diarahkan ke
  8081. Port 8080 baru dipakai Portal di Fase 2; sebelum itu 8080 kosong. Perlu juga diubah: `APP_PORT` di `.env`
  Toko, `ExecStart`/`WorkingDirectory` service menjadi `/opt/pos-mahenz/TOKO/BE`, dan nginx. Port dev di laptop
  mengikuti BE yang sama.

- **[SEPAKAT]** Konfigurasi nginx di VPS dikonfirmasi user **sama persis** dengan `nginx/pos-web.conf` (jadi acuan
  sudah benar). Catatan AI: file itu hanya `listen 80`, sedangkan APK memakai `https://`; kemungkinan TLS diurus di
  luar nginx ini. Cara TLS diurus perlu diketahui sebelum cutover, karena mempengaruhi `X-Forwarded-Proto`.
- **Alasan 3 service:** tiap aplikasi punya binary, `.env`, database, dan port sendiri; satu aplikasi gagal atau
  di-restart tidak mematikan yang lain.
- **[SEPAKAT]** HTTPS (hasil cek VPS lewat screenshot user): nginx mendengarkan port **80 dan 443**, sertifikat
  Let's Encrypt `mahenz.duckdns.org` ada di `/etc/letsencrypt/live/`. Jadi TLS diurus **nginx sendiri lewat
  certbot**. Port 8080 dipakai `pos_api` (Toko lama); **tidak ada** yang mendengarkan di 8081 (tidak ada
  `pos-backend-old`), jadi 8081 bebas untuk `pos-toko-backend`.
- **[SEPAKAT]** User mengonfirmasi SSL memang memakai certbot.
- **[TERBUKA]** `nginx/pos-web.conf` di repo hanya `listen 80`, padahal server memakai 443. Certbot biasanya
  menambahkan blok `listen 443 ssl` dan redirect ke file di server, sehingga file di server kemungkinan berbeda
  dari repo walau user menilai sama. Saat cutover, nginx baru wajib mempertahankan blok SSL (atau dijalankan
  ulang `certbot --nginx`), jika tidak HTTPS dan APK putus. Perlu dicek: `sudo grep -n ssl /etc/nginx/sites-enabled/*`.

Detail lain (isi nginx final) disusun setelah pengecekan VPS.

---

## 11. Alur Perpindahan Pengguna (Transisi Portal ke Aplikasi)

- Klik kartu aplikasi: pindah di **tab yang sama** (bukan tab baru, bukan iframe). Alasan: tablet/APK hanya
  punya satu WebView; iframe bermasalah untuk cetak, Bluetooth, popup, dan izin perangkat.
- Layar antara singkat "Memasuki POS Toko..." saat kode ditukar. Bila gagal: pesan jelas, tombol
  "Kembali ke Portal", dan "Login langsung" bila toggle on.
- Ganti aplikasi: tombol "Portal" di header tiap aplikasi; sesi Portal masih hidup sehingga pindah
  aplikasi cukup satu klik tanpa login ulang.
- Aturan logout dan sesi habis: lihat bagian 6.8.

---

## 12. APK Tablet

- **[SEPAKAT]** Satu APK untuk semua aplikasi (APK hanya jembatan cetak Bluetooth; aplikasinya tetap web).
- APK **tidak perlu diubah dan tidak perlu dibuat ulang**: alamat awalnya (`https://mahenz.duckdns.org`)
  otomatis menjadi Portal.
- Konsekuensi: kasir menekan satu kartu "POS Toko" tiap membuka aplikasi. Bila terasa mengganggu,
  nanti Portal bisa diberi opsi "buka aplikasi terakhir otomatis" per perangkat (tanpa menyentuh APK).
- Perlu diuji: berpindah path di domain yang sama harus tetap berjalan di dalam WebView APK (tidak terlempar
  ke Chrome) dan jembatan Bluetooth tetap aktif di Toko maupun Coffee.
- Rancangan APK lengkap: `docs/RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`; panduan build:
  `docs/PANDUAN_BUILD_APK_TABLET.md`.

---

## 13. Fakta Kode Saat Ini (dasar rancangan)

Terverifikasi dari kode/konfigurasi repo:
- Password memakai **Argon2id** (package bernama `pkg/bcrypt`, tetapi isinya Argon2id format PHC; koreksi dari
  catatan awal yang menulis bcrypt). Karena Portal menyalin package yang sama, hash bisa disalin apa adanya.
- Token JWT berisi klaim `user_id`, `username`, role, dan `apps` (isinya `DeviceInfo`). Masa berlaku token
  dari `TokenExpire` (28800 detik = 8 jam); sesi di tabel `sessions` berakhir sesuai `RefreshTokenExpire`
  (604800 detik = 7 hari), diperbarui tiap refresh.
- Tabel `sessions`: `UNIQUE KEY unique_user (user_id)`, FK ke `users(id)`.
- Role dan izin menu ada di database Toko; ada fitur PIN per user.
- FE menyimpan sesi lewat zustand persist (key `auth-session`), plus `cashier-draft`,
  `backdate-cashier-draft`, dan `pos.native_bt_printer_address`.
- Nginx dan alur deploy VPS seperti dijelaskan di bagian 5.1 dan 10.
- Konfigurasi APK (`mobile/capacitor.config.ts`) menunjuk ke `https://mahenz.duckdns.org`.

---

## 14. Fitur POS Coffee (sedang dibahas)

> Semua butir di bagian ini masih **[USULAN]** atau **[TERBUKA]** sampai user memutuskan.

### 14.1 Pemetaan modul Toko ke Coffee

Berdasarkan daftar domain BE Toko (`BE/domain/`):

| Kelompok | Modul Toko | Usulan untuk Coffee |
|----------|------------|---------------------|
| Infrastruktur umum | `auth`, `user`, `role`, `access`, `menu`, `route_registry`, `pin`, `setting`, `print`, `backup` | Disalin apa adanya |
| Operasional kasir | `cash_drawer`, `shift`, `expense`, `payment_method`, `payment_status` | Disalin, disesuaikan seperlunya |
| Pelaporan | `dashboard`, `report`, `business_summary`, `finance` | Kerangka disalin, isi laporan diganti ke konteks kopi |
| Penjualan | `transaction` | Diubah besar: pesanan dengan varian/tambahan, tipe pesanan, status |
| Katalog | `product`, `product_category`, `product_unit` | Diubah menjadi menu kopi (varian dan tambahan); satuan bertingkat tidak dibawa |
| Stok dan pengadaan | `stock_mutation`, `stock_reconciliation`, `expiry_batch`, `supplier`, `supplier_purchase`, `supplier_return` | `supplier`, `supplier_purchase`, `stock_mutation`, `stock_reconciliation` disalin lalu diubah untuk bahan baku (P3); `expiry_batch` dan `supplier_return` dipangkas |
| Pelanggan | `customer`, `customer_balance`, `receivable` | `customer` dan `receivable` dibawa (sistem utang seperti Toko); `customer_balance` dipangkas; member/loyalty dibahas terpisah |
| Lain | `backdate`, `capital_reconciliation`, `version` | Tidak dibawa kecuali diminta |

Modul cetak struk (9 metode) ikut tersalin dan bisa dipakai juga untuk tiket dapur.

**[USULAN]** Penamaan: modul yang khas kopi memakai istilah kopi (mis. `menu_items`, `menu_variants`,
`modifiers`, `orders`), sedangkan modul infrastruktur tetap bernama sama dengan Toko agar mudah
dibandingkan dan perbaikan bug mudah diterapkan di kedua sisi.

### 14.2 Kandidat fitur khas kedai kopi

| Fitur | Keterangan |
|-------|------------|
| Menu dengan varian | Ukuran (S/M/L), panas/dingin; harga bisa berbeda per varian |
| Tambahan (modifier) | Extra shot, topping, level gula, level es; ada yang berbayar, ada yang gratis |
| Catatan per item | Mis. "kurang manis" |
| Tipe pesanan | Dine-in, take away; nomor meja atau nomor antrean |
| Status pesanan dan antrean | Baru, dibuat, siap, selesai; layar antrean untuk bar/dapur |
| Tiket dapur/bar | Cetak tiket pesanan ke printer bar (terpisah dari struk pelanggan) |
| Open bill | Pesan dulu, bayar belakangan, per meja |
| Split / gabung bill, pindah meja | Lanjutan open bill |
| Resep dan stok bahan baku | Tiap menu punya resep; stok bahan (biji kopi, susu, sirup, gelas) berkurang otomatis saat terjual |
| Service charge | Selain pajak yang sudah ada |
| Promo / voucher / member | Belum jelas kebutuhannya |
| Laporan khas kopi | Penjualan per menu, kategori, dan jam; bahan terpakai; laba per menu (butuh resep dan harga bahan) |
| Tampilan sentuh | Grid menu besar, bergambar, nyaman di tablet |

### 14.3 Usulan tahapan

| Tahap | Isi |
|-------|-----|
| Versi 1 (inti) | Menu + kategori + varian + tambahan, pesanan dine-in (nomor meja) dan take away (nomor antrean), pembayaran, struk, kas dan shift, pengeluaran, user/role/izin, cetak struk, **katalog bahan baku (satu satuan stok per bahan), pembelian bahan dan barang jadi dari supplier, mode stok per menu (potong langsung untuk barang jadi, manual untuk racikan), pencatatan pemakaian bahan manual, opname bahan, resep (catatan takaran standar)**, open bill (bayar sekarang atau bayar nanti), gabung bill, pindah meja, sistem utang untuk pesanan yang tidak dibayar di hari yang sama, laporan penjualan dasar dan pemakaian bahan |
| Versi berikutnya (tergantung jawaban 14.7 dan 14.9) | Layar antrean, tiket bar, promo/member, laba per menu (HPP), perkiraan pemakaian dari resep |

### 14.4 Keputusan yang memengaruhi struktur database sejak awal ([TERBUKA])

Pertanyaan ini perlu dijawab sebelum desain database Coffee, karena mengubah tabel inti:
1. Dine-in dengan nomor meja, atau lebih ke order di kasir dengan nomor antrean?
2. Bar/dapur terpisah dari kasir (perlu layar atau printer sendiri)?
3. Bayar di depan, atau perlu bayar belakangan (open bill)?
4. Stok bahan lewat resep dari awal, atau stok per menu dulu?
5. Pembelian bahan baku dicatat di aplikasi atau tidak?
6. Pelanggan/member dibutuhkan?
7. Pajak saja, atau juga service charge?
8. Foto menu diunggah dan ditampilkan?
9. Satu outlet atau banyak cabang?
10. Perangkat: tablet kasir saja, printer yang sama (EP58M), printer bar tambahan?

### 14.5 Pendekatan: salin dulu dari Toko

**[USULAN user]** Coffee dimulai dengan menyalin BE dan FE Toko, lalu yang diubah lebih dulu hanya
menu **Kasir** dan **Pembelian**.

Catatan AI (menunggu konfirmasi user):
1. **Waktu penyalinan:** di Fase 4, setelah Toko terhubung ke Portal. Dengan begitu Coffee otomatis
   mewarisi SSO, toggle login langsung, role "menunggu akses", kolom `portal_user_id`, dan penamaan key.
2. **Pangkas langsung setelah salin:** modul yang jelas tidak dipakai dihapus (bukan dibiarkan).
   Kandidat: `supplier_return`, `expiry_batch`, `customer_balance`, `capital_reconciliation`,
   `backdate`, `version`. (`receivable` semula ada di daftar ini, kemudian dikeluarkan karena sistem utang
   dibawa ke Coffee; lihat 14.10.)
3. **Database Coffee memakai migrasi baseline baru** yang hanya berisi tabel yang dibutuhkan, bukan
   menyalin 10 migrasi retail Toko (agar tidak ada tabel dan data seed retail yang tidak terpakai).
4. **Daftar yang wajib diganti setelah salin:** module Go, nama package FE, `APP_ID`, port, nama database
   dan user MySQL, `SecretKey` JWT (harus berbeda per aplikasi), prefix key browser `toko.` menjadi
   `coffee.`, Vite `base` dan router `basename`, path API, judul/branding, file config.
5. **Menu Pembelian bergantung pada model stok bahan baku** (di Toko, satu produk dibeli sekaligus dijual;
   di Coffee bahan yang dibeli berbeda dari menu yang dijual). Pilihan:
   - **P1** pembelian bahan dicatat sebagai pengeluaran saja (modul `expense` sudah ada), tanpa stok;
   - **P2** katalog bahan baku + pembelian menambah stok, stok dikoreksi lewat opname (tanpa resep);
   - **P3** seperti P2 ditambah resep, sehingga stok bahan berkurang otomatis saat menu terjual.
   Bila pembelian menambah stok tetapi tidak ada resep atau opname, stok bahan tidak pernah berkurang.
6. **Daftar modul infrastruktur yang tersalin dicatat** agar perbaikan bug di Toko dicek juga di Coffee, dan
   sebaliknya.

### 14.6 Jawaban user atas 14.5

- **[SEPAKAT]** Coffee dimulai dengan menyalin BE dan FE Toko; yang diubah lebih dulu adalah menu Kasir dan Pembelian.
- **[SEPAKAT]** Menu "Pembelian" (grup Pengadaan) di Coffee berarti pembelian **bahan baku** dari supplier.
  Prinsipnya berbeda dari Toko, di mana barang yang dibeli dari supplier adalah barang yang dijual.
- **[SEPAKAT]** Model stok bahan: **P3** (katalog bahan baku, pembelian bahan menambah stok, resep per menu, stok
  berkurang otomatis saat menu terjual). Masuk versi pertama karena daftar menu sudah siap.
- **[SEPAKAT]** Menu Kasir versi pertama: varian, tambahan, nomor meja, dan nomor antrean.
- **[SEPAKAT]** Daftar modul yang dipangkas (14.5 poin 2) sudah sesuai untuk saat ini.
- Belum dikonfirmasi user: 14.5 poin 1 (waktu salin di Fase 4), poin 3 (baseline migrasi baru), poin 4
  (daftar yang diganti setelah salin), poin 6 (daftar modul yang dicatat sebagai salinan).

### 14.7 Dampak P3 dan pertanyaan lanjutan ([TERBUKA])

Catatan desain ([USULAN]):
- Bahan baku dan menu adalah dua katalog terpisah: bahan dibeli, menu dijual.
- Satuan beli dan satuan pakai berbeda (beli per kg atau dus, pakai gram atau ml). Konsep satuan bertingkat
  (`product_packages`) di Toko bisa menjadi dasar untuk konversi ini.
- Resep terikat ke menu dan varian (takaran S/M/L bisa berbeda); tambahan (modifier) bisa menambah atau
  mengganti bahan.
- Pembelian bahan memakai modul `supplier` dan `supplier_purchase` Toko yang diubah target barangnya dari
  produk menjadi bahan baku; opname memakai `stock_reconciliation`; pengeluaran lain memakai `expense`.
- Pengurangan stok bahan perlu aturan jelas soal waktu dan pembatalan (lihat pertanyaan 2).

Pertanyaan lanjutan:
1. Apakah resep berbeda per ukuran? Apakah tambahan menambah bahan (extra shot) atau mengganti bahan (susu oat)?
2. Stok bahan dikurangi kapan: saat pesanan dibuat, saat dibayar, atau saat siap? Bagaimana bila pesanan dibatalkan?
3. Bila stok bahan tidak cukup: menu otomatis ditandai habis, atau tetap boleh dijual dengan peringatan?
4. Bahan terbuang, tumpah, atau kedaluwarsa perlu dicatat?
5. HPP dan laba per menu masuk versi pertama? Bila ya, biaya bahan memakai harga rata-rata atau harga beli terakhir?
6. Daftar menu "sudah siap" dalam bentuk apa (mis. Excel)? Perlu fitur impor untuk menu, bahan, dan resep?
7. Dine-in: bayar di kasir dulu, atau pesan dulu lalu bayar belakangan (open bill)?
8. Nomor antrean: formatnya bagaimana dan di-reset kapan (harian atau per shift)? Daftar meja dikelola di aplikasi?
9. Ada layar atau printer bar terpisah yang menerima pesanan?

### 14.8 Jawaban user atas 14.7 (tahap 1)

Ringkasan jawaban user:
- Pertanyaan 1 sampai 3 (resep per ukuran dan tambahan, waktu pengurangan stok, stok tidak cukup): user
  **tidak ingin konversi dan pengurangan stok terlalu ketat**, karena takaran di lapangan bisa berbeda dan
  sistem yang kaku akan memunculkan kesalahan. Waktu pengurangan stok akan didiskusikan lebih lanjut.
- Pertanyaan 7 (dine-in): user ingin **dua opsi**, bayar sekarang atau pesan dulu bayar belakangan.

**Tingkat otomatisasi stok bahan ([USULAN], menunggu pilihan user):**

| Pilihan | Cara kerja | Risiko |
|---------|------------|--------|
| Ketat | Stok berkurang otomatis persis sesuai resep; satuan beli dikonversi tetap oleh sistem | Selisih menumpuk; sistem terlihat "salah" saat takaran di lapangan berbeda |
| Lentur (rekomendasi AI) | Resep berfungsi sebagai standar/perkiraan. Sistem menampilkan stok perkiraan, stok sebenarnya dikoreksi lewat opname, dan laporan selisih menunjukkan perbedaannya. Stok minus diperbolehkan (peringatan, tidak memblokir penjualan); menu ditandai habis secara manual. Pengurangan otomatis bisa dinyalakan atau dimatikan per bahan | Perlu opname berkala |
| Manual | Tidak ada pengurangan otomatis; resep hanya untuk perkiraan pemakaian dan HPP; stok berubah lewat pembelian dan opname atau pemakaian harian | Lebih banyak input manual |

**Pengaman agar tidak muncul "kesalahan" ([USULAN]):**
- Satu bahan hanya punya **satu satuan stok** (satuan pakai: gram, ml, pcs).
- Konversi saat pembelian **tidak dikunci permanen**: user mengisi jumlah dalam satuan pakai (atau jumlah
  kemasan dikali isi aktual), karena isi tiap pembelian bisa berbeda.
- Setiap penambahan atau pengurangan stok tercatat di riwayat mutasi dengan alasan dan referensi (pesanan,
  pembelian, opname), sehingga bisa dilacak dan dikoreksi.
- Pemakaian bahan per pesanan **disimpan sebagai salinan saat pesanan dibuat**; mengubah resep hanya
  memengaruhi pesanan berikutnya, tidak mengubah riwayat.
- Koreksi lewat opname tidak dianggap error; selisihnya dilaporkan.
- Bila memilih Lentur, pertanyaan 3 (stok tidak cukup) terjawab: tidak memblokir penjualan.

**Dua opsi pembayaran / open bill dasar (implikasi, [USULAN]):**
- Pesanan menjadi entitas sendiri yang dibuat lebih dulu, pembayaran menyusul. Pada "bayar sekarang",
  pesanan dan pembayaran terjadi berurutan dalam satu alur.
- Pesanan "bayar nanti" berstatus terbuka dan terikat ke nomor meja atau antrean; item bisa ditambah
  sebelum dibayar.
- Pesanan diteruskan ke bar saat dibuat, bukan saat dibayar.
- Pendapatan dan kas diakui saat pembayaran.
- Perlu aturan untuk pesanan yang masih terbuka saat shift atau hari ditutup.
- Split bill, gabung bill, dan pindah meja tetap di tahap berikutnya.

Pertanyaan lanjutan ([TERBUKA]):
1. Tingkat otomatisasi stok: Ketat, Lentur, atau Manual?
2. Saat shift ditutup dan masih ada pesanan terbuka: ditolak sampai diselesaikan, atau boleh dilanjutkan ke
   shift berikutnya?
3. Open bill dasar (tambah item, bayar di akhir) cukup untuk versi pertama, dengan split bill dan pindah meja
   menyusul?

### 14.9 Jawaban user atas 14.8

- **[SEPAKAT]** Stok bahan: **Lentur dengan pencatatan pemakaian manual**. Contoh dari user: beli 10 bungkus kopi;
  hari pertama terpakai 1 bungkus dan diinput manual "terpakai 1"; stok tersimpan menjadi sisa 9, dengan riwayat
  terpakai 1.
- **[SEPAKAT]** Pesanan terbuka saat shift ditutup **boleh dilanjutkan** ke shift berikutnya.
- **[SEPAKAT]** Split bill, gabung bill, dan pindah meja **ikut masuk versi pertama** ("sekalian").

Tafsiran AI atas contoh stok ([USULAN], menunggu konfirmasi):
- Stok sisa = pembelian dikurangi pemakaian yang diinput manual, dikoreksi lewat opname.
- Satuan stok = satuan yang dipakai staf sehari-hari (mis. bungkus), tanpa konversi ke gram.
- Karena itu **resep tidak mengurangi stok sisa**, supaya tidak terjadi pengurangan ganda. Pengurangan otomatis dari
  resep bisa disediakan per bahan dengan default mati (opsional).
- Catatan pemakaian berisi tanggal, bahan, jumlah (dalam satuan stok), catatan, dan siapa yang menginput; tersimpan
  di riwayat mutasi stok.
- Stok minus diperbolehkan dengan peringatan; menu ditandai habis secara manual.

Hal yang perlu diluruskan: bila stok dikelola dari input pemakaian manual dalam satuan seperti bungkus, resep
(takaran dalam gram per cangkir) tidak bisa dibandingkan dengan pemakaian tanpa angka konversi (isi per bungkus).
Laba per menu (HPP) sudah ditunda ke versi berikutnya, sehingga resep di versi pertama belum punya "pemakai" selain
laporan perkiraan.

Pertanyaan lanjutan ([TERBUKA]):
1. Benar stok sisa hanya berkurang lewat input pemakaian manual dan opname, bukan dari resep?
2. Peran resep di versi pertama: (a) tetap dimasukkan sebagai catatan takaran standar per menu, sementara perkiraan
   pemakaian dan HPP menyusul; atau (b) ditunda penuh ke versi berikutnya?
3. Pencatatan pemakaian dilakukan per hari, per shift, atau kapan saja? Satu bahan per input, atau banyak bahan sekaligus?
4. Pesanan terbuka boleh lintas hari, atau hanya antar shift di hari yang sama? (Di Toko ada scheduler laci kas,
   `StartCashDrawerScheduler`; perilakunya terhadap pesanan terbuka perlu dicek.)
5. Split bill yang dimaksud: pisah per item, bagi rata, atau keduanya?

### 14.10 Jawaban user atas 14.9

- **[SEPAKAT]** Stok sisa hanya berkurang lewat input pemakaian manual dan opname, bukan dari resep.
- **[SEPAKAT]** Peran resep di versi pertama: **catatan takaran standar per menu** (data resep sudah ada).
  Perkiraan pemakaian dan HPP menyusul.
- **[SEPAKAT]** Pemakaian stok dicatat **kapan saja**, **satu bahan per input**. Maknanya: setiap kali stok mulai
  dipakai (mis. satu bungkus kopi dibuka dan baru habis setelah 4 hari), user menginputnya sebagai "stok terpakai".
  Berlaku juga untuk krimer, es, dan bahan lain. Per bahan ditampilkan stok sisa dan total terpakai.
- **[SEPAKAT]** Pesanan terbuka **tidak boleh lintas hari**. Yang belum dibayar saat batas hari menjadi **utang**,
  seperti sistem utang di Toko.
- **[SEPAKAT]** **Split bill dihapus** dari rencana. Gabung bill dan pindah meja tetap di versi pertama.

Dampak ([USULAN]):
- Modul `receivable` dan `customer` tidak jadi dipangkas dari Coffee karena sistem utang dibawa. `customer_balance`
  (saldo pelanggan) tetap dipangkas kecuali diminta.
- Terverifikasi di kode Toko: `scheduler/cash_drawer_scheduler.go` berjalan tiap pukul 00:00 (zona Asia/Jakarta) dan
  menutup otomatis kas hari sebelumnya (`AutoCloseYesterday`). Jadi "hari" di Toko adalah hari kalender dengan batas
  tengah malam. Untuk kedai yang buka sampai lewat tengah malam, batas ini perlu ditinjau.
- Utang membutuhkan pelanggan, sehingga pesanan terbuka yang dijadikan utang harus dikaitkan ke pelanggan.
- Pemakaian stok memakai satuan stok bahan dan tercatat di riwayat mutasi (siapa, kapan, catatan).

Pertanyaan lanjutan ([TERBUKA]):
1. Setuju `receivable` dan `customer` dibawa ke Coffee, sedangkan `customer_balance` tetap dipangkas?
2. Aturan utang di Coffee sama seperti Toko apa adanya (batas kredit, jatuh tempo, cicilan)?
3. Definisi "hari": tengah malam seperti Toko, atau mengikuti jam tutup kedai? Jam operasional kedai?
4. Pesanan terbuka yang belum dibayar saat batas hari dan belum punya pelanggan: (A) tutup kas manual diblokir sampai
   kasir menyelesaikan (bayar, jadikan utang dengan memilih pelanggan, atau batalkan); (B) untuk penutupan otomatis
   tengah malam, otomatis dijadikan utang atas nama "Pelanggan Umum" yang dilengkapi kemudian; (C) otomatis
   dibatalkan. Saran AI: A untuk tutup manual, B untuk penutupan otomatis.
5. Jumlah pemakaian stok boleh desimal (mis. 0,5 bungkus)?
6. Data menu dan resep yang sudah siap berbentuk apa (mis. Excel), dan perlu fitur impor? (Pertanyaan 6 di 14.7
   masih terbuka.)

### 14.11 Jawaban user atas 14.10 dan kebutuhan baru: barang jadi

Jawaban:
- **[SEPAKAT]** `receivable` dan `customer` dibawa ke Coffee; `customer_balance` tetap dipangkas.
- **[SEPAKAT]** Aturan utang sama seperti Toko apa adanya.
- **[SEPAKAT]** Batas hari mengikuti Toko (tengah malam).
- **[SEPAKAT]** Opsi A: tutup kas manual diblokir sampai pesanan terbuka diselesaikan (bayar, jadikan utang dengan
  memilih pelanggan, atau batalkan).
- **[SEPAKAT]** Jumlah pemakaian stok boleh desimal.
- **[SEPAKAT]** Tidak ada fitur impor menu dan resep di versi pertama; data dicatat manual, resep diisi dari
  form menu (menu punya pilihan resep).

Celah yang masih terbuka ([TERBUKA]): opsi A hanya mengatur tutup kas **manual**. Penutupan otomatis tengah malam
(scheduler Toko) tetap membutuhkan aturan untuk pesanan yang masih terbuka.

**Kebutuhan baru: snack dan roti (barang jadi).** Kedai juga menjual snack dan roti; sebagian dibeli jadi dari
supplier sehingga tidak memakai bahan. Karena itu stok diatur per menu:

| Mode stok menu | Untuk | Perilaku stok |
|----------------|-------|---------------|
| Potong stok langsung | Barang jadi dari supplier (roti, snack tertentu) | Stok ada di menu itu sendiri; pembelian dari supplier menambah, penjualan mengurangi otomatis |
| Manual (racikan) | Dibuat sendiri dari bahan | Tidak ada stok di level menu; stok bahan dikurangi lewat pemakaian manual; resep hanya catatan takaran |

Tidak semua makanan memakai mode langsung; modenya dipilih per menu.

Implikasi ([USULAN]):
- Pembelian dari supplier di Coffee punya dua jenis tujuan: bahan baku, dan menu bermode langsung. Untuk barang jadi,
  perilakunya mirip produk di Toko (dibeli dari supplier, dijual, stok berkurang otomatis), sehingga logika stok
  produk Toko bisa dipakai ulang.
- Opname dan riwayat mutasi stok berlaku untuk bahan baku maupun barang jadi.
- Mode langsung di versi pertama hanya untuk menu tanpa varian dan tanpa tambahan (agar stok tetap sederhana).

Pertanyaan lanjutan ([TERBUKA]):
1. Barang jadi dibeli per pcs, atau per dus/bungkus lalu dijual per pcs? (Bila per dus, jumlah dan isi diisi saat
   pembelian, tidak dikunci permanen.)
2. Kapan stok barang jadi dipotong: saat item masuk ke pesanan (disarankan karena ada open bill; dikembalikan bila
   item dihapus atau pesanan dibatalkan) atau saat dibayar?
3. Stok barang jadi habis: penjualan diblokir, atau hanya peringatan?
4. Mode langsung hanya untuk menu tanpa varian dan tambahan di versi pertama, cukup?
5. Penutupan otomatis tengah malam saat masih ada pesanan terbuka: (i) kas yang masih punya pesanan terbuka tidak
   ditutup otomatis dan ditandai untuk diselesaikan, (ii) pesanan otomatis dijadikan utang atas nama "Pelanggan
   Umum" untuk dilengkapi kemudian, atau (iii) pesanan otomatis dibatalkan?

### 14.12 Jawaban user atas 14.11

- **[SEPAKAT]** Satuan barang jadi mengikuti konsep Toko: satuan bertingkat **dus dan pcs** (konversi tetap seperti di Toko).
- **[SEPAKAT]** Waktu potong stok barang jadi: mengikuti rekomendasi AI, yaitu **saat item masuk ke pesanan**; stok
  dikembalikan bila item dihapus atau pesanan dibatalkan.
- **[SEPAKAT]** Stok barang jadi habis: mengikuti rekomendasi AI, yaitu **penjualan diblokir**.
- **[SEPAKAT]** Penutupan otomatis tengah malam: **opsi i**, kas yang masih punya pesanan terbuka tidak ditutup
  otomatis dan ditandai untuk diselesaikan.
- Pertanyaan 4 di 14.11 (varian dan tambahan pada barang jadi) belum jelas bagi user; penjelasan diberikan dan
  jawabannya ditunggu (lihat pertanyaan di bawah).

Implikasi ([USULAN]):
- Dua model satuan berdampingan: **bahan baku** (satu satuan stok, tanpa konversi tetap) dan **barang jadi**
  (dus dan pcs seperti Toko). Logika satuan bertingkat (`product_packages`) Toko disalin untuk barang jadi.
- Toko memotong stok saat checkout (pembayaran). Di Coffee pemotongan barang jadi dipindah ke saat item masuk
  pesanan, karena pesanan bisa terbuka berjam-jam; ini mencegah dua pelanggan memesan stok terakhir yang sama. Alur
  transaksi Toko yang disalin perlu diubah pada bagian ini.
- Karena penjualan diblokir saat stok habis, stok barang jadi harus sudah diinput (pembelian atau opname) sebelum
  dijual; ketidaksesuaian dikoreksi lewat opname.
- Opsi i berarti kas bisa tetap terbuka melewati tengah malam. Aturan pendamping yang diusulkan: pesanan yang
  melewati batas hari ditandai "melewati batas hari" dan hanya bisa diselesaikan (bayar, jadikan utang, atau
  batalkan), tidak bisa ditambah item; pembukaan kas baru diblokir sampai kas yang tertunda ditutup; admin dapat
  menyelesaikan atas nama kasir; ada daftar pesanan tertunda yang terlihat admin. Perilaku Toko saat kas kemarin
  belum tertutup perlu dicek saat desain.

Jawaban atas pertanyaan lanjutan:
1. **[SEPAKAT]** Barang jadi: satu menu = satu stok, tanpa varian; rasa atau jenis berbeda dibuat menu terpisah
   (mis. "Roti Cokelat" dan "Roti Keju"). Tambahan (mis. "dihangatkan") boleh ada tetapi tidak memengaruhi stok.
2. **[SEPAKAT]** Pesanan berisi barang jadi yang dibatalkan setelah barang diserahkan: kasir memilih saat
   membatalkan, **"kembali ke stok"** atau **"dibuang"**. Bila dibuang, stok tidak kembali dan tercatat sebagai
   barang terbuang (terkait pertanyaan 14.7 nomor 4, bahan terbuang).
3. **[SEPAKAT]** Aturan pendamping opsi i disetujui: pesanan "melewati batas hari" hanya bisa diselesaikan (bayar,
   utang, batal) dan tidak bisa ditambah item; buka kas baru diblokir sampai kas tertunda ditutup; admin dapat
   menyelesaikan atas nama kasir; ada daftar pesanan tertunda. Tersisa saat desain: cek perilaku Toko saat kas
   kemarin belum tertutup.

### 14.13 Jawaban user atas pertanyaan terbuka sebelumnya

- **[SEPAKAT]** Bahan terbuang (14.7 nomor 4): dicatat, mengikuti usulan AI, yaitu satu mekanisme "barang terbuang"
  yang sama untuk bahan baku dan barang jadi (barang jadi lewat pilihan "dibuang" saat pembatalan, bahan lewat
  input manual).
- **[SEPAKAT]** Nomor antrean (14.7 nomor 8): **direset tiap hari** (00:00 Asia/Jakarta, sama dengan batas hari).
  Format nomor dan master meja belum dijawab.
- **[SEPAKAT]** Layar atau printer bar/dapur (14.7 nomor 9): **tanpa tiket dan tanpa layar bar** di versi pertama.
  Kasir menyebutkan pesanan langsung; nomor antrean di struk dipakai untuk memanggil pelanggan. Dapat ditambah
  nanti tanpa mengubah rancangan.
- **[SEPAKAT]** Foto menu (14.4 nomor 8): **tanpa foto dan tanpa fitur unggah**; menu berupa tombol teks.
- **[SEPAKAT]** Banyak cabang (14.4 nomor 9): **satu kedai saja** untuk saat ini (tanpa pilih cabang dan laporan
  per cabang).
- **[SEPAKAT]** Perangkat Coffee (14.4 nomor 10): kasir dapat membuka lewat **browser (laptop) maupun APK tablet**,
  karena diakses sebagai website. Tidak ada pekerjaan tambahan untuk APK (satu APK untuk semua aplikasi). Tombol
  menu dibuat cukup besar agar nyaman disentuh di tablet.

Usulan ([USULAN]) untuk bagian yang belum dijawab:
- Format antrean: nomor urut sederhana per hari (001, 002, ...), tanpa awalan.
- Master meja: daftar meja sederhana (nama atau nomor meja) yang dikelola admin, dipilih di kasir; meja yang
  sedang punya pesanan terbuka ditandai.

---

## 15. Hal yang Belum Dibahas ([TERBUKA])

1. **Fitur POS Coffee**: sedang dibahas di bagian 14 (rincian di bawah dari catatan awal). Mana yang wajib di versi pertama (varian/topping menu, resep dan pemotongan stok bahan,
   tipe pesanan dan nomor meja, antrean/tiket dapur, open bill, split bill) dan fitur retail mana yang dibuang
   (barcode, satuan bertingkat, pembelian supplier, batch kedaluwarsa, piutang).
2. Perangkat Coffee: dipakai di tablet (printer, APK) atau cukup web; satu outlet atau banyak cabang.
3. Detail nginx: konfigurasi path, maintenance mode per aplikasi, bypass versi.
4. Mekanisme kunci aplikasi di UI "kelola aplikasi": cara membuat, ditampilkan sekali, rotasi.
5. ~~Aturan password Portal~~ **[SELESAI]**, lihat 6.9.
6. ~~Salin vs pustaka bersama~~ **[SEPAKAT]** Salin saja; tiap aplikasi berdiri sendiri, tanpa pustaka bersama
   (termasuk modul cetak struk). Perbaikan bug di satu aplikasi dicek manual di yang lain memakai daftar modul salinan.
7. ~~CI/workflow dan SOP deploy per aplikasi~~ **[SEPAKAT]** SOP deploy dipecah per aplikasi (tarik kode, build BE
   dan restart service aplikasi itu saja, build FE dengan base path dan salin ke foldernya saja). Deploy tetap
   manual; CI otomatis belum sekarang. Tidak ada pengecekan otomatis lain yang dipakai saat ini.
8. Opsi "buka aplikasi terakhir otomatis" dan logout serentak (tahap lanjutan).
9. ~~Cek hook git (Husky)~~ **[SEPAKAT, ditunda]** Husky baru diimplementasi dan belum jalan di production; tidak
   menjadi prioritas. Dahulukan Portal dan Coffee; cek Husky dilakukan seperlunya saat folder dipindah.

---

## 16. Riwayat Pembaruan

| # | Pembaruan |
|---|-----------|
| 1 | Dokumen dibuat. Memuat keputusan struktur folder, database, urutan kerja, URL berbasis path, desain SSO dengan `app_id` dan kunci aplikasi, migrasi user, role default "menunggu akses", toggle login langsung (admin saja), sesi 8 jam/7 hari, role Portal `admin`/`user`, UI kelola aplikasi (admin), penyeragaman key `<aplikasi>.<nama>` termasuk Toko (`toko.*`), dan satu APK untuk semua aplikasi. Ditambahkan usulan migrasi key Toko sekali jalan. |
| 2 | Mulai membahas fitur POS Coffee. Ditambahkan bagian 14: pemetaan modul Toko ke Coffee, kandidat fitur khas kopi, usulan tahapan (versi 1, 1.5, 2), dan 10 pertanyaan yang memengaruhi struktur database. Semua masih USULAN/TERBUKA menunggu jawaban user. |
| 3 | Usulan user: Coffee dimulai dengan menyalin BE dan FE Toko, yang diubah dulu hanya menu Kasir dan Pembelian. Ditambahkan bagian 14.5 berisi catatan AI: waktu penyalinan, pangkas modul, migrasi baseline baru, daftar yang diganti setelah salin, dan pilihan model stok bahan (P1/P2/P3). Menunggu konfirmasi user. |
| 4 | Jawaban user atas 14.5: salin dulu dari Toko (disepakati), Pembelian di Coffee = pembelian bahan baku, stok bahan memakai P3 (resep dan pengurangan otomatis masuk versi pertama), Kasir versi pertama memakai varian, tambahan, nomor meja, dan nomor antrean, daftar modul yang dipangkas disetujui. Tahapan 14.3 diperbarui, ditambahkan bagian 14.6 (jawaban user) dan 14.7 (dampak P3 serta 9 pertanyaan lanjutan). |
| 5 | Jawaban user atas 14.7: pengurangan stok dan konversi takaran tidak boleh terlalu ketat, waktu pengurangan stok dibahas lebih lanjut; dine-in punya dua opsi (bayar sekarang atau bayar nanti). Ditambahkan bagian 14.8: tiga tingkat otomatisasi stok (Ketat, Lentur, Manual) dengan rekomendasi Lentur, daftar pengaman, implikasi open bill dasar, dan tiga pertanyaan lanjutan. Keputusan 24 diubah menjadi "sepakat sebagian", keputusan 27 ditambahkan, tahapan 14.3 disesuaikan. |
| 6 | Jawaban user atas 14.8: stok bahan memakai model Lentur dengan pencatatan pemakaian manual (contoh: beli 10 bungkus, terpakai 1, sisa 9); pesanan terbuka boleh dilanjutkan ke shift berikutnya; split bill, gabung bill, dan pindah meja ikut versi pertama. Ditambahkan bagian 14.9 berisi tafsiran AI (resep tidak mengurangi stok sisa) dan 5 pertanyaan lanjutan, termasuk peran resep di versi pertama. Keputusan 24 dan 27 diperbarui, keputusan 28 dan 29 ditambahkan, tahapan 14.3 disesuaikan. |
| 7 | Jawaban user atas 14.9: stok sisa hanya berubah lewat pemakaian manual dan opname; resep versi pertama = catatan takaran standar; pemakaian dicatat kapan saja, satu bahan per input; pesanan terbuka tidak boleh lintas hari dan yang belum dibayar menjadi utang seperti di Toko; split bill dihapus. Ditambahkan bagian 14.10 (termasuk temuan scheduler kas Toko yang menutup kas tiap pukul 00:00 dan 6 pertanyaan lanjutan). `receivable` dan `customer` dikeluarkan dari daftar pangkas; keputusan 24, 26, 28, 29 diperbarui dan keputusan 30 sampai 32 ditambahkan. |
| 8 | Jawaban user atas 14.10: `receivable` dan `customer` dibawa, aturan utang sama seperti Toko, batas hari mengikuti Toko, tutup kas manual diblokir selama ada pesanan terbuka (opsi A), pemakaian stok boleh desimal, tanpa fitur impor menu/resep. Kebutuhan baru: snack dan roti dari supplier, sehingga tiap menu punya mode stok (potong langsung untuk barang jadi, manual untuk racikan). Ditambahkan bagian 14.11 (5 pertanyaan lanjutan), keputusan 23 dan 30 diperbarui, keputusan 33 sampai 38 ditambahkan, tahapan 14.3 disesuaikan. |
| 9 | Jawaban user atas 14.11: barang jadi memakai satuan dus/pcs seperti Toko; potong stok saat item masuk pesanan dan blokir saat stok habis (mengikuti rekomendasi AI); penutupan otomatis tengah malam memakai opsi i. Ditambahkan bagian 14.12, keputusan 39-42 ditambahkan dan keputusan 34 diperbarui. |
| 10 | Jawaban user atas pertanyaan lanjutan 14.12: barang jadi tanpa varian (menu terpisah per rasa), pembatalan setelah barang diserahkan dipilih kasir (kembali ke stok atau dibuang), aturan pendamping opsi i disetujui. Keputusan 43 sampai 45 ditambahkan. |
| 11 | Jawaban user: bahan terbuang dicatat dengan mekanisme yang sama seperti barang jadi; nomor antrean direset tiap hari. Pertanyaan layar/printer bar dan foto menu/cabang/perangkat dijelaskan ulang. Ditambahkan bagian 14.13. |
| 12 | Jawaban user: tanpa tiket/layar bar, menu tombol teks tanpa foto, satu kedai. Perangkat Coffee belum dipahami user, dijelaskan ulang. Bagian 14.13 diperbarui. |
| 13 | Jawaban user: Coffee dapat dibuka lewat browser maupun APK tablet (diakses sebagai website). Bagian 14.13 diperbarui; pertanyaan 14.4 nomor 8-10 selesai. |
| 14 | Diskusi nginx/VPS: acuan `nginx/pos-web.conf`, tiga FE di folder terpisah, maintenance satu saklar semua tapi bisa per aplikasi, cutover rename + clone ulang. Penjelasan 8081 (Mode B) ditambahkan. Ditambahkan bagian 10.1. |
| 15 | Jawaban user: nginx VPS sama persis dengan repo; pertanyaan jumlah service dijawab (usulan 3 service, satu per BE); Mode B untuk aplikasi baru dijelaskan ulang. Bagian 10.1 diperbarui. |
| 16 | Jawaban user: Mode B (versi lama/baru) dihapus, hanya satu versi. Port dan service ditetapkan: Portal `pos-portal-backend` 8080, Toko `pos-toko-backend` 8081, Coffee `pos-coffee-backend` 8082 (menggantikan 8083/8080/8082 sebelumnya). Bagian 3 dan 10.1 diperbarui. |
| 17 | Hasil cek VPS (screenshot user): nginx di 80 dan 443, sertifikat certbot ada, `pos_api` di 8080, 8081 bebas. Perbedaan repo vs server soal blok SSL dicatat sebagai terbuka di bagian 10.1. |
| 18 | User mengonfirmasi SSL memakai certbot. Bagian 10.1 diperbarui. |
| 19 | Mulai bahas aturan password Portal; ditambahkan bagian 6.9 berisi fakta aturan password Toko. |
| 20 | Jawaban user soal password Portal: usulan 1-4 dan 6 disetujui (min 8 karakter, ganti sendiri, reset admin dengan password sementara, lupa lewat admin, password Toko terpisah), poin 5 ditolak: tanpa batas percobaan login. Bagian 6.9 diperbarui. |
| 21 | Jawaban user: kode Toko dan Coffee disalin, tanpa pustaka bersama (aplikasi berdiri sendiri). Bagian 15 diperbarui. |
| 22 | Jawaban user: SOP deploy dipecah per aplikasi, deploy tetap manual tanpa CI; Husky belum berjalan di production dan ditunda, prioritas Portal dan Coffee. Bagian 15 diperbarui. |
| 23 | Jawaban user: kunci aplikasi disimpan di database aplikasi (tabel baru), bukan `.env`. Ditambahkan bagian 6.10. |
| 24 | Jawaban user: kunci diisi lewat SQL manual; `APP_ID` ikut pindah ke tabel kunci di database aplikasi. Keputusan 9 dan bagian 6.10 diperbarui. |
| 25 | Implementasi Fase 1 (di repo): `BE/` dan `FE/` dipindah ke `TOKO/BE` dan `TOKO/FE`; `.gitignore` memakai pola `**/BE/...` dan `**/FE/...`; port dev/prod Toko 8081 (`.env*`, proxy Vite, `.env.development`); `nginx/pos-web.conf` menjadi sumber tunggal (Mode B dan `ver_bypass` dihapus, `proxy_pass` ke 8081), `nginx/pos-web-old.conf` dan `TOKO/FE/nginx.conf` dihapus; README, `mobile/README.md`, `docs/SOP_DEPLOY.md`, `docs/DEPLOYMENT_PROD.md`, `docs/DEPLOYMENT_REDEPLOY_FULL.md` diperbarui. Verifikasi: `go build ./...`, `type-check`, `lint`, `build` FE semua lolos. Belum: commit, deploy VPS, uji jalan BE (MySQL). |
| 26 | Port dev diselaraskan: BE Portal 8080, Toko 8081, Coffee 8082; FE Portal 3000, Toko 3001, Coffee 3002. Port FE Toko diubah dari 3000 ke 3001 (vite.config, APP_HOST di .env, CORS dev, README). |
| 27 | Usul user: Portal dibuat dengan menyalin Toko lalu memangkas, desain dirombak total dan bebas diserahkan ke AI, implementasi bertahap. Ditambahkan keputusan 46-48. Tanggapan AI dan tahapan diberikan di chat (belum implementasi). |
| 28 | Tahapan Fase 2 (2.1-2.7) dicatat di bagian 5.0. User: database `pos_portal_db` sudah siap, migrasi dibuat AI, akun awal samakan dengan seed Toko. Koreksi fakta: hash password Toko adalah Argon2id (package `pkg/bcrypt` hanya namanya), bukan bcrypt; bagian 7, 13, dan tabel Portal diperbarui. |
| 29 | Tahap 2.1 selesai: `PORTAL/BE` disalin dari `TOKO/BE` (module `portal_api`, port 8080, DB `pos_portal_db`), domain tersisa `auth`, `user`, `role` (role belum dirouting), scheduler/cmd/deploy/permission dibuang, `/users` hanya untuk role `admin`, migrasi 001 (roles, users, sessions, log_*) dan 002 (role `admin`/`user`, akun `owner` dan `admin` seperti seed Toko, keduanya role `admin`). Verifikasi: `go build`, `go vet`, BE hidup, `/api/health` 200, login `admin` berhasil. Ditemukan: pembatas login Toko ada (5 gagal/15 menit); bagian 6.9 poin 5 dibuka kembali. |
| 30 | Keputusan user: Portal memakai pembatas login seperti Toko (5 gagal/15 menit); bagian 6.9 poin 5 diperbarui. Tidak ada perubahan kode (pembatas sudah ikut tersalin). |
| 31 | Tahap 2.2 selesai (BE Portal). Migrasi 003 menambah `users.must_change_password`. Endpoint baru `POST /api/auth/change-password` (password lama wajib, baru minimal 8 dan harus berbeda). Admin reset lewat `POST /api/users/reset-password/:id` (menggantikan `change-password/:id`): menetapkan password sementara, menandai wajib ganti, menghapus sesi user itu, dan ditolak untuk akun sendiri. Selama wajib ganti, token hanya boleh ke `/auth/me`, `/auth/logout`, `/auth/change-password` (ditegakkan di middleware auth). Password baru (buat user dan reset) minimal 8 karakter. Istilah `owner` dibersihkan (admin terakhir dihitung dari role `admin`). Domain `role` hanya menyisakan `POST /api/roles/options` (admin). Terverifikasi lewat 16 skenario API; build dan vet lolos. Catatan: user baru yang dibuat admin belum ditandai wajib ganti (hanya hasil reset). |
| 32 | Tahap 2.3 selesai (FE Portal): `PORTAL/FE` disalin dari `TOKO/FE` (package `portal-web`, port 3000, proxy ke 8080, key browser `portal.auth-session`, env produksi `VITE_API_URL=/portal-api`), fitur Toko, sistem menu/permission, cetak, dan dependency berat (recharts, xlsx, jspdf, react-barcode, husky) dibuang. Sistem desain baru (indigo/violet, token di `styles/tokens.css`, animasi di `tailwind.config.ts`). Halaman: login (panel branding + form), ganti password (wajib bila password sementara; route lain dikunci oleh `ProtectedRoute`), `/apps` (placeholder, diisi tahap 2.4), layout `PortalLayout` dengan menu akun. Terverifikasi: `type-check`, `lint`, `build` lolos; uji UI otomatis (login admin, login salah, password sementara dipaksa ganti, validasi form, logout, tanpa overflow di 390px). Catatan: nginx produksi nanti memetakan `/portal-api/` ke `http://127.0.0.1:8080/api/`. |
| 33 | Tahap 2.5 selesai (BE aplikasi dan SSO), dikerjakan sebelum 2.4. Migrasi 004: `apps` (kunci disimpan sebagai hash SHA-256), `user_apps`, `sso_codes` (hanya hash kode); seed aplikasi `toko` (path `/toko/`, kunci kosong) dan akses akun `owner`/`admin`. Endpoint: `POST /api/apps/my` (aplikasi aktif milik user), `POST /api/sso/launch` (kode acak 32 byte, berlaku 60 detik, mengembalikan `redirect_url` = `<path aplikasi>sso?code=...`), `POST /api/sso/exchange` (publik, header `X-App-Key`, body `{app_id, code}`, balasan `{valid, app_id, user{portal_user_id, username, full_name, is_active}}`). Penolakan seragam 401. Terverifikasi 22 skenario: sekali pakai (3 permintaan paralel, hanya 1 lolos), kedaluwarsa, kunci/app_id salah, aplikasi atau user nonaktif, akses dicabut setelah launch. User baru buatan admin kini wajib ganti password saat login pertama (keputusan 49). Catatan: nama kolom kunci di Portal berbeda dari sisi aplikasi (`portal_credentials` baru dibuat di Fase 3); nilai `path` boleh URL penuh untuk lingkungan dev. |
| 34 | Tahap 2.4 selesai (FE): `/apps` menampilkan kartu aplikasi dari `POST /apps/my` (ikon dari kolom `icon`: `store`, `coffee`, lainnya ikon umum), klik kartu memanggil `POST /sso/launch` lalu pindah di tab yang sama ke `redirect_url`. Ada skeleton loading, state error dengan tombol "Coba lagi", dan state kosong untuk user tanpa akses. Terverifikasi: `type-check`, `lint`, `build`, dan uji UI otomatis (2 kartu tampil, klik mengarah ke `<alamat aplikasi>/sso?code=<64 karakter>`, error lalu coba lagi, state kosong). Data uji sudah dikembalikan. |
| 35 | Tahap 2.6a selesai (BE admin Portal), semua endpoint khusus role `admin` (non-admin 403). Aplikasi: `POST /api/apps/list`, `create` (kode huruf kecil/angka/strip, alamat berupa path `/...` atau URL http/https, menolak `//`; mengembalikan kunci `pak_<64 hex>` SEKALI, yang tersimpan hanya hash), `update/:id` (kode tidak bisa diubah), `toggle-status/:id`, `rotate-key/:id` (kunci lama langsung mati, kunci baru tampil sekali). Akses user: `POST /api/apps/access/:user_id` (semua aplikasi + status granted) dan `/access/:user_id/set` (mengganti seluruh akses secara transaksional, id ganda dibuang). Tidak ada endpoint hapus aplikasi (cukup nonaktifkan). Terverifikasi 24 skenario, termasuk kunci hasil create/rotate dipakai di `/sso/exchange`. Catatan: id aplikasi yang tidak ada pada `set` diabaikan diam-diam, bukan ditolak. |
| 36 | Tahap 2.6b selesai (FE admin Portal). Header memuat menu "Aplikasi / Kelola User / Kelola Aplikasi" hanya untuk admin (di layar kecil masuk ke menu akun); route `/admin/users` dan `/admin/apps` dibungkus `ProtectedRoute roles=[admin]` (non-admin dialihkan ke `/apps`). Kelola Aplikasi: tabel, tambah (kunci `pak_...` tampil sekali dengan tombol Salin), ubah (kode terkunci), aktif/nonaktif, buat kunci baru (dengan konfirmasi). Kelola User: tabel dengan pencarian dan paginasi, tambah (password awal dengan tombol "Acak"), ubah, aktif/nonaktif, reset password sementara, dialog akses aplikasi (centang), hapus; akun sendiri tidak punya aksi reset/nonaktif/hapus. Terverifikasi: `type-check`, `lint`, `build`, dan uji UI otomatis 15 skenario (termasuk salin kunci ke clipboard, hash di DB bukan kunci asli, akses tersimpan, non-admin tidak melihat menu admin dan dialihkan dari URL admin). Data uji dibersihkan. |
| 37 | Fase 3 dimulai atas persetujuan user; tahap 2.7 digabung ke Fase 3. Tahapan 3.1-3.6 dicatat di bagian 5.0b. |
| 38 | Tahap 3.1 selesai (BE Toko, aditif). Migrasi `011_portal_sso.sql`: `users.portal_user_id` (unik, nullable), role `menunggu_akses` (tanpa izin menu), tabel `portal_credentials` (satu baris: `app_id`, `app_key`, diisi lewat SQL), setting `direct_login_enabled=true`. Config baru: `PortalUrl` dan `DirectLoginForceEnabled` (flag darurat, dipakai tahap 3.2). `POST /api/auth/sso {code, device_info?}` (publik): menukar kode ke Portal lewat `pkg/portal` (server ke server, header `X-App-Key`), memastikan jawaban untuk `app_id` aplikasi ini, mencari profil lokal lewat `portal_user_id` (dibuat otomatis dengan role `menunggu_akses` dan password acak tak dikenal bila belum ada), menyegarkan nama dari Portal, menolak akun lokal nonaktif, lalu membuat sesi lewat `issueSession` (hasil refactor Login, perilaku tidak berubah). Akun lokal dengan username sama yang belum ditautkan TIDAK otomatis dipakai (ditolak, harus lewat skrip migrasi tahap 3.4). Terverifikasi 12 skenario dengan kode asli dari Portal; login langsung lama tetap 200. Catatan: DB dev Toko sudah ter-migrasi 011. Penamaan: dokumen 6.5 menyebut `APP_ID` di `.env`; sesuai keputusan 9 dan 6.10, `app_id` dan kunci kini di tabel `portal_credentials`. |
| 39 | Keputusan user: toggle login langsung di Toko hanya boleh diubah role `admin` (role `owner` tidak). Tahap 3.2 dimulai. |
| 40 | Tahap 3.2 selesai (BE Toko). `GET /api/auth/login-options` (publik) mengembalikan `direct_login_enabled`, `portal_login_enabled`, `direct_login_forced`. `POST /api/auth/direct-login {enabled}` hanya role `admin` (owner 403). Saat off, `POST /auth/login` menolak dengan 403 "Login langsung dinonaktifkan"; SSO tetap jalan. Pengaman anti-terkunci: flag darurat `DirectLoginForceEnabled` di config server mengalahkan setelan DB; hanya nilai `false` yang mematikan (kosong/hilang/rusak = aktif); bila login Portal belum siap (alamat Portal atau `portal_credentials` kosong) login langsung tetap aktif dan permintaan mematikannya ditolak 400. Key `direct_login_enabled` dilindungi dari endpoint pengaturan umum: `POST /settings` yang memuat key itu ditolak 400 tanpa menyimpan apa pun, dan Reset pengaturan tidak menimpa nilainya. Terverifikasi 17 skenario API ditambah uji flag darurat (config dikembalikan ke `false` sesudahnya). |