import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { SyncConflict, SyncFilter, SyncHistoryItem } from './sync.types'

export function useSyncStatusQuery(enabled = true) {
  return useQuery({
    queryKey: queryKeys.sync.status(),
    queryFn: () => api.get<{ count: number }>('/sync/conflicts/count'),
    refetchInterval: 30_000,
    enabled,
  })
}

export function useSyncHistoryQuery(filter?: SyncFilter) {
  return useQuery({
    queryKey: queryKeys.sync.history(filter as unknown as Record<string, unknown>),
    queryFn: () => api.get<PaginatedData<SyncHistoryItem>>('/sync/history', filter),
  })
}

export function useSyncConflictsQuery() {
  return useQuery({
    queryKey: queryKeys.sync.conflicts(),
    queryFn: () => api.get<SyncConflict[]>('/sync/conflicts'),
  })
}

export function useResolveConflictMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, action }: { id: number; action: 'approve' | 'reject' }) =>
      api.post<void>(`/sync/conflicts/${id}/resolve`, { action }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.sync.conflicts() })
      qc.invalidateQueries({ queryKey: queryKeys.sync.status() })
      toast.success('Konflik berhasil diselesaikan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useTriggerSyncMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>('/sync/push', {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.sync.status() })
      qc.invalidateQueries({ queryKey: queryKeys.sync.history() })
      toast.success('Sinkronisasi manual dimulai')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
