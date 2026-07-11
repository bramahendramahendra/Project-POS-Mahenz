import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  ConflictListResponse,
  QueueListResponse,
  SyncFilter,
  SyncHistoryListResponse,
} from './sync.types'

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
    queryFn: () => api.get<SyncHistoryListResponse>('/sync/history', filter),
  })
}

export function useSyncConflictsQuery() {
  return useQuery({
    queryKey: queryKeys.sync.conflicts(),
    queryFn: () => api.get<ConflictListResponse>('/sync/conflicts'),
  })
}

export function useSyncQueueQuery(filter?: SyncFilter) {
  return useQuery({
    queryKey: queryKeys.sync.queue(filter as unknown as Record<string, unknown>),
    queryFn: () => api.get<QueueListResponse>('/sync/queue', filter),
  })
}

// Semantik backend (BE/domain/sync/service/sync_service.go): 'approve' = terapkan data
// offline/lokal ke server, 'reject' = buang data offline, pertahankan data server (dan untuk
// transaksi, kembalikan stok yang sempat dikurangi). Ini kebalikan dari istilah "approve"
// sehari-hari, jadi mapping ke pilihan user ("terima data server" vs "pakai data lokal")
// dibalik di sini — lihat ConflictList.tsx.
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
      // Query nyata selalu menyertakan filter {page, limit} di key (lihat
      // queryKeys.sync.history di atas) — invalidate pakai prefix tanpa filter
      // supaya cocok dengan SEMUA variasi filter yang mungkin ter-cache, bukan
      // hanya key persis ['sync','history',undefined] yang tidak pernah ada.
      qc.invalidateQueries({ queryKey: ['sync', 'history'] })
      toast.success('Sinkronisasi manual dimulai')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
