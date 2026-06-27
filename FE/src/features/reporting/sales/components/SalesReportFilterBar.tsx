import { Download, RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import type { SalesFilter, SalesReport } from '../sales.types'
import { exportSalesCSV } from '../sales.utils'

const PAYMENT_METHODS = [
  { value: 'cash', label: 'Tunai' },
  { value: 'transfer', label: 'Transfer' },
  { value: 'card', label: 'Kartu' },
  { value: 'qris', label: 'QRIS' },
  { value: 'kredit', label: 'Kredit' },
]

interface SalesReportFilterBarProps {
  filter: SalesFilter
  onChange: (filter: SalesFilter) => void
  onReset: () => void
  exportData: SalesReport[]
}

export function SalesReportFilterBar({ filter, onChange, onReset, exportData }: SalesReportFilterBarProps) {
  return (
    <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Dari</Label>
        <Input
          type="date"
          value={filter.date_from ?? ''}
          onChange={(e) => onChange({ ...filter, date_from: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Sampai</Label>
        <Input
          type="date"
          value={filter.date_to ?? ''}
          onChange={(e) => onChange({ ...filter, date_to: e.target.value || undefined })}
          className="w-36 h-9"
        />
      </div>
      <div className="space-y-1">
        <Label className="text-xs text-gray-500">Metode Bayar</Label>
        <Select
          value={filter.payment_method ?? 'all'}
          onValueChange={(v) => onChange({ ...filter, payment_method: v === 'all' ? undefined : v })}
        >
          <SelectTrigger className="w-40 h-9">
            <SelectValue placeholder="Semua" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua</SelectItem>
            {PAYMENT_METHODS.map((m) => (
              <SelectItem key={m.value} value={m.value}>
                {m.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <Button variant="outline" size="sm" onClick={onReset} className="h-9 gap-1">
        <RotateCcw size={13} />
        Reset
      </Button>
      <Button
        variant="outline"
        size="sm"
        className="h-9 gap-1.5"
        onClick={() => exportSalesCSV(exportData)}
        disabled={exportData.length === 0}
      >
        <Download className="h-4 w-4" />
        Export CSV
      </Button>
    </div>
  )
}
