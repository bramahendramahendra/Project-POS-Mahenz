import type { MutationType } from './ingredients.types'

const qtyFormat = new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 })

/** Jumlah stok bahan dengan pemisah ribuan Indonesia, maksimal 3 angka desimal. */
export function formatQty(value: number, unit?: string): string {
  const text = qtyFormat.format(value)
  return unit ? `${text} ${unit}` : text
}

/** Jumlah perubahan stok dengan tanda + / -. */
export function formatQtyChange(value: number, unit?: string): string {
  const sign = value > 0 ? '+' : value < 0 ? '-' : ''
  return `${sign}${formatQty(Math.abs(value), unit)}`
}

export const MUTATION_TYPE_LABEL: Record<MutationType, string> = {
  usage: 'Terpakai',
  waste: 'Terbuang',
  cancel: 'Pembatalan',
  purchase: 'Pembelian',
  purchase_void: 'Void pembelian',
}

export const MUTATION_TYPE_CLASS: Record<MutationType, string> = {
  usage: 'bg-blue-100 text-blue-700',
  waste: 'bg-red-100 text-red-700',
  cancel: 'bg-gray-100 text-gray-700',
  purchase: 'bg-green-100 text-green-700',
  purchase_void: 'bg-amber-100 text-amber-700',
}

/** Mengubah teks input desimal ("1,5" atau "1.5") menjadi angka; kosong atau tidak valid = NaN. */
export function parseDecimalInput(text: string): number {
  const cleaned = text.trim().replace(/\s/g, '').replace(',', '.')
  if (cleaned === '') return Number.NaN
  return Number(cleaned)
}
