# SOP Deploy — Perbaikan Menu Laba Rugi

SOP ini KHUSUS untuk merilis paket perbaikan **menu Laba Rugi** ke server produksi.
Ditulis untuk orang yang tidak biasa koding — ikuti langkah dari atas ke bawah, salin-tempel
perintahnya ke terminal server, jangan lompat-lompat.

> Panduan deploy umum (update rutin biasa) ada di `PANDUAN_DEPLOY.md`. SOP ini adalah
> **tambahan** untuk rilis Laba Rugi karena ada satu langkah istimewa yang tidak ada di
> update biasa: **koreksi data modal lama**. Langkah itu WAJIB dilakukan, kalau tidak,
> laporan Laba Rugi di produksi tetap salah (tampak rugi padahal untung).

---

## Ringkasan: Apa yang berubah di rilis ini?

Dengan bahasa sederhana, rilis Laba Rugi berisi 5 hal:

1. **Data modal lama dibetulkan.** Dulu sebagian transaksi mencatat modal barang dengan
   satuan yang salah, sehingga laporan Laba Rugi tampak rugi padahal sebenarnya untung.
   Rilis ini membetulkan data lama itu.
2. **Menu baru "Rekonsiliasi Modal"** (khusus admin) untuk membetulkan manual segelintir
   data yang tidak bisa dibetulkan otomatis.
3. **Perbaikan transaksi mundur tanggal (backdate)** supaya mencatat modal dengan benar.
4. **Tampilan Laba Rugi baru** yang lebih mudah dibaca (kartu kesimpulan, pembanding
   periode, daftar produk untung/rugi).
5. **Filter tanggal baru** dengan keterangan periode aktif dan tombol pilihan cepat.

Poin 2–5 ikut otomatis saat backend & frontend di-update biasa. Yang butuh perhatian
ekstra adalah **poin 1** (langkah 4 di bawah).

---

## Isi SOP

- [Sebelum Mulai — Cek Dulu](#sebelum-mulai--cek-dulu)
- [Langkah 1 — Aktifkan Maintenance](#langkah-1--aktifkan-maintenance)
- [Langkah 2 — Backup Database (WAJIB)](#langkah-2--backup-database-wajib)
- [Langkah 3 — Update Backend & Frontend](#langkah-3--update-backend--frontend)
- [Langkah 4 — Koreksi Data Modal Lama (langkah istimewa)](#langkah-4--koreksi-data-modal-lama-langkah-istimewa)
- [Langkah 5 — Cek Hasil di Aplikasi](#langkah-5--cek-hasil-di-aplikasi)
- [Langkah 6 — Tangani Sisa Data Manual (Rekonsiliasi Modal)](#langkah-6--tangani-sisa-data-manual-rekonsiliasi-modal)
- [Langkah 7 — Matikan Maintenance](#langkah-7--matikan-maintenance)
- [Kalau Gagal — Cara Balik ke Semula (Rollback)](#kalau-gagal--cara-balik-ke-semula-rollback)

---

## Sebelum Mulai — Cek Dulu

- Pastikan Anda bisa masuk ke server (SSH) dan punya akses `sudo`.
- Pastikan tahu password database `pos_user`.
- Lakukan rilis ini di **jam sepi** (misalnya malam / sebelum toko buka), karena aplikasi
  akan di-maintenance sebentar.
- Sisihkan waktu tenang ± 30–45 menit. Jangan buru-buru.

---

## Langkah 1 — Aktifkan Maintenance

Supaya tidak ada kasir yang bertransaksi saat data sedang dibetulkan.

```bash
sudo maintenance-on.sh 139.180.214.187
```

---

## Langkah 2 — Backup Database (WAJIB)

**Jangan lewati langkah ini.** Ini jaring pengaman kalau ada yang salah. Simpan salinan
seluruh data ke satu file bertanggal:

```bash
mkdir -p /opt/pos-mahenz/BE/backups
mysqldump -u pos_user -p pos_retail_db > /opt/pos-mahenz/BE/backups/backup_sebelum_laba_rugi_$(date +%Y%m%d_%H%M%S).sql
```

Pastikan file backup benar-benar terisi (ukuran tidak 0 byte):

```bash
ls -lh /opt/pos-mahenz/BE/backups/
```

> Kalau ukuran file 0 atau perintah error, **STOP**. Jangan lanjut sebelum backup berhasil.

---

## Langkah 3 — Update Backend & Frontend

Ini mengambil kode terbaru dan menyalakan menu/tampilan baru. Saat backend restart,
struktur database baru (tabel `transaction_cost_review` + menu Rekonsiliasi Modal)
**dibuat otomatis** — tidak perlu langkah manual.

**Backend (mesin):**

```bash
cd /opt/pos-mahenz/BE
git pull
go build -o pos_api main.go
sudo systemctl restart pos-backend
sudo systemctl status pos-backend    # harus "active (running)"
```

**Frontend (tampilan):**

```bash
cd /opt/pos-mahenz/FE
git pull
npm install
npm run type-check
npm run build
sudo rm -rf /var/www/pos-web/dist
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

---

## Langkah 4 — Koreksi Data Modal Lama (langkah istimewa)

Ini inti dari rilis Laba Rugi. Sebuah skrip khusus akan menelusuri semua transaksi lama
dan membetulkan modal yang salah satuan. Skrip ini **hanya dijalankan sekali** saat rilis.

Skrip punya dua mode:
- **Mode uji (dry-run):** cuma menghitung & melapor, TIDAK mengubah data. Aman.
- **Mode terapkan (`--apply`):** benar-benar membetulkan data.

### 4.1 Jalankan mode UJI dulu (tidak mengubah apa pun)

```bash
cd /opt/pos-mahenz/BE
go run ./cmd/backfill_transaction_cogs
```

Baca laporan yang muncul. Perhatikan angka ini:
- **Akan dikoreksi otomatis** — jumlah baris yang akan dibetulkan.
- **Perlu tinjauan MANUAL** — jumlah baris yang tidak bisa otomatis (nanti dibereskan
  lewat menu di Langkah 6). Wajar kalau ada beberapa.
- **Total COGS SEBELUM vs SESUDAH** — angka modal total sebelum & sesudah dikoreksi.

> Kalau laporan terlihat wajar (banyak baris dikoreksi, sisa manual sedikit), lanjut ke 4.2.
> Kalau muncul error atau angkanya aneh sekali, **STOP** dan hubungi pengembang. Belum ada
> data yang berubah, jadi aman.

### 4.2 Jalankan mode TERAPKAN

> Pastikan **Langkah 2 (backup) sudah berhasil** sebelum menjalankan ini.

```bash
cd /opt/pos-mahenz/BE
go run ./cmd/backfill_transaction_cogs --apply
```

Di akhir akan muncul, misalnya: `APPLY selesai. Baris diperbarui: <angka>`. Setiap
perubahan otomatis dicatat ke tabel audit `transaction_items_cogs_audit` (nilai lama
disimpan), jadi bisa ditelusuri bila perlu.

---

## Langkah 5 — Cek Hasil di Aplikasi

1. Sementara matikan maintenance dulu atau buka lewat akses admin, lalu login sebagai admin.
2. Buka menu **Pelaporan → Laba Rugi**.
3. Pilih periode (misalnya "Bulan ini"). Pastikan:
   - Angka **Laba Bersih** masuk akal (tidak lagi rugi besar yang janggal).
   - Keterangan periode aktif muncul ("Menampilkan: tanggal ... – ...").
   - Daftar produk untung/rugi tampil normal.
4. Buka menu **Pelaporan → Rekonsiliasi Modal** (khusus admin). Di sini muncul daftar
   baris yang perlu dibetulkan manual (dari "Perlu tinjauan MANUAL" tadi).

---

## Langkah 6 — Tangani Sisa Data Manual (Rekonsiliasi Modal)

Sebagian kecil transaksi lama tidak bisa dibetulkan otomatis (misalnya master paket
barangnya tidak lengkap). Baris-baris ini masuk ke menu **Rekonsiliasi Modal**.

- Ini **tidak mendesak** — aplikasi tetap jalan normal. Bisa dikerjakan bertahap.
- Untuk tiap baris: admin melihat harga jual & modal yang tercatat, lalu memasukkan modal
  yang benar (atau "lewati" bila memang tidak bisa dipastikan).
- Setelah dibereskan, angka Laba Rugi untuk transaksi tersebut ikut akurat.

> Boleh dikosongkan dulu dan dibereskan lain hari. Tidak menghalangi rilis.

---

## Langkah 7 — Matikan Maintenance

Kalau semua sudah dicek dan beres:

```bash
sudo maintenance-off.sh 139.180.214.187
```

Lalu buka aplikasi seperti biasa, coba login & buka beberapa halaman untuk memastikan
semua normal. Refresh keras browser (Ctrl+Shift+R) bila tampilan masih versi lama.

---

## Kalau Gagal — Cara Balik ke Semula (Rollback)

Kalau setelah koreksi data ternyata ada yang tidak beres dan ingin kembali ke kondisi
sebelum rilis:

1. Pastikan maintenance masih aktif (Langkah 1).
2. Kembalikan database dari file backup yang dibuat di Langkah 2:

   ```bash
   mysql -u pos_user -p pos_retail_db < /opt/pos-mahenz/BE/backups/NAMA_FILE_BACKUP_ANDA.sql
   ```

   (Ganti `NAMA_FILE_BACKUP_ANDA.sql` dengan nama file backup yang tadi dibuat.)

3. Restart backend:

   ```bash
   sudo systemctl restart pos-backend
   ```

4. Cek aplikasi kembali seperti sebelum rilis.

> Karena data sebelum koreksi sudah di-backup di Langkah 2 dan tiap perubahan tercatat di
> tabel audit, rollback aman dilakukan. Kalau ragu, hubungi pengembang sambil menjaga file
> backup jangan terhapus.

---

> Checklist singkat sebelum bilang "beres":
> 1. Backup database berhasil (file terisi, bukan 0 byte)
> 2. Backend `active (running)`, frontend ter-build & tersalin
> 3. Skrip koreksi mode UJI wajar → mode TERAPKAN sukses
> 4. Menu Laba Rugi angkanya masuk akal & menu Rekonsiliasi Modal muncul
> 5. Maintenance dimatikan, aplikasi bisa login & dipakai normal
