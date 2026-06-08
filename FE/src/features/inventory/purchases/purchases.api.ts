import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services/api.client'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateSupplierPurchasePayload,
  PurchasePayment,
  SupplierPurchase,
  SupplierPurchaseFilter,
  SupplierPurchasePayment,
} from './purchases.types'

const QK = {
  all: () => ['supplierPurchases'] as const,
  list: (filter?: SupplierPurchaseFilter) => ['supplierPurchases', 'list', filter] as const,
  detail: (id: number) => ['supplierPurchases', 'detail', id] as const,
}

export function useSupplierPurchasesQuery(filter?: SupplierPurchaseFilter) {
  return useQuery({
    queryKey: QK.list(filter),
    queryFn: () => api.post<PaginatedData<SupplierPurchase>>('/supplier-purchases/list', filter ?? {}),
  })
}

export function useSupplierPurchaseDetailQuery(id: number | null) {
  return useQuery({
    queryKey: QK.detail(id ?? 0),
    queryFn: async () => api.get<SupplierPurchase>(`/supplier-purchases/${id}`),
    enabled: id !== null && id > 0,
  })
}

export function useGeneratePurchaseCodeQuery(enabled: boolean) {
  return useQuery({
    queryKey: ['purchaseGenerateCode', enabled],
    queryFn: () => api.get<{ purchase_code: string }>('/supplier-purchases/generate-code'),
    enabled,
    staleTime: 0,
  })
}

export function useCreateSupplierPurchaseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateSupplierPurchasePayload) =>
      api.post<SupplierPurchase>('/supplier-purchases', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteSupplierPurchaseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete<void>(`/supplier-purchases/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useSupplierPurchasePaymentsQuery(id: number | null) {
  return useQuery({
    queryKey: [...QK.detail(id ?? 0), 'payments'],
    queryFn: () => api.get<PurchasePayment[]>(`/supplier-purchases/${id}/payments`),
    enabled: id !== null && id > 0,
  })
}

export function usePaySupplierPurchaseMutation(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: SupplierPurchasePayment) =>
      api.post<void>(`/supplier-purchases/${id}/pay`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
