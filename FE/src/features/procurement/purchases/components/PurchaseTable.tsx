import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierPurchase } from '../purchases.types'
import { buildPurchaseColumns } from './PurchaseTableColumns'

interface PurchaseTableProps {
  data: SupplierPurchase[]
  isLoading: boolean
  pagination: PaginationProps
  onDetail: (purchase: SupplierPurchase) => void
  onEdit: (purchase: SupplierPurchase) => void
  onPay: (purchase: SupplierPurchase) => void
  onDelete: (purchase: SupplierPurchase) => void
}

export function PurchaseTable({
  data,
  isLoading,
  pagination,
  onDetail,
  onEdit,
  onPay,
  onDelete,
}: PurchaseTableProps) {
  const columns = buildPurchaseColumns({ onDetail, onEdit, onPay, onDelete })

  return (
    <DataTable<SupplierPurchase & Record<string, unknown>>
      columns={columns}
      data={data as (SupplierPurchase & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada data pembelian"
      emptyDescription="Data pembelian supplier akan muncul sesuai filter yang dipilih."
      pagination={pagination}
    />
  )
}
