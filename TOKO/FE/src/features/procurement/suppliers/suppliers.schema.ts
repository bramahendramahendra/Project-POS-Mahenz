import { z } from 'zod'

export const supplierSchema = z.object({
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(100, 'Nama maksimal 100 karakter'),
  contact_person: z.string().max(100, 'Nama kontak maksimal 100 karakter').optional().or(z.literal('')),
  phone: z.string().max(20, 'No. telepon maksimal 20 karakter').optional().or(z.literal('')),
  email: z.string().email('Format email tidak valid').max(100).optional().or(z.literal('')),
  address: z.string().max(255, 'Alamat maksimal 255 karakter').optional().or(z.literal('')),
  notes: z.string().max(500, 'Catatan maksimal 500 karakter').optional().or(z.literal('')),
})

export type SupplierFormValues = z.infer<typeof supplierSchema>
