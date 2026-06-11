import { Button } from '@/shared/components/ui/button'

import type { DashboardPeriod } from '../dashboard.types'

const PERIODS: { label: string; value: DashboardPeriod }[] = [
  { label: 'Hari Ini', value: 'today' },
  { label: 'Minggu Ini', value: 'week' },
  { label: 'Bulan Ini', value: 'month' },
]

interface DashboardPeriodSelectorProps {
  period: DashboardPeriod
  onChange: (period: DashboardPeriod) => void
}

export function DashboardPeriodSelector({ period, onChange }: DashboardPeriodSelectorProps) {
  return (
    <div className="flex gap-1 rounded-lg border p-1 bg-gray-50">
      {PERIODS.map((p) => (
        <Button
          key={p.value}
          size="sm"
          variant={period === p.value ? 'default' : 'ghost'}
          className="h-7 text-xs"
          onClick={() => onChange(p.value)}
        >
          {p.label}
        </Button>
      ))}
    </div>
  )
}
