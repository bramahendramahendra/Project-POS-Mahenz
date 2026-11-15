import { useMemo, useState } from 'react'
import { Coffee, Search } from 'lucide-react'
import { toast } from 'sonner'

import { Input } from '@/shared/components/ui/input'
import { formatRupiah } from '@/shared/utils'

import { useMenuCatalogQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'
import type { CatalogItem } from '../cashier.types'
import { cartLineKey } from '../cashier.utils'
import { buildCartItem, isMenuUnavailable } from './menu-option.utils'
import { MenuOptionDialog } from './MenuOptionDialog'

function priceLabel(item: CatalogItem): string {
  const prices = item.variants.map((v) => v.price)
  const min = Math.min(...prices)
  return prices.length > 1 && Math.max(...prices) !== min ? `mulai ${formatRupiah(min)}` : formatRupiah(min)
}

/**
 * Katalog menu kasir (tahap 4.7b): pilih kategori atau cari nama, lalu ketuk tombol menu.
 * Menu satu kombinasi tanpa tambahan langsung masuk keranjang; lainnya membuka pilihan.
 */
export function MenuCatalog() {
  const { data, isLoading, isError, refetch } = useMenuCatalogQuery()
  const { cart, addToCart } = useCashierStore()
  const [categoryId, setCategoryId] = useState<number | 'all'>('all')
  const [search, setSearch] = useState('')
  const [picking, setPicking] = useState<CatalogItem | null>(null)

  const items = useMemo(() => {
    const keyword = search.trim().toLowerCase()
    return (data?.items ?? []).filter(
      (it) => (keyword ? it.name.toLowerCase().includes(keyword) : categoryId === 'all' || it.category_id === categoryId),
    )
  }, [data, categoryId, search])

  // Jumlah barang jadi yang sudah ada di keranjang untuk menu tertentu (cek stok sebelum menambah).
  const qtyInCart = (menuId: number) => cart.filter((c) => c.menu_item_id === menuId).reduce((s, c) => s + c.qty, 0)

  const handlePick = (item: CatalogItem) => {
    if (isMenuUnavailable(item)) return
    if (item.stock !== null && qtyInCart(item.id) + 1 > item.stock) {
      toast.error(`Stok ${item.name} tidak cukup (sisa ${item.stock})`)
      return
    }
    const simple = item.variants.length === 1 && item.modifier_groups.length === 0
    if (simple) {
      const variant = item.variants[0]
      addToCart(buildCartItem(item, variant, [], 1, '', cartLineKey(variant.id, [], '')))
      toast.success(`${item.name} ditambahkan`, { duration: 1200, id: 'cart-add' })
      return
    }
    setPicking(item)
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="shrink-0 space-y-3 px-4 pb-3 pt-4">
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari menu..."
            aria-label="Cari menu"
            className="pl-9"
          />
        </div>
        {!search && (
          <div className="flex gap-2 overflow-x-auto pb-1" role="tablist" aria-label="Kategori menu">
            {[{ id: 'all' as const, name: 'Semua' }, ...(data?.categories ?? [])].map((c) => (
              <button
                key={c.id}
                type="button"
                role="tab"
                aria-selected={categoryId === c.id}
                onClick={() => setCategoryId(c.id)}
                className={`shrink-0 rounded-full border px-4 py-1.5 text-sm font-medium transition-colors ${
                  categoryId === c.id ? 'border-[#2c3e50] bg-[#2c3e50] text-white' : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50'
                }`}
              >
                {c.name}
              </button>
            ))}
          </div>
        )}
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        {isLoading ? (
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-4">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="h-24 animate-pulse rounded-xl bg-gray-200" />
            ))}
          </div>
        ) : isError ? (
          <div className="py-10 text-center text-sm text-gray-500">
            Gagal memuat menu.{' '}
            <button type="button" className="text-blue-600 underline" onClick={() => refetch()}>
              Coba lagi
            </button>
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-col items-center gap-2 py-10 text-gray-400">
            <Coffee size={28} className="opacity-40" />
            <p className="text-sm">{search ? 'Menu tidak ditemukan' : (data?.categories ?? []).length === 0 ? 'Belum ada menu yang dijual — tambahkan di menu Daftar Menu' : 'Belum ada menu di kategori ini'}</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-4">
            {items.map((item) => {
              const unavailable = isMenuUnavailable(item)
              return (
                <button
                  key={item.id}
                  type="button"
                  disabled={unavailable}
                  onClick={() => handlePick(item)}
                  aria-label={`Menu ${item.name}`}
                  className={`flex min-h-24 flex-col justify-between rounded-xl border bg-white p-3 text-left shadow-sm transition-colors ${
                    unavailable ? 'cursor-not-allowed opacity-50' : 'hover:border-[#2c3e50] active:bg-gray-50'
                  }`}
                >
                  <span className="line-clamp-2 text-sm font-semibold text-gray-800">{item.name}</span>
                  <span className="mt-2 flex items-end justify-between gap-1">
                    <span className="text-sm font-medium text-[#2c3e50]">{priceLabel(item)}</span>
                    {unavailable ? (
                      <span className="rounded-full bg-red-100 px-2 py-0.5 text-[11px] font-semibold text-red-600">Habis</span>
                    ) : item.stock !== null ? (
                      <span className="text-[11px] text-gray-500">Sisa {item.stock}</span>
                    ) : null}
                  </span>
                </button>
              )
            })}
          </div>
        )}
      </div>

      {picking && (
        <MenuOptionDialog
          key={picking.id}
          item={picking}
          maxQty={picking.stock !== null ? picking.stock - qtyInCart(picking.id) : undefined}
          onClose={() => setPicking(null)}
          onAdd={(line) => {
            addToCart(line)
            toast.success(`${picking.name} ditambahkan`, { duration: 1200, id: 'cart-add' })
            setPicking(null)
          }}
        />
      )}
    </div>
  )
}
