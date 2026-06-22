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
