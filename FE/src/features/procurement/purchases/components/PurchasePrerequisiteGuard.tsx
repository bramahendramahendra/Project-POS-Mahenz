import type { ReactNode } from 'react'
import { Package, Truck } from 'lucide-react'

import { useSupplierListQuery } from '@/features/procurement/suppliers'
import { useProductListQuery } from '@/features/products/products'

interface PurchasePrerequisiteGuardProps {
  children: ReactNode
}

export function PurchasePrerequisiteGuard({ children }: PurchasePrerequisiteGuardProps) {
  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 1, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })

  if (isSuppliersLoading || isProductsLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  const hasSuppliers = (suppliersData?.total ?? 0) > 0
  const hasProducts = (productsData?.total ?? 0) > 0

  if (!hasSuppliers || !hasProducts) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4 flex gap-3">
          <div className={`rounded-full p-3 ${!hasSuppliers ? 'bg-amber-50' : 'bg-green-50'}`}>
            <Truck size={24} className={!hasSuppliers ? 'text-amber-500' : 'text-green-500'} />
          </div>
          <div className={`rounded-full p-3 ${!hasProducts ? 'bg-amber-50' : 'bg-green-50'}`}>
            <Package size={24} className={!hasProducts ? 'text-amber-500' : 'text-green-500'} />
          </div>
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa menambah pembelian</h3>
        <p className="mb-1 text-sm text-gray-500">
          Sebelum menambah pembelian, pastikan data berikut sudah tersedia:
        </p>
        <ul className="mb-6 text-sm">
          <li className={`flex items-center gap-2 ${hasSuppliers ? 'text-green-600' : 'text-amber-600'}`}>
            <span>{hasSuppliers ? '✓' : '!'}</span>
            {hasSuppliers ? 'Supplier sudah tersedia' : 'Belum ada supplier — tambahkan di menu Supplier'}
          </li>
          <li className={`flex items-center gap-2 ${hasProducts ? 'text-green-600' : 'text-amber-600'}`}>
            <span>{hasProducts ? '✓' : '!'}</span>
            {hasProducts ? 'Produk sudah tersedia' : 'Belum ada produk — tambahkan di menu Produk'}
          </li>
        </ul>
      </div>
    )
  }

  return <>{children}</>
}
