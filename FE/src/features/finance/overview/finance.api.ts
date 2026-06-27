import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { CashflowFilter, CashflowItem, FinanceDateFilter, FinanceSummary } from './finance.types'

export function useFinanceSummaryQuery(filter?: FinanceDateFilter) {
  return useQuery({
    queryKey: queryKeys.finance.summary(filter as Record<string, unknown>),
    queryFn: () => api.post<FinanceSummary>('/finance/summary', filter ?? {}),
  })
}

export function useCashflowQuery(filter: CashflowFilter) {
  return useQuery({
    queryKey: queryKeys.finance.cashflow(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<CashflowItem>>('/finance/cashflow', filter),
  })
}
