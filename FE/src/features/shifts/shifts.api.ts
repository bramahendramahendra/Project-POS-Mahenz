import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { CloseShiftPayload, OpenShiftPayload, Shift, ShiftListFilter } from './shifts.types'

export function useShiftListQuery(filter?: ShiftListFilter) {
  return useQuery({
    queryKey: queryKeys.shifts.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Shift>>('/shifts/list', filter ?? {}),
  })
}

export function useActiveShiftQuery() {
  return useQuery({
    queryKey: queryKeys.shifts.active(),
    queryFn: () => api.post<Shift | null>('/shifts/active', {}),
  })
}

export function useShiftDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.shifts.detail(id),
    queryFn: () => api.post<Shift>(`/shifts/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useOpenShiftMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: OpenShiftPayload) => api.post<Shift>('/shifts/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
      qc.invalidateQueries({ queryKey: queryKeys.shifts.active() })
      toast.success('Shift berhasil dibuka')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useCloseShiftMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: CloseShiftPayload }) =>
      api.post<Shift>(`/shifts/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
      qc.invalidateQueries({ queryKey: queryKeys.shifts.active() })
      toast.success('Shift berhasil ditutup')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
