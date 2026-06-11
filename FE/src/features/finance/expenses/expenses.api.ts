import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { CreateExpensePayload, Expense, ExpenseListFilter, UpdateExpensePayload } from './expenses.types'

export function useExpensesQuery(filter?: ExpenseListFilter) {
  return useQuery({
    queryKey: queryKeys.expenses.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Expense>>('/expenses/list', filter ?? {}),
  })
}

export function useCreateExpenseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateExpensePayload) => api.post<Expense>('/expenses/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.expenses.all() })
      toast.success('Pengeluaran berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateExpenseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateExpensePayload & { id: number }) =>
      api.post<Expense>(`/expenses/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.expenses.all() })
      toast.success('Pengeluaran berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteExpenseMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/expenses/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.expenses.all() })
      toast.success('Pengeluaran berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
