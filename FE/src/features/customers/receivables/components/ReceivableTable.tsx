import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { Receivable } from '../receivables.types'
import { buildReceivableColumns } from './ReceivableTableColumns'

interface ReceivableTableProps {
  data: Receivable[]
  isLoading: boolean
  pagination: PaginationProps
  onPay: (receivable: Receivable) => void
}

export function ReceivableTable({ data, isLoading, pagination, onPay }: ReceivableTableProps) {
  const columns = buildReceivableColumns({ onPay })

  return (
    <DataTable<Receivable & Record<string, unknown>>
      columns={columns}
      data={data as (Receivable & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada piutang"
      emptyDescription="Piutang akan muncul saat transaksi dilakukan dengan metode kredit."
      pagination={pagination}
    />
  )
}
