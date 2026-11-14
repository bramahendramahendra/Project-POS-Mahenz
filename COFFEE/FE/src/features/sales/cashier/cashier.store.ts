import { create } from 'zustand'
import { persist } from 'zustand/middleware'

import type { CartItem, Discount, DiscountType, OrderType, Tax } from './cashier.types'
import { calcDiscountAmount, calcItemSubtotal, calcTaxAmount, calculateItemDiscount } from './cashier.utils'

const DEFAULT_DISCOUNT: Discount = { type: 'none', value: 0, amount: 0 }
const DEFAULT_TAX: Tax = { percent: 0, amount: 0 }

interface CashierState {
  // Data (disimpan sebagai draf di perangkat)
  cart: CartItem[]
  discount: Discount
  tax: Tax
  selectedCustomer: { id: number; name: string } | null
  // Info pesanan Coffee (keputusan 25, 95-99)
  orderType: OrderType
  table: { id: number; name: string } | null
  orderName: string

  // UI State (not persisted)
  paymentModalOpen: boolean

  // Actions — Cart (baris dikenali dari key: kombinasi + tambahan + catatan)
  addToCart: (item: CartItem) => void
  removeFromCart: (key: string) => void
  updateQty: (key: string, qty: number) => void
  setItemDiscount: (key: string, type: 'percent' | 'nominal', value: number) => void
  clearCart: () => void

  // Actions — Discount & Tax
  setDiscount: (discount: Omit<Discount, 'amount'>) => void
  setTax: (percent: number) => void

  // Actions — Customer & pesanan
  setCustomer: (customer: { id: number; name: string } | null) => void
  setOrderType: (type: OrderType) => void
  setTable: (table: { id: number; name: string } | null) => void
  setOrderName: (name: string) => void

  // Actions — Modal
  openPaymentModal: () => void
  closePaymentModal: () => void
}

function recalcDiscountAndTax(
  cart: CartItem[],
  discount: Discount,
  tax: Tax
): { discount: Discount; tax: Tax } {
  const subtotal = cart.reduce((s, i) => s + i.subtotal, 0)
  const discountAmount = calcDiscountAmount(subtotal, discount)
  const taxAmount = calcTaxAmount(subtotal, discountAmount, tax.percent)
  return {
    discount: { ...discount, amount: discountAmount },
    tax: { ...tax, amount: taxAmount },
  }
}

/** Hitung ulang subtotal satu baris setelah jumlah atau harga berubah (diskon item ikut). */
function withQty(item: CartItem, qty: number): CartItem {
  if (item.discount_type && item.discount_value) {
    const { discount_amount, effective_price, subtotal } = calculateItemDiscount(item.price, qty, item.discount_type, item.discount_value)
    return { ...item, qty, discount_amount, effective_price, subtotal }
  }
  return { ...item, qty, subtotal: calcItemSubtotal(qty, item.price) }
}

export const useCashierStore = create<CashierState>()(
  persist(
    (set, get) => ({
      cart: [],
      discount: DEFAULT_DISCOUNT,
      tax: DEFAULT_TAX,
      selectedCustomer: null,
      orderType: 'take_away',
      table: null,
      orderName: '',

      paymentModalOpen: false,

      // ── Cart Actions ──

      addToCart: (item) => {
        const { cart, discount, tax } = get()
        const existing = cart.find((i) => i.key === item.key)
        const newCart = existing
          ? cart.map((i) => (i.key === item.key ? withQty(i, i.qty + item.qty) : i))
          : [...cart, withQty(item, item.qty)]
        set({ cart: newCart, ...recalcDiscountAndTax(newCart, discount, tax) })
      },

      removeFromCart: (key) => {
        const { cart, discount, tax } = get()
        const newCart = cart.filter((i) => i.key !== key)
        set({ cart: newCart, ...recalcDiscountAndTax(newCart, discount, tax) })
      },

      updateQty: (key, qty) => {
        if (qty <= 0) {
          get().removeFromCart(key)
          return
        }
        const { cart, discount, tax } = get()
        const newCart = cart.map((i) => (i.key === key ? withQty(i, Math.floor(qty)) : i))
        set({ cart: newCart, ...recalcDiscountAndTax(newCart, discount, tax) })
      },

      setItemDiscount: (key, type, value) => {
        const { cart, discount, tax } = get()
        const newCart = cart.map((i) => {
          if (i.key !== key) return i
          if (value <= 0) {
            return {
              ...i,
              discount_type: undefined,
              discount_value: undefined,
              discount_amount: undefined,
              effective_price: undefined,
              subtotal: calcItemSubtotal(i.qty, i.price),
            }
          }
          const { discount_amount, effective_price, subtotal } = calculateItemDiscount(i.price, i.qty, type, value)
          return { ...i, discount_type: type, discount_value: value, discount_amount, effective_price, subtotal }
        })
        set({ cart: newCart, ...recalcDiscountAndTax(newCart, discount, tax) })
      },

      clearCart: () =>
        set({
          cart: [],
          discount: DEFAULT_DISCOUNT,
          tax: DEFAULT_TAX,
          selectedCustomer: null,
          orderType: 'take_away',
          table: null,
          orderName: '',
        }),

      // ── Discount & Tax ──

      setDiscount: ({ type, value }: { type: DiscountType; value: number }) => {
        const { cart, tax } = get()
        const subtotal = cart.reduce((s, i) => s + i.subtotal, 0)
        const amount = calcDiscountAmount(subtotal, { type, value })
        const taxAmount = calcTaxAmount(subtotal, amount, tax.percent)
        set({
          discount: { type, value, amount },
          tax: { ...tax, amount: taxAmount },
        })
      },

      setTax: (percent) => {
        const { cart, discount } = get()
        const subtotal = cart.reduce((s, i) => s + i.subtotal, 0)
        const taxAmount = calcTaxAmount(subtotal, discount.amount, percent)
        set({ tax: { percent, amount: taxAmount } })
      },

      // ── Customer & pesanan ──

      setCustomer: (customer) => set({ selectedCustomer: customer }),
      // Bawa pulang tidak memakai meja.
      setOrderType: (orderType) => set(orderType === 'take_away' ? { orderType, table: null } : { orderType }),
      setTable: (table) => set({ table }),
      setOrderName: (orderName) => set({ orderName }),

      // ── Modals ──

      openPaymentModal: () => set({ paymentModalOpen: true }),
      closePaymentModal: () => set({ paymentModalOpen: false }),
    }),
    {
      // Versi 2: bentuk keranjang Coffee (menu, kombinasi, tambahan). Draf versi lama (produk Toko) dibuang.
      name: 'coffee.cashier-draft',
      version: 2,
      migrate: () => ({
        cart: [],
        discount: DEFAULT_DISCOUNT,
        tax: DEFAULT_TAX,
        selectedCustomer: null,
        orderType: 'take_away' as OrderType,
        table: null,
        orderName: '',
      }),
      partialize: (state) => ({
        cart: state.cart,
        discount: state.discount,
        tax: state.tax,
        selectedCustomer: state.selectedCustomer,
        orderType: state.orderType,
        table: state.table,
        orderName: state.orderName,
      }),
    }
  )
)
