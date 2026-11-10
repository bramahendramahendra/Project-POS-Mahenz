export interface CapitalReconListItem {
  review_id: number
  transaction_item_id: number
  transaction_id: number
  transaction_code: string
  transaction_date: string
  product_id: number
  product_name: string
  unit: string
  quantity: number
  sell_price: number
  old_purchase_price: number
  suggested_price: number
  reason: string
  status: 'pending' | 'resolved' | 'skipped'
  resolved_price?: number
  note: string
  resolved_by?: string
  resolved_at?: string
}

export interface CapitalReconSummary {
  pending: number
  resolved: number
  skipped: number
  total: number
}

export interface CapitalReconFilter {
  search?: string
  status?: 'pending' | 'resolved' | 'skipped' | 'all'
}

export interface CapitalReconListFilter extends CapitalReconFilter {
  page: number
  limit: number
}

export interface ResolvePayload {
  review_id: number
  correct_cost: number
  note?: string
}

export interface SkipPayload {
  review_id: number
  note?: string
}
