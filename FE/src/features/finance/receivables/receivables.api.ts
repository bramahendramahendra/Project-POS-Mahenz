import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { CreatePaymentPayload, Receivable, ReceivableListFilter } from './receivables.types'

export function useReceivableListQuery(filter?: ReceivableListFilter) {
  return useQuery({
    queryKey: queryKeys.receivables.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Receivable>>('/receivables/list', filter ?? {}),
  })
}

export function useReceivableDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.receivables.detail(id),
    queryFn: () => api.post<Receivable>(`/receivables/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useAddPaymentMutation(receivableId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreatePaymentPayload) =>
      api.post<Receivable>(`/receivables/pay/${receivableId}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.receivables.all() })
      qc.invalidateQueries({ queryKey: queryKeys.receivables.detail(receivableId) })
      toast.success('Pembayaran berhasil dicatat')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
