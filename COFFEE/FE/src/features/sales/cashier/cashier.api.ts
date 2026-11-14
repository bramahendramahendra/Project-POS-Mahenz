import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { CheckoutResponse, MenuCatalog, PaymentPayload } from './cashier.types'

export const CATALOG_QUERY_KEY = ['cashier', 'catalog'] as const

/** Katalog kasir: menu aktif, kombinasi, tambahan, status habis, sisa stok barang jadi. */
export function useMenuCatalogQuery() {
  return useQuery({
    queryKey: CATALOG_QUERY_KEY,
    queryFn: () => api.post<MenuCatalog>('/menu-items/catalog', {}),
    // Stok barang jadi dan tanda habis bisa berubah dari kasir lain.
    refetchInterval: 60_000,
  })
}

export { useCustomerListQuery } from '@/features/customers'

export function useCustomerCreditQuery(customerId: number | null) {
  return useQuery({
    queryKey: ['customers', 'credit', customerId],
    queryFn: () =>
      api.post<{ id: number; name: string; credit_limit: number; balance: number; outstanding_amount: number }>(
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
      api.post<CheckoutResponse>('/transactions/create', payload),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: queryKeys.transactions.all() })
      qc.invalidateQueries({ queryKey: queryKeys.cashDrawer.all() })
      qc.invalidateQueries({ queryKey: queryKeys.myCash.data() })
      qc.invalidateQueries({ queryKey: CATALOG_QUERY_KEY })
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success(`Transaksi ${data.transaction_code} berhasil, antrean ${data.queue_label ?? ''}`)
    },
    onError: (e: Error) => {
      // Stok atau status menu mungkin berubah: muat ulang katalog agar tampilan sesuai.
      qc.invalidateQueries({ queryKey: CATALOG_QUERY_KEY })
      toast.error(e.message)
    },
  })
}
