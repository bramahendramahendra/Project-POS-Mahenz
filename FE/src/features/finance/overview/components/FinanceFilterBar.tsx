import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'

import type { FinanceFilter } from '../finance.types'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function weekStartString(): string {
  const d = new Date()
  d.setDate(d.getDate() - d.getDay() + 1)
  return d.toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

interface FinanceFilterBarProps {
  filter: FinanceFilter
  onChange: (filter: FinanceFilter) => void
  onReset: () => void
}

export function FinanceFilterBar({ filter, onChange, onReset }: FinanceFilterBarProps) {
  const applyPreset = (from: string, to: string) => {
    onChange({ ...filter, date_from: from, date_to: to })
    onReset()
  }

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Dari</label>
        <Input
          type="date"
          value={filter.date_from ?? ''}
          onChange={(e) => { onChange({ ...filter, date_from: e.target.value || undefined }); onReset() }}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Sampai</label>
        <Input
          type="date"
          value={filter.date_to ?? ''}
          onChange={(e) => { onChange({ ...filter, date_to: e.target.value || undefined }); onReset() }}
          className="w-40 h-9"
        />
      </div>
      <div className="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          className="h-9"
          onClick={() => applyPreset(todayString(), todayString())}
        >
          Hari ini
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="h-9"
          onClick={() => applyPreset(weekStartString(), todayString())}
        >
          Minggu ini
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="h-9"
          onClick={() => applyPreset(monthStartString(), todayString())}
        >
          Bulan ini
        </Button>
      </div>
    </div>
  )
}
