import { formatRupiah } from '@/shared/utils'

import type { StockSummary } from '../stock.types'

interface StockReportSummaryCardProps {
  summary: StockSummary | undefined
  isLoading: boolean
}

interface CardProps {
  label: string
  value: string
  isLoading: boolean
}

function SummaryCard({ label, value, isLoading }: CardProps) {
  return (
    <div className="rounded-lg border bg-white p-4 space-y-1">
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className="text-xl font-bold text-gray-800">{value}</p>
      )}
    </div>
  )
}

export function StockReportSummaryCard({ summary, isLoading }: StockReportSummaryCardProps) {
  return (
    <div className="grid grid-cols-3 gap-3">
      <SummaryCard
        label="Total Produk"
        value={String(summary?.total_products ?? 0)}
        isLoading={isLoading}
      />
      <SummaryCard
        label="Produk Stok Rendah"
        value={String(summary?.low_stock_count ?? 0)}
        isLoading={isLoading}
      />
      <SummaryCard
        label="Total Nilai Stok"
        value={formatRupiah(summary?.total_stock_value ?? 0)}
        isLoading={isLoading}
      />
    </div>
  )
}
