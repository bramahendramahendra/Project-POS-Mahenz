import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, downloadReportExport } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CashierPerformance,
  CashierPerformanceDateFilter,
  CashierPerformanceListFilter,
} from './cashier-performance.types'

export function useCashierPerformanceListQuery(filter: CashierPerformanceListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.cashierPerformanceList(filter as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<CashierPerformance>>('/reports/cashier/list', filter),
  })
}

export function useExportCashierPerformanceMutation() {
  return useMutation({
    mutationFn: (filter: CashierPerformanceDateFilter) =>
      downloadReportExport(
        '/reports/cashier/export',
        { date_from: filter.date_from, date_to: filter.date_to },
        'laporan-kasir.xlsx'
      ),
    onError: (e: Error) => toast.error(e.message),
  })
}
