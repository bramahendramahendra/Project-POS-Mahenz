import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  AdjustPayload,
  MarkReviewedPayload,
  ReconDetail,
  ReconFilter,
  ReconListFilter,
  ReconListItem,
  ReconSummary,
} from './stock-reconciliation.types'

export function useReconListQuery(filter: ReconListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.stockReconList(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<ReconListItem>>('/stock-reconciliation/list', filter),
  })
}

export function useReconSummaryQuery() {
  return useQuery({
    queryKey: queryKeys.reports.stockReconSummary(),
    queryFn: () => api.post<ReconSummary>('/stock-reconciliation/summary', {}),
  })
}

export function useReconDetailQuery(productId: number | null) {
  return useQuery({
    queryKey: queryKeys.reports.stockReconDetail(productId ?? 0),
    queryFn: () => api.post<ReconDetail>(`/stock-reconciliation/detail/${productId}`, {}),
    enabled: productId != null,
  })
}

function useInvalidateRecon() {
  const qc = useQueryClient()
  return (productId: number) => {
    qc.invalidateQueries({ queryKey: ['reports', 'stockRecon'] })
    qc.invalidateQueries({ queryKey: queryKeys.reports.stockReconDetail(productId) })
  }
}

export function useAdjustStockMutation() {
  const invalidate = useInvalidateRecon()
  return useMutation({
    mutationFn: (payload: AdjustPayload) =>
      api.post<ReconDetail>('/stock-reconciliation/adjust', payload),
    onSuccess: (_data, payload) => {
      invalidate(payload.product_id)
      toast.success('Stok berhasil dikoreksi')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useMarkReviewedMutation() {
  const invalidate = useInvalidateRecon()
  return useMutation({
    mutationFn: (payload: MarkReviewedPayload) =>
      api.post<ReconDetail>('/stock-reconciliation/mark-reviewed', payload),
    onSuccess: (_data, payload) => {
      invalidate(payload.product_id)
      toast.success('Produk ditandai sudah ditinjau')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
