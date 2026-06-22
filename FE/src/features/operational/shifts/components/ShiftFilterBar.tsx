import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'

import type { ShiftListFilter } from '../shifts.types'

interface ShiftFilterBarProps {
  filter: ShiftListFilter
  onChange: (filter: ShiftListFilter) => void
  onReset: () => void
}

export function ShiftFilterBar({ filter, onChange, onReset }: ShiftFilterBarProps) {
  return (
    <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
      <Input
        className="w-56 h-9"
        placeholder="Cari nama shift..."
        value={filter.search ?? ''}
        onChange={(e) => onChange({ ...filter, search: e.target.value || undefined })}
      />
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
