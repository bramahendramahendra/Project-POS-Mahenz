import { useEffect, useState } from 'react'
import { Search } from 'lucide-react'

import { Input } from '@/shared/components/ui/input'
import { useDebounce } from '@/shared/hooks'

import type { RoleListFilter } from '../roles.types'

interface RoleFilterBarProps {
  filter: RoleListFilter
  onChange: (filter: RoleListFilter) => void
}

export function RoleFilterBar({ filter, onChange }: RoleFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch ?? '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch])

  return (
    <div className="flex items-center gap-3">
      <div className="relative min-w-[220px] max-w-xs">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari role..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8"
        />
      </div>
    </div>
  )
}
