/* eslint-disable react-refresh/only-export-components */
import { lazy } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'

import { LoginPage, ProtectedRoute, RootRedirect } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'
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

// Procurement
const SuppliersPage    = lazy(() => import('@/features/procurement/suppliers/SuppliersPage').then(m => ({ default: m.SuppliersPage })))
const PurchasesPage    = lazy(() => import('@/features/procurement/purchases/PurchasesPage').then(m => ({ default: m.PurchasesPage })))
const ReturnsPage      = lazy(() => import('@/features/procurement/returns/ReturnsPage').then(m => ({ default: m.ReturnsPage })))

// Finance
const FinancePage           = lazy(() => import('@/features/finance/overview/FinancePage').then(m => ({ default: m.FinancePage })))
const CashDrawerPage        = lazy(() => import('@/features/finance/cash-drawer/CashDrawerPage').then(m => ({ default: m.CashDrawerPage })))
const MyCashPage            = lazy(() => import('@/features/finance/my-cash/MyCashPage').then(m => ({ default: m.MyCashPage })))
const ExpensesPage          = lazy(() => import('@/features/finance/expenses/ExpensesPage').then(m => ({ default: m.ExpensesPage })))
const ReceivablesPage       = lazy(() => import('@/features/customers/receivables/ReceivablesPage').then(m => ({ default: m.ReceivablesPage })))

// Reporting
const DashboardPage          = lazy(() => import('@/features/reporting/dashboard/DashboardPage').then(m => ({ default: m.DashboardPage })))
const SalesReportPage        = lazy(() => import('@/features/reporting/sales/SalesReportPage').then(m => ({ default: m.SalesReportPage })))
const ProfitLossPage         = lazy(() => import('@/features/reporting/profit-loss/ProfitLossPage').then(m => ({ default: m.ProfitLossPage })))
const StockReportPage        = lazy(() => import('@/features/reporting/stock/StockReportPage').then(m => ({ default: m.StockReportPage })))
const CashierPerformancePage = lazy(() => import('@/features/reporting/cashier-performance/CashierPerformancePage').then(m => ({ default: m.CashierPerformancePage })))

// Operational
const ShiftsPage       = lazy(() => import('@/features/operational/shifts/ShiftsPage').then(m => ({ default: m.ShiftsPage })))
const SyncCenterPage   = lazy(() => import('@/features/operational/sync/SyncCenterPage').then(m => ({ default: m.SyncCenterPage })))

// Settings
const SettingsPage         = lazy(() => import('@/features/settings/SettingsPage').then(m => ({ default: m.SettingsPage })))
const StoreProfilePage     = lazy(() => import('@/features/settings/store/StoreProfilePage').then(m => ({ default: m.StoreProfilePage })))
const UserManagementPage   = lazy(() => import('@/features/settings/users/UserManagementPage').then(m => ({ default: m.UserManagementPage })))
const PrinterSettingsPage  = lazy(() => import('@/features/settings/printer/PrinterSettingsPage').then(m => ({ default: m.PrinterSettingsPage })))
const AppVersionPage       = lazy(() => import('@/features/settings/versions/AppVersionPage').then(m => ({ default: m.AppVersionPage })))
const RolesPage            = lazy(() => import('@/features/settings/roles/RolesPage').then(m => ({ default: m.RolesPage })))
const RoleAccessPage       = lazy(() => import('@/features/settings/roles/RoleAccessPage').then(m => ({ default: m.RoleAccessPage })))
const MenusPage            = lazy(() => import('@/features/settings/menus/MenusPage').then(m => ({ default: m.MenusPage })))

const ALL_ROLES        = [ROLES.OWNER, ROLES.ADMIN, ROLES.KASIR] as const
const MANAGEMENT_ROLES = [ROLES.OWNER, ROLES.ADMIN] as const
const OWNER_ONLY       = [ROLES.OWNER] as const

export const router = createBrowserRouter([
  { path: '/', element: <RootRedirect /> },
  { path: ROUTES.LOGIN, element: <LoginPage /> },

  // Protected — semua role
  {
    element: <ProtectedRoute allowedRoles={[...ALL_ROLES]} />,
    children: [
      { path: ROUTES.KASIR,           element: <LazyRoute><CashierPage /></LazyRoute> },
      { path: ROUTES.FINANCE_MY_CASH, element: <LazyRoute><MyCashPage /></LazyRoute> },
    ],
  },

  // Protected — owner & admin
  {
    element: <ProtectedRoute allowedRoles={[...MANAGEMENT_ROLES]} />,
    children: [
      // Dashboard
      { path: ROUTES.DASHBOARD,    element: <LazyRoute><DashboardPage /></LazyRoute> },

      // Inventori
      { path: ROUTES.PRODUCTS,            element: <LazyRoute><ProductsPage /></LazyRoute> },
      { path: ROUTES.PRODUCTS_CATEGORIES, element: <LazyRoute><CategoryPage /></LazyRoute> },
      { path: ROUTES.PRODUCTS_UNITS,      element: <LazyRoute><UnitPage /></LazyRoute> },
      { path: ROUTES.SUPPLIERS,           element: <LazyRoute><SuppliersPage /></LazyRoute> },
      { path: ROUTES.SUPPLIER_PURCHASES,  element: <LazyRoute><PurchasesPage /></LazyRoute> },
      { path: ROUTES.SUPPLIER_RETURNS,    element: <LazyRoute><ReturnsPage /></LazyRoute> },

      // Penjualan
      { path: ROUTES.TRANSACTIONS, element: <LazyRoute><TransactionsPage /></LazyRoute> },

      // Pelanggan
      { path: ROUTES.CUSTOMERS,   element: <LazyRoute><CustomersPage /></LazyRoute> },
      { path: ROUTES.RECEIVABLES, element: <LazyRoute><ReceivablesPage /></LazyRoute> },

      // Keuangan
      { path: ROUTES.FINANCE,                   element: <LazyRoute><FinancePage /></LazyRoute> },
      { path: ROUTES.FINANCE_CASH_DRAWER, element: <LazyRoute><CashDrawerPage /></LazyRoute> },
      { path: ROUTES.FINANCE_EXPENSES,    element: <LazyRoute><ExpensesPage /></LazyRoute> },

      // Laporan
      { path: ROUTES.REPORTS_SALES,     element: <LazyRoute><SalesReportPage /></LazyRoute> },
      { path: ROUTES.REPORTS_PROFIT_LOSS, element: <LazyRoute><ProfitLossPage /></LazyRoute> },
      { path: ROUTES.REPORTS_STOCK,     element: <LazyRoute><StockReportPage /></LazyRoute> },
      { path: ROUTES.REPORTS_CASHIER,   element: <LazyRoute><CashierPerformancePage /></LazyRoute> },

      // Operasional
      { path: ROUTES.SHIFTS, element: <LazyRoute><ShiftsPage /></LazyRoute> },
      { path: ROUTES.SYNC,   element: <LazyRoute><SyncCenterPage /></LazyRoute> },

      // Pengaturan — owner & admin
      { path: ROUTES.SETTINGS,              element: <LazyRoute><SettingsPage /></LazyRoute> },
      { path: ROUTES.SETTINGS_STORE,        element: <LazyRoute><StoreProfilePage /></LazyRoute> },
      { path: ROUTES.SETTINGS_PRINTER,      element: <LazyRoute><PrinterSettingsPage /></LazyRoute> },
      { path: ROUTES.SETTINGS_ROLES,        element: <LazyRoute><RolesPage /></LazyRoute> },
      { path: ROUTES.SETTINGS_ROLES_ACCESS, element: <LazyRoute><RoleAccessPage /></LazyRoute> },
      { path: ROUTES.SETTINGS_MENUS,        element: <LazyRoute><MenusPage /></LazyRoute> },
      { path: ROUTES.SETTINGS_VERSIONS,     element: <LazyRoute><AppVersionPage /></LazyRoute> },
    ],
  },

  // Protected — owner only
  {
    element: <ProtectedRoute allowedRoles={[...OWNER_ONLY]} />,
    children: [
      { path: ROUTES.SETTINGS_USERS, element: <LazyRoute><UserManagementPage /></LazyRoute> },
    ],
  },

  // 404
  { path: '*', element: <Navigate to="/" replace /> },
])
