import { useState } from 'react'

import { PageHeader } from '@/shared/components'

import { useDashboardStatsQuery, useSalesTrendQuery, useTopProductsQuery } from './business-summary.api'
import type { DashboardPeriod } from './business-summary.types'
import { DashboardPeriodSelector } from './components/DashboardPeriodSelector'
import { SalesChart } from './components/SalesChart'
import { SummaryCards } from './components/SummaryCards'
import { TopProductsTable } from './components/TopProductsTable'

export function BusinessSummaryPage() {
  const [period, setPeriod] = useState<DashboardPeriod>('today')

  const { data: statsData, isLoading: statsLoading } = useDashboardStatsQuery(period)
  const { data: trendData, isLoading: trendLoading } = useSalesTrendQuery(period)
  const { data: topData, isLoading: topLoading } = useTopProductsQuery(period)

  const trendPoints = trendData ?? []
  const topProducts = topData ?? []

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between flex-wrap gap-3">
        <PageHeader title="Ringkasan Bisnis" breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Ringkasan Bisnis' }]} />
        <DashboardPeriodSelector period={period} onChange={setPeriod} />
      </div>

      <SummaryCards stats={statsData} isLoading={statsLoading} period={period} />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 rounded-lg border bg-white p-4 shadow-sm space-y-3">
          <h3 className="font-semibold text-gray-700 text-sm">Grafik Penjualan</h3>
          <SalesChart data={trendPoints} isLoading={trendLoading} />
        </div>
        <div className="lg:col-span-1">
          <TopProductsTable data={topProducts} isLoading={topLoading} />
        </div>
      </div>
    </div>
  )
}
