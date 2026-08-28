# Panduan Migrasi Stok: Data Prod (Skema Lama) → Skema Baru

> Panduan langkah demi langkah untuk memindahkan data produksi (yang masih pakai skema lama) ke skema baru, **tanpa kehilangan data stok**.

---

## 1. Latar Belakang (Kenapa Perlu Migrasi Ini)

### Konsep stok berubah
Dulu, stok disimpan sebagai **1 angka desimal** di kolom `products.stock`. Masalahnya: karena desimal, ada stok yang **sebenarnya masih ada tapi dianggap habis** akibat pembulatan/pemotongan presisi.

Konsep baru: stok disimpan **per level satuan** di `product_packages.stock` (misal: stok Pack sendiri, stok Slop sendiri). Lebih akurat, tidak ada lagi masalah presisi.

### Masalah saat restore data prod
Database prod masih skema **lama** (hanya migrasi `001` & `002`). Migrasi `003`–`008` belum pernah jalan di prod.

Saat data prod di-restore ke local lalu Backend (BE) dijalankan, BE **otomatis menjalankan migrasi** `003`–`008`. Dua migrasi ini berbahaya untuk data stok:

- **Migrasi `004`**: menambah kolom `product_packages.stock` dengan nilai **default 0**
- **Migrasi `005`**: **menghapus (DROP)** kolom `products.stock`

Artinya: kalau BE langsung dijalankan setelah restore, **angka stok lama hilang** (kolomnya di-drop) padahal `product_packages.stock` masih 0 semua → **semua stok jadi kosong**.

### Solusi
**Selamatkan angka stok lama ke tabel cadangan SEBELUM BE dijalankan.** Lalu setelah skema baru siap, hitung ulang stok dari riwayat transaksi (pembelian, retur, penjualan) dan validasi silang dengan angka stok lama tadi.

---

## 2. Prinsip Utama

| Prinsip | Penjelasan |
|---------|-----------|
| **Backup dulu** | Selalu backup DB sebelum mengubah apapun |
| **Jangan jalankan BE sebelum tabel cadangan dibuat** | Kalau BE jalan duluan, `products.stock` keburu ke-drop |
| **Stok dihitung ulang dari RIWAYAT, bukan disalin buta** | Sumber kebenaran = pembelian + retur + penjualan. `products.stock` lama hanya untuk **validasi silang** |
| **Kalau ragu, tandai untuk ditinjau** | Produk yang hasil hitungnya tidak cocok/aneh ditandai `needs_stock_review`, bukan diisi angka asal |

---

## 3. Alur Migrasi (Diagram Sederhana)

```
[1] Restore data prod (skema lama) ke pos_retail_db
        │
        ▼
[2] Backup DB local (pengaman)
        │
        ▼
[3] Buat tabel cadangan stok lama:
    products_stock_backup  ← menyalin products.stock
        │  (WAJIB sebelum BE jalan!)
        ▼
[4] Jalankan BE sekali
    → migrasi 003-008 jalan
    → products.stock ter-drop (tapi sudah aman di backup)
    → product_packages.stock terbentuk (masih 0)
        │
        ▼
[5] Jalankan backfill_purchase_package_id
    → isi purchase_items.package_id (prasyarat)
        │
        ▼
[6] Jalankan backfill_stock_restore
    → replay riwayat → hitung stok per level
    → validasi silang vs products_stock_backup
        │
        ▼
[7] Verifikasi hasil + cek produk needs_stock_review
```

---

## 4. Langkah Detail

### Prasyarat
- DB local MySQL aktif (WAMP: `C:\wamp64\bin\mysql\mysql8.4.7\bin\`)
- Database bernama `pos_retail_db`
- Punya file dump data prod (skema lama)

---

### Langkah 1 — Matikan Backend

Pastikan BE **tidak berjalan**. Kalau pakai systemd/terminal, hentikan dulu. Ini penting agar tidak ada auto-migrasi yang jalan sebelum kita siap.

---

### Langkah 2 — Restore Data Prod ke `pos_retail_db`

Restore dump prod (skema lama) ke database `pos_retail_db`. Contoh via command line:

```bash
mysql -u root pos_retail_db < path/ke/dump_prod.sql
```

> Setelah langkah ini, `pos_retail_db` berisi skema lama: `products.stock` masih ada dan berisi angka stok asli prod.

---

### Langkah 3 — Backup DB Local (Pengaman)

Sebelum apapun diubah, backup dulu:

```bash
mysqldump -u root pos_retail_db > backups/pos_retail_db_sebelum_migrasi.sql
```

> Kalau ada yang salah di tengah jalan, Anda bisa restore balik dari file ini.

---

### Langkah 4 — Buat Tabel Cadangan Stok Lama (WAJIB, SEBELUM BE JALAN)

Ini kunci utamanya. Salin angka stok lama ke tabel terpisah yang **tidak akan disentuh** migrasi:

```sql
CREATE TABLE products_stock_backup AS
SELECT id, stock, reserved_qty FROM products;
```

Jalankan via mysql client:

```bash
mysql -u root pos_retail_db -e "CREATE TABLE products_stock_backup AS SELECT id, stock, reserved_qty FROM products;"
```

Verifikasi tabelnya terisi:

```bash
mysql -u root pos_retail_db -e "SELECT COUNT(*) AS total, SUM(stock) AS total_stok FROM products_stock_backup;"
```

> Harus muncul jumlah produk dan total stok > 0. Kalau `total_stok` = 0, berarti data prod belum ke-restore dengan benar — jangan lanjut.

---

### Langkah 5 — Jalankan Backend Sekali (Biar Migrasi Jalan)

```bash
cd BE
go run main.go
```

Tunggu sampai log menampilkan semua route terdaftar (migrasi selesai). Lalu **matikan lagi** BE (Ctrl+C).

Apa yang terjadi di balik layar:
- Migrasi `003`–`008` dijalankan otomatis
- `products.stock` di-drop → **tapi angkanya sudah aman di `products_stock_backup`**
- `product_packages.stock` terbentuk (nilai 0)

---

### Langkah 6 — Backfill `purchase_items.package_id`

Prasyarat wajib sebelum hitung ulang stok. Skrip ini mengisi kolom `package_id` di data pembelian lama (dicocokkan dari nama satuan):

```bash
cd BE
go run ./cmd/backfill_purchase_package_id/
```

Hasil yang diharapkan:
```
Total baris purchase_items.package_id NULL: <N>
Berhasil diisi otomatis : <N>
Dilewati (unit tak ketemu): 0
Dilewati (ambigu)         : 0
```

> Kalau ada yang "dilewati", catat produk mana — nanti stoknya perlu dicek manual.

---

### Langkah 7 — Backfill Stok (Hitung Ulang dari Riwayat)

Skrip ini me-replay seluruh riwayat (pembelian, retur, penjualan, void, expired) untuk menghitung stok per level, lalu memvalidasi silang dengan `products_stock_backup`:

```bash
cd BE
go run ./cmd/backfill_stock_restore/
```

Hasil yang diharapkan:
```
Total produk: <N>
=== Laporan Backfill Restore (dari products_stock_backup) ===
Berhasil direkonstruksi & disimpan : <X>
Ditandai perlu ditinjau (needs_stock_review): <Y>
```

**Penjelasan hasil:**
- **Berhasil direkonstruksi**: stok dihitung dari riwayat DAN cocok dengan stok lama (selisih ≤ 0.01). Ini stok yang bisa dipercaya.
- **Ditandai perlu ditinjau**: hasil hitung tidak cocok dengan stok lama, atau ada anomali (stok minus saat replay, satuan bercabang, dll). Produk ini stoknya **tidak diubah** dan diberi flag agar muncul peringatan di aplikasi untuk dicek manual.

---

### Langkah 8 — Verifikasi

Cek berapa produk yang stoknya terisi dan berapa yang perlu ditinjau:

```bash
mysql -u root pos_retail_db -e "SELECT (SELECT COUNT(*) FROM product_packages WHERE stock > 0) AS paket_ada_stok, (SELECT COUNT(*) FROM products WHERE needs_stock_review = 1) AS perlu_ditinjau, (SELECT COUNT(*) FROM products) AS total_produk;"
```

Cek daftar produk yang perlu ditinjau beserta alasannya:

```bash
mysql -u root pos_retail_db -e "SELECT id, name, stock_review_note FROM products WHERE needs_stock_review = 1;"
```

---

### Langkah 9 — Jalankan BE & FE Normal

Setelah stok terisi, jalankan aplikasi seperti biasa:

```bash
# Backend
cd BE
go run main.go

# Frontend (terminal terpisah)
cd FE
node node_modules/vite/bin/vite.js --port 3000
```

Login ulang di browser (session lama dari DB prod sudah tidak valid). Stok sekarang seharusnya sudah muncul.

---

## 5. Menangani Produk `needs_stock_review`

Produk yang ditandai perlu ditinjau punya 2 kemungkinan penyebab:

### A. "Rantai satuan bercabang"
Produk punya lebih dari 1 satuan turunan yang menunjuk ke induk yang sama. Skrip sengaja tidak memproses otomatis karena bisa salah hitung. **Solusi:** cek struktur satuan produk di aplikasi, perbaiki bila perlu, lalu input stok manual atau jalankan ulang backfill.

### B. "Stok tidak mencukupi" / "paket tidak ditemukan"
Saat replay riwayat, di suatu titik stok jadi minus (misal ada penjualan tercatat sebelum pembelian), atau satuan yang dipakai di histori sudah berubah/dihapus. **Solusi:** cek riwayat produk tersebut secara manual, tentukan stok yang benar, input manual lewat aplikasi.

> Setelah stok diperbaiki manual, flag `needs_stock_review` akan hilang otomatis atau bisa di-clear lewat fitur "tandai sudah ditinjau" di aplikasi.

---

## 6. Kalau Terjadi Masalah (Rollback)

Kalau hasil migrasi tidak sesuai atau ada error, restore dari backup langkah 3:

```bash
# Drop database, buat ulang, restore backup
mysql -u root -e "DROP DATABASE pos_retail_db; CREATE DATABASE pos_retail_db;"
mysql -u root pos_retail_db < backups/pos_retail_db_sebelum_migrasi.sql
```

Lalu ulangi dari langkah yang sesuai.

---

## 7. Ringkasan Perintah (Cheat Sheet)

```bash
# 1. (matikan BE dulu)

# 2. Restore prod
mysql -u root pos_retail_db < path/ke/dump_prod.sql

# 3. Backup pengaman
mysqldump -u root pos_retail_db > backups/pos_retail_db_sebelum_migrasi.sql

# 4. Tabel cadangan stok (WAJIB sebelum BE jalan)
mysql -u root pos_retail_db -e "CREATE TABLE products_stock_backup AS SELECT id, stock, reserved_qty FROM products;"

# 5. Jalankan BE sekali (migrasi jalan), lalu matikan
cd BE && go run main.go   # tunggu selesai, Ctrl+C

# 6. Backfill package_id
go run ./cmd/backfill_purchase_package_id/

# 7. Backfill stok
go run ./cmd/backfill_stock_restore/

# 8. Verifikasi
mysql -u root pos_retail_db -e "SELECT (SELECT COUNT(*) FROM product_packages WHERE stock > 0) AS paket_ada_stok, (SELECT COUNT(*) FROM products WHERE needs_stock_review = 1) AS perlu_ditinjau;"

# 9. Jalankan BE & FE normal
```

---

## 8. File Terkait

| File | Fungsi |
|------|--------|
| `BE/cmd/backfill_purchase_package_id/main.go` | Isi `purchase_items.package_id` dari nama satuan (prasyarat) |
| `BE/cmd/backfill_stock_restore/main.go` | Hitung ulang stok dari riwayat + validasi vs `products_stock_backup` (untuk skenario restore prod) |
| `BE/cmd/backfill_stock_per_level/main.go` | Versi asli (untuk DB yang `products.stock`-nya masih ada, belum ter-drop) |
| `BE/database/migrations/004_stock_per_package_level.sql` | Migrasi: tambah `product_packages.stock` |
| `BE/database/migrations/005_drop_legacy_product_stock_columns.sql` | Migrasi: drop `products.stock` |

> **Catatan:** `backfill_stock_restore` dipakai kalau `products.stock` **sudah** ter-drop (skenario restore prod ke skema baru). `backfill_stock_per_level` dipakai kalau kolomnya **masih ada**. Keduanya logika intinya sama, beda cuma sumber baca stok lama.

---

## 9. Prompt Siap-Pakai untuk Migrasi Ulang (Restore DB Prod dari Awal)

> Salin-tempel prompt di bawah ini ke Kiro saat Anda sudah restore DB prod (skema lama) dari awal ke `pos_retail_db` dan ingin migrasi stok dijalankan ulang. Prompt ini sudah memuat semua konteks penting supaya tidak ada langkah yang keliru.

```
Saya baru restore DB prod (skema lama) dari awal ke pos_retail_db di local.
Tolong jalankan migrasi stok lengkap sesuai docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md,
dan kabari saya tiap tahap yang dijalankan beserta hasilnya.

Konteks penting (JANGAN dilanggar):
1. DB: pos_retail_db di 127.0.0.1:3306, user root TANPA password (WAMP).
   mysql client: C:\wamp64\bin\mysql\mysql8.4.7\bin\mysql.exe (dan mysqldump.exe).
2. Konsep stok: stok TIDAK PERNAH di-set manual, selalu naik/turun by sistem.
   Stok baru dihitung ULANG dari RIWAYAT (pembelian + retur + penjualan + void +
   expired), BUKAN disalin buta dari products.stock lama. products.stock lama
   HANYA untuk validasi silang (toleransi 0.01).
3. JANGAN ubah kode yang sudah benar. Skrip backfill (backfill_stock_restore &
   backfill_stock_per_level) dan logika stok live TIDAK BOLEH disentuh. Kalau ada
   kendala, buat query/skrip BARU terpisah — konfirmasi dulu ke saya sebelum
   mengubah kode manapun yang sudah jalan.
4. URUTAN WAJIB (kalau kebalik, stok lama hilang):
   a. Pastikan BE MATI dulu.
   b. Backup pengaman: mysqldump pos_retail_db -> backups/.
   c. Buat tabel cadangan SEBELUM BE jalan:
      CREATE TABLE products_stock_backup AS SELECT id, stock, reserved_qty FROM products;
      (verifikasi total_stok > 0; kalau 0 berarti restore gagal, STOP).
   d. Jalankan BE sekali (migrasi 003-009 jalan, products.stock ke-drop oleh 005,
      product_packages.stock terbentuk), lalu matikan BE.
   e. go run ./cmd/backfill_purchase_package_id/  (isi purchase_items.package_id).
   f. go run ./cmd/backfill_stock_restore/         (hitung ulang stok + validasi).
   g. Verifikasi: paket_ada_stok + perlu_ditinjau, lalu daftar produk
      needs_stock_review beserta alasannya.
   h. Jalankan BE & FE normal.
5. Catatan lingkungan:
   - PowerShell sering echo/buffer output — verifikasi lewat query terpisah atau
     tulis output ke file lalu baca, jangan andalkan output terminal langsung.
   - Migrasi dijalankan otomatis oleh BE saat start (database/migrate.go), termasuk
     009_stock_reconciliation_menu.sql (menu Rekonsiliasi Stok, admin-only).
   - Setelah selesai, 32 produk (atau berapa pun) yang needs_stock_review bisa
     ditinjau & dikoreksi manual lewat menu "Rekonsiliasi Stok" (Pelaporan) —
     tidak perlu edit DB manual.

Tampilkan ringkasan akhir: berapa produk berhasil direkonstruksi, berapa
needs_stock_review, dan daftar produk yang perlu ditinjau + alasannya.
```

### Versi singkat (kalau Anda sudah paham alurnya)

```
Restore prod ke pos_retail_db sudah selesai. Jalankan migrasi stok ulang sesuai
docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md (langkah 1-9). Ingat: buat
products_stock_backup SEBELUM BE dijalankan; stok dihitung ulang dari riwayat;
jangan ubah skrip backfill yang sudah ada. Kabari hasil tiap tahap.
```

> Setelah migrasi selesai, produk yang `needs_stock_review` sekarang bisa langsung
> ditangani lewat menu **Rekonsiliasi Stok** (admin-only) — lihat
> `docs/DESAIN_MENU_REKONSILIASI_STOK.md`.
