import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  AppVersion,
  AppUser,
  ChangePasswordPayload,
  CreateAppVersionPayload,
  CreateUserPayload,
  PrinterSettings,
  StoreProfile,
  UpdateUserPayload,
} from './settings.types'

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

export function useUserListQuery() {
  return useQuery({
    queryKey: queryKeys.settings.users(),
    queryFn: () => api.post<AppUser[]>('/users/list', {}),
  })
}

export function useCreateUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => api.post<AppUser>('/users/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.users() })
      toast.success('User berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateUserPayload }) =>
      api.post<AppUser>(`/users/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.users() })
      toast.success('User berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useChangePasswordMutation() {
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: ChangePasswordPayload }) =>
      api.post<void>(`/users/update/${id}`, payload),
    onSuccess: () => toast.success('Password berhasil diubah'),
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/users/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.users() })
      toast.success('User berhasil dinonaktifkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

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

export function usePrinterSettingsQuery() {
  return useQuery({
    queryKey: queryKeys.settings.printer(),
    queryFn: () => api.get<PrinterSettings>('/settings/printer'),
  })
}

export function useUpdatePrinterSettingsMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: PrinterSettings) => api.post<PrinterSettings>('/settings/printer', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.printer() })
      toast.success('Pengaturan printer berhasil disimpan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
