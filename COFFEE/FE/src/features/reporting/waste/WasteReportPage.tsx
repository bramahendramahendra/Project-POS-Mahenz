import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'

import { PageHeader } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { api } from '@/services'
import { formatRupiah, monthStart, todayStr } from '@/shared/utils'

// Laporan barang terbuang (tahap 4.10, keputusan 100, 116): bahan baku dan barang jadi yang terbuang beserta
// alasan, pelaku, waktu, dan nilainya.

interface WasteRow {
  created_at: string
  kind: 'ingredient' | 'product'
  item_name: string
  qty: number
  unit: string
  value: number
  reason: string
  user_name: string
}

interface WasteReport {
  total_value: number
  total_ingredient_value: number
  total_product_value: number
  items: WasteRow[]
}

const th = 'px-3 py-2 text-xs font-medium text-gray-500'

export function WasteReportPage() {
  const [range, setRange] = useState({ date_from: monthStart(), date_to: todayStr() })
  const { data, isLoading } = useQuery({
    queryKey: ['reports', 'waste', range],
    queryFn: () => api.post<WasteReport>('/reports/waste', range),
  })
  const items = data?.items ?? []

  return (
    <div className="space-y-4">
      <PageHeader title="Laporan Barang Terbuang" breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Barang Terbuang' }]} />
      <div className="flex flex-wrap items-end gap-3 rounded-lg border bg-white p-3">
        <div className="space-y-1">
          <Label className="text-xs text-gray-500">Dari</Label>
          <Input type="date" aria-label="Dari tanggal" value={range.date_from} onChange={(e) => setRange((r) => ({ ...r, date_from: e.target.value }))} className="h-9 w-[160px]" />
        </div>
        <div className="space-y-1">
          <Label className="text-xs text-gray-500">Sampai</Label>
          <Input type="date" aria-label="Sampai tanggal" value={range.date_to} onChange={(e) => setRange((r) => ({ ...r, date_to: e.target.value }))} className="h-9 w-[160px]" />
        </div>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        {[
          ['Total terbuang', data?.total_value],
          ['Bahan baku', data?.total_ingredient_value],
          ['Barang jadi', data?.total_product_value],
        ].map(([label, value]) => (
          <div key={label as string} className="rounded-lg border bg-white p-4">
            <p className="text-xs text-gray-500">{label}</p>
            <p className="text-lg font-bold text-red-600">{isLoading ? '…' : formatRupiah(Number(value ?? 0))}</p>
          </div>
        ))}
      </div>

      <div className="overflow-x-auto rounded-lg border bg-white">
        <table className="w-full min-w-[760px] text-sm">
          <thead className="border-b bg-gray-50">
            <tr>
              <th className={`${th} text-left`}>Waktu</th>
              <th className={`${th} text-left`}>Jenis</th>
              <th className={`${th} text-left`}>Barang</th>
              <th className={`${th} text-right`}>Jumlah</th>
              <th className={`${th} text-right`}>Nilai</th>
              <th className={`${th} text-left`}>Alasan</th>
              <th className={`${th} text-left`}>Oleh</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {!isLoading && items.length === 0 && (
              <tr>
                <td colSpan={7} className="px-3 py-6 text-center text-gray-400">Tidak ada barang terbuang di periode ini</td>
              </tr>
            )}
            {items.map((w, i) => (
              <tr key={i} aria-label={`Terbuang ${w.item_name}`}>
                <td className="whitespace-nowrap px-3 py-2 text-gray-600">{w.created_at.slice(0, 16)}</td>
                <td className="px-3 py-2">{w.kind === 'ingredient' ? 'Bahan baku' : 'Barang jadi'}</td>
                <td className="px-3 py-2 font-medium text-gray-800">{w.item_name}</td>
                <td className="px-3 py-2 text-right">
                  {w.qty.toLocaleString('id-ID', { maximumFractionDigits: 3 })} {w.unit}
                </td>
                <td className="px-3 py-2 text-right text-red-600">{formatRupiah(w.value)}</td>
                <td className="px-3 py-2 text-gray-600">{w.reason || '—'}</td>
                <td className="px-3 py-2 text-gray-600">{w.user_name || '—'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
