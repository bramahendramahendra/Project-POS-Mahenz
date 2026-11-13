import type { ReactNode } from 'react'
import { ArrowDownRight, ArrowUpRight, Banknote, CreditCard, PackageX, ShoppingBag, TrendingUp } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

import { formatRupiah } from '@/shared/utils'

import type { DashboardPeriod, DashboardStats } from '../business-summary.types'

interface SummaryCardsProps {
  stats: DashboardStats | undefined
  isLoading: boolean
  period: DashboardPeriod
}

const PERIOD_LABEL: Record<DashboardPeriod, string> = {
  today: 'Hari Ini',
  week: 'Minggu Ini',
  month: 'Bulan Ini',
}

// Label periode sebelumnya, untuk teks pembanding ("vs kemarin", dst).
const PREV_LABEL: Record<DashboardPeriod, string> = {
  today: 'kemarin',
  week: 'minggu lalu',
  month: 'bulan lalu',
}

// Hitung persentase perubahan current vs previous. null = tak bisa dibandingkan.
function pctChange(current: number, previous: number): number | null {
  if (previous <= 0) return null
  return ((current - previous) / previous) * 100
}

interface DeltaProps {
  current: number
  previous: number
  available: boolean
  prevLabel: string
}

// Baris kecil "▲ 12% vs kemarin" / "▼ 5% vs kemarin" / "belum ada pembanding".
function DeltaLine({ current, previous, available, prevLabel }: DeltaProps) {
  if (!available) {
    return <p className="text-[11px] leading-tight text-gray-400">belum ada data {prevLabel}</p>
  }
  const pct = pctChange(current, previous)
  if (pct === null) {
    return <p className="text-[11px] leading-tight text-gray-400">belum ada data {prevLabel}</p>
  }
  const up = pct >= 0
  const Icon = up ? ArrowUpRight : ArrowDownRight
  const color = up ? 'text-emerald-600' : 'text-red-600'
  return (
    <p className={`flex items-center gap-0.5 text-[11px] leading-tight font-medium ${color}`}>
      <Icon size={12} />
      {Math.abs(pct).toFixed(0)}% vs {prevLabel}
    </p>
  )
}

interface StatCardProps {
  icon: LucideIcon
  label: string
  value: string
  isLoading: boolean
  accent?: 'default' | 'red' | 'amber'
  hint?: string
  delta?: ReactNode
}

const ACCENT: Record<NonNullable<StatCardProps['accent']>, { border: string; value: string }> = {
  default: { border: '', value: 'text-gray-900' },
  red: { border: 'border-red-200 bg-red-50', value: 'text-red-700' },
  amber: { border: 'border-amber-200 bg-amber-50', value: 'text-amber-700' },
}

function StatCard({ icon: Icon, label, value, isLoading, accent = 'default', hint, delta }: StatCardProps) {
  const a = ACCENT[accent]
  return (
    <div className={`rounded-lg border p-4 shadow-sm space-y-1 ${accent === 'default' ? 'bg-white' : a.border}`}>
      <div className="flex items-center gap-2 text-gray-500 text-sm">
        <Icon size={15} />
        <span>{label}</span>
      </div>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className={`text-xl font-bold ${a.value}`}>{value}</p>
      )}
      {!isLoading && delta}
      {!isLoading && hint && <p className="text-[11px] leading-tight text-gray-400">{hint}</p>}
    </div>
  )
}

export function SummaryCards({ stats, isLoading, period }: SummaryCardsProps) {
  const periodLabel = PERIOD_LABEL[period]
  const prevLabel = PREV_LABEL[period]
  const today = stats?.today
  const prev = stats?.prev
  const prevAvailable = prev?.available ?? false

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-6">
      <StatCard
        icon={ShoppingBag}
        label={`Transaksi ${periodLabel}`}
        value={String(today?.total_transactions ?? 0)}
        isLoading={isLoading}
        delta={
          <DeltaLine
            current={today?.total_transactions ?? 0}
            previous={prev?.total_transactions ?? 0}
            available={prevAvailable}
            prevLabel={prevLabel}
          />
        }
      />
      <StatCard
        icon={Banknote}
        label={`Pendapatan ${periodLabel}`}
        value={formatRupiah(today?.total_sales ?? 0)}
        isLoading={isLoading}
        delta={
          <DeltaLine
            current={today?.total_sales ?? 0}
            previous={prev?.total_sales ?? 0}
            available={prevAvailable}
            prevLabel={prevLabel}
          />
        }
      />
      <StatCard
        icon={TrendingUp}
        label={`Laba Kotor ${periodLabel}`}
        value={formatRupiah(today?.gross_profit ?? 0)}
        isLoading={isLoading}
        hint="pendapatan − modal, belum potong biaya"
        delta={
          <DeltaLine
            current={today?.gross_profit ?? 0}
            previous={prev?.gross_profit ?? 0}
            available={prevAvailable}
            prevLabel={prevLabel}
          />
        }
      />
      <StatCard
        icon={PackageX}
        label="Stok Habis"
        value={String(stats?.out_of_stock_count ?? 0)}
        isLoading={isLoading}
        accent="red"
        hint="perlu segera dibeli"
      />
      <StatCard
        icon={PackageX}
        label="Stok Menipis"
        value={String(stats?.low_only_count ?? 0)}
        isLoading={isLoading}
        accent="amber"
        hint="mulai perlu diperhatikan"
      />
      <StatCard
        icon={CreditCard}
        label="Piutang Terbuka"
        value={String(stats?.open_receivables ?? 0)}
        isLoading={isLoading}
      />
    </div>
  )
}
