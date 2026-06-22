import { z } from 'zod'

export const openShiftSchema = z.object({
  opening_balance: z.number({ error: 'Modal awal wajib diisi' }).min(0, 'Modal awal tidak boleh negatif'),
  notes: z.string().optional(),
})

export const closeShiftSchema = z.object({
  closing_balance: z.number({ error: 'Uang akhir wajib diisi' }).min(0, 'Uang akhir tidak boleh negatif'),
  notes: z.string().optional(),
})

export type OpenShiftFormValues = z.infer<typeof openShiftSchema>
export type CloseShiftFormValues = z.infer<typeof closeShiftSchema>
