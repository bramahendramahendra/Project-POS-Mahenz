import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  Ingredient,
  IngredientListFilter,
  IngredientMutation,
  IngredientOption,
  IngredientPayload,
  MutationListFilter,
  RecordResult,
  RecordUsagePayload,
  RecordWastePayload,
} from './ingredients.types'

export function useIngredientListQuery(filter: IngredientListFilter) {
  return useQuery({
    queryKey: queryKeys.ingredients.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Ingredient>>('/ingredients/list', filter),
  })
}

export function useIngredientOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.ingredients.options(),
    queryFn: () => api.post<IngredientOption[]>('/ingredients/options', {}),
  })
}

export function useCreateIngredientMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: IngredientPayload) => api.post<Ingredient>('/ingredients/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      toast.success('Bahan berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateIngredientMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: IngredientPayload & { id: number }) =>
      api.post<void>(`/ingredients/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      toast.success('Bahan berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleIngredientStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id }: { id: number; willActivate: boolean }) =>
      api.post<void>(`/ingredients/toggle-status/${id}`, {}),
    onSuccess: (_data, { willActivate }) => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      toast.success(`Bahan berhasil ${willActivate ? 'diaktifkan' : 'dinonaktifkan'}`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteIngredientMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/ingredients/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      toast.success('Bahan berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// Stok minus tetap dicatat (keputusan 14.9), tetapi pengguna diberi peringatan.
function notifyRecorded(label: string, result: RecordResult) {
  if (result.negative_stock) {
    toast.warning(`${label} dicatat. Peringatan: stok ${result.mutation.ingredient_name} menjadi minus`)
  } else {
    toast.success(`${label} dicatat`)
  }
}

export function useRecordUsageMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: RecordUsagePayload) => api.post<RecordResult>('/ingredients/usage', payload),
    onSuccess: (result) => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      notifyRecorded('Pemakaian', result)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useRecordWasteMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: RecordWastePayload) => api.post<RecordResult>('/ingredients/waste', payload),
    onSuccess: (result) => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      notifyRecorded('Barang terbuang', result)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useMutationListQuery(filter: MutationListFilter) {
  return useQuery({
    queryKey: queryKeys.ingredients.mutations(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<IngredientMutation>>('/ingredients/mutations/list', filter),
  })
}

export function useCancelMutationMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      api.post<IngredientMutation>(`/ingredients/mutations/cancel/${id}`, { reason }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.ingredients.all() })
      toast.success('Entri berhasil dibatalkan, stok dikembalikan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
