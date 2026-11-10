export interface RecentTransactionItem {
  id: number
  transaction_code: string
  transaction_date: string
  total_amount: number
  payment_method: string
  status: string
}

export interface TodaySummary {
  total_transactions: number
  total_sales: number
}
