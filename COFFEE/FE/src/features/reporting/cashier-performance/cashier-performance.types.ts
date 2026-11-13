export interface CashierPerformanceDateFilter {
  date_from?: string
  date_to?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface CashierPerformanceListFilter extends CashierPerformanceDateFilter {
  page: number
  limit: number
}

export interface CashierPerformance {
  user_id: number
  cashier_name: string
  total_transactions: number
  total_sales: number
  total_cash: number
  total_non_cash: number
  avg_per_transaction: number
  void_count: number
}
