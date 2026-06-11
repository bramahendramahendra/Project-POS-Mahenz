import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'

import type { CashDrawerListFilter } from '../cash-drawer.types'

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

interface CashDrawerFilterBarProps {
  filter: CashDrawerListFilter
  onChange: (filter: CashDrawerListFilter) => void
}

export function CashDrawerFilterBar({ filter, onChange }: CashDrawerFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Dari</label>
        <Input
          type="date"
          value={filter.start_date ?? ''}
          onChange={(e) => onChange({ ...filter, start_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Sampai</label>
        <Input
          type="date"
          value={filter.end_date ?? ''}
          onChange={(e) => onChange({ ...filter, end_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <Button
        variant="outline"
        size="sm"
        className="h-9"
        onClick={() =>
          onChange({ ...filter, start_date: monthStartString(), end_date: todayString() })
        }
      >
        Bulan ini
      </Button>
    </div>
  )
}
