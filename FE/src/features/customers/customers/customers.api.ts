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
      toast.success('Pelanggan berhasil ditambahkan')
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
      toast.success('Pelanggan berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleCustomerStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id }: { id: number; isActive: boolean }) =>
      api.post<void>(`/customers/toggle-status/${id}`, {}),
    onSuccess: (_data, { isActive }) => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
      toast.success(`Pelanggan berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`)
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

// ═══ Customer Balance API ═══

export interface BalanceMutation {
  id: number
  customer_id: number
  amount: number
  balance_after: number
  type: 'topup' | 'usage' | 'refund' | 'adjustment'
  reference_type?: string
  reference_id?: number
  notes?: string
  user_name: string
  created_at: string
}

export interface BalanceResponse {
  customer_id: number
  balance: number
  balance_after: number
}

export function useCustomerBalanceHistoryQuery(customerId: number | null, page = 1, limit = 10) {
  return useQuery({
    queryKey: [...queryKeys.customers.all(), 'balance-history', customerId, page],
    queryFn: () =>
      api.post<PaginatedData<BalanceMutation>>(`/customers/${customerId}/balance/history`, { page, limit }),
    enabled: customerId !== null && customerId > 0,
  })
}

export function useTopupBalanceMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, amount, notes }: { id: number; amount: number; notes?: string }) =>
      api.post<BalanceResponse>(`/customers/${id}/balance/topup`, { amount, notes }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
      toast.success('Top-up saldo berhasil')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useRefundBalanceMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, amount, notes }: { id: number; amount: number; notes?: string }) =>
      api.post<BalanceResponse>(`/customers/${id}/balance/refund`, { amount, notes }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.customers.all() })
      toast.success('Tarik saldo berhasil')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
