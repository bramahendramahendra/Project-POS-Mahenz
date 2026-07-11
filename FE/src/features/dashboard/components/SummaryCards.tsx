import { Banknote, CreditCard, PackageX, ShoppingBag, TrendingUp } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

import { formatRupiah } from '@/shared/utils'

import type { DashboardPeriod, DashboardStats } from '../dashboard.types'

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

interface StatCardProps {
  icon: LucideIcon
  label: string
  value: string
  isLoading: boolean
}

function StatCard({ icon: Icon, label, value, isLoading }: StatCardProps) {
  return (
    <div className="rounded-lg border bg-white p-4 shadow-sm space-y-1">
      <div className="flex items-center gap-2 text-gray-500 text-sm">
        <Icon size={15} />
        <span>{label}</span>
      </div>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className="text-xl font-bold text-gray-900">{value}</p>
      )}
    </div>
  )
}

export function SummaryCards({ stats, isLoading, period }: SummaryCardsProps) {
  const periodLabel = PERIOD_LABEL[period]
  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-5">
      <StatCard
        icon={ShoppingBag}
        label={`Transaksi ${periodLabel}`}
        value={String(stats?.today.total_transactions ?? 0)}
        isLoading={isLoading}
      />
      <StatCard
        icon={Banknote}
        label={`Pendapatan ${periodLabel}`}
        value={formatRupiah(stats?.today.total_sales ?? 0)}
        isLoading={isLoading}
      />
      <StatCard
        icon={TrendingUp}
        label={`Laba Kotor ${periodLabel}`}
        value={formatRupiah(stats?.today.gross_profit ?? 0)}
        isLoading={isLoading}
      />
      <StatCard
        icon={PackageX}
        label="Stok Menipis"
        value={String(stats?.low_stock_count ?? 0)}
        isLoading={isLoading}
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
