/**
 * Membuat struk sebagai file PDF di sisi BROWSER (tanpa server, tanpa Bluetooth,
 * tanpa RawBT). Jalur ini jalan di HTTP maupun HTTPS, di tablet maupun desktop.
 *
 * Kegunaan:
 *  - Simpan/arsip e-nota.
 *  - Bagikan ke pelanggan (mis. WhatsApp).
 *  - Dari file PDF, user bisa mencetak ke printer thermal lewat aplikasi bawaan
 *    tablet atau aplikasi printer (termasuk membuka PDF di RawBT/printer app).
 *
 * Ukuran halaman mengikuti lebar kertas (58mm atau 80mm) dengan tinggi otomatis
 * menyesuaikan jumlah baris, sehingga tidak ada area kosong berlebih.
 */
import { jsPDF } from 'jspdf'

import { formatRupiah } from '@/shared/utils'

import type { ReceiptData } from './blePrinter'
import type { PaymentMethod } from './cashier.types'

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Hutang',
  balance: 'Saldo',
}

// Dimensi dalam milimeter.
const PAGE_WIDTH: Record<'58mm' | '80mm', number> = { '58mm': 58, '80mm': 80 }
const MARGIN_X = 3 // margin kiri-kanan (mm)
const LINE_H = 4 // tinggi baris (mm)
const FONT_SIZE = 8 // pt — pas untuk thermal 58mm
const FONT_SIZE_TITLE = 10

/** Bangun dokumen PDF struk dan kembalikan instance jsPDF. */
export function buildReceiptPdf(data: ReceiptData): jsPDF {
  const pageWidth = PAGE_WIDTH[data.paperSize]
  const contentWidth = pageWidth - MARGIN_X * 2
  const change = Math.max(0, data.amountPaid - data.summary.grandTotal)

  // Dokumen sementara untuk menghitung tinggi; jsPDF butuh ukuran halaman di awal,
  // jadi kita hitung dulu total baris lalu buat dokumen dengan tinggi final.
  const lines: Array<{ text: string; align?: 'left' | 'center'; bold?: boolean; right?: string; size?: number }> = []

  const pushCenter = (text: string, opts?: { bold?: boolean; size?: number }) =>
    lines.push({ text, align: 'center', bold: opts?.bold, size: opts?.size })
  const pushRow = (left: string, right: string, opts?: { bold?: boolean }) =>
    lines.push({ text: left, right, bold: opts?.bold })
  const pushLeft = (text: string) => lines.push({ text })
  const pushSep = () => lines.push({ text: '--------------------------------' })

  pushCenter(data.storeName, { bold: true, size: FONT_SIZE_TITLE })
  if (data.storeSub) pushCenter(data.storeSub)
  pushSep()
  pushRow('No', data.checkoutData.transaction_code)
  pushRow('Tgl', new Date(data.checkoutData.transaction_date).toLocaleString('id-ID'))
  if (data.customerName) pushRow('Plg', data.customerName)
  pushRow('Bayar', PAYMENT_LABELS[data.paymentMethod])
  pushSep()

  for (const item of data.cart) {
    const price = item.effective_price ?? item.price
    pushLeft(item.product_name)
    pushRow(`  ${item.unit_name} x${item.qty} @ ${formatRupiah(price)}`, formatRupiah(item.subtotal))
  }
  pushSep()

  pushRow('Subtotal', formatRupiah(data.summary.subtotal))
  if (data.summary.discountAmount > 0) {
    const label = `Diskon${data.discount.type === 'percent' ? ` (${data.discount.value}%)` : ''}`
    pushRow(label, `-${formatRupiah(data.summary.discountAmount)}`)
  }
  if (data.summary.taxAmount > 0) {
    pushRow(`Pajak (${data.tax.percent}%)`, `+${formatRupiah(data.summary.taxAmount)}`)
  }
  pushRow('TOTAL', formatRupiah(data.summary.grandTotal), { bold: true })
  pushRow('Dibayar', formatRupiah(data.amountPaid))
  pushRow('Kembalian', formatRupiah(change))
  pushSep()
  if (data.footer) pushCenter(data.footer)

  const topBottomMargin = 4
  const pageHeight = topBottomMargin * 2 + lines.length * LINE_H

  const doc = new jsPDF({ unit: 'mm', format: [pageWidth, pageHeight] })
  doc.setFont('courier', 'normal')

  let y = topBottomMargin + LINE_H
  for (const ln of lines) {
    doc.setFontSize(ln.size ?? FONT_SIZE)
    doc.setFont('courier', ln.bold ? 'bold' : 'normal')

    if (ln.align === 'center') {
      doc.text(ln.text, pageWidth / 2, y, { align: 'center' })
    } else if (ln.right !== undefined) {
      doc.text(ln.text, MARGIN_X, y)
      doc.text(ln.right, pageWidth - MARGIN_X, y, { align: 'right' })
    } else {
      // Teks kiri; bungkus bila melebihi lebar konten.
      const wrapped = doc.splitTextToSize(ln.text, contentWidth) as string[]
      for (let i = 0; i < wrapped.length; i++) {
        if (i > 0) y += LINE_H
        doc.text(wrapped[i], MARGIN_X, y)
      }
    }
    y += LINE_H
  }

  return doc
}

/** Nama file PDF struk berdasarkan kode transaksi. */
function receiptFileName(data: ReceiptData): string {
  const code = data.checkoutData.transaction_code || 'struk'
  return `struk-${code}.pdf`
}

/**
 * Simpan struk sebagai file PDF (unduh). Di tablet/HP, file ini bisa dibuka lalu
 * dicetak ke printer thermal lewat aplikasi bawaan atau aplikasi printer.
 */
export function saveReceiptPdf(data: ReceiptData): void {
  const doc = buildReceiptPdf(data)
  doc.save(receiptFileName(data))
}

/**
 * Buka struk PDF di tab baru (untuk pratinjau / bagikan / cetak dari viewer PDF).
 * Mengembalikan false bila popup diblokir.
 */
export function openReceiptPdf(data: ReceiptData): boolean {
  const doc = buildReceiptPdf(data)
  const url = doc.output('bloburl')
  const win = window.open(url as unknown as string, '_blank')
  return !!win
}
