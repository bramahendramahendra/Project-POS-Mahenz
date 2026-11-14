/**
 * Kirim ringkasan struk lewat WhatsApp. Tidak mencetak ke kertas — ini jalur
 * "bagikan e-nota" ke pelanggan. Menyusun teks struk polos lalu membuka
 * https://wa.me/ dengan teks tersebut (jalan di desktop maupun tablet/HP).
 */
import { formatRupiah } from '@/shared/utils'

import type { ReceiptData } from './blePrinter'
import type { PaymentMethod } from './cashier.types'
import { cartItemDetails, cartItemName, cartItemQtyLine, orderInfoRows } from './cashier.utils'

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Hutang',
  balance: 'Saldo',
}

/** Susun teks ringkasan struk (polos) untuk dikirim via WhatsApp. */
export function buildReceiptText(data: ReceiptData): string {
  const change = Math.max(0, data.amountPaid - data.summary.grandTotal)
  const tanggal = new Date(data.checkoutData.transaction_date).toLocaleString('id-ID')
  const lines: string[] = []

  lines.push(`*${data.storeName}*`)
  if (data.storeSub) lines.push(data.storeSub)
  lines.push('--------------------------------')
  if (data.checkoutData.queue_label) lines.push(`*Antrean ${data.checkoutData.queue_label}*`)
  lines.push(`No    : ${data.checkoutData.transaction_code}`)
  lines.push(`Tgl   : ${tanggal}`)
  if (data.customerName) lines.push(`Plg   : ${data.customerName}`)
  for (const [label, value] of orderInfoRows(data.checkoutData)) lines.push(`${label.padEnd(6)}: ${value}`)
  if (data.unpaid) lines.push('*TAGIHAN - BELUM DIBAYAR*')
  else lines.push(`Bayar : ${PAYMENT_LABELS[data.paymentMethod]}`)
  lines.push('--------------------------------')

  for (const item of data.cart) {
    lines.push(cartItemName(item))
    for (const d of cartItemDetails(item)) lines.push(`  ${d}`)
    lines.push(`  ${cartItemQtyLine(item)} = ${formatRupiah(item.subtotal)}`)
  }
  lines.push('--------------------------------')

  lines.push(`Subtotal : ${formatRupiah(data.summary.subtotal)}`)
  if (data.summary.discountAmount > 0) {
    lines.push(`Diskon   : -${formatRupiah(data.summary.discountAmount)}`)
  }
  if (data.summary.taxAmount > 0) {
    lines.push(`Pajak    : +${formatRupiah(data.summary.taxAmount)}`)
  }
  lines.push(`*TOTAL*  : *${formatRupiah(data.summary.grandTotal)}*`)
  if (!data.unpaid) {
    lines.push(`Dibayar  : ${formatRupiah(data.amountPaid)}`)
    lines.push(`Kembalian: ${formatRupiah(change)}`)
  }
  lines.push('--------------------------------')
  if (data.footer) lines.push(data.footer)

  return lines.join('\n')
}

/**
 * Buka WhatsApp dengan teks ringkasan struk. `phone` opsional (format
 * internasional tanpa + / spasi, mis. "628123456789"); bila kosong, WhatsApp
 * membuka pemilih kontak agar user memilih tujuan sendiri.
 */
export function shareReceiptViaWhatsApp(data: ReceiptData, phone?: string): void {
  const text = encodeURIComponent(buildReceiptText(data))
  const normalizedPhone = (phone ?? '').replace(/[^0-9]/g, '')
  const base = normalizedPhone ? `https://wa.me/${normalizedPhone}` : 'https://wa.me/'
  const url = `${base}?text=${text}`
  window.open(url, '_blank')
}
