import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateSupplierPayload,
  Supplier,
  SupplierDetail,
  SupplierListFilter,
  UpdateSupplierPayload,
} from './suppliers.types'

export function useSupplierListQuery(filter: SupplierListFilter) {
  return useQuery({
    queryKey: queryKeys.suppliers.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Supplier>>('/suppliers/list', filter),
  })
}

export function useSupplierOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.suppliers.options(),
    queryFn: () => api.post<Supplier[]>('/suppliers/options', {}),
  })
}

export function useSupplierDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.suppliers.detail(id),
    queryFn: () => api.post<SupplierDetail>(`/suppliers/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useCreateSupplierMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateSupplierPayload) => api.post<Supplier>('/suppliers/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.suppliers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateSupplierMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateSupplierPayload & { id: number }) =>
      api.post<Supplier>(`/suppliers/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.suppliers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteSupplierMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/suppliers/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.suppliers.all() })
      toast.success('Supplier berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleSupplierStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/suppliers/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.suppliers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
