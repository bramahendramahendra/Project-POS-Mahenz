import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { PaymentMethodOption } from './payment-methods.types'

export function usePaymentMethodsQuery() {
  return useQuery({
    queryKey: queryKeys.paymentMethods.options(),
    queryFn: () => api.post<PaymentMethodOption[]>('/payment-methods/list', {}),
    staleTime: 5 * 60 * 1000,
  })
}
