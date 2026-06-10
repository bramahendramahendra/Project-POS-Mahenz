export interface Category {
  id: number
  name: string
  code: string
  description: string
  is_active: boolean
  product_count: number
  active_product_count: number
  created_at: string
}

export interface CategoryOption {
  id: number
  name: string
}

export interface CategoryListFilter {
  page: number
  limit: number
  search: string
  is_active?: boolean
}

export interface CreateCategoryPayload {
  name: string
  description?: string
}

export type UpdateCategoryPayload = Partial<CreateCategoryPayload>
