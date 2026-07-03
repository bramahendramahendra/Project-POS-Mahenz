import { z } from 'zod'

export const createUserSchema = z.object({
  username: z
    .string()
    .trim()
    .min(3, 'Username minimal 3 karakter')
    .max(50, 'Username maksimal 50 karakter')
    .regex(/^[a-zA-Z0-9]+$/, 'Username hanya boleh huruf dan angka'),
  password: z.string().min(6, 'Password minimal 6 karakter'),
  full_name: z.string().trim().min(1, 'Nama wajib diisi').max(100, 'Nama maksimal 100 karakter'),
  role_id: z.number({ error: 'Role wajib dipilih' }).int().positive('Role wajib dipilih'),
})

export const updateUserSchema = z.object({
  full_name: z.string().trim().min(1, 'Nama wajib diisi').max(100, 'Nama maksimal 100 karakter'),
  role_id: z.number({ error: 'Role wajib dipilih' }).int().positive('Role wajib dipilih'),
})

export const changePasswordSchema = z
  .object({
    password: z.string().min(6, 'Password minimal 6 karakter'),
    confirm_password: z.string().min(6, 'Konfirmasi password minimal 6 karakter'),
  })
  .refine((d) => d.password === d.confirm_password, {
    message: 'Konfirmasi password tidak cocok',
    path: ['confirm_password'],
  })

export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>
