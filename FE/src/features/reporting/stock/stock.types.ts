export interface StockReport {
  product_code: string
  product_name: string
  category_name: string
  unit: string
  current_stock: number
  min_stock: number
  cost_price: number
  stock_value: number
}

export interface StockReportFilter {
  category_id?: number
  search?: string
  page?: number
  page_size?: number
}
