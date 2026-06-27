import { RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Combobox } from '@/shared/components/ui/combobox'
import { Input } from '@/shared/components/ui/input'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'

import type { SupplierReturnFilter } from '../returns.types'

interface Supplier {
  id: number
  name: string
}

interface ReturnFilterBarProps {
  filter: SupplierReturnFilter
  suppliers: Supplier[]
  onChange: (filter: SupplierReturnFilter) => void
  onReset: () => void
}

export function ReturnFilterBar({ filter, suppliers, onChange, onReset }: ReturnFilterBarProps) {
  const supplierOptions = suppliers.map((s) => ({ value: String(s.id), label: s.name }))

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

      <Combobox
        options={supplierOptions}
        value={filter.supplier_id ? String(filter.supplier_id) : undefined}
        onValueChange={(v) => onChange({ ...filter, supplier_id: v ? Number(v) : undefined })}
        placeholder="Semua Supplier"
        searchPlaceholder="Cari supplier..."
        emptyText="Supplier tidak ditemukan."
        className="w-48"
      />

      <ToggleGroup
        type="single"
        value={filter.status ?? 'all'}
        onValueChange={(v) =>
          v && onChange({ ...filter, status: v === 'all' ? undefined : v })
        }
      >
        <ToggleGroupItem value="all">Semua</ToggleGroupItem>
        <ToggleGroupItem value="pending">Pending</ToggleGroupItem>
        <ToggleGroupItem value="approved">Disetujui</ToggleGroupItem>
        <ToggleGroupItem value="rejected">Ditolak</ToggleGroupItem>
      </ToggleGroup>

      <Button variant="outline" size="sm" onClick={onReset} className="h-9 gap-1">
        <RotateCcw size={13} />
        Reset
      </Button>
    </div>
  )
}
