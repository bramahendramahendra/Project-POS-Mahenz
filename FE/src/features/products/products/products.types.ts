export interface ProductPackage {
  id: number
  product_id: number
  unit_id: number
  unit_name: string
  abbreviation: string
  package_name: string
  /** null hanya untuk paket anchor (is_default true) */
  ref_package_id: number | null
  qty: number
  ref_qty: number | null
  /** 1 paket ini = resolved_factor x satuan anchor produk, dihitung server-side */
  resolved_factor: number
  purchase_price: number
  selling_price: number
  is_default: boolean
  /** Stok LEVEL INI SENDIRI (bukan dikonversi ke anchor) — sumber breakdown per-level, Fase 5/6. */
  stock: number
  reserved_qty: number
  is_active: boolean
}

export interface PriceTier {
  id: number
  product_id: number
  unit_id: number
  unit_name: string
  tier_name: string
  min_qty: number
  price: number
}

export interface Product {
  id: number
  name: string
  sku?: string
  barcode?: string
  category_id?: number
  category_name?: string
  purchase_price: number
  selling_price: number
  stock: number
  reserved_qty: number
  min_stock: number
  unit_id: number
  unit_name: string
  unit_abbreviation: string
  is_active: boolean
  created_at: string
  units: ProductPackage[]
  prices: PriceTier[]
  extra_packages: number
  price_tiers_count: number
  is_low_stock: boolean
  needs_stock_review: boolean
  stock_review_note: string
}

export interface ProductListFilter {
  page: number
  limit: number
  search: string
  category_id?: number
  is_active?: boolean
  low_stock?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface PackageDraftPayload {
  temp_id: number
  unit_id: number
  package_name?: string
  ref_temp_id: number
  qty: number
  ref_qty: number
  purchase_price: number
  selling_price: number
}

export interface CreateProductPayload {
  name: string
  sku: string
  barcode: string
  category_id: number
  purchase_price: number
  selling_price: number
  stock: number
  min_stock: number
  unit_id: number
  is_active: boolean
  packages?: PackageDraftPayload[]
}

export type UpdateProductPayload = Partial<Omit<CreateProductPayload, 'unit_id'>>

export interface CreatePackagePayload {
  unit_id: number
  package_name?: string
  ref_package_id: number
  qty: number
  ref_qty: number
  purchase_price: number
  selling_price: number
}

export type UpdatePackagePayload = CreatePackagePayload

export type ProductFilter = ProductListFilter

export interface ProductOption {
  id: number
  name: string
}

export interface ProductSearchOption {
  id: number
  barcode: string
  name: string
  selling_price: number
  stock: number
  min_stock: number
  unit_id: number
  unit_name: string
}

// ─── Import Types ─────────────────────────────────────────────────────────────

export interface ImportPreviewRow {
  no: number
  nama: string
  barcode: string
  kategori: string
  harga_beli: number
  harga_jual: number
  margin: number
  stok: number
  stok_minimum: number
  satuan: string
  satuan_id: number
  valid: boolean
  errors: string[]
  warnings: string[]
}

export interface ImportPreviewGrosirRow {
  no_produk: number
  nama_paket: string
  satuan: string
  satuan_id: number
  qty: number
  ref_qty: number
  harga_beli: number
  harga_jual: number
  valid: boolean
  errors: string[]
}

export interface ImportPreviewResponse {
  rows: ImportPreviewRow[]
  grosir: ImportPreviewGrosirRow[]
}

export interface ImportBulkRow {
  no: number
  nama: string
  barcode: string
  kategori: string
  harga_beli: number
  harga_jual: number
  stok: number
  stok_minimum: number
  satuan: string
  satuan_id: number
}

export interface GrosirImportRow {
  no_produk: number
  nama_paket: string
  satuan: string
  satuan_id: number
  qty: number
  ref_qty: number
  harga_beli: number
  harga_jual: number
}

export interface ImportBulkResult {
  success: number
  failed: { baris: number; data: ImportBulkRow; alasan: string }[]
}

export interface ImportBulkPayload {
  rows: ImportBulkRow[]
  grosir: GrosirImportRow[]
}
