type Filter = object | undefined

export const queryKeys = {
  auth: {
    profile: () => ['auth', 'profile'] as const,
  },

  products: {
    all: () => ['products'] as const,
    list: (filter?: Filter) => ['products', 'list', filter] as const,
    options: () => ['products', 'options'] as const,
    detail: (id: number) => ['products', 'detail', id] as const,
    productUnits: (id: number) => ['products', 'productUnits', id] as const,
    priceTiers: (id: number) => ['products', 'priceTiers', id] as const,
    barcode: (code: string) => ['products', 'barcode', code] as const,
  },

  categories: {
    all: () => ['categories'] as const,
    list: (filter?: Filter) => ['categories', 'list', filter] as const,
    options: () => ['categories', 'options'] as const,
  },

  units: {
    all: () => ['units'] as const,
    list: (filter?: Filter) => ['units', 'list', filter] as const,
    options: () => ['units', 'options'] as const,
  },

  suppliers: {
    all: () => ['suppliers'] as const,
    list: (filter?: Filter) => ['suppliers', 'list', filter] as const,
    options: () => ['suppliers', 'options'] as const,
    detail: (id: number) => ['suppliers', 'detail', id] as const,
  },

  transactions: {
    all: () => ['transactions'] as const,
    list: (filter?: Filter) => ['transactions', 'list', filter] as const,
    detail: (id: number) => ['transactions', 'detail', id] as const,
  },

  customers: {
    all: () => ['customers'] as const,
    list: (filter?: Filter) => ['customers', 'list', filter] as const,
    detail: (id: number) => ['customers', 'detail', id] as const,
  },

  receivables: {
    all: () => ['receivables'] as const,
    list: (filter?: Filter) => ['receivables', 'list', filter] as const,
    detail: (id: number) => ['receivables', 'detail', id] as const,
  },

  expenses: {
    all: () => ['expenses'] as const,
    list: (filter?: Filter) => ['expenses', 'list', filter] as const,
  },

  finance: {
    all: () => ['finance'] as const,
    summary: (filter?: Filter) => ['finance', 'summary', filter] as const,
    cashflow: (filter?: Filter) => ['finance', 'cashflow', filter] as const,
  },

  myCash: {
    data: () => ['myCash', 'data'] as const,
  },

  dashboard: {
    all: () => ['dashboard'] as const,
    stats: () => ['dashboard', 'stats'] as const,
    salesTrend: (period: string) => ['dashboard', 'salesTrend', period] as const,
    topProducts: (period: string) => ['dashboard', 'topProducts', period] as const,
  },

  reports: {
    all: () => ['reports'] as const,
    sales: (filter?: Filter) => ['reports', 'sales', filter] as const,
    salesList: (filter?: Filter) => ['reports', 'sales', 'list', filter] as const,
    salesSummary: (filter?: Filter) => ['reports', 'sales', 'summary', filter] as const,
    profitLoss: (filter?: Filter) => ['reports', 'profitLoss', filter] as const,
    stock: (filter?: Filter) => ['reports', 'stock', filter] as const,
    stockList: (filter?: Filter) => ['reports', 'stock', 'list', filter] as const,
    stockSummary: (filter?: Filter) => ['reports', 'stock', 'summary', filter] as const,
    cashierPerformance: (filter?: Filter) => ['reports', 'cashierPerformance', filter] as const,
    cashierPerformanceList: (filter?: Filter) => ['reports', 'cashierPerformance', 'list', filter] as const,
  },

  shifts: {
    all: () => ['shifts'] as const,
    list: (filter?: Filter) => ['shifts', 'list', filter] as const,
    active: () => ['shifts', 'active'] as const,
    detail: (id: number) => ['shifts', 'detail', id] as const,
  },

  settings: {
    store: () => ['settings', 'store'] as const,
    appVersions: () => ['settings', 'appVersions'] as const,
    pageSizeOptions: () => ['settings', 'pageSizeOptions'] as const,
    printer: () => ['settings', 'printer'] as const,
  },

  users: {
    all: () => ['users'] as const,
    list: (filter?: Filter) => ['users', 'list', filter] as const,
  },

  sync: {
    status: () => ['sync', 'status'] as const,
    history: (filter?: Filter) => ['sync', 'history', filter] as const,
    conflicts: () => ['sync', 'conflicts'] as const,
  },

  roles: {
    all: () => ['roles'] as const,
    list: (filter?: Filter) => ['roles', 'list', filter] as const,
    detail: (id: number) => ['roles', 'detail', id] as const,
    menus: (id: number) => ['roles', 'menus', id] as const,
    options: () => ['roles', 'options'] as const,
  },

  menus: {
    all: () => ['menus'] as const,
    list: (filter?: Filter) => ['menus', 'list', filter] as const,
    detail: (id: number) => ['menus', 'detail', id] as const,
    my: () => ['menus', 'my'] as const,
    options: () => ['menus', 'options'] as const,
  },

  cashDrawer: {
    all: () => ['cashDrawer'] as const,
    current: () => ['cashDrawer', 'current'] as const,
    list: (filter?: Filter) => ['cashDrawer', 'list', filter] as const,
    detail: (id: number) => ['cashDrawer', 'detail', id] as const,
    summary: (filter?: Filter) => ['cashDrawer', 'summary', filter] as const,
    kasirOptions: () => ['cashDrawer', 'kasirOptions'] as const,
  },

  supplierPurchases: {
    all: () => ['supplierPurchases'] as const,
    list: (filter?: Filter) => ['supplierPurchases', 'list', filter] as const,
    detail: (id: number) => ['supplierPurchases', 'detail', id] as const,
    payments: (id: number) => ['supplierPurchases', 'detail', id, 'payments'] as const,
    generateCode: () => ['supplierPurchases', 'generateCode'] as const,
  },

  paymentMethods: {
    options: () => ['paymentMethods', 'options'] as const,
  },

  paymentStatuses: {
    options: () => ['paymentStatuses', 'options'] as const,
  },

  supplierReturns: {
    all: () => ['supplierReturns'] as const,
    list: (filter?: Filter) => ['supplierReturns', 'list', filter] as const,
    detail: (id: number) => ['supplierReturns', 'detail', id] as const,
  },
}
