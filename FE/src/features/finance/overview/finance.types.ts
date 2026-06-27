export interface FinanceSummary {
  total_income: number
  total_expense: number
  net_profit: number
  total_receivable: number
  period_label: string
}

export interface CashflowItem {
  id: number
  type: 'income' | 'expense'
  category: string
  amount: number
  description: string
  date: string
}

export interface FinanceDateFilter {
  date_from?: string
  date_to?: string
}

export interface CashflowFilter extends FinanceDateFilter {
  type?: 'income' | 'expense'
  page: number
  limit: number
}
