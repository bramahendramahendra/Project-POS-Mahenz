# Rancangan Perbaikan Menu Piutang

> Status: **RANCANGAN — belum ada implementasi.**
> Tujuan: papan kerja bersama untuk merapikan & menambah fitur menu Piutang, dibahas dulu
> sampai disepakati sebelum menulis kode. Ditulis dengan bahasa sederhana + rincian teknis
> per file agar bisa direview.

---

## 1. Latar Belakang

Menu Piutang (`/receivables`) saat ini hanya menampilkan daftar dengan kolom: Kode Transaksi,
Pelanggan, Total Piutang, Sudah Dibayar, Sisa, Status, Jatuh Tempo, dan satu tombol **Bayar**.
Tidak ada cara melihat **rincian barang** dari piutang, **riwayat cicilan**, atau **cetak nota**.

Kebutuhan owner (disepakati): dari sebuah baris piutang harus bisa melihat detail lengkap
(barang apa yang dibeli + sudah nyicil berapa kali) dan mencetak nota. Ditambah beberapa fitur
pendukung penagihan.

---

## 2. Temuan Analisis — Inkonsistensi FE ↔ BE (WAJIB dirapikan dulu)

Ditemukan saat menelusuri kode (baca `receivables.types.ts`, `dto_receivable.go`,
`receivable_repo.go`, `receivable_handler.go`, `ReceivableTable*.tsx`):

1. **FE mendeklarasikan field yang BE tidak kirim.** Type FE `Receivable` punya
   `transaction_id`, `customer_id`, `payments[]`, `created_at`. BE (`ReceivableResponse` &
   `ReceivableDetailResponse`) TIDAK mengirim satu pun dari itu → saat runtime nilainya
   `undefined`. TypeScript tak protes karena tipe terlanjur dideklarasikan.
2. **Endpoint detail = daftar.** `ReceivableResponse` dan `ReceivableDetailResponse` isinya
   IDENTIK. Hook `useReceivableDetailQuery` sudah ada tapi belum dipakai di UI mana pun.
3. **Riwayat cicilan: ada di DB & ada endpoint, tapi FE belum menyambung.** Endpoint
   `POST /receivables/payments/:id` mengembalikan cicilan lengkap, tapi tidak ada hook FE
   yang memanggilnya.
4. **Bentuk `payments` beda FE vs BE.** FE `ReceivablePayment` = {id, amount, payment_date,
   notes}. BE `PaymentResponse` = {id, payment_date, amount, payment_method, user_name,
   notes}. Perlu disamakan (punya BE lebih lengkap).

**Bahan yang sudah tersedia & bisa dipakai ulang (tidak bikin dari nol):**
- Setiap `receivables` punya kolom `transaction_id` (FK ke `transactions`) → rincian barang bisa ditarik.
- Endpoint `POST /transactions/detail/:id` sudah mengembalikan transaksi + `items[]`.
- Komponen `ReceiptPrint` (`FE/src/features/sales/cashier/components/ReceiptPrint.tsx`)
  sudah reusable — dipakai di `TransactionDetailModal` mode `reprint`. Butuh data:
  `checkoutData`, `cart[]`, `summary`, `discount`, `tax`, `paymentMethod`, `amountPaid`,
  `customerName`, `mode`.
- Endpoint `POST /receivables/summary` (rekap per pelanggan) sudah ada tapi belum dipakai FE.
- `queryKeys.receivables` sudah ada (`all`, `list`, `detail`) — tinggal tambah `payments`.
- Kolom indikator overdue `⚠` sudah ada di `ReceivableTableColumns.tsx`.

---

## 3. Prinsip Desain

1. **Rapikan fondasi dulu, baru tambah fitur.** Samakan FE↔BE sebelum menumpuk fitur di
   atas data yang belum konsisten.
2. **Pakai ulang yang sudah ada.** Detail piutang meniru pola `TransactionDetailModal`;
   cetak nota pakai `ReceiptPrint`.
3. **Jangan ubah cara data dibuat.** Piutang tetap lahir dari transaksi kredit
   (`is_credit`), tidak ada input piutang manual, tidak ada bunga/denda.
4. **Bahasa awam untuk owner** (mis. "Sisa Piutang", "Bukti Pembayaran"), istilah teknis
   dihindari di UI.
5. **Konsistensi lintas menu.** Kartu ringkasan & filter meniru pola menu Stok/Penjualan
   yang sudah kita kerjakan.

---

## 4. Ruang Lingkup — Dikerjakan vs TIDAK

**Dikerjakan (disepakati):**
- Detail Piutang lengkap (ringkasan + rincian barang + riwayat cicilan) + tombol Cetak Nota
- Bukti pembayaran tiap cicilan (struk kecil)
- Kartu ringkasan di atas daftar
- Filter "jatuh tempo"

**TIDAK dikerjakan (sengaja, biar tidak berlebihan):**
- Pengingat WhatsApp/SMS otomatis
- Edit/hapus piutang manual
- Bunga/denda keterlambatan
- Ringkasan per-pelanggan sebagai tampilan utama (endpoint `summary` boleh dipakai untuk
  kartu ringkasan agregat saja, bukan tabel per pelanggan) — bisa dibahas terpisah nanti.

---

## 5. Rencana Bertahap

### TAHAP 0 — Rapikan fondasi (samakan FE ↔ BE)

Tujuan: FE & BE bicara data yang sama. Tanpa ini, fitur baru menampilkan kosong.

**Backend:**
- `dto_receivable.go`:
  - `ReceivableResponse` (daftar): tambah `transaction_id` + `created_at`. (agar FE bisa
    buka detail transaksi & tahu tanggal piutang)
  - Rombak `ReceivableDetailResponse` jadi benar-benar "detail": tambah `transaction_id`,
    `customer_id`, `created_at`, `notes`, dan **sisipkan** `items []TransactionItemBrief`
    + `payments []PaymentResponse`. (satu panggilan dapat semua)
  - Tambah tipe ringkas `TransactionItemBrief` {product_name, unit, quantity, price,
    subtotal, discount_item} — cukup untuk tampilan & nota.
- `receivable_repo.go`:
  - `getAllReceivablesQuery` & `getReceivableDetailQuery`: SELECT tambah `r.transaction_id`,
    `r.created_at`, `r.notes`.
  - `GetDetailByID`: setelah ambil header, ambil juga item transaksi (JOIN ke
    `transaction_items` via `transaction_id`) + panggil `GetPayments`. Gabung ke response.
- `dto`/`repo`: pastikan status enum `void` ikut ditangani (DB enum hanya
  'unpaid','partial','paid' — `void` di-set lewat update saat transaksi induk di-void;
  perlu dicek apakah kolom status perlu diperlebar. **CATATAN: perlu verifikasi migrasi.**)

**Frontend:**
- `receivables.types.ts`:
  - Samakan `ReceivablePayment` dengan BE: tambah `payment_method`, `user_name`.
  - `Receivable`: pertahankan field, pastikan `transaction_id`/`created_at` kini benar terisi.
  - Tambah `ReceivableDetail` (extends Receivable) dengan `items[]` + `payments[]`.
  - Tambah tipe `ReceivableItem` {product_name, unit, quantity, price, subtotal, discount_item}.
- `receivables.api.ts`:
  - Perbaiki `useReceivableDetailQuery` agar mengembalikan `ReceivableDetail`.
  - Tambah `useReceivablePaymentsQuery(id)` → `POST /receivables/payments/:id` (kalau
    memilih ambil terpisah; atau cukup ikut dari detail bila digabung — **diputuskan:**
    gabung di detail agar 1 request).
- `queryKeys.ts`: tambah `receivables.payments(id)`.

**Verifikasi Tahap 0:** `go build ./...` + `npm run type-check` + `npm run lint` bersih;
cek via browser daftar piutang masih tampil normal (tidak ada regresi).

---

### TAHAP 1 — Detail Piutang + Cetak Nota (kebutuhan utama)

**Frontend (baru):**
- `components/ReceivableDetailModal.tsx` (baru) — meniru pola `TransactionDetailModal`:
  - **Bagian 1 Ringkasan:** pelanggan, kode transaksi, tanggal, jatuh tempo, status,
    total / sudah dibayar / sisa.
  - **Bagian 2 Rincian Barang:** tabel item (produk, qty, harga, subtotal) dari
    `detail.items`.
  - **Bagian 3 Riwayat Cicilan:** daftar `detail.payments` (tanggal, jumlah, metode,
    petugas, catatan). Kalau kosong → "Belum ada pembayaran".
  - **Footer:** tombol **Cetak Nota** (buka `ReceiptPrint` mode `reprint`, data dibangun
    dari `detail.items` seperti `buildReceiptData`), dan tombol **Bayar** (buka
    `PaymentRecordModal` yang sudah ada).
- `ReceivableTableColumns.tsx`: tambah aksi **Detail** (ikon mata) di kolom Aksi, di samping
  tombol Bayar. Untuk semua status (termasuk lunas/void — biar tetap bisa lihat riwayat).
- `ReceivableTable.tsx`: tambah state `detailTarget` + render `ReceivableDetailModal`.

**Backend:** tidak ada tambahan (sudah disiapkan di Tahap 0).

**Verifikasi Tahap 1:** buka detail dari daftar → 3 bagian tampil benar; klik Cetak Nota →
struk muncul dengan item benar; build/type-check/lint bersih; screenshot.

---

### TAHAP 2 — Fitur pendukung

**2a. Bukti Pembayaran cicilan (struk kecil).**
- Setelah `PaymentRecordModal` sukses catat cicilan, tawarkan cetak "Bukti Pembayaran"
  berisi: nama toko, pelanggan, tanggal, jumlah dibayar, sisa piutang setelah bayar,
  petugas. Pakai `ReceiptPrint` dengan varian teks, atau komponen struk ringkas baru
  `PaymentReceiptPrint` (diputuskan saat implementasi — utamakan pakai ulang).
- Data sisa setelah bayar sudah dikembalikan `PayResponse` (remaining_amount).

**2b. Kartu ringkasan di atas daftar.**
- Komponen `ReceivableSummaryCards.tsx` (baru): **Total Sisa Piutang**, **Jumlah Pelanggan
  Berutang**, **Jatuh Tempo Terlewat** (nominal/among). Sumber: bisa dari endpoint
  `summary` (rekap per pelanggan, tinggal diagregasi) atau endpoint ringkas baru. Pola
  visual meniru kartu menu Stok.
- **CATATAN:** perlu putuskan apakah "jatuh tempo terlewat" dihitung di BE (query) atau FE.
  Rekomendasi: BE, biar akurat lintas halaman (bukan cuma halaman aktif).

**2c. Filter "jatuh tempo".**
- `ReceivableFilterBar.tsx`: tambah pilihan filter selain status: "Semua / Belum jatuh
  tempo / Sudah jatuh tempo (overdue)". Backend `GetAll` (`getAllReceivablesQuery`) tambah
  kondisi `r.due_date < CURDATE() AND r.status NOT IN ('paid','void')` untuk overdue.
- Konsisten dengan indikator `⚠` yang sudah ada di kolom.

**Verifikasi Tahap 2:** tiap sub-fitur dites via browser + build/type-check/lint bersih.

---

## 6. Hal yang Perlu Diputuskan / Diverifikasi Sebelum Koding

1. **Status `void` di kolom DB.** Enum `receivables.status` di migrasi awal hanya
   'unpaid','partial','paid'. Kode meng-set 'void' saat transaksi di-void. Perlu cek:
   apakah sudah ada migrasi yang memperlebar enum, atau kolom sebenarnya VARCHAR. Kalau
   belum, siapkan migrasi kecil. (Akan diverifikasi saat mulai Tahap 0.)
2. **Ambil payments: gabung di detail vs endpoint terpisah.** Rekomendasi: **gabung di
   detail** (1 request, lebih simpel untuk modal). Endpoint `/payments/:id` tetap
   dipertahankan (tidak dihapus) untuk kompatibilitas.
3. **Bukti pembayaran: pakai ulang `ReceiptPrint` atau komponen baru.** Diputuskan saat
   implementasi; prioritas pakai ulang.
4. **Sumber kartu ringkasan (2b):** endpoint `summary` yang ada vs endpoint agregat baru.
5. **Hak akses:** semua fitur baca (detail, ringkasan) pakai permission `pelanggan.piutang`
   `can_view`; cetak mengikuti `can_view`; pembayaran tetap `can_edit` (sudah ada).

---

## 7. Aturan Kerja (sama seperti sesi sebelumnya)

- Tidak ada perubahan kode sampai rancangan ini disepakati.
- Kerjakan bertahap; tiap tahap diverifikasi (build BE, type-check + lint FE, tes browser)
  sebelum lanjut.
- Testing browser via Playwright dari folder `testing`; cache di folder testing (bukan disk C).
- Data DB saat ini = data production; koreksi/rollback dibicarakan dulu dengan user.
- Setelah selesai: siapkan catatan deploy (perubahan ini FE + BE biasa; bila ada migrasi
  enum status, catat sebagai langkah DB).
