import { Download, RotateCcw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import type { SalesReport, SalesReportFilter } from '../reports.types'
import { formatDate, monthStart, todayStr } from '../reports.utils'

const PAYMENT_METHODS = [
  { value: 'cash', label: 'Tunai' },
  { value: 'transfer', label: 'Transfer' },
  { value: 'card', label: 'Kartu' },
  { value: 'qris', label: 'QRIS' },
  { value: 'kredit', label: 'Kredit' },
]

function exportCSV(data: SalesReport[], todayDate: string) {
  const headers = ['Tanggal', 'Kode Transaksi', 'Kasir', 'Customer', 'Total', 'Metode Bayar', 'Status']
  const rows = data.map((r) => [
    formatDate(r.transaction_date),
    r.transaction_code,
    r.cashier_name,
    r.customer_name ?? '-',
    r.total_amount,
    r.payment_method,
    r.status,
  ])
  const csv = [headers, ...rows].map((row) => row.join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `laporan-penjualan-${todayDate}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

interface SalesReportFilterBarProps {
  filter: SalesReportFilter
  onChange: (filter: SalesReportFilter) => void
  onReset: () => void
  exportData: SalesReport[]
}

export function SalesReportFilterBar({ filter, onChange, onReset, exportData }: SalesReportFilterBarProps) {
  const today = new Date().toISOString().split('T')[0]

  const handleReset = () => {
    onChange({ date_from: monthStart(), date_to: todayStr() })
    onReset()
  }

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
        <label className="text-xs text-gray-500">Metode Bayar</label>
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
      <Button
        variant="outline"
        size="sm"
        onClick={handleReset}
        className="h-9 gap-1"
      >
        <RotateCcw size={13} />
        Reset
      </Button>
      <Button
        variant="outline"
        size="sm"
        className="h-9 gap-1.5"
        onClick={() => exportCSV(exportData, today)}
        disabled={exportData.length === 0}
      >
        <Download className="h-4 w-4" />
        Export CSV
      </Button>
    </div>
  )
}
