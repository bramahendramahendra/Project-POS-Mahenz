import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreateRolePayload,
  Role,
  RoleListFilter,
  RoleMenuAccessItem,
  RoleOption,
  SetRoleAccessPayload,
  UpdateRolePayload,
} from './roles.types'

export function useRoleListQuery(filter: RoleListFilter) {
  return useQuery({
    queryKey: queryKeys.roles.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Role>>('/roles/list', filter),
  })
}

export function useRoleOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.roles.options(),
    queryFn: () => api.post<RoleOption[]>('/roles/options', {}),
  })
}

export function useRoleDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.roles.detail(id),
    queryFn: () => api.post<Role>(`/roles/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useCreateRoleMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateRolePayload) => api.post<Role>('/roles/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.roles.all() })
      toast.success('Role berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateRoleMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateRolePayload & { id: number }) =>
      api.post<Role>(`/roles/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.roles.all() })
      toast.success('Role berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteRoleMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post(`/roles/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.roles.all() })
      toast.success('Role berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useRoleMenuAccessQuery(roleId: number) {
  return useQuery({
    queryKey: queryKeys.roles.menus(roleId),
    queryFn: () => api.post<RoleMenuAccessItem[]>(`/roles/${roleId}/menus/list`, {}),
    enabled: roleId > 0,
  })
}

export function useSetRoleAccessMutation(roleId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: SetRoleAccessPayload) =>
      api.post(`/roles/${roleId}/menus/set`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.roles.menus(roleId) })
      toast.success('Akses menu role berhasil disimpan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleRoleStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post(`/roles/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.roles.all() })
      toast.success('Status role berhasil diubah')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
