export type DiscountType = 'none' | 'percent' | 'amount'
export type PaymentMethod = 'cash' | 'transfer' | 'qris' | 'card' | 'kredit' | 'balance'

// Nilai enum yang diterima backend (sesuai validasi oneof)
export const PAYMENT_METHOD_VALUES = ['cash', 'transfer', 'qris', 'card', 'kredit'] as const

// ─── Katalog kasir (GET /menu-items/catalog, tahap 4.7b) ──────────────────────

export type Temperature = 'hot' | 'iced'
export type OrderType = 'take_away' | 'dine_in'

export const ORDER_TYPE_LABEL: Record<OrderType, string> = {
  take_away: 'Bawa pulang',
  dine_in: 'Makan di tempat',
}

// Label metode bayar (keputusan 135: kode `kredit` ditampilkan "Hutang" di semua tempat). Satu sumber untuk kasir,
// struk, riwayat transaksi, laporan, dan pengeluaran.
export const PAYMENT_METHOD_LABEL: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Hutang',
  balance: 'Saldo',
}

export interface CatalogVariant {
  id: number
  size_id: number | null
  size_name: string
  temperature: Temperature | null
  price: number
}

export interface CatalogModifierOption {
  id: number
  name: string
  price: number
  is_default: boolean
}

export interface CatalogModifierGroup {
  id: number
  name: string
  selection_type: 'single' | 'multiple'
  is_required: boolean
  options: CatalogModifierOption[]
}

export interface CatalogItem {
  id: number
  name: string
  category_id: number
  stock_mode: 'recipe' | 'direct'
  has_variants: boolean
  is_sold_out: boolean
  /** Sisa stok barang jadi; null untuk racikan. */
  stock: number | null
  variants: CatalogVariant[]
  modifier_groups: CatalogModifierGroup[]
}

export interface MenuCatalog {
  categories: { id: number; name: string }[]
  items: CatalogItem[]
}

// ─── Keranjang ────────────────────────────────────────────────────────────────

export interface CartModifier {
  option_id: number
  group_name: string
  option_name: string
  price: number
}

/**
 * Satu baris pesanan. Baris yang sama (kombinasi, tambahan, catatan sama) digabung jumlahnya.
 * product_name / unit_name dipertahankan untuk struk (nama menu; satuan kosong di Coffee).
 */
export interface CartItem {
  key: string
  menu_item_id: number
  variant_id: number
  product_name: string // nama menu
  variant_name: string // mis. "Regular · Dingin", kosong bila menu tanpa varian
  unit_name: string
  modifiers: CartModifier[]
  notes?: string
  is_stock_item: boolean // barang jadi (stok dipotong)
  qty: number
  price: number // harga satuan = kombinasi + tambahan
  subtotal: number
  // Diskon per item (seperti Toko)
  discount_type?: 'percent' | 'nominal'
  discount_value?: number
  discount_amount?: number
  effective_price?: number
}

export interface Discount {
  type: DiscountType
  value: number
  amount: number
}

export interface Tax {
  percent: number
  amount: number
}

export interface CartSummary {
  subtotal: number
  discountAmount: number
  taxAmount: number
  grandTotal: number
}

/** Info pesanan yang dicetak di struk (nomor antrean besar, jenis, meja, pemesan). */
export interface OrderInfo {
  queue_label?: string
  order_type?: OrderType | ''
  table_name?: string
  order_name?: string
}

export interface PaymentPayload {
  order_type: OrderType
  table_id?: number
  order_name?: string
  customer_id?: number
  shift_id?: number
  is_credit: boolean
  device_source: 'web'
  items: Array<{
    variant_id: number
    quantity: number
    modifier_option_ids: number[]
    notes?: string
    discount_item?: number
  }>
  discount: number
  tax: number
  payment_method: PaymentMethod
  payment_amount: number
}

export interface CheckoutResponse extends OrderInfo {
  id: number
  transaction_code: string
  total_amount: number
  payment_amount: number
  change_amount: number
  transaction_date: string
}
