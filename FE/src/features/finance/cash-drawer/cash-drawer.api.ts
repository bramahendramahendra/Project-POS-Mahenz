import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CashDrawer,
  CashDrawerDetail,
  CashDrawerListFilter,
  CashDrawerSummary,
  CloseCashDrawerBody,
  CurrentCashDrawer,
  KasirOption,
  OpenCashDrawerPayload,
} from './cash-drawer.types'

export function useCashDrawerCurrentQuery() {
  return useQuery({
    queryKey: queryKeys.cashDrawer.current(),
    queryFn: () => api.post<CurrentCashDrawer | null>('/cash-drawer/current', {}),
  })
}

export function useCashDrawerListQuery(filter?: CashDrawerListFilter) {
  return useQuery({
    queryKey: queryKeys.cashDrawer.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<CashDrawer>>('/cash-drawer/list', filter ?? {}),
  })
}

export function useCashDrawerDetailQuery(id: number | null) {
  return useQuery({
    queryKey: queryKeys.cashDrawer.detail(id ?? 0),
    queryFn: () => api.post<CashDrawerDetail>(`/cash-drawer/detail/${id}`, {}),
    enabled: id !== null,
  })
}

export function useOpenCashDrawerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: OpenCashDrawerPayload) => api.post<void>('/cash-drawer/open', body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.cashDrawer.all() })
      qc.invalidateQueries({ queryKey: queryKeys.myCash.data() })
      toast.success('Kas berhasil dibuka')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useCloseCashDrawerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, closing_balance, notes }: CloseCashDrawerBody & { id: number }) =>
      api.post<void>(`/cash-drawer/close/${id}`, { closing_balance, notes }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.cashDrawer.all() })
      qc.invalidateQueries({ queryKey: queryKeys.myCash.data() })
      toast.success('Kas berhasil ditutup')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useCashDrawerSummaryQuery(filter?: CashDrawerListFilter) {
  return useQuery({
    queryKey: queryKeys.cashDrawer.summary(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<CashDrawerSummary>('/cash-drawer/summary', filter ?? {}),
  })
}

export function useKasirOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.cashDrawer.kasirOptions(),
    queryFn: () => api.post<KasirOption[]>('/cash-drawer/kasir-options', {}),
  })
}
