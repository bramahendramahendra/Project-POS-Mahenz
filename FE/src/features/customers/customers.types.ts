export interface Customer {
  id: number
  customer_code: string
  name: string
  phone: string
  address: string
  credit_limit: number
  notes?: string
  is_active: boolean
  created_at: string
}

export interface CustomerOption {
  id: number
  name: string
  customer_code: string
  credit_limit: number
}

export interface CustomerListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
}

export interface CreateCustomerPayload {
  name: string
  phone?: string
  address?: string
  credit_limit?: number
  notes?: string
}

export type UpdateCustomerPayload = Partial<CreateCustomerPayload>
