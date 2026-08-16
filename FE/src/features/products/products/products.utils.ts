import * as XLSX from 'xlsx'

import type { PriceTier, Product, ProductPackage } from './products.types'

export function calcMargin(purchasePrice: number, sellingPrice: number): number {
  if (purchasePrice <= 0 || sellingPrice <= 0) return 0
  return Math.round(((sellingPrice - purchasePrice) / sellingPrice) * 100)
}

export function getApplicablePrice(
  prices: PriceTier[],
  unitId: number,
  qty: number
): number | null {
  const tiersByUnit = prices
    .filter((p) => p.unit_id === unitId && p.min_qty <= qty)
    .sort((a, b) => b.min_qty - a.min_qty)

  return tiersByUnit[0]?.price ?? null
}

export function formatProductPackage(pkg: ProductPackage): string {
  return pkg.package_name ? `${pkg.unit_name} (${pkg.package_name})` : pkg.unit_name
}

/**
 * resolved_factor bisa berupa hasil pembagian panjang (mis. 1/120 = 0.008333...) —
 * dibulatkan ke 4 desimal & trailing zero dibuang supaya enak dibaca, tanpa
 * menyembunyikan presisi yang masih berarti untuk satuan yang jauh lebih kecil.
 */
export function formatResolvedFactor(factor: number): string {
  return Number(factor.toFixed(4)).toString()
}

export interface StockBreakdownLevel {
  unit_name: string
  qty: number
}

/**
 * Pecah angka stok (desimal, satuan anchor) jadi kombinasi satuan yang lebih kecil
 * mengikuti resolved_factor tiap paket (sudah dihitung server, sudah menangani rantai
 * konversi berjenjang). Level diurutkan dari satuan fisik terbesar ke terkecil; semua
 * level KECUALI yang terkecil dibulatkan ke bawah (sisa pecahannya "diteruskan" ke level
 * berikutnya) — level terkecil boleh tetap desimal kalau stok tidak habis pas di situ
 * (relevan untuk produk berbasis berat/volume). Level bernilai 0 disaring oleh caller
 * (lihat formatStockBreakdown), bukan di sini, supaya breakdownStock tetap dipakai ulang
 * untuk kebutuhan lain yang butuh angka mentah per level.
 */
export function breakdownStock(stock: number, packages: ProductPackage[]): StockBreakdownLevel[] {
  if (packages.length === 0) return []

  const sorted = [...packages].sort((a, b) => b.resolved_factor - a.resolved_factor)
  const EPSILON = 1e-9
  let remaining = stock
  const levels: StockBreakdownLevel[] = []

  sorted.forEach((pkg, i) => {
    const isLast = i === sorted.length - 1
    const countInThisUnit = remaining / pkg.resolved_factor
    if (isLast) {
      levels.push({ unit_name: pkg.unit_name, qty: countInThisUnit })
    } else {
      const qty = Math.floor(countInThisUnit + EPSILON)
      levels.push({ unit_name: pkg.unit_name, qty })
      remaining -= qty * pkg.resolved_factor
    }
  })

  return levels
}

/**
 * Format hasil breakdownStock jadi teks siap-tampil, mis. "5 Pack" atau "1 Krak 5 Pieces".
 * Level bernilai 0 disembunyikan (di posisi manapun — depan/tengah/akhir); kalau semua
 * level 0 (stok benar-benar habis), tampilkan level terkecil saja dengan angka 0 supaya
 * tidak kosong.
 */
export function formatStockBreakdown(stock: number, packages: ProductPackage[]): string {
  const levels = breakdownStock(stock, packages)
  if (levels.length === 0) return String(stock)

  const nonZero = levels.filter((l) => l.qty !== 0)
  const toShow = nonZero.length > 0 ? nonZero : [levels[levels.length - 1]]

  return toShow.map((l) => `${formatResolvedFactor(l.qty)} ${l.unit_name}`).join(' ')
}

export function getDisplayPrice(product: Product): number {
  const defaultUnit = (product.units ?? []).find((u) => u.is_default)
  if (!defaultUnit) return product.selling_price
  const tiers = (product.prices ?? [])
    .filter((p) => p.unit_id === defaultUnit.unit_id)
    .sort((a, b) => a.min_qty - b.min_qty)
  return tiers[0]?.price ?? product.selling_price
}

export function exportProductsToExcel(products: Product[]): void {
  const rows = products.map((p) => ({
    'Nama Produk': p.name,
    Barcode: p.barcode ?? '',
    SKU: p.sku ?? '',
    Kategori: p.category_name ?? '',
    'Harga Beli': p.purchase_price,
    'Harga Jual': p.selling_price,
    Stok: p.stock,
    'Stok Minimum': p.min_stock,
    Satuan: p.unit_name ?? '',
    Status: p.is_active ? 'Aktif' : 'Nonaktif',
  }))
  const ws = XLSX.utils.json_to_sheet(rows)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Produk')
  XLSX.writeFile(wb, `produk-export-${Date.now()}.xlsx`)
}
