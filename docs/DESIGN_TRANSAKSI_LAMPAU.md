# Desain Fitur: Input Transaksi Lampau (Backdate)

> Status: **FINAL** — siap implementasi.

---

## 1. Latar Belakang

Ada kalanya transaksi penjualan belum sempat diinput ke sistem pada hari terjadinya (misal: sistem mati, lupa input, dll). Fitur ini memungkinkan Admin untuk menginput transaksi lampau dengan tanggal yang benar.

---

## 2. Konsep

Tambah **2 menu baru** (khusus Admin):

| Menu Baru | Serupa Dengan | Perbedaan |
|-----------|---------------|-----------|
| **Kas Historis** | Kas Saya | Ada input **kasir** + **tanggal** + pilih **shift** saat buka kas |
| **Kasir Historis** | Kasir | Ada **edit harga** per item + input **jam** saat pembayaran |

Prinsip: **logic sama dengan kasir biasa** — hanya beda di input kasir/tanggal/jam dan edit harga. Dibuat sebagai page & endpoint terpisah agar tidak merusak fitur existing.

---

## 3. Alur

### 3.1 Buka "Kas Historis"

```
Admin buka menu "Kas Historis"
→ Form Buka Kas:
  - Kasir [dropdown, pilih user yang menjadi kasir saat itu]
  - Tanggal [date picker, maks = kemarin, tidak boleh hari ini/masa depan]
  - Shift [dropdown, pilih dari master shift yang ada]
  - Saldo Awal [nominal]
  - Catatan (opsional)
→ Klik "Buka Kas"
→ Kas Historis aktif untuk kasir + tanggal + shift tersebut
→ user_id kas = kasir yang dipilih, created_by = admin yang login
```

### 3.2 Input Transaksi di "Kasir Historis"

```
Admin buka menu "Kasir Historis"
→ UI sama persis dengan Kasir biasa:
  - Search produk, tambah ke keranjang
  - Pilih pelanggan, diskon, pajak
  - PERBEDAAN 1: di keranjang, setiap item ada icon EDIT di sebelah harga
    → klik → bisa ubah harga live (untuk kasus harga waktu itu berbeda)
    → default tetap harga saat ini dari master

→ Klik "Bayar" → Modal Pembayaran:
  - Semua fitur sama (metode, saldo, hutang, uang pas, dll)
  - PERBEDAAN 2: ada input "Jam Transaksi" [time picker, HH:mm]
  - Info tanggal + nama kasir (dari Kas Historis aktif) ditampilkan sebagai label read-only

→ Klik "Proses Bayar"
→ Transaksi tersimpan dengan:
  - transaction_date = tanggal (dari kas) + jam (dari input)
  - user_id = kasir yang dipilih
  - created_by = admin yang login
```

### 3.3 Tutup "Kas Historis"

```
Admin buka menu "Kas Historis"
→ Ringkasan: total transaksi yang sudah diinput untuk tanggal tersebut
→ Tutup Kas
```

---

## 4. Keputusan yang Sudah Disepakati

| # | Keputusan | Catatan |
|---|-----------|---------|
| 1 | 2 menu baru: **Kas Historis** + **Kasir Historis** | Khusus Admin |
| 2 | Stok dikurangi dari stok **saat ini** | Sama seperti checkout biasa |
| 3 | Tanggal dari Kas Historis, jam dari form pembayaran | Digabung jadi `transaction_date` |
| 4 | Tidak ada batasan hari ke belakang | Unlimited |
| 5 | Harga default dari master (saat ini), tapi bisa **edit live** per item | Icon edit di sebelah harga |
| 6 | Shift: pilih dari dropdown master shift saat buka kas | Admin pilih shift yang sesuai |
| 7 | Semua metode bayar sama (tunai, transfer, QRIS, kartu, hutang, saldo) | Tidak ada yang dibatasi |
| 8 | Fitur simpan kembalian ke saldo **berlaku** | Sama seperti kasir biasa |
| 9 | Dibuat sebagai page & endpoint **terpisah** | Tidak merusak fitur existing |
| 10 | Kas biasa (hari ini) dan Kas Historis bisa aktif **bersamaan** | Independen, dibedakan via `is_backdate` |
| 11 | Admin **pilih kasir** saat buka Kas Historis | `user_id` = kasir yang dipilih, `created_by` = admin yang login |
| 12 | Track admin yang input via **`created_by`** (Opsi A) | Tambah kolom `created_by` di `cash_drawer` dan `transactions` — untuk audit multi-admin |
| 13 | Nama menu: **Kas Historis** & **Kasir Historis** | Grup sendiri: "Historis" |
| 14 | Posisi sidebar: grup **"Historis"** setelah "Penjualan" (khusus Admin) | Sidebar existing: Beranda, Penjualan, Produk, Pengadaan, Pelanggan, Keuangan, Pelaporan, Operasional, Sistem |
| 15 | Kode transaksi: **format sama** (`WEB-{tanggal_lampau}-XXX`) | Pembeda backdate: `created_at` ≠ `transaction_date` |
| 16 | Dropdown kasir: tampilkan **semua user aktif** (role apapun) | Pakai API `POST /users/list` yang sudah ada |
| 17 | `created_by` diisi dari **token admin yang login** | Untuk audit: siapa admin yang menginput backdate |

---

## 5. Perbedaan dengan Kasir Biasa (Ringkas)

| Aspek | Kasir Biasa | Kasir Historis |
|-------|-------------|----------------|
| Tanggal | Otomatis hari ini | Dari Kas Historis (manual) |
| Jam | Otomatis saat checkout | Input manual (time picker) |
| Harga | Dari master, read-only | Dari master + bisa edit live |
| Shift | Otomatis dari kas aktif | Pilih manual saat buka kas |
| Kasir (user_id) | User yang login | Admin pilih saat buka kas |
| created_by | NULL (tidak ada) | Admin yang login & menginput |
| Akses | Semua role | Admin only |
| Kas | Kas Saya | Kas Historis |

---

## 6. Implementasi — Scope Perubahan

### 6.1 FE (page baru, reuse komponen)

| Item | Keterangan |
|------|-----------|
| `features/sales/backdate-cash/` | Page "Kas Historis" — form buka/tutup kas dengan **dropdown kasir** + date picker + shift dropdown |
| `features/sales/backdate-cashier/` | Page "Kasir Historis" — reuse komponen kasir, tambah edit harga + time picker |
| Route baru | `/backdate/cash` + `/backdate/cashier` |
| Menu sidebar | Grup "Historis" (admin only) |

### 6.2 BE (endpoint baru, reuse service logic)

| Endpoint | Keterangan |
|----------|-----------|
| `POST /backdate/cash-drawer/open` | Buka kas historis. Body: `{ user_id, date, shift_id, opening_balance, notes }`. Simpan: `user_id` = kasir dipilih, `created_by` = admin login, `is_backdate` = 1, `open_time` = tanggal input |
| `GET /backdate/cash-drawer/current` | Get kas historis aktif milik admin yang sedang login (query by `created_by` + `is_backdate=1` + `status=open`) |
| `POST /backdate/cash-drawer/close/:id` | Tutup kas historis |
| `POST /backdate/transactions/create` | Checkout backdate. Body tambahan: `{ transaction_time }`. `transaction_date` = tanggal kas + jam input, `user_id` = dari kas backdate, `created_by` = admin login |
| Permission | Menu key baru: `historis.kas_historis` + `historis.kasir_historis` — hanya Admin |

**Catatan endpoint GET current**: Karena `user_id` di kas historis adalah kasir yang dipilih (bukan admin), maka query kas historis aktif harus pakai `created_by = admin_login AND is_backdate = 1 AND status = 'open'`.

### 6.3 DB (Migration 008)

```sql
-- 008_backdate_feature.sql

-- Pembeda kas historis vs kas biasa
ALTER TABLE cash_drawer ADD COLUMN is_backdate TINYINT(1) NOT NULL DEFAULT 0;

-- Audit: siapa admin yang buka kas historis / input transaksi
ALTER TABLE cash_drawer ADD COLUMN created_by INT NULL;
ALTER TABLE transactions ADD COLUMN created_by INT NULL;

-- Relasi transaksi → kas (untuk void rollback yang akurat)
ALTER TABLE transactions ADD COLUMN cash_drawer_id INT NULL;

-- Index untuk performa query backdate
ALTER TABLE cash_drawer ADD INDEX idx_cd_backdate (is_backdate, status);
ALTER TABLE transactions ADD INDEX idx_tx_cash_drawer (cash_drawer_id);
```

---

## 7. Analisis Gap & Potensi Bug

Setelah review kode existing, berikut masalah yang HARUS ditangani saat implementasi:

### 7.1 GAP KRITIS

| # | Area | Kode Existing | Masalah | Solusi |
|---|------|---------------|---------|--------|
| 1 | **Transaction date** | `createOnce()` → `now := time_helper.GetTimeNow()` | Tanggal selalu hari ini | Endpoint backdate pass `target_date` parameter, bukan pakai `now` |
| 2 | **Lock key** | `lockName := fmt.Sprintf("txcode:%s:%s", time_helper.ToSQLDate(now), req.DeviceSource)` | Lock pakai hari ini, backdate akan tabrakan kode | Lock key pakai tanggal target: `txcode:{target_date}:{device}` |
| 3 | **Code generation** | `generateTransactionCodeQuery` → `WHERE DATE(transaction_date) = ? AND device_source = ?` menggunakan `now` | Kode sequential berdasarkan hari ini | Untuk backdate, kirim tanggal target ke query — kode jadi `WEB-{tanggal_lampau}-NNN` |
| 4 | **Cash drawer conflict** | `GetOpenCashDrawer` → `WHERE user_id=? AND status='open' LIMIT 1` | Kasir X bisa punya kas biasa + kas historis open → query ambil salah satu | Tambah filter `is_backdate`: kas biasa = `is_backdate=0`, kas historis = `is_backdate=1` |
| 5 | **Auto-close scheduler** | `getOpenYesterdayQuery` → `WHERE status='open' AND DATE(open_time) < ?` | Scheduler tutup SEMUA termasuk kas historis yang open_time lampau | Exclude: `AND (is_backdate = 0 OR is_backdate IS NULL)` |
| 6 | **Void rollback kas** | `Void()` → `cashDrawerRepo.GetOpenCashDrawer(t.UserID)` | Void transaksi backdate → cari kas OPEN milik kasir saat ini → rollback ke kas biasa yang salah | Void harus cek: jika transaksi punya `created_by` != NULL → ini backdate → skip kas rollback jika kas historis sudah closed |

### 7.2 GAP MINOR (Acceptable)

| # | Area | Masalah | Keputusan |
|---|------|---------|-----------|
| 7 | `stock_mutations.created_at` | Mutasi stok `created_at` = hari ini, tapi `transaction_date` = lampau | **OK** — `created_at` = kapan fisik diinput, `transaction_date` = kapan bisnis terjadi |
| 8 | `liveExpectedBalanceExpr` | Formula: `SUM(total_amount) FROM transactions WHERE user_id = cd.user_id AND transaction_date >= cd.open_time` | **OK** — Kas historis punya `open_time` = tanggal lampau, jadi transaksi backdate yg `transaction_date` = tanggal lampau akan terhitung di kas historis ini, bukan kas biasa |

### 7.3 Solusi Void yang Dipilih

```
SKENARIO VOID TRANSAKSI BACKDATE:
1. Void transaksi yg punya created_by != NULL → ini transaksi backdate
2. Cek cash_drawer via cash_drawer_id di transaksi (KOLOM BARU: transactions.cash_drawer_id)
3. Jika kas tersebut masih open → rollback sales di kas tersebut
4. Jika kas sudah closed → SKIP rollback kas (stok tetap dikembalikan, hanya kas tidak di-update)
5. Saldo pelanggan tetap di-rollback seperti biasa (tidak terpengaruh)
```

**Alasan tambah `cash_drawer_id`**: Void harus tahu transaksi ini masuk ke kas mana. Tanpa field ini, harus lookup berdasarkan tanggal & user yang rawan error. Kolom ini sudah termasuk di migration 008 (section 6.3).

---

## 8. Detail Modifikasi Kode

### 8.1 Query EXISTING yang Perlu Ditambah Filter `is_backdate`

| # | Query | File | Perubahan |
|---|-------|------|-----------|
| 1 | `getCurrentCashDrawerQuery` | `cash_drawer_repo.go:16` | Tambah `AND (cd.is_backdate = 0 OR cd.is_backdate IS NULL)` |
| 2 | `getOpenCashDrawerQuery` | `cash_drawer_repo.go:17` | Tambah `AND (cd.is_backdate = 0 OR cd.is_backdate IS NULL)` |
| 3 | `getMyCashQuery` | `cash_drawer_repo.go:111` | Tambah `AND (cd.is_backdate = 0 OR cd.is_backdate IS NULL)` |
| 4 | `getOpenYesterdayQuery` | `cash_drawer_repo.go:68` | Tambah `AND (is_backdate = 0 OR is_backdate IS NULL)` |

**Aman karena**: default `is_backdate = 0`, semua data lama tetap masuk filter.

### 8.2 File EXISTING yang Perlu Dimodifikasi

| # | File | Perubahan | Resiko |
|---|------|-----------|--------|
| 1 | `BE/domain/cash_drawer/repo/cash_drawer_repo.go` | 4 query ditambah filter (lihat 8.1) | Rendah — default=0, backward-safe |
| 2 | `BE/domain/cash_drawer/model/cash_drawer.go` | Tambah field `IsBackdate`, `CreatedBy` di struct | Rendah — GORM scan aman |
| 3 | `BE/domain/transaction/model/transaction.go` | Tambah field `BalanceUsed`, `CreatedBy`, `CashDrawerID` di struct | Rendah — field sudah ada di DB, struct belum |
| 4 | `BE/domain/transaction/service/transaction_service.go` | Void: tambah cek `created_by != nil` → skip kas jika closed | Sedang — logic void berubah |
| 5 | `BE/routes.go` (atau `protected_routes.go`) | Tambah route grup `/api/backdate/...` | Rendah — tambah saja |
| 6 | `FE/src/routes/` | Tambah route `/backdate/cash`, `/backdate/cashier` | Rendah — tambah saja |
| 7 | `FE (sidebar via backend menu)` | Tambah menu record di DB seeder/migration | Rendah — menu diload dari DB |

### 8.3 File BARU yang Perlu Dibuat

| # | File | Tujuan |
|---|------|--------|
| 1 | `BE/database/migrations/008_backdate_feature.sql` | ALTER TABLE: is_backdate, created_by, cash_drawer_id + menu seed |
| 2 | `BE/domain/backdate/handler/backdate_handler.go` | Handler endpoint backdate (open/close/current kas + create transaction) |
| 3 | `BE/domain/backdate/service/backdate_service.go` | Service logic: validasi, buka/tutup kas, create transaksi backdate |
| 4 | `BE/domain/backdate/repo/backdate_repo.go` | Repo: query kas backdate + insert transaksi dengan custom date |
| 5 | `BE/domain/backdate/dto/backdate_dto.go` | Request/response DTOs |
| 6 | `FE/src/features/sales/backdate-cash/` | Page Kas Historis (form + tampilan aktif) |
| 7 | `FE/src/features/sales/backdate-cashier/` | Page Kasir Historis (reuse + edit harga + time picker) |
| 8 | `FE/src/features/sales/backdate-cashier/components/EditablePrice.tsx` | Komponen edit harga live per item |

---

## 9. Checklist Keamanan (Pastikan Tidak Merusak Existing)

| # | Fitur Existing | Terpengaruh? | Tindakan |
|---|----------------|:------------:|----------|
| 1 | Kasir biasa | ❌ | Endpoint terpisah, logic terpisah |
| 2 | Kas Saya | ❌ | Query ditambah `AND is_backdate=0` (backward-safe) |
| 3 | Laporan penjualan | ✅ Auto benar | Filter by `transaction_date` — backdate masuk tanggal aslinya |
| 4 | Void transaksi | ⚠️ Perlu handle | Void backdate → cari kas via `cash_drawer_id`, skip jika closed |
| 5 | Piutang/Saldo | ❌ | Logic sama, tidak peduli tanggal |
| 6 | Stok | ❌ | Dikurangi dari stok saat ini (sama seperti checkout biasa) |
| 7 | Scheduler auto-close | ⚠️ Perlu update | Exclude `is_backdate=1` dari auto-close |
| 8 | Dashboard hari ini | ❌ | Filter `DATE(transaction_date) = today` — backdate tidak muncul |
| 9 | Riwayat transaksi | ✅ Auto benar | List by `transaction_date`, backdate muncul di tanggal aslinya |
| 10 | Laporan Kinerja Kasir | ✅ Auto benar | Filter by `user_id` — backdate muncul di kasir yang dipilih |
| 11 | Buka kas biasa | ❌ | `GetOpenCashDrawer` exclude is_backdate=1, tidak konflik |

---

## 10. Catatan Implementasi Penting

### 10.1 `liveExpectedBalanceExpr` untuk Kas Historis

Formula existing:
```sql
(cd.opening_balance 
  + SUM(transactions WHERE user_id = cd.user_id AND payment_method='cash' AND status='completed' AND transaction_date >= cd.open_time) 
  - SUM(expenses WHERE user_id = cd.user_id AND created_at >= cd.open_time))
```

**Untuk kas historis ini AMAN** karena:
- `cd.open_time` kas historis = tanggal lampau (misal `2026-08-10 08:00:00`)
- Transaksi backdate yang kita buat punya `transaction_date` = tanggal lampau juga
- `transaction_date >= cd.open_time` → **match** (transaksi masuk perhitungan kas historis)
- Kas biasa punya `open_time` = hari ini → transaksi backdate dengan tanggal lampau **TIDAK** masuk kas biasa (karena `transaction_date < cd.open_time` kas biasa)

**TAPI ada edge case**: Formula juga filter by `user_id = cd.user_id`. Jika kasir yang dipilih (user_id) JUGA sedang punya kas biasa yang open, maka:
- Kas biasa kasir X: `open_time = hari ini` → transaksi backdate (tanggal lampau) TIDAK masuk (karena `transaction_date < open_time`)
- Kas historis kasir X: `open_time = tanggal lampau` → transaksi backdate MASUK (karena `transaction_date >= open_time`)

**Kesimpulan: AMAN tanpa modifikasi formula.**

### 10.2 Endpoint GET Current Kas Historis

Karena admin yang login ≠ user_id kas, query kas historis aktif harus:
```sql
SELECT ... FROM cash_drawer 
WHERE created_by = ? AND is_backdate = 1 AND status = 'open' 
LIMIT 1
```
(bukan `WHERE user_id = ?`)

### 10.3 Flow `cash_drawer_id` pada Transaksi

Saat **create transaksi biasa**: `cash_drawer_id` diisi dari kas yang open (via `GetOpenCashDrawer`).
Saat **create transaksi backdate**: `cash_drawer_id` diisi dari kas historis yang open.
Saat **void**: lookup kas via `cash_drawer_id`, jika masih open → rollback, jika closed → skip.

**Catatan**: Untuk transaksi lama (sebelum migration), `cash_drawer_id` = NULL. Void transaksi lama tetap pakai logic existing (`GetOpenCashDrawer`).

### 10.4 `openCashDrawerQuery` untuk Backdate

Query existing:
```sql
INSERT INTO cash_drawer (user_id, shift_id, open_time, opening_balance, open_notes, status) VALUES (?, ?, ?, ?, ?, 'open')
```

Untuk backdate, buat query baru:
```sql
INSERT INTO cash_drawer (user_id, shift_id, open_time, opening_balance, open_notes, status, is_backdate, created_by) 
VALUES (?, ?, ?, ?, ?, 'open', 1, ?)
```

**Tidak perlu modifikasi query existing** — cukup buat query baru di backdate repo.

### 10.5 `createTransactionQuery` untuk Backdate

Query existing:
```sql
INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, ..., device_source) 
VALUES (?, ?, ?, ?, ..., ?)
```

Untuk backdate, buat query baru yang tambah `created_by` + `cash_drawer_id`:
```sql
INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, ..., device_source, created_by, cash_drawer_id) 
VALUES (?, ?, ?, ?, ..., ?, ?, ?)
```

---

## 11. Pertanyaan Terbuka

Tidak ada — desain sudah final, siap implementasi.
