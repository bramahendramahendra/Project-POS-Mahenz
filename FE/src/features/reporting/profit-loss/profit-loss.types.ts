export interface ProfitLossDateFilter {
  date_from?: string
  date_to?: string
}

export interface ProfitLossItem {
  product_id: number
  product_name: string
  qty_sold: number
  purchase_price: number
  total_cogs: number
  total_revenue: number
  gross_profit: number
}

export interface ExpenseSummaryItem {
  category: string
  total: number
}

export interface ProfitLossReport {
  total_revenue: number
  total_cogs: number
  gross_profit: number
  total_expenses: number
  net_profit: number
  items: ProfitLossItem[]
  expenses: ExpenseSummaryItem[]
}
