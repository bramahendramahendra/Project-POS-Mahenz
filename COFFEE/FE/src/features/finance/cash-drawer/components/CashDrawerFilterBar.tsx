import { CalendarDays, RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Combobox } from '@/shared/components/ui/combobox'
import { Input } from '@/shared/components/ui/input'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { monthStart, todayStr } from '@/shared/utils'
import { useShiftOptionsQuery } from '@/features/operational/shifts'

import { useKasirOptionsQuery } from '../cash-drawer.api'
import type { CashDrawerListFilter } from '../cash-drawer.types'


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
  const { data: shiftOptions = [] } = useShiftOptionsQuery()
  const { data: kasirOptions = [] } = useKasirOptionsQuery()

  const kasirComboboxOptions = kasirOptions.map((k) => ({ value: String(k.id), label: k.full_name }))

  return (
    <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
      <Input
        type="date"
        value={filter.start_date ?? ''}
        onChange={(e) => onChange({ ...filter, start_date: e.target.value || undefined })}
        className="w-36 h-9 text-sm"
      />
      <span className="text-xs text-gray-400">—</span>
      <Input
        type="date"
        value={filter.end_date ?? ''}
        onChange={(e) => onChange({ ...filter, end_date: e.target.value || undefined })}
        className="w-36 h-9 text-sm"
      />

      <Select
        value={filter.shift_id ? String(filter.shift_id) : 'all'}
        onValueChange={(v) =>
          onChange({ ...filter, shift_id: v === 'all' ? undefined : Number(v) })
        }
      >
        <SelectTrigger className="w-44 h-9 text-sm">
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

      <ToggleGroup
        type="single"
        value={filter.status ?? 'all'}
        onValueChange={(v) => v && onChange({ ...filter, status: v === 'all' ? undefined : v })}
      >
        <ToggleGroupItem value="all">Semua</ToggleGroupItem>
        <ToggleGroupItem value="open">Buka</ToggleGroupItem>
        <ToggleGroupItem value="closed">Tutup</ToggleGroupItem>
      </ToggleGroup>

      {showKasirFilter && (
        <Combobox
          options={kasirComboboxOptions}
          value={filter.user_id ? String(filter.user_id) : undefined}
          onValueChange={(v) => onChange({ ...filter, user_id: v ? Number(v) : undefined })}
          placeholder="Semua kasir"
          searchPlaceholder="Cari kasir..."
          emptyText="Kasir tidak ditemukan."
          className="w-44"
        />
      )}

      <Button
        variant="outline"
        size="sm"
        className="h-9 gap-1"
        onClick={() => onChange({ ...filter, start_date: monthStart(), end_date: todayStr() })}
      >
        <CalendarDays size={13} />
        Bulan ini
      </Button>

      <Button variant="outline" size="sm" onClick={onReset} className="h-9 gap-1">
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
