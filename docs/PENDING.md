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
Tiap tahap: type-check, lint, build, regresi browser, lalu deploy FE terkait. Perlu keputusan: setuju urutan? xlsx pilih mana?

## 6. Hapus dead code — **Coffee SELESAI** (16 Nov 2026, rincian `docs/DEADCODE.md`); berikutnya Portal → BE User → Toko

Temuan Coffee yang perlu keputusan user (tidak dihapus, tidak dibuatkan fitur):
- Ringkasan Bisnis belum menampilkan **metode bayar**, padahal keputusan 117 menulis "metode bayar tetap" (endpoint
  `/reports/business-summary/payment-methods` sudah ada di BE; metode bayar saat ini hanya di Laporan Penjualan).
- Endpoint `/cash-drawer/auto-close-yesterday` (keputusan 4.9) belum punya tombol di layar.

Alat: Go `deadcode`/`staticcheck`, FE `knip` (file, export, library npm tak terpakai). Analisis dulu, laporkan daftar,
hapus setelah user setuju. Per aplikasi, satu tahap sekali jalan: Coffee → Portal → BE User → Toko (paling akhir;
Toko produksi dan aturan "aditif", butuh izin eksplisit, hanya yang pasti tidak dipakai).

---

## Sudah selesai

- 15 Nov 2026: nomor 11 uji Go hijau di keempat BE (lihat riwayat dokumen rancangan). Uji DB dijalankan dengan
  `POS_DB_TEST=1 go test ./domain/product/repo/` (Toko, Coffee) dan `./domain/queue_number/repo/` (Coffee).

- 15 Nov 2026: cutover VPS (D.4), deploy FE Toko/Coffee (perbaikan logout dan peringatan data kosong),
  syarat Pembelian/Daftar Menu Coffee, akses Coffee untuk owner dan kasir di produksi, hapus data gladi di laptop.
