import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import type { ShiftStatus } from '../shifts.types'

interface ShiftFilterBarProps {
  status: ShiftStatus | 'all'
  onChange: (status: ShiftStatus | 'all') => void
}

export function ShiftFilterBar({ status, onChange }: ShiftFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <Select
        value={status}
        onValueChange={(v) => onChange(v as ShiftStatus | 'all')}
      >
        <SelectTrigger className="w-40 h-9">
          <SelectValue placeholder="Semua Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Semua</SelectItem>
          <SelectItem value="open">Berjalan</SelectItem>
          <SelectItem value="closed">Selesai</SelectItem>
        </SelectContent>
      </Select>
    </div>
  )
}
