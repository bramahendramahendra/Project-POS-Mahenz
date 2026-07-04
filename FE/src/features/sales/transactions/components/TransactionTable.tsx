import { DataTable } from '@/shared/components'
import type { PaginationProps, SortState } from '@/shared/components/DataTable/DataTable.types'

import type { Transaction } from '../transactions.types'
import { buildTransactionColumns } from './TransactionTableColumns'

interface TransactionTableProps {
  data: Transaction[]
  isLoading: boolean
  pagination: PaginationProps
  currentSort?: SortState
  onSort?: (sort: SortState) => void
  onDetail: (transaction: Transaction) => void
  onVoid: (transaction: Transaction) => void
}

export function TransactionTable({
  data,
  isLoading,
  pagination,
  currentSort,
  onSort,
  onDetail,
  onVoid,
}: TransactionTableProps) {
  const columns = buildTransactionColumns({ onDetail, onVoid })

  return (
    <DataTable<Transaction & Record<string, unknown>>
      columns={columns}
      data={data as (Transaction & Record<string, unknown>)[]}
      isLoading={isLoading}
      currentSort={currentSort}
      onSort={onSort}
      emptyMessage="Belum ada transaksi"
      emptyDescription="Transaksi akan muncul setelah proses kasir selesai."
      pagination={pagination}
    />
  )
}
