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
AUTH/BE/   Backend Go BE User (layanan akun internal, tanpa FE, hanya 127.0.0.1)
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
# Backend Toko (port 8082)
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

Urutan menjalankan di laptop (MySQL harus hidup; migrasi DB otomatis saat startup):

| Layanan | Folder | Port | Database |
|---------|--------|------|----------|
| BE User (akun, password, akses aplikasi, kode SSO) | `AUTH/BE` | 8080 (hanya 127.0.0.1) | `pos_auth_db` |
| BE Portal | `PORTAL/BE` | 8081 | `pos_portal_db` |
| BE Toko | `TOKO/BE` | 8082 | `pos_retail_db` |
| FE Portal / FE Toko | `PORTAL/FE` / `TOKO/FE` | 3000 / 3001 | - |

Buat database kosong dulu (`pos_auth_db`, `pos_portal_db`, `pos_retail_db`). BE User dijalankan pertama: Portal
dan Toko memanggilnya untuk setiap login. Setiap aplikasi butuh kunci aplikasi dari BE User:

```bash
cd AUTH/BE
go run ./cmd/set_app_key --app portal --confirm   # ulangi untuk toko (dan coffee nanti)
```

Kunci tampil sekali; pasang di tabel `auth_credentials` database aplikasinya (contoh perintah ikut dicetak).
Untuk mencoba SSO Portal ke Toko di laptop, ubah path aplikasi `toko` di `pos_auth_db.apps` menjadi
`http://localhost:3001/`. Akun dan password dikelola di Portal; Toko hanya mengatur role, PIN, dan status.
User Toko lama disalin ke BE User dengan `TOKO/BE/cmd/migrate_users_to_auth` (dry-run dulu, lihat komentar skrip).

**Prasyarat:** Go 1.24+, Node 20+, MySQL 8.0

**Kredensial default (setelah seed):**
- BE User (dipakai login Portal dan Toko): Admin `admin` / `admin123`
- Toko: profil lokal `owner` dan `admin` (password lokal tidak dipakai lagi; akun disalin ke BE User dengan skrip di atas)

Ganti semua password bawaan sebelum production.

## Dokumentasi

- [Rancangan Portal, Toko, dan Coffee](docs/RANCANGAN_PORTAL_TOKO_COFFEE.md) (keputusan dan riwayat pengerjaan)
- [SOP Deploy per aplikasi](docs/SOP_DEPLOY.md) (deploy rutin)
- [Panduan Cutover VPS](docs/PANDUAN_CUTOVER_PORTAL.md) (sudah dijalankan 15 Nov 2026)
- [Pekerjaan pending](docs/PENDING.md)
- [Setup User Server](docs/SETUP_USER_SERVER.md)
- Arsip (usang): [Panduan Deploy Production](docs/DEPLOYMENT_PROD.md), [Redeploy Guide](docs/DEPLOYMENT_REDEPLOY_FULL.md)
