import { z } from 'zod'

export const createRoleSchema = z.object({
  name: z
    .string()
    .min(2, 'Minimal 2 karakter')
    .regex(/^[a-z0-9_]+$/, 'Hanya huruf kecil, angka, dan underscore'),
  display_name: z.string().min(1, 'Label wajib diisi'),
  description: z.string().optional(),
})

export const editRoleSchema = z.object({
  display_name: z.string().min(1, 'Label wajib diisi'),
  description: z.string().optional(),
})

export type CreateRoleFormValues = z.infer<typeof createRoleSchema>
export type EditRoleFormValues = z.infer<typeof editRoleSchema>
