import { useEffect, useState } from 'react'
import { RotateCcw, Search } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { useDebounce } from '@/shared/hooks'

import type { CapitalReconFilter } from '../capital-reconciliation.types'

interface Props {
  filter: CapitalReconFilter
  onChange: (filter: CapitalReconFilter) => void
  onReset: () => void
}

const STATUS_TABS: { value: NonNullable<CapitalReconFilter['status']>; label: string }[] = [
  { value: 'pending', label: 'Perlu Ditinjau' },
  { value: 'resolved', label: 'Sudah Dikoreksi' },
  { value: 'skipped', label: 'Ditandai Selesai' },
  { value: 'all', label: 'Semua' },
]

export function CapitalReconFilterBar({ filter, onChange, onReset }: Props) {
  const [search, setSearch] = useState(filter.search ?? '')
  const debouncedSearch = useDebounce(search, 300)

  useEffect(() => {
    onChange({ ...filter, search: debouncedSearch || undefined })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch])

  const handleReset = () => {
    setSearch('')
    onReset()
  }

  const status = filter.status ?? 'pending'

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="relative">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Cari produk / kode transaksi..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="h-9 w-64 pl-8"
        />
      </div>

      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Status</Label>
        <div className="inline-flex h-9 overflow-hidden rounded-md border">
          {STATUS_TABS.map((tab) => (
            <button
              key={tab.value}
              type="button"
              onClick={() => onChange({ ...filter, status: tab.value })}
              className={`px-3 text-sm ${status === tab.value ? 'bg-primary text-primary-foreground' : 'bg-white text-gray-600'}`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <Button variant="outline" size="sm" className="h-9 gap-1" onClick={handleReset}>
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
