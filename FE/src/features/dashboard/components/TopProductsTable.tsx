import { formatRupiah } from '@/shared/utils'

import type { TopProductItem } from '../dashboard.types'

interface TopProductsTableProps {
  data: TopProductItem[]
  isLoading: boolean
}

export function TopProductsTable({ data, isLoading }: TopProductsTableProps) {
  return (
    <div className="rounded-lg border bg-white overflow-hidden">
      <div className="px-4 py-3 border-b">
        <h3 className="font-semibold text-gray-700 text-sm">Top Produk Terlaris</h3>
      </div>
      {isLoading ? (
        <div className="space-y-3 p-4">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-8 animate-pulse rounded bg-gray-100" />
          ))}
        </div>
      ) : data.length === 0 ? (
        <p className="p-6 text-center text-sm text-gray-400">Belum ada data</p>
      ) : (
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-3 py-2.5 text-center font-medium text-gray-600 w-8">#</th>
              <th className="px-3 py-2.5 text-left font-medium text-gray-600">Produk</th>
              <th className="px-3 py-2.5 text-right font-medium text-gray-600">Qty</th>
              <th className="px-3 py-2.5 text-right font-medium text-gray-600">Revenue</th>
            </tr>
          </thead>
          <tbody>
            {data.map((p, i) => (
              <tr key={p.product_id} className="border-b last:border-0 hover:bg-gray-50">
                <td className="px-3 py-2.5 text-center text-gray-400 font-mono">{i + 1}</td>
                <td className="px-3 py-2.5 font-medium text-gray-800">{p.product_name}</td>
                <td className="px-3 py-2.5 text-right text-gray-600">{p.total_qty}</td>
                <td className="px-3 py-2.5 text-right font-semibold text-blue-600">
                  {formatRupiah(p.total_value)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
