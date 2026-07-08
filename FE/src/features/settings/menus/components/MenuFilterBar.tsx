import { useEffect, useState } from 'react'
import { Search } from 'lucide-react'

import { Input } from '@/shared/components/ui/input'
import { useDebounce } from '@/shared/hooks'

import type { MenuListFilter } from '@/features/menu/menu.types'

interface MenuFilterBarProps {
  filter: MenuListFilter
  onChange: (filter: MenuListFilter) => void
}

export function MenuFilterBar({ filter, onChange }: MenuFilterBarProps) {
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
          placeholder="Cari menu..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8"
        />
      </div>
    </div>
  )
}
