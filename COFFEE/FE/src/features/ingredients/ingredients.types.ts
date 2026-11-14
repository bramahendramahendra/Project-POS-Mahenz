// Bahan baku Coffee (tahap 4.3, keputusan 24, 80-84).

export interface Ingredient {
  id: number
  name: string
  unit_id: number
  unit_name: string
  unit_abbr: string
  stock: number
  total_used: number
  total_wasted: number
  min_stock: number | null
  is_low_stock: boolean
  is_active: boolean
  has_history: boolean
}

export interface IngredientOption {
  id: number
  name: string
  unit_abbr: string
  stock: number
}

export interface IngredientListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
  low_stock?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface IngredientPayload {
  name: string
  unit_id: number
  min_stock: number | null
}

/**
 * usage = terpakai, waste = terbuang, cancel = pembatalan entri,
 * purchase = pembelian dari supplier, purchase_void = pembelian di-void atau dikurangi saat edit (tahap 4.5).
 */
export type MutationType = 'usage' | 'waste' | 'cancel' | 'purchase' | 'purchase_void'

export interface IngredientMutation {
  id: number
  ingredient_id: number
  ingredient_name: string
  unit_abbr: string
  type: MutationType
  qty_change: number
  stock_before: number
  stock_after: number
  notes: string
  reference_id: number | null
  is_cancelled: boolean
  cancelled_at?: string
  cancelled_by_name?: string
  cancel_reason?: string
  created_by_name: string
  created_at: string
}

export interface MutationListFilter {
  page: number
  limit: number
  ingredient_id?: number
  type?: MutationType
  start_date?: string
  end_date?: string
}

export interface RecordUsagePayload {
  ingredient_id: number
  qty: number
  notes: string
}

export interface RecordWastePayload {
  ingredient_id: number
  qty: number
  reason: string
}

export interface RecordResult {
  mutation: IngredientMutation
  stock_after: number
  negative_stock: boolean
}
