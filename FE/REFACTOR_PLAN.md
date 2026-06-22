# Rencana Refactor: Struktur Folder `FE/src/features`

Tujuan: Menyelaraskan struktur folder FE dengan group menu di database (`002_seed_data.sql`),
dengan tetap mempertahankan penamaan teknis (English) di sisi FE.

**Status: SELESAI** — Dikerjakan 2026-06-21

---

## Struktur Akhir (Hasil Implementasi)

```
features/
├── auth/                        (tidak berubah)
├── menu/                        (tidak berubah)
│
├── sales/                       ← penjualan (tidak berubah)
│   ├── cashier/
│   └── transactions/
│
├── products/                    ← produk (dipecah dari inventory/)
│   ├── products/
│   ├── categories/
│   └── units/
│
├── procurement/                 ← pengadaan (dipecah dari inventory/)
│   ├── suppliers/
│   ├── purchases/
│   └── returns/
│
├── customers/                   ← pelanggan
│   ├── customers/               ← dijadikan sub-modul dari flat
│   └── receivables/             ← pindah dari finance/receivables/
│
├── finance/                     ← keuangan (tidak berubah, receivables sudah keluar)
│   ├── overview/
│   ├── cash-drawer/
│   ├── expenses/
│   └── my-cash/
│
├── reporting/                   ← pelaporan + beranda
│   ├── dashboard/               (tidak berubah)
│   ├── sales/                   ← dipecah dari reporting/reports/
│   ├── profit-loss/             ← dipecah dari reporting/reports/
│   ├── stock/                   ← dipecah dari reporting/reports/
│   └── cashier-performance/     ← dipecah dari reporting/reports/
│
├── operational/                 ← operasional (folder baru)
│   ├── shifts/                  ← pindah dari features/shifts/
│   └── sync/                    ← pindah dari features/sync/
│
└── settings/                    ← sistem
    ├── store/                   ← dipecah dari settings/ root
    ├── users/                   ← dipecah dari settings/ root
    ├── roles/                   (tidak berubah)
    ├── menus/                   (tidak berubah)
    ├── printer/                 ← dipecah dari settings/ root
    └── versions/                ← dipecah dari settings/ root
```

---

## Ringkasan Perubahan Per Tahap

### TAHAP 1 — `settings/` sub-modul ✅
- `StoreProfilePage` + `StoreProfileForm` → `settings/store/`
- `UserManagementPage` + `UserManagementTab` → `settings/users/`
- `PrinterSettingsPage` + `PrinterSettingsTab` → `settings/printer/`
- `AppVersionPage` + `AppVersionTab` → `settings/versions/`
- `router.tsx` diupdate: 4 path

### TAHAP 2 — `reporting/reports/` dipecah ✅
- `SalesReportPage` + komponen → `reporting/sales/`
- `ProfitLossPage` + komponen → `reporting/profit-loss/`
- `StockReportPage` + komponen → `reporting/stock/`
- `CashierPerformancePage` + komponen → `reporting/cashier-performance/`
- `reports.utils.ts` dihapus — `monthStart()` & `todayStr()` dipindah ke `shared/utils/date.ts`
- Legacy types dihapus (tidak dipakai)
- `router.tsx` diupdate: 5 path (ReportsPage dihapus dari route)

### TAHAP 3 — `operational/` ✅
- `features/shifts/` → `features/operational/shifts/`
- `features/sync/` → `features/operational/sync/`
- `Navbar.tsx` diupdate: import `useSyncStatus`
- `cashier.api.ts` diupdate: re-export `useActiveShiftQuery`
- `router.tsx` diupdate: 2 path

### TAHAP 4 — `customers/` nested + `receivables/` pindah ✅
- `customers/` flat → `customers/customers/` sub-modul
- `finance/receivables/` → `customers/receivables/`
- `customers/index.ts` diupdate agar barrel export tetap bekerja
- `router.tsx` diupdate: 2 path (CustomersPage, ReceivablesPage)

### TAHAP 5 — `inventory/` dipecah jadi `products/` + `procurement/` ✅
- `inventory/products/` → `products/products/`
- `inventory/categories/` → `products/categories/`
- `inventory/units/` → `products/units/`
- `inventory/suppliers/` → `procurement/suppliers/`
- `inventory/purchases/` → `procurement/purchases/`
- `inventory/returns/` → `procurement/returns/`
- Cross-import diupdate: 9 file dalam features + 5 file cashier
- `router.tsx` diupdate: 6 path

---

## File yang Terdampak (Total)

| File | Perubahan |
|---|---|
| `app/router.tsx` | 19 import path diupdate |
| `shared/components/layouts/Navbar.tsx` | 1 import diupdate |
| `shared/utils/date.ts` | Tambah `todayStr()`, `monthStart()` |
| `shared/utils/index.ts` | Export 2 fungsi baru |
| `features/customers/index.ts` | Barrel export diupdate |
| `sales/cashier/cashier.api.ts` | 2 import diupdate |
| `sales/cashier/cashier.store.ts` | 1 import diupdate |
| `sales/cashier/cashier.utils.ts` | 1 import diupdate |
| `sales/cashier/hooks/useBarcodeScan.ts` | 1 import diupdate |
| `sales/cashier/components/ProductSearch.tsx` | 1 import diupdate |
| `products/products/components/ProductTable.tsx` | 1 import diupdate |
| `products/products/components/ProductPrerequisiteGuard.tsx` | 2 import diupdate |
| `products/products/components/ProductFormModal.tsx` | 2 import diupdate |
| `products/products/components/ProductFilterBar.tsx` | 1 import diupdate |
| `products/products/components/PriceTierTab.tsx` | 1 import diupdate |
| `procurement/purchases/PurchasesPage.tsx` | 2 import diupdate |
| `procurement/purchases/components/PurchaseFormModal.tsx` | 3 import diupdate |
| `procurement/returns/ReturnsPage.tsx` | 2 import diupdate |
| `procurement/returns/components/ReturnFormModal.tsx` | 1 import diupdate |
