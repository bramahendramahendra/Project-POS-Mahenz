export type ExpirySeverity = 'near' | 'expired'

export interface ExpiryWarning {
  id: number
  product_id: number
  product_name: string
  qty: number
  expired_date: string
  severity: ExpirySeverity
  days_left: number
}

export interface ProductExpirySeverity {
  product_id: number
  severity: ExpirySeverity
  warning_count: number
}

export interface ExpiryWarningsResponse {
  warnings: ExpiryWarning[] | null
  product_severity: ProductExpirySeverity[] | null
}

export type ExpiryBatchStatus = 'active' | 'cleared' | 'written_off'

export interface ExpiryBatchHistory {
  id: number
  qty: number
  expired_date: string
  status: ExpiryBatchStatus
}
