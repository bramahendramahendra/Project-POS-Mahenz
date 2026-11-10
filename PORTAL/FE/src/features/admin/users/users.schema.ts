import { z } from 'zod'

const passwordSchema = z
  .string()
  .min(8, 'Password minimal 8 karakter')
  .max(72, 'Password maksimal 72 karakter')

export const createUserSchema = z.object({
  username: z
    .string()
    .trim()
    .min(3, 'Username minimal 3 karakter')
    .max(50, 'Username maksimal 50 karakter')
    .regex(/^[A-Za-z0-9]+$/, 'Username hanya huruf dan angka'),
  password: passwordSchema,
  full_name: z.string().trim().min(1, 'Nama lengkap wajib diisi').max(100, 'Maksimal 100 karakter'),
  role_id: z.number().int().gt(0, 'Pilih role'),
})

export const updateUserSchema = createUserSchema.omit({ username: true, password: true })

export const resetPasswordSchema = z.object({ password: passwordSchema })

export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>
