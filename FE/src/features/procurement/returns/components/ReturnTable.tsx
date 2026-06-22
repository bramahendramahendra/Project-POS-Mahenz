import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierReturn } from '../returns.types'
import { buildReturnColumns } from './ReturnTableColumns'

interface ReturnTableProps {
  data: SupplierReturn[]
  isLoading: boolean
  pagination: PaginationProps
  onDetail: (row: SupplierReturn) => void
  onDelete: (row: SupplierReturn) => void
}

export function ReturnTable({ data, isLoading, pagination, onDetail, onDelete }: ReturnTableProps) {
  const columns = buildReturnColumns({ onDetail, onDelete })

  return (
    <DataTable<SupplierReturn & Record<string, unknown>>
      columns={columns}
      data={data as (SupplierReturn & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada data retur"
      emptyDescription="Data retur pembelian akan muncul sesuai filter yang dipilih."
      pagination={pagination}
    />
  )
}
