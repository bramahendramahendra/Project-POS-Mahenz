import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { StoreProfile } from './store.types'

export function useStoreProfileQuery() {
  return useQuery({
    queryKey: queryKeys.settings.store(),
    queryFn: () => api.get<StoreProfile>('/settings/store'),
  })
}

export function useUpdateStoreProfileMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: StoreProfile) => api.post<StoreProfile>('/settings/store', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.store() })
      toast.success('Profil toko berhasil disimpan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
