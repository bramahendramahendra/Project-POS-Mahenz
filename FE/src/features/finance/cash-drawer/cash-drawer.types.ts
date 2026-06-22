export interface CashDrawer {
  id: number
  open_time: string
  close_time?: string
  user_name: string
  shift_name?: string
  opening_balance: number
  closing_balance?: number
  expected_balance: number
  difference?: number
  total_sales: number
  total_cash_sales: number
  total_expenses: number
  status: 'open' | 'closed'
}

export interface CashDrawerDetail {
  id: number
  open_time: string
  close_time?: string
  cashier_name: string
  shift_name?: string
  shift_start?: string
  shift_end?: string
  opening_balance: number
  closing_balance?: number
  expected_balance: number
  difference?: number
  total_cash_sales: number
  total_expenses: number
  status: 'open' | 'closed'
  notes?: string
  open_notes?: string
  transactions: CashDrawerTransaction[]
  expenses: CashDrawerExpenseItem[]
}

export interface CashDrawerTransaction {
  transaction_date: string
  transaction_code: string
  customer_name: string
  total_amount: number
}

export interface CashDrawerExpenseItem {
  category: string
  description: string
  amount: number
}

export interface CashDrawerListFilter {
  page: number
  limit: number
  start_date?: string
  end_date?: string
  status?: string
  shift_id?: number
  user_id?: number
  [key: string]: unknown
}

export interface KasirOption {
  id: number
  full_name: string
  username: string
}

export interface CloseCashDrawerBody {
  closing_balance: number
  notes?: string
}

export interface OpenCashDrawerPayload {
  shift_id: number
  opening_balance: number
  notes?: string
}

export interface CurrentCashDrawer {
  id: number
  status: 'open' | 'closed'
  shift_id?: number
  shift_name?: string
  shift_start?: string
  shift_end?: string
}

export interface CashDrawerSummary {
  total_opening: number
  total_closing: number
  total_expenses: number
  net: number
  records: CashDrawer[]
}
