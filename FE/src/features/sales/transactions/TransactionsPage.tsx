import { useState } from 'react'

import { PageHeader } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useTransactionListQuery } from './transactions.api'
import type { TransactionListFilter } from './transactions.types'
import { TransactionFilterBar } from './components/TransactionFilterBar'
import { TransactionTable } from './components/TransactionTable'
import { TransactionDetailModal } from './components/TransactionDetailModal'

const DEFAULT_FILTER: TransactionListFilter = {}

export function TransactionsPage() {
  const [filter, setFilter] = useState<TransactionListFilter>(DEFAULT_FILTER)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)
  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()

  const pageSizeOptions = usePageSizeOptions()
  const { data: txData, isLoading } = useTransactionListQuery({
    ...filter,
    page,
    limit: pageSize,
  })

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

  return (
    <div className="space-y-4">
      <PageHeader
        title="Transaksi"
        breadcrumbs={[{ label: 'Penjualan' }, { label: 'Transaksi' }]}
      />

      <TransactionFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <TransactionTable
        data={transactions}
        isLoading={isLoading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        currentSort={sortState}
        onSort={handleSort}
        onDetail={(t) => setSelectedId(t.id)}
        onVoid={(t) => setSelectedId(t.id)}
      />

      <TransactionDetailModal transactionId={selectedId} onClose={() => setSelectedId(null)} />
    </div>
  )
}
