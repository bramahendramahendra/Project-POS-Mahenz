import { z } from 'zod'

const today = new Date().toISOString().slice(0, 10)

export const returnSchema = z.object({
  purchase_id: z.number({ error: 'Pilih pembelian' }).positive('Pilih pembelian'),
  return_date: z
    .string()
    .min(1, 'Tanggal wajib diisi')
    .refine((v) => v <= today, 'Tanggal retur tidak boleh lebih dari hari ini'),
  reason: z.string().min(1, 'Alasan wajib diisi'),
  notes: z.string().optional(),
})

export type ReturnFormValues = z.infer<typeof returnSchema>
