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
}

export interface StockSummary {
  total_products: number
  low_stock_count: number
  total_stock_value: number
}

export interface StockFilter {
  search?: string
  category_id?: number
}

export interface StockListFilter extends StockFilter {
  page: number
  limit: number
}
