import type { DiscountType, OrderType, PaymentMethod } from '@/features/sales/cashier'

export type { DiscountType, PaymentMethod }

export interface TransactionItemModifier {
  option_id: number | null
  group_name: string
  option_name: string
  price: number
}

export interface TransactionItem {
  id: number
  menu_item_id: number | null
  variant_id: number | null
  /** Terisi = barang jadi (stok dipotong); null = racikan. */
  product_id: number | null
  product_name: string // nama menu saat dijual
  variant_name: string
  notes: string
  unit: string
  conversion_qty: number
  quantity: number
  price: number
  subtotal: number
  discount_item: number
  modifiers: TransactionItemModifier[]
}

export interface Transaction {
  id: number
  transaction_code: string
  user_id: number
  kasir_name: string
  customer_id?: number
  customer_name?: string
  shift_id?: number
  transaction_date: string
  items: TransactionItem[]
  subtotal: number
  discount: number
  tax: number
  total_amount: number
  payment_method: PaymentMethod
  payment_amount: number
  change_amount: number
  is_credit: boolean
  status: 'completed' | 'void'
  device_source: string
  order_type: OrderType | ''
  table_id: number | null
  table_name: string
  queue_number: number
  queue_label: string
  order_name: string
}

export interface TransactionListFilter {
  search?: string
  date_from?: string
  date_to?: string
  payment_method?: PaymentMethod | ''
  status?: 'completed' | 'void' | ''
  order_type?: OrderType | ''
  page: number
  limit: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/** Pilihan per item barang jadi saat void (keputusan 100). */
export interface VoidTransactionPayload {
  id: number
  items: { item_id: number; action: 'return' | 'waste' }[]
  reason?: string
}
