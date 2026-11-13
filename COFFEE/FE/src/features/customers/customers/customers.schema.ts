import { z } from 'zod'

export const customerSchema = z.object({
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(100, 'Nama maksimal 100 karakter'),
  phone: z.string().max(20, 'Nomor telepon maksimal 20 karakter').optional().or(z.literal('')),
  address: z.string().max(255, 'Alamat maksimal 255 karakter').optional().or(z.literal('')),
  credit_limit: z.number().min(0, 'Limit kredit tidak boleh negatif').optional(),
  notes: z.string().max(500, 'Catatan maksimal 500 karakter').optional().or(z.literal('')),
})

export type CustomerFormValues = z.infer<typeof customerSchema>
