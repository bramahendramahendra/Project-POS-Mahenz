import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateSupplierReturnPayload,
  SupplierReturn,
  SupplierReturnFilter,
  UpdateReturnStatusPayload,
} from './returns.types'

export function useSupplierReturnsQuery(filter?: SupplierReturnFilter) {
  return useQuery({
    queryKey: queryKeys.supplierReturns.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<SupplierReturn>>('/supplier-returns/list', filter ?? {}),
  })
}

export function useSupplierReturnDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.supplierReturns.detail(id),
    queryFn: () => api.post<SupplierReturn>(`/supplier-returns/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useCreateSupplierReturnMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateSupplierReturnPayload) =>
      api.post<SupplierReturn>('/supplier-returns/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierReturns.all() })
      toast.success('Retur berhasil dibuat')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateSupplierReturnStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateReturnStatusPayload & { id: number }) =>
      api.post<SupplierReturn>(`/supplier-returns/update-status/${id}`, payload),
    onSuccess: (_data, { status }) => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierReturns.all() })
      toast.success(status === 'approved' ? 'Retur berhasil disetujui' : 'Retur ditolak')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteSupplierReturnMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/supplier-returns/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierReturns.all() })
      toast.success('Retur berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
