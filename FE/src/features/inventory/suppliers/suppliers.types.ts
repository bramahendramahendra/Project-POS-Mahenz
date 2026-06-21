export interface Supplier {
  id: number
  supplier_code: string
  name: string
  address: string
  phone: string
  email: string
  contact_person: string
  notes: string
  is_active: boolean
  created_at: string
}

export interface SupplierPurchaseItem {
  id: number
  purchase_code: string
  purchase_date: string
  total_amount: number
  payment_status: string
  remaining_amount: number
}

export interface SupplierReturnHistoryItem {
  id: number
  return_code: string
  return_date: string
  total_return: number
  reason: string
  status: string
}

export interface SupplierDetail {
  id: number
  supplier_code: string
  name: string
  contact_person: string
  phone: string
  email: string
  address: string
  notes: string
  is_active: boolean
  total_purchases: number
  total_amount: number
  total_debt: number
  total_return: number
  purchase_history: SupplierPurchaseItem[]
  return_history: SupplierReturnHistoryItem[]
}

export interface SupplierListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface CreateSupplierPayload {
  name: string
  address?: string
  phone?: string
  email?: string
  contact_person?: string
  notes?: string
}

export interface UpdateSupplierPayload {
  name: string
  address?: string
  phone?: string
  email?: string
  contact_person?: string
  notes?: string
}
