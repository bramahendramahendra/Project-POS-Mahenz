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
  /** Menu racikan: biaya bahannya masuk biaya bahan periode, bukan per menu (keputusan 114). */
  is_recipe: boolean
}

export interface ExpenseSummaryItem {
  category: string
  total: number
}

export interface ProfitLossReport {
  total_revenue: number
  total_cogs: number
  product_cogs: number
  ingredient_cost: number
  gross_profit: number
  total_expenses: number
  net_profit: number
  items: ProfitLossItem[]
  expenses: ExpenseSummaryItem[]
  prev_net_profit: number
  prev_available: boolean
}
