import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  AdminApp,
  CreateAppPayload,
  CreateAppResult,
  RotateKeyResult,
  UpdateAppPayload,
} from './apps.types'

export function useAdminAppsQuery() {
  return useQuery({
    queryKey: queryKeys.apps.adminList(),
    queryFn: () => api.post<AdminApp[]>('/apps/list', {}),
  })
}

// Daftar admin dan daftar kartu di beranda sama-sama berubah saat aplikasi diubah.
function useInvalidateApps() {
  const queryClient = useQueryClient()
  return () => {
    void queryClient.invalidateQueries({ queryKey: queryKeys.apps.adminList() })
    void queryClient.invalidateQueries({ queryKey: queryKeys.apps.my() })
  }
}

export function useCreateAppMutation() {
  const invalidate = useInvalidateApps()
  return useMutation({
    mutationFn: (payload: CreateAppPayload) => api.post<CreateAppResult>('/apps/create', payload),
    onSuccess: invalidate,
  })
}

export function useUpdateAppMutation() {
  const invalidate = useInvalidateApps()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateAppPayload }) =>
      api.post<void>(`/apps/update/${id}`, payload),
    onSuccess: () => {
      invalidate()
      toast.success('Aplikasi diperbarui')
    },
  })
}

export function useToggleAppMutation() {
  const invalidate = useInvalidateApps()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/apps/toggle-status/${id}`, {}),
    onSuccess: () => {
      invalidate()
      toast.success('Status aplikasi diubah')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useRotateKeyMutation() {
  const invalidate = useInvalidateApps()
  return useMutation({
    mutationFn: (id: number) => api.post<RotateKeyResult>(`/apps/rotate-key/${id}`, {}),
    onSuccess: invalidate,
    onError: (e: Error) => toast.error(e.message),
  })
}
