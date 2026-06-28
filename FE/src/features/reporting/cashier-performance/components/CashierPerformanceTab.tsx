import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { monthStart, todayStr } from '@/shared/utils'

import { useCashierPerformanceListQuery } from '../cashier-performance.api'
import type { CashierPerformance, CashierPerformanceDateFilter, CashierPerformanceListFilter } from '../cashier-performance.types'
import { CashierPerformanceFilterBar } from './CashierPerformanceFilterBar'
import { buildCashierPerformanceColumns } from './CashierPerformanceTableColumns'

export function CashierPerformanceTab() {
  const [filter, setFilter] = useState<CashierPerformanceDateFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const listFilter: CashierPerformanceListFilter = { ...filter, page, limit: pageSize }
  const { data, isLoading } = useCashierPerformanceListQuery(listFilter)
  const items: CashierPerformance[] = data?.data ?? []
  const total = data?.total ?? 0

  const handleFilterChange = (newFilter: CashierPerformanceDateFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ date_from: monthStart(), date_to: todayStr() })
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const columns = buildCashierPerformanceColumns()

  return (
    <div className="space-y-4">
      <CashierPerformanceFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
      />

      <DataTable<CashierPerformance & Record<string, unknown>>
        columns={columns}
        data={items as (CashierPerformance & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data kinerja kasir"
        emptyDescription="Data kinerja kasir akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
    </div>
  )
}
