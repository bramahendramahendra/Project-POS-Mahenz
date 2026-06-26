import { z } from 'zod'

export const openCashDrawerSchema = z.object({
  shift_id: z.number({ error: 'Shift wajib dipilih' }).min(1, 'Shift wajib dipilih'),
  opening_balance: z.number({ error: 'Saldo awal tunai wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  notes: z.string().optional(),
})

export const closeCashDrawerSchema = z.object({
  closing_balance: z.number({ error: 'Saldo akhir tunai wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  notes: z.string().optional(),
})

export type OpenCashDrawerFormValues = z.infer<typeof openCashDrawerSchema>
export type CloseCashDrawerFormValues = z.infer<typeof closeCashDrawerSchema>
