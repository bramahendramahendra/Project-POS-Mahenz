import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { SalesFilter, SalesListFilter, SalesReport, SalesReportSummary } from './sales.types'

export function useSalesListQuery(filter: SalesListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.salesList(filter as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<SalesReport>>('/reports/sales/list', filter),
  })
}

export function useSalesSummaryQuery(filter: SalesFilter) {
  return useQuery({
    queryKey: queryKeys.reports.salesSummary(filter as Record<string, unknown>),
    queryFn: () => api.post<SalesReportSummary>('/reports/sales/summary', filter),
  })
}
