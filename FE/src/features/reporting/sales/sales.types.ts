export interface SalesReport {
  transaction_code: string
  transaction_date: string
  cashier_name: string
  customer_name?: string
  total_amount: number
  payment_method: string
  status: 'completed' | 'void'
}

export interface SalesReportSummary {
  total_transactions: number
  total_revenue: number
  avg_per_transaction: number
}

export interface SalesReportFilter {
  date_from?: string
  date_to?: string
  user_id?: number
  payment_method?: string
  page?: number
  page_size?: number
}
