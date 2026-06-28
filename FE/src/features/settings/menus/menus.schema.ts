import { z } from 'zod'

export const createMenuSchema = z.object({
  parent_id:   z.number().nullable().optional(),
  key_name:    z.string().min(2, 'Minimal 2 karakter'),
  label:       z.string().min(1, 'Label wajib diisi'),
  icon:        z.string().optional(),
  path:        z.string().optional(),
  order_index: z.number().optional(),
})

export const editMenuSchema = z.object({
  parent_id:   z.number().nullable().optional(),
  label:       z.string().min(1, 'Label wajib diisi'),
  icon:        z.string().optional(),
  path:        z.string().optional(),
  order_index: z.number().optional(),
})

export type CreateMenuFormValues = z.infer<typeof createMenuSchema>
export type EditMenuFormValues = z.infer<typeof editMenuSchema>
