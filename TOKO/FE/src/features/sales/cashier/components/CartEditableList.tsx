import { useState } from 'react'
import { ShoppingCart, Trash2 } from 'lucide-react'

import { ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'

import { useCashierStore } from '../cashier.store'
import { CartItemRow } from './CartItemRow'

export function CartEditableList() {
  const [clearConfirmOpen, setClearConfirmOpen] = useState(false)
  const { cart, clearCart } = useCashierStore()
  const itemCount = cart.reduce((s, i) => s + i.qty, 0)

  return (
    <div className="flex flex-col min-h-0 flex-1">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-t bg-white shrink-0">
        <div className="flex items-center gap-2">
          <ShoppingCart size={16} className="text-gray-500" />
          <span className="text-sm font-semibold text-gray-800">
            Keranjang
          </span>
          {itemCount > 0 && (
            <span className="rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700">
              {itemCount}
            </span>
          )}
        </div>
        {cart.length > 0 && (
          <Button
            variant="ghost"
            size="sm"
            className="h-7 gap-1 text-xs text-red-500 hover:text-red-600 hover:bg-red-50"
            onClick={() => setClearConfirmOpen(true)}
          >
            <Trash2 size={12} />
            Kosongkan
          </Button>
        )}
      </div>

      {/* List item editable — scrollable */}
      <ScrollArea className="flex-1">
        {cart.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 py-10 text-gray-400">
            <ShoppingCart size={28} className="opacity-30" />
            <p className="text-sm">Belum ada item</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-100">
            {cart.map((item) => (
              <CartItemRow key={`${item.product_id}-${item.unit_id}`} item={item} />
            ))}
          </ul>
        )}
      </ScrollArea>

      {/* Subtotal ringkas di bawah list */}
      {cart.length > 0 && (
        <div className="border-t px-4 py-2 shrink-0 bg-gray-50">
          <div className="flex justify-between text-sm text-gray-600">
            <span>
              {cart.length} produk
              {cart.length !== itemCount && ` · ${itemCount} item`}
            </span>
            <span className="font-semibold text-gray-800">
              {formatRupiah(cart.reduce((s, i) => s + i.subtotal, 0))}
            </span>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={clearConfirmOpen}
        onOpenChange={setClearConfirmOpen}
        title="Kosongkan Keranjang"
        description="Semua item di keranjang akan dihapus. Yakin?"
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
