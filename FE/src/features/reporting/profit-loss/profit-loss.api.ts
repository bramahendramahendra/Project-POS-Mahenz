import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { ProfitLossReport } from './profit-loss.types'

export function useProfitLossReportQuery(filter: { date_from?: string; date_to?: string }) {
  return useQuery({
    queryKey: queryKeys.reports.profitLoss(filter as Record<string, unknown>),
    queryFn: () => api.get<ProfitLossReport>('/reports/profit-loss', filter),
  })
}
