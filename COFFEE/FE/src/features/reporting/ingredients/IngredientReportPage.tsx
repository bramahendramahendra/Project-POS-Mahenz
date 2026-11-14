import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'

import { PageHeader } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { api } from '@/services'
import { formatRupiah, monthStart, todayStr } from '@/shared/utils'

// Laporan bahan baku (tahap 4.10, keputusan 114, 116): per bahan di periode — pembelian, terpakai, terbuang,
// stok sisa, dan nilainya (harga beli rata-rata).

interface IngredientRow {
  ingredient_id: number
  name: string
  unit_abbr: string
  purchased_qty: number
  purchase_value: number
  used_qty: number
  wasted_qty: number
  stock: number
  avg_price: number
  used_value: number
  wasted_value: number
}

interface IngredientReport {
  total_purchase_value: number
  total_used_value: number
  total_wasted_value: number
  items: IngredientRow[]
}

const qtyFmt = (n: number, unit: string) => `${n.toLocaleString('id-ID', { maximumFractionDigits: 3 })} ${unit}`
const th = 'px-3 py-2 text-xs font-medium text-gray-500'

export function IngredientReportPage() {
  const [range, setRange] = useState({ date_from: monthStart(), date_to: todayStr() })
  const { data, isLoading } = useQuery({
    queryKey: ['reports', 'ingredients', range],
    queryFn: () => api.post<IngredientReport>('/reports/ingredients', range),
  })
  // Bahan tanpa aktivitas di periode tetap tampil bila masih punya stok.
  const items = (data?.items ?? []).filter((i) => i.purchased_qty || i.used_qty || i.wasted_qty || i.stock)

  return (
    <div className="space-y-4">
      <PageHeader title="Laporan Bahan Baku" breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Bahan Baku' }]} />
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
          ['Pembelian bahan', data?.total_purchase_value],
          ['Nilai bahan terpakai', data?.total_used_value],
          ['Nilai bahan terbuang', data?.total_wasted_value],
        ].map(([label, value]) => (
          <div key={label as string} className="rounded-lg border bg-white p-4">
            <p className="text-xs text-gray-500">{label}</p>
            <p className="text-lg font-bold text-gray-900">{isLoading ? '…' : formatRupiah(Number(value ?? 0))}</p>
          </div>
        ))}
      </div>

      <div className="overflow-x-auto rounded-lg border bg-white">
        <table className="w-full min-w-[820px] text-sm">
          <thead className="border-b bg-gray-50">
            <tr>
              <th className={`${th} text-left`}>Bahan</th>
              <th className={`${th} text-right`}>Dibeli</th>
              <th className={`${th} text-right`}>Nilai beli</th>
              <th className={`${th} text-right`}>Terpakai</th>
              <th className={`${th} text-right`}>Terbuang</th>
              <th className={`${th} text-right`}>Stok sisa</th>
              <th className={`${th} text-right`}>Harga rata-rata</th>
              <th className={`${th} text-right`}>Nilai terpakai + terbuang</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {!isLoading && items.length === 0 && (
              <tr>
                <td colSpan={8} className="px-3 py-6 text-center text-gray-400">Belum ada data bahan di periode ini</td>
              </tr>
            )}
            {items.map((i) => (
              <tr key={i.ingredient_id} aria-label={`Bahan ${i.name}`}>
                <td className="px-3 py-2 font-medium text-gray-800">{i.name}</td>
                <td className="px-3 py-2 text-right">{qtyFmt(i.purchased_qty, i.unit_abbr)}</td>
                <td className="px-3 py-2 text-right">{formatRupiah(i.purchase_value)}</td>
                <td className="px-3 py-2 text-right">{qtyFmt(i.used_qty, i.unit_abbr)}</td>
                <td className="px-3 py-2 text-right text-red-600">{qtyFmt(i.wasted_qty, i.unit_abbr)}</td>
                <td className={`px-3 py-2 text-right ${i.stock < 0 ? 'text-red-600' : ''}`}>{qtyFmt(i.stock, i.unit_abbr)}</td>
                <td className="px-3 py-2 text-right text-gray-500">{i.avg_price ? formatRupiah(i.avg_price) : '—'}</td>
                <td className="px-3 py-2 text-right font-medium">{formatRupiah(i.used_value + i.wasted_value)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="text-[11px] text-gray-400">
        Terpakai dan terbuang yang dibatalkan tidak dihitung. Nilai memakai harga beli rata-rata dari semua pembelian bahan; stok sisa adalah stok saat ini.
      </p>
    </div>
  )
}
