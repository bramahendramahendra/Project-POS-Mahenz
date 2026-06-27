import { formatRupiah } from '@/shared/utils'

import type { SalesReportSummary } from '../sales.types'

interface SalesReportSummaryCardProps {
  summary: SalesReportSummary | undefined
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

export function SalesReportSummaryCard({ summary, isLoading }: SalesReportSummaryCardProps) {
  return (
    <div className="grid grid-cols-3 gap-3">
      <SummaryCard
        label="Total Transaksi"
        value={String(summary?.total_transactions ?? 0)}
        isLoading={isLoading}
      />
      <SummaryCard
        label="Total Pendapatan"
        value={formatRupiah(summary?.total_revenue ?? 0)}
        isLoading={isLoading}
      />
      <SummaryCard
        label="Rata-rata per Transaksi"
        value={formatRupiah(summary?.avg_per_transaction ?? 0)}
        isLoading={isLoading}
      />
    </div>
  )
}
