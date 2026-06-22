export interface Shift {
  id: number
  name: string
  start_time: string
  end_time: string
  is_active: boolean
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
  search?: string
}

export interface ShiftFormPayload {
  name: string
  start_time: string
  end_time: string
}
