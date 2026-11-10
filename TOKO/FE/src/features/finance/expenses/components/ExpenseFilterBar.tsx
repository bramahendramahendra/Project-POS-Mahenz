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

import type { ExpenseCategory, ExpenseListFilter } from '../expenses.types'
import { EXPENSE_CATEGORIES } from '../expenses.schema'

interface ExpenseFilterBarProps {
  filter: ExpenseListFilter
  onChange: (filter: ExpenseListFilter) => void
  onReset: () => void
}

export function ExpenseFilterBar({ filter, onChange, onReset }: ExpenseFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Dari</span>
        <Input
          type="date"
          value={filter.start_date ?? ''}
          onChange={(e) => onChange({ ...filter, start_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Sampai</span>
        <Input
          type="date"
          value={filter.end_date ?? ''}
          onChange={(e) => onChange({ ...filter, end_date: e.target.value || undefined })}
          className="w-40 h-9"
        />
      </div>
      <div className="space-y-1">
        <span className="text-xs text-gray-500">Kategori</span>
        <Select
          value={filter.category ?? 'all'}
          onValueChange={(v) =>
            onChange({ ...filter, category: v === 'all' ? undefined : (v as ExpenseCategory) })
          }
        >
          <SelectTrigger className="w-44 h-9">
            <SelectValue placeholder="Semua Kategori" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Kategori</SelectItem>
            {EXPENSE_CATEGORIES.map((cat) => (
              <SelectItem key={cat.value} value={cat.value}>
                {cat.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
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
