import { z } from 'zod'

export const openCashDrawerSchema = z.object({
  opening_balance: z.number({ error: 'Saldo awal wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  shift: z.enum(['pagi', 'siang', 'malam']).optional(),
  notes: z.string().optional(),
})

export const closeCashDrawerSchema = z.object({
  closing_balance: z.number({ error: 'Saldo penutupan wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  notes: z.string().optional(),
})

export type OpenCashDrawerFormValues = z.infer<typeof openCashDrawerSchema>
export type CloseCashDrawerFormValues = z.infer<typeof closeCashDrawerSchema>
