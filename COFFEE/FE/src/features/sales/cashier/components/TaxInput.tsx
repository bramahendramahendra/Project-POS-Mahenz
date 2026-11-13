import { formatRupiah } from '@/shared/utils'
import { Input } from '@/shared/components/ui/input'

import { useCashierStore } from '../cashier.store'
import { calcDiscountAmount, calcSubtotal, calcTaxAmount } from '../cashier.utils'
import type { CartItem, Discount, Tax } from '../cashier.types'

// Slice minimal yang dibutuhkan. Kedua store (kasir reguler & historis)
// memenuhinya, jadi komponen bisa dipakai di kedua konteks.
type TaxInputStore = () => {
  cart: CartItem[]
  discount: Discount
  tax: Tax
  setTax: (percent: number) => void
}

interface TaxInputProps {
  // Default useCashierStore agar pemakaian lama tanpa prop tetap jalan.
  store?: TaxInputStore
}

export function TaxInput({ store = useCashierStore }: TaxInputProps) {
  const { cart, discount, tax, setTax } = store()

  const subtotal = calcSubtotal(cart)
  const discountAmount = calcDiscountAmount(subtotal, discount)
  const previewTax = calcTaxAmount(subtotal, discountAmount, tax.percent)

  const handleChange = (raw: string) => {
    const v = parseFloat(raw)
    if (isNaN(v) || v < 0) return
    setTax(Math.min(v, 100))
  }

  return (
    <div className="flex items-center gap-2 px-4 py-2 border-t text-sm">
      <span className="text-gray-600 shrink-0 w-12">Pajak</span>

      <Input
        type="number"
        min={0}
        max={100}
        value={tax.percent}
        onChange={(e) => handleChange(e.target.value)}
        className="h-7 w-20 text-sm text-right px-2"
      />
      <span className="text-gray-400 text-xs">%</span>

      {previewTax > 0 && (
        <span className="ml-auto text-gray-600 text-xs shrink-0">+{formatRupiah(previewTax)}</span>
      )}
    </div>
  )
}
