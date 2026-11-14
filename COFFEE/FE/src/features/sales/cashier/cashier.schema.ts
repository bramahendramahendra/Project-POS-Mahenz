import { z } from 'zod'

import type { PAYMENT_METHOD_VALUES } from './cashier.types'

const paymentMethodEnum = z.enum(['cash', 'transfer', 'qris', 'card', 'kredit'] as const)

export function createPaymentSchema(grandTotal: number, isKredit: boolean) {
  if (grandTotal === 0 || isKredit) {
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
    .refine((d) => d.amount_paid >= grandTotal, {
      message: 'Jumlah bayar kurang dari total',
      path: ['amount_paid'],
    })
}

export type PaymentFormValues = {
  payment_method: (typeof PAYMENT_METHOD_VALUES)[number]
  amount_paid?: number
}
