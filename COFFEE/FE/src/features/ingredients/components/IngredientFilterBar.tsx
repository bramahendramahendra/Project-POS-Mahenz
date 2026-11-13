import { useEffect, useState } from 'react'
import { RotateCcw, Search } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Switch } from '@/shared/components/ui/switch'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'
import { useDebounce } from '@/shared/hooks'

import type { IngredientListFilter } from '../ingredients.types'

interface IngredientFilterBarProps {
  filter: IngredientListFilter
  onChange: (filter: IngredientListFilter) => void
  onReset: () => void
}

export function IngredientFilterBar({ filter, onChange, onReset }: IngredientFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch ?? '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch])

  return (
    <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
      <div className="relative min-w-[220px] flex-1">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari bahan..."
          aria-label="Cari bahan"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8 h-9 text-sm"
        />
      </div>

      <ToggleGroup
        type="single"
        value={filter.is_active === undefined ? 'all' : filter.is_active ? 'active' : 'inactive'}
        onValueChange={(v) => v && onChange({ ...filter, is_active: v === 'all' ? undefined : v === 'active' })}
      >
        <ToggleGroupItem value="all">Semua</ToggleGroupItem>
        <ToggleGroupItem value="active">Aktif</ToggleGroupItem>
        <ToggleGroupItem value="inactive">Nonaktif</ToggleGroupItem>
      </ToggleGroup>

      <label className="flex items-center gap-2 text-sm text-gray-700">
        <Switch
          checked={!!filter.low_stock}
          onCheckedChange={(checked) => onChange({ ...filter, low_stock: checked || undefined })}
          aria-label="Hanya stok menipis"
        />
        Stok menipis
      </label>

      <Button
        variant="outline"
        size="sm"
        onClick={() => {
          setSearch('')
          onReset()
        }}
        className="h-9 gap-1"
      >
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
