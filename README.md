# Project POS Mahenz

Aplikasi **Point of Sale (POS)** untuk retail — mengelola transaksi penjualan, stok produk, supplier, keuangan, hingga laporan, dengan backend Go dan frontend React.

## Tech Stack

**Backend**
- Go (Gin) — REST API
- MySQL 8.0 — database utama, migrasi otomatis saat startup
- JWT — autentikasi
- Redis — caching/permission cache

**Frontend**
- React 19 + TypeScript
- Vite — build tool
- TanStack Query — data fetching & caching
- Redux + Zustand — state management
- Tailwind CSS + Radix UI — styling & komponen UI

## Fitur Utama

- Kasir / transaksi penjualan
- Manajemen produk, kategori, dan satuan
- Manajemen stok & mutasi stok
- Supplier (pembelian & retur)
- Kas & shift kasir (cash drawer)
- Piutang (receivable) & pengeluaran (expense)
- Manajemen user, role & hak akses (menu/permission)
- Laporan & dashboard
- Sinkronisasi data (mode offline-online)
- Maintenance mode untuk aplikasi (lihat [docs/DEPLOYMENT_PROD.md](docs/DEPLOYMENT_PROD.md))

## Struktur Project

```
BE/      Backend Go (domain-driven: domain/, routes/, middleware/, dll)
FE/      Frontend React + Vite
docs/    Dokumentasi deployment & audit
```

## Dokumentasi

- [Panduan Instalasi & Deploy Production](docs/DEPLOYMENT_PROD.md)
- [Setup User Server](docs/SETUP_USER_SERVER.md)

## Menjalankan Secara Lokal

### Backend
```bash
cd BE
go run main.go
```

### Frontend
```bash
cd FE
npm install
npm run dev
```

Lihat [docs/DEPLOYMENT_PROD.md](docs/DEPLOYMENT_PROD.md) untuk konfigurasi environment, database, dan panduan deploy ke production.
