import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateUnitPayload,
  Unit,
  UnitListFilter,
  UnitOption,
  UpdateUnitPayload,
} from './units.types'

export function useUnitListQuery(filter: UnitListFilter) {
  return useQuery({
    queryKey: queryKeys.units.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Unit>>('/units/list', filter),
  })
}

export function useUnitOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.units.options(),
    queryFn: () => api.post<UnitOption[]>('/units/options', {}),
  })
}

export function useCreateUnitMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateUnitPayload) => api.post<Unit>('/units/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.units.all() })
      toast.success('Satuan berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateUnitMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateUnitPayload & { id: number }) =>
      api.post<Unit>(`/units/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.units.all() })
      toast.success('Satuan berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteUnitMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/units/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.units.all() })
      toast.success('Satuan berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleUnitStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/units/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.units.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
