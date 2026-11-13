import { z } from 'zod'

import type { PaymentMethod } from './cashier.types'

const paymentMethodEnum = z.enum(['cash', 'transfer', 'qris', 'card', 'kredit', 'balance'] as const)

export function createPaymentSchema(grandTotal: number, isKredit: boolean, effectiveTotal: number = grandTotal) {
  // Jika saldo cukup (effectiveTotal = 0), tidak perlu validasi amount
  if (effectiveTotal === 0 || isKredit) {
    return z.object({
      payment_method: paymentMethodEnum,
      amount_paid: z.number().default(0),
    })
  }
  return z
    .object({
      payment_method: paymentMethodEnum,
      amount_paid: z.number().min(0, 'Jumlah bayar wajib diisi'),
    })
    .refine((d) => d.amount_paid >= effectiveTotal, {
      message: 'Jumlah bayar kurang dari total',
      path: ['amount_paid'],
    })
}

export type PaymentFormValues = {
  payment_method: PaymentMethod
  amount_paid?: number
}
