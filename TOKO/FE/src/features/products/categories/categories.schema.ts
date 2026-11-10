import { z } from 'zod'

export const categorySchema = z.object({
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(100, 'Nama maksimal 100 karakter'),
  description: z.string().max(500, 'Deskripsi maksimal 500 karakter').optional().or(z.literal('')),
})

export type CategoryFormValues = z.infer<typeof categorySchema>
