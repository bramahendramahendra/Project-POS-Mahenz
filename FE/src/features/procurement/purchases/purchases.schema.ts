import { z } from 'zod'

export const purchaseItemSchema = z.object({
  product_id: z.number({ error: 'Pilih produk' }).positive('Pilih produk'),
  product_name: z.string().optional(),
  quantity: z.number({ error: 'Wajib diisi' }).positive('Harus lebih dari 0'),
  price: z.number({ error: 'Wajib diisi' }).nonnegative(),
  unit: z.string().min(1, 'Wajib diisi'),
  conversion_qty: z.number().min(1).catch(1),
})

export const purchaseSchema = z.object({
  purchase_date: z.string().min(1, 'Tanggal wajib diisi'),
  invoice_number: z.string().min(1, 'No. faktur wajib diisi'),
  supplier_id: z.number({ error: 'Pilih supplier' }).positive('Pilih supplier'),
  items: z.array(purchaseItemSchema).min(1, 'Minimal 1 item'),
  discount_amount: z.number().nonnegative(),
  notes: z.string().optional(),
  payment_status: z.enum(['paid', 'unpaid', 'partial']),
  paid_amount: z.number().nonnegative(),
  payment_method: z.string().optional(),
})

export type PurchaseItemFormValues = z.infer<typeof purchaseItemSchema>
export type PurchaseFormValues = z.infer<typeof purchaseSchema>
