import { Banknote, ShoppingBag } from 'lucide-react'

import { Card, CardContent } from '@/shared/components/ui/card'
import { formatRupiah } from '@/shared/utils'

import { useTodaySummaryQuery } from '../dashboard.api'

export function MyTodayPerformance() {
  const { data, isLoading } = useTodaySummaryQuery()

  return (
    <div className="grid grid-cols-2 gap-3">
      <Card>
        <CardContent className="p-4 space-y-1">
          <div className="flex items-start gap-2 text-gray-500 text-sm">
            <ShoppingBag size={15} className="shrink-0 mt-0.5" />
            <span>Transaksi Saya Hari Ini</span>
          </div>
          {isLoading ? (
            <div className="h-7 w-16 animate-pulse rounded bg-gray-100" />
          ) : (
            <p className="text-xl font-bold text-gray-900">{data?.total_transactions ?? 0}</p>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4 space-y-1">
          <div className="flex items-start gap-2 text-gray-500 text-sm">
            <Banknote size={15} className="shrink-0 mt-0.5" />
            <span>Penjualan Saya Hari Ini</span>
          </div>
          {isLoading ? (
            <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
          ) : (
            <p className="text-xl font-bold text-gray-900">{formatRupiah(data?.total_sales ?? 0)}</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
