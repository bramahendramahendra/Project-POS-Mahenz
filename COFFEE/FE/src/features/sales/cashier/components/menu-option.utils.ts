import { variantLabel } from '@/features/menu-items/menu-items.utils'

import type { CartItem, CatalogItem, CatalogModifierOption, CatalogVariant } from '../cashier.types'

/** Nama kombinasi di keranjang/struk; kosong untuk menu tanpa varian. */
export function catalogVariantName(item: CatalogItem, variant: CatalogVariant): string {
  return item.has_variants ? variantLabel(variant.size_name, variant.temperature) : ''
}

/** Susun baris keranjang dari pilihan kasir. Harga satuan = kombinasi + tambahan (sama dengan BE). */
export function buildCartItem(
  item: CatalogItem,
  variant: CatalogVariant,
  options: Array<CatalogModifierOption & { group_name: string }>,
  qty: number,
  notes: string,
  key: string,
): CartItem {
  const price = variant.price + options.reduce((s, o) => s + o.price, 0)
  return {
    key,
    menu_item_id: item.id,
    variant_id: variant.id,
    product_name: item.name,
    variant_name: catalogVariantName(item, variant),
    unit_name: '',
    modifiers: options.map((o) => ({ option_id: o.id, group_name: o.group_name, option_name: o.name, price: o.price })),
    notes: notes || undefined,
    is_stock_item: item.stock_mode === 'direct',
    qty,
    price,
    subtotal: price * qty,
  }
}

/** Menu tidak bisa dijual: ditandai habis, atau barang jadi dengan stok habis. */
export function isMenuUnavailable(item: CatalogItem): boolean {
  return item.is_sold_out || (item.stock !== null && item.stock <= 0)
}
