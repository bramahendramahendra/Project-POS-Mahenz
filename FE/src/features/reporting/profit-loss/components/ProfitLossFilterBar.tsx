import { AlertTriangle, CalendarRange, Download } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { formatDate, getWIBNow, todayStr } from '@/shared/utils'

import { useExportProfitLossMutation } from '../profit-loss.api'
import type { ProfitLossDateFilter } from '../profit-loss.types'

interface ProfitLossFilterBarProps {
  filter: ProfitLossDateFilter
  onChange: (filter: ProfitLossDateFilter) => void
}

// ── Daftar preset periode. Tiap preset punya id, label, dan fungsi rentang. ──
type Preset = { id: string; label: string; range: () => { from: string; to: string } }

const PRESETS: Preset[] = [
  {
    id: 'today',
    label: 'Hari ini',
    range: () => ({ from: todayStr(), to: todayStr() }),
  },
  {
    id: 'yesterday',
    label: 'Kemarin',
    range: () => {
      const y = getWIBNow().subtract(1, 'day').format('YYYY-MM-DD')
      return { from: y, to: y }
    },
  },
  {
    id: 'this_week',
    label: 'Minggu ini',
    range: () => {
      const now = getWIBNow()
      const from = now.subtract(now.day() === 0 ? 6 : now.day() - 1, 'day').format('YYYY-MM-DD')
      return { from, to: todayStr() }
    },
  },
  {
    id: 'this_month',
    label: 'Bulan ini',
    range: () => ({ from: getWIBNow().startOf('month').format('YYYY-MM-DD'), to: todayStr() }),
  },
  {
    id: 'last_month',
    label: 'Bulan lalu',
    range: () => {
      const lastMonth = getWIBNow().subtract(1, 'month')
      return {
        from: lastMonth.startOf('month').format('YYYY-MM-DD'),
        to: lastMonth.endOf('month').format('YYYY-MM-DD'),
      }
    },
  },
  {
    id: 'last_30',
    label: '30 hari terakhir',
    range: () => ({ from: getWIBNow().subtract(29, 'day').format('YYYY-MM-DD'), to: todayStr() }),
  },
]

// Cari preset mana yang cocok dengan filter aktif sekarang (untuk highlight).
function activePresetId(filter: ProfitLossDateFilter): string | null {
  for (const p of PRESETS) {
    const r = p.range()
    if (r.from === filter.date_from && r.to === filter.date_to) return p.id
  }
  return null
}

// Label periode aktif dalam bahasa manusia.
function periodLabel(filter: ProfitLossDateFilter): string {
  if (!filter.date_from || !filter.date_to) return 'Semua periode'
  if (filter.date_from === filter.date_to) return formatDate(filter.date_from)
  return `${formatDate(filter.date_from)} – ${formatDate(filter.date_to)}`
}

export function ProfitLossFilterBar({ filter, onChange }: ProfitLossFilterBarProps) {
  const { mutate: exportReport, isPending: isExporting } = useExportProfitLossMutation()

  const applyPreset = (p: Preset) => {
    const r = p.range()
    onChange({ date_from: r.from, date_to: r.to })
  }

  const active = activePresetId(filter)
  const invalidRange =
    !!filter.date_from && !!filter.date_to && filter.date_from > filter.date_to

  return (
    <div className="space-y-3 rounded-lg border bg-white p-3">
      {/* Baris 1: preset periode + export */}
      <div className="flex flex-wrap items-center gap-2">
        {PRESETS.map((p) => {
          const isActive = active === p.id
          return (
            <button
              key={p.id}
              type="button"
              onClick={() => applyPreset(p)}
              className={`h-9 rounded-md border px-3 text-sm transition ${
                isActive
                  ? 'border-primary bg-primary text-primary-foreground font-semibold'
                  : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              {p.label}
            </button>
          )
        })}
        <Button
          variant="outline"
          size="sm"
          className="ml-auto h-9 gap-1.5"
          onClick={() => exportReport(filter)}
          disabled={isExporting}
        >
          <Download className="h-4 w-4" />
          {isExporting ? 'Mengekspor...' : 'Export Excel'}
        </Button>
      </div>

      {/* Baris 2: tanggal manual + label periode aktif */}
      <div className="flex flex-wrap items-end gap-3">
        <div className="space-y-1">
          <Label className="text-xs text-gray-500">Dari</Label>
          <Input
            type="date"
            max={filter.date_to ?? undefined}
            value={filter.date_from ?? ''}
            onChange={(e) => onChange({ ...filter, date_from: e.target.value || undefined })}
            className="h-9 w-40"
          />
        </div>
        <div className="space-y-1">
          <Label className="text-xs text-gray-500">Sampai</Label>
          <Input
            type="date"
            min={filter.date_from ?? undefined}
            max={todayStr()}
            value={filter.date_to ?? ''}
            onChange={(e) => onChange({ ...filter, date_to: e.target.value || undefined })}
            className="h-9 w-40"
          />
        </div>

        {/* Label periode aktif — menjawab "sedang lihat tanggal berapa" */}
        <div className="flex items-center gap-2 rounded-md bg-sky-50 px-3 py-2 text-sm text-sky-700">
          <CalendarRange size={15} />
          <span>
            Menampilkan: <span className="font-semibold">{periodLabel(filter)}</span>
            {active && <span className="text-sky-500"> ({PRESETS.find((p) => p.id === active)?.label})</span>}
          </span>
        </div>
      </div>

      {invalidRange && (
        <div className="flex items-center gap-1.5 text-xs text-red-600">
          <AlertTriangle size={13} />
          Tanggal &quot;Dari&quot; lebih besar dari &quot;Sampai&quot; — perbaiki rentang tanggal.
        </div>
      )}
    </div>
  )
}
