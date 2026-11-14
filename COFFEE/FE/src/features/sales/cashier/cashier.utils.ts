import { formatRupiah } from '@/shared/utils'

import type { CartItem, CartSummary, Discount, OrderInfo, Tax } from './cashier.types'
import { ORDER_TYPE_LABEL } from './cashier.types'

export const calcItemSubtotal = (qty: number, price: number): number => {
  return qty * price
}

export const calcSubtotal = (items: CartItem[]): number => {
  return items.reduce((sum, item) => sum + item.subtotal, 0)
}

export const calcDiscountAmount = (
  subtotal: number,
  discount: Omit<Discount, 'amount'>
): number => {
  if (discount.type === 'none' || discount.value <= 0) return 0
  if (discount.type === 'percent') {
    return Math.round((subtotal * Math.min(discount.value, 100)) / 100)
  }
  return Math.min(discount.value, subtotal)
}

export const calcTaxAmount = (
  subtotal: number,
  discountAmount: number,
  taxPercent: number
): number => {
  if (taxPercent <= 0) return 0
  const taxableAmount = subtotal - discountAmount
  return Math.round((taxableAmount * taxPercent) / 100)
}

export const calcGrandTotal = (
  subtotal: number,
  discountAmount: number,
  taxAmount: number
): number => {
  return Math.max(0, subtotal - discountAmount + taxAmount)
}

export const calcCartSummary = (items: CartItem[], discount: Discount, tax: Tax): CartSummary => {
  const subtotal = calcSubtotal(items)
  const discountAmount = calcDiscountAmount(subtotal, discount)
  const taxAmount = calcTaxAmount(subtotal, discountAmount, tax.percent)
  const grandTotal = calcGrandTotal(subtotal, discountAmount, taxAmount)
  return { subtotal, discountAmount, taxAmount, grandTotal }
}

export const calcChange = (grandTotal: number, amountPaid: number): number => {
  return Math.max(0, amountPaid - grandTotal)
}

export const isPaymentSufficient = (grandTotal: number, amountPaid: number): boolean => {
  return amountPaid >= grandTotal
}

export function calculateItemDiscount(
  price: number,
  qty: number,
  type: 'percent' | 'nominal',
  value: number
): { discount_amount: number; effective_price: number; subtotal: number } {
  if (value <= 0) {
    return { discount_amount: 0, effective_price: price, subtotal: price * qty }
  }

  if (type === 'percent') {
    const pct = Math.min(value, 100)
    const effective_price = Math.round(price * (1 - pct / 100))
    const discount_amount = (price - effective_price) * qty
    return { discount_amount, effective_price, subtotal: effective_price * qty }
  }

  // nominal: value adalah potongan per unit
  const effective_price = Math.max(0, price - value)
  const discount_amount = (price - effective_price) * qty
  return { discount_amount, effective_price, subtotal: effective_price * qty }
}

// ─── Coffee: baris pesanan dan struk (tahap 4.7b) ─────────────────────────────

/** Kunci baris keranjang: kombinasi + tambahan (urut) + catatan sama = baris sama (jumlah digabung). */
export function cartLineKey(variantId: number, optionIds: number[], notes?: string): string {
  return `${variantId}|${[...optionIds].sort((a, b) => a - b).join(',')}|${(notes ?? '').trim().toLowerCase()}`
}

/** Nama item di struk dan keranjang: "Latte (Regular · Dingin)" atau "Teh". */
export function cartItemName(item: Pick<CartItem, 'product_name' | 'variant_name'>): string {
  return item.variant_name ? `${item.product_name} (${item.variant_name})` : item.product_name
}

/** Baris rincian di bawah nama item: tambahan, catatan. */
export function cartItemDetails(item: Pick<CartItem, 'modifiers' | 'notes'>): string[] {
  const lines = item.modifiers.map((m) => `+ ${m.option_name}${m.price > 0 ? ` ${formatRupiah(m.price)}` : ''}`)
  if (item.notes) lines.push(`Catatan: ${item.notes}`)
  return lines
}

/** "x2 @ Rp 29.000" (dengan satuan bila ada, untuk barang lama). */
export function cartItemQtyLine(item: Pick<CartItem, 'unit_name' | 'qty' | 'price' | 'effective_price'>): string {
  const unit = item.unit_name ? `${item.unit_name} ` : ''
  return `${unit}x${item.qty} @ ${formatRupiah(item.effective_price ?? item.price)}`
}

/** Baris info pesanan (label, nilai) untuk struk: jenis pesanan, meja, nama pemesan. */
export function orderInfoRows(info: OrderInfo): Array<[string, string]> {
  const rows: Array<[string, string]> = []
  if (info.order_type) rows.push(['Pesanan', ORDER_TYPE_LABEL[info.order_type]])
  if (info.table_name) rows.push(['Meja', info.table_name])
  if (info.order_name) rows.push(['Nama', info.order_name])
  return rows
}

/**
 * Baris struk untuk cetak ulang dari data tersimpan (riwayat transaksi, piutang). Harga satuan
 * tersimpan sudah termasuk tambahan; diskon item dalam rupiah.
 */
export function savedReceiptLine(p: {
  key: string
  name: string
  variant_name?: string
  unit?: string
  qty: number
  price: number
  subtotal: number
  discount?: number
  modifiers?: CartItem['modifiers']
  notes?: string
  is_stock_item?: boolean
}): CartItem {
  return {
    key: p.key,
    menu_item_id: 0,
    variant_id: 0,
    product_name: p.name,
    variant_name: p.variant_name ?? '',
    unit_name: p.unit ?? '',
    modifiers: p.modifiers ?? [],
    notes: p.notes || undefined,
    is_stock_item: !!p.is_stock_item,
    qty: p.qty,
    price: p.price,
    subtotal: p.subtotal,
    discount_amount: p.discount && p.discount > 0 ? p.discount : undefined,
  }
}
