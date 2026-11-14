import { useState } from 'react'

import { AsyncCombobox } from '@/shared/components/ui/async-combobox'
import { Combobox } from '@/shared/components/ui/combobox'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useProductSearchQuery } from '@/features/products/products'
import type { ProductSearchOption } from '@/features/products/products/products.types'
import { useIngredientOptionsQuery, type IngredientOption } from '@/features/ingredients'
import { PURCHASE_ITEM_TYPE_LABEL, type PurchaseItemType } from '../purchases.types'

interface PurchaseItemTargetCellProps {
  itemType: PurchaseItemType
  productId: number
  productName?: string
  ingredientId?: number
  ingredientName?: string
  error?: boolean
  onTypeChange: (type: PurchaseItemType) => void
  onProductChange: (id: number, name: string) => void
  onIngredientChange: (ingredient: IngredientOption) => void
}

function ProductPicker({
  value,
  label,
  onChange,
  error,
}: {
  value: number
  label?: string
  onChange: (id: number, name: string) => void
  error?: boolean
}) {
  const [keyword, setKeyword] = useState('')
  const { data: results, isFetching } = useProductSearchQuery(keyword)

  return (
    <AsyncCombobox<ProductSearchOption>
      className={`h-8 w-full text-xs ${error ? 'border-red-500' : ''}`}
      value={value > 0 ? value : undefined}
      selectedLabel={label}
      onSearch={setKeyword}
      onValueChange={(_v, item) => {
        if (item) onChange(item.id, item.name)
      }}
      options={results}
      getOptionValue={(p) => p.id}
      getOptionLabel={(p) => p.name}
      isLoading={isFetching}
      placeholder="Pilih barang jadi"
      searchPlaceholder="Cari barang jadi..."
      emptyText="Barang jadi tidak ditemukan."
    />
  )
}

function IngredientPicker({
  value,
  label,
  onChange,
  error,
}: {
  value: number
  label?: string
  onChange: (ingredient: IngredientOption) => void
  error?: boolean
}) {
  const { data: ingredients = [] } = useIngredientOptionsQuery()
  const options = ingredients.map((i) => ({ value: String(i.id), label: `${i.name} (${i.unit_abbr})` }))
  // Bahan yang sudah nonaktif tidak ada di pilihan; tetap tampilkan namanya pada baris lama.
  if (value > 0 && !ingredients.some((i) => i.id === value)) {
    options.unshift({ value: String(value), label: label || `Bahan #${value}` })
  }

  return (
    <Combobox
      className={`h-8 w-full text-xs ${error ? 'border-red-500' : ''}`}
      options={options}
      value={value > 0 ? String(value) : undefined}
      onValueChange={(v) => {
        const picked = ingredients.find((i) => String(i.id) === v)
        if (picked) onChange(picked)
      }}
      placeholder="Pilih bahan baku"
      searchPlaceholder="Cari bahan baku..."
      emptyText="Bahan baku tidak ditemukan."
    />
  )
}

/**
 * Kolom pertama baris pembelian: jenis baris (barang jadi / bahan baku, keputusan 92)
 * lalu pilihan barang atau bahan sesuai jenisnya.
 */
export function PurchaseItemTargetCell({
  itemType,
  productId,
  productName,
  ingredientId,
  ingredientName,
  error,
  onTypeChange,
  onProductChange,
  onIngredientChange,
}: PurchaseItemTargetCellProps) {
  return (
    <div className="flex flex-col gap-1">
      <Select value={itemType} onValueChange={(v) => onTypeChange(v as PurchaseItemType)}>
        <SelectTrigger className="h-7 w-[130px] text-[11px]" aria-label="Jenis baris">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {(Object.keys(PURCHASE_ITEM_TYPE_LABEL) as PurchaseItemType[]).map((t) => (
            <SelectItem key={t} value={t}>
              {PURCHASE_ITEM_TYPE_LABEL[t]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {itemType === 'ingredient' ? (
        <IngredientPicker value={ingredientId ?? 0} label={ingredientName} onChange={onIngredientChange} error={error} />
      ) : (
        <ProductPicker value={productId} label={productName} onChange={onProductChange} error={error} />
      )}
    </div>
  )
}
