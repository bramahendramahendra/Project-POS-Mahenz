import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { ProfitLossDateFilter, ProfitLossReport } from './profit-loss.types'

export function useProfitLossReportQuery(filter: ProfitLossDateFilter) {
  return useQuery({
    queryKey: queryKeys.reports.profitLoss(filter as Record<string, unknown>),
    queryFn: () => api.post<ProfitLossReport>('/reports/profit-loss/data', filter),
  })
}
