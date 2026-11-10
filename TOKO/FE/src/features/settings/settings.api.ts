import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

export function usePageSizeOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.settings.pageSizeOptions(),
    queryFn: async () => {
      const res = await api.get<{ key: string; value: string }>('/settings/pagination_sizes')
      return (JSON.parse(res.value) as number[])
    },
    staleTime: Infinity,
  })
}
