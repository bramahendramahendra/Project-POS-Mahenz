import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { StockFilter, StockListFilter, StockReport, StockSummary } from './stock.types'

export function useStockListQuery(filter: StockListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stockList(filter as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<StockReport>>('/reports/stock/list', filter),
  })
}

export function useStockSummaryQuery(filter: StockFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stockSummary(filter as Record<string, unknown>),
    queryFn: () => api.post<StockSummary>('/reports/stock/summary', filter),
  })
}
