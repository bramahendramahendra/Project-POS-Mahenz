import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type { Transaction, TransactionListFilter } from './transactions.types'

export function useTransactionListQuery(filter: TransactionListFilter) {
  return useQuery({
    queryKey: queryKeys.transactions.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Transaction>>('/transactions/list', filter),
  })
}

export function useTransactionDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.transactions.detail(id),
    queryFn: () => api.post<Transaction>(`/transactions/detail/${id}`, {}),
    enabled: id > 0,
  })
}

export function useVoidTransactionMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/transactions/void/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.transactions.all() })
      toast.success('Transaksi berhasil dibatalkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}
