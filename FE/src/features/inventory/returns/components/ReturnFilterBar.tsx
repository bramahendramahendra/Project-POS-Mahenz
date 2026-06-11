import { Input } from '@/shared/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import type { SupplierReturnFilter } from '../returns.types'

interface Supplier {
  id: number
  name: string
}

interface ReturnFilterBarProps {
  filter: SupplierReturnFilter
  suppliers: Supplier[]
  onChange: (filter: SupplierReturnFilter) => void
}

export function ReturnFilterBar({ filter, suppliers, onChange }: ReturnFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Dari</label>
        <Input
          type="date"
          value={filter.date_from ?? ''}
          onChange={(e) => onChange({ ...filter, date_from: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Sampai</label>
        <Input
          type="date"
          value={filter.date_to ?? ''}
          onChange={(e) => onChange({ ...filter, date_to: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Supplier</label>
        <Select
          value={filter.supplier_id ? String(filter.supplier_id) : 'all'}
          onValueChange={(v) =>
            onChange({ ...filter, supplier_id: v === 'all' ? undefined : Number(v) })
          }
        >
          <SelectTrigger className="w-44 h-9">
            <SelectValue placeholder="Semua Supplier" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Supplier</SelectItem>
            {suppliers.map((s) => (
              <SelectItem key={s.id} value={String(s.id)}>
                {s.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-1">
        <label className="text-xs text-gray-500">Status</label>
        <Select
          value={filter.status ?? 'all'}
          onValueChange={(v) => onChange({ ...filter, status: v === 'all' ? undefined : v })}
        >
          <SelectTrigger className="w-40 h-9">
            <SelectValue placeholder="Semua Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Status</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="approved">Disetujui</SelectItem>
            <SelectItem value="rejected">Ditolak</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
