import { useEffect, useState } from 'react'
import { RotateCcw, Search } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { useDebounce } from '@/shared/hooks'

import type { UnitListFilter } from '../units.types'

interface UnitFilterBarProps {
  filter: UnitListFilter
  onChange: (filter: UnitListFilter) => void
  onReset: () => void
}

export function UnitFilterBar({ filter, onChange, onReset }: UnitFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch ?? '' })
  }, [debouncedSearch])

  const hasFilter = !!search

  return (
    <div className="flex items-center gap-2 rounded-lg border bg-white p-3">
      <div className="relative min-w-[220px] flex-1">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari satuan..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8 h-9 text-sm"
        />
      </div>
      {hasFilter && (
        <Button
          variant="outline"
          size="sm"
          onClick={() => { setSearch(''); onReset() }}
          className="h-9 gap-1"
        >
          <RotateCcw size={13} />
          Reset
        </Button>
      )}
    </div>
  )
}
