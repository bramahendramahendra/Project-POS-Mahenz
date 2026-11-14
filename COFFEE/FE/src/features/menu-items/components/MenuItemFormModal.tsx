import { Fragment, useMemo, useState } from 'react'
import { ChevronDown, ChevronRight } from 'lucide-react'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Checkbox } from '@/shared/components/ui/checkbox'
import { Combobox } from '@/shared/components/ui/combobox'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'
import { Switch } from '@/shared/components/ui/switch'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'
import { useIngredientOptionsQuery } from '@/features/ingredients/ingredients.api'
import { useProductOptionsQuery } from '@/features/products/products'
import { formatRupiah } from '@/shared/utils'

import {
  menuCategoryApi,
  menuSizeApi,
  useMenuItemDetailQuery,
  useModifierGroupOptionsQuery,
  useSaveMenuItemMutation,
} from '../menu-items.api'
import type { MenuItemDetail, MenuItemPayload, RecipeLine, StockMode, Temperature, VariantInput } from '../menu-items.types'
import { STOCK_MODE_LABEL, TEMPERATURE_LABEL, TEMPERATURES, variantKey, variantLabel } from '../menu-items.utils'
import { RecipeEditor } from './RecipeEditor'

interface MenuItemFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  itemId: number | null // null = menu baru
}

// Pembungkus: saat edit, ambil detail dulu lalu tampilkan form dengan nilai awal dari detail.
export function MenuItemFormModal({ open, onOpenChange, itemId }: MenuItemFormModalProps) {
  const { data: detail, isLoading } = useMenuItemDetailQuery(open ? itemId : null)
  if (!open) return null
  if (itemId != null && (isLoading || !detail)) {
    return (
      <FormModal open onOpenChange={onOpenChange} title="Edit Menu" size="xl" isLoading hideFooter>
        <p className="py-6 text-center text-sm text-gray-500">Memuat data menu...</p>
      </FormModal>
    )
  }
  return <MenuItemForm key={detail?.id ?? 'new'} detail={itemId != null ? (detail ?? null) : null} onOpenChange={onOpenChange} />
}

interface ComboState {
  price: number
  is_active: boolean
  recipe: RecipeLine[]
}

function MenuItemForm({ detail, onOpenChange }: { detail: MenuItemDetail | null; onOpenChange: (open: boolean) => void }) {
  const { data: categories = [] } = menuCategoryApi.useOptions()
  const { data: sizes = [] } = menuSizeApi.useOptions()
  const { data: groups = [] } = useModifierGroupOptionsQuery()
  const { data: ingredients = [] } = useIngredientOptionsQuery()
  const { data: products = [] } = useProductOptionsQuery()
  const { mutate: save, isPending } = useSaveMenuItemMutation()

  const single = detail && !detail.has_variants ? detail.variants[0] : undefined
  const [name, setName] = useState(detail?.name ?? '')
  const [categoryId, setCategoryId] = useState<number | undefined>(detail?.category_id)
  const [stockMode, setStockMode] = useState<StockMode>(detail?.stock_mode ?? 'recipe')
  const [hasVariants, setHasVariants] = useState(detail?.has_variants ?? false)
  const [productId, setProductId] = useState<number | undefined>(detail?.product_id ?? undefined)
  const [price, setPrice] = useState(single?.price ?? 0)
  const [recipe, setRecipe] = useState<RecipeLine[]>(single?.recipe ?? [])
  const [sizeIds, setSizeIds] = useState<number[]>(() =>
    detail?.has_variants ? [...new Set(detail.variants.map((v) => v.size_id).filter((x): x is number => x != null))] : []
  )
  const [temps, setTemps] = useState<Temperature[]>(() =>
    detail?.has_variants ? TEMPERATURES.filter((t) => detail.variants.some((v) => v.temperature === t)) : []
  )
  const [combos, setCombos] = useState<Record<string, ComboState>>(() => {
    const out: Record<string, ComboState> = {}
    if (detail?.has_variants) {
      for (const v of detail.variants) out[variantKey(v.size_id, v.temperature)] = { price: v.price, is_active: v.is_active, recipe: v.recipe }
    }
    return out
  })
  const [groupIds, setGroupIds] = useState<number[]>(detail?.modifier_groups.map((g) => g.id) ?? [])
  const [openRecipe, setOpenRecipe] = useState<string | null>(null)
  const [error, setError] = useState('')
  const [confirming, setConfirming] = useState<MenuItemPayload | null>(null)

  const isDirect = stockMode === 'direct'

  // Kombinasi = ukuran terpilih (urutan daftar ukuran) x suhu terpilih. Ukuran/suhu boleh kosong
  // salah satunya (mis. teh hanya panas/dingin), tetapi tidak keduanya.
  const comboList = useMemo(() => {
    const orderedSizes = sizes.filter((s) => sizeIds.includes(s.id))
    const sizeAxis: ({ id: number; name: string } | null)[] = orderedSizes.length > 0 ? orderedSizes : [null]
    const tempAxis: (Temperature | null)[] = temps.length > 0 ? TEMPERATURES.filter((t) => temps.includes(t)) : [null]
    const list: { key: string; sizeId: number | null; sizeName: string; temperature: Temperature | null }[] = []
    for (const s of sizeAxis) {
      for (const t of tempAxis) {
        if (s == null && t == null) continue
        list.push({ key: variantKey(s?.id ?? null, t), sizeId: s?.id ?? null, sizeName: s?.name ?? '', temperature: t })
      }
    }
    return list
  }, [sizes, sizeIds, temps])

  const comboOf = (key: string): ComboState => combos[key] ?? { price: 0, is_active: true, recipe: [] }
  const updateCombo = (key: string, patch: Partial<ComboState>) =>
    setCombos((prev) => ({ ...prev, [key]: { ...comboOf(key), ...patch } }))

  const toggleIn = <T,>(list: T[], value: T, on: boolean) => (on ? [...list, value] : list.filter((x) => x !== value))

  const buildPayload = (): MenuItemPayload | string => {
    if (!name.trim()) return 'Nama menu wajib diisi'
    if (!categoryId) return 'Kategori wajib dipilih'
    const cleanRecipe = (lines: RecipeLine[]) => lines.map((l) => ({ ingredient_id: l.ingredient_id, amount: l.amount.trim() }))
    const recipeError = (lines: RecipeLine[]) =>
      lines.some((l) => !l.ingredient_id) ? 'Pilih bahan di setiap baris resep' : lines.some((l) => !l.amount.trim()) ? 'Takaran resep wajib diisi' : ''
    const activeGroupIds = groupIds.filter((id) => groups.some((g) => g.id === id))

    if (!hasVariants || isDirect) {
      if (isDirect && !productId) return 'Pilih produk barang jadi'
      if (price <= 0) return 'Harga wajib diisi'
      const lines = isDirect ? [] : recipe
      const re = recipeError(lines)
      if (re) return re
      return { name: name.trim(), category_id: categoryId, stock_mode: stockMode, product_id: isDirect ? productId : null, has_variants: false, price, recipe: cleanRecipe(lines), modifier_group_ids: activeGroupIds }
    }

    if (comboList.length === 0) return 'Pilih minimal satu ukuran atau suhu'
    const variants: VariantInput[] = comboList.map((c) => {
      const st = comboOf(c.key)
      return { size_id: c.sizeId, temperature: c.temperature, price: st.price, is_active: st.is_active, recipe: cleanRecipe(st.recipe) }
    })
    if (!variants.some((v) => v.is_active)) return 'Minimal satu kombinasi dijual'
    const noPrice = comboList.find((c) => comboOf(c.key).is_active && comboOf(c.key).price <= 0)
    if (noPrice) return `Harga ${variantLabel(noPrice.sizeName, noPrice.temperature)} wajib diisi`
    for (const c of comboList) {
      const re = recipeError(comboOf(c.key).recipe)
      if (re) return `${re} (${variantLabel(c.sizeName, c.temperature)})`
    }
    return { name: name.trim(), category_id: categoryId, stock_mode: stockMode, has_variants: true, variants, modifier_group_ids: activeGroupIds }
  }

  const submit = () => {
    const payload = buildPayload()
    if (typeof payload === 'string') return setError(payload)
    setError('')
    setConfirming(payload)
  }

  return (
    <>
      <FormModal
        open={confirming == null}
        onOpenChange={(v) => {
          if (!v) onOpenChange(false)
        }}
        title={detail ? 'Edit Menu' : 'Tambah Menu'}
        size="xl"
        isLoading={isPending}
        onSubmit={submit}
        submitLabel="Simpan"
      >
        <div className="space-y-5">
          {/* Info dasar */}
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="menu-name">
                Nama Menu <span className="text-red-500">*</span>
              </Label>
              <Input id="menu-name" value={name} onChange={(e) => setName(e.target.value)} placeholder="Contoh: Caffe Latte" />
            </div>
            <div className="space-y-1.5">
              <Label>
                Kategori <span className="text-red-500">*</span>
              </Label>
              <Select value={categoryId ? String(categoryId) : ''} onValueChange={(v) => setCategoryId(Number(v))}>
                <SelectTrigger className="w-full" aria-label="Kategori menu">
                  <SelectValue placeholder="Pilih kategori..." />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((c) => (
                    <SelectItem key={c.id} value={String(c.id)}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="flex flex-wrap items-end gap-6">
            <div className="space-y-1.5">
              <Label>Mode Stok</Label>
              <ToggleGroup
                type="single"
                value={stockMode}
                onValueChange={(v) => {
                  if (!v) return
                  setStockMode(v as StockMode)
                  if (v === 'direct') setHasVariants(false)
                }}
              >
                <ToggleGroupItem value="recipe">{STOCK_MODE_LABEL.recipe}</ToggleGroupItem>
                <ToggleGroupItem value="direct">{STOCK_MODE_LABEL.direct}</ToggleGroupItem>
              </ToggleGroup>
            </div>
            <label className={`flex items-center gap-2 text-sm ${isDirect ? 'text-gray-400' : 'text-gray-700'}`}>
              <Switch checked={hasVariants && !isDirect} disabled={isDirect} onCheckedChange={setHasVariants} aria-label="Punya varian" />
              Punya varian (ukuran / panas-dingin)
            </label>
          </div>
          <p className="-mt-3 text-xs text-gray-500">
            {isDirect
              ? 'Barang jadi (mis. roti, snack): ditautkan ke satu produk di menu Barang Jadi, tanpa varian dan resep. Stok produk bertambah lewat Pembelian.'
              : 'Racikan: dibuat dari bahan. Stok bahan dicatat lewat menu Bahan Baku; resep hanya catatan takaran.'}
          </p>

          {isDirect && (
            <div className="max-w-md space-y-1.5">
              <Label>
                Produk Barang Jadi <span className="text-red-500">*</span>
              </Label>
              <Combobox
                className="w-full"
                options={(() => {
                  const opts = products.map((p) => ({ value: String(p.id), label: p.name }))
                  // Produk yang sudah nonaktif tetap tampil namanya pada menu lama.
                  if (productId && !products.some((p) => p.id === productId)) {
                    opts.unshift({ value: String(productId), label: detail?.product_name || `Produk #${productId}` })
                  }
                  return opts
                })()}
                value={productId ? String(productId) : undefined}
                onValueChange={(v) => setProductId(v ? Number(v) : undefined)}
                placeholder="Pilih produk..."
                searchPlaceholder="Cari produk..."
                emptyText="Produk tidak ditemukan. Tambahkan di menu Barang Jadi."
              />
              <p className="text-xs text-gray-500">Satu produk hanya untuk satu menu. Harga jual diatur di sini, bukan di produk.</p>
            </div>
          )}

          {/* Harga */}
          {!hasVariants || isDirect ? (
            <div className="space-y-3">
              <div className="max-w-xs space-y-1.5">
                <Label htmlFor="menu-price">
                  Harga <span className="text-red-500">*</span>
                </Label>
                <RupiahInput id="menu-price" value={price} onChange={setPrice} placeholder="0" />
              </div>
              {!isDirect && (
                <div className="space-y-1.5">
                  <Label>Resep (opsional)</Label>
                  <RecipeEditor lines={recipe} onChange={setRecipe} ingredients={ingredients} label="menu" />
                </div>
              )}
            </div>
          ) : (
            <div className="space-y-3">
              <div className="flex flex-wrap gap-6">
                <div className="space-y-1.5">
                  <Label>Ukuran</Label>
                  <div className="flex flex-wrap gap-3">
                    {sizes.length === 0 && <p className="text-xs text-amber-600">Belum ada ukuran. Tambahkan di menu Ukuran.</p>}
                    {sizes.map((s) => (
                      <label key={s.id} className="flex items-center gap-2 text-sm text-gray-700">
                        <Checkbox checked={sizeIds.includes(s.id)} onCheckedChange={(v) => setSizeIds((prev) => toggleIn(prev, s.id, v === true))} aria-label={`Ukuran ${s.name}`} />
                        {s.name}
                      </label>
                    ))}
                  </div>
                </div>
                <div className="space-y-1.5">
                  <Label>Suhu</Label>
                  <div className="flex gap-3">
                    {TEMPERATURES.map((t) => (
                      <label key={t} className="flex items-center gap-2 text-sm text-gray-700">
                        <Checkbox checked={temps.includes(t)} onCheckedChange={(v) => setTemps((prev) => toggleIn(prev, t, v === true))} aria-label={`Suhu ${TEMPERATURE_LABEL[t]}`} />
                        {TEMPERATURE_LABEL[t]}
                      </label>
                    ))}
                  </div>
                </div>
              </div>

              {comboList.length > 0 && (
                <div className="overflow-x-auto rounded-md border">
                  <table className="w-full text-sm">
                    <thead className="bg-gray-50 text-left text-xs text-gray-600">
                      <tr>
                        <th className="px-3 py-2 font-medium">Kombinasi</th>
                        <th className="px-3 py-2 font-medium">Harga</th>
                        <th className="px-3 py-2 text-center font-medium">Dijual</th>
                        <th className="px-3 py-2 font-medium">Resep</th>
                      </tr>
                    </thead>
                    <tbody>
                      {comboList.map((c) => {
                        const st = comboOf(c.key)
                        const label = variantLabel(c.sizeName, c.temperature)
                        const expanded = openRecipe === c.key
                        return (
                          <Fragment key={c.key}>
                            <tr className="border-t">
                              <td className="px-3 py-2 font-medium text-gray-800">{label}</td>
                              <td className="px-3 py-2">
                                <RupiahInput value={st.price} onChange={(v) => updateCombo(c.key, { price: v })} placeholder="0" className="h-9 w-36" />
                              </td>
                              <td className="px-3 py-2 text-center">
                                <Checkbox checked={st.is_active} onCheckedChange={(v) => updateCombo(c.key, { is_active: v === true })} aria-label={`Dijual ${label}`} />
                              </td>
                              <td className="px-3 py-2">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="sm"
                                  className="gap-1"
                                  aria-label={`Resep ${label}`}
                                  onClick={() => setOpenRecipe(expanded ? null : c.key)}
                                >
                                  {expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                                  Resep ({st.recipe.length})
                                </Button>
                              </td>
                            </tr>
                            {expanded && (
                              <tr className="bg-gray-50/60">
                                <td colSpan={4} className="px-3 pb-3 pt-1">
                                  <RecipeEditor lines={st.recipe} onChange={(lines) => updateCombo(c.key, { recipe: lines })} ingredients={ingredients} label={label} />
                                </td>
                              </tr>
                            )}
                          </Fragment>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}

          {/* Kelompok tambahan */}
          <div className="space-y-1.5">
            <Label>Kelompok Tambahan</Label>
            {groups.length === 0 ? (
              <p className="text-xs text-gray-500">Belum ada kelompok tambahan aktif.</p>
            ) : (
              <div className="grid gap-2 sm:grid-cols-2">
                {groups.map((g) => (
                  <label key={g.id} className="flex items-start gap-2 rounded-md border p-2 text-sm">
                    <Checkbox
                      checked={groupIds.includes(g.id)}
                      onCheckedChange={(v) => setGroupIds((prev) => toggleIn(prev, g.id, v === true))}
                      aria-label={`Kelompok ${g.name}`}
                    />
                    <span>
                      <span className="font-medium text-gray-800">{g.name}</span>
                      <span className="block text-xs text-gray-500">
                        {g.selection_type === 'single' ? 'Pilih satu' : 'Boleh banyak'}, {g.is_required ? 'wajib' : 'opsional'}:{' '}
                        {g.options.map((o) => (o.price > 0 ? `${o.name} (+${formatRupiah(o.price)})` : o.name)).join(', ')}
                      </span>
                    </span>
                  </label>
                ))}
              </div>
            )}
          </div>

          {error && <p className="text-sm text-red-500">{error}</p>}
        </div>
      </FormModal>

      <ConfirmDialog
        open={confirming != null}
        onOpenChange={(v) => {
          if (!v) setConfirming(null)
        }}
        title={detail ? 'Update Menu' : 'Tambah Menu'}
        description={`Yakin ingin menyimpan menu "${confirming?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={() =>
          confirming &&
          save({ id: detail?.id, ...confirming }, { onSuccess: () => onOpenChange(false), onError: () => setConfirming(null) })
        }
      />
    </>
  )
}
