import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import { getWIBNow } from '@/shared/utils/date'

import type {
  DashboardPeriod,
  DashboardStats,
  SalesTrendItem,
  TopProductItem,
} from './business-summary.types'

function periodToTrendParam(period: DashboardPeriod): string {
  switch (period) {
    case 'month': return '30days'
    default:      return '7days'
  }
}

// Rentang kotak di bawah kartu angka (keputusan 134): Hari Ini = hari ini saja, Minggu Ini = 7 hari terakhir,
// Bulan Ini = awal bulan s.d. hari ini.
function periodToDateRange(period: DashboardPeriod): { start_date: string; end_date: string } {
  const now = getWIBNow()
  const end = now.format('YYYY-MM-DD')
  switch (period) {
    case 'today': return { start_date: end, end_date: end }
    case 'month': return { start_date: now.startOf('month').format('YYYY-MM-DD'), end_date: end }
    default:      return { start_date: now.subtract(6, 'day').format('YYYY-MM-DD'), end_date: end }
  }
}

export function useDashboardStatsQuery(period: DashboardPeriod) {
  return useQuery({
    queryKey: queryKeys.businessSummary.stats(period),
    queryFn: () => api.get<DashboardStats>('/reports/business-summary/stats', { period }),
  })
}

export function useSalesTrendQuery(period: DashboardPeriod) {
  return useQuery({
    queryKey: queryKeys.businessSummary.salesTrend(period),
    queryFn: () =>
      api.get<SalesTrendItem[]>('/reports/business-summary/sales-trend', { period: periodToTrendParam(period) }),
  })
}

export function useTopProductsQuery(period: DashboardPeriod) {
  const range = periodToDateRange(period)
  return useQuery({
    queryKey: queryKeys.businessSummary.topProducts(period),
    queryFn: () => api.get<TopProductItem[]>('/reports/business-summary/top-products', { ...range, limit: 10 }),
  })
}
