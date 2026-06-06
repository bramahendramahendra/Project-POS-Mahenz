export interface Unit {
  id: number
  name: string
  abbreviation: string
  is_active: boolean
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
}

export interface ProductPackage {
  id: number
  product_id: number
  unit_id: number
  unit_name: string
  abbreviation: string
  package_name: string
  conversion_qty: number
  purchase_price: number
  selling_price: number
  is_default: boolean
}

export interface PriceTier {
  id: number
  product_id: number
  unit_id: number
  unit_name: string
  tier_name: string
  min_qty: number
  price: number
}

export interface Product {
  id: number
  name: string
  sku?: string
  barcode?: string
  category_id?: number
  category_name?: string
  description?: string
  purchase_price: number
  selling_price: number
  stock: number
  min_stock: number
  unit_id: number
  unit_name: string
  unit_abbreviation: string
  is_active: boolean
  created_at: string
  units: ProductPackage[]
  prices: PriceTier[]
  extra_packages: number
  price_tiers_count: number
}

export interface ProductFilter {
  search?: string
  category_id?: number
  is_active?: boolean
  low_stock?: boolean
  page?: number
  page_size?: number
}

export interface CreateProductPayload {
  name: string
  sku: string
  barcode: string
  category_id: number
  description?: string
  purchase_price: number
  selling_price: number
  stock: number
  min_stock: number
  unit_id: number
  is_active: boolean
}

export type UpdateProductPayload = Partial<CreateProductPayload>

export interface CreateUnitPayload {
  name: string
  abbreviation: string
}
export type UpdateUnitPayload = Partial<CreateUnitPayload>

export interface CreateProductPackagePayload {
  unit_id: number
  package_name?: string
  conversion_qty: number
  purchase_price: number
  selling_price: number
  is_default: boolean
}

export interface CreatePriceTierPayload {
  unit_id: number
  tier_name: string
  min_qty: number
  price: number
}

export interface UpdatePriceTierPayload extends Partial<CreatePriceTierPayload> {
  priceId: number
}
