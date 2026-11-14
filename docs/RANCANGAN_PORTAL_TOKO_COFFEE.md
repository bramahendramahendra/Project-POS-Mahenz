# Rancangan Portal + POS Toko + POS Coffee

> Status: **Fase 1-3 (Portal dan Toko) sudah diimplementasi; Fase 4 (Coffee) masih rancangan.**
> Status tiap tahap ada di bagian 5.0 dan 5.0b. Dokumen ini catatan hidup: setiap kali ada diskusi atau
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
| 24 | Stok bahan baku | Model Lentur dengan pencatatan pemakaian manual: stok sisa = pembelian dikurangi pemakaian yang diinput ~~, dikoreksi lewat opname~~ (**diperbarui keputusan 84: tanpa opname**). Resep hanya catatan takaran standar dan tidak mengurangi stok sisa | SEPAKAT |
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
| 41 | Stok barang jadi habis | Penjualan diblokir; ~~ketidaksesuaian stok dikoreksi lewat opname~~ (opname dihapus, keputusan 84: stok lewat pembelian, kekurangan dicatat sebagai terbuang) (mengikuti rekomendasi AI) | SEPAKAT |
| 42 | Penutupan otomatis tengah malam dengan pesanan terbuka | Opsi i: kas yang masih punya pesanan terbuka tidak ditutup otomatis dan ditandai untuk diselesaikan | SEPAKAT |
| 43 | Varian barang jadi | Satu menu = satu stok, tanpa varian; rasa berbeda = menu terpisah. Tambahan boleh, tidak memengaruhi stok | SEPAKAT |
| 44 | Pembatalan setelah barang diserahkan | Kasir memilih "kembali ke stok" atau "dibuang" | SEPAKAT |
| 45 | Aturan pesanan melewati batas hari (pendamping opsi i) | Hanya bisa diselesaikan, tanpa tambah item; buka kas baru diblokir sampai kas tertunda ditutup; admin bisa menyelesaikan atas nama; ada daftar pesanan tertunda | SEPAKAT |
| 49 | User baru buatan admin | Wajib ganti password saat login pertama (sama seperti hasil reset); diserahkan user ke AI | SEPAKAT |
| 50 | Urutan tahap 2.4/2.5 | Tahap 2.5 (BE aplikasi dan SSO) dikerjakan lebih dulu, baru 2.4 (halaman pilih aplikasi) agar memakai data asli; user menyerahkan urutan ke rekomendasi AI | SEPAKAT |
| 46 | Cara membangun Portal | Usul user: struktur dan konsistensi kode BE dan FE Portal wajib sama dengan Toko; caranya menyalin `TOKO/BE` dan `TOKO/FE` ke `PORTAL/`, menghapus fitur yang tidak relevan, dan desain UI/UX diubah total. Implementasi dibuat bertahap | SEPAKAT (usulan user, ditanggapi AI) |
| 47 | Desain Portal | Diserahkan ke AI (bebas); target tampilan modern dan canggih | SEPAKAT |
| 52 | Pemetaan role Toko ke Portal saat migrasi user | `admin` Toko menjadi `admin` Portal; semua role lain (owner, kasir, dst) menjadi `user` Portal | SEPAKAT |
| 51 | Siapa boleh mengubah toggle login langsung | Hanya role `admin` (di Toko juga: `owner` tidak boleh) | SEPAKAT |
| 48 | Skrip salin user Toko ke Portal | Dijalankan sebagai perintah terpisah (bukan otomatis saat Portal hidup); user menyerahkan keputusan ke AI | SEPAKAT |
| 53 | Varian menu Coffee | Dua pilihan tetap: **ukuran × panas/dingin**; harga per kombinasi; kombinasi yang tidak dijual bisa dinonaktifkan; pilihan lain (gula, shot, susu) lewat tambahan; menu tanpa varian (roti, snack) satu harga | SEPAKAT |
| 54 | Tambahan (modifier) | Kelompok tambahan dipasang ke menu (mis. "Level gula": pilih satu, gratis; "Topping": boleh banyak, berbayar) | SEPAKAT (rekomendasi AI) |
| 55 | Resep | Catatan takaran per kombinasi varian, opsional | SEPAKAT (rekomendasi AI) |
| 56 | Antrean dan meja | Nomor antrean urut harian 001, 002, ... tanpa awalan, reset 00:00 Asia/Jakarta; meja = daftar sederhana (nama/nomor) dikelola admin | SEPAKAT |
| 57 | Diskon, pajak, service charge | Diskon dan pajak seperti kasir Toko; service charge tidak dulu | SEPAKAT (rekomendasi AI) |
| 58 | Role awal Coffee | `admin`, `owner`, `kasir` seperti Toko; toggle login langsung hanya `admin` (keputusan 51) | SEPAKAT |
| 59 | Open bill | Berlaku untuk dine-in dan take away | SEPAKAT (rekomendasi AI) |
| 60 | Cara salin Toko ke Coffee (14.5 poin 1, 3, 4, 6) | Disalin sekarang (Fase 4); migrasi baseline baru; identitas diganti; modul salinan dicatat (bagian 3.1) | SEPAKAT |
| 61 | Urutan kerja | Bangun Coffee di dev sekarang; **seluruh urusan deploy VPS (Portal, Toko, Coffee) dikerjakan belakangan** | SEPAKAT |
| 62 | Akun awal Coffee | Sama seperti Toko: seed akun lokal `owner` (`owner123`) dan `admin` (`admin123`), password bawaan wajib diganti sebelum production; penautan ke Portal memakai skrip `migrate_users_to_portal` yang ikut disalin ke Coffee (aturan sama dengan 3.4: dry-run default, `--apply --i-have-backup`, `admin` ke `admin` Portal, lainnya `user`, akses aplikasi `coffee`, KONFLIK bila password berbeda dan tidak menebak). Memperbarui §7 ("Coffee tidak punya akun lama"). **Diperbarui oleh 63-72:** seed profil lokal `owner`/`admin` tetap, tetapi skrip tautan tidak diperlukan untuk Coffee (tautan otomatis, keputusan 67) | SEPAKAT, diperbarui |
| 63 | Layanan akun terpisah ("BE User") | Aplikasi keempat `AUTH/BE` (tanpa FE), database `pos_auth_db`, service `pos-auth-backend`; **hanya mendengarkan `127.0.0.1`, tidak dipasang di nginx**. FE tidak pernah menghubungi BE User; FE Portal/Toko/Coffee ke BE masing-masing, lalu BE aplikasi ke BE User memakai kunci aplikasi masing-masing (Portal juga punya kunci). Data akun (username, nama, password, aktif, wajib ganti password, akses aplikasi, kunci aplikasi, kode SSO, pembatas login) pindah dari Portal ke BE User. Setiap aplikasi tetap membuat sesi/token sendiri dan menyimpan role, PIN, riwayat sendiri. Muncul dari masalah KONFLIK tautan admin Coffee (riwayat #55) | SEPAKAT |
| 64 | Port BE baru | BE User **8080**, Portal **8081**, Toko **8082**, Coffee **8083**; FE tetap Portal 3000, Toko 3001, Coffee 3002. Menggantikan port di bagian 3 dan 10.1 | SEPAKAT |
| 65 | Role pengelola akun | Role `admin`/`user` (semula milik Portal) disimpan di BE User | SEPAKAT |
| 66 | Bila BE User mati | **Tanpa jalur darurat**: tidak ada login baru di aplikasi mana pun sampai BE User hidup (yang sudah login tetap bekerja sampai sesinya habis). Password lokal aplikasi tidak dipakai lagi. Risiko diterima user demi satu tempat pengelolaan akun | SEPAKAT |
| 67 | Tautan otomatis | Login pertama yang **terverifikasi BE User** menautkan otomatis ke profil lokal dengan username sama yang **belum tertaut**. Memperbarui "Ditolak: auto-link username" (bagian 4 HANDOFF / riwayat #38) untuk kasus terverifikasi ini | SEPAKAT |
| 68 | Toggle login langsung | **Dipertahankan**, hanya `admin` yang mengubah (keputusan 51). Nyala: form username/password tampil (diperiksa BE User lewat BE aplikasi) + tombol "Masuk lewat Portal". Mati: hanya tombol Portal | SEPAKAT |
| 69 | Kelola akun | Buat user, ganti password, reset password **hanya lewat Portal** (Portal meneruskan ke BE User). Toko/Coffee hanya mengatur role, PIN, aktif/nonaktif profil lokal | SEPAKAT |
| 70 | Wajib ganti password | Login di aplikasi ditolak dengan pesan "Ganti password dulu di Portal" | SEPAKAT |
| 71 | Kolom tautan | `users.portal_user_id` di Toko/Coffee diganti nama `auth_user_id` (aman: Toko baru belum di-deploy) | SEPAKAT |
| 72 | Urutan kerja | Fase 3B (BE User) dikerjakan **sebelum Fase 4 dilanjutkan**; tahap 4.1 Coffee diulang dengan menyalin Toko versi baru (hasil 4.1 sebelumnya, mis. skema dan seed, dipakai ulang sebagai bahan) | SEPAKAT (rancangan 5.0d disetujui user) |
| 74 | Password di log permintaan | Ke depan: penyamaran rahasia di log (`helper/logmask`, seperti BE User) dipasang di Portal (3B.3), Toko (3B.4), dan Coffee (saat 4.1 diulang). Log lama yang sudah berisi password **dibiarkan** (tidak dibersihkan, baik dev maupun production), tanpa saran ganti password ke kasir | SEPAKAT |
| 75 | Hapus user di Toko/Coffee | **Tetap ada**: hanya menghapus profil lokal aplikasi (akun BE User/Portal tetap). Bila orangnya login lagi, profil baru `menunggu_akses` dibuat. Aturan lama tetap (bukan diri sendiri, bukan admin/owner aktif terakhir) | SEPAKAT (rekomendasi AI) |
| 76 | Edit user di Toko/Coffee | **Hanya role.** Nama lengkap dikelola di Portal/BE User dan disalin ke profil lokal setiap kali user login (form langsung maupun SSO) | SEPAKAT (rekomendasi AI) |
| 77 | Tautan lama di dev | Migrasi 012 Toko: `portal_user_id` diganti nama `auth_user_id` dan **dikosongkan** (isinya id Portal lama, tidak sama dengan id BE User); tautan dibuat ulang otomatis saat login pertama yang terverifikasi BE User (keputusan 67). Production tidak terdampak (Toko baru belum di-deploy); DB dev di-backup dulu | SEPAKAT (rekomendasi AI) |
| 78 | Data dev Portal (3B.5) | **Tidak dipindah** ke `pos_auth_db` (Portal belum pernah di-deploy, isinya hanya data uji). Akun penting datang dari skrip salin user Toko (yang juga dipakai saat cutover production); akun uji khusus Portal ditinggalkan. Memperbarui rencana 3B.5 | SEPAKAT (rekomendasi AI) |
| 79 | Kolom/tabel lama | **Dihapus sekarang** lewat migrasi (backup DB dev dulu): Portal `users.password`, `users.role_id`, `users.must_change_password`, tabel `roles`, `apps`, `user_apps`, `sso_codes`; Toko tabel `portal_credentials`. Kolom `users.password` Toko **tetap** (dibaca skrip salin saat cutover production) | SEPAKAT (rekomendasi AI) |
| 80 | Input pemakaian dan barang terbuang (bahan baku) | Boleh oleh **kasir, admin, owner** (izin menu "Riwayat Stok Bahan", aksi tambah); kelola daftar bahan oleh admin/owner. Dapat diubah per role lewat pengaturan izin menu | SEPAKAT (rekomendasi AI) |
| 81 | Koreksi salah input | **Admin/owner membatalkan entri** pemakaian/terbuang dengan alasan: stok dan total terpakai/terbuang kembali, entri asli tetap terlihat bertanda dibatalkan (siapa, kapan, alasan) dan pembatalan ikut tercatat di riwayat. Entri tidak bisa diubah. Opsi koreksi lewat opname ditolak (merusak angka total terpakai) | SEPAKAT (rekomendasi AI) |
| 82 | Satuan stok bahan | Dipilih dari **master Unit** yang sudah ada | SEPAKAT (rekomendasi AI) |
| 83 | Data bahan baku | Nama, satuan, stok sisa, total terpakai, total terbuang, status aktif, **batas stok menipis** (opsional, tanda "menipis"), **alasan wajib** saat barang terbuang. Tanpa kategori bahan. **Tanpa kolom stok awal**: stok bahan baru 0 | SEPAKAT |
| 84 | Opname | **Dihapus** untuk bahan baku dan barang jadi. Stok = pembelian - terpakai/terjual - terbuang (dengan pembatalan entri). Konsekuensi diterima user: stok yang sudah ada saat mulai memakai aplikasi masuk lewat pembelian, kelebihan stok fisik tidak bisa ditambahkan selain lewat pembelian, kekurangan dicatat sebagai barang terbuang. Memperbarui keputusan 24, 14.11, 14.12, tahap 4.3 dan 4.5 | SEPAKAT |
| 85 | Ukuran varian menu | **Daftar ukuran dikelola admin** (mis. Regular, Large); tiap menu memilih ukuran yang dijual. Panas/dingin tetap dua pilihan (keputusan 53) | SEPAKAT (rekomendasi AI) |
| 86 | Bentuk resep | **Baris bahan + takaran**: pilih bahan dari Daftar Bahan + takaran teks bebas (mis. "18 gram"), per kombinasi varian, opsional, tidak mengurangi stok (keputusan 24, 55) | SEPAKAT (rekomendasi AI) |
| 87 | Aturan kelompok tambahan | Nama, pilih satu atau boleh banyak, **wajib atau opsional**, pilihan dengan harga (0 = gratis), **satu pilihan bawaan** (opsional). Tanpa batas maksimal pilihan (bisa ditambah nanti) | SEPAKAT (rekomendasi AI) |
| 88 | Izin menu jualan | Kelola kategori, menu, ukuran, kelompok tambahan: admin dan owner. Kasir: melihat daftar menu dan **menandai habis/tersedia** | SEPAKAT (rekomendasi AI) |
| 89 | Barang jadi di 4.4 | Menu mode langsung di 4.4 hanya pilihan mode dan harga (tanpa varian, keputusan 43); stok dus/pcs dan pembelian di 4.5. Menu "Produk" salinan Toko dibiarkan sampai 4.5 lalu diputuskan | SEPAKAT (rekomendasi AI) |
| 90 | Perbaikan Combobox Toko | Pencarian Combobox Toko memakai label, bukan id (sama dengan perbaikan Coffee di 4.3) | SEPAKAT |
| 91 | Katalog barang jadi | Modul **Produk** salinan Toko menjadi katalog barang jadi (satuan bertingkat dus/pcs, stok, pembelian, riwayat mutasi Toko). Menu bermode barang jadi **ditautkan ke tepat satu produk** (satu produk hanya untuk satu menu); harga jual diatur di Menu; stok berkurang dari produk saat dijual (4.7). Grup menu navigasi "Produk" diganti nama "Barang Jadi". Alternatif stok langsung di menu ditolak | SEPAKAT (rekomendasi AI) |
| 92 | Pembelian dua jenis | Satu nota pembelian boleh berisi baris **bahan baku** dan **barang jadi**. Baris bahan: pilih bahan, jumlah dalam satuan stok bahan (desimal), harga per satuan; stok bahan bertambah dan tercatat di Riwayat Stok Bahan sebagai "Pembelian". Baris barang jadi: seperti Toko (produk + kemasan). Harga beli bahan disimpan untuk HPP nanti | SEPAKAT (rekomendasi AI) |
| 93 | Bayar, void, edit pembelian | Sama seperti Toko (utang supplier, bayar bertahap, void, edit = stok lama dibatalkan lalu stok baru diterapkan). Void boleh membuat stok bahan minus (model lentur) | SEPAKAT (rekomendasi AI) |
| 94 | Batas tahap 4.5 | Form Produk salinan Toko tidak diubah (hanya keterangan harga jual diatur di Menu); kasir Toko yang masih ada di Coffee belum disentuh (dibangun ulang di 4.7) | SEPAKAT (rekomendasi AI) |
| 95 | Siapa mendapat nomor antrean | **Semua pesanan** (dine-in dan take away) mendapat nomor antrean harian; dine-in juga memakai meja. Nomor di struk dipakai memanggil/mengantar (tanpa layar bar) | SEPAKAT (rekomendasi AI) |
| 96 | Kapan nomor antrean diambil | **Saat pesanan dibuat** (termasuk open bill). Pesanan batal tetap memakai nomornya (dilewati, tidak dipakai ulang); pesanan yang gagal disimpan tidak memakai nomor | SEPAKAT (rekomendasi AI) |
| 97 | Data meja | **Nama/nomor unik + urutan tampil + aktif**; tanpa area dan kapasitas | SEPAKAT (rekomendasi AI) |
| 98 | Letak dan izin halaman Meja | Item **"Meja" di grup navigasi Menu** (setelah Ukuran); admin dan owner mengelola, kasir hanya melihat (memilih meja di kasir, 4.7) | SEPAKAT (rekomendasi AI) |
| 99 | Nama pemesan | Kolom **nama pemesan opsional** di pesanan, tampil di struk dan riwayat transaksi di samping nomor antrean | SEPAKAT (rekomendasi AI) |
| 100 | Void transaksi lunas (barang jadi) | Usul user: void memberi **pilihan per item barang jadi "kembali ke stok" atau "dibuang"** (bawaan kembali ke stok), **alasan wajib** bila dibuang; dibuang tercatat sebagai mutasi stok produk jenis baru "Terbuang" (`waste`). Mekanisme yang sama dipakai batal pesanan di 4.8 (keputusan 44). Barang racikan tanpa pilihan (stok bahan dicatat manual) | SEPAKAT (usulan user, cara pilih rekomendasi AI) |
| 101 | Utang saat bayar sekarang | Pilihan utang (pilih pelanggan) di pembayaran **dipertahankan seperti Toko** (keputusan 30, 36) | SEPAKAT (rekomendasi AI) |
| 102 | Pembagian tahap 4.7 | Dipecah **4.7a** (BE pesanan + uji API) dan **4.7b** (kasir FE, struk, riwayat + uji browser), masing-masing dilaporkan dan disetujui | SEPAKAT (rekomendasi AI) |
| 73 | Flag darurat `DirectLoginForceEnabled` | **Dibuang.** Setelan toggle login langsung tetap per aplikasi di tabel `settings` (`direct_login_enabled`) database aplikasi itu. Konsekuensi diterima: bila toggle mati dan halaman Portal tidak bisa dibuka, pemulihan lewat SQL `UPDATE settings SET setting_value = 'true' WHERE setting_key = 'direct_login_enabled';` di database aplikasi (dicatat di dokumentasi deploy) | SEPAKAT |

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

### 3.1 Standar konsistensi kode antaraplikasi [SEPAKAT]

Toko, Portal, dan Coffee wajib sama dalam bahasa, framework, versi pustaka, struktur folder, dan gaya kode.
Yang boleh berbeda hanya fitur sesuai kebutuhan tiap aplikasi. Kode disalin per aplikasi (tanpa pustaka
bersama), sehingga perbaikan bug di satu aplikasi wajib dicek manual di aplikasi lain.

| Lapisan | Standar (terverifikasi sama di Toko dan Portal) |
|---------|-------------------------------------------------|
| BE stack | Go 1.24.5, Gin 1.11.0, GORM 1.31.1, golang-jwt v5, Viper 1.21.0; module `pos_api` / `portal_api` / `coffee_api` |
| BE folder | `config`, `database/migrations`, `domain`, `dto`, `errors`, `helper`, `middleware`, `model`, `pkg`, `repository`, `routes`, `server`, `validation`; tambahan per kebutuhan (Toko: `scheduler`, `cmd`, `deploy`, `backups`) |
| BE domain | Tiap domain berisi `dto/`, `handler/`, `model/`, `repo/`, `service/` (repo dan service memakai interface) |
| BE config | Port dari `APP_PORT` di `.env`; database dan setelan lain di `config/config_{dev,prod}.json` |
| FE stack | React 19, Vite 8, TypeScript 6, TanStack Query 5, Zustand 5, React Router 7, Tailwind 3, axios, zod, react-hook-form; ESLint flat config `eslint.config.js` |
| FE folder | `src/app` (App, router, providers, LazyRoute), `src/features/<fitur>`, `src/shared/{components,constants,hooks,types,utils}`, `src/services` (api client), `src/styles` |
| FE fitur | Pola file `<fitur>.api.ts`, `.schema.ts`, `.store.ts`, `.types.ts`, `components/`, `hooks/`, `index.ts` |
| Browser | Key localStorage berawalan nama aplikasi (`toko.*`, `portal.*`, `coffee.*`) |

Daftar modul infrastruktur yang disalin antaraplikasi (acuan cek bug lintas aplikasi): `auth`, `user`, `role`,
middleware auth/role/pembatas login, `pkg/bcrypt` (Argon2id), `errors`, `helper`, `validation`, `services/api.client.ts`,
`services/authClient.ts`, `shared/components/ui`. Daftar ini diperluas saat Coffee dibuat.

Perbedaan yang disengaja (desain Portal berbeda, keputusan 47): di Portal, komponen `ui/dialog` (tombol tutup 40px,
judul diberi ruang), `ui/dropdown-menu` (item min. 40px), `ui/password-input` (tombol mata 40px), `DataTable`
(prop `cardBreakpoint`, aksi card di kanan atas), dan paginasi (tombol 40px) sudah disesuaikan untuk sentuh (riwayat #49).
Toko memakai pendekatan berbeda agar tampilan laptop production tidak berubah: varian Tailwind `touch:` (`@media (pointer: coarse)`)
yang hanya memperbesar target sentuh di HP/tablet/APK, plus kolom aksi `DataTable` yang dibekukan di kanan (riwayat #50).
**[USULAN]** Coffee kelak mengikuti pola Toko (varian `touch:` dan kolom aksi beku), karena Coffee juga dipakai di tablet.

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
| 1 | Pindahkan BE/FE ke `TOKO/` | **[KODE SELESAI dan sudah di-commit user, belum di-deploy]** Tanpa perubahan fungsi, tetapi port Toko menjadi 8081 dan service menjadi `pos-toko-backend`. Deploy ke VPS dan verifikasi POS Toko jalan sama seperti sebelumnya |
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
| 2.1 **[SELESAI]** | Kerangka BE Portal: salin `TOKO/BE`, pangkas, ganti identitas, migrasi awal ke `pos_portal_db` | BE hidup di 8080, `/api/health` membalas, tabel terbentuk |
| 2.2 **[SELESAI]** | Auth BE: login, sesi, user, role `admin`/`user`, seed akun awal | Login dan refresh/logout lewat API |
| 2.3 **[SELESAI]** | Kerangka FE Portal: salin `TOKO/FE`, pangkas, sistem desain baru, halaman login | Login berjalan di 3000 dengan tampilan baru |
| 2.4 **[SELESAI]** | Halaman pilih aplikasi: kartu aplikasi, menu pengguna | Daftar aplikasi tampil setelah login. **Dikerjakan setelah 2.5** (urutan ditukar agar memakai data asli, bukan daftar sementara) |
| 2.5 **[SELESAI]** | Aplikasi dan akses (BE): tabel `apps`, `user_apps`, `sso_codes`, kunci, `/sso/exchange` | Kode SSO dibuat dan ditukar lewat API |
| 2.6 **[SELESAI]** | Admin Portal: kelola user, kelola aplikasi, atur akses, reset password | Admin mengatur semuanya lewat UI. Dibagi: 2.6a BE, 2.6b FE |
| 2.7 **[DIGABUNG ke Fase 3]** | Ganti password sendiri (sudah ada sejak 2.2) dan toggle login langsung (tahap 3.2-3.3) | Fitur akun berjalan |

Keputusan teknis tahap 2.1-2.2:
- Database `pos_portal_db` sudah disiapkan user; migrasi dibuat oleh AI.
- Akun awal Portal semula sama dengan seed Toko (`owner` dan `admin`, hash Argon2id sama). **Diperbarui (riwayat
  #42-#43):** seed Portal kini hanya akun `admin` (role `admin`), sesuai keputusan 52. Role Portal hanya `admin` dan `user`.

### 5.0b Tahapan Fase 3 (hubungkan Portal ke Toko), dikerjakan bertahap

Tahap 2.7 (toggle login langsung) digabung ke Fase 3 (disetujui user) karena penegakannya ada di BE Toko.
Perubahan di Toko bersifat **aditif**: login lama tetap berfungsi sampai toggle dimatikan.

| Tahap | Pekerjaan | Hasil yang dicek |
|-------|-----------|------------------|
| 3.1 **[SELESAI]** | BE Toko: migrasi (`users.portal_user_id`, role "menunggu akses", tabel `portal_credentials`, setting `direct_login_enabled`), config alamat Portal, `POST /api/auth/sso` (tukar kode ke Portal, cari/buat user lokal, buat sesi) | Kode asli dari Portal menghasilkan sesi Toko lewat API |
| 3.2 **[SELESAI]** | BE Toko: toggle login langsung (endpoint status publik, login ditolak saat off, endpoint ubah untuk admin, flag darurat di config) | Toggle bekerja, flag darurat mengalahkan setelan DB |
| 3.3 **[SELESAI]** | FE Toko: migrasi key browser ke `toko.*`, halaman `/sso`, tombol "Masuk lewat Portal" dan tombol kembali ke Portal, form login mengikuti toggle, UI toggle di Pengaturan, base path `/toko/` (lewat env) | Alur Portal ke Toko berjalan di browser (dev) |
| 3.4 **[SELESAI]** | Skrip sekali jalan: salin user Toko ke Portal dan tautkan `portal_user_id` (dry-run, backup) | Dry-run dan eksekusi pada DB dev |
| 3.5 **[SELESAI, disetujui user]** | Nginx (path `/`, `/portal-api/`, `/toko/`, `/toko/api/`), SOP deploy, panduan cutover VPS | Konfigurasi dan panduan siap, divalidasi `nginx -t` bila tersedia (belum: nginx tidak ada di mesin dev) |
| 3.6 **[SELESAI]** | Uji ujung ke ujung di dev (Portal 3000/8080 ke Toko 3001/8081) | 66 pemeriksaan API + browser; semua lolos kecuali bug pesan galat SSO di Toko (riwayat #46, diperbaiki di 3.6b) |
| 3.6b **[SELESAI, menunggu persetujuan]** | Perbaikan bug pesan galat SSO Toko + hapus 4 file sisa (riwayat #47) | Pesan "belum diatur" / "Portal tidak dapat dihubungi" sampai ke pengguna (503); regresi API 33/33 dan browser 21/21 lolos |

### 5.0c Tahapan Fase 4 (POS Coffee) [SEPAKAT, disetujui user; dikerjakan satu per satu]

Dasar: keputusan 22-45 dan bagian 14. Cara membangun sama seperti Portal: salin `TOKO/BE` dan `TOKO/FE` ke `COFFEE/`,
pangkas, ganti identitas, lalu tambah fitur kopi bertahap. Perubahan **hanya di `COFFEE/`** (Toko, Portal, dan BE User tidak
disentuh, kecuali data dev: path dan kunci aplikasi `coffee` di BE User). Identitas: module Go `coffee_api`, package FE
`coffee-web`, DB `pos_coffee_db`, **BE 8083** (keputusan 64), FE 3002, base `/coffee/`, API `/coffee/api`, key browser `coffee.*`,
kode aplikasi `coffee`, `SecretKey` sendiri. **Diperbarui setelah Fase 3B (keputusan 72):** Coffee disalin dari Toko versi
baru, jadi otomatis menjadi klien BE User seperti Toko (login form dan SSO lewat BE User, `auth_user_id`, `auth_credentials`,
tautan otomatis, Manajemen User hanya role/status/hapus, penyamaran log); skrip tautan user tidak diperlukan (keputusan 62
diperbarui, Coffee tidak punya akun lama). Pola tampilan mengikuti Toko (kolom aksi beku, varian `touch:`) karena dipakai di tablet.
Penamaan: domain `menu` milik Toko adalah menu navigasi/izin (tetap); menu jualan kopi memakai `menu_item` agar tidak bentrok.
Setiap tahap: verifikasi build/lint/type-check, uji API, uji browser (laptop + tablet sentuh), catat di riwayat, lalu berhenti.

| Tahap | Pekerjaan | Hasil yang dicek |
|-------|-----------|------------------|
| 4.1 **(diulang setelah Fase 3B) [SELESAI]** | Kerangka BE: salin ulang `TOKO/BE` versi baru (file ter-track git) menggantikan isi `COFFEE/BE` hasil 4.1 lama, pangkas modul (`supplier_return`, `expiry_batch`, `customer_balance`, `capital_reconciliation`, `backdate`, `version`, skrip `cmd/backfill_*` dan `cmd/migrate_users_to_auth`, `_sample_reference`), ganti identitas, **migrasi baseline baru** dari bahan 4.1 lama disesuaikan 3B (infrastruktur, `auth_user_id` + `auth_credentials`, `users.password` boleh kosong, role `admin`/`owner`/`kasir`/`menunggu_akses`, seed profil `owner`/`admin`, menu navigasi Coffee, kas/shift, pengeluaran, pelanggan, piutang, metode/status bayar tanpa "Saldo", `direct_login_enabled`, log), `AuthUrl` di config | BE hidup di **8083**, `/api/health`, tabel terbentuk di `pos_coffee_db` (dibuat ulang), login form `admin` lewat BE User menautkan profil seed, akun baru menjadi `menunggu_akses`, SSO dari Portal (path `coffee` dev `http://localhost:3002/`, kunci di `auth_credentials`) menghasilkan sesi, log tanpa password |
| 4.2 **[SELESAI]** | Kerangka FE: salin Toko, pangkas halaman modul yang dibuang, ganti identitas (port, base path, key `coffee.*`, judul/branding), menu sidebar Coffee | `type-check`/`lint`/`build`; di browser: kartu Coffee di Portal ke `/sso` ke dashboard Coffee, halaman Menunggu Akses, toggle login langsung, logout ke Portal |
| 4.3 **[SELESAI]** | Bahan baku (keputusan 80-84): katalog bahan (satuan dari master Unit, desimal, batas stok menipis, stok awal 0), stok sisa, total terpakai, total terbuang, input **pemakaian manual** (satu bahan per input), **barang terbuang** (alasan wajib), **pembatalan entri** oleh admin/owner, riwayat mutasi (siapa, kapan, catatan); stok minus diizinkan dengan peringatan; **tanpa opname**. Menu grup "Bahan Baku": Daftar Bahan, Riwayat Stok Bahan | CRUD bahan; pemakaian/terbuang mengubah stok dan total serta tercatat di riwayat; pembatalan mengembalikan stok dan total; izin kasir sesuai keputusan 80 |
| 4.4 **[SELESAI]** | Menu jualan: kategori menu, menu (mode stok **langsung** atau **racikan**), varian **ukuran × panas/dingin** dengan harga per kombinasi dan kombinasi bisa dinonaktifkan (menu tanpa varian: satu harga), **kelompok tambahan** (pilih satu/banyak, gratis/berbayar) dipasang ke menu, **resep sebagai catatan takaran per kombinasi** (opsional), tandai habis manual (keputusan 53-55) | CRUD menu/varian/tambahan/resep lewat UI; resep tidak mengubah stok |
| 4.5 **[SELESAI]** | Barang jadi dan pembelian (keputusan 91-94): katalog barang jadi = modul Produk salinan Toko (satuan bertingkat dus/pcs, logika `product_packages`), menu mode langsung ditautkan ke satu produk; Pembelian supplier dua jenis baris (bahan baku, barang jadi) dalam satu nota, bayar, edit, void; tanpa opname (keputusan 84) | Pembelian menambah stok bahan / barang jadi; void mengurangi; utang supplier seperti Toko |
| 4.6 **[SELESAI]** | Meja dan antrean (keputusan 56, 95-98): master meja (nama unik, urutan, aktif; dikelola admin/owner, kasir melihat), nomor antrean harian 001, 002, ... reset 00:00 Asia/Jakarta, diambil saat pesanan dibuat (dipakai pesanan di 4.7) | CRUD meja; nomor antrean urut dan reset harian |
| 4.7a **[SELESAI]** | Pesanan (BE): pesanan = tabel `transactions` Toko yang diperluas (jenis pesanan dine-in/take away, meja, nomor antrean, nama pemesan) + item menu (varian, tambahan, catatan); harga dihitung server dari varian + tambahan; diskon/pajak seperti Toko; barang jadi **dipotong saat item masuk pesanan** (bayar sekarang), **penjualan diblokir bila stok habis**; menu habis/nonaktif ditolak; utang seperti Toko; void dengan pilihan per item barang jadi kembali ke stok/dibuang (keputusan 99-102) | Uji API: harga, stok barang jadi, antrean, meja, void kembali/dibuang, kas |
| 4.7b **[SELESAI, menunggu persetujuan]** | Kasir (FE): halaman kasir baru (tombol menu per kategori, popup varian/tambahan/catatan, keranjang, jenis pesanan, meja, nama pemesan), pembayaran (modal Toko), struk (modul cetak Toko) dengan nomor antrean, riwayat transaksi dan void dengan pilihan kembali/dibuang | Transaksi bayar sekarang di laptop, tablet, dan HP; struk tercetak |
| 4.8 | Open bill: pesan dulu bayar belakangan, tambah item, **gabung bill**, **pindah meja**, batal (barang jadi: pilih kembali ke stok atau dibuang), jadikan **utang** dengan memilih pelanggan | Semua alur open bill; data stok/utang konsisten |
| 4.9 | Aturan kas dan batas hari: tutup kas manual **diblokir** selama ada pesanan terbuka; scheduler tengah malam **opsi i** (kas dengan pesanan terbuka tidak ditutup, ditandai); pesanan melewati batas hari hanya bisa diselesaikan; buka kas baru diblokir sampai kas tertunda ditutup; admin menyelesaikan atas nama kasir; daftar pesanan tertunda | Skenario lintas tengah malam diuji (jam disimulasikan di dev) |
| 4.10 | Laporan dan dashboard: penjualan (per menu, kategori, metode bayar), pemakaian bahan, barang terbuang, ringkasan kas | Angka laporan cocok dengan data uji |
| 4.11 | Uji e2e penuh Coffee di dev (laptop, tablet, HP dengan ketukan sentuh) dan perbaikan temuan | Semua alur lolos |

Usulan pembagian bila satu tahap terlalu besar: 4.7 dan 4.8 dapat dipecah BE/FE (seperti 2.6a/2.6b).
Sesuai keputusan 61, persiapan dan pelaksanaan deploy (nginx `/coffee/`, maintenance `coffee`, SOP, `.env.production`,
cutover VPS Portal/Toko/Coffee) **tidak masuk Fase 4**; dikerjakan belakangan sebagai fase deploy tersendiri.
Detail kecil yang baru muncul di suatu tahap (mis. aturan gabung bill, isi laporan) ditanyakan ke user di tahap itu sebelum kode.

### 5.0d Fase 3B: Layanan akun terpisah "BE User" [SEPAKAT, disetujui user; dikerjakan satu per satu]

Dasar: keputusan 63-72. Fase 4 (Coffee) dihentikan sementara di tengah 4.1 sampai fase ini selesai.

**Susunan**

```
FE Portal (3000) -> BE Portal (8081) -+
FE Toko   (3001) -> BE Toko   (8082) -+-> BE User (127.0.0.1:8080) -> pos_auth_db
FE Coffee (3002) -> BE Coffee (8083) -+    tidak ada di nginx, tanpa FE
```

**Letak data**

| Tempat | Isi |
|--------|-----|
| BE User (`pos_auth_db`) | `users` (username, nama, hash password Argon2id, aktif, wajib ganti password), `roles` (`admin`/`user`, pengelola akun), `apps` (kode, nama, path, ikon, hash kunci; termasuk aplikasi `portal`), `user_apps`, `sso_codes`, log permintaan |
| Portal (`pos_portal_db`) | Profil lokal `users` (`auth_user_id`), `sessions`, log. Tidak lagi menyimpan password, aplikasi, akses, kode SSO |
| Toko / Coffee | Profil lokal `users` (`auth_user_id`, role lokal, PIN), `sessions`, `auth_credentials` (pengganti `portal_credentials`: `app_id` + kunci untuk BE User), data bisnis. Kolom password lokal tidak dipakai lagi (keputusan 66) |

**API BE User (internal, setiap permintaan wajib `X-App-Key`; aplikasi pemanggil dikenali dari kuncinya)** — bentuk akhir yang dibangun di 3B.2

Endpoint Portal membawa header `X-Acting-User-Id` (user yang bertindak; Portal mengisinya setelah memeriksa sesi Portal-nya sendiri).

| Kelompok | Endpoint | Dipanggil oleh |
|----------|----------|----------------|
| Verifikasi login | `POST /api/login/verify` `{username, password, client_ip}` -> `{id, username, full_name, role_id, role_name, is_active, must_change_password}`. Salah/tidak ada: 401 pesan seragam; nonaktif: 401 "Akun tidak aktif"; aplikasi selain Portal: wajib ganti password 403 "Ganti password dulu di Portal", tanpa akses 403; pembatas 5 gagal/15 menit per `client_ip`+username (429) | Portal, Toko, Coffee |
| SSO | `POST /api/sso/launch` `{app_code}` -> `{redirect_url, expires_in_seconds}` (Portal bukan tujuan); `POST /api/sso/exchange` `{code}` -> `{valid, app_id, user{id, username, full_name, is_active}}` (kode harus milik aplikasi pemanggil, sekali pakai, 60 detik) | launch: Portal; exchange: Toko, Coffee |
| Akun sendiri | `POST /api/account/profile`, `POST /api/account/change-password` `{old_password, new_password}`, `POST /api/apps/my` | Portal (user mana pun) |
| Kelola akun dan aplikasi | `/api/users/{list, detail/:id, create, update/:id, reset-password/:id, delete/:id, toggle-status/:id}`, `/api/roles/options`, `/api/apps/{list, create, update/:id, toggle-status/:id, rotate-key/:id, access/:user_id, access/:user_id/set}`; hanya role `admin` yang tidak sedang wajib ganti password. Aplikasi `portal` tidak bisa dinonaktifkan atau diputar kuncinya lewat API (pakai `cmd/set_app_key`) | Portal (admin) |

Kunci Toko/Coffee hanya untuk verifikasi login dan penukaran SSO (endpoint lain 403 "hanya untuk Portal").
Kunci awal dibuat dengan `go run ./cmd/set_app_key --app <kode> --confirm` di `AUTH/BE` (tampil sekali, tersimpan hash).

**Alur**
- Login form di Toko/Coffee (toggle nyala): FE -> BE aplikasi -> `login/verify` -> BE aplikasi mencari profil lewat `auth_user_id`; bila belum ada: tautkan otomatis ke profil lokal username sama yang belum tertaut (keputusan 67), atau buat profil baru `menunggu_akses` -> sesi aplikasi.
- Login Portal: FE Portal -> BE Portal -> `login/verify` -> profil lokal Portal + sesi Portal (role dari BE User).
- SSO: kartu di Portal -> BE Portal -> `sso/launch` -> browser ke `<app>/sso?code=` -> BE aplikasi -> `sso/exchange` -> sesi aplikasi.
- Kelola user/aplikasi dan ganti password: halaman Portal -> BE Portal -> BE User.
- Flag server `DirectLoginForceEnabled` **dibuang** (keputusan 73). Toggle disimpan per aplikasi di `settings.direct_login_enabled`; pemulihan bila terkunci lewat SQL.

**Tahapan**

| Tahap | Pekerjaan | Hasil yang dicek |
|-------|-----------|------------------|
| 3B.1 **[SELESAI]** | Kerangka BE User: struktur sama dengan aplikasi lain (salin dari BE Portal: domain auth/user/role/app/sso), module `auth_api`, port 8080 hanya `127.0.0.1`, DB `pos_auth_db`, migrasi + seed (role, akun `admin`, aplikasi `portal`/`toko`/`coffee`) | BE hidup, hanya bisa diakses dari localhost, tabel terbentuk |
| 3B.2 **[SELESAI]** | API BE User lengkap (tabel API di atas), pembatas login, pengecekan kunci per aplikasi dan `acting_user_id` | Uji API: verifikasi benar/salah/nonaktif/tanpa akses/wajib ganti, pembatas login, SSO sekali pakai dan kedaluwarsa, kunci salah, kelola akun hanya admin |
| 3B.3 **[SELESAI]** | BE Portal menjadi klien BE User: port 8081, login/ganti password/SSO launch/kelola user-aplikasi-akses diteruskan ke BE User, tabel akun Portal dipangkas jadi profil lokal; FE Portal hanya ganti alamat dev | Uji API + browser semua fitur Portal (skrip `v_portal.js`) |
| 3B.4 **[SELESAI]** | BE dan FE Toko menjadi klien BE User: port 8082, login form lewat BE User (toggle tetap), SSO exchange ke BE User, tautan otomatis, `portal_user_id` -> `auth_user_id`, `portal_credentials` -> `auth_credentials` (migrasi 012), Manajemen User tanpa buat user/ganti password (keputusan 69), password lokal tidak dipakai | Uji API + browser Toko (skrip `v_toko.js`, uji sentuh tablet/HP) |
| 3B.5 **[SELESAI]** | Skrip migrasi user Toko disesuaikan (`cmd/migrate_users_to_auth`): salin user Toko ke BE User (hash apa adanya, kasir tetap login dengan password lama) dan isi `auth_user_id`; ~~data dev Portal dipindah ke `pos_auth_db`~~ tidak dipindah (keputusan 78); kolom/tabel lama dibuang (keputusan 79) | Dry-run dan `--apply` pada salinan DB / DB dev |
| 3B.6 **[SELESAI]** | Konfigurasi menyeluruh: port di `.env`/config/proxy dev ketiga aplikasi, nginx (`/portal-api/` -> 8081, `/toko/api/` -> 8082; BE User tidak dipasang), README, HANDOFF; uji ujung ke ujung Portal + Toko (laptop, tablet, HP) | Semua skrip uji lolos |

Setelah 3B selesai: kembali ke 4.1 Coffee (port 8083, salin dari Toko versi baru). Perubahan production tetap belum ada (Portal, Toko baru, BE User belum pernah di-deploy); panduan cutover diperbarui di fase deploy (keputusan 61).

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
     Body  : { "app_id": "toko", "code": "XXXX" }
7. BE Portal menjawab:
     { "valid": true, "app_id": "toko",
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
2. Saat penukaran kode: BE aplikasi mengirim `app_id` dan kunci rahasia dari tabel `portal_credentials` di
   database aplikasi (bukan `.env`, lihat 6.10). Portal mencocokkan keduanya dengan kode tadi.
3. Jawaban Portal memuat `app_id`; BE aplikasi membandingkannya dengan `app_id` di `portal_credentials`. Tidak
   sama = ditolak.

Catatan:
- `app_id` adalah **kode teks** aplikasi di Portal (kolom `apps.code`, mis. `"toko"`, `"coffee"`), bukan id angka.
- Kode aplikasi bukan rahasia, jadi **kunci rahasia (`X-App-Key`) tetap wajib**.
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
- Alamat Portal di config (`PortalUrl`); `app_id` dan kunci rahasia aplikasi di tabel `portal_credentials` database
  aplikasi (diperbarui sesuai 6.10; semula direncanakan di `.env`).
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
- ~~Coffee tidak punya akun lama, jadi tidak ada migrasi; profil dibuat otomatis saat SSO pertama.~~ Diperbarui
  (keputusan 62): Coffee di-seed akun `owner` dan `admin` seperti Toko, lalu ditautkan ke Portal dengan skrip
  `migrate_users_to_portal` yang sama; staf baru tetap dibuat di Portal dan mendapat profil "menunggu akses" saat SSO pertama.

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

Berdasarkan daftar domain BE Toko (`TOKO/BE/domain/`):

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
4. **Daftar yang wajib diganti setelah salin:** module Go, nama package FE, `app_id` (baris `portal_credentials`), port, nama database
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
| 41 | Tahap 3.3 selesai (FE Toko). Base path lewat env: `VITE_BASE_PATH` (default `/`, perilaku lama) dibaca `vite.config.ts` (`loadEnv`), React Router memakai `basename`, helper `appUrl()` dipakai untuk navigasi penuh (login, reload, refresh token gagal); `VITE_PORTAL_URL` (kosong = fitur Portal tersembunyi). `.env.production` SENGAJA belum diubah (tetap root dan `/api`) supaya deploy Fase 1 tidak berubah sebelum cutover; nilai produksi (`VITE_BASE_PATH=/toko/`, `VITE_API_URL=/toko/api`, `VITE_PORTAL_URL=/`) diterapkan saat tahap 3.5. Migrasi key browser: `bootstrap/migrateStorageKeys.ts` di-import paling awal di `main.tsx` (sebelum store dibuat) menyalin `auth-session`, `cashier-draft`, `backdate-cashier-draft`, `pos.native_bt_printer_address` ke `toko.*` lalu menghapus yang lama; sesi, keranjang, dan printer terjaga. Halaman baru: `/sso` (menukar kode; galat disimpan di state lokal karena `useMutation` kehilangan hasil di StrictMode), `/menunggu-akses` (akun tanpa izin menu; `getDefaultRoute` kini mengarah ke sini, bukan `/login`, agar tidak loop). Halaman login mengikuti `login-options` (form + tombol "Masuk lewat Portal", atau hanya tombol Portal saat login langsung mati; bila opsi gagal dimuat, form tetap tampil). Header Toko memiliki tombol "Kembali ke Portal"; logout dari Toko mengembalikan ke Portal. Kartu toggle login langsung di Manajemen User (hanya admin, konfirmasi saat mematikan). Terverifikasi: `type-check`, `lint`, `build`, dan uji browser dev (Portal 3000 ke Toko 3001): SSO pertama ke halaman menunggu akses, setelah diberi role kasir masuk ke `/dashboard`, logout kembali ke Portal, toggle lewat UI mengubah DB dan tampilan login, kode ngawur menampilkan galat, migrasi key menjaga isi keranjang. Catatan: di DB dev, owner asli adalah `Banana` dan `BramaGanteng` (tidak ada akun `owner`). |
| 42 | Keputusan user soal pemetaan role saat migrasi user Toko ke Portal: role `admin` di Toko menjadi `admin` di Portal; semua role lain (termasuk `owner` dan `kasir`) menjadi `user` di Portal. Konsekuensi yang dicatat AI: akun seed `owner` di Portal (role admin) tidak sesuai aturan ini sehingga seed Portal dipangkas hanya `admin`. Tahap 3.4 dimulai. |
| 43 | Tahap 3.4 selesai: skrip `TOKO/BE/cmd/migrate_users_to_portal` (dijalankan dari `TOKO/BE`). Default dry-run (hanya laporan); `--apply` wajib disertai `--i-have-backup` (backup `pos_retail_db` dan `pos_portal_db` dibuat manual dulu, perintah mysqldump ada di komentar file). Koneksi Toko seperti skrip cmd lain; Portal lewat `--portal-db` (kredensial Toko) atau `--portal-dsn`/`PORTAL_MIGRATION_DSN`. Aturan: hash Argon2id disalin apa adanya; role Toko `admin` menjadi `admin` Portal, semua role lain `user` (keputusan 52); akun Portal baru tidak wajib ganti password; akses aplikasi `toko` diberikan; `users.portal_user_id` diisi. Username sudah ada di Portal: password bawaan seed (admin123/owner123) diambil alih (password, nama, role, status disamakan, sesi lama dihapus); password sama langsung ditautkan; password lain = KONFLIK, dilewati dan dilaporkan (tidak menebak). Portal ditulis lebih dulu (satu transaksi), lalu Toko (satu transaksi); aman diulang. Peringatan dicetak untuk akun Portal berpassword bawaan yang tersisa. Seed Portal (`002_seed_data.sql`) dipangkas menjadi hanya akun `admin` (akun `owner` dihapus dari seed agar sesuai aturan; hanya memengaruhi database baru). Terverifikasi pada SALINAN database: BUAT, AMBIL ALIH, KONFLIK, penolakan `--apply` tanpa konfirmasi backup, hash identik, role sesuai aturan, akses diberikan, dan `--apply` kedua tidak mengubah apa pun; salinan dihapus. Catatan: DB dev Portal masih memuat akun `owner` (role admin, password sudah bukan bawaan) dari seed lama; tidak terkait user Toko mana pun dan tidak disentuh. |
| 44 | Tahap 3.5 selesai: `nginx/pos-web.conf` baru (Portal di `/`, API Portal `/portal-api/`, Toko di `/toko/`, API Toko `/toko/api/`, `/api/` lama tetap ke Toko 8081, cache aset 1 tahun), maintenance global dan per aplikasi (`nginx/maintenance/maintenance-{on,off}.sh` menerima `all|portal|toko|coffee`), `TOKO/FE/.env.production` kini `/toko/` (`VITE_API_URL=/toko/api`, `VITE_BASE_PATH=/toko/`, `VITE_PORTAL_URL=/`), `docs/PANDUAN_CUTOVER_PORTAL.md`, bagian "SOP Deploy per Aplikasi" di `SOP_DEPLOY.md`, README dan catatan pemetaan di dokumen deploy lama. Terverifikasi: build produksi Toko (aset `/toko/assets/`) dan Portal (aset `/assets/`) lolos. BELUM terverifikasi: `nginx -t` (nginx tidak terpasang di mesin dev) dan panduan cutover belum dijalankan di VPS. Peringatan: build Toko dari commit ini hanya cocok dengan nginx baru; jangan dipakai untuk deploy gaya Fase 1. |
| 45 | Serah terima ke AI baru (Claude); user menyetujui tahap 3.5. Dokumen dirapikan agar tidak rancu dengan kode (dicek langsung ke kode): kepala dokumen (implementasi Fase 1-3 sudah ada); status tahap 2.1-2.6 dan 3.1-3.5 diseragamkan menjadi SELESAI (sudah di-commit user), 2.7 ditandai digabung ke Fase 3; contoh `app_id` di 6.2 diganti kode teks `"toko"` (sesuai `apps.code` dan `portal_credentials.app_id` VARCHAR); 6.3, 6.5, dan 14.5 tidak lagi menyebut `APP_ID` di `.env`; catatan seed Portal di 5.0 disesuaikan (hanya `admin`); path 14.1 menjadi `TOKO/BE/domain/`. Ditambahkan bagian 3.1 "Standar konsistensi kode antaraplikasi" atas permintaan user (stack dan versi, struktur folder BE/FE, pola domain/fitur, daftar modul salinan), hasil cek kode Toko vs Portal: stack, versi, dan struktur sudah sama. Ditemukan file sisa yang membingungkan: `eslint.config.ts` di Toko dan Portal (tidak terbaca ESLint karena `eslint.config.js` didahulukan; berisi aturan lebih ketat yang tidak pernah aktif), `TOKO/FE/src/App.tsx` (re-export sisa, tidak dipakai), `TOKO/FE/nginx example.conf` (pola nginx lama). User setuju dihapus. |
| 46 | Tahap 3.6 selesai (uji ujung ke ujung di dev, keempat server dijalankan; backup DB dev dibuat lebih dulu). Persiapan data dev: alamat aplikasi `toko` di Portal diubah ke `http://localhost:3001/` lewat API admin, kunci dibuat lewat rotate-key lalu diisi ke `portal_credentials` Toko lewat SQL, skrip `migrate_users_to_portal` dijalankan (dry-run, penolakan `--apply` tanpa `--i-have-backup`, `--apply` pada DB dev: Banana/BramaGanteng BUAT sebagai `user`, `admin` AMBIL ALIH; hash identik; `--apply` kedua tidak mengubah apa pun). **Lolos (API):** login-options sebelum/sesudah kunci terpasang; mematikan login langsung saat Portal belum siap ditolak 400; username sama yang belum tertaut ditolak (tidak auto-link); admin tertaut masuk dengan role lokal; kode sekali pakai, kode ngawur, kode kedaluwarsa >60 detik ditolak; user baru buatan admin wajib ganti password (launch ditolak 403 selama itu), tanpa akses kartu kosong dan launch ditolak, setelah diberi akses SSO pertama membuat profil `menunggu_akses` dengan `portal_user_id`, role lokal kasir dipakai pada SSO berikutnya, nama dari Portal diperbarui di Toko; akses dicabut / user nonaktif / aplikasi nonaktif: kode yang sudah dibuat pun ditolak; kunci salah dan `app_id` salah ditolak seragam 401; owner tidak boleh mengubah toggle (403), admin mematikan -> login langsung 403 dan SSO tetap jalan, dinyalakan lagi; logout Toko mengakhiri sesi Toko tetapi sesi Portal tetap; satu user satu sesi (SSO kedua mematikan token pertama); sesi Toko yang sudah berjalan tetap hidup walau user dinonaktifkan di Portal (sesuai 6.8). **Lolos (browser Chrome):** login Portal -> kartu -> `/sso` -> dashboard Toko di tab yang sama; key `portal.*` dan `toko.*` tanpa key lama; tombol "Kembali ke Portal" dan buka Toko lagi tanpa login ulang; logout Toko kembali ke Portal; kode ngawur menampilkan galat + tombol kembali; user baru dipaksa ganti password lewat UI lalu sampai ke halaman Menunggu Akses; kartu toggle hanya untuk admin (owner tidak melihat); halaman login mengikuti toggle; tanpa scroll horizontal di 390px; tanpa error JavaScript. **BUG (belum diperbaiki, menunggu keputusan user):** di BE Toko, pesan SSO "Login lewat Portal belum diatur di server ini" dan "Portal tidak dapat dihubungi" dibungkus `InternalServerError`, sedangkan middleware error sengaja mengganti pesan semua galat 500 menjadi "Internal Server Error"; pengguna di `/sso` hanya melihat "Internal Server Error" (bertentangan dengan 6.10 dan bagian 11). Kasus "Portal tidak dapat dihubungi" dipastikan dari kode (jalur sama), tidak diuji dengan mematikan Portal. Tidak diuji ulang di tahap ini (sudah di 3.2): flag darurat `DirectLoginForceEnabled` dan proteksi key toggle di `POST /settings`. Data uji dibiarkan di DB dev: user Portal+Toko `kasir72589`, `web12286`, user Toko `owner72589`. |
| 47 | Tahap 3.6b selesai atas persetujuan user. **Perbaikan bug:** tipe galat baru `ServiceUnavailableError` di package `errors` (HTTP 503, kode respons `"53"` = `helper.StatusServiceUnavailable` yang sudah ada); berbeda dari `InternalServerError`, pesannya diteruskan ke pengguna (tetap dicatat di log). Ditambahkan sama persis di Toko dan Portal (package `errors`, tesnya, dan `middleware/errorHandler_middleware.go` tetap identik antaraplikasi, sesuai bagian 3.1); di Portal belum dipakai. Toko `auth_service_portal.go` memakai tipe ini untuk "Login lewat Portal belum diatur di server ini" dan "Portal tidak dapat dihubungi. Coba lagi beberapa saat". Galat internal lain (mis. gagal membuat profil lokal) tetap `InternalServerError`. Catatan nginx: `pos-web.conf` tidak memakai `proxy_intercept_errors`, jadi 503 dari BE diteruskan apa adanya dan tidak tertimpa halaman maintenance (berdasarkan perilaku nginx, belum diuji karena nginx tidak ada di mesin dev). **Hapus file sisa** (disetujui user): `PORTAL/FE/eslint.config.ts`, `TOKO/FE/eslint.config.ts`, `TOKO/FE/src/App.tsx`, `TOKO/FE/nginx example.conf`. Terverifikasi: `go build`, `go vet`, `go test ./errors/` (Toko dan Portal); `type-check`, `lint`, `build` FE Toko dan Portal; kredensial kosong -> 503 + pesan "belum diatur" (API dan tampilan `/sso` di browser); Portal BE dimatikan -> 503 + "Portal tidak dapat dihubungi" (diuji langsung); regresi penuh API 33/33 dan browser 21/21 lolos. Data uji tambahan di DB dev: `kasir63112`, `owner63112`, `web16889`. |
| 48 | Uji visual e2e atas permintaan user (Chrome tampil, jeda 2 detik per aksi; skrip `testing/portal_toko_e2e/v_*.js`; backup DB dev dibuat lebih dulu). **Portal 40/40 lolos:** login (validasi, salah, benar), menu akun, Kelola User (daftar/paginasi/cari, akun sendiri tanpa aksi berbahaya, tambah + Acak, ubah, akses aplikasi, reset password, nonaktif/aktif, hapus), Kelola Aplikasi (tambah dengan kunci tampil sekali + Salin, ubah dengan kode terkunci, buat kunci baru, nonaktifkan), wajib ganti password, ganti password sendiri, user biasa tidak bisa membuka menu admin, user tanpa akses, pembatas 5 kali gagal, logout, tampilan HP 390px, tanpa error JavaScript. **Toko semua skenario lolos** (setelah 3 kali pengulangan sebagian karena kesalahan skrip uji, bukan aplikasi: dropdown yang sudah berisi nilai default, halaman yang belum dibuka, dan selektor ikon `ban` yang ikut mengenai `banknote`): masuk lewat Portal; tambah/ubah/hapus kategori, pelanggan, supplier, pengeluaran; tambah produk dan ubah harga; buka kas; pembelian supplier (nota 10 pcs menambah stok, bayar lunas, nota kedua lalu void mengurangi stok kembali); kasir tunai dan hutang (stok berkurang, piutang tercatat); riwayat transaksi, piutang, laporan penjualan dan stok; toggle login langsung; user baru ke Menunggu Akses; Kembali ke Portal; logout; tanpa error JavaScript. Stok dan status dicek ke DB. Tidak ditemukan bug aplikasi. Data uji dibiarkan di DB dev (akhiran `15936` Portal, `81382` Toko); kas hari uji (11 Nov 2026) dibiarkan terbuka. |
| 49 | Uji tampilan Portal di HP 390px, tablet tegak 768px, tablet mendatar 1024px, laptop 1366px (15 layar/dialog per ukuran, skrip `testing/portal_toko_e2e/responsive/r_portal.js`). Temuan: di tablet tegak header terlalu padat (teks dan menu terpotong dua baris) dan tombol aksi Kelola Aplikasi tersembunyi karena tabel harus digeser; di semua ukuran beberapa target sentuh terlalu kecil (tombol tutup dialog 16px, ikon mata 16px, item menu dropdown dan paginasi 32px); di HP kartu user menyisakan baris kosong untuk tombol aksi. Atas persetujuan user, diperbaiki di FE Portal saja: menu navigasi header mulai `lg` (1024px), di bawahnya masuk menu akun; `DataTable` mendapat prop `cardBreakpoint` (default `md`, Kelola Aplikasi `lg`); tombol aksi card di kanan atas; tombol tutup dialog, ikon mata (login dan `PasswordInput`, kini ber-`aria-label`), item dropdown, dan paginasi menjadi 40px; judul dialog diberi ruang agar tidak tertimpa tombol tutup. Terverifikasi: `type-check`, `lint`, `build`; uji tampilan ulang 60 layar tanpa masalah (termasuk pemeriksaan baru "isi harus digeser"); regresi fungsional Portal 39/40, satu-satunya gagal adalah selektor skrip uji (dua baris "(diubah)" dari data uji sebelumnya; nama aplikasi uji terbukti berubah di DB), skrip sudah diperbaiki. Belum diuji: HP posisi mendatar, perangkat asli/APK. Data uji baru: `e2eu84716`, aplikasi `uji-e2e-84716` (nonaktif). |
| 50 | Uji tampilan Toko (22 layar/dialog × HP 390, tablet 768 dan 1024, laptop 1366; skrip `testing/portal_toko_e2e/responsive/r_toko.js`, pengukur bersama `r_common.js`). Temuan: HP baik (laci menu, card, tab kasir); di tablet 11 tabel lebih lebar dari layar sehingga kolom kanan termasuk tombol aksi terpotong tanpa tanda bisa digeser; di laptop 1366 Produk dan Kas Harian juga terpotong (tombol cetak/kunci/hapus Produk tersembunyi); banyak target sentuh 24-32px (ikon aksi tabel, filter status, ikon header, paginasi, tombol tutup dialog 16px, tombol kasir). Atas persetujuan user (ikut rekomendasi AI), diperbaiki di FE Toko, murni tampilan: (1) `DataTable` membekukan kolom `actions` di kanan (sticky, latar hover/terpilih ikut, bayangan hanya bila ada kolom tertutup) sehingga aksi selalu terlihat di tablet dan laptop; kolom aksi `CashDrawerTableColumns` dan `ReturnTableColumns` diberi key `actions` (sebelumnya `id`). (2) Varian Tailwind baru `touch:` (`@media (pointer: coarse)`) di `tailwind.config.ts`: di perangkat sentuh saja, tombol aksi tabel/card, paginasi, filter status (`toggle-group`), ikon header Navbar (kini ber-`aria-label`), tombol tutup dialog (+ ruang judul), ikon mata (login dan `PasswordInput`), serta seluruh tombol/input kasir dan dialog Bayar (kecuali kotak centang) minimal 40px; laptop dengan mouse tidak berubah. Terverifikasi: `type-check`, `lint`, `build` (CSS memuat `@media (pointer: coarse)`); uji tampilan ulang 86 layar: semua tabel utama "digeser, kolom aksi tetap terlihat", target sentuh tersisa hanya kotak centang 16px, kunci barcode 14px di form Produk/Pembelian, kontrol form Pembelian 28-32px, dan saklar Login langsung; sisa "isi harus digeser": tabel item dialog Pembelian (sudah ada petunjuk geser) dan Laporan Penjualan di tablet (±100px, tanpa kolom aksi). Catatan uji: emulasi sentuh Chrome kehilangan `pointer: coarse` setelah screenshot satu halaman penuh, jadi untuk HP/tablet screenshot dibatasi layar terlihat dan interaksi memakai `tap`. Regresi fungsional Toko (mode cepat) 30/32: dua gagal adalah pemeriksaan skrip yang terlalu cepat setelah pencarian (data terbukti tersimpan di DB, skrip diperbaiki). Kas 11 Nov ditutup otomatis tengah malam (perilaku benar); regresi membuka kas 12 Nov. Data uji baru akhiran `30143`. Belum diuji: perangkat asli/APK, HP mendatar. |
| 51 | Uji fungsi dengan ketukan sentuh (bukan klik mouse) setelah perubahan tampilan #49-#50, jendela Chrome tampil ukuran tablet 768x1024 dan HP 390x844 dengan mode sentuh, jeda 0,5 detik, jendela ditutup tiap bagian selesai (`testing/portal_toko_e2e/run_sentuh.sh`). **Portal tablet 40/40, Portal HP 40/40** (menu admin lewat menu akun, kartu Kelola Aplikasi, tombol "..." di kartu, semua dialog). **Toko tablet 32/32** (aksi diketuk dari kolom beku, kasir lewat tab Produk/Keranjang, tombol tutup dialog, paginasi, filter). **Toko HP 26/31 lalu 10/10 setelah perbaikan skrip uji**: skrip mencari nota pembelian lewat No. Faktur, padahal kolom itu sengaja disembunyikan di card HP (`mobileHidden`); nota terbukti tersimpan dan tampil (screenshot dan DB), T10-T11 gagal berantai karena nota kedua belum di-void. Skrip kini mencari nota lewat supplier + total, pengecekan stok kasir dibuat relatif; T8-T12 diulang di HP: bayar lunas, void, kasir tunai dan hutang, riwayat, piutang, laporan lolos. Tidak ditemukan bug aplikasi. Data uji baru: Portal akhiran `96313`, `19348`; Toko akhiran `41235`, `77133`. |
| 52 | User memutuskan lanjut ke Fase 4. Disusun usulan tahapan 4.1-4.11 di bagian 5.0c (belum disetujui) beserta daftar pertanyaan terbuka; belum ada kode Coffee. |
| 53 | Jawaban user atas pertanyaan Fase 4: keputusan 53-61 ditambahkan (varian ukuran × panas/dingin, kelompok tambahan, resep per kombinasi, antrean harian dan meja, diskon/pajak seperti Toko tanpa service charge, role admin/owner/kasir, open bill untuk dine-in dan take away, cara salin 14.5 poin 1/3/4/6 disetujui, deploy VPS seluruh aplikasi belakangan). Tahapan 5.0c disesuaikan: 4.4 memuat varian/tambahan/resep sesuai keputusan, 4.11 menjadi uji e2e penuh di dev (persiapan deploy dikeluarkan dari Fase 4). Bagian 14.13 [USULAN] format antrean dan master meja kini SEPAKAT (keputusan 56). Tahapan 5.0c menunggu persetujuan user sebelum kode. |
| 54 | User menyetujui tahapan 5.0c dan meminta dikerjakan satu per satu. Keputusan 62 (akun awal Coffee sama seperti Toko: seed `owner`/`admin`, tautan lewat skrip `migrate_users_to_portal`); baris 4.1 dan bagian 7 diperbarui. Tahap 4.1 dimulai. |
| 55 | Tahap 4.1 berjalan sebagian (belum selesai): `COFFEE/BE` disalin dari file Toko yang ter-track git, modul dipangkas, module `coffee_api`, ketergantungan ke modul dipangkas diputus (saldo pelanggan di transaksi, batch kedaluwarsa dan retur di pembelian/supplier), baseline `001_init_schema.sql` (30 tabel diurutkan sesuai foreign key; kolom sisa `customers.balance`, `transactions.balance_used`, `cash_drawer.is_backdate` dipertahankan sementara) dan `002_seed_data.sql`, konfigurasi (`APP_PORT` 8082, DB `pos_coffee_db`, `SecretKey` dev baru, config prod tanpa rahasia Toko), skrip tautan diarahkan ke `coffee`. Terverifikasi: `go build`/`go vet`, BE hidup, migrasi bersih, uji API 9/9 (uji asap 59 endpoint baca tanpa 500, login langsung admin/owner, aplikasi `coffee` terdaftar di Portal dev, kunci terpasang). Ditemukan: skrip tautan membandingkan hash secara persis, sehingga `admin` Coffee selalu KONFLIK dengan `admin` Portal (hash Portal = hash Toko hasil ambil alih 3.6); penjelasan AI sebelumnya bahwa KONFLIK selesai dengan "menyamakan password" keliru (hash ber-salt). Diskusi berlanjut: user memilih login satu pintu lewat layanan akun terpisah. Keputusan 63-72 dicatat; rancangan Fase 3B di bagian 5.0d (menunggu persetujuan). Belum: `--apply` skrip tautan, uji SSO Coffee (tidak dilanjutkan karena digantikan 3B). Data dev: DB `pos_coffee_db` dibuat, aplikasi `coffee` terdaftar di Portal dev. |
| 56 | User menyetujui rancangan 5.0d, flag `DirectLoginForceEnabled` dibuang (keputusan 73; toggle tetap di UI, pemulihan lewat SQL). **Tahap 3B.1 selesai:** `AUTH/BE` disalin dari file BE Portal yang ter-track git (module `auth_api`); **hanya mendengarkan `127.0.0.1:8080`** (di `main.go`); middleware CORS dilepas (tidak dipanggil browser); rute hanya `GET /api/health` (API internal di 3B.2; domain auth/user/role/app/sso dari Portal tetap ada sebagai bahan, belum dirouting); baseline `pos_auth_db` (`roles`, `users` + `must_change_password` + FK role, `apps` + kolom baru `show_in_portal`, `user_apps`, `sso_codes`, log; tanpa tabel sesi) dan seed (role `admin`/`user`, akun `admin`/`admin123`, aplikasi `portal` (tidak tampil sebagai kartu), `toko`, `coffee`, tanpa kunci; akses admin ke ketiganya); konfigurasi port 8080, DB `pos_auth_db`, `SecretKey` dev baru, config prod tanpa rahasia, CORS kosong. **Keamanan baru:** paket `helper/logmask` menyamarkan `password`, `old_password`, `new_password`, `app_key`, `token`, `refresh_token`, `redirect_url`, dan `code` panjang (kode SSO; kode status respons tetap) sebelum body permintaan/respons disimpan di `log_requests` (3 uji unit). Terverifikasi: `go build`/`go vet`, uji unit `logmask` dan `errors`, BE hidup, port 8080 hanya `127.0.0.1` (dari IP LAN tidak terjangkau, dari localhost 200), migrasi dan seed sesuai, log uji menyimpan `***` untuk rahasia dan tidak ada nilai asli. **Temuan (belum ditangani, menunggu keputusan user):** middleware log yang sama di Toko, Portal, dan Coffee menyimpan password **dalam bentuk asli** di `log_requests` (DB dev: Toko 313 baris, Portal 253, Coffee 2; dari login, buat user, reset password); Toko production hampir pasti juga (kode sama), belum bisa diperiksa. |
| 57 | Keputusan user atas temuan log password: log lama dibiarkan, fokus ke depan; penyamaran dipasang di Portal (3B.3), Toko (3B.4), Coffee (4.1 diulang) — keputusan 74. |
| 58 | User menyetujui 3B.1 (mengikuti rekomendasi AI). **Tahap 3B.2 selesai (API BE User):** domain `auth` ditulis ulang tanpa sesi/JWT (verifikasi login, profil, ganti password sendiri); middleware baru `AppKeyMiddleware` (hash SHA-256 `X-App-Key` -> aplikasi aktif pemanggil), `PortalOnlyMiddleware`, `ActingUserMiddleware` (`X-Acting-User-Id`; mode admin menolak non-admin dan admin yang wajib ganti password); pembatas login memakai `client_ip` dari aplikasi; SSO exchange tidak lagi butuh `app_id` di body (kode harus milik aplikasi pemanggil), respons `user.id` (bukan `portal_user_id`); `apps/my` dan SSO launch mengecualikan aplikasi `show_in_portal = 0`; aplikasi `portal` tidak bisa dinonaktifkan/diputar kuncinya lewat API; domain `user` tidak lagi menghapus sesi (tanggung jawab Portal/aplikasi); `SecretKey` tidak wajib (tanpa JWT), `middleware/auth`, `RoleMiddleware`, `pkg/jwt` dibuang; perintah server `cmd/set_app_key` (lihat dulu, tulis dengan `--confirm`, kunci tampil sekali) beserta `cmd/internal/migrationdb` (salinan pola Toko). Terverifikasi: `go build`/`go vet`, uji unit, **uji API 33/33** (`testing/auth_3b2/api_3b2.js`): kunci tidak ada/salah, verifikasi benar/salah/seragam, wajib ganti (Toko ditolak, Portal diterima), ganti password, akses per aplikasi, nonaktif, pembatas per IP pengguna (IP lain tetap bisa), endpoint khusus Portal dan admin, profil, reset password (termasuk akun sendiri ditolak), kartu tanpa `portal`, pengaman aplikasi Portal, rotasi kunci coffee, SSO (kode lintas aplikasi ditolak, sekali pakai, Portal bukan tujuan, launch hanya Portal, kedaluwarsa 60 detik), hapus user, log tanpa password/kunci asli. Diulang setelah perbaikan kecil terakhir: 33/33. Data dev: kunci uji portal/toko/coffee di `testing/auth_3b2/keys.txt` (gitignore), user uji `kasir72090` di `pos_auth_db`. Belum: Portal dan Toko belum memakai BE User (3B.3, 3B.4). |
| 59 | User menyetujui 3B.2. **Tahap 3B.3 selesai (BE Portal menjadi klien BE User):** BE Portal pindah ke port **8081** (`.env*`), FE Portal dev menunjuk 8081 (`.env.development`, proxy Vite); klien `pkg/authclient` (header `X-App-Key` + `X-Acting-User-Id`, batas waktu 10 detik; penolakan 4xx BE User diteruskan apa adanya lewat `errors.UpstreamError`, termasuk 201 dan rincian validasi 422; BE User mati/5xx -> 503 "Layanan akun tidak dapat dihubungi"; kunci belum dipasang/ditolak -> 503 "Layanan akun belum diatur di server ini", bukan 401 ke pengguna); kunci Portal disimpan di tabel baru **`auth_credentials`** (sama dengan rencana Toko/Coffee), alamat BE User di config `AuthUrl`. Login: verifikasi ke BE User (dengan `client_ip`), lalu profil lokal (`users.auth_user_id`; profil lama dengan username sama yang belum tertaut otomatis ditautkan, selain itu dibuat baru) dan sesi Portal; `user.id` di respons = id BE User (cocok dengan halaman Kelola User). Status wajib ganti password disalin ke `sessions.must_change_password` (dihapus setelah ganti password berhasil). `/auth/me` dan refresh token memeriksa akun ke BE User (akun nonaktif/dihapus tidak bisa memperpanjang sesi). Kelola user/role/aplikasi/akses dan SSO launch diteruskan ke BE User (helper `helper/forward`; id URL wajib angka); pemeriksaan role admin pindah ke BE User (role terbaru, `RoleMiddleware` dibuang); setelah reset password atau ubah status, sesi Portal user itu dihapus; setelah hapus user, profil lokal ikut dihapus. `/api/sso/exchange` dibuang dari Portal (aplikasi menukar kode langsung ke BE User). Domain user/role/app/sso Portal dipangkas jadi handler + service penerus (dto/model/repo lama dibuang), `pkg/bcrypt` dibuang. Penyamaran log `helper/logmask` dipasang (keputusan 74). Migrasi `005_auth_service_client.sql` (aditif): `users.auth_user_id` unik, `password` boleh kosong, indeks unik `username` diganti indeks biasa, `sessions.must_change_password`, tabel `auth_credentials`; kolom/tabel lama (`password`, `role_id`, `roles`, `apps`, `user_apps`, `sso_codes`) **belum dihapus**, menunggu data dev dipindah di 3B.5. Terverifikasi: `go build`/`go vet`, uji unit logmask, FE type-check/lint/build; migrasi 005 di DB dev (backup sebelumnya: `testing/auth_3b3/backup/`); **uji API 53/53** (`testing/auth_3b3/api_3b3.js`); BE User dimatikan -> 503; **uji browser `v_portal.js` (laptop) 40/40**; `log_requests` Portal tanpa password/token/kunci/kode SSO asli. Uji gagal yang sudah ada sebelum tahap ini (bukan dari perubahan ini): `pkg/binder` TestBindJSON_LargeNumbers dan paket `helper` (butuh `.env` saat uji), juga terjadi di BE User. Belum: uji sentuh tablet/HP (bagian 3B.6); **SSO Portal -> Toko di dev terputus sampai 3B.4** (Toko masih menukar kode ke alamat lama); aplikasi `toko`/`coffee` di `pos_auth_db` masih memakai path produksi (`/toko/`, `/coffee/`), dirapikan di 3B.5. Data dev: kunci uji portal dari `testing/auth_3b2/keys.txt` terpasang di `auth_credentials` Portal dev, user uji `e2eu39038` (dan lainnya) dibiarkan, aplikasi uji `uji-e2e-39038`/`ujiNNNNN` dibiarkan nonaktif. |
| 60 | User menyetujui 3B.3. Sebelum 3B.4 user memutuskan keputusan 75 (hapus user di Toko tetap ada, hanya profil lokal), 76 (edit user Toko hanya role, nama disalin dari BE User saat login), 77 (tautan lama di dev dikosongkan lalu ditautkan ulang otomatis). **Tahap 3B.4 selesai (Toko menjadi klien BE User):** BE Toko pindah ke port **8082** (`.env*`), FE Toko dev menunjuk 8082 (`.env.development`, proxy Vite). Config `PortalUrl` diganti `AuthUrl`, flag `DirectLoginForceEnabled` dibuang dari config dan respons `login-options` (keputusan 73). Klien `pkg/authclient` (salinan Portal, ditambah `Configured()` di kedua aplikasi) menggantikan `pkg/portal`; `errors.UpstreamError`, `repository/auth_credentials_repository.go`, `helper/logmask` disalin dari Portal. Login form: diverifikasi BE User (dengan `client_ip`); pesan BE User diteruskan (salah, nonaktif, tanpa akses, "Ganti password dulu di Portal"); kunci belum dipasang -> 503. Login form dan SSO memakai satu alur profil lokal: cari lewat `auth_user_id` (nama disalin), atau tautkan ke profil dengan username sama yang belum tertaut (keputusan 67), atau buat profil `menunggu_akses` tanpa password lokal; username yang sudah dipakai profil tertaut ke akun lain ditolak; profil yang dinonaktifkan admin Toko ditolak ("Akun tidak aktif di Toko"). SSO menukar kode langsung ke BE User (`/sso/exchange`, kode lintas aplikasi ditolak). Toggle login langsung tetap (aktif bila layanan akun belum diatur). Refresh token tetap lokal (sesi berjalan tetap bekerja saat BE User mati, keputusan 66). Manajemen User BE: endpoint `create` dan `change-password` dibuang, `update` hanya `role_id`; FE: tombol Tambah User dan Ganti Password dibuang, Edit User menampilkan username/nama terkunci + role, keterangan "Akun baru dan password diatur di Portal". Migrasi `012_auth_service_client.sql`: `portal_user_id` -> `auth_user_id` (indeks ikut diganti nama) lalu dikosongkan, `password` boleh kosong, tabel `auth_credentials` baru; `portal_credentials` tidak dipakai lagi tetapi **belum dihapus**. Terverifikasi: `go build`/`go vet`, uji unit logmask, FE type-check/lint/build; migrasi 012 di DB dev (backup sebelumnya: `testing/auth_3b4/backup/`); **uji API 31/31** (`testing/auth_3b4/api_3b4.js`: kunci belum dipasang, login form benar/salah/tanpa akses/wajib ganti, profil baru `menunggu_akses`, tautan otomatis ke profil lama dengan role tetap, nama tersalin, username bentrok ditolak, SSO Portal -> Toko termasuk sekali pakai dan kode coffee ditolak, profil nonaktif ditolak di form dan SSO, endpoint user lama hilang, edit hanya role, toggle login langsung + SSO tetap jalan, satu sesi per user, log tanpa rahasia); **browser `v_toko.js` (laptop) 31/32**: satu gagal (T5 "kas hari ini terbuka") karena data dev, kas tanggal 12 Nov masih terbuka dan uji berjalan lewat tengah malam (tidak terkait login; transaksi kasir tetap lolos); **browser Manajemen User + login form 4/4** (`testing/auth_3b4/b_toko_users.js`); log Toko tanpa password/token asli. Belum: uji sentuh tablet/HP (3B.6); skrip `cmd/migrate_users_to_portal` masih memakai kolom/tabel lama (ditulis ulang di 3B.5); sesi Toko user yang dinonaktifkan/dihapus di BE User tetap berjalan sampai habis atau dinonaktifkan admin Toko (konsekuensi keputusan 66). Data dev: aplikasi `toko` di `pos_auth_db` diubah ke path `http://localhost:3001/`, kunci uji toko terpasang di `auth_credentials` Toko dev, akun uji `tokoNNNNN` dan profil Toko-nya dibiarkan. |
| 61 | User menyetujui 3B.4. Sebelum 3B.5 user memutuskan keputusan 78 (data dev Portal tidak dipindah) dan 79 (kolom/tabel akun lama dibuang sekarang). **Tahap 3B.5 selesai:** skrip `TOKO/BE/cmd/migrate_users_to_auth` menggantikan `migrate_users_to_portal` (dihapus): aturan sama dengan 3.4 tetapi tujuannya BE User (`--auth-db`/`--auth-dsn`/`AUTH_MIGRATION_DSN`); dry-run default, `--apply --i-have-backup`; hash disalin apa adanya; role `admin` Toko -> `admin` akun, lainnya `user` (role Toko tidak berubah); akses `toko`; aksi BUAT / AMBIL ALIH (akun seed) / TAUTKAN (hash sama) / SUDAH TERTAUT / KONFLIK (tidak menebak; bila orang yang sama cukup login ke Toko dengan password BE User, keputusan 67) / TAUTAN YATIM / TANPA PASSWORD (profil yang dibuat lewat BE User). Migrasi Portal `006_drop_legacy_account_tables.sql` (buang `sso_codes`, `user_apps`, `apps`, `roles`, kolom `users.password`, `role_id`, `must_change_password`) dan Toko `013_drop_portal_credentials.sql`; kolom `users.password` Toko tetap (dibaca skrip saat cutover). Terverifikasi di DB dev (backup sebelum `--apply` dan sebelum migrasi pembersihan di `testing/auth_3b5/backup/`): dry-run dengan fixture uji untuk semua jenis aksi (`testing/auth_3b5/setup_fixtures.sh`), `--apply` (12 dibuat, 1 diambil alih, 1 ditautkan, 14 profil tertaut; 1 konflik, 2 tautan yatim dari data uji 3B.4, 1 tanpa password dilaporkan dan dilewati), dijalankan ulang tanpa perubahan (aman diulang), hash akun BE User sama dengan Toko dan tanpa wajib ganti password; **uji login 8/8** (`testing/auth_3b5/verify_logins.js`: password lama Toko berlaku di Toko dan Portal, profil dan role Toko tetap, SSO Portal -> Toko, akun seed yang diambil alih tidak lagi menerima `admin123`, konflik tidak ditimpa); migrasi 006 dan 013 berjalan; **regresi API Portal 53/53, Toko 31/31, login 8/8**. Belum: uji browser setelah migrasi pembersihan (bagian 3B.6); `docs/PANDUAN_CUTOVER_PORTAL.md` masih menyebut skrip lama, diperbarui di fase deploy (keputusan 61). Catatan dev: akun `admin` di BE User masih memakai password bawaan `admin123` (sengaja untuk dev, wajib diganti di production). Data dev: akun hasil salin (Banana, BramaGanteng, kasir/owner uji, e2e...) dan fixture `migl0121`/`migs0121`/`migk0121`/`mign0121` dibiarkan. |
| 62 | User menyetujui 3B.5. **Tahap 3B.6 selesai (konfigurasi menyeluruh + uji ujung ke ujung):** `nginx/pos-web.conf`: `/portal-api/` -> 127.0.0.1:8081, `/toko/api/` dan jalur lama `/api/` -> 127.0.0.1:8082, Coffee nanti 8083, catatan bahwa BE User (127.0.0.1:8080) sengaja tidak dipasang. README utama (struktur dengan `AUTH/BE`, tabel layanan/port/database, urutan menjalankan, kunci lewat `cmd/set_app_key` ke `auth_credentials`, path dev `toko`, skrip `migrate_users_to_auth`, kredensial bawaan), peta di `CLAUDE.md`, README `TOKO/BE` dan `TOKO/FE`, Dockerfile/docker-compose lama Toko (8082), README dan `.env*` `COFFEE/BE` (8083, agar tidak bentrok dengan Toko; isi Coffee lain diulang di 4.1), HANDOFF. `.env*`, config, dan proxy Vite Portal/Toko sudah disesuaikan di 3B.3/3B.4. Uji ujung ke ujung (`testing/auth_3b6/run_all.sh`, jendela tampil, jeda 0,5 detik, ditutup tiap bagian): **Portal 40/40 di laptop, tablet, dan HP**; **Toko 31/32 di laptop, tablet, dan HP** (satu-satunya gagal T5 "kas hari ini terbuka" karena kas dev tanggal 12 Nov masih terbuka, sama seperti 3B.4, tidak terkait akun; kasir dan transaksi lolos); **Manajemen User + login form Toko 4/4**; `log_requests` BE User, Portal, Toko tanpa password asli. Belum diubah (menunggu fase deploy, keputusan 61): `docs/SOP_DEPLOY.md`, `DEPLOYMENT_PROD.md`, `DEPLOYMENT_REDEPLOY_FULL.md`, `PANDUAN_DEPLOY.md`, `PANDUAN_CUTOVER_PORTAL.md` (port lama, skrip lama, service `pos-auth-backend` belum ada); uji di tablet/HP fisik dan APK. **Fase 3B selesai bila 3B.6 disetujui**; berikutnya kembali ke Fase 4 tahap 4.1 Coffee (diulang, port 8083, salin Toko versi baru). |
| 63 | User menyetujui 3B.6 (Fase 3B selesai) dan rencana 4.1 yang diperbarui (bagian 5.0c), termasuk mengganti isi `COFFEE/BE`, membuat ulang `pos_coffee_db` dev, dan mengubah path `coffee` di BE User dev. **Tahap 4.1 (diulang) selesai:** isi `COFFEE/BE` lama dipindah ke `testing/coffee_4_1_redo/old_COFFEE_BE` sebagai bahan (tetap ada di riwayat git), lalu `TOKO/BE` versi HEAD disalin ulang; pemangkasan 4.1 lama diterapkan ulang lewat patch yang dibuat dari selisih Toko dan Coffee pada commit `a18932f` (10 file kode: rute, saldo pelanggan di transaksi, batch kedaluwarsa dan retur di pembelian/supplier, nama toko bawaan), modul/folder dipangkas, `cmd/` dibuang seluruhnya (termasuk `migrate_users_to_auth`, tidak diperlukan Coffee), module `coffee_api`, teks "Toko" di kode akun diganti "Coffee". Baseline `001`/`002` dari 4.1 lama disesuaikan Fase 3B (`auth_user_id`, `password` boleh kosong, `auth_credentials` menggantikan `portal_credentials`, profil seed `owner`/`admin` tanpa password lokal); `.env*`/config dari 4.1 lama dengan `AuthUrl` menggantikan `PortalUrl` dan tanpa flag darurat (kunci config sama dengan Toko). Data dev: backup `pos_coffee_db` lama di `testing/coffee_4_1_redo/backup/`, DB dibuat ulang (31 tabel), path aplikasi `coffee` di BE User dev `http://localhost:3002/`, kunci uji coffee terpasang di `auth_credentials`. **Temuan dan perbaikan di Toko (bug 3B.4):** `TOKO/BE/domain/user/dto/dto_user.go` (satu-satunya file CRLF di domain user) tidak ikut terubah di 3B.4, sehingga `UpdateRequest` masih mewajibkan `full_name` dan Edit User dari FE Toko (yang hanya mengirim `role_id`) akan ditolak 422; uji 3B.4 tidak menangkapnya karena T23 ikut mengirim `full_name` dan uji browser B4 tidak menyimpan. Diperbaiki (DTO hanya `role_id`, tipe buat user/ganti password dibuang) di Toko dan Coffee; audit penggantian otomatis Fase 3B lainnya tidak menemukan yang terlewat; uji T23 diperketat (body hanya `role_id`) dan uji browser B5 (simpan Edit User) ditambahkan. Kesalahan kecil selama tahap: komentar seed berisi titik koma memutus runner migrasi (sudah diketahui), diperbaiki. Terverifikasi: `go build`/`go vet` Coffee dan Toko; BE Coffee hidup di **8083**, migrasi 001-002 bersih; **uji API Coffee 20/20** (`testing/coffee_4_1_redo/api_4_1_redo.js`: kunci belum dipasang -> 503, login form admin menautkan profil seed, tanpa akses/wajib ganti ditolak, akun baru `menunggu_akses` tanpa password lokal, SSO Portal -> Coffee sekali pakai, kode Toko ditolak, `/users/create` tidak ada, edit role, toggle login langsung, uji asap 59 endpoint baca tanpa 500, menu tanpa modul dipangkas, endpoint modul dipangkas tidak ada, log tanpa rahasia); **regresi Toko API 31/31**, **browser Manajemen User Toko 5/5**. Belum: FE Coffee (4.2), uji browser Coffee. |
| 64 | User menyetujui 4.1 (diulang). **Tahap 4.2 selesai (kerangka FE Coffee):** `COFFEE/FE` disalin dari `TOKO/FE` versi HEAD. Identitas: package `coffee-web`, port 3002, proxy dan `VITE_API_URL` dev ke 8083, `VITE_APP_NAME` "POS Coffee", production base `/coffee/` dan API `/coffee/api`, judul tab "POS Coffee", key browser `coffee.*` (migrasi key lama `bootstrap/migrateStorageKeys` dibuang, Coffee tidak punya key lama), teks "Toko" -> "Kedai"/"Coffee" (label mengikuti seed BE: "Profil Kedai"), README disesuaikan, `DEPLOYMENT.md` Toko tidak disalin (deploy di fase deploy). Dipangkas: halaman historis (backdate kas/kasir), retur supplier, rekonsiliasi modal, versi aplikasi (folder, rute, konstanta, ekspor, kunci query); batch kedaluwarsa (peringatan dan riwayat di Produk, rincian di form/detail Pembelian); saldo pelanggan (kolom dan hook API di Pelanggan, "simpan kembalian ke saldo" di struk kasir); kartu dan tab retur di detail Supplier. Logika saldo di layar bayar kasir dibiarkan (tidak aktif karena saldo selalu 0, dibersihkan bersama kolom BE-nya di 4.7, sesuai catatan baseline 4.1). Menu sidebar mengikuti seed menu Coffee dari BE (tanpa modul dipangkas). Penggantian teks memakai skrip Node dengan pengecekan cocok tepat satu kali (`testing/coffee_4_2/prune_*.js`) setelah dua kesalahan delimiter regex yang segera diperbaiki. **Temuan:** skrip `type-check` (`tsc --noEmit`) di Toko, Portal, dan Coffee **tidak memeriksa file apa pun** karena `tsconfig.json` berisi `"files": []` + references; pemeriksaan tipe sebenarnya selama ini hanya lewat `npm run build` (`tsc -b`). Diperbaiki di Coffee (`tsc -p tsconfig.app.json --noEmit`); Toko dan Portal belum diubah (menunggu keputusan user), tetapi sudah diperiksa dengan perintah yang benar: Toko, Portal, Coffee tanpa galat. Terverifikasi: type-check (yang benar), lint, build (base `/coffee/` di production); **uji browser 12/12 di laptop, tablet, dan HP** (`testing/coffee_4_2/b_coffee.js`): login form Coffee lewat BE User, 27 halaman Coffee terbuka tanpa panggilan API gagal, halaman modul dipangkas dialihkan, Manajemen User tanpa Tambah User, toggle login langsung, logout ke Portal, kartu POS Coffee di Portal -> `/sso` -> dashboard, tombol Kembali ke Portal, user baru -> Menunggu Akses, tanpa error JavaScript; log Coffee tanpa password asli. Data dev: akun uji `coffNNNNN` (3 buah, satu per perangkat) dibiarkan. Belum: fitur kopi (4.3 dan seterusnya). |
| 65 | User menyetujui 4.2 dan perbaikan skrip `type-check` di Toko dan Portal (kini `tsc -p tsconfig.app.json --noEmit`, memeriksa 441 dan 122 file; Coffee 395; ketiganya lolos). Sebelum 4.3 user memutuskan keputusan 80-84 (input pemakaian/terbuang oleh kasir/admin/owner, koreksi dengan pembatalan entri oleh admin/owner, satuan dari master Unit, batas stok menipis dan alasan wajib saat terbuang, tanpa kolom stok awal, **opname dihapus** untuk bahan baku dan barang jadi). **Tahap 4.3 selesai (bahan baku):** BE migrasi `003_ingredients.sql` (tabel `ingredients` dan `ingredient_mutations` DECIMAL(15,3), entri tidak pernah diubah atau dihapus kecuali penanda pembatalan, menu grup "Bahan Baku" urutan 4 dengan "Daftar Bahan" `/ingredients` dan "Riwayat Stok Bahan" `/ingredients/history`, izin sesuai keputusan 80). Domain `ingredient` (pola sama dengan Toko): CRUD bahan (nama unik, satuan aktif, satuan terkunci bila sudah ada riwayat, hapus hanya bila tanpa riwayat, nonaktif tidak bisa dicatat), catat terpakai dan terbuang (dibulatkan 3 desimal, stok minus diizinkan dengan peringatan), riwayat dengan filter bahan/jenis/tanggal, pembatalan entri (stok dan total dikembalikan, entri asal ditandai, entri pembatalan baru merujuk entri asal); setiap perubahan stok dalam transaksi dengan baris bahan dikunci (`FOR UPDATE OF i`). FE fitur `ingredients`: halaman Daftar Bahan (filter cari/status/menipis, kolom stok sisa/total terpakai/total terbuang/batas menipis, tanda "Menipis", aksi per izin: catat terpakai, catat terbuang, riwayat, edit, aktif/nonaktif, hapus bila tanpa riwayat), form bahan (satuan dari master Unit, keterangan stok mulai 0), form catat (input desimal koma/titik, peringatan stok akan minus, alasan wajib untuk terbuang), halaman Riwayat (filter, jumlah bertanda +/-, entri dibatalkan dicoret dengan siapa/kapan/alasan, tombol batalkan hanya untuk izin hapus). **Temuan:** komponen bersama `Combobox` mencari berdasarkan id, bukan label (ketik nama tidak menemukan apa pun); diperbaiki di salinan Coffee (`keywords` = label), **Toko belum diubah** (di luar tahap ini, menunggu keputusan user). Terverifikasi: `go build`/`go vet`, FE type-check/lint/build; migrasi 003 di DB dev; **uji API 33/33** (`testing/coffee_4_3/api_4_3.js`: CRUD dan validasi, izin kasir/admin/menunggu akses, terpakai/terbuang termasuk minus dan pembulatan, riwayat dan filter, pembatalan termasuk aturan sekali batal dan hanya terpakai/terbuang, satuan terkunci, hapus/nonaktif, **20 pencatatan bersamaan konsisten**); **uji browser 10/10 di laptop, tablet, dan HP** (`testing/coffee_4_3/b_4_3.js`: admin tambah bahan, tanda menipis, kasir tanpa tambah/edit, kasir catat terpakai 0,5 dengan peringatan minus, terbuang wajib alasan, riwayat terfilter tanpa tombol batalkan untuk kasir, admin membatalkan entri, pencarian nama di filter riwayat, satuan terkunci, tanpa error JavaScript); tata letak kolom riwayat dirapikan agar muat di laptop. Data dev: satuan "Bungkus" dan "Liter" (dibuat uji), bahan dan akun uji (`kasirNNNNN`, `baristaNNNNN`, `tungguNNNNN`) dibiarkan. Belum: stok bertambah dari pembelian (tahap 4.5). |
| 66 | User menyetujui 4.3 dan perbaikan Combobox Toko, lalu memutuskan keputusan 85-90 (ukuran dikelola admin, resep baris bahan + takaran, kelompok tambahan wajib/opsional + satu pilihan bawaan tanpa batas maksimal, kasir boleh menandai habis, barang jadi di 4.4 hanya mode + harga, perbaikan Combobox Toko). Keputusan 41 diberi catatan opname dihapus (keputusan 84). **Tahap 4.4 selesai (menu jualan):** BE migrasi `004_menu_items.sql` (tabel `menu_categories`, `menu_sizes`, `menu_items`, `menu_item_variants` (setiap menu punya minimal satu baris varian, menu tanpa varian = satu baris tanpa ukuran/suhu), `menu_recipe_lines`, `modifier_groups`, `modifier_options`, `menu_item_modifier_groups`; grup menu navigasi "Menu" urutan 3 berisi Daftar Menu, Kategori Menu, Kelompok Tambahan, Ukuran; izin keputusan 88). Domain `menu_category`, `menu_size`, `modifier_group`, `menu_item` (pola Toko). Aturan: nama unik, kategori aktif, barang jadi tanpa varian dan tanpa resep (keputusan 43), menu tanpa varian wajib harga, kombinasi varian = ukuran x suhu (salah satu boleh kosong, tidak keduanya), tidak ganda, minimal satu dijual, ukuran aktif, bahan resep ada dan tidak ganda per kombinasi, kelompok tambahan aktif; ubah menu mempertahankan id kombinasi yang sama (dipakai pesanan di 4.7), kombinasi yang dihapus dari form dibuang beserta resepnya; kelompok tambahan: paling banyak satu pilihan bawaan dan harus aktif, minimal satu pilihan aktif, nama pilihan unik, id pilihan harus milik kelompok; hapus kategori/ukuran/kelompok hanya bila tidak dipakai; tandai habis cukup izin lihat. Bahan baku kini juga tidak bisa dihapus bila masih dipakai resep (pesan jelas, sebelumnya akan jatuh ke galat foreign key). FE fitur `menu-items`: Daftar Menu (filter cari/kategori/mode/habis, harga atau rentang harga dari kombinasi yang dijual, aksi tandai habis/tersedia, edit, aktif/nonaktif, hapus sesuai izin), form menu (nama, kategori, mode racikan/barang jadi, saklar varian, pilih ukuran dan suhu menghasilkan tabel kombinasi dengan harga, centang dijual, resep per kombinasi; menu tanpa varian: harga dan resep; centang kelompok tambahan), halaman Kategori Menu dan Ukuran (satu komponen master bersama), Kelompok Tambahan (aturan, pilihan dengan harga, bawaan, aktif). **Perbaikan Toko (keputusan 90):** Combobox Toko mencari berdasarkan label. **Temuan (Toko, belum diubah):** filter supplier di halaman Pembelian Toko hanya berisi 10 supplier pertama (daftar dikirim dari halaman induk), sehingga supplier lain tidak bisa dipilih di filter itu. Terverifikasi: `go build`/`go vet`, FE type-check/lint/build Coffee dan Toko; migrasi 004 di DB dev; **uji API 40/40** (`testing/coffee_4_4/api_4_4.js`: kategori, ukuran, kelompok tambahan dan semua aturannya, menu tanpa/dengan varian, varian hanya suhu, barang jadi, id kombinasi dipertahankan saat ubah, ubah ke tanpa varian, hapus dan pengamannya, izin kasir/menunggu akses, tandai habis); regresi bahan baku 33/33; **uji browser 7/7 di laptop, tablet, dan HP** (`testing/coffee_4_4/b_4_4.js`: kategori, ukuran, kelompok tambahan, menu Latte 2 ukuran x panas/dingin dengan satu kombinasi tidak dijual + resep + kelompok, barang jadi, edit menu mempertahankan id kombinasi dan resep, kasir menandai habis, tanpa error JavaScript); **uji browser Combobox Toko 1/1** (`testing/coffee_4_4/b_toko_combobox.js`, ketik sebagian nama supplier: 11 -> 1 opsi). Data dev: kategori, ukuran, kelompok, menu, bahan, dan akun uji dibiarkan. Belum: stok barang jadi dan pembelian (4.5), pemakaian menu di kasir (4.7). |
| 67 | User menyetujui 4.4 dan memutuskan keputusan 91-94 (rekomendasi AI): modul Produk salinan Toko menjadi katalog **barang jadi**, menu barang jadi ditautkan ke tepat satu produk, pembelian dua jenis baris, bayar/void/edit seperti Toko dengan stok bahan boleh minus, form Produk tidak diubah selain keterangan harga jual. User juga meminta perbaikan filter supplier Toko. **Tahap 4.5 selesai (barang jadi dan pembelian):** BE migrasi `005_purchase_ingredients.sql` (`purchase_items.item_type` bawaan `product` + `ingredient_id` FK RESTRICT; `menu_items.product_id` unik FK RESTRICT; label grup navigasi "Produk" menjadi "Barang Jadi"). Pembelian (`supplier_purchase`): baris `product` tetap logika Toko (produk + kemasan, stok per kemasan dalam satuan dasar), baris `ingredient` memakai bahan aktif, satuan diisi singkatan satuan stok bahan, konversi 1, jumlah desimal; stok bahan bertambah lewat `ingredient_repo.ApplyStock` (baris bahan dikunci, riwayat jenis baru `purchase`/`purchase_void` dengan rujukan nota). Edit nota menerapkan selisih per bahan (tambah = `purchase`, kurang = `purchase_void`), void mengurangi stok bahan (boleh minus, keputusan 93), tambah item juga mendukung bahan; sisa kolom batch kedaluwarsa dibuang dari DTO/model/service. Entri pembelian tidak bisa dibatalkan dari Riwayat Stok Bahan (aturan 4.3 tetap: hanya terpakai/terbuang). Menu (`menu_item`): mode barang jadi wajib produk aktif yang belum dipakai menu lain, mode racikan selalu tanpa produk; daftar dan detail menu memuat `product_id`/`product_name`. Produk yang dipakai menu tidak bisa dihapus (pesan jelas, sebelumnya akan jatuh ke galat foreign key). FE: form Tambah/Edit Pembelian dan Tambah Item memakai komponen baru `PurchaseItemTargetCell` (pilih jenis baris Barang jadi/Bahan baku, lalu pilih produk atau bahan), baris bahan desimal dengan satuan stok bahan, cek ganda per jenis, peringatan stok minus hanya untuk barang jadi; detail nota menampilkan jenis baris; riwayat bahan punya label dan filter "Pembelian" dan "Void pembelian"; daftar bahan diperbarui setelah perubahan pembelian; form menu barang jadi memilih produk (Combobox), daftar menu menampilkan nama produk; form Produk diberi keterangan bahwa harga jual kasir diatur di Daftar Menu. **Perbaikan komponen Coffee:** `Combobox` bersama kini bisa dipakai di dalam modal (z-index dan pointer-events sama dengan `AsyncCombobox`, tombol `type="button"`); sebelumnya hanya dipakai di bar filter. **Perbaikan Toko (diminta user):** filter supplier di halaman Pembelian dan Retur Toko memakai `useSupplierOptionsQuery` (semua supplier aktif), sebelumnya `limit: 200` melebihi batas BE 100 sehingga jatuh ke 10 supplier. Terverifikasi: `go build`/`go vet`, FE type-check/lint/build Coffee dan Toko; migrasi 005 di DB dev; **uji API 38/38** (`testing/coffee_4_5/api_4_5.js`: validasi baris, nota campuran, stok bahan dan barang jadi, riwayat dan filter, edit dengan selisih, bayar sebagian, tambah item bahan, void dengan stok minus, baris lama tanpa `item_type`, bahan yang pernah dibeli tidak bisa dihapus, aturan tautan menu-produk, produk dipakai menu tidak bisa dihapus, label navigasi); regresi API bahan baku 33/33 dan menu 40/40 (uji 4.4 disesuaikan: menu barang jadi kini memakai produk); **uji browser 7/7 di laptop, tablet, dan HP** (`testing/coffee_4_5/b_4_5.js`: nota campuran bahan 1.500,5 ml + 12 croissant, detail nota, riwayat "Pembelian" tanpa tombol batalkan, menu barang jadi menolak simpan tanpa produk lalu tersimpan tertaut, filter supplier Toko memuat 49/49 supplier aktif dan supplier terakhir bisa dipilih, tanpa error JavaScript); regresi browser menu 7/7 (laptop, tablet, HP) dan bahan baku 10/10 (laptop). Data dev: supplier, bahan, produk, nota, dan menu uji dibiarkan. Belum: stok barang jadi berkurang saat dijual dan pemakaian menu di kasir (4.7); retur supplier tidak ada di Coffee (dipangkas). |
| 68 | User menyetujui 4.5 dan memutuskan keputusan 95-98 (semua rekomendasi AI): nomor antrean untuk semua pesanan, diambil saat pesanan dibuat dan tidak dipakai ulang, data meja nama + urutan + aktif, halaman Meja di grup Menu dengan izin admin/owner kelola dan kasir lihat. **Tahap 4.6 selesai (meja dan antrean):** BE migrasi `006_tables_queue.sql` (tabel `dining_tables` nama unik/urutan/aktif; tabel `queue_counters` satu baris per tanggal dengan nomor terakhir; menu navigasi `menu_jual.meja` "Meja" `/menu-items/tables` urutan 5; izin owner/admin penuh, kasir lihat). Domain `dining_table` (pola sama dengan `menu_size`: daftar dengan cari/filter aktif, pilihan meja aktif, tambah, ubah, aktif/nonaktif, hapus; nama dirapikan dan unik; `usage_count` sementara 0, mulai 4.7 dihitung dari pesanan dan meja yang sudah dipakai tidak bisa dihapus). Rute `/api/dining-tables/*` dengan izin `menu_jual.meja` (pilihan meja cukup izin lihat agar kasir bisa memakainya di 4.7). Nomor antrean: `domain/queue_number/repo` berisi `NextQueueNumber(tx, now)` (INSERT ... ON DUPLICATE KEY UPDATE lalu SELECT ... FOR UPDATE pada baris tanggal, di dalam transaksi pesanan; tanggal mengikuti zona waktu `now` = `time_helper.GetTimeNow()` sehingga reset 00:00 WIB) dan `FormatQueueNumber` (001, 002, ...); belum ada tampilan karena baru dipakai pesanan di 4.7. FE: halaman Meja (komponen `SimpleMasterPage` yang sama dengan Ukuran/Kategori Menu), rute, query key, API. Terverifikasi: `go build`/`go vet`, FE type-check/lint/build; migrasi 006 di DB dev; **`go test` nomor antrean 5/5** (20 pengambilan bersamaan menghasilkan 1-20 tanpa ganda, tanggal baru mulai dari 1, transaksi yang di-rollback tidak memakai nomor, batas 00:00 WIB, format tiga digit; baris uji memakai tanggal 1990-an yang kosong dan dibiarkan di DB dev); **uji API 17/17** (`testing/coffee_4_6/api_4_6.js`: tambah, urutan, nama dirapikan/ganda/kosong/terlalu panjang, ubah, nonaktif tidak muncul di pilihan, filter, izin kasir lihat saja dan 403 untuk kelola, menunggu akses 403, isi izin, hapus, 404, posisi navigasi); regresi API 33/33, 40/40, 38/38; **uji browser 6/6 di laptop, tablet, dan HP** (`testing/coffee_4_6/b_4_6.js`: admin menambah tiga meja dan urutannya, ubah nama, nonaktifkan, item "Meja" di navigasi, kasir melihat tanpa tombol kelola, tanpa error JavaScript). **Temuan (belum diubah):** uji integrasi salinan Toko `COFFEE/BE/domain/product/repo/stock_delta_repo_integration_test.go` masih terhubung ke `pos_retail_db` (DB dev Toko), sehingga `go test ./...` di Coffee akan menulis data uji ke DB Toko; sebaiknya diarahkan ke `pos_coffee_db`. Data dev: meja dan akun uji (`kasirtNNNNN`, `tunggutNNNNN`, `mejakNNNNN`) dibiarkan. Belum: pemakaian meja dan nomor antrean di pesanan/kasir (4.7), tanda meja yang sedang terisi (4.7/4.8). |
| 69 | User menyetujui 4.6 dan perbaikan uji integrasi stok produk Coffee. Sebelum 4.7 user memutuskan keputusan 99-102 (nama pemesan opsional; void lunas dengan pilihan per item barang jadi kembali ke stok/dibuang, alasan wajib, usulan user; utang saat bayar sekarang seperti Toko; 4.7 dipecah 4.7a/4.7b). **Perbaikan uji integrasi stok produk (salinan Toko):** `COFFEE/BE/domain/product/repo/stock_delta_repo_integration_test.go` kini memakai `pos_coffee_db` (sebelumnya `pos_retail_db`) dan dua satuan dari master Unit Coffee (sebelumnya id satuan Toko 4 dan 11); id produk/kemasan dibaca berdasarkan nama/kunci, bukan `LAST_INSERT_ID()` lewat pool koneksi (bisa terbaca dari koneksi lain); nama produk uji unik per jalan; teardown menghapus kemasan anak lebih dulu (sebelumnya penghapusan ditolak foreign key `ref_package_id` sehingga produk uji tertinggal). Paket ini mengimpor `config` yang membaca `./.env` saat init, jadi dijalankan sebagai biner dari folder `COFFEE/BE` (cara di komentar uji). Lolos dua kali berturut-turut tanpa sisa; satu produk uji lama ("TEST STOCKDELTA CONCURRENCY 1", id 19) dari jalan sebelum perbaikan dibiarkan di DB dev. Catatan: perbaikan semula di riwayat #68 menyebut uji ini menulis ke DB Toko; kenyataannya uji tidak pernah jalan lewat `go test ./...` (panik membaca `.env`). **Tahap 4.7a selesai (pesanan BE):** migrasi `007_orders.sql`: `transactions` + `order_type`, `table_id` (FK `dining_tables` RESTRICT), `table_name`, `queue_date`, `queue_number`, `order_name`; `transaction_items` + `menu_item_id` (FK RESTRICT), `variant_id` (FK SET NULL), `variant_name`, `notes`; tabel baru `transaction_item_modifiers` (nama kelompok, nama pilihan, harga saat dijual); `stock_mutations.mutation_type` + `waste`. Domain `transaction` (pola Toko, kode transaksi dan kas tetap): permintaan pesanan berisi jenis pesanan, meja, nama pemesan, item (kombinasi varian, jumlah bulat, pilihan tambahan, catatan, diskon item), diskon nota, pajak, pembayaran; **harga dihitung server** (harga kombinasi + tambahan, dikali jumlah, dikurangi diskon), subtotal/total/kembalian dihitung server. Aturan: makan di tempat wajib meja aktif, bawa pulang tanpa meja; menu nonaktif/habis dan kombinasi yang tidak dijual ditolak; menu barang jadi wajib tertaut produk; tambahan harus dari kelompok yang terpasang di menu, kelompok dan pilihan aktif, tidak ganda, kelompok pilih satu paling banyak satu, kelompok wajib minimal satu; diskon tidak melebihi harga; utang wajib memilih pelanggan (piutang tercatat seperti Toko). Barang jadi: stok produk dipotong (kemasan anchor) di transaksi yang sama, **penjualan ditolak bila stok kurang** dengan pesan nama menu; nomor antrean diambil di transaksi yang sama sehingga pesanan gagal tidak memakai nomor. Void (admin/owner): body opsional berisi pilihan per item barang jadi `return`/`waste` + alasan (wajib bila ada yang dibuang); semua barang jadi kembali ke stok lalu yang dibuang dikurangi lagi sebagai mutasi `waste` beserta alasannya; racikan tidak bisa dipilih dibuang; kas tunai dan piutang ikut dikembalikan seperti Toko. Daftar transaksi: cari juga nama pemesan, meja, nomor antrean; filter jenis pesanan; detail memuat kombinasi, catatan, tambahan, antrean. Menu yang pernah dipesan tidak bisa dihapus; pemakaian meja dihitung dari pesanan dan meja yang pernah dipakai tidak bisa dihapus; rumus rekonsiliasi stok ikut mengurangi `waste`. `pkg/pricing` (harga produk ala Toko) dihapus karena tidak dipakai lagi. Terverifikasi: `go build`/`go vet`; migrasi 007 di DB dev; **uji API 49/49** (`testing/coffee_4_7a/api_4_7a.js`: hitung harga dan kembalian, nama kombinasi/catatan/tambahan tersimpan, antrean urut, kas hanya tunai, aturan meja, aturan tambahan, menu habis/nonaktif/tidak dijual/tidak ada, jumlah pecahan, diskon item/nota/pajak, pembayaran kurang, potong stok barang jadi, stok kurang diblokir tanpa sisa transaksi dan tanpa memakai nomor antrean, utang, detail/daftar/cari/filter, izin void, void kembali/dibuang dengan alasan, kas berkurang saat void, void ulang ditolak, void tanpa body, void utang, keseimbangan mutasi stok, pengaman hapus menu dan meja); **uji bersamaan** (`api_4_7a_bersamaan.js`: 10 pesanan serentak untuk stok 5 -> tepat 5 berhasil, stok 0, nomor antrean unik); regresi API 33/33, 40/40, 38/38, 17/17; `go test` antrean 5/5 dan uji integrasi stok produk lolos. **Temuan (belum diubah):** `pkg/binder` `TestBindJSON_LargeNumbers` gagal juga di Toko (kegagalan lama, bukan dari Coffee); paket `helper` tidak bisa diuji lewat `go test` karena `config` membaca `./.env` saat init (sama di Toko); bug teardown uji integrasi stok yang sama masih ada di salinan Toko. **Belum:** halaman kasir Coffee masih salinan Toko dan belum memakai permintaan pesanan baru, sehingga transaksi dari halaman kasir lama akan ditolak sampai 4.7b; label "Terbuang" di riwayat stok produk/rekonsiliasi di FE (4.7b). Data dev: menu, produk, meja, pelanggan, transaksi, dan akun uji (`kasiroNNNNN`) dibiarkan. |
| 70 | User menyetujui 4.7a. **Tahap 4.7b selesai (kasir FE, struk, riwayat transaksi):** BE: endpoint katalog kasir `POST /menu-items/catalog` (izin lihat `penjualan.kasir`): kategori aktif yang punya menu, menu aktif beserta kombinasi yang dijual (urut ukuran lalu panas/dingin), kelompok tambahan aktif dengan pilihan aktif dan bawaannya, tanda habis, dan sisa stok barang jadi (cara hitung sama dengan daftar produk); lima query tanpa query per menu. Cetak via server/LAN: payload struk ditambah nomor antrean (dicetak besar), jenis pesanan, meja, nama pemesan, dan baris rincian item. FE kasir (pola Toko, fitur `sales/cashier`): kiri katalog menu (cari, tab kategori, tombol besar dengan harga atau "mulai", "Sisa N" untuk barang jadi, "Habis" dan tidak bisa diketuk bila habis atau stok 0); menu satu kombinasi tanpa tambahan langsung masuk keranjang, lainnya membuka popup pilihan (kombinasi, tambahan wajib/opsional dengan bawaan terpilih, pilih satu/banyak, jumlah, catatan, tombol "Tambah" dengan harga); jumlah barang jadi di keranjang tidak bisa melebihi sisa stok. Kanan panel pesanan: jenis pesanan (bawa pulang/makan di tempat), meja (wajib untuk makan di tempat, tombol Bayar nonaktif tanpa meja), nama pemesan opsional, pelanggan, daftar item (kombinasi, tambahan, catatan, jumlah, diskon item), diskon/pajak, total, Bayar. Baris yang sama (kombinasi + tambahan + catatan) digabung jumlahnya; draf keranjang tersimpan di perangkat (versi baru, draf lama salinan Toko dibuang). Tablet/HP: tab Menu dan Pesanan. Pembayaran memakai modal Toko tanpa saldo pelanggan (dipangkas), utang tetap (keputusan 101); payload hanya pilihan (harga dihitung server). Struk (pratinjau, browser, PDF, WhatsApp, Bluetooth/serial/RawBT ESC/POS, server/LAN): nomor antrean besar, jenis pesanan, meja, nama pemesan, nama menu + kombinasi, tambahan, catatan. Riwayat transaksi: kolom antrean, jenis pesanan + meja, pemesan; cari kode/antrean/pemesan/meja/pelanggan; filter jenis pesanan; detail memuat antrean, meja, pemesan, kombinasi, tambahan, catatan, tanda barang jadi; cetak ulang struk lengkap; void lewat dialog baru dengan pilihan per item barang jadi "Kembali ke stok"/"Dibuang" dan alasan wajib bila dibuang (keputusan 100). Cetak nota piutang memakai pemeta struk yang sama. Label "Terbuang" di kartu stok dan rincian rekonsiliasi stok. Perbaikan tata letak: tinggi halaman kasir dikurangi padding halaman (sebelumnya tombol Bayar terdorong keluar layar laptop, warisan salinan Toko); diskon/pajak ikut tergulir agar daftar item lega. Dibuang: pencarian produk/barcode kasir Toko (`ProductSearch`, `useProductSearch`, `useBarcodeScan`) dan daftar keranjang lama. Terverifikasi: `go build`/`go vet`, FE type-check/lint/build; **uji API katalog 6/6** (`testing/coffee_4_7b/api_4_7b.js`); regresi API 49/49 (4.7a), 17/17, 38/38; **uji browser 7/7 di laptop, tablet, dan HP** (`testing/coffee_4_7b/b_4_7b.js`: menu sederhana langsung masuk, popup Latte dengan bawaan Normal dan harga 58.000, batas stok donat 3 dengan peringatan, Bayar nonaktif tanpa meja, bayar uang pas, struk berisi antrean/meja/tambahan/catatan, data tersimpan benar dan stok donat 0, keranjang kembali kosong dan donat tampil Habis, admin mencari di riwayat, detail, cetak ulang berisi antrean, void dengan donat dibuang + alasan, tanpa error JavaScript); regresi browser pembelian 7/7, meja 6/6, menu 7/7 (laptop). **Belum diuji:** cetak fisik ke printer (server/LAN/Bluetooth/RawBT/APK) karena tidak ada printer di mesin dev, hanya pratinjau struk yang diuji; cetak nota piutang di halaman Piutang belum diuji di browser. Data dev: menu, produk, meja, transaksi, akun uji (`kaskNNNNN`, `kasircNNNNN`, `tunggucNNNNN`) dibiarkan. |
