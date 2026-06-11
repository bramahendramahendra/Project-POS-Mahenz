import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { Transaction } from '../transactions.types'
import { buildTransactionColumns } from './TransactionTableColumns'

interface TransactionTableProps {
  data: Transaction[]
  isLoading: boolean
  pagination: PaginationProps
  onDetail: (transaction: Transaction) => void
  onVoid: (transaction: Transaction) => void
}

export function TransactionTable({
  data,
  isLoading,
  pagination,
  onDetail,
  onVoid,
}: TransactionTableProps) {
  const columns = buildTransactionColumns({ onDetail, onVoid })

  return (
    <DataTable<Transaction & Record<string, unknown>>
      columns={columns}
      data={data as (Transaction & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada transaksi"
      emptyDescription="Transaksi akan muncul setelah proses kasir selesai."
      pagination={pagination}
    />
  )
}
