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

import type { StockFilter } from '../stock.types'

interface StockReportFilterBarProps {
  filter: StockFilter
  onChange: (filter: StockFilter) => void
  onReset: () => void
}

export function StockReportFilterBar({ filter, onChange, onReset }: StockReportFilterBarProps) {
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

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="relative">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari produk..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8 h-9 w-52"
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
          <SelectTrigger className="w-40 h-9">
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
      <Button variant="outline" size="sm" className="h-9 gap-1" onClick={handleReset}>
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
