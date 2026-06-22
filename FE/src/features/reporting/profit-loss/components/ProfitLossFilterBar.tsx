import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { monthStart, todayStr } from '@/shared/utils'

interface DateFilter {
  date_from?: string
  date_to?: string
}

interface ProfitLossFilterBarProps {
  filter: DateFilter
  onChange: (filter: DateFilter) => void
}

export function ProfitLossFilterBar({ filter, onChange }: ProfitLossFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Dari</Label>
        <Input
          type="date"
          value={filter.date_from ?? ''}
          onChange={(e) => onChange({ ...filter, date_from: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Sampai</Label>
        <Input
          type="date"
          value={filter.date_to ?? ''}
          onChange={(e) => onChange({ ...filter, date_to: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <Button
        variant="outline"
        size="sm"
        className="h-9"
        onClick={() => onChange({ date_from: monthStart(), date_to: todayStr() })}
      >
        Bulan ini
      </Button>
      <Button
        variant="outline"
        size="sm"
        onClick={() => onChange({ date_from: monthStart(), date_to: todayStr() })}
        className="h-9 gap-1"
      >
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
