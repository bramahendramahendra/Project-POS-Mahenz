import { z } from 'zod'

export const loginSchema = z.object({
  username: z.string().trim().min(1, 'Username wajib diisi'),
  password: z.string().min(1, 'Password wajib diisi'),
})

export type LoginFormValues = z.infer<typeof loginSchema>

export const changePasswordSchema = z
  .object({
    oldPassword: z.string().min(1, 'Password saat ini wajib diisi'),
    newPassword: z
      .string()
      .min(8, 'Password baru minimal 8 karakter')
      .max(72, 'Password baru maksimal 72 karakter'),
    confirmPassword: z.string().min(1, 'Konfirmasi password wajib diisi'),
  })
  .refine((v) => v.newPassword === v.confirmPassword, {
    path: ['confirmPassword'],
    message: 'Konfirmasi password tidak sama',
  })
  .refine((v) => v.newPassword !== v.oldPassword, {
    path: ['newPassword'],
    message: 'Password baru harus berbeda dari password saat ini',
  })

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>
