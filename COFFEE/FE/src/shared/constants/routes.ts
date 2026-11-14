export const ROUTES = {
  LOGIN: '/login',
  SSO: '/sso',
  PENDING_ACCESS: '/menunggu-akses',
  DASHBOARD: '/dashboard',
  KASIR: '/kasir',

  // Inventori
  PRODUCTS: '/products',
  PRODUCTS_CATEGORIES: '/products/categories',
  PRODUCTS_UNITS: '/products/units',

  // Menu jualan
  MENU_ITEMS: '/menu-items',
  MENU_CATEGORIES: '/menu-items/categories',
  MENU_MODIFIERS: '/menu-items/modifiers',
  MENU_SIZES: '/menu-items/sizes',
  MENU_TABLES: '/menu-items/tables',

  // Bahan Baku
  INGREDIENTS: '/ingredients',
  INGREDIENTS_HISTORY: '/ingredients/history',
  SUPPLIERS: '/suppliers',
  SUPPLIER_PURCHASES: '/suppliers/purchases',

  // Penjualan
  TRANSACTIONS: '/transactions',

  // Pelanggan
  CUSTOMERS: '/customers',
  RECEIVABLES: '/receivables',

  // Keuangan
  FINANCE: '/finance',
  FINANCE_CASH_DRAWER: '/finance/cash-drawer',
  FINANCE_EXPENSES: '/finance/expenses',
  FINANCE_MY_CASH: '/finance/my-cash',

  // Laporan
  REPORTS: '/reports',
  REPORTS_SALES: '/reports/sales',
  REPORTS_PROFIT_LOSS: '/reports/profit-loss',
  REPORTS_STOCK: '/reports/stock',
  REPORTS_STOCK_RECONCILIATION: '/reports/stock-reconciliation',
  REPORTS_CASHIER: '/reports/cashier',
  REPORTS_INGREDIENTS: '/reports/ingredients',
  REPORTS_WASTE: '/reports/waste',
  REPORTS_BUSINESS_SUMMARY: '/reports/business-summary',

  // Operasional
  SHIFTS: '/shifts',

  // Sistem
  SETTINGS: '/settings',
  SETTINGS_STORE: '/settings/store',
  SETTINGS_USERS: '/settings/users',
  SETTINGS_PRINTER: '/settings/printer',
  SETTINGS_ROLES: '/settings/roles',
  SETTINGS_ROLES_ACCESS: '/settings/roles/:id/access',
  SETTINGS_MENUS: '/settings/menus',
  SETTINGS_BACKUP: '/settings/backup',
} as const
