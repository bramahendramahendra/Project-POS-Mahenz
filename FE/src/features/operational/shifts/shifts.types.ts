export type ShiftStatus = 'open' | 'closed'

export interface Shift {
  id: number
  kasir_id: number
  kasir_name: string
  opening_balance: number
  closing_balance?: number
  total_transactions: number
  total_revenue: number
  status: ShiftStatus
  notes?: string
  opened_at: string
  closed_at?: string
}

export interface ShiftListFilter {
  page: number
  limit: number
  search?: string
  status?: ShiftStatus | ''
}

export interface OpenShiftPayload {
  opening_balance: number
  notes?: string
}

export interface CloseShiftPayload {
  closing_balance: number
  notes?: string
}
