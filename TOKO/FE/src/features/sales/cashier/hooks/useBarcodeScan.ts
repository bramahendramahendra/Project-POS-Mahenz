import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import { fetchProductPackages, fetchProductPrices } from '@/features/products/products'
import type { Product, ProductPackage } from '@/features/products/products'

export const useBarcodeScan = () => {
  const [isScanning, setIsScanning] = useState(false)
  const qc = useQueryClient()

  const handleBarcodeEnter = async (
    code: string
  ): Promise<{ product: Product; units: ProductPackage[] }> => {
    setIsScanning(true)
    try {
      const product = await qc.fetchQuery({
        queryKey: queryKeys.products.barcode(code),
        queryFn: () => api.post<Product>(`/products/by-barcode/${code}`, {}),
        staleTime: 30_000,
      })
      // Endpoint by-barcode hanya mengembalikan unit default (unit_id/unit_name flat),
      // bukan daftar units/prices — fetch terpisah seperti alur search produk.
      const [units, prices] = await Promise.all([
        qc.fetchQuery({
          queryKey: queryKeys.products.productUnits(product.id),
          queryFn: () => fetchProductPackages(product.id),
          staleTime: 60_000,
        }),
        qc.fetchQuery({
          queryKey: queryKeys.products.priceTiers(product.id),
          queryFn: () => fetchProductPrices(product.id),
          staleTime: 60_000,
        }),
      ])
      const resolvedUnits = Array.isArray(units) ? units : []
      const resolvedPrices = Array.isArray(prices) ? prices : []
      return {
        product: { ...product, units: resolvedUnits, prices: resolvedPrices } as Product,
        units: resolvedUnits,
      }
    } finally {
      setIsScanning(false)
    }
  }

  return { handleBarcodeEnter, isScanning }
}
