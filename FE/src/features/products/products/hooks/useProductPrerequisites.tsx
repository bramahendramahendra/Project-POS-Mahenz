import { Ruler, Tag } from 'lucide-react'

import { useCategoryOptionsQuery } from '@/features/products/categories'
import { useUnitOptionsQuery } from '@/features/products/units'

export function useProductPrerequisites() {
  const { data: categoriesData, isLoading: isCatLoading } = useCategoryOptionsQuery()
  const { data: unitsData, isLoading: isUnitLoading } = useUnitOptionsQuery()

  const hasCategories = (categoriesData ?? []).length > 0
  const hasActiveUnits = (unitsData ?? []).length > 0

  return {
    isLoading: isCatLoading || isUnitLoading,
    items: [
      {
        label: 'Belum ada kategori — tambahkan di tab Kategori',
        metLabel: 'Kategori sudah tersedia',
        met: hasCategories,
        icon: <Tag size={24} />,
      },
      {
        label: 'Belum ada satuan aktif — tambahkan di tab Satuan',
        metLabel: 'Satuan sudah tersedia',
        met: hasActiveUnits,
        icon: <Ruler size={24} />,
      },
    ],
  }
}
