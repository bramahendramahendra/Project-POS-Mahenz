import { z } from 'zod'

export const shiftFormSchema = z.object({
  name: z
    .string()
    .min(2, 'Nama shift minimal 2 karakter')
    .max(100, 'Nama shift maksimal 100 karakter'),
  start_time: z.string().min(1, 'Jam mulai wajib diisi'),
  end_time: z.string().min(1, 'Jam selesai wajib diisi'),
})

export type ShiftFormValues = z.infer<typeof shiftFormSchema>
