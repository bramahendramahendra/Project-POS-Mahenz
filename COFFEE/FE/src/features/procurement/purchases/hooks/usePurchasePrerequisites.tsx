import { Package, Truck } from 'lucide-react'

import { useIngredientOptionsQuery } from '@/features/ingredients/ingredients.api'
import { useSupplierListQuery } from '@/features/procurement/suppliers'
import { useProductListQuery } from '@/features/products/products'

// Syarat halaman Pembelian Coffee: supplier, dan barang yang bisa dibeli. Nota Coffee berisi bahan baku
// dan/atau barang jadi (keputusan 91-94), jadi cukup salah satu: produk ATAU bahan baku aktif (keputusan 126).
export function usePurchasePrerequisites() {
  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 1, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })
  const { data: ingredientsData, isLoading: isIngredientsLoading } = useIngredientOptionsQuery()

  const hasSuppliers = (suppliersData?.total ?? 0) > 0
  const hasItems = (productsData?.total ?? 0) > 0 || (ingredientsData ?? []).length > 0

  return {
    isLoading: isSuppliersLoading || isProductsLoading || isIngredientsLoading,
    items: [
      {
        label: 'Belum ada supplier — tambahkan di menu Supplier',
        metLabel: 'Supplier sudah tersedia',
        met: hasSuppliers,
        icon: <Truck size={24} />,
      },
      {
        label: 'Belum ada barang yang bisa dibeli — tambahkan bahan baku (menu Daftar Bahan) atau produk (menu Barang Jadi → Produk)',
        metLabel: 'Bahan baku atau produk sudah tersedia',
        met: hasItems,
        icon: <Package size={24} />,
      },
    ],
  }
}
