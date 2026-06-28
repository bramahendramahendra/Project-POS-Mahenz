import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { CashierPerformance, CashierPerformanceListFilter } from './cashier-performance.types'

export function useCashierPerformanceListQuery(filter: CashierPerformanceListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.cashierPerformanceList(filter as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<CashierPerformance>>('/reports/cashier/list', filter),
  })
}
