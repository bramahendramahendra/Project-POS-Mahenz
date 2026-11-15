import { Ruler } from 'lucide-react'

import { useUnitOptionsQuery } from '@/features/products/units'

// Syarat halaman Daftar Bahan: satuan harus ada dulu (pola sama dengan useProductPrerequisites).
export function useIngredientPrerequisites() {
  const { data: unitsData, isLoading } = useUnitOptionsQuery()
  const hasActiveUnits = (unitsData ?? []).length > 0

  return {
    isLoading,
    items: [
      {
        label: 'Belum ada satuan aktif — tambahkan di menu Barang Jadi → Unit',
        metLabel: 'Satuan sudah tersedia',
        met: hasActiveUnits,
        icon: <Ruler size={24} />,
      },
    ],
  }
}
