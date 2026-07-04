import { z } from 'zod'

export const productSchema = z
  .object({
    name: z.string().min(1, 'Nama produk wajib diisi').max(100, 'Nama maksimal 100 karakter'),
    sku: z.string().min(1, 'SKU wajib digenerate'),
    barcode: z.string().min(1, 'Barcode wajib digenerate'),
    category_id: z.number({ error: 'Kategori wajib dipilih' }).min(1, 'Kategori wajib dipilih'),
    purchase_price: z.number().min(0, 'Harga beli tidak boleh negatif'),
    selling_price: z.number().min(1, 'Harga jual harus lebih dari 0'),
    stock: z.number().min(0, 'Stok tidak boleh negatif'),
    min_stock: z.number().min(0, 'Stok minimum tidak boleh negatif'),
    unit_id: z.number({ error: 'Satuan wajib dipilih' }).min(1, 'Satuan wajib dipilih'),
    is_active: z.boolean(),
  })
  .refine((v) => v.selling_price >= v.purchase_price, {
    message: 'Harga jual tidak boleh lebih rendah dari harga beli',
    path: ['selling_price'],
  })

export type ProductFormValues = z.infer<typeof productSchema>

export const grosirSchema = z.object({
  unit_id: z.number({ error: 'Satuan wajib dipilih' }).min(1, 'Satuan wajib dipilih'),
  package_name: z.string().max(100).optional(),
  conversion_qty: z.number().min(1, 'Konversi harus > 0'),
  purchase_price: z.number().min(0),
  selling_price: z.number().min(1, 'Harga jual harus > 0'),
})

export type GrosirFormValues = z.infer<typeof grosirSchema>

export const priceTierSchema = z.object({
  unit_id: z.number({ error: 'Satuan wajib dipilih' }),
  tier_name: z.string().min(1, 'Nama tier wajib diisi'),
  min_qty: z.number().min(1, 'Minimal qty harus > 0'),
  price: z.number().min(0, 'Harga tidak boleh negatif'),
})
export type PriceTierFormValues = z.infer<typeof priceTierSchema>
