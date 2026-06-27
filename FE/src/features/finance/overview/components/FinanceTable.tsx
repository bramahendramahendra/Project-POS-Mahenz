import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import { monthStart, todayStr } from '@/shared/utils'

import { useCashflowQuery, useFinanceSummaryQuery } from '../finance.api'
import type { CashflowFilter, CashflowItem, FinanceDateFilter } from '../finance.types'
import { FinanceFilterBar } from './FinanceFilterBar'
import { FinanceSummaryCard } from './FinanceSummaryCard'
import { buildFinanceColumns } from './FinanceTableColumns'

const defaultDateFilter: FinanceDateFilter = {
  date_from: monthStart(),
  date_to: todayStr(),
}

export function FinanceTable() {
  const [dateFilter, setDateFilter] = useState<FinanceDateFilter>(defaultDateFilter)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: summaryData, isLoading: summaryLoading } = useFinanceSummaryQuery(dateFilter)

  const cashflowFilter: CashflowFilter = {
    ...dateFilter,
    page,
    limit: pageSize,
  }
  const { data: cashflowData, isLoading: cashflowLoading } = useCashflowQuery(cashflowFilter)

  const cashflows = cashflowData?.data ?? []
  const total = cashflowData?.total ?? 0

  const handleFilterChange = (newFilter: FinanceDateFilter) => {
    setDateFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setDateFilter({ date_from: monthStart(), date_to: todayStr() })
    reset()
  }

  const columns = buildFinanceColumns()

  return (
    <div className="space-y-4">
      <FinanceFilterBar filter={dateFilter} onChange={handleFilterChange} onReset={handleReset} />

      <FinanceSummaryCard summary={summaryData} isLoading={summaryLoading} />

      <div className="space-y-3">
        <h2 className="font-semibold text-gray-700">Arus Kas</h2>
        <DataTable<CashflowItem & Record<string, unknown>>
          columns={columns}
          data={cashflows as (CashflowItem & Record<string, unknown>)[]}
          isLoading={cashflowLoading}
          emptyMessage="Belum ada data arus kas"
          emptyDescription="Data arus kas akan muncul sesuai filter periode yang dipilih."
          pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        />
      </div>
    </div>
  )
}
