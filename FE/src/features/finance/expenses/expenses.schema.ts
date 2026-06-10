import { z } from 'zod'

export const EXPENSE_PAYMENT_METHODS = [
  { value: 'cash', label: 'Tunai' },
  { value: 'transfer', label: 'Transfer' },
  { value: 'card', label: 'Kartu' },
  { value: 'qris', label: 'QRIS' },
  { value: 'kredit', label: 'Kredit' },
] as const

export const EXPENSE_CATEGORIES = [
  { value: 'operasional', label: 'Operasional' },
  { value: 'pembelian', label: 'Pembelian' },
  { value: 'gaji', label: 'Gaji' },
  { value: 'lainnya', label: 'Lainnya' },
] as const

export const expenseSchema = z.object({
  expense_date: z.string().min(1, 'Tanggal wajib diisi'),
  category: z.enum(['operasional', 'pembelian', 'gaji', 'lainnya'], {
    required_error: 'Kategori wajib dipilih',
  }),
  description: z.string().trim().min(1, 'Keterangan wajib diisi').max(255, 'Keterangan maksimal 255 karakter'),
  amount: z.number({ invalid_type_error: 'Jumlah wajib diisi' }).positive('Jumlah harus lebih dari 0'),
  payment_method: z.enum(['cash', 'transfer', 'card', 'qris', 'kredit'], {
    required_error: 'Metode pembayaran wajib dipilih',
  }),
  notes: z.string().max(500, 'Catatan maksimal 500 karakter').optional().or(z.literal('')),
})

export type ExpenseFormValues = z.infer<typeof expenseSchema>
