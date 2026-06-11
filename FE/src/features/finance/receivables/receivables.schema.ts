import { z } from 'zod'

import { formatRupiah } from '@/shared/utils'

export function createPaymentSchema(remaining: number) {
  return z.object({
    amount: z
      .number({ error: 'Jumlah bayar wajib diisi' })
      .min(1, 'Jumlah bayar wajib diisi')
      .max(remaining, `Maksimal pembayaran ${formatRupiah(remaining)}`),
    payment_date: z.string().min(1, 'Tanggal wajib diisi'),
    notes: z.string().optional(),
  })
}

export type PaymentFormValues = z.infer<ReturnType<typeof createPaymentSchema>>
