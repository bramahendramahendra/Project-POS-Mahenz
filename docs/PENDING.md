# Pending

Pekerjaan yang ditunda. Copas perintahnya ke server saat waktunya tiba. Selesai: pindahkan ke "Sudah selesai" di bawah.

---

## 1. Uji tablet/APK dan cetak Bluetooth (D.5)

Di-skip dulu (keputusan user, 15 Nov 2026). Saat sempat, di tablet:

- [ ] Buka APK → tampil halaman login Portal
- [ ] Login → kartu POS Toko → masuk Toko tanpa login ulang
- [ ] Toko: transaksi kecil → cetak struk Bluetooth → void
- [ ] Pengaturan printer: alamat printer Bluetooth masih tersimpan
- [ ] Logout → kembali ke Portal
- [ ] Coffee (bila dipakai di tablet): masuk lewat kartu POS Coffee → cetak struk

Kalau ada yang gagal: catat langkahnya, kirim ke AI.

## 2. Bersih-bersih server (tunggu stabil beberapa hari setelah cutover 15 Nov 2026)

Isi `TGL` dengan tanggal cutover (nama folder lama: `ls /opt`).

```bash
TGL=20261115
ls /opt                                   # harus: ada pos-mahenz_lama_$TGL
sudo rm -rf /opt/pos-mahenz_lama_$TGL
sudo rm -rf /var/www/pos-web/dist         # FE Toko lama (jalur rollback)
sudo systemctl disable pos-backend 2>/dev/null
sudo rm -f /etc/systemd/system/pos-backend.service && sudo systemctl daemon-reload
systemctl list-units --type=service | grep pos-   # harus: hanya pos-auth/portal/toko/coffee-backend
```

Setelah ini rollback ke versi lama tidak bisa lagi. Backup `~/cutover_$TGL/` boleh disimpan dulu.

## 3. Hapus jalur API lama `/api/` di nginx

Blok `location ^~ /api/` di `nginx/pos-web.conf` hanya untuk perangkat yang masih memakai alamat lama
(FE/APK lama yang ter-cache). Cek dulu apakah masih ada yang memakainya:

```bash
sudo grep -cE '"[A-Z]+ /api/' /var/log/nginx/access.log          # jumlah request ke /api/ lama (hari ini)
sudo grep -cE '"[A-Z]+ /api/' /var/log/nginx/access.log.1        # kemarin
sudo zgrep -hcE '"[A-Z]+ /api/' /var/log/nginx/access.log.*.gz   # hari-hari sebelumnya, satu angka per file
sudo grep -E '"[A-Z]+ /api/' /var/log/nginx/access.log | tail -5 # contoh: IP dan waktu perangkat yang masih lama
```

Pola ini hanya menghitung `/api/...` di awal alamat, tidak ikut menghitung `/toko/api/`, `/coffee/api/`, `/portal-api/`.

- Masih ada angka > 0: buka ulang aplikasi di perangkat itu (tutup penuh APK/tab, buka lagi) lalu cek beberapa hari lagi.
- Semua 0 selama beberapa hari: minta AI menghapus blok itu di `nginx/pos-web.conf`, lalu deploy nginx
  (SOP_DEPLOY.md, bagian "Nginx").

## Urutan yang disepakati (15 Nov 2026)

**Nomor 11 (uji Go, bagian 5) → hapus dead code (bagian 6) → nomor 10 (library FE, bagian 4).**
Uji Go hijau dulu sebagai jaring pengaman; dead code dihapus sebelum memperbarui library agar library yang
ternyata tidak dipakai ikut terbuang.

## 4. Library FE dengan celah keamanan (`npm audit`) — sudah dianalisis, menunggu keputusan

Analisis 15 Nov 2026 (belum ada yang diubah). Portal 24 temuan, Toko 27 (1 critical), Coffee 27 (1 critical).

| Kelompok | Library | Catatan |
|---|---|---|
| A. aman (versi kecil) | axios 1.16→1.21.1, zod 4.4.3→4.6.7, react-router-dom 7.15→7.18.5, vite, postcss, eslint, library turunan | `npm audit fix` TANPA `--force` |
| B. critical | jspdf 2.5.2→4.2.1 (versi besar; Toko, Coffee; simpan struk PDF) | Hanya pakai `text/setFont/setFontSize/splitTextToSize/save/output`; risiko nyata rendah, tetap uji PDF struk |
| C. tanpa perbaikan di npm | xlsx 0.18.5 (Toko, Coffee; hanya MENULIS export produk) | Celah hanya saat membaca file; import produk dibaca BE (excelize). Pilihan: biarkan (saran AI) / SheetJS 0.20.3 dari cdn.sheetjs.com / ganti exceljs |
| D. alat lint | typescript-eslint | Jangan ikuti saran audit menurunkan ke 6.15; naikkan ke 8.72.2 |

Usulan urutan: 10a (kelompok A + typescript-eslint, ketiga FE) → 10b (jspdf) → 10c (xlsx sesuai pilihan).
**Keputusan user 17 Nov 2026 (mengikuti saran AI):** urutan 10a → 10b → 10c; xlsx **dibiarkan** (keputusan 131).
**10a [SELESAI] 17 Nov 2026:** `npm audit fix` (tanpa --force) di Portal, Toko, Coffee. Hanya `package-lock.json` yang berubah
(`package.json` tetap; server memakai `npm ci`). Library langsung naik versi kecil: axios 1.16.0→1.21.1, zod 4.4.3→4.6.7,
react-router-dom 7.15.0→7.18.5, vite 8.0.12→8.3.8, postcss 8.5.14→8.5.33, eslint 10.3.0→10.15.0, typescript-eslint 8.59.3→8.73.0
(bukan turun ke 6.15 seperti saran audit lama). Lompatan versi besar hanya di paket turunan alat dev (file-entry-cache, flat-cache,
keyv, minimatch milik eslint; @emnapi/* hanya untuk platform WebAssembly), tidak ikut ke hasil build. Hasil audit: Portal 24→0,
Toko 27→3, Coffee 27→3 (sisa: jspdf + dompurify = 10b, xlsx = 10c dibiarkan). Catatan: Vite 8.3 memberi peringatan `__dirname` di
`vite.config.ts` (akan tidak didukung di Vite versi besar berikutnya); bukan error, belum diubah.
Uji: type-check/lint/build ketiga FE; server dev dijalankan ulang dengan library baru; Portal `v_portal` 40/40 di laptop/tablet/HP,
login 9/9, profil/validasi 18/18; Toko `v_toko` 32/32 di laptop/tablet/HP, `b_toko_k3` 8/8; Coffee crawl baru
`testing/deadcode/c_coffee_crawl.js` 37 halaman tanpa error/API gagal/error JS, login 5/5, kelompok2 7/7, e2e `b_4_11` 22/22 di laptop
dan tablet, HP: 22/22. **Deploy:** FE ketiga aplikasi (`npm ci && npm run build`, SOP_DEPLOY), BE tidak berubah.
**10b [SELESAI] 17 Nov 2026:** jspdf 2.5.2 → 4.2.1 di Toko dan Coffee (`--save-exact`, versi tetap dipin seperti sebelumnya).
Turunan: dompurify 2.5.9 → 3.4.18 (tidak dipakai langsung oleh kode aplikasi), baru pako/fast-png/iobuffer, atob/btoa hilang.
Kode aplikasi hanya memakai `new jsPDF`, `setFont`, `setFontSize`, `text`, `splitTextToSize`, `save` (`receiptPdf.ts`); tidak ada
perubahan kode. Celah Toko dan Coffee 3 → 1 (sisa xlsx, dibiarkan sesuai keputusan 131).
Uji: `testing/deadcode/jspdf_snapshot.js` membuat PDF contoh dengan panggilan yang sama seperti struk (courier normal/tebal, tengah,
kanan, teks dibungkus, kertas 58 mm dan 80 mm) di versi lama dan baru: ukuran halaman, jumlah halaman, pembungkusan teks, dan 30
operasi teks (isi, font, ukuran, posisi) **identik**; type-check/lint/build; server dev dijalankan ulang; Simpan / Bagikan PDF dari
cetak ulang transaksi: Toko `b_toko_k3` 9/9 dan Coffee `b_kelompok2` 8/8 di laptop/tablet/HP, file yang diunduh diperiksa isinya
(header %PDF, %%EOF, kode transaksi, TOTAL — pemeriksaan isi baru ditambahkan); regresi `v_toko` 32/32 dan crawl Coffee 4/4.
**10c:** xlsx dibiarkan (keputusan 131), tanpa perubahan. **Nomor 10 SELESAI.** **Deploy:** FE Toko dan Coffee (bersama 10a).
Tiap tahap: type-check, lint, build, regresi browser, lalu deploy FE terkait. Perlu keputusan: setuju urutan? xlsx pilih mana?

## 6. Hapus dead code — **Coffee SELESAI**; **Portal dan BE User SELESAI**; **Toko SELESAI** (17 Nov 2026; log nginx produksi dicek: hanya 1 panggilan pemindai). Deploy BE + FE Toko oleh user. 6 tipe sisa Coffee sudah dihapus (17 Nov 2026). **Dead code semua aplikasi SELESAI.** Berikutnya nomor 10 (npm audit, bagian 4). Temuan baru: Coffee masih punya 6 tipe sisa (lihat DEADCODE.md Toko kelompok 1); lalu BE User → Toko (rincian `docs/DEADCODE.md`)

Temuan Coffee dari dead code **[SELESAI 17 Nov 2026]**: metode bayar di Ringkasan Bisnis (keputusan 132) dan tombol Tutup Kas
Hari Sebelumnya untuk admin (keputusan 133), riwayat #99. Deploy: FE Coffee saja.

Temuan baru saat mengerjakan 132-133 **[SELESAI 18 Nov 2026]** (keputusan 134 "Hari Ini" dan 135 "Hutang", riwayat #100).
Deploy: BE + FE Toko dan Coffee (BE menjalankan migrasi label otomatis saat start; backup DB dulu sesuai SOP).
- Ringkasan Bisnis Coffee **dan Toko**: pilihan periode "Hari Ini" hanya berlaku untuk kartu angka di atas; kotak di bawah (menu/produk
  terlaris, kategori, jenis pesanan, metode bayar) memakai **7 hari terakhir** karena `periodToDateRange` hanya membedakan "bulan" dan
  "selain bulan". Saran: "Hari Ini" = hari ini saja untuk kotak-kotak itu (perubahan kecil di FE, Toko produksi diuji penuh).
- Struk Bluetooth (`blePrinter.ts`, Toko dan Coffee) menulis metode `kredit` sebagai "Kredit", sedangkan struk layar dan laporan
  menulis "Hutang". Saran: samakan menjadi "Hutang" (mengubah teks struk cetak).

Alat: Go `deadcode`/`staticcheck`, FE `knip` (file, export, library npm tak terpakai). Analisis dulu, laporkan daftar,
hapus setelah user setuju. Per aplikasi, satu tahap sekali jalan: Coffee → Portal → BE User → Toko (paling akhir;
Toko produksi dan aturan "aditif", butuh izin eksplisit, hanya yang pasti tidak dipakai).

---

## Sudah selesai

- 15 Nov 2026: nomor 11 uji Go hijau di keempat BE (lihat riwayat dokumen rancangan). Uji DB dijalankan dengan
  `POS_DB_TEST=1 go test ./domain/product/repo/` (Toko, Coffee) dan `./domain/queue_number/repo/` (Coffee).

- 15 Nov 2026: cutover VPS (D.4), deploy FE Toko/Coffee (perbaikan logout dan peringatan data kosong),
  syarat Pembelian/Daftar Menu Coffee, akses Coffee untuk owner dan kasir di produksi, hapus data gladi di laptop.
