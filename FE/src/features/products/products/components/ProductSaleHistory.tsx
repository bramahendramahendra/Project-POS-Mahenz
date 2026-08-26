import { useState } from 'react'
import { ShoppingCart } from 'lucide-react'

import { formatDate, formatRupiah } from '@/shared/utils'

import { useProductSaleHistoryQuery } from '../products.api'

interface Props {
  productId: number
}

export function ProductSaleHistory({ productId }: Props) {
  const [page, setPage] = useState(1)
  const [status, setStatus] = useState('')
  const limit = 10
  const { data, isLoading } = useProductSaleHistoryQuery(productId, page, limit, status)

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-8 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  if (!data || data.total === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-10 text-gray-400">
        <ShoppingCart size={28} className="opacity-30 mb-2" />
        <p className="text-sm">Belum ada riwayat penjualan</p>
      </div>
    )
  }

  const totalPages = Math.ceil(data.total / limit)

  return (
    <div className="space-y-4 text-sm">
      {/* Summary */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4 rounded-lg border border-green-200 bg-green-50 p-3">
        <SummaryItem label="Total Transaksi" value={String(data.summary.total_transactions)} />
        <SummaryItem label="Total Qty Terjual" value={`${data.summary.total_qty}`} />
        <SummaryItem label="Total Pendapatan" value={formatRupiah(data.summary.total_revenue)} />
        <SummaryItem label="Rata-rata/Hari" value={`${data.summary.average_per_day}`} />
      </div>

      {/* Filter */}
      <div className="flex items-center gap-2">
        <span className="text-xs text-gray-500">Status:</span>
        <select
          value={status}
          onChange={(e) => { setStatus(e.target.value); setPage(1) }}
          className="rounded border border-gray-200 px-2 py-1 text-xs bg-white"
        >
          <option value="">Semua</option>
          <option value="completed">Completed</option>
          <option value="void">Void</option>
        </select>
      </div>

      {/* Table */}
      <div className="rounded-lg border overflow-x-auto">
        <table className="w-full text-xs">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-2.5 py-2 text-left font-semibold text-gray-600">Tanggal</th>
              <th className="px-2.5 py-2 text-left font-semibold text-gray-600">Kode Transaksi</th>
              <th className="px-2.5 py-2 text-left font-semibold text-gray-600">Pelanggan</th>
              <th className="px-2.5 py-2 text-left font-semibold text-gray-600">Satuan</th>
              <th className="px-2.5 py-2 text-right font-semibold text-gray-600">Qty</th>
              <th className="px-2.5 py-2 text-right font-semibold text-gray-600">Harga</th>
              <th className="px-2.5 py-2 text-right font-semibold text-gray-600">Subtotal</th>
              <th className="px-2.5 py-2 text-center font-semibold text-gray-600">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {data.items.map((item, idx) => (
              <tr key={idx} className={`hover:bg-gray-50 ${item.status === 'void' ? 'opacity-50' : ''}`}>
                <td className={`px-2.5 py-2 text-gray-600 ${item.status === 'void' ? 'line-through' : ''}`}>
                  {formatDate(item.transaction_date)}
                </td>
                <td className="px-2.5 py-2">
                  <span className={`font-medium ${item.status === 'void' ? 'text-gray-400 line-through' : 'text-blue-600'}`}>
                    {item.transaction_code}
                  </span>
                </td>
                <td className={`px-2.5 py-2 text-gray-700 ${item.status === 'void' ? 'line-through' : ''}`}>
                  {item.customer_name || '-'}
                </td>
                <td className="px-2.5 py-2">
                  <span className="inline-flex items-center rounded-full bg-blue-100 px-2 py-0.5 text-[10px] font-medium text-blue-700">
                    {item.unit}
                  </span>
                </td>
                <td className={`px-2.5 py-2 text-right ${item.status === 'void' ? 'line-through' : ''}`}>
                  {item.quantity}
                </td>
                <td className={`px-2.5 py-2 text-right ${item.status === 'void' ? 'line-through' : ''}`}>
                  {formatRupiah(item.price)}
                </td>
                <td className={`px-2.5 py-2 text-right font-semibold ${item.status === 'void' ? 'line-through' : ''}`}>
                  {formatRupiah(item.subtotal)}
                </td>
                <td className="px-2.5 py-2 text-center">
                  <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-medium ${
                    item.status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-600'
                  }`}>
                    {item.status === 'completed' ? 'Completed' : 'Void'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between text-xs text-gray-500">
          <span>{(page - 1) * limit + 1}-{Math.min(page * limit, Number(data.total))} dari {data.total}</span>
          <div className="flex gap-1">
            {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => i + 1).map((p) => (
              <button
                key={p}
                onClick={() => setPage(p)}
                className={`px-2.5 py-1 rounded border text-xs ${
                  p === page ? 'bg-blue-600 text-white border-blue-600' : 'bg-white border-gray-200 hover:bg-gray-50'
                }`}
              >
                {p}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

function SummaryItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-[10px] text-gray-500 uppercase tracking-wide">{label}</p>
      <p className="text-sm font-bold text-gray-900 mt-0.5">{value}</p>
    </div>
  )
}
