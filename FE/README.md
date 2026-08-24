# POS Mahenz — Frontend

Frontend aplikasi POS berbasis **React 19 + TypeScript + Vite**.

## Tech Stack

- **React 19** + TypeScript 6
- **Vite 8** — dev server & build
- **TanStack Query** — server state management
- **Zustand** — client state (auth, cashier cart, menu)
- **React Hook Form** + Zod — form & validasi
- **Tailwind CSS** + Radix UI — styling & komponen
- **react-router-dom v7** — routing (lazy-loaded, permission-guarded)

## Setup

```bash
npm install
npm run dev         # dev server di http://localhost:3000
```

## Scripts

| Command | Keterangan |
|---------|-----------|
| `npm run dev` | Development server (port 3000, proxy API ke :8080) |
| `npm run build` | Production build (tsc + vite build) |
| `npm run type-check` | TypeScript check tanpa emit |
| `npm run lint` | ESLint |
| `npm run lint:fix` | ESLint auto-fix |

## Arsitektur

```
src/
├── app/            Providers, router, layout
├── features/       Feature modules (per domain)
│   ├── auth/
│   ├── sales/      Kasir, transaksi
│   ├── products/   Produk, kategori, unit
│   ├── procurement/  Supplier, pembelian, retur
│   ├── customers/  Pelanggan, piutang
│   ├── finance/    Kas, pengeluaran
│   ├── reporting/  Laporan
│   ├── settings/   Pengaturan, user, role
│   └── ...
├── shared/         Komponen, hooks, utils, types bersama
└── services/       API client (axios + interceptor)
```

Setiap feature punya: `*.types.ts`, `*.api.ts`, `*.schema.ts`, `components/`, dan page.

## Environment

Buat file `.env.local`:
```env
VITE_API_URL=http://localhost:8080/api
```

Production build mengambil dari `.env.production`.

## Build

```bash
npm run build       # output di dist/
```

Hasil build di-serve sebagai static files (Nginx/Caddy). Lihat `docs/DEPLOYMENT_PROD.md`.
