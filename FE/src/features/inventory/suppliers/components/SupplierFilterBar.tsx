import { useEffect, useState } from 'react'
import { RotateCcw, Search } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useDebounce } from '@/shared/hooks'

import type { SupplierListFilter } from '../suppliers.types'

interface SupplierFilterBarProps {
  filter: SupplierListFilter
  onChange: (filter: SupplierListFilter) => void
  onReset: () => void
}

export function SupplierFilterBar({ filter, onChange, onReset }: SupplierFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch ?? '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch])

  return (
    <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
      {/* Search */}
      <div className="relative min-w-[220px] flex-1">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari supplier..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8 h-9 text-sm"
        />
      </div>

      {/* Status */}
      <Select
        value={filter.is_active === undefined ? 'all' : filter.is_active ? 'active' : 'inactive'}
        onValueChange={(v) =>
          onChange({
            ...filter,
            is_active: v === 'all' ? undefined : v === 'active',
          })
        }
      >
        <SelectTrigger className="h-9 w-[160px] text-sm">
          <SelectValue placeholder="Semua Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Semua Status</SelectItem>
          <SelectItem value="active">Aktif</SelectItem>
          <SelectItem value="inactive">Nonaktif</SelectItem>
        </SelectContent>
      </Select>

      {/* Reset */}
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
