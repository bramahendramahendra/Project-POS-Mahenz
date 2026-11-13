export interface Shift {
  id: number
  name: string
  start_time: string
  end_time: string
  is_active: boolean
  created_at: string
}

export interface ShiftOption {
  id: number
  name: string
  start_time: string
  end_time: string
}

export interface ShiftListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface ShiftFormPayload {
  name: string
  start_time: string
  end_time: string
}
