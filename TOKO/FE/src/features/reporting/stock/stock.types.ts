export interface StockReport {
  id: number
  product_code: string
  product_name: string
  category_name: string
  current_stock: number
  min_stock: number
  unit: string
  cost_price: number
  stock_value: number
  is_low_stock: boolean
  stock_status: StockStatus
}

// "out" = habis, "low" = menipis, "ok" = aman, "" = data diragukan (tanpa alarm)
export type StockStatus = 'out' | 'low' | 'ok' | ''

export interface StockSummary {
  total_products: number
  low_stock_count: number // habis + menipis
  out_of_stock_count: number // habis saja
  low_only_count: number // menipis saja
  total_stock_value: number
}

export interface StockFilter {
  search?: string
  category_id?: number
  stock_status?: StockStatus
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface StockListFilter extends StockFilter {
  page: number
  limit: number
}
