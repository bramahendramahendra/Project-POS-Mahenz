import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateCustomerPayload,
  Customer,
  CustomerListFilter,
  UpdateCustomerPayload,
} from './customers.types'

export function useCustomerListQuery(filter: CustomerListFilter) {
  return useQuery({
    queryKey: queryKeys.customers.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Customer>>('/customers/list', filter),
  })
}

export function useCreateCustomerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateCustomerPayload) =>
      api.post<Customer>('/customers/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateCustomerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateCustomerPayload & { id: number }) =>
      api.post<Customer>(`/customers/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleCustomerStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/customers/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteCustomerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/customers/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
      toast.success('Pelanggan berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
