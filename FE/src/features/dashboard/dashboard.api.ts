import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { RecentTransactionItem, TodaySummary } from './dashboard.types'

export function useRecentTransactionsQuery(limit = 5) {
  return useQuery({
    queryKey: queryKeys.dashboard.recentTransactions(),
    queryFn: () => api.get<RecentTransactionItem[]>('/dashboard/recent-transactions', { limit }),
  })
}

export function useTodaySummaryQuery() {
  return useQuery({
    queryKey: queryKeys.dashboard.todaySummary(),
    queryFn: () => api.get<TodaySummary>('/dashboard/today-summary'),
  })
}
