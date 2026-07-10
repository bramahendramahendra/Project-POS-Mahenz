import { Printer, ShoppingCart } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { formatRupiah } from '@/shared/utils'

import { useStoreProfileQuery } from '@/features/settings/store'
import type { StoreProfile } from '@/features/settings/store'

import type {
  CartItem,
  CartSummary,
  CheckoutResponse,
  Discount,
  PaymentMethod,
  Tax,
} from '../cashier.types'

interface ReceiptPrintProps {
  open: boolean
  onClose: () => void
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
}

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Kredit',
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function buildReceiptHtml(params: {
  storeProfile?: StoreProfile
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
}): string {
  const { storeProfile, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName } = params
  const storeName = storeProfile?.name || 'POS System'
  const storeSub = [storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • ')
  const change = Math.max(0, amountPaid - summary.grandTotal)

  const itemRows = cart.map((item) => {
    const price = formatRupiah(item.effective_price ?? item.price)
    const discountRow = item.discount_amount && item.discount_amount > 0
      ? `<div class="row muted">
           <span>Diskon ${item.discount_type === 'percent' ? `${item.discount_value}%` : formatRupiah(item.discount_value ?? 0)}</span>
           <span>-${formatRupiah(item.discount_amount)}</span>
         </div>`
      : ''
    return `
      <div class="item">
        <div class="row">
          <span class="item-name">${item.product_name}</span>
          <span class="bold">${formatRupiah(item.subtotal)}</span>
        </div>
        <div class="row muted">
          <span>${item.unit_name} &times; ${item.qty} @ ${price}</span>
        </div>
        ${discountRow}
      </div>`
  }).join('')

  const discountRow = summary.discountAmount > 0
    ? `<div class="row green">
         <span>Diskon${discount.type === 'percent' ? ` (${discount.value}%)` : ''}</span>
         <span>-${formatRupiah(summary.discountAmount)}</span>
       </div>`
    : ''

  const taxRow = summary.taxAmount > 0
    ? `<div class="row muted">
         <span>Pajak (${tax.percent}%)</span>
         <span>+${formatRupiah(summary.taxAmount)}</span>
       </div>`
    : ''

  const customerRow = customerName
    ? `<div class="row"><span class="label">Pelanggan</span><span>${customerName}</span></div>`
    : ''

  return `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8" />
  <title>Struk - ${checkoutData.transaction_code}</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Courier New', monospace;
      font-size: 12px;
      color: #111;
      width: 300px;
      margin: 0 auto;
      padding: 16px 12px;
    }
    .center { text-align: center; }
    .store-name { font-size: 15px; font-weight: 700; letter-spacing: 0.5px; }
    .store-sub { font-size: 11px; color: #555; margin-top: 2px; }
    .divider { border: none; border-top: 1px dashed #aaa; margin: 10px 0; }
    .divider-solid { border: none; border-top: 1px solid #ccc; margin: 8px 0; }
    .row { display: flex; justify-content: space-between; margin-bottom: 3px; }
    .label { color: #777; }
    .muted { color: #777; font-size: 11px; }
    .green { color: #16a34a; }
    .bold { font-weight: 700; }
    .item { margin-bottom: 8px; }
    .item-name { font-weight: 600; }
    .total-row { display: flex; justify-content: space-between; font-size: 14px; font-weight: 700; margin: 4px 0; }
    .kembalian { display: flex; justify-content: space-between; font-weight: 600; color: #16a34a; margin-top: 2px; }
    .footer { text-align: center; color: #888; font-size: 11px; margin-top: 8px; }
    @media print {
      body { width: 100%; }
    }
  </style>
</head>
<body>
  <div class="center">
    <div class="store-name">${storeName}</div>
    ${storeSub ? `<div class="store-sub">${storeSub}</div>` : ''}
  </div>

  <hr class="divider" />

  <div class="row">
    <span class="label">No. Transaksi</span>
    <span class="bold">${checkoutData.transaction_code}</span>
  </div>
  <div class="row">
    <span class="label">Tanggal</span>
    <span>${formatDate(checkoutData.transaction_date)}</span>
  </div>
  ${customerRow}
  <div class="row">
    <span class="label">Pembayaran</span>
    <span>${PAYMENT_LABELS[paymentMethod]}</span>
  </div>

  <hr class="divider" />

  ${itemRows}

  <hr class="divider" />

  <div class="row muted">
    <span>Subtotal</span>
    <span>${formatRupiah(summary.subtotal)}</span>
  </div>
  ${discountRow}
  ${taxRow}

  <hr class="divider-solid" />

  <div class="total-row">
    <span>TOTAL</span>
    <span>${formatRupiah(summary.grandTotal)}</span>
  </div>
  <div class="row muted">
    <span>Dibayar (${PAYMENT_LABELS[paymentMethod]})</span>
    <span>${formatRupiah(amountPaid)}</span>
  </div>
  <div class="kembalian">
    <span>Kembalian</span>
    <span>${formatRupiah(change)}</span>
  </div>

  <hr class="divider" />

  <div class="footer">Terima kasih telah berbelanja!</div>

  <script>
    window.onload = function() {
      window.print()
      window.onafterprint = function() { window.close() }
    }
  </script>
</body>
</html>`
}

export function ReceiptPrint({
  open,
  onClose,
  checkoutData,
  cart,
  summary,
  discount,
  tax,
  paymentMethod,
  amountPaid,
  customerName,
}: ReceiptPrintProps) {
  const change = amountPaid - summary.grandTotal
  const { data: storeProfile } = useStoreProfileQuery()

  const handlePrint = () => {
    const html = buildReceiptHtml({
      storeProfile, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName,
    })
    const win = window.open('', '_blank', 'width=380,height=600,toolbar=0,menubar=0,scrollbars=1')
    if (!win) return
    win.document.write(html)
    win.document.close()
  }

  return (
    <Dialog open={open} onOpenChange={(val) => { if (!val) onClose() }}>
      <DialogContent className="flex flex-col gap-0 p-0 max-w-sm">
        <DialogHeader className="border-b px-6 py-4">
          <DialogTitle>Struk Transaksi</DialogTitle>
          <DialogDescription className="sr-only">Detail struk transaksi</DialogDescription>
        </DialogHeader>

        {/* Receipt preview */}
        <ScrollArea style={{ maxHeight: '65vh' }}>
        <div className="px-6 py-5 space-y-4">

          {/* Toko header */}
          <div className="text-center space-y-0.5">
            <p className="text-base font-bold text-gray-900 tracking-wide">
              {storeProfile?.name || 'POS System'}
            </p>
            {[storeProfile?.address, storeProfile?.phone].filter(Boolean).length > 0 && (
              <p className="text-sm text-gray-500">
                {[storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • ')}
              </p>
            )}
          </div>

          <hr className="border-dashed border-gray-300" />

          {/* Info transaksi */}
          <div className="space-y-1 text-xs text-gray-600">
            <div className="flex justify-between">
              <span className="text-gray-400">No. Transaksi</span>
              <span className="font-medium text-gray-800">{checkoutData.transaction_code}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Tanggal</span>
              <span className="text-gray-700">{formatDate(checkoutData.transaction_date)}</span>
            </div>
            {customerName && (
              <div className="flex justify-between">
                <span className="text-gray-400">Pelanggan</span>
                <span className="text-gray-700">{customerName}</span>
              </div>
            )}
            <div className="flex justify-between">
              <span className="text-gray-400">Pembayaran</span>
              <span className="text-gray-700">{PAYMENT_LABELS[paymentMethod]}</span>
            </div>
          </div>

          <hr className="border-dashed border-gray-300" />

          {/* Item list */}
          <div className="space-y-2.5">
            {cart.map((item) => (
              <div key={`${item.product_id}-${item.unit_id}`} className="space-y-0.5">
                <div className="flex justify-between items-start gap-2">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-800 truncate">{item.product_name}</p>
                    <p className="text-xs text-gray-400">
                      {item.unit_name} &times; {item.qty} &nbsp;@&nbsp;
                      {formatRupiah(item.effective_price ?? item.price)}
                    </p>
                  </div>
                  <span className="text-sm font-semibold text-gray-800 shrink-0">
                    {formatRupiah(item.subtotal)}
                  </span>
                </div>
                {item.discount_amount && item.discount_amount > 0 && (
                  <div className="flex justify-between text-xs text-red-500 pl-2">
                    <span>
                      Diskon{' '}
                      {item.discount_type === 'percent'
                        ? `${item.discount_value}%`
                        : formatRupiah(item.discount_value ?? 0)}
                    </span>
                    <span>-{formatRupiah(item.discount_amount)}</span>
                  </div>
                )}
              </div>
            ))}
          </div>

          <hr className="border-dashed border-gray-300" />

          {/* Summary */}
          <div className="space-y-1.5 text-sm">
            <div className="flex justify-between text-gray-600">
              <span>Subtotal</span>
              <span>{formatRupiah(summary.subtotal)}</span>
            </div>
            {summary.discountAmount > 0 && (
              <div className="flex justify-between text-green-600">
                <span>Diskon{discount.type === 'percent' ? ` (${discount.value}%)` : ''}</span>
                <span>-{formatRupiah(summary.discountAmount)}</span>
              </div>
            )}
            {summary.taxAmount > 0 && (
              <div className="flex justify-between text-gray-600">
                <span>Pajak ({tax.percent}%)</span>
                <span>+{formatRupiah(summary.taxAmount)}</span>
              </div>
            )}
            <hr className="border-gray-200" />
            <div className="flex justify-between font-bold text-base text-gray-900">
              <span>TOTAL</span>
              <span>{formatRupiah(summary.grandTotal)}</span>
            </div>
            <div className="flex justify-between text-gray-600">
              <span>Dibayar</span>
              <span>{formatRupiah(amountPaid)}</span>
            </div>
            <div className="flex justify-between font-semibold">
              <span>Kembalian</span>
              <span className="text-green-600">{formatRupiah(Math.max(0, change))}</span>
            </div>
          </div>

          <div className="text-center pt-1">
            <p className="text-xs text-gray-400">Terima kasih telah berbelanja!</p>
          </div>
        </div>
        </ScrollArea>

        <DialogFooter className="border-t px-6 py-4">
          <Button variant="outline" onClick={handlePrint} className="gap-1.5">
            <Printer size={14} />
            Cetak
          </Button>
          <Button onClick={onClose} className="gap-1.5">
            <ShoppingCart size={14} />
            Transaksi Baru
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
