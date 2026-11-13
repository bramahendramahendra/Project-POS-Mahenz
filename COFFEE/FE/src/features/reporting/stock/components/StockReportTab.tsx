import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useStockListQuery, useStockSummaryQuery } from '../stock.api'
import type { StockFilter, StockReport } from '../stock.types'
import { StockReportFilterBar } from './StockReportFilterBar'
import { StockReportSummaryCard } from './StockReportSummaryCard'
import { buildStockReportColumns } from './StockReportTableColumns'

export function StockReportTab() {
  const [filter, setFilter] = useState<StockFilter>({})
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)
  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: listData, isLoading: listLoading } = useStockListQuery({ ...filter, page, limit: pageSize })
  const { data: summary, isLoading: summaryLoading } = useStockSummaryQuery(filter)

  const items: StockReport[] = listData?.data ?? []
  const total = listData?.total ?? 0

  const handleFilterChange = (newFilter: StockFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({})
    setSortState(undefined)
    reset()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    reset()
  }

  const columns = buildStockReportColumns()

  return (
    <div className="space-y-4">
      <StockReportFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <StockReportSummaryCard summary={summary} isLoading={summaryLoading} />

      <DataTable<StockReport & Record<string, unknown>>
        columns={columns}
        data={items as (StockReport & Record<string, unknown>)[]}
        isLoading={listLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data stok"
        emptyDescription="Data stok produk akan muncul sesuai filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
    </div>
  )
}
