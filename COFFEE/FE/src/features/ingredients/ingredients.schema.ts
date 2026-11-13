import { z } from 'zod'

// Jumlah stok boleh desimal sampai 3 angka di belakang koma (kolom DECIMAL(15,3)).
const qtySchema = z
  .number({ error: 'Jumlah wajib diisi' })
  .min(0.001, 'Jumlah minimal 0,001')
  .max(999999999, 'Jumlah terlalu besar')

export const ingredientSchema = z.object({
  name: z.string().trim().min(1, 'Nama bahan wajib diisi').max(100, 'Nama maksimal 100 karakter'),
  unit_id: z.number({ error: 'Satuan wajib dipilih' }).int().positive('Satuan wajib dipilih'),
  min_stock: z.number().min(0, 'Tidak boleh minus').max(999999999).nullable(),
})

export const recordStockSchema = z.object({
  qty: qtySchema,
  notes: z.string().trim().max(255, 'Maksimal 255 karakter'),
})

// Barang terbuang: alasan wajib (keputusan 83).
export const recordWasteSchema = recordStockSchema.extend({
  notes: z.string().trim().min(1, 'Alasan wajib diisi').max(255, 'Maksimal 255 karakter'),
})

// Pembatalan entri: alasan wajib (keputusan 81).
export const cancelMutationSchema = z.object({
  reason: z.string().trim().min(1, 'Alasan wajib diisi').max(255, 'Maksimal 255 karakter'),
})

export type IngredientFormValues = z.infer<typeof ingredientSchema>
export type RecordStockFormValues = z.infer<typeof recordStockSchema>
export type CancelMutationFormValues = z.infer<typeof cancelMutationSchema>
