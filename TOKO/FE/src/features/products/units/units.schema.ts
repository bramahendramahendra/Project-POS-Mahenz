import { z } from 'zod'

export const unitSchema = z.object({
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(100, 'Nama maksimal 100 karakter'),
  abbreviation: z.string().trim().min(2, 'Singkatan minimal 2 karakter').max(20, 'Singkatan maksimal 20 karakter'),
})

export type UnitFormValues = z.infer<typeof unitSchema>
