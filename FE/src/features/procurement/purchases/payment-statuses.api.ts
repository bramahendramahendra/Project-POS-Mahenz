import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { PaymentStatusOption } from './payment-statuses.types'

export function usePaymentStatusesQuery() {
  return useQuery({
    queryKey: queryKeys.paymentStatuses.options(),
    queryFn: () => api.post<PaymentStatusOption[]>('/payment-statuses/list', {}),
    staleTime: 5 * 60 * 1000,
  })
}
