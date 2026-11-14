import { savedReceiptLine } from '@/features/sales/cashier/cashier.utils'
import type { CartItem, CartSummary, CheckoutResponse, Discount, Tax } from '@/features/sales/cashier'

import type { PaymentMethod, Transaction } from './transactions.types'

export const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Kredit',
  balance: 'Saldo',
}

export function formatDateTimeShort(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDateTimeLong(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** Data struk cetak ulang dari transaksi tersimpan (nomor antrean, info pesanan, rincian item). */
export function transactionReceipt(t: Transaction): {
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
} {
  const cart = t.items.map((item) =>
    savedReceiptLine({
      key: String(item.id),
      name: item.product_name,
      variant_name: item.variant_name,
      qty: item.quantity,
      price: item.price,
      subtotal: item.subtotal,
      discount: item.discount_item,
      notes: item.notes,
      is_stock_item: item.product_id !== null,
      modifiers: item.modifiers.map((m) => ({ option_id: m.option_id ?? 0, group_name: m.group_name, option_name: m.option_name, price: m.price })),
    }),
  )
  return {
    checkoutData: {
      id: t.id,
      transaction_code: t.transaction_code,
      total_amount: t.total_amount,
      payment_amount: t.payment_amount,
      change_amount: t.change_amount,
      transaction_date: t.transaction_date,
      queue_label: t.queue_label,
      order_type: t.order_type,
      table_name: t.table_name,
      order_name: t.order_name,
    },
    cart,
    summary: { subtotal: t.subtotal, discountAmount: t.discount, taxAmount: t.tax, grandTotal: t.total_amount },
    discount: { type: 'amount', value: t.discount, amount: t.discount },
    tax: { percent: 0, amount: t.tax },
  }
}
