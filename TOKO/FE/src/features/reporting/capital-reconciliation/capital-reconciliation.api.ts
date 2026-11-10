import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CapitalReconListFilter,
  CapitalReconListItem,
  CapitalReconSummary,
  ResolvePayload,
  SkipPayload,
} from './capital-reconciliation.types'

export function useCapitalReconListQuery(filter: CapitalReconListFilter) {
  return useQuery({
    queryKey: queryKeys.reports.capitalReconList(filter as unknown as Record<string, unknown>),
    queryFn: () =>
      api.post<PaginatedData<CapitalReconListItem>>('/capital-reconciliation/list', filter),
  })
}

export function useCapitalReconSummaryQuery() {
  return useQuery({
    queryKey: queryKeys.reports.capitalReconSummary(),
    queryFn: () => api.post<CapitalReconSummary>('/capital-reconciliation/summary', {}),
  })
}

function useInvalidateCapitalRecon() {
  const qc = useQueryClient()
  return () => {
    qc.invalidateQueries({ queryKey: ['reports', 'capitalRecon'] })
  }
}

export function useResolveCostMutation() {
  const invalidate = useInvalidateCapitalRecon()
  return useMutation({
    mutationFn: (payload: ResolvePayload) =>
      api.post<void>('/capital-reconciliation/resolve', payload),
    onSuccess: () => {
      invalidate()
      toast.success('Modal berhasil dikoreksi')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useSkipCostMutation() {
  const invalidate = useInvalidateCapitalRecon()
  return useMutation({
    mutationFn: (payload: SkipPayload) => api.post<void>('/capital-reconciliation/skip', payload),
    onSuccess: () => {
      invalidate()
      toast.success('Baris ditandai sudah ditinjau')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
