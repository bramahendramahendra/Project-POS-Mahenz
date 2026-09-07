# Panduan Migrasi Stok (Skema Lama → Skema Baru)

Panduan memindahkan data lama (stok di `products.stock`) ke skema baru (stok per satuan di `product_packages.stock`), **tanpa kehilangan data stok**.

- **Untuk coba manual di LOCAL** → ikuti [Bagian A](#bagian-a--langkah-di-local).
- **Untuk deploy di PRODUCTION** → ikuti [Bagian B](#bagian-b--langkah-di-production).

---

## Ringkasan Singkat

Alur intinya:

1. **Backup** database dulu (pengaman).
2. **Jalankan Backend sekali** → migrasi otomatis jalan. Stok lama otomatis disalin ke tabel `products_stock_backup` sebelum dihapus.
3. **Jalankan 3 skrip backfill** (urut): `backfill_purchase_package_id` → `backfill_missing_purchase_in` → `backfill_stock_restore`. Yang tengah menambal mutasi `in` pembelian yang hilang dari ledger; kalau dilewati, banyak produk salah ditandai perlu ditinjau.
4. **Cek hasil** → produk yang stoknya tidak yakin ditandai untuk ditinjau manual.

> **Yang penting Anda tahu:** stok baru **dihitung ulang dari riwayat** (pembelian + retur + penjualan), bukan disalin buta. Stok lama hanya dipakai untuk **validasi silang**. Kalau hasil hitung tidak cocok, produk ditandai `needs_stock_review` (tidak diisi angka asal).

---

## Bagian A — Langkah di LOCAL

Prasyarat: MySQL local aktif (WAMP), database `pos_retail_db` sudah berisi data lama (hasil restore dump prod).

Path tool MySQL WAMP: `C:\wamp64\bin\mysql\mysql8.4.7\bin\`

### Langkah 1 — Matikan Backend

Pastikan Backend (BE) **tidak sedang jalan**. Ini penting: kalau BE jalan, migrasi keburu jalan sebelum kita backup.

### Langkah 2 — Backup Database (pengaman)

```bash
sudo mysqldump -u root pos_retail_db > backups/backup_sebelum_migrasi.sql
```

> Kalau ada yang salah, Anda bisa balik ke kondisi ini (lihat [Rollback](#rollback)).

### Langkah 3 — Cek Stok Lama (baseline)

Catat angka ini untuk dibandingkan nanti:

```bash
sudo mysql -u root pos_retail_db -e "SELECT COUNT(*) AS produk, ROUND(SUM(stock),3) AS total_stok FROM products;"
```

> Kalau `total_stok` = 0, berarti data belum ke-restore dengan benar — **STOP**, jangan lanjut.

### Langkah 4 — Jalankan Backend Sekali (migrasi otomatis)

```bash
cd BE
go run main.go
```

Tunggu sampai muncul log "Server running on port :8080" (tanda migrasi selesai). Lalu **matikan** BE (Ctrl+C).

Apa yang terjadi otomatis di balik layar:
- Migrasi 003–009 dijalankan.
- **Migrasi 005 otomatis menyalin stok lama ke `products_stock_backup` sebelum menghapus `products.stock`.** (Anda tidak perlu buat tabel ini manual.)
- `product_packages.stock` terbentuk (masih 0, akan diisi backfill).

Verifikasi backup stok otomatis berhasil:

```bash
mysql -u root pos_retail_db -e "SELECT COUNT(*) AS baris, ROUND(SUM(stock),3) AS total FROM products_stock_backup;"
```

> `total` harus sama dengan `total_stok` di Langkah 3.

### Langkah 5 — Backfill package_id

```bash
cd BE
go run ./cmd/backfill_purchase_package_id/
```

Hasil yang diharapkan: semua terisi, 0 dilewati.

### Langkah 5b — Tambal mutasi `in` pembelian yang hilang

Sebagian item pembelian (`purchase_items`) di data prod ternyata TIDAK punya baris
mutasi `in` pasangannya di ledger `stock_mutations` — barang yang benar-benar
dibeli, tapi jejak stok masuknya tidak pernah tercatat. Kalau dibiarkan, backfill
stok (langkah 6) akan menghitung barang itu seolah tidak pernah masuk → stok minus
→ produk salah ditandai `needs_stock_review`. Skrip ini menambal baris `in` yang
hilang itu dari data `purchase_items` yang memang ada.

```bash
cd BE
go run ./cmd/backfill_missing_purchase_in/
```

Hasil yang diharapkan (contoh): `Berhasil disisipkan 'in' : 23` (angka bisa beda).
Yang `[SKIP]` biasanya produk rantai satuan bercabang (celah #10) — wajar, ditangani
manual belakangan.

> Aman & idempotent: skrip HANYA menyentuh item yang benar-benar belum punya baris
> `in` (cek `NOT EXISTS`), jadi tidak ada dobel-hitung dan boleh dijalankan berulang.
> Baris `in` yang disisipkan sengaja di-backdate `created_at`-nya ke tanggal PO,
> supaya saat langkah 6 me-replay riwayat, stok masuk itu dihitung SEBELUM penjualan.
> Bisa dicoba dulu tanpa menulis apa pun: tambahkan flag `--dry-run`.

### Langkah 6 — Backfill Stok (hitung ulang dari riwayat)

```bash
cd BE
go run ./cmd/backfill_stock_restore/
```

Hasil yang diharapkan (contoh, setelah langkah 5b dijalankan):
```
Berhasil direkonstruksi & disimpan : 238
Ditandai perlu ditinjau (needs_stock_review): 12
```

> Kalau langkah 5b DILEWATI, angka `needs_stock_review` akan jauh lebih besar (mis. 32)
> karena banyak pembelian yang mutasi `in`-nya belum ditambal.

### Langkah 7 — Verifikasi

```bash
mysql -u root pos_retail_db -e "SELECT (SELECT COUNT(*) FROM product_packages WHERE stock>0) AS paket_ada_stok, (SELECT COUNT(*) FROM products WHERE needs_stock_review=1) AS perlu_ditinjau;"
```

### Langkah 8 — Jalankan BE & FE normal

```bash
# Terminal 1
cd BE
go run main.go

# Terminal 2
cd FE
node node_modules/vite/bin/vite.js --port 3000
```

Login ulang di browser. Produk yang `needs_stock_review` bisa ditangani lewat menu **Rekonsiliasi Stok** (grup Pelaporan, khusus admin).

---

## Bagian B — Langkah di PRODUCTION

Sama seperti local, dengan **2 perbedaan**: cara menjalankan BE (build + systemd) dan **wajib set `MIGRATION_DSN`** sebelum backfill (karena DB prod pakai user/password, bukan root tanpa password).

### 1. Backup dulu

```bash
mysqldump -u pos_user -p pos_retail_db > backup_sebelum_migrasi.sql
```

### 2. Deploy kode + restart BE (migrasi otomatis jalan)

```bash
git pull
go build -o pos_api .
systemctl restart pos-backend
```

Saat BE restart, migrasi 003–009 jalan otomatis, dan migrasi 005 otomatis membuat `products_stock_backup` sebelum drop.

### 3. Set koneksi DB untuk backfill

Skrip backfill defaultnya connect ke `root@127.0.0.1` tanpa password (untuk local). Di production, set `MIGRATION_DSN` dulu sesuai `config_prod.json`:

```bash
export MIGRATION_DSN='pos_user:P@ssw0rd@tcp(127.0.0.1:3306)/pos_retail_db?charset=utf8&parseTime=True&loc=Local'
```

> Format: `USER:PASSWORD@tcp(HOST:PORT)/pos_retail_db?charset=utf8&parseTime=True&loc=Local`

### 4. Jalankan backfill

```bash
cd BE
go run ./cmd/backfill_purchase_package_id/
go run ./cmd/backfill_missing_purchase_in/
go run ./cmd/backfill_stock_restore/
```

> `backfill_missing_purchase_in` WAJIB dijalankan SEBELUM `backfill_stock_restore`
> (menambal mutasi `in` pembelian yang hilang dari ledger). Kalau dilewati, banyak
> produk akan salah ditandai `needs_stock_review`. Skrip ini idempotent & baca
> `MIGRATION_DSN` yang sama.

### 5. Verifikasi & tangani produk yang perlu ditinjau

```bash
mysql -u pos_user -p pos_retail_db -e "SELECT (SELECT COUNT(*) FROM product_packages WHERE stock>0) AS paket_ada_stok, (SELECT COUNT(*) FROM products WHERE needs_stock_review=1) AS perlu_ditinjau;"
```

Produk `needs_stock_review` ditangani lewat menu **Rekonsiliasi Stok** di aplikasi.

---

## Menangani Produk `needs_stock_review`

Setelah langkah 5b (`backfill_missing_purchase_in`) dijalankan, penyebab paling
umum — pembelian yang mutasi `in`-nya hilang dari ledger — sudah ditambal otomatis.
Sisa produk yang masih ditandai biasanya karena hal yang TIDAK bisa ditambal
otomatis:

- **Rantai satuan bercabang** — struktur satuan produk tidak linear (>1 anak menunjuk
  satu paket induk). Perlu dirapikan strukturnya dulu, tidak bisa direkonstruksi otomatis.
- **Paket satuan tidak ditemukan** — ada penjualan lama yang memakai `package_id` yang
  kini tidak ada lagi pada produk itu (satuan sudah diubah/dihapus sejak transaksi terjadi).
- **Selisih di luar wajar** — hasil hitung ulang masih beda jauh dari stok lama, mis. ada
  penjualan yang benar-benar melebihi total pembelian tercatat (bukan sekadar pembulatan).

Cara menangani sisa ini: buka menu **Rekonsiliasi Stok** di aplikasi (login sebagai admin).
Di sana Anda bisa lihat stok lama vs baru, rincian perhitungan, dan koreksi manual per
satuan. Setelah dikoreksi, flag hilang otomatis.

---

## Rollback

Kalau hasil migrasi tidak sesuai, balik ke backup:

```bash
mysql -u root -e "DROP DATABASE pos_retail_db; CREATE DATABASE pos_retail_db;"
mysql -u root pos_retail_db < backups/backup_sebelum_migrasi.sql
```

Lalu ulangi dari awal.

---

## Prompt Siap-Pakai (kalau minta bantuan Kiro)

Salin-tempel ini saat ingin Kiro menjalankan migrasi ulang untuk Anda:

```
Restore prod ke pos_retail_db sudah selesai. Jalankan migrasi stok sesuai
docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md Bagian A (local), langkah 1-8 (termasuk
langkah 5b: backfill_missing_purchase_in SEBELUM backfill_stock_restore).
Kabari hasil tiap langkah. Ingat: backup dulu, stok dihitung ulang dari riwayat,
jalankan ketiga skrip backfill sesuai urutan.
```

---

## Catatan Teknis

| Hal | Keterangan |
|-----|-----------|
| Migrasi jalan otomatis | Saat BE start (`database/migrate.go`), file `001`–`009` di `BE/database/migrations/` |
| Backup stok otomatis | Migrasi `005` membuat `products_stock_backup` sebelum drop `products.stock` |
| Skrip backfill (urut) | `backfill_purchase_package_id` → `backfill_missing_purchase_in` → `backfill_stock_restore` (semua di `BE/cmd/`) |
| Skrip tambal `in` hilang | `BE/cmd/backfill_missing_purchase_in` — sisip mutasi `in` untuk `purchase_items` PO active yang bolong ledger; idempotent (`NOT EXISTS`), `created_at` di-backdate ke tanggal PO. Punya flag `--dry-run` |
| Koneksi DB backfill | Env `MIGRATION_DSN` (kalau kosong → default `root@127.0.0.1` untuk local) |
| Tabel backup stok | `products_stock_backup` (kolom: `id`, `stock`, `reserved_qty`) |
| Flag tinjau | Kolom `products.needs_stock_review` + `products.stock_review_note` |
| Urutan replay backfill | `backfill_stock_restore` me-replay `stock_mutations` urut `created_at, id` (bukan murni `id`) — supaya baris `in` hasil tambal yang di-backdate diproses SEBELUM penjualan lama |
