/* eslint-disable react-refresh/only-export-components */
import { lazy } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'

import { LoginPage, PendingAccessPage, ProtectedRoute, RootRedirect, SsoPage } from '@/features/auth'
import { routerBasename } from '@/shared/utils/appUrl'
import { ROUTES } from '@/shared/constants/routes'

import { LazyRoute } from './LazyRoute'

// Sales
const CashierPage      = lazy(() => import('@/features/sales/cashier/CashierPage').then(m => ({ default: m.CashierPage })))
const TransactionsPage = lazy(() => import('@/features/sales/transactions/TransactionsPage').then(m => ({ default: m.TransactionsPage })))

// Customers
const CustomersPage    = lazy(() => import('@/features/customers/customers/CustomersPage').then(m => ({ default: m.CustomersPage })))

// Products
const ProductsPage     = lazy(() => import('@/features/products/products/ProductsPage').then(m => ({ default: m.ProductsPage })))
const CategoryPage     = lazy(() => import('@/features/products/categories/CategoryPage').then(m => ({ default: m.CategoryPage })))
const UnitPage         = lazy(() => import('@/features/products/units/UnitPage').then(m => ({ default: m.UnitPage })))

// Menu jualan
const MenuItemsPage      = lazy(() => import('@/features/menu-items/MenuItemsPage').then(m => ({ default: m.MenuItemsPage })))
const MenuCategoriesPage = lazy(() => import('@/features/menu-items/MenuCategoriesPage').then(m => ({ default: m.MenuCategoriesPage })))
const ModifierGroupsPage = lazy(() => import('@/features/menu-items/ModifierGroupsPage').then(m => ({ default: m.ModifierGroupsPage })))
const MenuSizesPage      = lazy(() => import('@/features/menu-items/MenuSizesPage').then(m => ({ default: m.MenuSizesPage })))
const DiningTablesPage   = lazy(() => import('@/features/menu-items/DiningTablesPage').then(m => ({ default: m.DiningTablesPage })))

// Bahan baku
const IngredientsPage       = lazy(() => import('@/features/ingredients/IngredientsPage').then(m => ({ default: m.IngredientsPage })))
const IngredientHistoryPage = lazy(() => import('@/features/ingredients/IngredientHistoryPage').then(m => ({ default: m.IngredientHistoryPage })))

// Procurement
const SuppliersPage    = lazy(() => import('@/features/procurement/suppliers/SuppliersPage').then(m => ({ default: m.SuppliersPage })))
const PurchasesPage    = lazy(() => import('@/features/procurement/purchases/PurchasesPage').then(m => ({ default: m.PurchasesPage })))

// Finance
const FinancePage           = lazy(() => import('@/features/finance/overview/FinancePage').then(m => ({ default: m.FinancePage })))
const CashDrawerPage        = lazy(() => import('@/features/finance/cash-drawer/CashDrawerPage').then(m => ({ default: m.CashDrawerPage })))
const MyCashPage            = lazy(() => import('@/features/finance/my-cash/MyCashPage').then(m => ({ default: m.MyCashPage })))
const ExpensesPage          = lazy(() => import('@/features/finance/expenses/ExpensesPage').then(m => ({ default: m.ExpensesPage })))
const ReceivablesPage       = lazy(() => import('@/features/customers/receivables/ReceivablesPage').then(m => ({ default: m.ReceivablesPage })))

// Beranda
const DashboardPage          = lazy(() => import('@/features/dashboard/DashboardPage').then(m => ({ default: m.DashboardPage })))

// Reporting
const BusinessSummaryPage    = lazy(() => import('@/features/reporting/business-summary/BusinessSummaryPage').then(m => ({ default: m.BusinessSummaryPage })))
const SalesReportPage        = lazy(() => import('@/features/reporting/sales/SalesReportPage').then(m => ({ default: m.SalesReportPage })))
const ProfitLossPage         = lazy(() => import('@/features/reporting/profit-loss/ProfitLossPage').then(m => ({ default: m.ProfitLossPage })))
const StockReportPage        = lazy(() => import('@/features/reporting/stock/StockReportPage').then(m => ({ default: m.StockReportPage })))
const StockReconciliationPage = lazy(() => import('@/features/reporting/stock-reconciliation').then(m => ({ default: m.StockReconciliationPage })))
const CashierPerformancePage = lazy(() => import('@/features/reporting/cashier-performance/CashierPerformancePage').then(m => ({ default: m.CashierPerformancePage })))

// Operational
const ShiftsPage       = lazy(() => import('@/features/operational/shifts/ShiftsPage').then(m => ({ default: m.ShiftsPage })))

// Settings
const SettingsPage         = lazy(() => import('@/features/settings/SettingsPage').then(m => ({ default: m.SettingsPage })))
const StoreProfilePage     = lazy(() => import('@/features/settings/store/StoreProfilePage').then(m => ({ default: m.StoreProfilePage })))
const UserManagementPage   = lazy(() => import('@/features/settings/users/UserManagementPage').then(m => ({ default: m.UserManagementPage })))
const PrinterSettingsPage  = lazy(() => import('@/features/settings/printer/PrinterSettingsPage').then(m => ({ default: m.PrinterSettingsPage })))
const RolesPage            = lazy(() => import('@/features/settings/roles/RolesPage').then(m => ({ default: m.RolesPage })))
const RoleAccessPage       = lazy(() => import('@/features/settings/roles/RoleAccessPage').then(m => ({ default: m.RoleAccessPage })))
const MenusPage            = lazy(() => import('@/features/settings/menus/MenusPage').then(m => ({ default: m.MenusPage })))
const BackupPage           = lazy(() => import('@/features/settings/backup/BackupPage').then(m => ({ default: m.BackupPage })))

// Setiap route dipetakan ke menu_key di tabel `menus` (lihat database/migrations/002_seed_data.sql).
// Akses ditentukan oleh role_menu_access (permission.can_view) via ProtectedRoute — bukan role hardcode.
interface RouteDef {
  path: string
  menuKey: string
  element: React.ReactNode
}

const PROTECTED_ROUTES: RouteDef[] = [
  // Beranda
  { path: ROUTES.DASHBOARD, menuKey: 'beranda.dashboard', element: <DashboardPage /> },

  // Penjualan
  { path: ROUTES.KASIR, menuKey: 'penjualan.kasir', element: <CashierPage /> },
  { path: ROUTES.TRANSACTIONS, menuKey: 'penjualan.transaksi', element: <TransactionsPage /> },


  // Produk
  { path: ROUTES.PRODUCTS, menuKey: 'produk.produk', element: <ProductsPage /> },
  { path: ROUTES.PRODUCTS_CATEGORIES, menuKey: 'produk.kategori', element: <CategoryPage /> },
  { path: ROUTES.PRODUCTS_UNITS, menuKey: 'produk.unit', element: <UnitPage /> },

  // Menu jualan
  { path: ROUTES.MENU_ITEMS, menuKey: 'menu_jual.daftar', element: <MenuItemsPage /> },
  { path: ROUTES.MENU_CATEGORIES, menuKey: 'menu_jual.kategori', element: <MenuCategoriesPage /> },
  { path: ROUTES.MENU_MODIFIERS, menuKey: 'menu_jual.tambahan', element: <ModifierGroupsPage /> },
  { path: ROUTES.MENU_SIZES, menuKey: 'menu_jual.ukuran', element: <MenuSizesPage /> },
  { path: ROUTES.MENU_TABLES, menuKey: 'menu_jual.meja', element: <DiningTablesPage /> },

  // Bahan Baku
  { path: ROUTES.INGREDIENTS, menuKey: 'bahan.daftar', element: <IngredientsPage /> },
  { path: ROUTES.INGREDIENTS_HISTORY, menuKey: 'bahan.riwayat', element: <IngredientHistoryPage /> },

  // Pengadaan
  { path: ROUTES.SUPPLIERS, menuKey: 'pengadaan.supplier', element: <SuppliersPage /> },
  { path: ROUTES.SUPPLIER_PURCHASES, menuKey: 'pengadaan.pembelian', element: <PurchasesPage /> },

  // Pelanggan
  { path: ROUTES.CUSTOMERS, menuKey: 'pelanggan.pelanggan', element: <CustomersPage /> },
  { path: ROUTES.RECEIVABLES, menuKey: 'pelanggan.piutang', element: <ReceivablesPage /> },

  // Keuangan
  { path: ROUTES.FINANCE, menuKey: 'keuangan.dashboard', element: <FinancePage /> },
  { path: ROUTES.FINANCE_CASH_DRAWER, menuKey: 'keuangan.kas_harian', element: <CashDrawerPage /> },
  { path: ROUTES.FINANCE_EXPENSES, menuKey: 'keuangan.pengeluaran', element: <ExpensesPage /> },
  { path: ROUTES.FINANCE_MY_CASH, menuKey: 'keuangan.kas_saya', element: <MyCashPage /> },

  // Pelaporan
  { path: ROUTES.REPORTS_BUSINESS_SUMMARY, menuKey: 'pelaporan.ringkasan_bisnis', element: <BusinessSummaryPage /> },
  { path: ROUTES.REPORTS_SALES, menuKey: 'pelaporan.penjualan', element: <SalesReportPage /> },
  { path: ROUTES.REPORTS_PROFIT_LOSS, menuKey: 'pelaporan.laba_rugi', element: <ProfitLossPage /> },
  { path: ROUTES.REPORTS_STOCK, menuKey: 'pelaporan.stok', element: <StockReportPage /> },
  { path: ROUTES.REPORTS_STOCK_RECONCILIATION, menuKey: 'pelaporan.rekonsiliasi_stok', element: <StockReconciliationPage /> },
  { path: ROUTES.REPORTS_CASHIER, menuKey: 'pelaporan.kinerja_kasir', element: <CashierPerformancePage /> },

  // Operasional
  { path: ROUTES.SHIFTS, menuKey: 'operasional.shift', element: <ShiftsPage /> },

  // Sistem
  { path: ROUTES.SETTINGS_STORE, menuKey: 'sistem.profil_toko', element: <StoreProfilePage /> },
  { path: ROUTES.SETTINGS_USERS, menuKey: 'sistem.users', element: <UserManagementPage /> },
  { path: ROUTES.SETTINGS_PRINTER, menuKey: 'sistem.printer', element: <PrinterSettingsPage /> },
  { path: ROUTES.SETTINGS_ROLES, menuKey: 'sistem.roles', element: <RolesPage /> },
  { path: ROUTES.SETTINGS_ROLES_ACCESS, menuKey: 'sistem.roles', element: <RoleAccessPage /> },
  { path: ROUTES.SETTINGS_MENUS, menuKey: 'sistem.menus', element: <MenusPage /> },
  { path: ROUTES.SETTINGS_BACKUP, menuKey: 'sistem.backup', element: <BackupPage /> },
  // Halaman pengaturan (grup) — mengikuti izin profil toko sebagai akses minimum settings
  { path: ROUTES.SETTINGS, menuKey: 'sistem.profil_toko', element: <SettingsPage /> },
]

export const router = createBrowserRouter([
  { path: '/', element: <RootRedirect /> },
  { path: ROUTES.LOGIN, element: <LoginPage /> },
  { path: ROUTES.SSO, element: <SsoPage /> },
  { path: ROUTES.PENDING_ACCESS, element: <PendingAccessPage /> },

  ...PROTECTED_ROUTES.map(({ path, menuKey, element }) => ({
    path,
    element: <ProtectedRoute menuKey={menuKey} />,
    children: [{ index: true, element: <LazyRoute>{element}</LazyRoute> }],
  })),

  // 404
  { path: '*', element: <Navigate to="/" replace /> },
], {
  // undefined bila aplikasi di root domain; '/coffee' saat disajikan di bawah Portal.
  basename: routerBasename,
})
