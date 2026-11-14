import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { Transaction } from '@/features/sales/transactions/transactions.types'

import { CATALOG_QUERY_KEY } from './cashier.api'
import type { OrderType } from './cashier.types'

// Open bill Coffee (tahap 4.8, keputusan 27-29, 103-109).

export interface OpenBillSummary {
  id: number
  transaction_code: string
  order_type: OrderType | ''
  table_id: number | null
  table_name: string
  queue_number: number
  queue_label: string
  order_name: string
  subtotal: number
  item_count: number
  created_by_name: string
  created_at: string
  /** Tanggal pesanan dan tanda tertunda (dari hari sebelumnya, tahap 4.9). */
  bill_date: string
  is_overdue: boolean
}

export interface TableOption {
  id: number
  name: string
  open_bill_id: number | null
  open_queue_label: string
}

export const OPEN_BILLS_KEY = ['cashier', 'open-bills'] as const

export function useOpenBillsQuery() {
  return useQuery({
    queryKey: OPEN_BILLS_KEY,
    queryFn: () => api.post<OpenBillSummary[]>('/transactions/open/list', {}),
    refetchInterval: 30_000,
  })
}

/** Pilihan meja kasir beserta tanda meja terisi. */
export function useTableOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.diningTables.options(),
    queryFn: () => api.post<TableOption[]>('/dining-tables/options', {}),
    refetchInterval: 30_000,
  })
}

export function useBillQuery(id: number | null) {
  return useQuery({
    queryKey: queryKeys.transactions.detail(id ?? 0),
    queryFn: () => api.post<Transaction>(`/transactions/detail/${id}`, {}),
    enabled: !!id,
  })
}

/** Setelah perubahan pesanan terbuka: daftar, detail, meja, stok katalog, riwayat ikut dimuat ulang. */
function useInvalidateBills() {
  const qc = useQueryClient()
  return () => {
    qc.invalidateQueries({ queryKey: OPEN_BILLS_KEY })
    qc.invalidateQueries({ queryKey: queryKeys.transactions.all() })
    qc.invalidateQueries({ queryKey: queryKeys.diningTables.all() })
    qc.invalidateQueries({ queryKey: CATALOG_QUERY_KEY })
  }
}

function useBillMutation<V>(fn: (v: V) => Promise<Transaction | void>, success: string) {
  const invalidate = useInvalidateBills()
  return useMutation({
    mutationFn: fn,
    onSuccess: () => {
      invalidate()
      toast.success(success)
    },
    onError: (e: Error) => {
      invalidate()
      toast.error(e.message)
    },
  })
}

export interface SaveOpenBillPayload {
  order_type: OrderType
  table_id?: number
  order_name?: string
  device_source: 'web'
  items: { variant_id: number; quantity: number; modifier_option_ids: number[]; notes?: string; discount_item?: number }[]
}

export function useSaveOpenBillMutation() {
  const invalidate = useInvalidateBills()
  return useMutation({
    mutationFn: (p: SaveOpenBillPayload) => api.post<Transaction>('/transactions/create', { ...p, open_bill: true }),
    onSuccess: (data) => {
      invalidate()
      toast.success(`Pesanan terbuka disimpan, antrean ${data.queue_label}`)
    },
    onError: (e: Error) => {
      invalidate()
      toast.error(e.message)
    },
  })
}

export const useAddBillItemsMutation = () =>
  useBillMutation(
    ({ id, items }: { id: number; items: SaveOpenBillPayload['items'] }) => api.post<Transaction>(`/transactions/open/${id}/add-items`, { items }),
    'Item ditambahkan ke pesanan',
  )

export const useRemoveBillItemMutation = () =>
  useBillMutation(
    ({ id, ...body }: { id: number; item_id: number; quantity: number; action?: 'return' | 'waste'; reason?: string }) =>
      api.post<Transaction>(`/transactions/open/${id}/remove-item`, body),
    'Item dihapus dari pesanan',
  )

export const useMoveBillTableMutation = () =>
  useBillMutation(({ id, table_id }: { id: number; table_id: number }) => api.post<Transaction>(`/transactions/open/${id}/move-table`, { table_id }), 'Pesanan dipindah meja')

export const useMergeBillMutation = () =>
  useBillMutation(({ id, source_id }: { id: number; source_id: number }) => api.post<Transaction>(`/transactions/open/${id}/merge`, { source_id }), 'Pesanan digabung')

export const useCancelBillMutation = () =>
  useBillMutation(
    ({ id, ...body }: { id: number; reason: string; items: { item_id: number; action: 'return' | 'waste' }[] }) => api.post<void>(`/transactions/open/${id}/cancel`, body),
    'Pesanan dibatalkan',
  )

export interface PayBillPayload {
  id: number
  shift_id?: number
  discount: number
  tax: number
  payment_method: string
  payment_amount: number
  customer_id?: number
  is_credit: boolean
}

export function usePayBillMutation() {
  const invalidate = useInvalidateBills()
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...body }: PayBillPayload) => api.post<Transaction>(`/transactions/open/${id}/pay`, body),
    onSuccess: (data) => {
      invalidate()
      qc.invalidateQueries({ queryKey: queryKeys.cashDrawer.all() })
      qc.invalidateQueries({ queryKey: queryKeys.myCash.data() })
      toast.success(`Pesanan antrean ${data.queue_label} ${data.is_credit ? 'dicatat sebagai utang' : 'dibayar'}`)
    },
    onError: (e: Error) => {
      invalidate()
      toast.error(e.message)
    },
  })
}
