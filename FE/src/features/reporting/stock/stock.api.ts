import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { StockReport, StockReportFilter } from './stock.types'

interface StockReportResponse {
  items: StockReport[]
  total: number
  total_stock_value: number
}

export function useStockReportQuery(filter?: StockReportFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stock(filter as Record<string, unknown>),
    queryFn: () => api.get<StockReportResponse>('/reports/stock', filter),
  })
}
