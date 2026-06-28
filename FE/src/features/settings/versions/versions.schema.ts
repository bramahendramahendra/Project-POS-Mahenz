import { z } from 'zod'

export const createAppVersionSchema = z.object({
  version:      z.string().min(1, 'Versi wajib diisi').regex(/^\d+\.\d+\.\d+$/, 'Format versi harus x.y.z'),
  download_url: z.string().url('URL download tidak valid'),
  release_notes: z.string().optional(),
  is_mandatory: z.boolean(),
})

export type CreateAppVersionFormValues = z.infer<typeof createAppVersionSchema>
