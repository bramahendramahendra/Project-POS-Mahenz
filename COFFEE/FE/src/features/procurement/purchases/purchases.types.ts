export type PaymentStatus = 'paid' | 'unpaid' | 'partial'

export type PurchaseStatus = 'active' | 'void'

/** Jenis baris pembelian (keputusan 92): product = barang jadi, ingredient = bahan baku. */
export type PurchaseItemType = 'product' | 'ingredient'

export const PURCHASE_ITEM_TYPE_LABEL: Record<PurchaseItemType, string> = {
  product: 'Barang jadi',
  ingredient: 'Bahan baku',
}

export interface SupplierPurchaseItem {
  id: number
  item_type: PurchaseItemType
  product_id: number
  product_name: string
  ingredient_id?: number
  ingredient_name?: string
  package_id?: number
  quantity: number
  unit: string
  conversion_qty: number
  purchase_price: number
  subtotal: number
}

export interface PurchasePayment {
  id: number
  payment_date: string
  amount: number
  payment_method: string
  notes: string
  user_name: string
  created_at: string
}

export interface SupplierPurchase {
  id: number
  purchase_code: string
  invoice_number: string
  purchase_date: string
  supplier_id: number
  supplier_name: string
  discount_amount: number
  total_amount: number
  paid_amount: number
  remaining_amount: number
  payment_status: PaymentStatus
  status: PurchaseStatus
  user_name: string
  notes?: string
  items: SupplierPurchaseItem[]
}

export interface SupplierPurchaseFilter {
  search?: string
  start_date?: string
  end_date?: string
  supplier_id?: number
  payment_status?: PaymentStatus
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  page: number
  limit: number
}

export interface SupplierPurchasePayment {
  amount: number
  payment_date: string
  payment_method: string
  notes?: string
}

export interface CreatePurchaseItemPayload {
  item_type: PurchaseItemType
  product_id?: number
  ingredient_id?: number
  package_id?: number
  quantity: number
  purchase_price: number
  unit: string
  conversion_qty: number
}

export interface CreateSupplierPurchasePayload {
  purchase_date: string
  invoice_number: string
  supplier_id: number
  items: CreatePurchaseItemPayload[]
  discount_amount: number
  notes?: string
  payment_status: PaymentStatus
  paid_amount: number
  payment_method?: string
}

export type UpdateSupplierPurchasePayload = CreateSupplierPurchasePayload & { id: number }

export interface AddPurchaseItemsPayload {
  id: number
  items: CreatePurchaseItemPayload[]
}
