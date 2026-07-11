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

// /sync/push adalah kontrak untuk CLIENT offline (desktop/Android) mem-push antrian lokalnya
// sendiri — endpoint ini mewajibkan device_id + minimal 1 item, dan dashboard web ini tidak
// pernah punya antrian lokal untuk dikirim. Jadi tombol di halaman ini TIDAK memanggil
// /sync/push (dulu begitu dan selalu gagal 400); ia hanya refetch semua data sync yang
// relevan supaya user bisa lihat status terbaru tanpa menunggu polling 30 detik.
export function useRefreshSyncMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async () => {
      await Promise.all([
        qc.invalidateQueries({ queryKey: queryKeys.sync.status() }),
        qc.invalidateQueries({ queryKey: queryKeys.sync.conflicts() }),
        qc.invalidateQueries({ queryKey: ['sync', 'queue'] }),
        qc.invalidateQueries({ queryKey: ['sync', 'history'] }),
      ])
    },
    onSuccess: () => {
      toast.success('Status sinkronisasi diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
