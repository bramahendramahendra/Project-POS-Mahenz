# POS Mahenz — Backend API

REST API untuk aplikasi POS, dibangun dengan **Go + Gin + GORM + MySQL**.

## Prasyarat

- Go 1.24+
- MySQL 8.0+

## Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Buat database
mysql -u root -p -e "CREATE DATABASE pos_retail_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 3. Sesuaikan config
cp .env.local .env
# Edit config/config_dev.json → sesuaikan Database.Password

# 4. Jalankan
go run main.go
# Migrasi otomatis, server di http://localhost:8080
```

## Konfigurasi

| File | Keterangan |
|------|-----------|
| `.env` | `RELEASE_MODE` (dev/prod), `APP_PORT` |
| `config/config_dev.json` | Config development (DB, log, CORS, backup) |
| `config/config_prod.json` | Config production |

## Arsitektur

```
BE/
├── config/             Konfigurasi (viper)
├── database/
│   └── migrations/     SQL migrasi (auto-run saat startup)
├── domain/             Business logic per domain
│   ├── auth/           Login, session, JWT
│   ├── product/        Produk, package, import
│   ├── transaction/    Penjualan, checkout
│   ├── supplier_purchase/  Pembelian supplier
│   ├── supplier_return/    Retur supplier
│   ├── customer/       Pelanggan
│   ├── customer_balance/   Saldo pelanggan (deposit)
│   ├── receivable/     Piutang
│   ├── cash_drawer/    Kas harian
│   ├── expense/        Pengeluaran
│   ├── report/         Laporan
│   └── ...
├── middleware/         Auth, CORS, permission, logging
├── pkg/                Internal packages (database, binder, logger, jwt, pricing)
├── routes/             Router & route segments
├── scheduler/          Background jobs (auto-close kas)
└── main.go
```

Setiap domain mengikuti pola: `handler → service → repo → model/dto`.

## Migrasi Database

Migrasi **otomatis** saat startup. File SQL di `database/migrations/` dijalankan berurutan; yang sudah pernah jalan di-skip (tracking via tabel `migrations_history`).

```
database/migrations/
├── 001_init_schema.sql           Schema lengkap (28+ tabel)
├── 002_seed_data.sql             Data awal (user, role, menu, setting)
├── 003_fix_receivables_status_enum.sql
├── 004_stock_per_package_level.sql
├── 005_drop_legacy_product_stock_columns.sql
├── 006_customer_balance.sql      Saldo pelanggan
└── 007_fix_balance_payment_method.sql
```

**Menambah migrasi baru:** buat file `0XX_nama.sql` berikutnya, restart BE.

**Reset DB:** drop database, buat ulang, jalankan BE.

## Build

```bash
go build -o pos_api main.go
./pos_api
```

## API

Semua endpoint di bawah `/api/`. Auth via `Authorization: Bearer <token>`.

| Grup | Prefix | Contoh |
|------|--------|--------|
| Auth | `/api/auth/` | login, refresh, logout |
| Produk | `/api/products/` | list, create, search, packages |
| Transaksi | `/api/transactions/` | create, void, list |
| Pembelian | `/api/supplier-purchases/` | create, pay, void |
| Pelanggan | `/api/customers/` | CRUD, balance/topup, balance/history |
| Kas | `/api/cash-drawer/` | open, close, summary |
| Laporan | `/api/reports/` | sales, stock, profit-loss |

Health check: `GET /api/health`
