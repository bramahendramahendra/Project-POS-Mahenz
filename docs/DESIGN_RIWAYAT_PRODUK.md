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

Mengubah modal detail produk yang sekarang (flat info tanpa tab) menjadi **tabbed modal** dengan 4 tab.

### Kondisi Saat Ini (Sebelum)

Modal detail produk berisi info flat dalam 1 halaman scroll:
- Identitas: Nama, Status, Barcode, SKU, Kategori, Satuan
- Harga: Beli, Jual, Margin
- Stok: Saat ini (breakdown per paket), Minimum, Reserved
- Grosiran/Satuan Lain: tabel rasio + harga per paket
- Batch Expired: tabel batch qty + tanggal + status

Semua ditampilkan tanpa tab, langsung scroll ke bawah.

### Kondisi Sesudah (Redesign)

```
Modal Detail Produk → 4 Tab:

  [Detail] [Paket Satuan] [Riwayat Pembelian] [Riwayat Penjualan]
```

### Tab 1: Detail
Isi sama dengan modal sekarang, ditata ulang ke format grid:
- Nama Produk, Status
- Barcode, SKU/Kode
- Kategori, Satuan Dasar
- Harga Beli, Harga Jual, Margin
- Stok Saat Ini (+ visual bar low stock), Stok Minimum
- Warning `needs_stock_review` (jika ada)
- Batch Expired (tabel, jika ada)

### Tab 2: Paket Satuan
Data dari API `products/:id/packages/list` yang sudah ada:
- Daftar semua satuan/kemasan produk (Pack, Slop, Dus, dll)
- Setiap paket: nama satuan, konversi (1 Slop = 10 Pack), harga beli, harga jual, stok per level
- Badge "Satuan Dasar" untuk default package
- Info ini sekarang ada di section "Grosiran / Satuan Lain" — dipindahkan ke tab sendiri agar lebih rapi

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

### 5.2 Query SQL (Verified dari schema DB)

**Riwayat Pembelian** (tabel `purchase_items` JOIN `purchases` JOIN `suppliers`):
```sql
SELECT pi.quantity, pi.purchase_price, pi.unit, pi.subtotal,
       p.purchase_code, p.invoice_number, p.purchase_date,
       s.name AS supplier_name
FROM purchase_items pi
JOIN purchases p ON pi.purchase_id = p.id
JOIN suppliers s ON p.supplier_id = s.id
WHERE pi.product_id = ?
  AND p.status = 'active'
ORDER BY p.purchase_date DESC
LIMIT ? OFFSET ?
```

**Count untuk pagination:**
```sql
SELECT COUNT(*) FROM purchase_items pi
JOIN purchases p ON pi.purchase_id = p.id
WHERE pi.product_id = ? AND p.status = 'active'
```

**Summary (total nota, total qty, total nilai):**
```sql
SELECT COUNT(DISTINCT pi.purchase_id) as total_notes,
       COALESCE(SUM(pi.quantity), 0) as total_qty,
       COALESCE(SUM(pi.subtotal), 0) as total_value
FROM purchase_items pi
JOIN purchases p ON pi.purchase_id = p.id
WHERE pi.product_id = ? AND p.status = 'active'
```

**Riwayat Penjualan** (tabel `transaction_items` JOIN `transactions` LEFT JOIN `customers`):
```sql
SELECT ti.quantity, ti.price, ti.unit, ti.subtotal, ti.discount_item,
       t.transaction_code, t.transaction_date, t.status,
       COALESCE(c.name, '') AS customer_name
FROM transaction_items ti
JOIN transactions t ON ti.transaction_id = t.id
LEFT JOIN customers c ON t.customer_id = c.id
WHERE ti.product_id = ?
ORDER BY t.transaction_date DESC
LIMIT ? OFFSET ?
```

**Dengan filter status (opsional):**
```sql
... AND t.status = ?
```

**Summary:**
```sql
SELECT COUNT(*) as total_transactions,
       COALESCE(SUM(ti.quantity), 0) as total_qty,
       COALESCE(SUM(ti.subtotal), 0) as total_revenue
FROM transaction_items ti
JOIN transactions t ON ti.transaction_id = t.id
WHERE ti.product_id = ? AND t.status = 'completed'
```

### 5.3 Catatan Schema

- `purchase_items.product_id` → FK ke `products(id)` ON DELETE SET NULL. Jika produk dihapus, `product_id` jadi NULL — query harus handle ini.
- `transaction_items.product_id` → sama, ON DELETE SET NULL.
- `purchases` tabel punya kolom `status` ENUM('active','void') — hanya tampilkan yang `active`.
- `purchase_items` tidak punya `product_name` snapshot (beda dari `transaction_items` yang punya). Jadi JOIN ke `products.name` dibutuhkan sebagai fallback display.

### 5.4 FE — Perubahan

| Item | Keterangan |
|------|-----------|
| Refactor `ProductDetailModal.tsx` | Ubah dari flat layout → tabbed (shadcn Tabs component) |
| Tab "Detail" | Isi yang sama seperti sekarang (hanya ditata ulang) |
| Tab "Paket Satuan" | Pindahkan section "Grosiran" ke tab sendiri, tambah info stok per level |
| Tab "Riwayat Pembelian" | Komponen baru: tabel + summary + pagination |
| Tab "Riwayat Penjualan" | Komponen baru: tabel + summary + filter + pagination |
| API hooks baru | `useProductPurchaseHistory(id, page)` + `useProductSaleHistory(id, page, status)` |

### 5.5 Tidak Ada Perubahan di

- Tabel database (query dari tabel yang sudah ada)
- Menu sidebar
- Action buttons di tabel produk (tetap pakai icon Eye yang sudah ada)
- Permission (mengikuti permission `produk.produk` yang sudah ada)

---

## 6. Analisis Gap & Potensi Masalah

| # | Area | Masalah | Solusi |
|---|------|---------|--------|
| 1 | `product_id` bisa NULL | Jika produk dihapus, `purchase_items.product_id` = NULL → record hilang dari riwayat | Query tetap pakai `WHERE pi.product_id = ?` — record lama yang NULL memang sudah tidak relevan |
| 2 | Tabel `purchases` punya status `void` | PO yang di-void tidak boleh muncul di riwayat pembelian | Filter `AND p.status = 'active'` |
| 3 | `purchase_items` tidak punya `product_name` | Jika join ke products gagal (NULL), nama produk tidak tersedia | Gunakan `COALESCE(prod.name, 'Produk Dihapus')` sebagai fallback — tapi ini edge case karena kita query by product_id yang valid |
| 4 | Performa query | Produk populer bisa punya 500+ transaction_items | Index `product_id` sudah ada (FK). Pagination LIMIT/OFFSET cukup. |
| 5 | Modal size | Konten tab riwayat bisa panjang (tabel banyak baris) | Modal pakai ScrollArea, tab content scrollable independent |
| 6 | Existing modal content | Modal sekarang flat tanpa tab — refactor diperlukan | Minimal breaking change: bungkus konten existing di Tab "Detail", tambah tab baru |

### Checklist Keamanan (Tidak Merusak Existing)

| Fitur | Terpengaruh? | Tindakan |
|-------|:------------:|----------|
| Modal detail produk existing | ⚠️ Refactor | Konten dipindah ke Tab "Detail" — isi sama, layout jadi tabbed |
| Action button Eye | ❌ | Tetap sama, trigger modal yang sama |
| Halaman produk (tabel) | ❌ | Tidak berubah |
| Edit produk | ❌ | Modal terpisah, tidak terpengaruh |
| API existing | ❌ | Tidak diubah, hanya tambah endpoint baru |
| Permission | ❌ | Ikut permission `produk.produk` yang sudah ada |

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
