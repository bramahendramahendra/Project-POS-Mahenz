# Desain Fitur: Riwayat Produk (Pembelian & Penjualan)

> Status: **FINAL** — siap implementasi.
> 
> Mockup: `mockup/detail_produk_redesign.html` (buka di browser)
> 
> Catatan: Implementasi mengikuti desain standar FE yang sudah ada (shadcn/ui, warna, spacing, komponen).

---

## 1. Latar Belakang

Saat ini, untuk mengetahui riwayat pembelian atau penjualan suatu produk, user harus membuka menu Pembelian/Transaksi lalu scroll satu-satu mencari produk tersebut di setiap nota. Ini tidak efisien, terutama saat nota sudah banyak.

**Kebutuhan:**
- Dari **1 produk**, bisa langsung lihat: pernah dibeli dari supplier mana, di nota/PO apa, tanggal berapa, harga berapa.
- Dari **1 produk**, bisa langsung lihat: pernah terjual di transaksi apa, tanggal berapa, ke pelanggan siapa, harga berapa.

---

## 2. Desain: Redesign Modal Detail Produk (Tabbed)

Mengubah modal detail produk yang sekarang (flat info) menjadi **tabbed modal** dengan 4 tab:

### Posisi di UI

```
Halaman Produk → Klik icon Eye (action button) → Modal Detail terbuka

SEKARANG: info flat tanpa tab
SESUDAH:  4 tab interactive

  [Detail] [Paket Satuan] [Riwayat Pembelian] [Riwayat Penjualan]
```

### Tab 1: Detail (info yang sudah ada, ditata ulang)
- Nama Produk, Status
- Barcode, SKU/Kode
- Kategori, Satuan Dasar
- Harga Beli, Harga Jual, Margin
- Stok Saat Ini (+ visual bar), Stok Minimum

### Tab 2: Paket Satuan
- Daftar semua satuan/kemasan produk (Pack, Slop, Dus, dll)
- Setiap paket menampilkan: nama satuan, konversi (1 Slop = 10 Pack), harga beli, harga jual, stok per level
- Badge "Satuan Dasar" untuk default package

---

## 3. Tab "Riwayat Pembelian"

Menampilkan daftar semua PO/nota pembelian yang mengandung produk ini.

### Informasi yang Ditampilkan (per baris)

| Kolom | Keterangan |
|-------|-----------|
| Tanggal | Tanggal PO (`purchase_date`) |
| Kode PO | Kode pembelian (`purchase_code`) — klikable, bisa buka detail PO |
| No. Faktur | Nomor invoice supplier |
| Supplier | Nama supplier |
| Satuan | Satuan yang dibeli (Pack, Slop, dll) |
| Qty | Jumlah yang dibeli di PO tersebut |
| Harga Beli | Harga per unit saat itu |
| Subtotal | Qty × Harga |

### Fitur Tambahan
- **Sortable** by tanggal (default: terbaru dulu)
- **Pagination** (jika banyak)
- **Ringkasan** di atas tabel: Total pembelian (berapa kali dibeli, total qty, total nilai)

### Contoh Tampilan

```
┌─────────────────────────────────────────────────────────────────────────┐
│ Riwayat Pembelian — 88 Kretek 12                                        │
│                                                                         │
│ Total: 15 nota pembelian | 450 Pack | Rp 4.500.000                      │
│                                                                         │
│ ┌─────────┬────────────┬──────────┬──────────────┬──────┬───┬────────┐ │
│ │ Tanggal │ Kode PO    │ Faktur   │ Supplier     │Satuan│Qty│ Harga  │ │
│ ├─────────┼────────────┼──────────┼──────────────┼──────┼───┼────────┤ │
│ │ 22 Aug  │ PO-0045    │ INV-112  │ Toko Sejati  │ Pack │ 50│ 10.000 │ │
│ │ 18 Aug  │ PO-0042    │ INV-108  │ Toko Maju    │ Slop │  5│100.000 │ │
│ │ 10 Aug  │ PO-0038    │ INV-095  │ Toko Sejati  │ Pack │ 30│ 10.200 │ │
│ │ ...     │            │          │              │      │   │        │ │
│ └─────────┴────────────┴──────────┴──────────────┴──────┴───┴────────┘ │
│                                                      Page 1 of 3        │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Tab "Riwayat Penjualan"

Menampilkan daftar semua transaksi penjualan yang mengandung produk ini.

### Informasi yang Ditampilkan (per baris)

| Kolom | Keterangan |
|-------|-----------|
| Tanggal | Tanggal transaksi (`transaction_date`) |
| Kode Transaksi | Kode transaksi (`transaction_code`) — klikable |
| Pelanggan | Nama pelanggan (atau "-" jika umum) |
| Satuan | Satuan yang dijual |
| Qty | Jumlah yang dijual |
| Harga Jual | Harga per unit saat itu |
| Subtotal | Qty × Harga |
| Status | completed / void |

### Fitur Tambahan
- **Sortable** by tanggal (default: terbaru dulu)
- **Pagination** (jika banyak)
- **Filter** status (completed / void / semua)
- **Ringkasan** di atas tabel: Total penjualan (berapa kali terjual, total qty, total nilai)
- Transaksi void ditampilkan dengan warna merah/strikethrough

### Contoh Tampilan

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Riwayat Penjualan — 88 Kretek 12                                            │
│                                                                             │
│ Total: 82 transaksi | 245 Pack | Rp 2.695.000                              │
│ [Filter: Semua ▼]                                                           │
│                                                                             │
│ ┌─────────┬────────────────┬──────────────┬──────┬───┬────────┬──────────┐ │
│ │ Tanggal │ Kode Transaksi │ Pelanggan    │Satuan│Qty│ Harga  │ Status   │ │
│ ├─────────┼────────────────┼──────────────┼──────┼───┼────────┼──────────┤ │
│ │ 22 Aug  │ WEB-20260822-5 │ Bu Sari      │ Pack │  2│ 11.000 │completed │ │
│ │ 22 Aug  │ WEB-20260822-3 │ -            │ Pack │  1│ 11.000 │completed │ │
│ │ 21 Aug  │ WEB-20260821-8 │ Pak Ahmad    │ Pack │  5│ 11.000 │  void    │ │
│ │ ...     │                │              │      │   │        │          │ │
│ └─────────┴────────────────┴──────────────┴──────┴───┴────────┴──────────┘ │
│                                                         Page 1 of 9         │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Implementasi Teknis

### 5.1 BE — Endpoint Baru

| Endpoint | Keterangan |
|----------|-----------|
| `POST /products/:id/purchase-history` | Riwayat pembelian produk. Body: `{ page, limit }` |
| `POST /products/:id/sale-history` | Riwayat penjualan produk. Body: `{ page, limit, status? }` |

### 5.2 Query SQL

**Riwayat Pembelian:**
```sql
SELECT pi.quantity, pi.purchase_price, pi.unit,
       sp.purchase_code, sp.invoice_number, sp.purchase_date,
       s.name AS supplier_name
FROM purchase_items pi
JOIN supplier_purchases sp ON pi.purchase_id = sp.id
JOIN suppliers s ON sp.supplier_id = s.id
WHERE pi.product_id = ?
ORDER BY sp.purchase_date DESC
LIMIT ? OFFSET ?
```

**Riwayat Penjualan:**
```sql
SELECT ti.quantity, ti.price, ti.unit, ti.subtotal,
       t.transaction_code, t.transaction_date, t.status,
       COALESCE(c.name, '') AS customer_name
FROM transaction_items ti
JOIN transactions t ON ti.transaction_id = t.id
LEFT JOIN customers c ON t.customer_id = c.id
WHERE ti.product_id = ?
ORDER BY t.transaction_date DESC
LIMIT ? OFFSET ?
```

### 5.3 FE — Perubahan

| Item | Keterangan |
|------|-----------|
| Buat komponen `ProductPurchaseHistory.tsx` | Tab riwayat pembelian (tabel + pagination) |
| Buat komponen `ProductSaleHistory.tsx` | Tab riwayat penjualan (tabel + pagination + filter status) |
| Integrasikan ke detail produk | Tambah 2 tab baru di modal/panel detail produk |
| API hooks baru | `useProductPurchaseHistory(id)` + `useProductSaleHistory(id)` |

### 5.4 Tidak Ada Perubahan di

- Tabel database (query dari tabel yang sudah ada: `purchase_items`, `transaction_items`)
- Menu sidebar (tidak perlu menu baru)
- Fitur existing (tidak ada modifikasi logic, hanya tambah endpoint read-only)

---

## 6. Pertimbangan

| Aspek | Keputusan | Alasan |
|-------|-----------|--------|
| Posisi fitur | Tab di detail produk | Tidak menambah menu, akses cepat |
| Akses/role | Sama dengan akses halaman Produk | Jika bisa lihat produk, bisa lihat riwayatnya |
| Performa | Pagination (10-20 per page) | Produk populer bisa punya ratusan record |
| Klik kode PO/Transaksi | Navigasi ke detail PO/Transaksi | Biar bisa lihat nota lengkap |

---

## 7. Keputusan Final

| # | Keputusan | Catatan |
|---|-----------|---------|
| 1 | Redesign modal detail produk jadi **tabbed** | 4 tab: Detail, Paket Satuan, Riwayat Pembelian, Riwayat Penjualan |
| 2 | Tidak ada menu sidebar baru | Cukup dari modal existing yang di-upgrade |
| 3 | Implementasi mengikuti **desain standar FE** | shadcn/ui components, warna, spacing |
| 4 | Akses mengikuti permission halaman Produk | Siapa bisa buka produk, bisa lihat riwayat |
| 5 | Kode PO / Transaksi bisa diklik | Navigasi ke detail nota |
| 6 | Pagination 10 per halaman | Produk populer bisa ratusan record |
| 7 | Riwayat Penjualan bisa filter status | Semua / Completed / Void |
| 8 | Tidak perlu migration DB baru | Query dari tabel existing (purchase_items + transaction_items) |
