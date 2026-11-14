import { formatRupiah } from '@/shared/utils'
import { ORDER_TYPE_LABEL, type OrderType } from '@/features/sales/cashier/cashier.types'

import { useOrderTypesQuery, useTopCategoriesQuery } from '../business-summary.api'
import type { DashboardPeriod } from '../business-summary.types'

/** Kategori menu terlaris dan penjualan per jenis pesanan (Coffee, keputusan 117). */
export function MenuMixCard({ period }: { period: DashboardPeriod }) {
  const { data: categories = [], isLoading: catLoading } = useTopCategoriesQuery(period)
  const { data: orderTypes = [], isLoading: otLoading } = useOrderTypesQuery(period)
  const otTotal = orderTypes.reduce((s, o) => s + o.total, 0)

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
      <div className="space-y-3 rounded-lg border bg-white p-4 shadow-sm">
        <h3 className="text-sm font-semibold text-gray-700">Kategori Menu Terlaris</h3>
        {catLoading ? (
          <div className="h-24 animate-pulse rounded bg-gray-50" />
        ) : categories.length === 0 ? (
          <p className="text-xs text-gray-400">Belum ada penjualan.</p>
        ) : (
          categories.map((c) => (
            <div key={c.category_id} className="space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-gray-700">{c.category_name}</span>
                <span className="font-medium">{formatRupiah(c.total_sales)}</span>
              </div>
              <div className="h-1.5 overflow-hidden rounded-full bg-gray-100">
                <div className="h-full rounded-full bg-[#2c3e50]" style={{ width: `${Math.max(3, c.percentage)}%` }} />
              </div>
            </div>
          ))
        )}
      </div>
      <div className="space-y-3 rounded-lg border bg-white p-4 shadow-sm">
        <h3 className="text-sm font-semibold text-gray-700">Jenis Pesanan</h3>
        {otLoading ? (
          <div className="h-24 animate-pulse rounded bg-gray-50" />
        ) : orderTypes.length === 0 ? (
          <p className="text-xs text-gray-400">Belum ada penjualan.</p>
        ) : (
          orderTypes.map((o) => (
            <div key={o.order_type} className="flex items-center justify-between text-sm">
              <span className="text-gray-700">{ORDER_TYPE_LABEL[o.order_type as OrderType] ?? 'Lainnya'}</span>
              <span className="text-gray-500">
                {o.count} transaksi · <span className="font-medium text-gray-800">{formatRupiah(o.total)}</span>
                {otTotal > 0 && ` (${Math.round((o.total / otTotal) * 100)}%)`}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  )
}
