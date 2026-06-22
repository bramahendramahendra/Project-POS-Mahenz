import type { ReactNode } from 'react'
import { Ruler, Tag } from 'lucide-react'

import { useCategoryOptionsQuery } from '@/features/products/categories'
import { useUnitOptionsQuery } from '@/features/products/units'

interface ProductPrerequisiteGuardProps {
  children: ReactNode
}

export function ProductPrerequisiteGuard({ children }: ProductPrerequisiteGuardProps) {
  const { data: categories = [], isLoading: isCatLoading } = useCategoryOptionsQuery()
  const { data: units = [], isLoading: isUnitLoading } = useUnitOptionsQuery()

  if (isCatLoading || isUnitLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  const hasCategories = categories.length > 0
  const hasActiveUnits = units.length > 0

  if (!hasCategories || !hasActiveUnits) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4 flex gap-3">
          <div className={`rounded-full p-3 ${!hasCategories ? 'bg-amber-50' : 'bg-green-50'}`}>
            <Tag size={24} className={!hasCategories ? 'text-amber-500' : 'text-green-500'} />
          </div>
          <div className={`rounded-full p-3 ${!hasActiveUnits ? 'bg-amber-50' : 'bg-green-50'}`}>
            <Ruler size={24} className={!hasActiveUnits ? 'text-amber-500' : 'text-green-500'} />
          </div>
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa menambah produk</h3>
        <p className="mb-1 text-sm text-gray-500">
          Sebelum menambah produk, pastikan data berikut sudah tersedia:
        </p>
        <ul className="mb-6 text-sm">
          <li className={`flex items-center gap-2 ${hasCategories ? 'text-green-600' : 'text-amber-600'}`}>
            <span>{hasCategories ? '✓' : '!'}</span>
            {hasCategories ? 'Kategori sudah tersedia' : 'Belum ada kategori — tambahkan di tab Kategori'}
          </li>
          <li className={`flex items-center gap-2 ${hasActiveUnits ? 'text-green-600' : 'text-amber-600'}`}>
            <span>{hasActiveUnits ? '✓' : '!'}</span>
            {hasActiveUnits ? 'Satuan sudah tersedia' : 'Belum ada satuan aktif — tambahkan di tab Satuan'}
          </li>
        </ul>
      </div>
    )
  }

  return <>{children}</>
}
