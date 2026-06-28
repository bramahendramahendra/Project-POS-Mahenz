import { z } from 'zod'

export const storeProfileSchema = z.object({
  name: z.string().min(1, 'Nama toko wajib diisi'),
  address: z.string().optional().or(z.literal('')),
  phone: z.string().optional().or(z.literal('')),
  email: z.string().email('Format email tidak valid').optional().or(z.literal('')),
  tax_default: z.number().min(0).max(100).optional(),
})

export type StoreProfileFormValues = z.infer<typeof storeProfileSchema>
