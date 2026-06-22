import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { Shift, ShiftFormPayload, ShiftListFilter, ShiftOption } from './shifts.types'

export function useShiftListQuery(filter: ShiftListFilter) {
  return useQuery({
    queryKey: queryKeys.shifts.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Shift>>('/shifts/list', filter),
  })
}

export function useShiftOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.shifts.active(),
    queryFn: () => api.post<ShiftOption[]>('/shifts/active', {}),
  })
}

export function useShiftDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.shifts.detail(id),
    queryFn: () => api.post<Shift>(`/shifts/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useCreateShiftMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: ShiftFormPayload) => api.post<Shift>('/shifts/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateShiftMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: ShiftFormPayload & { id: number }) =>
      api.post<void>(`/shifts/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteShiftMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/shifts/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
      toast.success('Shift berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleShiftStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/shifts/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.shifts.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
