import { Download, RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { monthStart, todayStr, weekStart } from '@/shared/utils'

import { useExportCashierPerformanceMutation } from '../cashier-performance.api'
import type { CashierPerformanceDateFilter } from '../cashier-performance.types'

interface CashierPerformanceFilterBarProps {
  filter: CashierPerformanceDateFilter
  onChange: (filter: CashierPerformanceDateFilter) => void
  onReset: () => void
}

export function CashierPerformanceFilterBar({ filter, onChange, onReset }: CashierPerformanceFilterBarProps) {
  const { mutate: exportReport, isPending: isExporting } = useExportCashierPerformanceMutation()

  const applyPreset = (from: string, to: string) => {
    onChange({ ...filter, date_from: from, date_to: to })
  }

  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Dari</Label>
        <Input
          type="date"
          value={filter.date_from ?? ''}
          onChange={(e) => onChange({ ...filter, date_from: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Sampai</Label>
        <Input
          type="date"
          value={filter.date_to ?? ''}
          onChange={(e) => onChange({ ...filter, date_to: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="flex gap-2">
        <Button variant="outline" size="sm" className="h-9" onClick={() => applyPreset(todayStr(), todayStr())}>
          Hari ini
        </Button>
        <Button variant="outline" size="sm" className="h-9" onClick={() => applyPreset(weekStart(), todayStr())}>
          Minggu ini
        </Button>
        <Button variant="outline" size="sm" className="h-9" onClick={() => applyPreset(monthStart(), todayStr())}>
          Bulan ini
        </Button>
        <Button variant="outline" size="sm" className="h-9 gap-1" onClick={onReset}>
          <RotateCcw size={13} />
          Reset
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="h-9 gap-1.5"
          onClick={() => exportReport(filter)}
          disabled={isExporting}
        >
          <Download className="h-4 w-4" />
          {isExporting ? 'Mengekspor...' : 'Export Excel'}
        </Button>
      </div>
    </div>
  )
}
