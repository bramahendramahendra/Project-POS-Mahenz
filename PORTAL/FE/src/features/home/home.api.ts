import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { LaunchResponse, PortalApp } from './home.types'

export function useMyAppsQuery() {
  return useQuery({
    queryKey: queryKeys.apps.my(),
    queryFn: () => api.post<PortalApp[]>('/apps/my', {}),
    // Daftar akses bisa diubah admin kapan saja, jadi selalu segarkan saat halaman dibuka.
    staleTime: 0,
  })
}

// Meminta kode SSO ke Portal, lalu pindah ke aplikasi tujuan di tab yang sama.
export function useLaunchAppMutation() {
  return useMutation({
    mutationFn: (appCode: string) => api.post<LaunchResponse>('/sso/launch', { app_code: appCode }),
    onSuccess: (data) => {
      window.location.href = data.redirect_url
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Gagal membuka aplikasi')
    },
  })
}
