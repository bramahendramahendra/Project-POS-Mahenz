import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { MyCashData } from './my-cash.types'

export function useMyCashQuery() {
  return useQuery({
    queryKey: queryKeys.myCash.data(),
    queryFn: () => api.post<MyCashData>('/cash-drawer/my-cash', {}),
  })
}
