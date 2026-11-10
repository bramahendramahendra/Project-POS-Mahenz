import type { PaymentMethod } from '../cashier/cashier.types'

export interface BackdatePaymentPayload {
  transaction_time: string   // HH:mm — combined with backdate kas date on backend
  shift_id?: number
  customer_id?: number
  is_credit: boolean
  balance_used?: number
  device_source: 'web'
  items: Array<{
    product_id: number
    product_name: string
    unit_id?: number
    unit: string
    conversion_qty: number
    quantity: number
    price: number
    subtotal: number
    discount_item?: number
  }>
  subtotal: number
  discount: number
  tax: number
  total_amount: number
  payment_method: PaymentMethod
  payment_amount: number
  change_amount: number
}

export interface BackdateCheckoutResponse {
  id: number
  transaction_code: string
  transaction_date: string
  total_amount: number
}
