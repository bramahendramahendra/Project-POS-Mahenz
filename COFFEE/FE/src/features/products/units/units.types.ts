export interface Unit {
  id: number
  name: string
  abbreviation: string
  is_active: boolean
  created_at: string
}

export interface UnitOption {
  id: number
  name: string
  abbreviation: string
}

export interface UnitListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface CreateUnitPayload {
  name: string
  abbreviation: string
}

export type UpdateUnitPayload = Partial<CreateUnitPayload>
