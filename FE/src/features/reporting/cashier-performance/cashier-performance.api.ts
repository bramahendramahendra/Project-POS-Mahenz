import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { CashierPerformance } from './cashier-performance.types'

export function useCashierPerformanceQuery(filter: { date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: queryKeys.reports.cashierPerformance(filter as Record<string, unknown>),
    queryFn: () => api.get<CashierPerformance[]>('/reports/cashier-performance', filter),
  })
}
