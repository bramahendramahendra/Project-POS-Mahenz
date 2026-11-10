import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { ExpiryBatchHistory, ExpiryWarningsResponse } from './expiry-batches.types'

export function useExpiryWarningsQuery(search?: string) {
  return useQuery({
    queryKey: queryKeys.expiryBatches.warnings(search),
    queryFn: () => api.post<ExpiryWarningsResponse>('/expiry-batches/warnings', { search: search ?? '' }),
  })
}

export function useExpiryBatchHistoryQuery(productId?: number) {
  return useQuery({
    queryKey: queryKeys.expiryBatches.history(productId ?? 0),
    queryFn: () => api.post<ExpiryBatchHistory[]>(`/expiry-batches/product/${productId}`, {}),
    enabled: !!productId && productId > 0,
  })
}

export function useConfirmExpiryBatchMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, notes }: { id: number; notes?: string }) =>
      api.post<void>(`/expiry-batches/confirm/${id}`, { notes: notes ?? '' }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.expiryBatches.all() })
      toast.success('Batch dikonfirmasi aman, warning dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useWriteOffExpiryBatchMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, notes }: { id: number; notes?: string }) =>
      api.post<void>(`/expiry-batches/write-off/${id}`, { notes: notes ?? '' }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.expiryBatches.all() })
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success('Stok expired berhasil dimusnahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
