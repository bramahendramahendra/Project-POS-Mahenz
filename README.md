# POS Mahenz

Versi APlikasi : 2.8.4
Aplikasi **Point of Sale (POS)** untuk retail — mengelola transaksi penjualan, stok produk, supplier, keuangan, dan laporan.

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Backend | Go 1.24 (Gin), MySQL 8.0, JWT, GORM |
| Frontend | React 19, TypeScript, Vite 8, TanStack Query, Zustand |
| UI | Tailwind CSS, Radix UI (shadcn-style) |

## Fitur

- **Kasir** — checkout multi-metode (tunai, transfer, QRIS, kartu, hutang, saldo pelanggan)
- **Produk** — multi-satuan/paket (Slop → Pack → Batang), barcode, SKU, import massal
- **Pembelian** — PO supplier, multi-satuan per produk dalam 1 nota, bayar bertahap, void
- **Retur** — retur ke supplier dengan reservasi stok
- **Pelanggan** — saldo deposit, piutang (hutang), riwayat mutasi
- **Keuangan** — kas harian, pengeluaran, rekonsiliasi shift
- **Stok** — mutasi otomatis, batch expired, stok per level satuan
- **Laporan** — penjualan, laba rugi, stok, kinerja kasir, ringkasan bisnis
- **RBAC** — role & permission per menu (owner, admin, kasir)
- **Sync** — offline-online dengan conflict resolution

## Struktur

```
PORTAL/BE/ Backend Go Portal (login terpusat / SSO)
PORTAL/FE/ Frontend React + Vite Portal
TOKO/BE/   Backend Go POS Toko (domain-driven)
TOKO/FE/   Frontend React + Vite POS Toko
mobile/    Pembungkus APK (Capacitor)
nginx/     Konfigurasi nginx
docs/      Dokumentasi
```

Rencana Portal dan POS Coffee: lihat [docs/RANCANGAN_PORTAL_TOKO_COFFEE.md](docs/RANCANGAN_PORTAL_TOKO_COFFEE.md).

## Quick Start

```bash
# Backend Toko (port 8081)
cd TOKO/BE
cp .env.local .env        # sesuaikan config
go run main.go            # migrasi DB otomatis saat startup

# Frontend Toko (port 3001)
cd TOKO/FE
npm install
npm run type-check
npm run lint
npm run dev
```

Portal (login terpusat) dijalankan dengan cara yang sama dari `PORTAL/BE` (port 8080, database `pos_portal_db`
harus dibuat dulu) dan `PORTAL/FE` (port 3000). Untuk mencoba alur Portal ke Toko di laptop, isi alamat aplikasi
POS Toko di Portal dengan `http://localhost:3001/` dan pasang kunci aplikasinya di tabel `portal_credentials` Toko
(lihat [SOP Deploy](docs/SOP_DEPLOY.md), bagian "Kunci SSO aplikasi").

**Prasyarat:** Go 1.24+, Node 20+, MySQL 8.0

**Kredensial default (setelah seed):**
- Toko: Owner `owner` / `owner123`, Admin `admin` / `admin123`
- Portal: Admin `admin` / `admin123`

Ganti semua password bawaan sebelum production.

## Dokumentasi

- [Rancangan Portal, Toko, dan Coffee](docs/RANCANGAN_PORTAL_TOKO_COFFEE.md) (keputusan dan riwayat pengerjaan)
- [Panduan Cutover VPS ke Portal + Toko](docs/PANDUAN_CUTOVER_PORTAL.md)
- [SOP Deploy](docs/SOP_DEPLOY.md)
- [Panduan Deploy Production](docs/DEPLOYMENT_PROD.md)
- [Setup User Server](docs/SETUP_USER_SERVER.md)
- [Redeploy Guide](docs/DEPLOYMENT_REDEPLOY_FULL.md)
- [Panduan Migrasi Stok ke Skema Baru](docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md)
