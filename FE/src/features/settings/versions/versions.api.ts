import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { AppVersion, CreateAppVersionPayload } from './versions.types'

export function useAppVersionListQuery() {
  return useQuery({
    queryKey: queryKeys.settings.appVersions(),
    queryFn: () => api.post<AppVersion[]>('/version/list', {}),
  })
}

export function useCreateAppVersionMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateAppVersionPayload) => api.post('/version/android', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.appVersions() })
      toast.success('Versi berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
