import { z } from 'zod'

// Toko hanya mengubah role profil lokal. Nama dan password dikelola di Portal (keputusan 69, 76).
export const updateUserSchema = z.object({
  role_id: z.number({ error: 'Role wajib dipilih' }).int().positive('Role wajib dipilih'),
})

export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
