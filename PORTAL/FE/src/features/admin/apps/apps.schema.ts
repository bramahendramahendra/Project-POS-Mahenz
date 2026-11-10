import { z } from 'zod'

const pathSchema = z
  .string()
  .trim()
  .min(1, 'Alamat wajib diisi')
  .max(255, 'Maksimal 255 karakter')
  .refine(
    (v) => (v.startsWith('/') && !v.startsWith('//')) || /^https?:\/\//.test(v),
    "Isi path diawali '/' (mis. /toko/) atau URL http/https"
  )

export const createAppSchema = z.object({
  code: z
    .string()
    .trim()
    .regex(/^[a-z][a-z0-9-]{1,49}$/, 'Huruf kecil, angka, dan strip; diawali huruf; 2-50 karakter'),
  name: z.string().trim().min(1, 'Nama wajib diisi').max(100, 'Maksimal 100 karakter'),
  description: z.string().trim().max(255, 'Maksimal 255 karakter'),
  path: pathSchema,
  icon: z.string().trim().max(50, 'Maksimal 50 karakter'),
  sort_order: z.number().int(),
})

export const updateAppSchema = createAppSchema.omit({ code: true })

export type CreateAppFormValues = z.infer<typeof createAppSchema>
export type UpdateAppFormValues = z.infer<typeof updateAppSchema>
