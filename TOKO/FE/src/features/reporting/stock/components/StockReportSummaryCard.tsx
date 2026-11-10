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
  accent?: 'default' | 'red' | 'amber'
  hint?: string
}

const ACCENT: Record<NonNullable<CardProps['accent']>, { border: string; value: string }> = {
  default: { border: '', value: 'text-gray-800' },
  red: { border: 'border-red-200 bg-red-50', value: 'text-red-700' },
  amber: { border: 'border-amber-200 bg-amber-50', value: 'text-amber-700' },
}

function SummaryCard({ label, value, isLoading, accent = 'default', hint }: CardProps) {
  const a = ACCENT[accent]
  return (
    <div className={`rounded-lg border p-4 space-y-1 ${accent === 'default' ? 'bg-white' : a.border}`}>
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className={`text-xl font-bold ${a.value}`}>{value}</p>
      )}
      {hint && <p className="text-[11px] leading-tight text-gray-400">{hint}</p>}
    </div>
  )
}

export function StockReportSummaryCard({ summary, isLoading }: StockReportSummaryCardProps) {
  return (
    <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <SummaryCard
        label="Total Produk"
        value={String(summary?.total_products ?? 0)}
        isLoading={isLoading}
      />
      <SummaryCard
        label="Stok Habis"
        value={String(summary?.out_of_stock_count ?? 0)}
        isLoading={isLoading}
        accent="red"
        hint="perlu segera dibeli"
      />
      <SummaryCard
        label="Stok Menipis"
        value={String(summary?.low_only_count ?? 0)}
        isLoading={isLoading}
        accent="amber"
        hint="mulai perlu diperhatikan"
      />
      <SummaryCard
        label="Total Nilai Stok"
        value={formatRupiah(summary?.total_stock_value ?? 0)}
        isLoading={isLoading}
        hint="modal barang tersimpan di toko"
      />
    </div>
  )
}
