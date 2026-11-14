import { useState } from 'react'
import { AlertTriangle, Receipt, Trash2 } from 'lucide-react'

import { ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'

import { useCashierStore } from '../cashier.store'
import { calcCartSummary } from '../cashier.utils'
import { CartItemRow } from './CartItemRow'
import { CustomerSelector } from './CustomerSelector'
import { DiscountInput } from './DiscountInput'
import { OrderInfoPanel } from './OrderInfoPanel'
import { TaxInput } from './TaxInput'

/** Panel pesanan kasir: info pesanan, keranjang yang bisa diubah, diskon/pajak, total, tombol Bayar. */
export function SummaryPanel() {
  const { cart, discount, tax, orderType, table, openPaymentModal, clearCart } = useCashierStore()
  const [clearConfirmOpen, setClearConfirmOpen] = useState(false)
  const summary = calcCartSummary(cart, discount, tax)
  const itemCount = cart.reduce((s, i) => s + i.qty, 0)

  const { data: currentDrawer } = useCashDrawerCurrentQuery()
  const kasOpen = currentDrawer?.status === 'open'
  const needsTable = orderType === 'dine_in' && !table

  return (
    <div className="flex h-full flex-col border-l border-gray-200 bg-white">
      <div className="flex shrink-0 items-center gap-2 border-b px-4 py-3">
        <Receipt size={16} className="text-gray-500" />
        <span className="text-sm font-semibold text-gray-800">Pesanan</span>
        {itemCount > 0 && <span className="rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700">{itemCount}</span>}
        {cart.length > 0 && (
          <Button
            variant="ghost"
            size="sm"
            className="ml-auto h-7 gap-1 text-xs text-red-500 hover:bg-red-50 hover:text-red-600"
            onClick={() => setClearConfirmOpen(true)}
          >
            <Trash2 size={12} />
            Kosongkan
          </Button>
        )}
      </div>

      <OrderInfoPanel />
      <CustomerSelector />

      <ScrollArea className="min-h-0 flex-1">
        {cart.length === 0 ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-sm text-gray-400">Belum ada item. Ketuk menu untuk menambah.</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-100">
            {cart.map((item) => (
              <CartItemRow key={item.key} item={item} />
            ))}
          </ul>
        )}
        {/* Diskon dan pajak ikut tergulir bersama item agar daftar item tetap lega di layar laptop. */}
        {cart.length > 0 && (
          <div className="border-t">
            <DiscountInput />
            <TaxInput />
          </div>
        )}
      </ScrollArea>

      <div className="shrink-0 space-y-1.5 border-t bg-gray-50 px-4 py-3 text-sm">
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
        <div className="flex items-baseline justify-between border-t pt-2 font-bold text-gray-900">
          <span className="text-sm">TOTAL</span>
          <span className="text-xl" aria-label="Total pesanan">{formatRupiah(summary.grandTotal)}</span>
        </div>
      </div>

      {!kasOpen && (
        <div className="mx-4 mb-2 flex shrink-0 items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700">
          <AlertTriangle size={13} className="mt-0.5 shrink-0" />
          <span>Kas belum dibuka. Buka kas terlebih dahulu untuk memproses transaksi.</span>
        </div>
      )}
      {needsTable && cart.length > 0 && (
        <p className="mx-4 mb-2 shrink-0 text-xs text-amber-700">Pilih meja untuk pesanan makan di tempat.</p>
      )}

      <div className="shrink-0 border-t px-4 py-3">
        <Button className="h-11 w-full gap-2 text-base" onClick={openPaymentModal} disabled={cart.length === 0 || !kasOpen || needsTable}>
          Bayar {cart.length > 0 && formatRupiah(summary.grandTotal)}
        </Button>
      </div>

      <ConfirmDialog
        open={clearConfirmOpen}
        onOpenChange={setClearConfirmOpen}
        title="Kosongkan Pesanan"
        description="Semua item, meja, dan nama pemesan akan dihapus. Yakin?"
        confirmLabel="Ya, Kosongkan"
        variant="destructive"
        onConfirm={() => {
          clearCart()
          setClearConfirmOpen(false)
        }}
      />
    </div>
  )
}
