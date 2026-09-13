# Bug: Transaksi "Cangkang" (Header Tanpa Item)

Status: **Kode SUDAH diperbaiki & diverifikasi (lokal). Skrip pembersihan data
SUDAH dijalankan di lokal; SIAP dijalankan di production.**
Tanggal analisis: 13 September 2026
Ditemukan pada: Production (transaksi kasir owner, user_id 3)

---

## 1. Ringkasan

Satu aksi kasir dapat menghasilkan beberapa baris transaksi di tabel `transactions`,
tetapi sebagian di antaranya hanya berupa **header tanpa item** (disebut di sini
"transaksi cangkang"). Header ini:

- Punya `total_amount` terisi, `status = 'completed'`.
- **Tidak** punya baris di `transaction_items`.
- **Tidak** punya baris di `stock_mutations`.

Akibat langsung: detail transaksi tampak kosong / tidak bisa dibuka. Akibat turunan
(lebih penting): angka transaksi cangkang **ikut terhitung** di laporan penjualan
dan perhitungan saldo kas, karena laporan mengagregasi `SUM(total_amount)`
berdasarkan `status = 'completed'` tanpa memeriksa keberadaan item.

Laporan awal dari owner: transaksi WEB-20260913-005/006/007 — tercatat 3, tapi hanya
detail WEB-20260913-007 yang bisa dibuka.

---

## 2. Akar Penyebab (Root Cause)

File: `BE/domain/transaction/repo/transaction_repo.go`, fungsi `createOnce`.

Alur yang diharapkan: `transaction_service.go` membungkus seluruh pembuatan
transaksi dalam **satu** DB transaction:

```go
txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
    txnRepo := s.repo.WithTx(tx)          // txnRepo.db == tx
    result, err := txnRepo.Create(req, userID)
    ...
})
```

Sehingga jika insert item gagal, header ikut ter-rollback.

Namun di dalam `createOnce`, insert **header** dilakukan lewat:

```go
r.db.Connection(func(conn *gorm.DB) error {
    // GET_LOCK, generate kode, INSERT header, LAST_INSERT_ID
})
```

Pada GORM v1.31.1 (`finisher_api.go:609`), `Connection()` memanggil
`sqlDB.Conn(...)` yang **mengambil koneksi baru dari pool** — bukan koneksi
transaksi (`tx`) milik service. Efeknya:

- **Insert header** berjalan di koneksi terpisah, mode **autocommit** →
  langsung ter-commit sendiri, lepas dari transaksi service.
- **Insert item + `ApplyStockDelta`** berjalan lewat `r.db.Exec(...)`
  (koneksi transaksi `tx`).

Ketika langkah item gagal (mis. stok tidak cukup / error validasi paket satuan),
service me-rollback transaksinya — **tetapi header sudah terlanjur ter-commit di
koneksi lain dan tidak ikut ter-rollback.** Hasilnya: header nyangkut, 0 item.

Blok `r.db.Connection` semula ditambahkan untuk serialisasi named lock
(`GET_LOCK`/`RELEASE_LOCK`) saat generate kode transaksi. Niatnya benar
(mencegah race pada nomor urut kode), tapi efek sampingnya memutus konteks
transaksi dan menciptakan celah atomicity ini.

---

## 3. Bukti dari Data Production (lokal = mirror production)

### 3.1 Header 3 transaksi

| Kode | ID | created_at | total_amount | Item | Mutasi stok |
|------|----|-----------|--------------|------|-------------|
| WEB-20260908-004 | 545 | 2026-09-08 02:25:49 | 37.000 | 0 | tidak ada |
| WEB-20260913-005 | 635 | 2026-09-13 03:41:40 | 315.500 | 0 | tidak ada |
| WEB-20260913-006 | 636 | 2026-09-13 03:41:44 | 315.500 | 0 | tidak ada |
| WEB-20260913-007 | 637 | 2026-09-13 03:43:02 | 315.500 | 9 (OK) | 9 (OK) |

Ketiga cangkang berpola identik: `device_source=web`, `payment_method=cash`,
`status=completed`, dibuat user_id 3. Percobaan berhasil hanya di WEB-007.

### 3.2 Dampak ke laporan penjualan cash (user 3)

| Tanggal | Total tercatat | Palsu (cangkang) | Seharusnya |
|---------|---------------|------------------|-----------|
| 2026-09-08 | 1.120.500 | 37.000 | 1.083.500 |
| 2026-09-13 | 1.244.500 | 631.000 | 613.500 |

### 3.3 Dampak ke cash drawer

- **Drawer 21 (8 Sep, `closed`)**: kolom `total_cash_sales` = 1.083.500 (BENAR,
  karena `UpdateSales` tidak jalan untuk cangkang). Namun `closing_balance` /
  `expected_balance` tersimpan **1.120.501**, sedangkan `opening_balance +
  total_cash_sales` = 1.083.501. Selisih persis **37.000** — nilai tersimpan
  tercemar karena saat closing `expected_balance` dihitung ulang dari
  `SUM(total_amount)` (`liveExpectedBalanceExpr`) yang termasuk cangkang.
- **Drawer 26 (13 Sep, `open`)**: kolom tersimpan BENAR (`total_cash_sales`
  613.500, `expected_balance` 713.500). Tetapi tampilan "kas saat ini" di
  layar kasir dihitung live dari `SUM(total_amount)` sehingga tampak +631.000.
  Begitu cangkang dibersihkan, tampilan otomatis benar tanpa mengubah kolom.

Catatan penting: `UpdateSales` **tidak** terpanggil untuk transaksi cangkang
(karena gagal sebelum baris itu), jadi kolom `total_cash_sales` akurat. Yang
tercemar hanya perhitungan berbasis `SUM(total_amount)`.

---

## 4. Audit Cakupan Bug (memastikan tidak ada yang terlewat)

Dilakukan scan seluruh tabel dengan 5 pola untuk mendeteksi segala bentuk
transaksi rusak, bukan hanya 3 yang dilaporkan:

| Pola | Definisi | Hasil |
|------|----------|-------|
| 1 | Header tanpa item sama sekali | **3 baris**: 545, 635, 636 |
| 2 | Ada item tapi `SUM(subtotal)` != `total_amount` (tanpa diskon/pajak) | 0 |
| 3 | Ada item tapi tanpa mutasi stok `out` | 0 |
| 4 | Jumlah item != jumlah mutasi stok `out` | 0 |
| 5 | Kredit `completed` tanpa baris `receivables` | 0 |

**Kesimpulan:** kerusakan data terbatas pada 3 transaksi cangkang (545, 635, 636).
Bug bermanifestasi konsisten sebagai "header saja tanpa item"; tidak ada bentuk
kerusakan parsial lain.

---

## 5. Rencana Perbaikan

### 5.1 Perbaikan kode (mencegah kejadian berulang)

Buat generate kode + insert header berada di dalam transaksi service yang sama,
bukan di koneksi pool baru. Yaitu **hilangkan `r.db.Connection(...)`** yang
memutus konteks transaksi, dan jalankan `GET_LOCK` + insert lewat koneksi `tx`
yang sama (`r.db`). Retry-on-duplicate-key di `Create()` tetap dipertahankan
sebagai lapis kedua.

Prinsip: header, item, mutasi stok, dan receivable harus **atomik** — semua
berhasil atau semua gagal.

### 5.2 Pembersihan data production

1. Backup DB sebelum eksekusi.
2. Void 3 transaksi cangkang (set `status = 'void'`, bukan DELETE — jejak audit
   tetap ada). Transaksi cangkang tidak punya item/mutasi/receivable, jadi tidak
   ada efek samping yang perlu dirollback selain status.
3. Koreksi drawer 21: `closing_balance` & `expected_balance` dari 1.120.501 →
   1.083.501.
4. Drawer 26: tidak perlu diubah; tampilan otomatis benar setelah cangkang
   di-void (karena filter laporan mengecualikan `status = 'void'`).

---

## 6. Verifikasi Setelah Perbaikan

### 6.1 Perbaikan kode — SELESAI & TERUJI

- File diubah: `BE/domain/transaction/repo/transaction_repo.go` (fungsi `createOnce`).
  `r.db.Connection(...)` dihapus; `GET_LOCK` + generate kode + insert header +
  `LAST_INSERT_ID` kini berjalan lewat `r.db` yang sama (= `tx` dari service),
  `RELEASE_LOCK` via `defer`. Alur sync (`ApplySyncTransaction`) TIDAK diubah
  karena polanya berbeda (`r.db.Connection` membungkus `conn.Transaction` di
  dalamnya = sudah atomik, dan `r.db` di sana adalah repo root, bukan `tx`).
- `go build ./...` → sukses (EXIT=0), tanpa error kompilasi.
- Integration test (sementara, sudah dihapus setelah lulus): meniru alur service
  `db.Transaction(func(tx){ repo.WithTx(tx).Create(...) })` dengan item yang
  sengaja gagal (product_id invalid → error TEPAT setelah insert header). Hasil:
  `Create` return error, jumlah baris `transactions` TIDAK bertambah, nomor kode
  web hari ini TIDAK berubah → header ikut rollback. (Dengan kode lama, test ini
  akan GAGAL karena header cangkang tertinggal.) `ok ... EXIT=0`.

### 6.2 Pembersihan data — SUDAH dijalankan di LOKAL (mirror production)

Skrip: `backups/cleanup_transaksi_cangkang_20260913.sql`. Semua verifikasi lulus:

- 545, 635, 636 → status `void`.
- Tidak ada cangkang `completed` tersisa (audit Pola 1 = 0).
- Total penjualan cash: 8 Sep = 1.083.500, 13 Sep = 613.500 (bersih).
- Drawer 21: `closing_balance` = `opening + total_cash_sales - total_expenses`
  (1.083.501) — konsisten.
- Drawer 26 (open): `expected_balance` live = 713.500 (100.000 + 613.500) — benar.

## 7. Prosedur Deploy ke Production

1. **Deploy BE baru** (yang memuat perbaikan `createOnce`). Ini menghentikan
   munculnya cangkang baru.
2. **Backup DB production** sebelum pembersihan, mis.:
   `mysqldump -u <user> -p pos_retail_db > backups/pos_retail_db_before_cleanup_cangkang_20260913.sql`
3. **Jalankan** `backups/cleanup_transaksi_cangkang_20260913.sql` terhadap DB
   production. Skrip berjalan dalam satu transaksi dan menampilkan verifikasi
   sebelum `COMMIT`. Jika muncul baris "VERIFIKASI GAGAL", ganti `COMMIT;` di
   akhir file menjadi `ROLLBACK;` lalu investigasi.
   - Catatan: id 545/635/636 dan angka drawer 21 (1.120.501) sudah dicek identik
     antara lokal dan production. Guard di tiap `UPDATE` membuat skrip idempotent
     dan tidak menyentuh baris lain.
4. **Verifikasi ulang** di production dengan audit Pola 1 (harus 0 cangkang
   `completed`) dan cek total penjualan 8 & 13 Sep.
