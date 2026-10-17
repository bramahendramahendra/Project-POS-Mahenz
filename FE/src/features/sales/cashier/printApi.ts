import { api } from '@/services'

import type { ReceiptData } from './blePrinter'

/**
 * Payload cetak struk ke BE, yang akan mengirim ESC/POS langsung ke printer
 * thermal via COM port (jalur paling tajam & tidak terpotong untuk desktop).
 */
interface PrintReceiptPayload {
  store_name: string
  store_sub: string
  footer: string
  transaction_code: string
  transaction_date: string
  customer_name: string
  payment_method: string
  items: Array<{
    name: string
    unit_name: string
    qty: number
    price: number
    subtotal: number
  }>
  subtotal: number
  discount_amount: number
  tax_amount: number
  tax_percent: number
  grand_total: number
  amount_paid: number
  change: number
}

function toPayload(data: ReceiptData): PrintReceiptPayload {
  const change = Math.max(0, data.amountPaid - data.summary.grandTotal)
  return {
    store_name: data.storeName,
    store_sub: data.storeSub ?? '',
    footer: data.footer ?? '',
    transaction_code: data.checkoutData.transaction_code,
    transaction_date: data.checkoutData.transaction_date,
    customer_name: data.customerName ?? '',
    payment_method: data.paymentMethod,
    items: data.cart.map((i) => ({
      name: i.product_name,
      unit_name: i.unit_name,
      qty: i.qty,
      price: i.effective_price ?? i.price,
      subtotal: i.subtotal,
    })),
    subtotal: data.summary.subtotal,
    discount_amount: data.summary.discountAmount,
    tax_amount: data.summary.taxAmount,
    tax_percent: data.tax.percent,
    grand_total: data.summary.grandTotal,
    amount_paid: data.amountPaid,
    change,
  }
}

/**
 * Kirim struk ke BE untuk dicetak langsung ke printer thermal (COM port).
 * Melempar error bila BE gagal mencetak — pemanggil bisa menyiapkan fallback.
 */
export async function printViaSerialApi(data: ReceiptData): Promise<void> {
  await api.post('/print/receipt', toPayload(data))
}
