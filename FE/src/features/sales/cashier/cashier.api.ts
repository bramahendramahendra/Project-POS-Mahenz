import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { Product } from '@/features/products/products'

import type { CheckoutResponse, PaymentPayload, ProductSearchResult } from './cashier.types'

export function useProductSearchQuery(keyword: string, enabled: boolean) {
  return useQuery({
    queryKey: ['cashier', 'search', keyword],
    queryFn: () => api.post<ProductSearchResult[]>('/products/search', { q: keyword, limit: 20 }),
    enabled: enabled && keyword.length >= 2,
  })
}

export function useProductBarcodeSearchQuery(code: string, enabled: boolean) {
  return useQuery({
    queryKey: queryKeys.products.barcode(code),
    queryFn: () => api.post<Product>(`/products/by-barcode/${code}`, {}),
    enabled: enabled && code.length > 0,
  })
}

export { useCustomerListQuery } from '@/features/customers'
export { useActiveShiftQuery } from '@/features/operational/shifts'

export function useCustomerCreditQuery(customerId: number | null) {
  return useQuery({
    queryKey: ['customers', 'credit', customerId],
    queryFn: () =>
      api.post<{ id: number; name: string; credit_limit: number; outstanding_amount: number }>(
        `/customers/detail/${customerId}`,
        {}
      ),
    enabled: customerId !== null && customerId > 0,
  })
}

export function useCheckoutMutation() {
  const qc = useQueryClient()

  return useMutation({
    mutationFn: (payload: PaymentPayload) =>
      api.post<CheckoutResponse>('/transactions', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.transactions.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
