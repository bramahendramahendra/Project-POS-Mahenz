# Desain Fitur: Saldo Pelanggan (Customer Deposit)

> Dokumen ini mencatat desain lengkap fitur Saldo Pelanggan sebelum implementasi.
> Status: **DISKUSI** — belum implementasi.

---

## 1. Latar Belakang

### Use Case Real

Customer A datang ke toko, beli air minum Rp 5.000 tapi bayar Rp 200.000. Sisa Rp 195.000 disimpan di toko sebagai "deposit/titipan" untuk belanja berikutnya tanpa perlu bayar lagi setiap kali datang.

### Kondisi Existing Saat Ini

| Fitur | Status | Catatan |
|-------|--------|---------|
| Checkout Tunai | ✅ Ada | Kas harian bertambah |
| Checkout Transfer/QRIS/Kartu | ✅ Ada | Kas harian TIDAK bertambah (uang masuk rekening) |
| Checkout Kredit (hutang pelanggan) | ✅ Ada | Kas TIDAK bertambah, piutang terbuat otomatis |
| Deposit/Saldo Pelanggan | ❌ Belum ada | — |

### Existing: Bagaimana Kas Harian Bekerja

Dari kode `transaction_service.go`:
```go
if req.PaymentMethod == "cash" {
    cashDrawerRepo.UpdateSales(drawer.ID, req.TotalAmount, req.TotalAmount, time)
}
```

**Hanya metode "cash" yang menambah kas harian.** Transfer, QRIS, Kartu, Kredit — tidak.

---

## 2. Keputusan Desain yang Sudah Disepakati

| # | Keputusan | Alasan |
|---|-----------|--------|
| 1 | **Opsi B: Kas harian = uang fisik** | Kas selalu cocok dengan isi laci saat tutup shift |
| 2 | Saldo pakai saldo → kas TIDAK bertambah | Uang sudah masuk kas saat top-up |
| 3 | Top-up saldo → kas BERTAMBAH | Uang fisik masuk laci saat customer titip |
| 4 | Kombinasi saldo + tunai dibolehkan | Saldo < total → sisa bayar tunai |
| 5 | Kombinasi saldo + hutang dibolehkan | Saldo < total → sisa jadi piutang |
| 6 | Default credit limit pelanggan = 0 (bukan tak terbatas) | Mencegah hutang tanpa batas |
| 7 | Rename "Kredit" → "Hutang" | Lebih jelas untuk kasir retail |
| 8 | "Gunakan Saldo" = checkbox, bukan metode pembayaran | Saldo dipotong dulu, sisa baru pilih metode |
| 9 | "Gunakan Saldo" HANYA muncul jika pelanggan dipilih | Tanpa pelanggan = UI seperti existing |
| 10 | "Hutang" HANYA muncul jika pelanggan dipilih | Sama seperti "Kredit" existing |
| 11 | Tanpa pelanggan → modal 100% seperti sekarang | Tidak ada perubahan visual jika tidak centang "Tambah Pelanggan" |
| 12 | "Simpan kembalian ke saldo" muncul SETELAH checkout di halaman struk | Bukan di modal pembayaran |

---

## 3. Desain Database

### Tambah kolom di `customers`

```sql
ALTER TABLE customers ADD COLUMN balance DECIMAL(15,2) NOT NULL DEFAULT 0;
```

### Tambah kolom di `transactions` (untuk track saldo yang dipakai)

```sql
ALTER TABLE transactions ADD COLUMN balance_used DECIMAL(15,2) NOT NULL DEFAULT 0;
```

Kolom ini mencatat berapa saldo pelanggan yang digunakan dalam transaksi ini.
- `balance_used = 0` → tidak pakai saldo (default, existing transactions tetap valid)
- `balance_used > 0` → sebagian/seluruh dibayar dari saldo

Contoh penyimpanan per skenario:

| Skenario | payment_method | total_amount | payment_amount | balance_used | change_amount | is_credit | Piutang |
|----------|---------------|-------------|---------------|-------------|--------------|-----------|---------|
| Saldo cukup | balance | 15.000 | 0 | 15.000 | 0 | false | — |
| Saldo + Tunai | cash | 50.000 | 20.000 | 30.000 | 0 | false | — |
| Saldo + Hutang | kredit | 30.000 | 0 | 10.000 | 0 | true | 20.000 |
| Tunai biasa | cash | 22.500 | 25.000 | 0 | 2.500 | false | — |
| Hutang penuh | kredit | 30.000 | 0 | 0 | 0 | true | 30.000 |

**Rumus piutang saat is_credit:**
```
piutang = total_amount - balance_used
```

### Tabel baru: `customer_balance_mutations`

```sql
CREATE TABLE IF NOT EXISTS customer_balance_mutations (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    customer_id     INT           NOT NULL,
    amount          DECIMAL(15,2) NOT NULL,  -- +positif = masuk, -negatif = keluar
    balance_after   DECIMAL(15,2) NOT NULL,  -- saldo setelah mutasi ini
    type            ENUM('topup', 'usage', 'refund', 'adjustment') NOT NULL,
    reference_type  VARCHAR(50)   NULL,      -- 'transaction', 'manual'
    reference_id    INT           NULL,      -- ID transaksi/piutang terkait
    notes           TEXT          NULL,
    user_id         INT           NOT NULL,
    created_at      DATETIME      DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    INDEX idx_balance_mut_customer (customer_id),
    INDEX idx_balance_mut_type (type),
    INDEX idx_balance_mut_ref (reference_type, reference_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Tipe Mutasi

| Type | Deskripsi | Amount | Contoh |
|------|-----------|--------|--------|
| `topup` | Customer setor/titip uang | + | Titip Rp 195.000 dari kembalian |
| `usage` | Pakai saldo untuk belanja | - | Bayar belanja Rp 15.000 dari saldo |
| `refund` | Pengembalian saldo ke customer | - | Customer minta uang balik Rp 50.000 |
| `adjustment` | Koreksi manual admin | +/- | Fix saldo salah |

---

## 4. Alur & Skenario

### Skenario A: Customer Titip Uang (Top-up dari Kasir)

```
Customer A beli air Rp 5.000, bayar tunai Rp 200.000
                    │
                    ▼
┌─────────────────────────────────────────────┐
│ Checkout tunai: Rp 200.000                  │
│ Kembalian: Rp 195.000                       │
│                                             │
│ ┌─────────────────────────────────────────┐ │
│ │ Simpan kembalian ke Saldo Pelanggan?    │ │
│ │                                         │ │
│ │ [Input: Rp 195.000] (default: penuh)    │ │
│ │                                         │ │
│ │ [Tidak, kembalian tunai] [Ya, Simpan]   │ │
│ └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
                    │ "Ya, Simpan"
                    ▼
Efek:
- Transaksi tersimpan (tunai Rp 200.000, kembalian Rp 195.000)
- Kas harian: +Rp 200.000 (semua uang fisik masuk)
- Saldo customer: +Rp 195.000
- Mutasi: type=topup, ref=transaction/ID
- Struk: "Saldo ditambahkan: Rp 195.000"
```

**Catatan penting**: Kas harian tetap +Rp 200.000 (bukan cuma +Rp 5.000), karena uang fisik Rp 200.000 memang masuk ke laci. Saldo customer hanya pencatatan "ini uangnya dia yang kita simpan".

### Skenario B: Customer Belanja Pakai Saldo (Saldo Cukup)

```
Customer A beli air Rp 15.000
Saldo saat ini: Rp 195.000
                    │
                    ▼
┌──────────────────────────────────────────────────┐
│ Modal Pembayaran                                 │
│                                                  │
│ ☑ Gunakan Saldo              Rp 195.000          │
│   Saldo terpakai: Rp 15.000                      │
│   Sisa saldo setelah: Rp 180.000                 │
│   Sisa bayar: Rp 0                               │
│                                                  │
│ Status: Lunas via Saldo ✓                        │
│                                                  │
│ [Proses Bayar]                                   │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi tersimpan (payment_method: 'balance')
- Kas harian: TIDAK bertambah (uang sudah di laci sejak top-up)
- Saldo customer: -Rp 15.000 → sisa Rp 180.000
- Mutasi: type=usage, amount=-15000, ref=transaction/ID
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
│ ☑ Gunakan Saldo              Rp 30.000           │
│   Saldo terpakai: Rp 30.000 (habis)              │
│   Sisa bayar: Rp 20.000                          │
│                                                  │
│ Metode Pembayaran (sisa Rp 20.000):              │
│ [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]       │
│                                                  │
│ Jumlah Bayar: [Rp 20.000]                        │
│ Kembalian: Rp 0                                  │
│                                                  │
│ [Proses Bayar]                                   │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi tersimpan (payment_method: 'cash', balance_used: 30000)
- Kas harian: +Rp 20.000 (hanya porsi tunai yang masuk laci)
- Saldo customer: -Rp 30.000 → sisa Rp 0
- Mutasi saldo: type=usage, amount=-30000, ref=transaction/ID
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
│ ☑ Gunakan Saldo              Rp 10.000           │
│   Saldo terpakai: Rp 10.000 (habis)              │
│   Sisa bayar: Rp 20.000                          │
│                                                  │
│ Metode Pembayaran (sisa Rp 20.000):              │
│ [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]       │
│                              pilih → [Hutang]     │
│                                                  │
│ ⚠️ Piutang: Rp 20.000                            │
│                                                  │
│ [Proses Bayar]                                   │
└──────────────────────────────────────────────────┘
                    │
                    ▼
Efek:
- Transaksi tersimpan (payment_method: 'kredit', balance_used: 10000, is_credit: true)
- Kas harian: TIDAK bertambah
- Saldo customer: -Rp 10.000 → sisa Rp 0
- Piutang terbuat: Rp 20.000 (BUKAN Rp 30.000 — hanya sisa)
- Mutasi saldo: type=usage, amount=-10000, ref=transaction/ID
```

### Skenario D: Top-up Manual (tanpa transaksi)

```
Admin buka menu Pelanggan → Detail → "Top-up Saldo"
                    │
                    ▼
Input: Rp 100.000, catatan: "Customer titip untuk beli air"
                    │
                    ▼
Efek:
- Saldo customer: +Rp 100.000
- Mutasi: type=topup, ref=manual
- Kas harian: +Rp 100.000 (uang fisik masuk laci)
```

### Skenario E: Pengembalian Saldo (Refund)

```
Admin buka menu Pelanggan → Detail → "Tarik Saldo"
                    │
                    ▼
Input: Rp 50.000, catatan: "Customer minta uang balik"
                    │
                    ▼
Efek:
- Saldo customer: -Rp 50.000
- Mutasi: type=refund
- Kas harian: -Rp 50.000 (uang fisik keluar dari laci)
  → Ini HARUS dicatat sebagai pengeluaran kas agar rekonsiliasi tutup shift cocok
```

---

## 5. Potensi Celah/Bug yang Harus Dicegah

| # | Celah | Solusi |
|---|-------|--------|
| 1 | Double-count kas: top-up masuk kas, pakai saldo masuk kas lagi | Saldo usage → kas TIDAK bertambah |
| 2 | Saldo negatif (pakai lebih dari yang ada) | Validasi BE: `balance >= balance_used` sebelum deduct, pakai `FOR UPDATE` |
| 3 | Race condition: 2 kasir proses saldo bersamaan | `SELECT ... FOR UPDATE` pada row customer saat deduct saldo |
| 4 | Refund lebih dari saldo | Validasi: refund amount ≤ balance |
| 5 | Top-up dari kembalian tapi customer belum dipilih | Opsi "Simpan ke Saldo" hanya muncul jika customer dipilih + kembalian > 0 |
| 6 | Void transaksi yang pakai saldo → saldo harus dikembalikan | Void handler: `balance += balance_used` + buat mutasi type=refund |
| 7 | Void transaksi yang top-up saldo → saldo harus dikurangi | Void handler: cari mutasi top-up terkait transaction_id, rollback balance |
| 8 | Rekonsiliasi tutup kas tidak cocok jika refund manual terjadi | Refund dicatat sebagai pengeluaran kas (expense) agar kas harian balance |
| 9 | Laporan laba rugi: transaksi saldo tetap dihitung sbg revenue | Ya — laporan penjualan mencatat SEMUA transaksi regardless metode |
| 10 | Pelanggan dihapus tapi masih punya saldo | Validasi: tidak bisa hapus/nonaktifkan pelanggan dengan balance > 0 |
| 11 | Piutang saldo+hutang: void harus rollback KEDUANYA | Void: kembalikan saldo (balance_used) DAN void piutang (sisa) |
| 12 | balance_used > total_amount (manipulasi API) | Validasi: balance_used ≤ total_amount |
| 13 | Saldo + Hutang tanpa pelanggan | Validasi: is_credit=true ATAU balance_used>0 → wajib customer_id |
| 14 | Top-up di halaman struk: kasir klik "Ya" 2x cepat | Disable button setelah klik pertama (FE) + idempotency check (BE) |

---

## 6. Perubahan UI

### A. Modal Pembayaran (Kasir)

```
SEBELUM:
[Tunai] [Transfer] [QRIS] [Kartu] [Kredit]

SESUDAH (tanpa pelanggan):
[Tunai] [Transfer] [QRIS] [Kartu]
→ Tidak ada "Hutang", tidak ada "Gunakan Saldo"
→ Persis seperti sekarang

SESUDAH (pelanggan dipilih):
┌────────────────────────────────────────┐
│ ☑ Gunakan Saldo          Rp 195.000   │  ← muncul jika saldo > 0
│   Saldo terpakai: ...                  │
│   Sisa bayar: ...                      │
└────────────────────────────────────────┘
Metode: [Tunai] [Transfer] [QRIS] [Kartu] [Hutang]
                                           ↑rename dari "Kredit"
```

Aturan visibilitas:
- **Checkbox "Gunakan Saldo"**: hanya tampil jika pelanggan dipilih DAN saldo > 0
- **Tombol "Hutang"**: hanya tampil jika pelanggan dipilih
- **Tanpa pelanggan**: modal persis seperti existing (Tunai/Transfer/QRIS/Kartu saja)

### B. After Checkout (jika ada kembalian + customer dipilih)

Tampilkan opsi **sebelum struk ditutup**:
```
Kembalian: Rp 195.000
┌──────────────────────────────────────────┐
│ Simpan ke Saldo Pelanggan?               │
│ [Rp 195.000        ] (bisa diedit)       │
│ [Tidak] [Ya, Simpan]                     │
└──────────────────────────────────────────┘
```

### C. Halaman Pelanggan (List)

Tambah kolom **"Saldo"** di tabel.

### D. Detail Pelanggan

Tab/section baru: **"Saldo & Riwayat"**
- Saldo saat ini (prominent)
- Tombol: "Top-up" + "Tarik Saldo"
- Tabel riwayat mutasi (tanggal, tipe, jumlah, saldo setelah, catatan, oleh siapa)

### E. Struk

Jika bayar pakai saldo:
```
Pembayaran: Saldo
Sisa Saldo Anda: Rp 180.000
```

---

## 7. Perubahan Backend

### Endpoint Baru

| Method | Path | Fungsi |
|--------|------|--------|
| POST | `/customers/:id/balance/topup` | Top-up manual (admin/owner) |
| POST | `/customers/:id/balance/refund` | Tarik/refund saldo (admin/owner) |
| POST | `/customers/:id/balance/history` | Riwayat mutasi saldo |
| POST | `/transactions/save-to-balance` | Simpan kembalian ke saldo (setelah checkout) |

### Perubahan Endpoint Existing

| Endpoint | Perubahan |
|----------|-----------|
| `POST /transactions/create` | Tambah field `balance_used` di request. Jika > 0: deduct saldo customer, buat mutasi, hitung piutang = total - balance_used |
| `POST /transactions/void/:id` | Jika `balance_used > 0`: kembalikan saldo ke customer + buat mutasi refund. Jika `is_credit`: void piutang (sudah existing) |
| `POST /customers/detail/:id` | Tambah field `balance` di response |
| `POST /customers/delete/:id` | Tambah validasi: tolak jika balance > 0 |

### Logic di `transaction_service.go` Create() — Pseudocode

```go
func Create(req) {
    // ... existing validasi ...

    // Hitung efektif
    effectiveTotal := req.TotalAmount - req.BalanceUsed
    
    // Validasi saldo
    if req.BalanceUsed > 0 {
        if req.CustomerID == nil { return error("Pilih pelanggan untuk gunakan saldo") }
        customer := getCustomer(req.CustomerID)
        if customer.Balance < req.BalanceUsed { return error("Saldo tidak cukup") }
    }
    
    // Validasi pembayaran (hanya untuk sisa setelah saldo)
    if !req.IsCredit && req.PaymentAmount < effectiveTotal {
        return error("Jumlah pembayaran kurang")
    }

    // Dalam DB transaction:
    tx {
        // 1. Insert transaksi (dengan balance_used)
        // 2. Insert items + kurangi stok
        // 3. Jika balance_used > 0: deduct customer.balance + insert mutasi
        // 4. Jika is_credit: insert piutang (amount = effectiveTotal, bukan total_amount)
        // 5. Jika payment_method == "cash": update kas harian (amount = payment_amount, bukan total)
    }
}
```

### Kas Harian — Aturan Lengkap

| Metode | Kas bertambah? | Jumlah yang ditambahkan |
|--------|---------------|------------------------|
| cash (tanpa saldo) | ✅ | total_amount |
| cash (dengan saldo) | ✅ | payment_amount (= total - balance_used) |
| transfer/qris/kartu | ❌ | — |
| kredit/hutang | ❌ | — |
| balance (saldo cukup) | ❌ | — |
| top-up saldo (dari kasir/manual) | ✅ | jumlah top-up |
| refund saldo | ❌ (-) | dicatat sebagai pengeluaran |

---

## 8. Permission / Role

| Aksi | Owner | Admin | Kasir |
|------|-------|-------|-------|
| Top-up dari kembalian (kasir) | ✅ | ✅ | ✅ |
| Top-up manual | ✅ | ✅ | ❌ |
| Tarik/Refund saldo | ✅ | ✅ | ❌ |
| Adjustment | ✅ | ❌ | ❌ |
| Lihat riwayat saldo | ✅ | ✅ | ✅ (hanya lihat) |

---

## 9. Pertanyaan Terbuka (Belum Diputuskan)

1. **Apakah perlu fitur "transfer saldo antar pelanggan"?** (Edge case: customer A mau kasih saldo-nya ke customer B)
   - Rekomendasi: Tidak perlu untuk MVP

2. **Expiry saldo?** (Saldo hangus setelah X bulan tidak aktif)
   - Rekomendasi: Tidak perlu, tidak umum di toko retail

3. **Notifikasi saldo rendah?** (Saat saldo hampir habis)
   - Rekomendasi: Tidak perlu, customer yang tahu sendiri

4. **Cetak kartu saldo / member card?**
   - Rekomendasi: Bisa ditambah nanti, bukan prioritas

---

## 10. Prioritas Implementasi

| Fase | Scope |
|------|-------|
| **Fase 1** | DB migration + backend top-up/deduct/history + kolom saldo di pelanggan |
| **Fase 2** | Metode "Saldo" di kasir + opsi "Simpan ke Saldo" setelah checkout |
| **Fase 3** | Void handler (rollback saldo) + refund manual + adjustment |
| **Fase 4** | Info saldo di struk + rekonsiliasi refund di kas harian |
