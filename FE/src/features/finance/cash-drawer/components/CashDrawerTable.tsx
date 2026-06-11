import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { CashDrawer } from '../cash-drawer.types'
import { buildCashDrawerColumns } from './CashDrawerTableColumns'

interface CashDrawerTableProps {
  data: CashDrawer[]
  isLoading: boolean
  pagination: PaginationProps
  onRowClick: (row: CashDrawer) => void
}

export function CashDrawerTable({ data, isLoading, pagination, onRowClick }: CashDrawerTableProps) {
  const columns = buildCashDrawerColumns({ onRowClick })

  return (
    <DataTable<CashDrawer & Record<string, unknown>>
      columns={columns}
      data={data as (CashDrawer & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada data kas harian"
      emptyDescription="Data kas harian akan muncul sesuai filter periode yang dipilih."
      pagination={pagination}
    />
  )
}
