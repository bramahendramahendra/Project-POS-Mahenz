import { z } from 'zod'
import { todayStr } from '@/shared/utils'

export const purchaseItemSchema = z.object({
  product_id: z.number({ error: 'Pilih produk' }).positive('Pilih produk'),
  product_name: z.string().optional(),
  quantity: z.number({ error: 'Wajib diisi' }).positive('Harus lebih dari 0').int('Qty harus bilangan bulat'),
  price: z.number({ error: 'Wajib diisi' }).positive('Harga harus lebih dari 0'),
  unit: z.string().min(1, 'Wajib diisi'),
  conversion_qty: z.number().min(1).catch(1),
})

export const purchaseSchema = z
  .object({
    purchase_date: z
      .string()
      .min(1, 'Tanggal wajib diisi')
      .refine((val) => val <= todayStr(), 'Tanggal tidak boleh lebih dari hari ini'),
    invoice_number: z.string().min(1, 'No. faktur wajib diisi').max(50, 'No. faktur maksimal 50 karakter'),
    supplier_id: z.number({ error: 'Pilih supplier' }).positive('Pilih supplier'),
    items: z.array(purchaseItemSchema).min(1, 'Minimal 1 item'),
    discount_amount: z.number().nonnegative(),
    notes: z.string().max(500, 'Catatan maksimal 500 karakter').optional(),
    payment_status: z.enum(['paid', 'unpaid', 'partial']),
    paid_amount: z.number().nonnegative(),
    payment_method: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    const subtotal = data.items.reduce((sum, item) => sum + (item.quantity || 0) * (item.price || 0), 0)
    const total = Math.max(0, subtotal - (data.discount_amount || 0))

    if (data.discount_amount > subtotal) {
      ctx.addIssue({
        code: 'custom',
        message: 'Diskon tidak boleh lebih besar dari subtotal',
        path: ['discount_amount'],
      })
    }

    if (data.payment_status === 'partial') {
      if (data.paid_amount <= 0) {
        ctx.addIssue({
          code: 'custom',
          message: 'Jumlah dibayar harus lebih dari 0',
          path: ['paid_amount'],
        })
      } else if (data.paid_amount >= total) {
        ctx.addIssue({
          code: 'custom',
          message: 'Jumlah dibayar harus kurang dari total untuk status partial',
          path: ['paid_amount'],
        })
      }
    }

    if (data.payment_status !== 'unpaid' && !data.payment_method) {
      ctx.addIssue({
        code: 'custom',
        message: 'Metode pembayaran wajib dipilih',
        path: ['payment_method'],
      })
    }

    const seen = new Map<number, number>()
    data.items.forEach((item, index) => {
      if (!item.product_id) return
      const firstIndex = seen.get(item.product_id)
      if (firstIndex !== undefined) {
        ctx.addIssue({
          code: 'custom',
          message: 'Produk sudah dipilih di baris lain',
          path: ['items', index, 'product_id'],
        })
      } else {
        seen.set(item.product_id, index)
      }
    })
  })

export type PurchaseItemFormValues = z.infer<typeof purchaseItemSchema>
export type PurchaseFormValues = z.infer<typeof purchaseSchema>
