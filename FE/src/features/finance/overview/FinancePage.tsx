import { useState } from 'react'

import { PageHeader } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useCashflowQuery, useFinanceSummaryQuery } from './finance.api'
import type { FinanceFilter } from './finance.types'
import { FinanceFilterBar } from './components/FinanceFilterBar'
import { FinanceSummaryCard } from './components/FinanceSummaryCard'
import { FinanceTable } from './components/FinanceTable'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

export function FinancePage() {
  const [filter, setFilter] = useState<FinanceFilter>({
    date_from: monthStartString(),
    date_to: todayString(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: summaryData, isLoading: summaryLoading } = useFinanceSummaryQuery({
    date_from: filter.date_from,
    date_to: filter.date_to,
  })
  const { data: cashflowData, isLoading: cashflowLoading } = useCashflowQuery({
    ...filter,
    page,
    page_size: pageSize,
  })

  const cashflows = cashflowData?.data ?? []
  const total = cashflowData?.total ?? 0

  return (
    <div className="space-y-4">
      <PageHeader title="Keuangan" breadcrumbs={[{ label: 'Finance' }, { label: 'Keuangan' }]} />

      <FinanceFilterBar filter={filter} onChange={setFilter} onReset={reset} />

      <FinanceSummaryCard summary={summaryData} isLoading={summaryLoading} />

      <div className="space-y-3">
        <h2 className="font-semibold text-gray-700">Arus Kas</h2>
        <FinanceTable
          data={cashflows}
          isLoading={cashflowLoading}
          pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        />
      </div>
    </div>
  )
}
