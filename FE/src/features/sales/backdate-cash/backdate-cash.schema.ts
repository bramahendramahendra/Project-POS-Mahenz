import { z } from 'zod'

export const openBackdateCashDrawerSchema = z.object({
  user_id: z.number({ required_error: 'Kasir wajib dipilih' }).min(1, 'Kasir wajib dipilih'),
  date: z.string({ required_error: 'Tanggal wajib diisi' }).min(1, 'Tanggal wajib diisi'),
  shift_id: z.number().nullable().optional(),
  opening_balance: z.number({ required_error: 'Saldo awal wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  notes: z.string().optional(),
})

export type OpenBackdateCashDrawerFormValues = z.infer<typeof openBackdateCashDrawerSchema>

export const closeBackdateCashDrawerSchema = z.object({
  closing_balance: z.number({ required_error: 'Saldo akhir wajib diisi' }).min(0, 'Saldo tidak boleh negatif'),
  notes: z.string().optional(),
})

export type CloseBackdateCashDrawerFormValues = z.infer<typeof closeBackdateCashDrawerSchema>
