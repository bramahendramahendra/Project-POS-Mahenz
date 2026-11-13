import { z } from 'zod'

export const paymentSchema = z.object({
  amount:         z.number({ error: 'Jumlah wajib diisi' }).positive('Jumlah harus lebih dari 0'),
  payment_date:   z.string().min(1, 'Tanggal wajib diisi'),
  payment_method: z.string().min(1, 'Metode wajib dipilih'),
  notes:          z.string().optional(),
})

export type PaymentFormValues = z.infer<typeof paymentSchema>
