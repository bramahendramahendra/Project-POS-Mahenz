import { useMemo } from 'react'

import { transactionReceipt } from '@/features/sales/transactions/transactions.utils'

import { useCashierStore } from './cashier.store'
import type { CartItem } from './cashier.types'
import { calcCartSummary } from './cashier.utils'
import { useBillQuery } from './openBill.api'

/**
 * Pesanan terbuka yang sedang dibuka di kasir beserta ringkasan pembayarannya (tahap 4.8).
 * Mode pesanan baru: ringkasan dari keranjang. Mode pesanan terbuka: dari item tersimpan,
 * diskon nota dan pajak diisi saat bayar (keputusan 104).
 */
export function useActiveBill() {
  const { cart, discount, tax, activeBillId } = useCashierStore()
  const { data: bill, isLoading } = useBillQuery(activeBillId)
  const billOpen = !!activeBillId && bill?.status === 'open'
  const billLines = useMemo<CartItem[]>(() => (bill ? transactionReceipt(bill).cart : []), [bill])
  const payLines = billOpen ? billLines : cart
  const summary = calcCartSummary(payLines, discount, tax)
  return { activeBillId, bill, isLoading, billOpen, billLines, payLines, summary }
}

/** Store untuk DiscountInput/TaxInput dalam mode pesanan terbuka: "keranjang" = item tersimpan. */
export function makeBillStore(lines: CartItem[]) {
  return function useBillStore() {
    return { ...useCashierStore(), cart: lines }
  }
}
