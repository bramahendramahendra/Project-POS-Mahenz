// Menu jualan Coffee (tahap 4.4, keputusan 43, 53-55, 85-89).

export type StockMode = 'recipe' | 'direct'
export type Temperature = 'hot' | 'iced'
export type SelectionType = 'single' | 'multiple'

// ─── Kategori dan ukuran (master sederhana) ───────────────────────────────────

export interface MenuCategory {
  id: number
  name: string
  sort_order: number
  is_active: boolean
  item_count: number
}

export interface MenuSize {
  id: number
  name: string
  sort_order: number
  is_active: boolean
  usage_count: number
}

export interface NamedOption {
  id: number
  name: string
}

export interface MasterPayload {
  name: string
  sort_order: number
}

// ─── Kelompok tambahan ────────────────────────────────────────────────────────

export interface ModifierOption {
  id: number
  name: string
  price: number
  is_default: boolean
  is_active: boolean
}

export interface ModifierGroup {
  id: number
  name: string
  selection_type: SelectionType
  is_required: boolean
  is_active: boolean
  menu_count: number
  options: ModifierOption[]
}

export interface ModifierOptionInput {
  id: number // 0 = pilihan baru
  name: string
  price: number
  is_default: boolean
  is_active: boolean
}

export interface ModifierGroupPayload {
  name: string
  selection_type: SelectionType
  is_required: boolean
  options: ModifierOptionInput[]
}

// ─── Menu ─────────────────────────────────────────────────────────────────────

export interface MenuItem {
  id: number
  category_id: number
  category_name: string
  name: string
  stock_mode: StockMode
  /** Produk barang jadi yang ditautkan (hanya mode barang jadi, keputusan 91). */
  product_id: number | null
  product_name: string
  has_variants: boolean
  is_active: boolean
  is_sold_out: boolean
  price_min: number
  price_max: number
  variant_count: number
}

export interface RecipeLine {
  ingredient_id: number
  ingredient_name?: string
  amount: string
}

export interface MenuVariant {
  id: number
  size_id: number | null
  size_name: string
  temperature: Temperature | null
  price: number
  is_active: boolean
  recipe: RecipeLine[]
}

export interface MenuItemDetail extends MenuItem {
  variants: MenuVariant[]
  modifier_groups: NamedOption[]
}

export interface MenuItemListFilter {
  page: number
  limit: number
  search: string
  category_id?: number
  stock_mode?: StockMode
  is_active?: boolean
  is_sold_out?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface VariantInput {
  size_id: number | null
  temperature: Temperature | null
  price: number
  is_active: boolean
  recipe: RecipeLine[]
}

export interface MenuItemPayload {
  category_id: number
  name: string
  stock_mode: StockMode
  product_id?: number | null
  has_variants: boolean
  price?: number
  recipe?: RecipeLine[]
  variants?: VariantInput[]
  modifier_group_ids: number[]
}
