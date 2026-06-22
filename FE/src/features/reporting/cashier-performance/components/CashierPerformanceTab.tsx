import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { monthStart, todayStr } from '@/shared/utils'

import { useCashierPerformanceQuery } from '../cashier-performance.api'
import type { CashierPerformance } from '../cashier-performance.types'
import { CashierPerformanceFilterBar } from './CashierPerformanceFilterBar'
import { buildCashierPerformanceColumns } from './CashierPerformanceTableColumns'

interface DateFilter {
  date_from?: string
  date_to?: string
}

export function CashierPerformanceTab() {
  const [filter, setFilter] = useState<DateFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })

  const { data, isLoading } = useCashierPerformanceQuery({
    date_from: filter.date_from,
    date_to: filter.date_to,
  })

  const items: CashierPerformance[] = (data ?? [])
    .slice()
    .sort((a, b) => b.total_sales - a.total_sales)

  const columns = buildCashierPerformanceColumns()

  return (
    <div className="space-y-4">
      <CashierPerformanceFilterBar filter={filter} onChange={setFilter} />

      <DataTable<CashierPerformance & Record<string, unknown>>
        columns={columns}
        data={items as (CashierPerformance & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data kinerja kasir"
        emptyDescription="Data kinerja kasir akan muncul sesuai filter periode yang dipilih."
      />
    </div>
  )
}
