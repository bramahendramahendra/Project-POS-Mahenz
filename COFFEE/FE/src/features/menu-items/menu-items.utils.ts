import { formatRupiah } from '@/shared/utils'

import type { StockMode, Temperature } from './menu-items.types'

export const TEMPERATURE_LABEL: Record<Temperature, string> = {
  hot: 'Panas',
  iced: 'Dingin',
}

export const TEMPERATURES: Temperature[] = ['hot', 'iced']

export const STOCK_MODE_LABEL: Record<StockMode, string> = {
  recipe: 'Racikan',
  direct: 'Barang jadi',
}

/** Nama kombinasi varian, mis. "Regular · Dingin". */
export function variantLabel(sizeName: string | null | undefined, temperature: Temperature | null): string {
  const parts = [sizeName || null, temperature ? TEMPERATURE_LABEL[temperature] : null].filter(Boolean)
  return parts.length > 0 ? parts.join(' · ') : 'Standar'
}

/** Harga satu atau rentang (mis. "Rp 22.000 - Rp 29.000"). */
export function formatPriceRange(min: number, max: number): string {
  if (min === max) return formatRupiah(min)
  return `${formatRupiah(min)} - ${formatRupiah(max)}`
}

/** Kunci kombinasi untuk mencocokkan varian lama dan baru di form. */
export function variantKey(sizeId: number | null, temperature: Temperature | null): string {
  return `${sizeId ?? 0}|${temperature ?? ''}`
}
