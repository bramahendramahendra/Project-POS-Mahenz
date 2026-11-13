import type { CashDrawerTransaction, CashDrawerExpenseItem, NonCashSaleItem, NonCashTransaction } from '@/features/finance/cash-drawer'

export type { CashDrawerTransaction, CashDrawerExpenseItem, NonCashSaleItem, NonCashTransaction }

export interface MyCashData {
  id?: number
  status: 'open' | 'closed'
  shift_name?: string
  shift_start?: string
  shift_end?: string
  open_time?: string
  opening_balance: number
  total_cash_sales: number
  total_expenses: number
  expected_balance: number
  open_notes?: string
  transactions: CashDrawerTransaction[]
  expenses: CashDrawerExpenseItem[]
  non_cash_sales: NonCashSaleItem[]
  non_cash_transactions: NonCashTransaction[]
}
