import { useState } from 'react'
import { Minus, Plus, Trash2, Tag } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { Input } from '@/shared/components/ui/input'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { formatRupiah } from '@/shared/utils'

import { useCashierStore } from '../cashier.store'
import type { CartItem } from '../cashier.types'

interface CartItemRowProps {
  item: CartItem
}

function DiscountBadge({ item }: { item: CartItem }) {
  if (!item.discount_type || !item.discount_value) return null
  const label =
    item.discount_type === 'percent'
      ? `Disc ${item.discount_value}%`
      : `-${formatRupiah(item.discount_value)}`
  return (
    <span className="inline-flex items-center rounded-full bg-red-100 px-1.5 py-0.5 text-xs font-medium text-red-600">
      {label}
    </span>
  )
}

interface DiscountInputProps {
  item: CartItem
  onClose: () => void
}

function DiscountInput({ item, onClose }: DiscountInputProps) {
  const { setItemDiscount } = useCashierStore()
  const [discType, setDiscType] = useState<'percent' | 'nominal'>(item.discount_type ?? 'percent')
  const [discValue, setDiscValue] = useState(item.discount_value ?? 0)

  function handleApply() {
    if (discValue >= 0) {
      setItemDiscount(item.product_id, item.unit_id, discType, discValue)
    }
    onClose()
  }

  function handleRemove() {
    setItemDiscount(item.product_id, item.unit_id, 'percent', 0)
    onClose()
  }

  return (
    <div className="flex items-center gap-1.5 mt-1 flex-wrap">
      {/* Type toggle */}
      <div className="flex rounded-md border border-gray-200 overflow-hidden">
        <button
          type="button"
          onClick={() => setDiscType('percent')}
          className={`px-2 py-1 text-xs font-medium transition-colors ${
            discType === 'percent' ? 'bg-[#2c3e50] text-white' : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          %
        </button>
        <button
          type="button"
          onClick={() => setDiscType('nominal')}
          className={`px-2 py-1 text-xs font-medium transition-colors ${
            discType === 'nominal' ? 'bg-[#2c3e50] text-white' : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          Rp
        </button>
      </div>

      {discType === 'percent' ? (
        <Input
          autoFocus
          type="number"
          min={0}
          max={100}
          value={discValue}
          onChange={(e) => setDiscValue(parseFloat(e.target.value) || 0)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') handleApply()
            if (e.key === 'Escape') onClose()
          }}
          className="h-6 w-20 text-xs px-2"
          placeholder="0–100"
        />
      ) : (
        <RupiahInput
          value={discValue}
          onChange={(v) => setDiscValue(v)}
          className="h-6 text-xs"
        />
      )}

      <button
        type="button"
        onClick={handleApply}
        className="h-6 px-2 text-xs rounded bg-[#2c3e50] text-white hover:bg-[#1a252f] transition-colors"
      >
        OK
      </button>

      {item.discount_value && item.discount_value > 0 && (
        <button
          type="button"
          onClick={handleRemove}
          className="h-6 px-2 text-xs rounded border border-gray-200 text-gray-500 hover:text-red-500 transition-colors"
        >
          Hapus
        </button>
      )}
    </div>
  )
}

export function CartItemRow({ item }: CartItemRowProps) {
  const { updateQty, updatePrice, removeFromCart } = useCashierStore()
  const [editingPrice, setEditingPrice] = useState(false)
  const [priceInput, setPriceInput] = useState(item.price)
  const [showDiscount, setShowDiscount] = useState(false)

  const hasDiscount = !!(item.discount_type && item.discount_value && item.discount_value > 0)
  const displayPrice = item.effective_price ?? item.price

  const handleQtyChange = (raw: string) => {
    const v = parseInt(raw, 10)
    if (!isNaN(v) && v > 0) updateQty(item.product_id, item.unit_id, v)
  }

  const handlePriceCommit = (v: number) => {
    if (v >= 0) updatePrice(item.product_id, item.unit_id, v)
    else setPriceInput(item.price)
    setEditingPrice(false)
  }

  return (
    <li className="px-4 py-3 space-y-1.5">
      {/* Row 1: Name + delete */}
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-gray-800 truncate">{item.product_name}</p>
          <p className="text-xs text-gray-400">{item.unit_name}</p>
        </div>
        <button
          onClick={() => removeFromCart(item.product_id, item.unit_id)}
          className="text-gray-300 hover:text-red-500 transition-colors mt-0.5 shrink-0"
        >
          <Trash2 size={14} />
        </button>
      </div>

      {/* Row 2: Price (editable) + discount badge + subtotal */}
      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-1.5 flex-wrap">
          {editingPrice ? (
            <RupiahInput
              autoFocus
              value={priceInput}
              onChange={(v) => setPriceInput(v)}
              onBlur={() => handlePriceCommit(priceInput)}
              className="h-6 text-xs"
            />
          ) : (
            <Tooltip>
              <TooltipTrigger asChild>
                <button
                  onClick={() => {
                    setPriceInput(item.price)
                    setEditingPrice(true)
                  }}
                  className={`text-xs hover:underline ${hasDiscount ? 'text-gray-400 line-through' : 'text-blue-600'}`}
                >
                  {formatRupiah(item.price)}
                </button>
              </TooltipTrigger>
              <TooltipContent>Klik untuk ubah harga</TooltipContent>
            </Tooltip>
          )}

          {hasDiscount && (
            <span className="text-xs font-medium text-blue-600">
              {formatRupiah(displayPrice)}
            </span>
          )}

          <DiscountBadge item={item} />

          {/* Discount toggle button */}
          <Tooltip>
            <TooltipTrigger asChild>
              <button
                type="button"
                onClick={() => setShowDiscount((v) => !v)}
                className={`flex items-center justify-center h-5 w-5 rounded transition-colors ${
                  hasDiscount
                    ? 'bg-red-100 text-red-500'
                    : 'text-gray-300 hover:text-gray-500 hover:bg-gray-100'
                }`}
              >
                <Tag size={11} />
              </button>
            </TooltipTrigger>
            <TooltipContent>Atur diskon item</TooltipContent>
          </Tooltip>
        </div>

        <span className="text-sm font-semibold text-gray-800 shrink-0">
          = {formatRupiah(item.subtotal)}
        </span>
      </div>

      {/* Inline discount input */}
      {showDiscount && (
        <DiscountInput item={item} onClose={() => setShowDiscount(false)} />
      )}

      {/* Row 3: Qty controls */}
      <div className="flex items-center gap-1.5">
        <Button
          variant="outline"
          size="icon"
          className="h-6 w-6 shrink-0"
          onClick={() => updateQty(item.product_id, item.unit_id, item.qty - 1)}
        >
          <Minus size={10} />
        </Button>
        <Input
          type="number"
          value={item.qty}
          onChange={(e) => handleQtyChange(e.target.value)}
          className="h-6 w-14 text-center text-sm px-1"
          min={1}
        />
        <Button
          variant="outline"
          size="icon"
          className="h-6 w-6 shrink-0"
          onClick={() => updateQty(item.product_id, item.unit_id, item.qty + 1)}
        >
          <Plus size={10} />
        </Button>
      </div>

    </li>
  )
}
