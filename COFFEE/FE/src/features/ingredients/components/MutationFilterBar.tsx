import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Combobox } from '@/shared/components/ui/combobox'
import { Input } from '@/shared/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'

import { useIngredientOptionsQuery } from '../ingredients.api'
import type { MutationListFilter, MutationType } from '../ingredients.types'
import { MUTATION_TYPE_LABEL } from '../ingredients.utils'

interface MutationFilterBarProps {
  filter: MutationListFilter
  onChange: (filter: MutationListFilter) => void
  onReset: () => void
}

export function MutationFilterBar({ filter, onChange, onReset }: MutationFilterBarProps) {
  const { data: options = [] } = useIngredientOptionsQuery()
  const ingredientOptions = options.map((o) => ({ value: String(o.id), label: o.name }))

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="min-w-[200px] flex-1 space-y-1">
        <span className="text-xs text-gray-500">Bahan</span>
        <Combobox
          options={ingredientOptions}
          value={filter.ingredient_id ? String(filter.ingredient_id) : undefined}
          onValueChange={(v) => onChange({ ...filter, ingredient_id: v ? Number(v) : undefined })}
          placeholder="Semua bahan"
          searchPlaceholder="Cari bahan..."
          emptyText="Bahan tidak ditemukan."
          className="w-full"
        />
      </div>
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Jenis</span>
        <Select
          value={filter.type ?? 'all'}
          onValueChange={(v) => onChange({ ...filter, type: v === 'all' ? undefined : (v as MutationType) })}
        >
          <SelectTrigger className="w-40 h-9" aria-label="Jenis riwayat">
            <SelectValue placeholder="Semua jenis" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua jenis</SelectItem>
            {(Object.keys(MUTATION_TYPE_LABEL) as MutationType[]).map((t) => (
              <SelectItem key={t} value={t}>
                {MUTATION_TYPE_LABEL[t]}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Dari</span>
        <Input
          type="date"
          aria-label="Dari tanggal"
          value={filter.start_date ?? ''}
          onChange={(e) => onChange({ ...filter, start_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Sampai</span>
        <Input
          type="date"
          aria-label="Sampai tanggal"
          value={filter.end_date ?? ''}
          onChange={(e) => onChange({ ...filter, end_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <Button variant="outline" size="sm" onClick={onReset} className="h-9 gap-1">
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
