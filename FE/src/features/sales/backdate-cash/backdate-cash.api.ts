import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  BackdateCashDrawer,
  CloseBackdateCashDrawerPayload,
  CloseBackdateCashDrawerResult,
  OpenBackdateCashDrawerPayload,
  UserOption,
} from './backdate-cash.types'

export function useBackdateCashDrawerCurrentQuery() {
  return useQuery({
    queryKey: queryKeys.backdateCash.current(),
    queryFn: () => api.post<BackdateCashDrawer | null>('/backdate/cash-drawer/current', {}),
  })
}

export function useOpenBackdateCashDrawerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: OpenBackdateCashDrawerPayload) =>
      api.post<{ id: number }>('/backdate/cash-drawer/open', body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.backdateCash.all() })
      toast.success('Kas historis berhasil dibuka')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useCloseBackdateCashDrawerMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...body }: CloseBackdateCashDrawerPayload & { id: number }) =>
      api.post<CloseBackdateCashDrawerResult>(`/backdate/cash-drawer/close/${id}`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.backdateCash.all() })
      toast.success('Kas historis berhasil ditutup')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUserOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.backdateCash.userOptions(),
    queryFn: () =>
      api.post<PaginatedData<UserOption>>('/users/list', { page: 1, limit: 100 }).then((res) => res.data),
  })
}
