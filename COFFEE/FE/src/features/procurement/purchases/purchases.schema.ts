import { z } from 'zod'
import { todayStr } from '@/shared/utils'

// Baris pembelian: item_type product = barang jadi (produk + kemasan), ingredient = bahan baku
// (jumlah desimal dalam satuan stok bahan, keputusan 92). Pilihan wajib per jenis dicek di
// refineItemTargets.
export const purchaseItemSchema = z.object({
  item_type: z.enum(['product', 'ingredient']),
  product_id: z.number(),
  product_name: z.string().optional(),
  ingredient_id: z.number().optional(),
  ingredient_name: z.string().optional(),
  package_id: z.number().optional(),
  quantity: z.number({ error: 'Wajib diisi' }).positive('Harus lebih dari 0'),
  price: z.number({ error: 'Wajib diisi' }).positive('Harga harus lebih dari 0'),
  unit: z.string().min(1, 'Wajib diisi'),
  conversion_qty: z.number().min(1).catch(1),
  // Satuan kontinu (mis. Kilogram) boleh qty pecahan; satuan diskrit (Pcs/Slop/dll.)
  // wajib bilangan bulat -- dicek di superRefine (bukan di sini) karena aturannya
  // beda per baris item, bukan aturan tetap untuk semua qty.
  is_continuous: z.boolean().optional(),
})

type ItemTarget = { item_type: 'product' | 'ingredient'; product_id: number; ingredient_id?: number; package_id?: number }
type IssueCtx = { addIssue: (issue: { code: 'custom'; message: string; path: (string | number)[] }) => void }

// Dipakai di superRefine purchaseSchema & addItemsSchema: barang/bahan wajib dipilih sesuai
// jenis baris, dan baris yang sama (produk + satuan, atau bahan yang sama) tidak boleh dobel.
// Kesalahan ditaruh di path product_id (kolom pilihan barang/bahan).
export function refineItemTargets(items: ItemTarget[], ctx: IssueCtx, basePath: (string | number)[] = ['items']) {
  const seen = new Set<string>()
  items.forEach((item, index) => {
    const path = [...basePath, index, 'product_id']
    const isIngredient = item.item_type === 'ingredient'
    const id = isIngredient ? item.ingredient_id || 0 : item.product_id
    if (!id) {
      ctx.addIssue({ code: 'custom', message: isIngredient ? 'Pilih bahan baku' : 'Pilih barang jadi', path })
      return
    }
    const key = isIngredient ? `i-${id}` : `p-${id}-${item.package_id || 0}`
    if (seen.has(key)) {
      ctx.addIssue({
        code: 'custom',
        message: isIngredient ? 'Bahan baku ini sudah dipilih di baris lain' : 'Produk dengan satuan ini sudah dipilih di baris lain',
        path,
      })
    } else {
      seen.add(key)
    }
  })
}

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

    refineItemTargets(data.items, ctx)
  })

export type PurchaseItemFormValues = z.infer<typeof purchaseItemSchema>
export type PurchaseFormValues = z.infer<typeof purchaseSchema>
