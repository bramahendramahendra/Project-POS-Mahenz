import { useEffect, useState } from 'react'
import { Search } from 'lucide-react'

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

import type { ReceivableListFilter, ReceivableStatus } from '../receivables.types'

interface ReceivableFilterBarProps {
  filter: ReceivableListFilter
  onChange: (filter: ReceivableListFilter) => void
  onReset: () => void
}

export function ReceivableFilterBar({ filter, onChange, onReset }: ReceivableFilterBarProps) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 400)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch || undefined })
    onReset()
  }, [debouncedSearch])

  const handleStatusChange = (value: string) => {
    onChange({ ...filter, status: value === 'all' ? undefined : (value as ReceivableStatus) })
    onReset()
  }

  const handleReset = () => {
    setSearch('')
    onChange({ ...filter, search: undefined, status: undefined })
    onReset()
  }

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="relative min-w-[220px] flex-1">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari kode transaksi / pelanggan..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-8 h-9 text-sm"
        />
      </div>
      <Select
        value={filter.status ?? 'all'}
        onValueChange={handleStatusChange}
      >
        <SelectTrigger className="w-44 h-9">
          <SelectValue placeholder="Semua Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Semua</SelectItem>
          <SelectItem value="unpaid">Belum Lunas</SelectItem>
          <SelectItem value="partial">Sebagian</SelectItem>
          <SelectItem value="paid">Lunas</SelectItem>
        </SelectContent>
      </Select>
      <Button variant="outline" size="sm" className="h-9" onClick={handleReset}>
        Reset
      </Button>
    </div>
  )
}
