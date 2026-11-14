import { useState } from 'react'
import { Check, Minus, Plus } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { formatRupiah } from '@/shared/utils'

import type { CartItem, CatalogItem, CatalogModifierGroup } from '../cashier.types'
import { cartLineKey } from '../cashier.utils'
import { buildCartItem, catalogVariantName } from './menu-option.utils'

interface MenuOptionDialogProps {
  item: CatalogItem
  /** Batas jumlah untuk barang jadi (sisa stok dikurangi yang sudah di keranjang). */
  maxQty?: number
  onClose: () => void
  onAdd: (line: CartItem) => void
}

function initialSelection(groups: CatalogModifierGroup[]): Record<number, number[]> {
  const sel: Record<number, number[]> = {}
  for (const g of groups) {
    const defaults = g.options.filter((o) => o.is_default).map((o) => o.id)
    sel[g.id] = g.selection_type === 'single' ? defaults.slice(0, 1) : defaults
  }
  return sel
}

const pillClass = (active: boolean) =>
  `flex min-h-11 items-center justify-between gap-2 rounded-lg border-2 px-3 py-2 text-left text-sm transition-colors ${
    active ? 'border-[#2c3e50] bg-[#2c3e50]/5 font-semibold text-gray-900' : 'border-gray-200 text-gray-700 hover:border-gray-300'
  }`

/**
 * Pilihan menu di kasir (tahap 4.7b): kombinasi varian, tambahan (wajib/opsional, bawaan sudah
 * terpilih, keputusan 87), jumlah, dan catatan. Aturan sama dengan BE.
 */
export function MenuOptionDialog({ item, maxQty, onClose, onAdd }: MenuOptionDialogProps) {
  const [variantId, setVariantId] = useState<number>(item.variants.length === 1 ? item.variants[0].id : 0)
  const [selected, setSelected] = useState<Record<number, number[]>>(() => initialSelection(item.modifier_groups))
  const [qty, setQty] = useState(1)
  const [notes, setNotes] = useState('')

  const variant = item.variants.find((v) => v.id === variantId)
  const chosen = item.modifier_groups.flatMap((g) =>
    g.options.filter((o) => selected[g.id]?.includes(o.id)).map((o) => ({ ...o, group_name: g.name })),
  )
  const missingGroups = item.modifier_groups.filter((g) => g.is_required && (selected[g.id]?.length ?? 0) === 0)
  const unitPrice = (variant?.price ?? 0) + chosen.reduce((s, o) => s + o.price, 0)
  const limit = maxQty ?? Infinity
  const canAdd = !!variant && missingGroups.length === 0 && qty >= 1 && qty <= limit

  const toggleOption = (group: CatalogModifierGroup, optionId: number) => {
    setSelected((prev) => {
      const current = prev[group.id] ?? []
      if (group.selection_type === 'single') {
        // Kelompok wajib tidak bisa dikosongkan dengan mengetuk ulang pilihan yang sama.
        const next = current.includes(optionId) ? (group.is_required ? current : []) : [optionId]
        return { ...prev, [group.id]: next }
      }
      return { ...prev, [group.id]: current.includes(optionId) ? current.filter((id) => id !== optionId) : [...current, optionId] }
    })
  }

  const submit = () => {
    if (!variant || !canAdd) return
    const trimmed = notes.trim()
    const optionIds = chosen.map((o) => o.id)
    onAdd(buildCartItem(item, variant, chosen, qty, trimmed, cartLineKey(variant.id, optionIds, trimmed)))
  }

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title={item.name}
      description={item.stock !== null ? `Barang jadi, sisa stok ${item.stock}` : 'Pilih kombinasi dan tambahan'}
      contentClassName="max-w-lg max-h-[90vh]"
      footer={
        <div className="flex items-center gap-2 border-t px-5 py-3">
          <div className="flex items-center gap-1.5">
            <Button type="button" variant="outline" size="icon" className="h-10 w-10" aria-label="Kurangi jumlah" onClick={() => setQty((q) => Math.max(1, q - 1))}>
              <Minus size={16} />
            </Button>
            <span className="w-8 text-center text-base font-semibold" aria-label="Jumlah">{qty}</span>
            <Button type="button" variant="outline" size="icon" className="h-10 w-10" aria-label="Tambah jumlah" disabled={qty + 1 > limit} onClick={() => setQty((q) => q + 1)}>
              <Plus size={16} />
            </Button>
          </div>
          <Button type="button" className="ml-auto h-10 flex-1" disabled={!canAdd} onClick={submit}>
            Tambah · {formatRupiah(unitPrice * qty)}
          </Button>
        </div>
      }
    >
      <div className="max-h-[65vh] space-y-5 overflow-y-auto px-5 py-4">
        {item.variants.length > 1 && (
          <div className="space-y-2">
            <Label>
              Pilihan <span className="text-red-500">*</span>
            </Label>
            <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
              {item.variants.map((v) => (
                <button key={v.id} type="button" aria-pressed={variantId === v.id} onClick={() => setVariantId(v.id)} className={pillClass(variantId === v.id)}>
                  <span>{catalogVariantName(item, v)}</span>
                  <span className="shrink-0 text-xs text-gray-500">{formatRupiah(v.price)}</span>
                </button>
              ))}
            </div>
          </div>
        )}

        {item.modifier_groups.map((g) => (
          <div key={g.id} className="space-y-2">
            <Label>
              {g.name}{' '}
              <span className="text-xs font-normal text-gray-500">
                {g.is_required ? '(wajib' : '(opsional'}
                {g.selection_type === 'single' ? ', pilih satu)' : ', boleh lebih dari satu)'}
              </span>
            </Label>
            <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
              {g.options.map((o) => {
                const active = selected[g.id]?.includes(o.id) ?? false
                return (
                  <button key={o.id} type="button" aria-pressed={active} onClick={() => toggleOption(g, o.id)} className={pillClass(active)}>
                    <span className="flex items-center gap-1.5">
                      {active && <Check size={14} className="text-[#2c3e50]" />}
                      {o.name}
                    </span>
                    <span className="shrink-0 text-xs text-gray-500">{o.price > 0 ? `+${formatRupiah(o.price)}` : 'Gratis'}</span>
                  </button>
                )
              })}
            </div>
          </div>
        ))}

        <div className="space-y-1.5">
          <Label htmlFor="menu-notes">Catatan</Label>
          <Input id="menu-notes" value={notes} maxLength={200} onChange={(e) => setNotes(e.target.value)} placeholder="Contoh: es sedikit" />
        </div>

        {!variant && <p className="text-xs text-amber-600">Pilih salah satu pilihan menu.</p>}
        {missingGroups.length > 0 && (
          <p className="text-xs text-amber-600">Wajib dipilih: {missingGroups.map((g) => g.name).join(', ')}.</p>
        )}
        {qty > limit && <p className="text-xs text-red-600">Stok tidak cukup.</p>}
      </div>
    </ActionModal>
  )
}
