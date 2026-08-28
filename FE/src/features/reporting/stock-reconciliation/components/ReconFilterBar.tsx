import { useEffect, useState } from 'react'
import { RotateCcw, Search } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useDebounce } from '@/shared/hooks'
import { useCategoryOptionsQuery } from '@/features/products/categories'

import type { ReconFilter } from '../stock-reconciliation.types'

interface ReconFilterBarProps {
  filter: ReconFilter
  onChange: (filter: ReconFilter) => void
  onReset: () => void
}

export function ReconFilterBar({ filter, onChange, onReset }: ReconFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  const { data: categories = [] } = useCategoryOptionsQuery()

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch || undefined })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch])

  const handleReset = () => {
    setSearch('')
    onReset()
  }

  const onlyReview = filter.only_review !== false // default true

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="relative">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari produk..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="h-9 w-52 pl-8"
        />
      </div>

      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Kategori</Label>
        <Select
          value={filter.category_id ? String(filter.category_id) : 'all'}
          onValueChange={(v) =>
            onChange({ ...filter, category_id: v === 'all' ? undefined : Number(v) })
          }
        >
          <SelectTrigger className="h-9 w-40">
            <SelectValue placeholder="Semua Kategori" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Kategori</SelectItem>
            {categories.map((cat) => (
              <SelectItem key={cat.id} value={String(cat.id)}>
                {cat.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Tampilkan</Label>
        <div className="inline-flex h-9 overflow-hidden rounded-md border">
          <button
            type="button"
            onClick={() => onChange({ ...filter, only_review: true })}
            className={`px-3 text-sm ${onlyReview ? 'bg-primary text-primary-foreground' : 'bg-white text-gray-600'}`}
          >
            Perlu Ditinjau
          </button>
          <button
            type="button"
            onClick={() => onChange({ ...filter, only_review: false })}
            className={`px-3 text-sm ${!onlyReview ? 'bg-primary text-primary-foreground' : 'bg-white text-gray-600'}`}
          >
            Semua
          </button>
        </div>
      </div>

      <Button variant="outline" size="sm" className="h-9 gap-1" onClick={handleReset}>
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
