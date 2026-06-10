export type ExpenseCategory = 'operasional' | 'pembelian' | 'gaji' | 'lainnya'
export type ExpensePaymentMethod = 'cash' | 'transfer' | 'card' | 'qris' | 'kredit'

export interface Expense {
  id: number
  expense_date: string
  category: ExpenseCategory
  description: string
  amount: number
  payment_method: ExpensePaymentMethod
  notes?: string
  user_name: string
  created_at: string
}

export interface ExpenseListFilter {
  page: number
  limit: number
  start_date?: string
  end_date?: string
  category?: ExpenseCategory
}

export interface CreateExpensePayload {
  expense_date: string
  category: ExpenseCategory
  description: string
  amount: number
  payment_method: ExpensePaymentMethod
  notes?: string
}

export type UpdateExpensePayload = Partial<CreateExpensePayload>
