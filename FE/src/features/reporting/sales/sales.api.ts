import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { SalesReport, SalesReportFilter, SalesReportSummary } from './sales.types'

interface SalesReportResponse {
  items: SalesReport[]
  total: number
  summary: SalesReportSummary
}

export function useSalesReportQuery(filter?: SalesReportFilter) {
  return useQuery({
    queryKey: queryKeys.reports.sales(filter as Record<string, unknown>),
    queryFn: () => api.get<SalesReportResponse>('/reports/sales', filter),
  })
}
