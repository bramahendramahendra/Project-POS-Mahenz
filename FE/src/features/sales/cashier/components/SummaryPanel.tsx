import { AlertTriangle, Receipt } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'

import { useCashierStore } from '../cashier.store'
import { calcCartSummary } from '../cashier.utils'
import { CustomerSelector } from './CustomerSelector'
import { DiscountInput } from './DiscountInput'
import { TaxInput } from './TaxInput'

export function SummaryPanel() {
  const { cart, discount, tax, openPaymentModal } = useCashierStore()
  const summary = calcCartSummary(cart, discount, tax)

  const { data: currentDrawer } = useCashDrawerCurrentQuery()
  const kasOpen = currentDrawer?.status === 'open'

  return (
    <div className="flex h-full flex-col bg-white border-l border-gray-200">
      {/* Header */}
      <div className="flex items-center gap-2 px-4 py-3 border-b shrink-0">
        <Receipt size={16} className="text-gray-500" />
        <span className="text-sm font-semibold text-gray-800">Ringkasan</span>
      </div>

      <CustomerSelector />

      {/* Read-only item list — 1 baris per item, compact */}
      <ScrollArea className="flex-1">
        {cart.length === 0 ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-sm text-gray-400">Keranjang kosong</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-100 px-4">
            {cart.map((item) => {
              const displayPrice = item.effective_price ?? item.price
              const hasDiscount = !!(item.discount_type && item.discount_value && item.discount_value > 0)
              return (
                <li
                  key={`${item.product_id}-${item.unit_id}`}
                  className="flex items-center gap-2 py-2"
                >
                  {/* Nama + unit */}
                  <div className="flex-1 min-w-0">
                    <span className="text-sm text-gray-800 truncate block">
                      {item.product_name}
                    </span>
                    <span className="text-xs text-gray-400">
                      {item.unit_name}
                      {hasDiscount && (
                        <span className="ml-1 text-red-500">
                          · {item.discount_type === 'percent'
                            ? `disc ${item.discount_value}%`
                            : `-${formatRupiah(item.discount_value ?? 0)}`}
                        </span>
                      )}
                    </span>
                  </div>

                  {/* Qty */}
                  <span className="text-xs text-gray-500 shrink-0 w-8 text-center">
                    {item.qty}×
                  </span>

                  {/* Harga satuan */}
                  <span className="text-xs text-gray-500 shrink-0 w-20 text-right">
                    {formatRupiah(displayPrice)}
                  </span>

                  {/* Subtotal */}
                  <span className="text-sm font-semibold text-gray-800 shrink-0 w-24 text-right">
                    {formatRupiah(item.subtotal)}
                  </span>
                </li>
              )
            })}
          </ul>
        )}
      </ScrollArea>

      {/* Diskon & Pajak */}
      {cart.length > 0 && (
        <div className="border-t shrink-0">
          <DiscountInput />
          <TaxInput />
        </div>
      )}

      {/* Summary total */}
      <div className="border-t px-4 py-3 space-y-1.5 bg-gray-50 text-sm shrink-0">
        <div className="flex justify-between text-gray-600">
          <span>Subtotal</span>
          <span>{formatRupiah(summary.subtotal)}</span>
        </div>
        {summary.discountAmount > 0 && (
          <div className="flex justify-between text-green-600">
            <span>
              Diskon{discount.type === 'percent' ? ` (${discount.value}%)` : ''}
            </span>
            <span>-{formatRupiah(summary.discountAmount)}</span>
          </div>
        )}
        {summary.taxAmount > 0 && (
          <div className="flex justify-between text-gray-600">
            <span>Pajak ({tax.percent}%)</span>
            <span>+{formatRupiah(summary.taxAmount)}</span>
          </div>
        )}
        <div className="flex justify-between border-t pt-2 font-bold text-gray-900 text-base">
          <span>TOTAL</span>
          <span>{formatRupiah(summary.grandTotal)}</span>
        </div>
      </div>

      {/* Warning kas belum buka */}
      {!kasOpen && (
        <div className="mx-4 mb-2 flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700 shrink-0">
          <AlertTriangle size={13} className="shrink-0 mt-0.5" />
          <span>Kas belum dibuka. Buka kas terlebih dahulu untuk memproses transaksi.</span>
        </div>
      )}

      {/* Tombol Bayar */}
      <div className="px-4 py-3 border-t shrink-0">
        <Button
          className="w-full gap-2"
          onClick={openPaymentModal}
          disabled={cart.length === 0 || !kasOpen}
        >
          💳 Bayar
        </Button>
      </div>
    </div>
  )
}
