import { z } from 'zod'
import { todayStr } from '@/shared/utils'

// expiryBatchDraftSchema: 1 baris rincian "sebagian qty item ini expired tanggal X".
// Total qty di semua baris untuk 1 item dicek sama dengan quantity item itu lewat
// superRefine di bawah (butuh baca field lain, tidak bisa lewat tag field tunggal).
export const expiryBatchDraftSchema = z.object({
  qty: z.number({ error: 'Wajib diisi' }).positive('Harus lebih dari 0'),
  expired_date: z.string().min(1, 'Wajib diisi'),
})

export const purchaseItemSchema = z.object({
  product_id: z.number({ error: 'Pilih produk' }).positive('Pilih produk'),
  product_name: z.string().optional(),
  package_id: z.number().optional(),
  quantity: z.number({ error: 'Wajib diisi' }).positive('Harus lebih dari 0'),
  price: z.number({ error: 'Wajib diisi' }).positive('Harga harus lebih dari 0'),
  unit: z.string().min(1, 'Wajib diisi'),
  conversion_qty: z.number().min(1).catch(1),
  // Satuan kontinu (mis. Kilogram) boleh qty pecahan; satuan diskrit (Pcs/Slop/dll.)
  // wajib bilangan bulat -- dicek di superRefine (bukan di sini) karena aturannya
  // beda per baris item, bukan aturan tetap untuk semua qty.
  is_continuous: z.boolean().optional(),
  expiry_batches: z.array(expiryBatchDraftSchema).optional(),
})

// Dipakai di superRefine purchaseSchema & addItemsSchema (2 tempat) -- qty item
// wajib bilangan bulat KECUALI satuan yang dipakai kontinu (mis. Kilogram).
export function refineItemQuantities(
  items: { quantity: number; is_continuous?: boolean }[],
  ctx: { addIssue: (issue: { code: 'custom'; message: string; path: (string | number)[] }) => void },
  basePath: (string | number)[] = ['items'],
) {
  items.forEach((item, index) => {
    if (!item.is_continuous && !Number.isInteger(item.quantity)) {
      ctx.addIssue({
        code: 'custom',
        message: 'Qty harus bilangan bulat',
        path: [...basePath, index, 'quantity'],
      })
    }
  })
}

export const purchaseSchema = z
  .object({
    purchase_date: z
      .string()
      .min(1, 'Tanggal wajib diisi')
      .refine((val) => val <= todayStr(), 'Tanggal tidak boleh lebih dari hari ini'),
    invoice_number: z.string().max(50, 'No. faktur maksimal 50 karakter'),
    no_invoice: z.boolean().optional(),
    supplier_id: z.number({ error: 'Pilih supplier' }).positive('Pilih supplier'),
    items: z.array(purchaseItemSchema).min(1, 'Minimal 1 item'),
    discount_amount: z.number().nonnegative(),
    notes: z.string().max(500, 'Catatan maksimal 500 karakter').optional(),
    payment_status: z.union([z.enum(['paid', 'unpaid', 'partial']), z.literal('')]),
    paid_amount: z.number().nonnegative(),
    payment_method: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    const subtotal = data.items.reduce((sum, item) => sum + (item.quantity || 0) * (item.price || 0), 0)
    const total = Math.max(0, subtotal - (data.discount_amount || 0))

    if (!data.no_invoice && data.invoice_number.trim().length === 0) {
      ctx.addIssue({
        code: 'custom',
        message: 'No. faktur wajib diisi',
        path: ['invoice_number'],
      })
    }

    if (data.discount_amount > subtotal) {
      ctx.addIssue({
        code: 'custom',
        message: 'Diskon tidak boleh lebih besar dari subtotal',
        path: ['discount_amount'],
      })
    }

    if (data.payment_status === '') {
      ctx.addIssue({
        code: 'custom',
        message: 'Pilih status pembayaran',
        path: ['payment_status'],
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

    if (data.payment_status !== '' && data.payment_status !== 'unpaid' && !data.payment_method) {
      ctx.addIssue({
        code: 'custom',
        message: 'Metode pembayaran wajib dipilih',
        path: ['payment_method'],
      })
    }

    refineItemQuantities(data.items, ctx)

    // Validasi duplikat: produk + satuan yang sama tidak boleh muncul di >1 baris.
    // Produk sama BOLEH muncul berkali-kali selama satuannya (package_id) berbeda —
    // contoh: beli 2 Slop + 5 Pack rokok yang sama dalam 1 nota.
    const seen = new Map<string, number>()
    data.items.forEach((item, index) => {
      if (!item.product_id) return
      // Kunci: product_id + package_id. Jika package_id belum dipilih (0/undefined),
      // fallback ke product_id saja supaya tetap ada pengecekan minimal.
      const key = `${item.product_id}-${item.package_id || 0}`
      const firstIndex = seen.get(key)
      if (firstIndex !== undefined) {
        ctx.addIssue({
          code: 'custom',
          message: 'Produk dengan satuan ini sudah dipilih di baris lain',
          path: ['items', index, 'product_id'],
        })
      } else {
        seen.set(key, index)
      }
    })

    data.items.forEach((item, index) => {
      if (!item.expiry_batches || item.expiry_batches.length === 0) return
      const totalBatchQty = item.expiry_batches.reduce((sum, b) => sum + (b.qty || 0), 0)
      if (Math.abs(totalBatchQty - item.quantity) > 0.0001) {
        ctx.addIssue({
          code: 'custom',
          message: `Total qty batch expired (${totalBatchQty}) harus sama dengan qty item (${item.quantity})`,
          path: ['items', index, 'expiry_batches'],
        })
      }
    })
  })

export type PurchaseItemFormValues = z.infer<typeof purchaseItemSchema>
export type PurchaseFormValues = z.infer<typeof purchaseSchema>
