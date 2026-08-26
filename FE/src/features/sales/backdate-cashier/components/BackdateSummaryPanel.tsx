/**
 * Summary panel for backdate cashier — shows totals and checkout button.
 * Uses backdate store.
 */
import { Button } from '@/shared/components/ui/button'
import { formatRupiah } from '@/shared/utils'

import { useBackdateCashierStore } from '../backdate-cashier.store'
import { calcCartSummary } from '../../cashier/cashier.utils'
import { CustomerSelector } from '../../cashier/components/CustomerSelector'
import { DiscountInput } from '../../cashier/components/DiscountInput'
import { TaxInput } from '../../cashier/components/TaxInput'

export function BackdateSummaryPanel() {
  const { cart, discount, tax, openPaymentModal } = useBackdateCashierStore()
  const summary = calcCartSummary(cart, discount, tax)
  const isEmpty = cart.length === 0

  return (
    <div className="flex flex-col h-full border-l bg-white">
      {/* Header */}
      <div className="px-4 py-3 border-b shrink-0">
        <h2 className="text-sm font-semibold text-gray-800">Ringkasan Historis</h2>
      </div>

      {/* Body - scrollable */}
      <div className="flex-1 overflow-y-auto px-4 py-3 space-y-4">
        {/* Customer selector */}
        <CustomerSelector
          value={useBackdateCashierStore.getState().selectedCustomer}
          onChange={(c) => useBackdateCashierStore.getState().setCustomer(c)}
        />

        {/* Discount */}
        <DiscountInput
          discount={discount}
          onChange={(d) => useBackdateCashierStore.getState().setDiscount(d)}
        />

        {/* Tax */}
        <TaxInput
          tax={tax}
          onChange={(percent) => useBackdateCashierStore.getState().setTax(percent)}
        />

        {/* Totals */}
        <div className="space-y-2 pt-2 border-t">
          <div className="flex justify-between text-sm">
            <span className="text-gray-500">Subtotal</span>
            <span className="font-medium">{formatRupiah(summary.subtotal)}</span>
          </div>
          {summary.discountAmount > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Diskon</span>
              <span className="font-medium text-red-600">-{formatRupiah(summary.discountAmount)}</span>
            </div>
          )}
          {summary.taxAmount > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Pajak</span>
              <span className="font-medium">+{formatRupiah(summary.taxAmount)}</span>
            </div>
          )}
          <div className="flex justify-between text-base font-bold pt-2 border-t">
            <span>Total</span>
            <span>{formatRupiah(summary.grandTotal)}</span>
          </div>
        </div>
      </div>

      {/* Footer - checkout button */}
      <div className="px-4 py-3 border-t shrink-0">
        <Button
          className="w-full h-12 text-base font-semibold"
          disabled={isEmpty}
          onClick={openPaymentModal}
        >
          Bayar {!isEmpty && formatRupiah(summary.grandTotal)}
        </Button>
      </div>
    </div>
  )
}
