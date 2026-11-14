import { useState } from 'react'
import { Minus, Plus, Tag, Trash2 } from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/shared/components/ui/button'
import { Collapsible, CollapsibleContent } from '@/shared/components/ui/collapsible'
import { Input } from '@/shared/components/ui/input'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { formatRupiah } from '@/shared/utils'

import { useMenuCatalogQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'
import type { CartItem } from '../cashier.types'
import { cartItemDetails } from '../cashier.utils'

interface CartItemRowProps {
  item: CartItem
}

function DiscountBadge({ item }: { item: CartItem }) {
  if (!item.discount_type || !item.discount_value) return null
  const label = item.discount_type === 'percent' ? `Disc ${item.discount_value}%` : `-${formatRupiah(item.discount_value)}`
  return <span className="inline-flex items-center rounded-full bg-red-100 px-1.5 py-0.5 text-xs font-medium text-red-600">{label}</span>
}

function ItemDiscountInput({ item, onClose }: { item: CartItem; onClose: () => void }) {
  const { setItemDiscount } = useCashierStore()
  const [discType, setDiscType] = useState<'percent' | 'nominal'>(item.discount_type ?? 'percent')
  const [discValue, setDiscValue] = useState(item.discount_value ?? 0)

  const apply = () => {
    if (discValue >= 0) setItemDiscount(item.key, discType, discValue)
    onClose()
  }

  return (
    <div className="mt-1 flex flex-wrap items-center gap-1.5">
      <div className="flex overflow-hidden rounded-md border border-gray-200">
        {(['percent', 'nominal'] as const).map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setDiscType(t)}
            className={`px-2 py-1 text-xs font-medium ${discType === t ? 'bg-[#2c3e50] text-white' : 'bg-white text-gray-600 hover:bg-gray-50'}`}
          >
            {t === 'percent' ? '%' : 'Rp'}
          </button>
        ))}
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
            if (e.key === 'Enter') apply()
            if (e.key === 'Escape') onClose()
          }}
          className="h-7 w-20 px-2 text-xs"
          aria-label="Diskon item persen"
        />
      ) : (
        <RupiahInput value={discValue} onChange={setDiscValue} className="h-7 text-xs" aria-label="Diskon item rupiah" />
      )}
      <button type="button" onClick={apply} className="h-7 rounded bg-[#2c3e50] px-2 text-xs text-white hover:bg-[#1a252f]">
        OK
      </button>
      {!!item.discount_value && (
        <button
          type="button"
          onClick={() => {
            setItemDiscount(item.key, 'percent', 0)
            onClose()
          }}
          className="h-7 rounded border border-gray-200 px-2 text-xs text-gray-500 hover:text-red-500"
        >
          Hapus
        </button>
      )}
    </div>
  )
}

/** Satu baris pesanan di panel kanan: nama + kombinasi, tambahan, catatan, jumlah, diskon item. */
export function CartItemRow({ item }: CartItemRowProps) {
  const { cart, updateQty, removeFromCart } = useCashierStore()
  const { data: catalog } = useMenuCatalogQuery()
  const [showDiscount, setShowDiscount] = useState(false)
  const hasDiscount = !!(item.discount_type && item.discount_value && item.discount_value > 0)
  const details = cartItemDetails(item)

  // Barang jadi: jumlah total menu ini di keranjang tidak boleh melebihi sisa stok.
  const increase = () => {
    const stock = catalog?.items.find((m) => m.id === item.menu_item_id)?.stock
    if (stock != null) {
      const inCart = cart.filter((c) => c.menu_item_id === item.menu_item_id).reduce((s, c) => s + c.qty, 0)
      if (inCart + 1 > stock) {
        toast.error(`Stok ${item.product_name} tidak cukup (sisa ${stock})`)
        return
      }
    }
    updateQty(item.key, item.qty + 1)
  }

  return (
    <li className="space-y-1.5 px-4 py-3" aria-label={`Item ${item.product_name}`}>
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1">
          <p className="text-sm font-semibold text-gray-800">
            {item.product_name}
            {item.variant_name && <span className="font-normal text-gray-500"> · {item.variant_name}</span>}
          </p>
          {details.map((d) => (
            <p key={d} className="text-xs text-gray-500">
              {d}
            </p>
          ))}
        </div>
        <button
          type="button"
          onClick={() => removeFromCart(item.key)}
          aria-label="Hapus item"
          className="mt-0.5 shrink-0 text-gray-300 transition-colors hover:text-red-500 touch:flex touch:min-w-10 touch:items-center touch:justify-center"
        >
          <Trash2 size={14} />
        </button>
      </div>

      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-1.5">
          <Button variant="outline" size="icon" className="h-9 w-9 shrink-0" aria-label="Kurangi" onClick={() => updateQty(item.key, item.qty - 1)}>
            <Minus size={14} />
          </Button>
          <span className="w-8 text-center text-sm font-semibold" aria-label="Jumlah item">
            {item.qty}
          </span>
          <Button variant="outline" size="icon" className="h-9 w-9 shrink-0" aria-label="Tambah" onClick={increase}>
            <Plus size={14} />
          </Button>
          <button
            type="button"
            onClick={() => setShowDiscount((v) => !v)}
            aria-label="Atur diskon item"
            className={`ml-1 flex h-7 w-7 items-center justify-center rounded transition-colors ${
              hasDiscount ? 'bg-red-100 text-red-500' : 'text-gray-300 hover:bg-gray-100 hover:text-gray-500'
            }`}
          >
            <Tag size={12} />
          </button>
          <DiscountBadge item={item} />
        </div>
        <div className="text-right">
          {hasDiscount && <p className="text-[11px] text-gray-400 line-through">{formatRupiah(item.price * item.qty)}</p>}
          <p className="text-sm font-semibold text-gray-800">{formatRupiah(item.subtotal)}</p>
        </div>
      </div>

      <Collapsible open={showDiscount}>
        <CollapsibleContent>
          <ItemDiscountInput item={item} onClose={() => setShowDiscount(false)} />
        </CollapsibleContent>
      </Collapsible>
    </li>
  )
}
