import type { ReactNode } from 'react'
import { ArrowDownRight, ArrowUpRight } from 'lucide-react'

import { formatRupiah } from '@/shared/utils'

import type { SalesReportSummary } from '../sales.types'

interface SalesReportSummaryCardProps {
  summary: SalesReportSummary | undefined
  isLoading: boolean
}

// Persentase perubahan current vs previous. null = tak bisa dibandingkan.
function pctChange(current: number, previous: number): number | null {
  if (previous <= 0) return null
  return ((current - previous) / previous) * 100
}

interface DeltaProps {
  current: number
  previous: number
  available: boolean
}

// Baris kecil "▲ 12% vs periode lalu" / "▼ 5% vs periode lalu".
function DeltaLine({ current, previous, available }: DeltaProps) {
  if (!available) {
    return <p className="text-[11px] leading-tight text-gray-400">belum ada data periode lalu</p>
  }
  const pct = pctChange(current, previous)
  if (pct === null) {
    return <p className="text-[11px] leading-tight text-gray-400">belum ada data periode lalu</p>
  }
  const up = pct >= 0
  const Icon = up ? ArrowUpRight : ArrowDownRight
  const color = up ? 'text-emerald-600' : 'text-red-600'
  return (
    <p className={`flex items-center gap-0.5 text-[11px] leading-tight font-medium ${color}`}>
      <Icon size={12} />
      {Math.abs(pct).toFixed(0)}% vs periode lalu
    </p>
  )
}

interface CardProps {
  label: string
  value: string
  isLoading: boolean
  delta?: ReactNode
}

function SummaryCard({ label, value, isLoading, delta }: CardProps) {
  return (
    <div className="rounded-lg border bg-white p-4 space-y-1">
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className="text-xl font-bold text-gray-800">{value}</p>
      )}
      {!isLoading && delta}
    </div>
  )
}

export function SalesReportSummaryCard({ summary, isLoading }: SalesReportSummaryCardProps) {
  const prevAvailable = summary?.prev_available ?? false
  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <SummaryCard
        label="Total Transaksi"
        value={String(summary?.total_transactions ?? 0)}
        isLoading={isLoading}
        delta={
          <DeltaLine
            current={summary?.total_transactions ?? 0}
            previous={summary?.prev_transactions ?? 0}
            available={prevAvailable}
          />
        }
      />
      <SummaryCard
        label="Total Pendapatan"
        value={formatRupiah(summary?.total_revenue ?? 0)}
        isLoading={isLoading}
        delta={
          <DeltaLine
            current={summary?.total_revenue ?? 0}
            previous={summary?.prev_revenue ?? 0}
            available={prevAvailable}
          />
        }
      />
      <SummaryCard
        label="Rata-rata per Transaksi"
        value={formatRupiah(summary?.avg_per_transaction ?? 0)}
        isLoading={isLoading}
        delta={
          <DeltaLine
            current={summary?.avg_per_transaction ?? 0}
            previous={summary?.prev_avg_per_transaction ?? 0}
            available={prevAvailable}
          />
        }
      />
    </div>
  )
}
