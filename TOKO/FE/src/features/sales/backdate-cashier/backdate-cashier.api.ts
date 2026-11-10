import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { BackdatePaymentPayload, BackdateCheckoutResponse } from './backdate-cashier.types'

export function useBackdateCheckoutMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: BackdatePaymentPayload) =>
      api.post<BackdateCheckoutResponse>('/backdate/transactions/create', payload),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: queryKeys.transactions.all() })
      qc.invalidateQueries({ queryKey: queryKeys.backdateCash.all() })
      qc.invalidateQueries({ queryKey: queryKeys.backdateTransaction.all() })
      toast.success(`Transaksi historis ${data.transaction_code} berhasil`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// Reuse product search from regular cashier
export { useProductSearchQuery, useProductBarcodeSearchQuery, useCustomerCreditQuery } from '../cashier/cashier.api'
