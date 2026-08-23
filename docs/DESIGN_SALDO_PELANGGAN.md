# Desain Fitur: Saldo Pelanggan (Customer Deposit)

> Dokumen ini mencatat desain lengkap fitur Saldo Pelanggan sebelum implementasi.
> Status: **FINAL** — siap implementasi.

---

## 1. Latar Belakang

### Use Case Real

Customer A datang ke toko, beli air minum Rp 5.000 tapi bayar Rp 200.000. Sisa Rp 195.000 disimpan di toko sebagai "deposit/titipan" untuk belanja berikutnya tanpa perlu bayar lagi setiap kali datang.

### Kondisi Existing Saat Ini

| Fitur | Status | Catatan |
|-------|--------|---------|
| Checkout Tunai | ✅ Ada | Kas harian bertambah senilai `total_amount` |
| Checkout Transfer/QRIS/Kartu | ✅ Ada | Kas harian TIDAK bertambah |
| Checkout Kredit (hutang pelanggan) | ✅ Ada | Kas TIDAK bertambah, piutang terbuat otomatis |
| Deposit/Saldo Pelanggan | ❌ Belum ada | — |

### Existing: Bagaimana Kas Harian Bekerja

Dari kode `transaction_service.go`:
```go
if req.PaymentMethod == "cash" {
    cashDrawerRepo.UpdateSales(drawer.ID, req.TotalAmount, req.TotalAmount, time)
}
```

Saat ini kas harian dicatat = `total_amount` (nilai transaksi). Ini **perlu diubah** menjadi `payment_amount` (uang yang benar-benar masuk laci) agar konsisten dengan Opsi B.

---

## 2. Keputusan Desain yang Sudah Disepakati

| # | Keputusan | Alasan |
|---|-----------|--------|
| 1 | **Opsi B: Kas harian = uang fisik masuk laci** | Kas selalu cocok dengan isi laci saat tutup shift |
| 2 | Pakai saldo → kas TIDAK bertambah | Uang sudah masuk kas saat top-up |
| 3 | Top-up saldo (dari kembalian) → kas TIDAK bertambah ekstra | Kas sudah +payment_amount saat checkout; top-up hanya pencatatan saldo |
| 4 | Top-up manual → kas BERTAMBAH | Uang fisik masuk laci di luar transaksi jual beli |
| 5 | Kombinasi saldo + tunai dibolehkan | Saldo < total → sisa bayar tunai |
| 6 | Kombinasi saldo + hutang dibolehkan | Saldo < total → sisa jadi piutang |
| 7 | Default credit limit pelanggan = 0 (bukan tak terbatas) | Mencegah hutang tanpa batas |
| 8 | Rename "Kredit" → "Hutang" | Lebih jelas untuk kasir retail |
| 9 | "Gunakan Saldo" = checkbox, bukan metode pembayaran | Saldo dipotong dulu, sisa baru pilih metode |
| 10 | "Gunakan Saldo" HANYA muncul jika pelanggan dipilih DAN saldo > 0 | Tanpa pelanggan = UI tanpa saldo |
| 11 | "Hutang" HANYA muncul jika pelanggan dipilih | Sama seperti "Kredit" existing |
| 12 | Tanpa pelanggan → modal hanya [Tunai] [Transfer] [QRIS] [Kartu] | Tidak ada "Hutang", tidak ada "Gunakan Saldo" |
| 13 | "Simpan kembalian ke saldo" muncul SETELAH checkout di halaman struk | Bukan di modal pembayaran |

---

## 3. Desain Database

### 3.1 Tambah kolom di `customers`

```sql
ALTER TABLE customers ADD COLUMN balance DECIMAL(15,2) NOT NULL DEFAULT 0;
```

### 3.2 Tambah kolom di `transactions`

```sql
ALTER TABLE transactions ADD COLUMN balance_used DECIMAL(15,2) NOT NULL DEFAULT 0;
```

- `balance_used = 0` → tidak pakai saldo (default, backward-compatible)
- `balance_used > 0` → sebagian/seluruh dibayar dari saldo

### 3.3 Tambah value di `payment_methods`

```sql
INSERT INTO payment_methods (code, label, is_active, sort_order) VALUES ('balance', 'Saldo', 1, 6);
```

Dipakai saat seluruh transaksi dibayar penuh dari saldo (tanpa metode lain).

### 3.4 Tabel baru: `customer_balance_mutations`

```sql
CREATE TABLE IF NOT EXISTS customer_balance_mutations (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    customer_id     INT           NOT NULL,
    amount          DECIMAL(15,2) NOT NULL,  -- +positif = masuk, -negatif = keluar
    balance_after   DECIMAL(15,2) NOT NULL,  -- saldo setelah mutasi ini
    type            ENUM('topup', 'usage', 'refund', 'adjustment') NOT NULL,
    reference_type  VARCHAR(50)   NULL,      -- 'transaction', 'manual', 'void'
    reference_id    INT           NULL,      -- ID transaksi terkait (nullable)
    notes           TEXT          NULL,
    user_id         INT           NOT NULL,
    created_at      DATETIME      DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    INDEX idx_cbm_customer (customer_id),
    INDEX idx_cbm_type (type),
    INDEX idx_cbm_ref (reference_type, reference_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.5 Contoh Data Tersimpan per Skenario

| Skenario | payment_method | total_amount | payment_amount | balance_used | change_amount | is_credit | Piutang |
|----------|---------------|-------------|---------------|-------------|--------------|-----------|---------|
| Saldo cukup | balance | 15.000 | 0 | 15.000 | 0 | false | — |
| Saldo + Tunai | cash | 50.000 | 20.000 | 30.000 | 0 | false | — |
| Saldo + Hutang | kredit | 30.000 | 0 | 10.000 | 0 | true | 20.000 |
| Tunai biasa | cash | 22.500 | 25.000 | 0 | 2.500 | false | — |
| Hutang penuh | kredit | 30.000 | 0 | 0 | 0 | true | 30.000 |
| Transfer biasa | transfer | 22.500 | 22.500 | 0 | 0 | false | — |

**Rumus:**
```
effective_total = total_amount - balance_used
piutang (jika is_credit) = effective_total
kas_harian (jika cash) = payment_amount
```

### 3.6 Tipe Mutasi Saldo

| Type | Deskripsi | Amount | Trigger |
|------|-----------|--------|---------|
| `topup` | Customer setor/titip uang | + | Simpan kembalian, top-up manual |
| `usage` | Pakai saldo untuk belanja | - | Checkout dengan saldo |
| `refund` | Pengembalian saldo ke customer | - | Tarik manual, void transaksi topup |
| `adjustment` | Koreksi manual owner | +/- | Fix saldo salah |

---

## 4. Alur & Skenario

### Skenario A: Customer Titip Uang (Top-up dari Kembalian)

```
Customer A beli air Rp 5.000, bayar tunai Rp 200.000
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│ Modal Pembayaran                                │
│ Total: Rp 5.000                                 │
│ Metode: Tunai                                   │
│ Jumlah Bayar: Rp 200.000                        │
│ Kembalian: Rp 195.000                           │
│ → Klik "Proses Bayar"                           │
└─────────────────────────────────────────────────┘
                    │ transaksi berhasil
                    ▼
┌─────────────────────────────────────────────────┐
│ Struk Transaksi                                 │
│ ...                                             │
│ Kembalian: Rp 195.000                           │
│                                                 │
│ ┌─────────────────────────────────────────────┐ │
│ │ 💰 Simpan kembalian ke Saldo Pelanggan?     │ │
│ │                                             │ │
│ │ Jumlah: [Rp 195.000] (bisa diedit)          │ │
│ │                                             │ │
│ │ [Tidak, kembalian tunai] [Ya, Simpan]       │ │
│ └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
                    │ "Ya, Simpan"
                    ▼
Efek:
- Transaksi: payment_method=cash, payment_amount=200000, change=195000
- Kas harian: +Rp 200.000 (uang fisik masuk laci, ini dari checkout biasa)
- Saldo customer: +Rp 195.000
- Mutasi: type=topup, amount=+195000, ref=transaction/{id}
- Struk update: "Saldo ditambahkan: Rp 195.000 | Saldo Anda: Rp 195.000"
```

**Catatan**: Kas harian +Rp 200.000 sudah terjadi dari proses checkout biasa (payment_amount). Top-up saldo TIDAK menambah kas lagi — hanya memindahkan pencatatan dari "kembalian fisik" menjadi "saldo digital".

### Skenario B: Customer Belanja Pakai Saldo (Saldo Cukup)

```
Customer A beli air Rp 15.000
Saldo saat ini: Rp 195.000
                    │
                    ▼
┌──────────────────────────────────────────────────┐
│ Modal Pembayaran                                 │
│                                                  │
│ Total Belanja: Rp 15.000                         │
│                                                  │
│ ┌──────────────────────────────────────────────┐ │
│ │ ☑ Gunakan Saldo              Rp 195.000      │ │
│ │   Saldo terpakai: Rp 15.000                  │ │
│ │   Sisa saldo setelah: Rp 180.000             │ │
│ │   Sisa bayar: Rp 0                           │ │
│ └──────────────────────────────────────────────┘ │
│                                                  │
│ Status: Lunas via Saldo ✓                        │
│                                                  │
│ [Batal] [✓ Proses Bayar]                         │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi: payment_method='balance', balance_used=15000, payment_amount=0
- Kas harian: TIDAK bertambah
- Saldo customer: -Rp 15.000 → sisa Rp 180.000
- Mutasi: type=usage, amount=-15000, ref=transaction/{id}
- Stok berkurang seperti biasa
```

### Skenario C: Kombinasi Saldo + Tunai

```
Customer A beli barang Rp 50.000
Saldo saat ini: Rp 30.000
                    │
                    ▼
┌──────────────────────────────────────────────────┐
│ Modal Pembayaran                                 │
│                                                  │
│ Total Belanja: Rp 50.000                         │
│                                                  │
│ ┌──────────────────────────────────────────────┐ │
│ │ ☑ Gunakan Saldo              Rp 30.000       │ │
│ │   Saldo terpakai: Rp 30.000 (habis)          │ │
│ │   Sisa bayar: Rp 20.000                      │ │
│ └──────────────────────────────────────────────┘ │
│                                                  │
│ Metode Pembayaran (sisa Rp 20.000):              │
│ [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]      │
│  ^^^^^ ← dipilih                                │
│                                                  │
│ Jumlah Bayar: [Rp 20.000]                        │
│ Kembalian: Rp 0                                  │
│                                                  │
│ [Batal] [✓ Proses Bayar]                         │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi: payment_method='cash', balance_used=30000, payment_amount=20000
- Kas harian: +Rp 20.000 (hanya porsi tunai)
- Saldo customer: -Rp 30.000 → sisa Rp 0
- Mutasi saldo: type=usage, amount=-30000, ref=transaction/{id}
```

### Skenario C2: Kombinasi Saldo + Hutang

```
Customer A beli barang Rp 30.000
Saldo saat ini: Rp 10.000
                    │
                    ▼
┌──────────────────────────────────────────────────┐
│ Modal Pembayaran                                 │
│                                                  │
│ Total Belanja: Rp 30.000                         │
│                                                  │
│ ┌──────────────────────────────────────────────┐ │
│ │ ☑ Gunakan Saldo              Rp 10.000       │ │
│ │   Saldo terpakai: Rp 10.000 (habis)          │ │
│ │   Sisa bayar: Rp 20.000                      │ │
│ └──────────────────────────────────────────────┘ │
│                                                  │
│ Metode Pembayaran (sisa Rp 20.000):              │
│ [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]      │
│                                     ^^^^^^       │
│                                     dipilih      │
│                                                  │
│ ⚠️ Piutang yang akan terbuat: Rp 20.000          │
│    Pelanggan: Customer A                         │
│                                                  │
│ [Batal] [✓ Proses Bayar]                         │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi: payment_method='kredit', balance_used=10000, is_credit=true, payment_amount=0
- Kas harian: TIDAK bertambah
- Saldo customer: -Rp 10.000 → sisa Rp 0
- Piutang terbuat: Rp 20.000 (= total_amount - balance_used)
- Mutasi saldo: type=usage, amount=-10000, ref=transaction/{id}
```

### Skenario D: Top-up Manual (tanpa transaksi belanja)

```
Admin buka menu Pelanggan → Detail → "Top-up Saldo"
                    │
                    ▼
Form: Jumlah [Rp 100.000], Catatan: "Customer titip uang"
                    │
                    ▼
Efek:
- Saldo customer: +Rp 100.000
- Mutasi: type=topup, amount=+100000, ref=manual, notes="Customer titip uang"
- Kas harian: +Rp 100.000 (uang fisik masuk laci — via update kas endpoint)
```

**Catatan**: Top-up manual berarti customer datang dan kasih uang fisik tanpa beli apa-apa. Uang masuk laci, maka kas harian harus bertambah. Ini diimplementasi via endpoint terpisah yang juga update cash drawer.

### Skenario E: Pengembalian Saldo (Refund)

```
Admin buka menu Pelanggan → Detail → "Tarik Saldo"
                    │
                    ▼
Form: Jumlah [Rp 50.000], Catatan: "Customer minta uang balik"
                    │
                    ▼
Efek:
- Saldo customer: -Rp 50.000
- Mutasi: type=refund, amount=-50000, ref=manual
- Kas harian: pengeluaran +Rp 50.000 (dicatat via expense supaya tutup shift cocok)
```

---

## 5. Potensi Celah/Bug yang Harus Dicegah

| # | Celah | Solusi |
|---|-------|--------|
| 1 | Double-count kas: top-up masuk kas, pakai saldo masuk kas lagi | Pakai saldo → kas TIDAK bertambah. Top-up dari kembalian → kas sudah masuk via checkout biasa |
| 2 | Saldo negatif (pakai lebih dari yang ada) | Validasi BE: `balance >= balance_used`, pakai `SELECT ... FOR UPDATE` |
| 3 | Race condition: 2 kasir proses saldo bersamaan | `FOR UPDATE` lock pada row customer di dalam DB transaction |
| 4 | Refund lebih dari saldo | Validasi: refund amount ≤ current balance |
| 5 | Top-up dari kembalian tapi customer belum dipilih | Opsi "Simpan ke Saldo" hanya render jika customer_id ada + kembalian > 0 |
| 6 | Void transaksi yang pakai saldo | Void handler: kembalikan `balance += balance_used` + buat mutasi type=refund ref=void |
| 7 | Void transaksi yang di-topup saldo dari kembaliannya | Void handler: cari mutasi top-up dengan ref=transaction/{id}, rollback: `balance -= mutasi.amount` |
| 8 | Tutup kas tidak cocok setelah refund manual | Refund dicatat sebagai expense di kas harian |
| 9 | Laporan penjualan harus tetap mencatat semua transaksi | Ya — laporan SEMUA transaksi termasuk yang bayar saldo (revenue tetap dihitung) |
| 10 | Pelanggan dihapus tapi masih punya saldo | Validasi: tolak hapus/nonaktifkan jika balance > 0 |
| 11 | Void transaksi saldo+hutang: harus rollback KEDUANYA | Void: kembalikan saldo (balance_used) DAN void piutang (effective_total) |
| 12 | balance_used > total_amount (manipulasi API langsung) | Validasi BE: `0 ≤ balance_used ≤ total_amount` |
| 13 | balance_used > 0 tanpa customer | Validasi BE: `balance_used > 0 → customer_id wajib ada` |
| 14 | Top-up di struk: kasir klik "Ya" 2x cepat (double topup) | FE: disable button setelah klik. BE: cek tidak ada mutasi topup dengan ref=transaction/{id} yang sama |

---

## 6. Perubahan UI

### A. Modal Pembayaran (Kasir)

**Tanpa pelanggan (checkbox "Tambah Pelanggan" tidak dicentang):**
```
Metode: [Tunai] [Transfer] [QRIS] [Kartu]
```
Tidak ada "Hutang", tidak ada "Gunakan Saldo". Persis seperti modal tunai biasa.

**Pelanggan dipilih, saldo > 0:**
```
┌────────────────────────────────────────┐
│ ☑ Gunakan Saldo          Rp 195.000   │
│   Saldo terpakai: Rp X                │
│   Sisa bayar: Rp Y                    │
└────────────────────────────────────────┘
Metode (untuk sisa): [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]
```

**Pelanggan dipilih, saldo = 0:**
```
Metode: [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]
```
Tidak ada checkbox saldo (karena saldo kosong).

**Aturan visibilitas:**
| Kondisi | "Gunakan Saldo" | Metode | "Hutang" |
|---------|----------------|--------|----------|
| Tanpa pelanggan | Hidden | Tunai/Transfer/QRIS/Kartu | Hidden |
| Pelanggan dipilih, saldo = 0 | Hidden | Tunai/Transfer/QRIS/Kartu/Hutang | Visible |
| Pelanggan dipilih, saldo > 0 | Visible (checkbox) | Tunai/Transfer/QRIS/Kartu/Hutang | Visible |

### B. After Checkout — Opsi Simpan ke Saldo

Muncul di halaman struk **hanya jika**:
- Pelanggan dipilih (customer_id ada)
- Metode pembayaran = tunai
- Kembalian > 0

```
┌──────────────────────────────────────────┐
│ 💰 Simpan kembalian ke Saldo Pelanggan?  │
│ Jumlah: [Rp 195.000] (bisa diedit)      │
│ [Tidak] [Ya, Simpan]                    │
└──────────────────────────────────────────┘
```

### C. Halaman Pelanggan (List)

Tambah kolom **"Saldo"** di tabel list pelanggan.

### D. Detail Pelanggan

Tambah section **"Saldo & Riwayat"**:
- Saldo saat ini (angka besar, prominent)
- Tombol: "Top-up Saldo" + "Tarik Saldo" (sesuai permission)
- Tabel riwayat mutasi: tanggal, tipe, jumlah, saldo setelah, catatan, oleh siapa

### E. Struk

Jika transaksi menggunakan saldo (`balance_used > 0`):
```
Saldo digunakan: Rp 15.000
Sisa Saldo Anda: Rp 180.000
```

---

## 7. Perubahan Backend

### 7.1 Endpoint Baru

| Method | Path | Fungsi | Role |
|--------|------|--------|------|
| POST | `/customers/:id/balance/topup` | Top-up manual | Admin, Owner |
| POST | `/customers/:id/balance/refund` | Tarik/refund saldo | Admin, Owner |
| POST | `/customers/:id/balance/history` | Riwayat mutasi saldo | Semua |
| POST | `/transactions/save-to-balance` | Simpan kembalian ke saldo | Semua |

### 7.2 Perubahan Endpoint Existing

| Endpoint | Perubahan |
|----------|-----------|
| `POST /transactions/create` | Tambah field `balance_used`. Validasi, deduct saldo, buat mutasi. Piutang = total - balance_used |
| `POST /transactions/void/:id` | Jika balance_used > 0: kembalikan saldo + mutasi. Cek juga top-up terkait (skenario void+topup) |
| `POST /customers/detail/:id` | Tambah field `balance` di response |
| `POST /customers/delete/:id` | Validasi: tolak jika balance > 0 |
| `POST /customers/toggle-status/:id` | Validasi: tolak nonaktif jika balance > 0 |

### 7.3 Logic `transaction_service.go` Create() — Pseudocode

```go
func Create(req) {
    // 1. Validasi basic (shift, stok, dll — existing)
    
    // 2. Hitung effective total
    effectiveTotal := req.TotalAmount - req.BalanceUsed
    
    // 3. Validasi saldo
    if req.BalanceUsed > 0 {
        if req.CustomerID == nil → error "Pilih pelanggan"
        if req.BalanceUsed > req.TotalAmount → error "Saldo melebihi total"
    }
    
    // 4. Validasi pembayaran (hanya cek sisa setelah saldo)
    if !req.IsCredit && effectiveTotal > 0 && req.PaymentAmount < effectiveTotal {
        → error "Jumlah pembayaran kurang"
    }

    // 5. DB Transaction
    tx {
        // a. Lock & deduct saldo customer (jika balance_used > 0)
        if req.BalanceUsed > 0 {
            customer = SELECT ... FOR UPDATE WHERE id = customer_id
            if customer.Balance < req.BalanceUsed → error "Saldo tidak cukup"
            UPDATE customers SET balance = balance - req.BalanceUsed
            INSERT customer_balance_mutations (type=usage, amount=-req.BalanceUsed, ...)
        }
        
        // b. Insert transaksi (termasuk balance_used)
        INSERT transactions (... balance_used = req.BalanceUsed ...)
        
        // c. Insert items + kurangi stok (existing logic)
        
        // d. Insert piutang (jika hutang)
        if req.IsCredit && effectiveTotal > 0 {
            INSERT receivables (total_amount = effectiveTotal, remaining = effectiveTotal)
        }
        
        // e. Update kas harian (hanya cash, dan hanya payment_amount)
        if req.PaymentMethod == "cash" && req.PaymentAmount > 0 {
            UPDATE cash_drawer ... total_sales += req.PaymentAmount
        }
    }
}
```

### 7.4 Logic Void — Pseudocode

```go
func Void(transactionID) {
    // ... existing void logic (rollback stok, void piutang) ...
    
    // Tambahan: rollback saldo
    if transaction.BalanceUsed > 0 {
        UPDATE customers SET balance = balance + transaction.BalanceUsed
        INSERT customer_balance_mutations (type=refund, amount=+BalanceUsed, ref=void/{id})
    }
    
    // Tambahan: rollback top-up (jika kembalian pernah disimpan ke saldo)
    topupMutation = SELECT FROM customer_balance_mutations 
                    WHERE reference_type='transaction' AND reference_id={id} AND type='topup'
    if topupMutation != nil {
        UPDATE customers SET balance = balance - topupMutation.Amount
        INSERT customer_balance_mutations (type=adjustment, amount=-topupMutation.Amount, ref=void/{id})
    }
}
```

### 7.5 Kas Harian — Aturan Lengkap

| Kejadian | Kas harian | Jumlah |
|----------|-----------|--------|
| Checkout tunai (tanpa saldo) | +✅ | payment_amount |
| Checkout tunai (dengan saldo) | +✅ | payment_amount (= total - balance_used) |
| Checkout transfer/qris/kartu | — | 0 |
| Checkout hutang | — | 0 |
| Checkout full saldo | — | 0 |
| Top-up manual (customer titip di luar transaksi) | +✅ | jumlah top-up |
| Refund/tarik saldo | -❗ | dicatat sebagai expense |

---

## 8. Permission / Role

| Aksi | Owner | Admin | Kasir |
|------|-------|-------|-------|
| Checkout pakai saldo | ✅ | ✅ | ✅ |
| Simpan kembalian ke saldo (di struk) | ✅ | ✅ | ✅ |
| Top-up manual (menu pelanggan) | ✅ | ✅ | ❌ |
| Tarik/Refund saldo | ✅ | ✅ | ❌ |
| Adjustment saldo | ✅ | ❌ | ❌ |
| Lihat riwayat saldo | ✅ | ✅ | ✅ |

---

## 9. Pertanyaan Terbuka

| # | Pertanyaan | Rekomendasi |
|---|-----------|-------------|
| 1 | Transfer saldo antar pelanggan? | Tidak perlu untuk MVP |
| 2 | Expiry saldo? | Tidak perlu |
| 3 | Notifikasi saldo rendah? | Tidak perlu |
| 4 | Cetak kartu saldo? | Nanti, bukan prioritas |

---

## 10. Fase Implementasi

| Fase | Scope | Detail |
|------|-------|--------|
| **1** | Database + BE core | Migration (balance, balance_used, mutations table, payment_method 'balance'). Endpoint top-up, refund, history. Kolom saldo di customer response |
| **2** | FE kasir + checkout | Checkbox "Gunakan Saldo" di modal. Rename Kredit→Hutang. Logic saldo+metode. Opsi simpan kembalian di struk |
| **3** | BE checkout integration | Handle balance_used di transactions/create. Deduct saldo, buat mutasi, piutang = effective. Kas harian = payment_amount |
| **4** | Void + edge cases | Rollback saldo saat void. Rollback top-up saat void. Refund = expense di kas |
| **5** | FE pelanggan + polish | Kolom saldo di list. Detail: section saldo + riwayat + top-up/tarik. Info saldo di struk |
