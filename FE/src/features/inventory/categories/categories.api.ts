import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  Category,
  CategoryListFilter,
  CategoryOption,
  CreateCategoryPayload,
  UpdateCategoryPayload,
} from './categories.types'

export function useCategoryListQuery(filter: CategoryListFilter) {
  return useQuery({
    queryKey: queryKeys.categories.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Category>>('/categories/list', filter),
  })
}

export function useCategoryOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.categories.options(),
    queryFn: () => api.post<CategoryOption[]>('/categories/options', {}),
  })
}

export function useCreateCategoryMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateCategoryPayload) =>
      api.post<Category>('/categories/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.categories.all() })
      toast.success('Kategori berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateCategoryMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateCategoryPayload & { id: number }) =>
      api.post<Category>(`/categories/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.categories.all() })
      toast.success('Kategori berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteCategoryMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/categories/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.categories.all() })
      toast.success('Kategori berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleCategoryStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/categories/toggle-status/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.categories.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
