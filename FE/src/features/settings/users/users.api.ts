import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  AppUser,
  ChangePasswordPayload,
  CreateUserPayload,
  UpdateUserPayload,
  UserListFilter,
} from './users.types'

export function useUserListQuery(filter: UserListFilter) {
  return useQuery({
    queryKey: queryKeys.users.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<AppUser>>('/users/list', filter),
  })
}

export function useCreateUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => api.post<AppUser>('/users/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.users.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateUserPayload & { id: number }) =>
      api.post<AppUser>(`/users/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.users.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useChangePasswordMutation() {
  return useMutation({
    mutationFn: ({ id, ...payload }: ChangePasswordPayload & { id: number }) =>
      api.post<void>(`/users/change-password/${id}`, payload),
    onSuccess: () => toast.success('Password berhasil diubah'),
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteUserMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/users/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.users.all() })
      toast.success('User berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleUserStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/users/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.users.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
