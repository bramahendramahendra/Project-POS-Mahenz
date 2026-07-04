# Audit Konsistensi & Keamanan Kode — FE/BE POS

Tanggal audit: 2026-07-05
Metode: analisis manual per modul (FE React/TS + BE Go), fokus pada konsistensi pola antar fitur sejenis (sibling), korektnes logika bisnis (terutama uang & stok), dan celah otorisasi.

## Daftar File Laporan

| # | File | Modul | Prioritas Tertinggi |
|---|------|-------|---------------------|
| 1 | [01-products.md](01-products.md) | Products (categories/products/units) | Medium |
| 2 | [02-customers.md](02-customers.md) | Customers (customers/receivables) | **Kritis** — fitur bayar piutang rusak |
| 3 | [03-finance.md](03-finance.md) | Finance (cash-drawer/expenses/my-cash) | Tinggi |
| 4 | [04-procurement.md](04-procurement.md) | Procurement (purchases/returns/suppliers) | Tinggi |
| 5 | [05-operational.md](05-operational.md) | Operational (shifts/sync) | Tinggi |
| 6 | [06-sales.md](06-sales.md) | Sales (cashier/transactions) | **Kritis** — checkout tidak aman |
| 7 | [07-reporting.md](07-reporting.md) | Reporting (4 jenis laporan) | Tinggi |
| 8 | [08-settings.md](08-settings.md) | Settings (roles/users/menus/dll) | **Kritis** — privilege escalation |
| 9 | [09-auth.md](09-auth.md) | Auth (login/token/session) | **Kritis** — refresh token rusak, secret placeholder |

## Cara Pakai Laporan Ini

Setiap file laporan berisi tabel temuan dengan kolom:
- **File & Baris** — lokasi kode
- **Temuan** — apa yang salah/tidak konsisten
- **Skenario Kegagalan** — kondisi konkret yang memicu bug
- **Kategori** — security / correctness / consistency
- **Severity** — Critical / High / Medium / Low

Untuk eksekusi perbaikan, gunakan [IMPLEMENTATION-PROMPTS.md](IMPLEMENTATION-PROMPTS.md) — berisi prompt siap pakai per fase, karena jumlah temuan terlalu banyak untuk dikerjakan sekaligus.

## Ringkasan Temuan Lintas Modul (Pola Berulang)

1. **`Update...Payload` di-duplikasi manual** alih-alih `Partial<Create...Payload>` — ditemukan di products, units, expenses, purchases, suppliers, roles, users. Risiko: field baru di Create tidak otomatis ikut ke Update.
2. **Kolom `sortable` tidak konsisten** antar tabel yang struktur mirip.
3. **Validasi FE vs BE tidak sinkron** pada field yang sama — rawan bypass via API langsung.
4. **`RoleGuard` FE tidak konsisten dipakai** di beberapa halaman settings.
5. **BE terlalu percaya nilai dari client** untuk data uang (transaksi, sync) — pola berulang di sales & operational.
6. **Permission middleware hilang** di beberapa route tulis/baca (transaction create, sync push, sales/stock report).

Total temuan: **~52** (semua sudah diverifikasi baca kode langsung, bukan asumsi).
