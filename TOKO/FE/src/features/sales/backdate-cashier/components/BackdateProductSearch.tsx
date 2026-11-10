/**
 * Product search for backdate cashier — reuses the same search/barcode logic
 * but adds to the backdate store instead of the regular cashier store.
 */
import { useEffect, useRef, useState } from 'react'
import { Loader2, Search } from 'lucide-react'
import { toast } from 'sonner'
import { useQueryClient } from '@tanstack/react-query'

import { Input } from '@/shared/components/ui/input'
import { formatRupiah } from '@/shared/utils'
import { queryKeys } from '@/shared/constants'
import {
  fetchProductDetail,
  fetchProductPackages,
  fetchProductPrices,
} from '@/features/products/products'
import type { Product, ProductPackage, PriceTier } from '@/features/products/products'

import { useBackdateCashierStore } from '../backdate-cashier.store'
import { getApplicablePrice } from '../../cashier/cashier.utils'
import { useBarcodeScan } from '../../cashier/hooks/useBarcodeScan'
import { useProductSearch } from '../../cashier/hooks/useProductSearch'
import type { ProductSearchResult } from '../../cashier/cashier.types'

interface ResolvedCard {
  product: Product
}

export function BackdateProductSearch() {
  const inputRef = useRef<HTMLInputElement>(null)
  const { keyword, setKeyword, results, isLoading, clearSearch } = useProductSearch()
  const { handleBarcodeEnter, isScanning } = useBarcodeScan()
  const { addToCart } = useBackdateCashierStore()
  const qc = useQueryClient()

  const [resolvedCards, setResolvedCards] = useState<Map<number, ResolvedCard>>(new Map())
  const mountedRef = useRef(true)
  useEffect(() => { mountedRef.current = true; return () => { mountedRef.current = false } }, [])
  useEffect(() => { inputRef.current?.focus() }, [])

  const fetchFullProduct = async (id: number): Promise<Product> => {
    const [product, packages, prices] = await Promise.all([
      qc.fetchQuery({ queryKey: queryKeys.products.detail(id), queryFn: () => fetchProductDetail(id), staleTime: 60_000 }) as Promise<Product>,
      qc.fetchQuery({ queryKey: queryKeys.products.productUnits(id), queryFn: () => fetchProductPackages(id), staleTime: 60_000 }) as Promise<ProductPackage[]>,
      qc.fetchQuery({ queryKey: queryKeys.products.priceTiers(id), queryFn: () => fetchProductPrices(id), staleTime: 60_000 }) as Promise<PriceTier[]>,
    ])
    return { ...product, units: Array.isArray(packages) ? packages : [], prices: Array.isArray(prices) ? prices : [] }
  }

  useEffect(() => {
    results.forEach((item) => {
      setResolvedCards((prev) => {
        if (prev.has(item.id)) return prev
        fetchFullProduct(item.id).then((product) => {
          if (!mountedRef.current) return
          setResolvedCards((p) => { const next = new Map(p); next.set(item.id, { product }); return next })
        }).catch(() => {})
        return prev
      })
    })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [results])

  const addItemToCart = (product: Product, pkg: ProductPackage) => {
    const label = pkg.package_name ? `${pkg.unit_name} (${pkg.package_name})` : pkg.unit_name
    const price = Number(getApplicablePrice(product.prices, pkg.unit_id, 1) ?? pkg.selling_price ?? 0)
    addToCart({
      product_id: product.id,
      product_name: product.name,
      unit_id: pkg.id,
      unit_name: label,
      conversion_qty: Number(pkg.resolved_factor ?? 1),
      is_continuous: Boolean(pkg.is_continuous),
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
      if (units.length === 0) { toast.error('Produk ini belum memiliki unit'); return }
      if (units.length === 1) { addItemToCart(product, units[0]) }
      else { setResolvedCards((prev) => { const next = new Map(prev); next.set(product.id, { product }); return next }); setKeyword(product.name) }
    } catch { toast.error('Produk dengan barcode tersebut tidak ditemukan') }
  }

  const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== 'Enter' || !keyword.trim()) return
    e.preventDefault()
    await handleBarcodeSubmit(keyword.trim())
  }

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        {(isLoading || isScanning) && <Loader2 size={16} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 animate-spin" />}
        <Input
          ref={inputRef}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Cari produk atau scan barcode..."
          className="pl-9 pr-9 h-11 text-base"
        />
      </div>

      {keyword.length >= 2 && (
        <div className="grid gap-3" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(160px, 1fr))' }}>
          {results.length === 0 && !isLoading ? (
            <p className="col-span-full text-center text-sm text-gray-400 py-8">Produk tidak ditemukan</p>
          ) : results.map((item: ProductSearchResult) => {
            const resolved = resolvedCards.get(item.id)
            const units = resolved?.product.units ?? []
            const isCardLoading = !resolved

            return (
              <div key={item.id} className="flex flex-col rounded-lg border bg-white shadow-sm hover:shadow-md transition-all hover:border-blue-300">
                <div className="flex flex-col items-center gap-1 px-3 pt-3 pb-2 text-center">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-gray-100 text-lg">
                    {isCardLoading ? <Loader2 size={16} className="animate-spin text-gray-400" /> : '📦'}
                  </div>
                  <p className="text-sm font-medium text-gray-800 line-clamp-2 leading-tight">{item.name}</p>
                </div>
                <div className="border-t px-2 py-2">
                  {isCardLoading ? (
                    <div className="flex justify-center py-1"><Loader2 size={14} className="animate-spin text-gray-300" /></div>
                  ) : units.length === 0 ? (
                    <p className="text-center text-xs text-gray-400 py-1">Belum ada unit</p>
                  ) : units.length === 1 ? (
                    <button
                      onClick={() => addItemToCart(resolved!.product, units[0])}
                      className="w-full rounded-md bg-blue-50 border border-blue-200 py-1.5 text-xs font-semibold text-blue-700 hover:bg-blue-100 active:scale-95 transition-all"
                    >
                      {units[0].package_name ? `${units[0].unit_name} (${units[0].package_name})` : units[0].unit_name}
                      {' — '}
                      {formatRupiah(getApplicablePrice(resolved!.product.prices, units[0].unit_id, 1) ?? units[0].selling_price)}
                    </button>
                  ) : (
                    <div className="flex flex-wrap gap-1.5 justify-center">
                      {units.map((unit) => (
                        <button
                          key={unit.id}
                          onClick={() => addItemToCart(resolved!.product, unit)}
                          className="rounded-md bg-blue-50 border border-blue-200 px-2 py-1 text-[11px] font-medium text-blue-700 hover:bg-blue-100 active:scale-95 transition-all"
                        >
                          {unit.package_name ? `${unit.unit_name} (${unit.package_name})` : unit.unit_name}
                          {' — '}
                          {formatRupiah(getApplicablePrice(resolved!.product.prices, unit.unit_id, 1) ?? unit.selling_price)}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
