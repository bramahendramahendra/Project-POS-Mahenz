# POS Mahenz

Versi APlikasi : 2.6.5
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
BE/      Backend Go (domain-driven)
FE/      Frontend React + Vite
docs/    Dokumentasi
```

## Quick Start

```bash
# Backend (port 8080)
cd BE
cp .env.local .env        # sesuaikan config
go run main.go            # migrasi DB otomatis saat startup

# Frontend (port 3000)
cd FE
npm install
npm run dev
```

**Prasyarat:** Go 1.24+, Node 20+, MySQL 8.0

**Kredensial default (setelah seed):**
- Owner: `owner` / `owner123`
- Admin: `admin` / `admin123`

## Dokumentasi

- [Panduan Deploy Production](docs/DEPLOYMENT_PROD.md)
- [Setup User Server](docs/SETUP_USER_SERVER.md)
- [Redeploy Guide](docs/DEPLOYMENT_REDEPLOY_FULL.md)
- [Desain Fitur Saldo Pelanggan](docs/DESIGN_SALDO_PELANGGAN.md)
