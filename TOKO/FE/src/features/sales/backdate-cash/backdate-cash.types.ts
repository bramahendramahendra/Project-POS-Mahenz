export interface BackdateCashDrawer {
  id: number
  user_id: number
  user_name: string
  shift_id: number | null
  shift_name: string | null
  open_time: string
  opening_balance: number
  total_sales: number
  total_cash_sales: number
  total_expenses: number
  expected_balance: number
  status: string
  open_notes: string | null
  created_by: number | null
  created_by_name: string
}

export interface OpenBackdateCashDrawerPayload {
  user_id: number
  date: string
  shift_id: number | null
  opening_balance: number
  notes?: string
}

export interface CloseBackdateCashDrawerPayload {
  closing_balance: number
  notes?: string
}

export interface CloseBackdateCashDrawerResult {
  expected_balance: number
  closing_balance: number
  difference: number
}

export interface UserOption {
  id: number
  full_name: string
  username: string
}
