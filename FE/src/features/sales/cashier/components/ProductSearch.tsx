import { useEffect, useRef, useState } from 'react'
import { Loader2, Search } from 'lucide-react'
import { toast } from 'sonner'
import { useQueryClient } from '@tanstack/react-query'

import { Input } from '@/shared/components/ui/input'
import { formatRupiah } from '@/shared/utils'
import { queryKeys } from '@/shared/constants'
import { api } from '@/services'
import type { Product, ProductPackage, PriceTier } from '@/features/products/products'

import { useCashierStore } from '../cashier.store'
import { getApplicablePrice } from '../cashier.utils'
import { useBarcodeScan } from '../hooks/useBarcodeScan'
import { useProductSearch } from '../hooks/useProductSearch'
import type { ProductSearchResult } from '../cashier.types'

interface ResolvedCard {
  product: Product
}

export function ProductSearch() {
  const inputRef = useRef<HTMLInputElement>(null)
  const { keyword, setKeyword, results, isLoading, clearSearch } = useProductSearch()
  const { handleBarcodeEnter, isScanning } = useBarcodeScan()
  const { addToCart } = useCashierStore()
  const qc = useQueryClient()

  // product id → resolved data (unit + harga sudah diload)
  const [resolvedCards, setResolvedCards] = useState<Map<number, ResolvedCard>>(new Map())
  // Ref untuk mencegah setState setelah komponen unmount
  const mountedRef = useRef(true)
  useEffect(() => {
    mountedRef.current = true
    return () => { mountedRef.current = false }
  }, [])

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const fetchFullProduct = async (id: number): Promise<Product> => {
    const [product, packages, prices] = await Promise.all([
      qc.fetchQuery({
        queryKey: queryKeys.products.detail(id),
        queryFn: () => api.post<Product>(`/products/detail/${id}`, {}),
        staleTime: 60_000,
      }) as Promise<Product>,
      qc.fetchQuery({
        queryKey: queryKeys.products.productUnits(id),
        queryFn: () => api.post<ProductPackage[]>(`/products/${id}/packages/list`, {}),
        staleTime: 60_000,
      }) as Promise<ProductPackage[]>,
      qc.fetchQuery({
        queryKey: queryKeys.products.priceTiers(id),
        queryFn: () => api.post<PriceTier[]>(`/products/${id}/prices/list`, {}),
        staleTime: 60_000,
      }) as Promise<PriceTier[]>,
    ])
    return {
      ...product,
      units: Array.isArray(packages) ? packages : [],
      prices: Array.isArray(prices) ? prices : [],
    }
  }

  // Fetch unit & harga semua produk di hasil search, update state hanya di async callback
  useEffect(() => {
    results.forEach(async (item) => {
      // Skip jika sudah resolved — stale entries di map tidak dirender, cukup aman
      setResolvedCards((prev) => {
        if (prev.has(item.id)) return prev
        // Mulai fetch tanpa mengubah state secara sinkron
        fetchFullProduct(item.id)
          .then((product) => {
            if (!mountedRef.current) return
            setResolvedCards((p) => {
              const next = new Map(p)
              next.set(item.id, { product })
              return next
            })
          })
          .catch(() => {
            // Biarkan kartu render tanpa unit buttons
          })
        return prev // tidak ubah state sinkron
      })
    })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [results])

  const addItemToCart = (product: Product, unitId: number, unitName: string) => {
    const pkg = product.units.find((u) => u.unit_id === unitId)
    // Pastikan price selalu number — selling_price dari API bisa datang sebagai string
    const price = Number(getApplicablePrice(product.prices, unitId, 1) ?? pkg?.selling_price ?? 0)
    addToCart({
      product_id: product.id,
      product_name: product.name,
      unit_id: unitId,
      unit_name: unitName,
      conversion_qty: Number(pkg?.conversion_qty ?? 1),
      qty: 1,
      price,
      subtotal: price,
    })
    clearSearch()
  }

  const handleBarcodeSubmit = async (value: string) => {
    try {
      const { product } = await handleBarcodeEnter(value)
      const units = product.units ?? []
      if (units.length === 0) {
        toast.error('Produk ini belum memiliki unit')
        return
      }
      if (units.length === 1) {
        addItemToCart(product, units[0].unit_id, units[0].unit_name)
      } else {
        // Multi-unit via barcode: masukkan ke resolved dan tampilkan di grid
        setResolvedCards((prev) => {
          const next = new Map(prev)
          next.set(product.id, { product })
          return next
        })
        setKeyword(product.name)
      }
    } catch {
      toast.error('Produk dengan barcode tersebut tidak ditemukan')
    }
  }

  const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== 'Enter' || !keyword.trim()) return
    e.preventDefault()
    await handleBarcodeSubmit(keyword.trim())
  }

  return (
    <div className="space-y-4">
      {/* Search Input */}
      <div className="relative">
        <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        {(isLoading || isScanning) && (
          <Loader2
            size={16}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 animate-spin"
          />
        )}
        <Input
          ref={inputRef}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Cari produk atau scan barcode..."
          className="pl-9 pr-9 h-11 text-base"
        />
      </div>

      {/* Results Grid */}
      {keyword.length >= 2 && (
        <div className="grid gap-3" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(160px, 1fr))' }}>
          {results.length === 0 && !isLoading ? (
            <p className="col-span-full text-center text-sm text-gray-400 py-8">
              Produk tidak ditemukan
            </p>
          ) : (
            results.map((item: ProductSearchResult) => {
              const resolved = resolvedCards.get(item.id)
              const units = resolved?.product.units ?? []
              // Loading: belum ada di resolvedCards
              const isCardLoading = !resolved
              const isOutOfStock = !isCardLoading && resolved.product.stock <= 0
              const isLowStock =
                !isCardLoading && !isOutOfStock && resolved.product.stock < resolved.product.min_stock

              return (
                <div
                  key={item.id}
                  className={`flex flex-col rounded-lg border bg-white shadow-sm hover:shadow-md transition-all ${
                    isOutOfStock
                      ? 'border-red-200 bg-red-50/40 hover:border-red-300'
                      : isLowStock
                        ? 'border-amber-200 bg-amber-50/30 hover:border-amber-300'
                        : 'hover:border-blue-300'
                  }`}
                >
                  {/* Info produk */}
                  <div className="flex flex-col items-center gap-1 px-3 pt-3 pb-2 text-center">
                    <div className="flex h-9 w-9 items-center justify-center rounded-full bg-gray-100 text-lg">
                      {isCardLoading ? (
                        <Loader2 size={16} className="animate-spin text-gray-400" />
                      ) : (
                        '📦'
                      )}
                    </div>
                    <p className="text-sm font-medium text-gray-800 line-clamp-2 leading-tight">
                      {item.name}
                    </p>
                  </div>

                  {/* Unit buttons — langsung tampil di kartu */}
                  <div className="border-t px-2 py-2">
                    {isLowStock && (
                      <p className="text-center text-[10px] font-medium text-amber-600 mb-1">
                        Sisa {resolved!.product.stock} — stok menipis
                      </p>
                    )}
                    {isCardLoading ? (
                      <div className="flex justify-center py-1">
                        <Loader2 size={14} className="animate-spin text-gray-300" />
                      </div>
                    ) : isOutOfStock ? (
                      <p className="text-center text-xs text-red-500 font-medium py-1.5">Stok Habis</p>
                    ) : units.length === 0 ? (
                      <p className="text-center text-xs text-gray-400 py-1">Belum ada unit</p>
                    ) : units.length === 1 ? (
                      <button
                        onClick={() =>
                          addItemToCart(resolved!.product, units[0].unit_id, units[0].unit_name)
                        }
                        className="w-full rounded-md bg-blue-50 border border-blue-200 py-1.5 text-xs font-semibold text-blue-700 hover:bg-blue-100 active:scale-95 transition-all"
                      >
                        {units[0].unit_name}
                        {' — '}
                        {formatRupiah(
                          getApplicablePrice(resolved!.product.prices, units[0].unit_id, 1) ??
                            units[0].selling_price
                        )}
                      </button>
                    ) : (
                      <div className="flex flex-wrap gap-1.5 justify-center">
                        {units.map((unit) => {
                          // Prioritas: price tier → selling_price di package
                          const price =
                            getApplicablePrice(resolved!.product.prices, unit.unit_id, 1) ??
                            unit.selling_price
                          return (
                            <button
                              key={unit.unit_id}
                              onClick={() =>
                                addItemToCart(resolved!.product, unit.unit_id, unit.unit_name)
                              }
                              className="flex flex-col items-center rounded-md border border-blue-200 bg-blue-50 px-3 py-1.5 hover:border-blue-500 hover:bg-blue-100 active:scale-95 transition-all"
                            >
                              <span className="text-xs font-semibold text-blue-800">
                                {unit.unit_name}
                              </span>
                              <span className="text-xs text-blue-600">
                                {formatRupiah(price)}
                              </span>
                            </button>
                          )
                        })}
                      </div>
                    )}
                  </div>
                </div>
              )
            })
          )}
        </div>
      )}
    </div>
  )
}
