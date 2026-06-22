import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
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
  onReset: () => void
}

export function ShiftFilterBar({ status, onChange, onReset }: ShiftFilterBarProps) {
  return (
    <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
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

      <Button
        variant="outline"
        size="sm"
        onClick={onReset}
        className="h-9 gap-1"
      >
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
