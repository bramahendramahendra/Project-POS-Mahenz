import { Tag } from 'lucide-react'

import { useCategoryOptionsQuery } from '@/features/products/categories'

// Syarat Laporan Stok dan Rekonsiliasi Stok: kategori produk harus ada dulu
// (pola sama dengan useProductPrerequisites).
export function useStockReportPrerequisites() {
  const { data: categoriesData, isLoading } = useCategoryOptionsQuery()
  const hasCategories = (categoriesData ?? []).length > 0

  return {
    isLoading,
    items: [
      {
        label: 'Belum ada kategori — tambahkan di menu Barang Jadi → Kategori',
        metLabel: 'Kategori sudah tersedia',
        met: hasCategories,
        icon: <Tag size={24} />,
      },
    ],
  }
}
