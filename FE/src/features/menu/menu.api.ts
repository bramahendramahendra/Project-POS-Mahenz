import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  CreateMenuPayload,
  MenuFilter,
  MenuResponse,
  MenuItem,
  ReorderMenuPayload,
  UpdateMenuPayload,
} from './menu.types'

export function useMyMenusQuery() {
  return useQuery({
    queryKey: queryKeys.menus.my(),
    queryFn: () => api.post<MenuItem[]>('/menus/my', {}),
    staleTime: 10 * 60 * 1000,
  })
}

export function useMenuListQuery(filter?: MenuFilter) {
  return useQuery({
    queryKey: queryKeys.menus.list(filter),
    queryFn: () => api.post<MenuResponse[]>('/menus/list', filter ?? {}),
  })
}

export function useMenuDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.menus.detail(id),
    queryFn: () => api.post<MenuResponse>(`/menus/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useCreateMenuMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateMenuPayload) => api.post<MenuResponse>('/menus/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menus.all() })
      toast.success('Menu berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateMenuMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateMenuPayload & { id: number }) =>
      api.post<MenuResponse>(`/menus/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menus.all() })
      toast.success('Menu berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteMenuMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post(`/menus/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menus.all() })
      toast.success('Menu berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useReorderMenuMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: ReorderMenuPayload) => api.post('/menus/reorder', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menus.all() })
      toast.success('Urutan menu berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
