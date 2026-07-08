import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, downloadReportExport } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { SalesFilter, SalesListFilter, SalesReport, SalesReportSummary } from './sales.types'

export function useSalesListQuery(filter: SalesListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.salesList(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<SalesReport>>('/reports/sales/list', filter),
  })
}

export function useSalesSummaryQuery(filter: SalesFilter) {
  return useQuery({
    queryKey: queryKeys.reports.salesSummary(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<SalesReportSummary>('/reports/sales/summary', filter),
  })
}

export function useExportSalesReportMutation() {
  return useMutation({
    mutationFn: (filter: SalesFilter) =>
      downloadReportExport(
        '/reports/sales/export',
        { date_from: filter.date_from, date_to: filter.date_to },
        'laporan-penjualan.xlsx'
      ),
    onError: (e: Error) => toast.error(e.message),
  })
}
