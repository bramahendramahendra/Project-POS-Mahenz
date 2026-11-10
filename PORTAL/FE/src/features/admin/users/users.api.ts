import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  AdminUser,
  CreateUserPayload,
  RoleOption,
  UpdateUserPayload,
  UserAccessItem,
  UsersFilter,
} from './users.types'

export function useUsersQuery(filter: UsersFilter) {
  return useQuery({
    queryKey: queryKeys.users.list(filter),
    queryFn: () => api.post<PaginatedData<AdminUser>>('/users/list', filter),
    placeholderData: (previous) => previous,
  })
}

export function useRoleOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.roles.options(),
    queryFn: () => api.post<RoleOption[]>('/roles/options', {}),
    staleTime: 10 * 60 * 1000,
  })
}

function useInvalidateUsers() {
  const queryClient = useQueryClient()
  return () => queryClient.invalidateQueries({ queryKey: queryKeys.users.all() })
}

export function useCreateUserMutation() {
  const invalidate = useInvalidateUsers()
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => api.post<AdminUser>('/users/create', payload),
    onSuccess: () => {
      void invalidate()
      toast.success('User ditambahkan. User wajib mengganti password saat login pertama')
    },
  })
}

export function useUpdateUserMutation() {
  const invalidate = useInvalidateUsers()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateUserPayload }) =>
      api.post<void>(`/users/update/${id}`, payload),
    onSuccess: () => {
      void invalidate()
      toast.success('User diperbarui')
    },
  })
}

export function useToggleUserMutation() {
  const invalidate = useInvalidateUsers()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/users/toggle-status/${id}`, {}),
    onSuccess: () => {
      void invalidate()
      toast.success('Status user diubah')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteUserMutation() {
  const invalidate = useInvalidateUsers()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/users/delete/${id}`, {}),
    onSuccess: () => {
      void invalidate()
      toast.success('User dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useResetPasswordMutation() {
  const invalidate = useInvalidateUsers()
  return useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) =>
      api.post<void>(`/users/reset-password/${id}`, { password }),
    onSuccess: () => {
      void invalidate()
      toast.success('Password sementara ditetapkan')
    },
  })
}

export function useUserAccessQuery(userId: number | null) {
  return useQuery({
    queryKey: queryKeys.apps.userAccess(userId ?? 0),
    queryFn: () => api.post<UserAccessItem[]>(`/apps/access/${userId}`, {}),
    enabled: userId !== null,
    staleTime: 0,
  })
}

export function useSetUserAccessMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ userId, appIds }: { userId: number; appIds: number[] }) =>
      api.post<void>(`/apps/access/${userId}/set`, { app_ids: appIds }),
    onSuccess: (_data, vars) => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.apps.userAccess(vars.userId) })
      void queryClient.invalidateQueries({ queryKey: queryKeys.apps.adminList() })
      toast.success('Akses aplikasi diperbarui')
    },
  })
}
