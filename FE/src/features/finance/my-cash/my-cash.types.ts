import type { CashDrawerTransaction, CashDrawerExpenseItem } from '@/features/finance/cash-drawer'

export type { CashDrawerTransaction, CashDrawerExpenseItem }

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
}
