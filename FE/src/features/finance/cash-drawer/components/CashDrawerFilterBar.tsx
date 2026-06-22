import { useMemo } from 'react'
import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useShiftOptionsQuery } from '@/features/operational/shifts'

import { useKasirOptionsQuery } from '../cash-drawer.api'
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
  onReset: () => void
  showKasirFilter?: boolean
}

export function CashDrawerFilterBar({
  filter,
  onChange,
  onReset,
  showKasirFilter = false,
}: CashDrawerFilterBarProps) {
  const { data: shiftOptionsRaw } = useShiftOptionsQuery()
  const { data: kasirOptionsRaw } = useKasirOptionsQuery()
  const shiftOptions = useMemo(() => shiftOptionsRaw ?? [], [shiftOptionsRaw])
  const kasirOptions = useMemo(() => kasirOptionsRaw ?? [], [kasirOptionsRaw])

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

      <div className="space-y-1">
        <label className="text-xs text-gray-500">Shift</label>
        <Select
          value={filter.shift_id ? String(filter.shift_id) : 'all'}
          onValueChange={(v) =>
            onChange({ ...filter, shift_id: v === 'all' ? undefined : Number(v) })
          }
        >
          <SelectTrigger className="w-44 h-9">
            <SelectValue placeholder="Semua shift" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua shift</SelectItem>
            {shiftOptions.map((s) => (
              <SelectItem key={s.id} value={String(s.id)}>
                {s.name} ({s.start_time}–{s.end_time})
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-1">
        <label className="text-xs text-gray-500">Status</label>
        <Select
          value={filter.status || 'all'}
          onValueChange={(v) => onChange({ ...filter, status: v === 'all' ? undefined : v })}
        >
          <SelectTrigger className="w-32 h-9">
            <SelectValue placeholder="Semua" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua</SelectItem>
            <SelectItem value="open">Buka</SelectItem>
            <SelectItem value="closed">Tutup</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {showKasirFilter && (
        <div className="space-y-1">
          <label className="text-xs text-gray-500">Kasir</label>
          <Select
            value={filter.user_id ? String(filter.user_id) : 'all'}
            onValueChange={(v) =>
              onChange({ ...filter, user_id: v === 'all' ? undefined : Number(v) })
            }
          >
            <SelectTrigger className="w-44 h-9">
              <SelectValue placeholder="Semua kasir" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Semua kasir</SelectItem>
              {kasirOptions.map((k) => (
                <SelectItem key={k.id} value={String(k.id)}>
                  {k.full_name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}

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

      <Button variant="outline" size="sm" onClick={onReset} className="h-9 gap-1">
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
