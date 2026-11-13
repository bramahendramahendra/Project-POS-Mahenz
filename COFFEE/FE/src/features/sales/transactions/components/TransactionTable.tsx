import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useTransactionListQuery } from '../transactions.api'
import type { Transaction, TransactionListFilter } from '../transactions.types'
import { TransactionDetailModal } from './TransactionDetailModal'
import { TransactionFilterBar } from './TransactionFilterBar'
import { buildTransactionColumns } from './TransactionTableColumns'

const DEFAULT_FILTER: TransactionListFilter = { page: 1, limit: 10 }

export function TransactionTable() {
  const [filter, setFilter] = useState<TransactionListFilter>(DEFAULT_FILTER)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)
  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()

  const pageSizeOptions = usePageSizeOptions()
  const { data: txData, isLoading } = useTransactionListQuery({ ...filter, page, limit: pageSize })

  const transactions = txData?.data ?? []
  const total = txData?.total ?? 0

  const handleFilterChange = (newFilter: TransactionListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter(DEFAULT_FILTER)
    setSortState(undefined)
    reset()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    reset()
  }

  const columns = buildTransactionColumns({
    onDetail: (t: Transaction) => setSelectedId(t.id),
    onVoid: (t: Transaction) => setSelectedId(t.id),
  })

  return (
    <div className="space-y-4">
      <TransactionFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Transaction & Record<string, unknown>>
        columns={columns}
        data={transactions as (Transaction & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada transaksi"
        emptyDescription="Transaksi akan muncul setelah proses kasir selesai."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <TransactionDetailModal transactionId={selectedId} onClose={() => setSelectedId(null)} />
    </div>
  )
}
