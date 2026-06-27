import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import { monthStart, todayStr } from '@/shared/utils'

import { useSalesListQuery, useSalesSummaryQuery } from '../sales.api'
import type { SalesFilter, SalesReport } from '../sales.types'
import { SalesReportFilterBar } from './SalesReportFilterBar'
import { SalesReportSummaryCard } from './SalesReportSummaryCard'
import { buildSalesReportColumns } from './SalesReportTableColumns'

export function SalesReportTab() {
  const [filter, setFilter] = useState<SalesFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: listData, isLoading: listLoading } = useSalesListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const { data: summary, isLoading: summaryLoading } = useSalesSummaryQuery(filter)

  const items: SalesReport[] = listData?.data ?? []
  const total = listData?.total ?? 0

  const handleFilterChange = (newFilter: SalesFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({ date_from: monthStart(), date_to: todayStr() })
    reset()
  }

  const columns = buildSalesReportColumns()

  return (
    <div className="space-y-4">
      <SalesReportFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        exportData={items}
      />

      <SalesReportSummaryCard summary={summary} isLoading={summaryLoading} />

      <DataTable<SalesReport & Record<string, unknown>>
        columns={columns}
        data={items as (SalesReport & Record<string, unknown>)[]}
        isLoading={listLoading}
        emptyMessage="Belum ada data penjualan"
        emptyDescription="Data penjualan akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
    </div>
  )
}
