import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { RouteOption } from './menu.types'

export function useRouteRegistryOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.routeRegistry.options(),
    queryFn: () => api.post<RouteOption[]>('/route-registry/options', {}),
    staleTime: 10 * 60 * 1000,
  })
}
