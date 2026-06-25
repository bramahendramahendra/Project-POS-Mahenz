import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateSupplierPurchasePayload,
  PurchasePayment,
  SupplierPurchase,
  SupplierPurchaseFilter,
  SupplierPurchasePayment,
  UpdateSupplierPurchasePayload,
} from './purchases.types'

export function useGeneratePurchaseCodeQuery(enabled: boolean) {
  return useQuery({
    queryKey: ['purchaseGenerateCode'],
    queryFn: () => api.post<{ purchase_code: string }>('/supplier-purchases/generate-code', {}),
    enabled,
    staleTime: 0,
  })
}

export function useSupplierPurchasesQuery(filter?: SupplierPurchaseFilter) {
  return useQuery({
    queryKey: queryKeys.supplierPurchases.list(filter as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<SupplierPurchase>>('/supplier-purchases/list', filter ?? {}),
  })
}

export function useSupplierPurchaseDetailQuery(id: number | null) {
  return useQuery({
    queryKey: queryKeys.supplierPurchases.detail(id ?? 0),
    queryFn: () => api.post<SupplierPurchase>(`/supplier-purchases/detail/${id}`, {}),
    enabled: id !== null && id > 0,
  })
}

export function useSupplierPurchasePaymentsQuery(id: number | null) {
  return useQuery({
    queryKey: queryKeys.supplierPurchases.payments(id ?? 0),
    queryFn: () => api.post<PurchasePayment[]>(`/supplier-purchases/detail/${id}/payments`, {}),
    enabled: id !== null && id > 0,
  })
}

export function useCreateSupplierPurchaseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateSupplierPurchasePayload) =>
      api.post<SupplierPurchase>('/supplier-purchases/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierPurchases.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateSupplierPurchaseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateSupplierPurchasePayload) =>
      api.post<SupplierPurchase>(`/supplier-purchases/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierPurchases.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteSupplierPurchaseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/supplier-purchases/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierPurchases.all() })
      toast.success('Pembelian berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function usePaySupplierPurchaseMutation(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: SupplierPurchasePayment) =>
      api.post<void>(`/supplier-purchases/pay/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.supplierPurchases.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
