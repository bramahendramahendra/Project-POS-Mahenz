import type { ReactNode } from 'react'
import { FileText } from 'lucide-react'

import { useSupplierPurchasesQuery } from '@/features/procurement/purchases'

interface ReturnPrerequisiteGuardProps {
  children: ReactNode
}

export function ReturnPrerequisiteGuard({ children }: ReturnPrerequisiteGuardProps) {
  const { data: purchasesData, isLoading: isPurchasesLoading } = useSupplierPurchasesQuery({ page: 1, limit: 1 })

  if (isPurchasesLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  const hasPurchases = (purchasesData?.total ?? 0) > 0

  if (!hasPurchases) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4 rounded-full bg-amber-50 p-3">
          <FileText size={24} className="text-amber-500" />
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa menambah retur</h3>
        <p className="mb-1 text-sm text-gray-500">
          Retur pembelian dibuat berdasarkan faktur pembelian yang sudah ada.
        </p>
        <ul className="mb-6 text-sm">
          <li className="flex items-center gap-2 text-amber-600">
            <span>!</span>
            Belum ada faktur pembelian — tambahkan di menu Pembelian
          </li>
        </ul>
      </div>
    )
  }

  return <>{children}</>
}
