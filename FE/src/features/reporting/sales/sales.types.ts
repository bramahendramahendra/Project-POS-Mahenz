export interface SalesReport {
  id: number
  transaction_code: string
  transaction_date: string
  cashier_name: string
  customer_name: string
  total_amount: number
  discount: number
  payment_method: string
  status: 'completed' | 'void'
}

export interface SalesReportSummary {
  total_transactions: number
  total_revenue: number
  avg_per_transaction: number
  total_discount: number
  total_tax: number
}

export interface SalesFilter {
  date_from?: string
  date_to?: string
  payment_method?: string
  user_id?: number
}

export interface SalesListFilter extends SalesFilter {
  page: number
  limit: number
}
