import { FileText } from 'lucide-react'

import { useSupplierPurchasesQuery } from '@/features/procurement/purchases'

export function useReturnPrerequisites() {
  const { data: purchasesData, isLoading } = useSupplierPurchasesQuery({ page: 1, limit: 1 })

  const hasPurchases = (purchasesData?.total ?? 0) > 0

  return {
    isLoading,
    items: [
      {
        label: 'Belum ada faktur pembelian — tambahkan di menu Pembelian',
        metLabel: 'Faktur pembelian sudah tersedia',
        met: hasPurchases,
        icon: <FileText size={24} />,
      },
    ],
  }
}
