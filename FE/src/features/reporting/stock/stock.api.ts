import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, downloadReportExport } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { StockFilter, StockListFilter, StockReport, StockSummary } from './stock.types'

export function useStockListQuery(filter: StockListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stockList(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<StockReport>>('/reports/stock/list', filter),
  })
}

export function useStockSummaryQuery(filter: StockFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stockSummary(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<StockSummary>('/reports/stock/summary', filter),
  })
}

export function useExportStockReportMutation() {
  return useMutation({
    mutationFn: () => downloadReportExport('/reports/stock/export', {}, 'laporan-stok.xlsx'),
    onError: (e: Error) => toast.error(e.message),
  })
}
