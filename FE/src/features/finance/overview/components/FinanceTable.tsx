import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { CashflowItem } from '../finance.types'
import { buildFinanceColumns } from './FinanceTableColumns'

interface FinanceTableProps {
  data: CashflowItem[]
  isLoading: boolean
  pagination: PaginationProps
}

export function FinanceTable({ data, isLoading, pagination }: FinanceTableProps) {
  const columns = buildFinanceColumns()

  return (
    <DataTable<CashflowItem & Record<string, unknown>>
      columns={columns}
      data={data as (CashflowItem & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada data arus kas"
      emptyDescription="Data arus kas akan muncul sesuai filter periode yang dipilih."
      pagination={pagination}
    />
  )
}
