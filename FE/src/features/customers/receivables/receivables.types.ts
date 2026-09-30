export type ReceivableStatus = 'unpaid' | 'partial' | 'paid' | 'void'

// Selaras dengan BE PaymentResponse (dto_receivable.go).
export interface ReceivablePayment {
  id: number
  payment_date: string
  amount: number
  payment_method: string
  user_name: string
  notes: string
}

// Rincian ringkas barang dari transaksi asal (BE TransactionItemBrief).
export interface ReceivableItem {
  product_name: string
  unit: string
  quantity: number
  price: number
  subtotal: number
  discount_item: number
}

// Baris daftar piutang (BE ReceivableResponse).
export interface Receivable {
  id: number
  transaction_id: number
  transaction_code: string
  customer_name: string
  total_amount: number
  paid_amount: number
  remaining_amount: number
  status: ReceivableStatus
  due_date?: string
  created_at: string
}

// Detail piutang lengkap (BE ReceivableDetailResponse).
export interface ReceivableDetail extends Receivable {
  customer_id: number
  notes: string
  items: ReceivableItem[]
  payments: ReceivablePayment[]
}

export type ReceivableDueTerm = '' | 'overdue' | 'not_due'

export interface ReceivableListFilter {
  page: number
  limit: number
  search: string
  status?: ReceivableStatus | ''
  due_term?: ReceivableDueTerm
}

// Angka agregat untuk kartu ringkasan (BE ReceivableStats).
export interface ReceivableStats {
  total_remaining: number
  customer_count: number
  overdue_count: number
  overdue_remaining: number
}

export interface CreatePaymentPayload {
  amount: number
  payment_date: string
  notes?: string
}
