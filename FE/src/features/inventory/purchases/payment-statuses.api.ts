import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'

import type { PaymentStatusOption } from './payment-statuses.types'

export function usePaymentStatusesQuery() {
  return useQuery({
    queryKey: ['paymentStatuses'],
    queryFn: () => api.post<PaymentStatusOption[]>('/payment-statuses/list', {}),
    staleTime: 5 * 60 * 1000,
  })
}
