import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  MasterPayload,
  MenuCategory,
  MenuItem,
  MenuItemDetail,
  MenuItemListFilter,
  MenuItemPayload,
  MenuSize,
  DiningTable,
  ModifierGroup,
  ModifierGroupPayload,
  NamedOption,
} from './menu-items.types'

// ─── Master sederhana: kategori menu dan ukuran (bentuk API sama) ─────────────

interface MasterKeys {
  all: () => readonly unknown[]
  list: (filter?: object) => readonly unknown[]
  options: () => readonly unknown[]
}

function createMasterApi<T>(resource: string, keys: MasterKeys, label: string) {
  return {
    useList(filter: { search: string; is_active?: boolean }) {
      return useQuery({ queryKey: keys.list(filter), queryFn: () => api.post<T[]>(`/${resource}/list`, filter) })
    },
    useOptions() {
      return useQuery({ queryKey: keys.options(), queryFn: () => api.post<NamedOption[]>(`/${resource}/options`, {}) })
    },
    useSave() {
      const qc = useQueryClient()
      return useMutation({
        mutationFn: ({ id, ...payload }: MasterPayload & { id?: number }) =>
          id ? api.post<void>(`/${resource}/update/${id}`, payload) : api.post<void>(`/${resource}/create`, payload),
        onSuccess: (_d, vars) => {
          qc.invalidateQueries({ queryKey: keys.all() })
          toast.success(`${label} berhasil ${vars.id ? 'diperbarui' : 'ditambahkan'}`)
        },
        onError: (e: Error) => toast.error(e.message),
      })
    },
    useToggle() {
      const qc = useQueryClient()
      return useMutation({
        mutationFn: (id: number) => api.post<void>(`/${resource}/toggle-status/${id}`, {}),
        onSuccess: () => {
          qc.invalidateQueries({ queryKey: keys.all() })
          toast.success(`Status ${label.toLowerCase()} berhasil diubah`)
        },
        onError: (e: Error) => toast.error(e.message),
      })
    },
    useDelete() {
      const qc = useQueryClient()
      return useMutation({
        mutationFn: (id: number) => api.post<void>(`/${resource}/delete/${id}`, {}),
        onSuccess: () => {
          qc.invalidateQueries({ queryKey: keys.all() })
          toast.success(`${label} berhasil dihapus`)
        },
        onError: (e: Error) => toast.error(e.message),
      })
    },
  }
}

export type MasterApi<T> = ReturnType<typeof createMasterApi<T>>

export const menuCategoryApi = createMasterApi<MenuCategory>('menu-categories', queryKeys.menuCategories, 'Kategori menu')
export const menuSizeApi = createMasterApi<MenuSize>('menu-sizes', queryKeys.menuSizes, 'Ukuran')
export const diningTableApi = createMasterApi<DiningTable>('dining-tables', queryKeys.diningTables, 'Meja')

// ─── Kelompok tambahan ────────────────────────────────────────────────────────

export function useModifierGroupListQuery(filter: { search: string; is_active?: boolean }) {
  return useQuery({
    queryKey: queryKeys.modifierGroups.list(filter),
    queryFn: () => api.post<ModifierGroup[]>('/modifier-groups/list', filter),
  })
}

export function useModifierGroupOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.modifierGroups.options(),
    queryFn: () => api.post<ModifierGroup[]>('/modifier-groups/options', {}),
  })
}

export function useSaveModifierGroupMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: ModifierGroupPayload & { id?: number }) =>
      id ? api.post<unknown>(`/modifier-groups/update/${id}`, payload) : api.post<unknown>('/modifier-groups/create', payload),
    onSuccess: (_d, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.modifierGroups.all() })
      qc.invalidateQueries({ queryKey: queryKeys.menuItems.all() })
      toast.success(`Kelompok tambahan berhasil ${vars.id ? 'diperbarui' : 'ditambahkan'}`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleModifierGroupMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/modifier-groups/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.modifierGroups.all() })
      toast.success('Status kelompok tambahan berhasil diubah')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteModifierGroupMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/modifier-groups/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.modifierGroups.all() })
      toast.success('Kelompok tambahan berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Menu ─────────────────────────────────────────────────────────────────────

export function useMenuItemListQuery(filter: MenuItemListFilter) {
  return useQuery({
    queryKey: queryKeys.menuItems.list(filter),
    queryFn: () => api.post<PaginatedData<MenuItem>>('/menu-items/list', filter),
  })
}

export function useMenuItemDetailQuery(id: number | null) {
  return useQuery({
    queryKey: queryKeys.menuItems.detail(id ?? 0),
    queryFn: () => api.post<MenuItemDetail>(`/menu-items/detail/${id}`, {}),
    enabled: id != null && id > 0,
  })
}

export function useSaveMenuItemMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: MenuItemPayload & { id?: number }) =>
      id ? api.post<unknown>(`/menu-items/update/${id}`, payload) : api.post<unknown>('/menu-items/create', payload),
    onSuccess: (_d, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.menuItems.all() })
      qc.invalidateQueries({ queryKey: queryKeys.menuCategories.all() })
      qc.invalidateQueries({ queryKey: queryKeys.menuSizes.all() })
      qc.invalidateQueries({ queryKey: queryKeys.modifierGroups.all() })
      toast.success(`Menu berhasil ${vars.id ? 'diperbarui' : 'ditambahkan'}`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleMenuItemMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/menu-items/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menuItems.all() })
      toast.success('Status menu berhasil diubah')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useSetSoldOutMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, isSoldOut }: { id: number; isSoldOut: boolean }) =>
      api.post<void>(`/menu-items/sold-out/${id}`, { is_sold_out: isSoldOut }),
    onSuccess: (_d, { isSoldOut }) => {
      qc.invalidateQueries({ queryKey: queryKeys.menuItems.all() })
      toast.success(isSoldOut ? 'Menu ditandai habis' : 'Menu ditandai tersedia')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteMenuItemMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/menu-items/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.menuItems.all() })
      qc.invalidateQueries({ queryKey: queryKeys.menuCategories.all() })
      toast.success('Menu berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
