import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, downloadReportExport } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { ProfitLossDateFilter, ProfitLossReport } from './profit-loss.types'

export function useProfitLossReportQuery(filter: ProfitLossDateFilter) {
  return useQuery({
    queryKey: queryKeys.reports.profitLoss(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<ProfitLossReport>('/reports/profit-loss/data', filter),
  })
}

export function useExportProfitLossMutation() {
  return useMutation({
    mutationFn: (filter: ProfitLossDateFilter) =>
      downloadReportExport(
        '/reports/profit-loss/export',
        { date_from: filter.date_from, date_to: filter.date_to },
        'laporan-laba-rugi.xlsx'
      ),
    onError: (e: Error) => toast.error(e.message),
  })
}
