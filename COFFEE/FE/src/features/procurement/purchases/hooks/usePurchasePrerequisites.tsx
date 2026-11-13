import { Package, Truck } from 'lucide-react'

import { useSupplierListQuery } from '@/features/procurement/suppliers'
import { useProductListQuery } from '@/features/products/products'

export function usePurchasePrerequisites() {
  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 1, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })

  const hasSuppliers = (suppliersData?.total ?? 0) > 0
  const hasProducts = (productsData?.total ?? 0) > 0

  return {
    isLoading: isSuppliersLoading || isProductsLoading,
    items: [
      {
        label: 'Belum ada supplier — tambahkan di menu Supplier',
        metLabel: 'Supplier sudah tersedia',
        met: hasSuppliers,
        icon: <Truck size={24} />,
      },
      {
        label: 'Belum ada produk — tambahkan di menu Produk',
        metLabel: 'Produk sudah tersedia',
        met: hasProducts,
        icon: <Package size={24} />,
      },
    ],
  }
}
