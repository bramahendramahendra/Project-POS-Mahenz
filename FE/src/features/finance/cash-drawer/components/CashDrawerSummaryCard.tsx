import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card'
import { formatRupiah } from '@/shared/utils'

import type { CashDrawerSummary } from '../cash-drawer.types'

interface CashDrawerSummaryCardProps {
  summary: CashDrawerSummary | undefined
  isLoading: boolean
}

interface SummaryCardProps {
  title: string
  value: number
  valueClass?: string
  isLoading: boolean
}

function SummaryCard({ title, value, valueClass = '', isLoading }: SummaryCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-gray-500">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="h-6 w-32 animate-pulse rounded bg-gray-200" />
        ) : (
          <p className={`text-lg font-semibold ${valueClass}`}>{formatRupiah(value)}</p>
        )}
      </CardContent>
    </Card>
  )
}

export function CashDrawerSummaryCard({ summary, isLoading }: CashDrawerSummaryCardProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      <SummaryCard
        title="Total Saldo Awal Tunai"
        value={summary?.total_opening ?? 0}
        isLoading={isLoading}
      />
      <SummaryCard
        title="Total Saldo Akhir Tunai"
        value={summary?.total_closing ?? 0}
        isLoading={isLoading}
      />
      <SummaryCard
        title="Total Pengeluaran"
        value={summary?.total_expenses ?? 0}
        valueClass="text-red-600"
        isLoading={isLoading}
      />
      <SummaryCard
        title="Selisih Bersih"
        value={summary?.net ?? 0}
        valueClass={(summary?.net ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'}
        isLoading={isLoading}
      />
    </div>
  )
}
