export interface ReconListItem {
  product_id: number
  product_code: string
  product_name: string
  base_unit: string
  old_stock: number
  old_stock_available: boolean
  new_stock: number
  diff: number
  needs_stock_review: boolean
  stock_review_note: string
}

export interface ReconSummary {
  total_products: number
  needs_review: number
  matched: number
  backup_available: boolean
}

export interface ReconBreakdown {
  purchase_in: number
  purchase_void: number
  sale_out: number
  sale_void: number
  supplier_return: number
  expired: number
  adjustment: number
  computed_new_stock: number
}

export interface ReconPackageLevel {
  package_id: number
  unit_name: string
  package_name: string
  is_default: boolean
  factor_to_base: number
  current_stock: number
}

export interface ReconKartuStokRow {
  id: number
  mutation_type: string
  quantity: number
  stock_before: number
  stock_after: number
  reference_type: string
  reference_id: number
  notes: string
  user_name: string
  created_at: string
}

export interface ReconDetail {
  product_id: number
  product_code: string
  product_name: string
  base_unit: string
  old_stock: number
  old_stock_available: boolean
  new_stock: number
  diff: number
  needs_stock_review: boolean
  stock_review_note: string
  breakdown: ReconBreakdown
  packages: ReconPackageLevel[]
  kartu_stok: ReconKartuStokRow[]
}

export interface ReconFilter {
  search?: string
  category_id?: number
  only_review?: boolean
}

export interface ReconListFilter extends ReconFilter {
  page: number
  limit: number
}

export interface AdjustLevelPayload {
  package_id: number
  new_stock: number
}

export interface AdjustPayload {
  product_id: number
  levels: AdjustLevelPayload[]
  note?: string
}

export interface MarkReviewedPayload {
  product_id: number
  note?: string
}
