import { formatRupiah } from '@/shared/utils'

import { useReceivableStatsQuery } from '../receivables.api'

interface CardProps {
  label: string
  value: string
  hint?: string
  accent?: 'default' | 'red'
  isLoading: boolean
}

function Card({ label, value, hint, accent = 'default', isLoading }: CardProps) {
  const border = accent === 'red' ? 'border-red-200 bg-red-50' : 'bg-white'
  const valueColor = accent === 'red' ? 'text-red-700' : 'text-gray-800'
  return (
    <div className={`rounded-lg border p-4 space-y-1 ${border}`}>
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className={`text-xl font-bold ${valueColor}`}>{value}</p>
      )}
      {!isLoading && hint && <p className="text-[11px] leading-tight text-gray-400">{hint}</p>}
    </div>
  )
}

export function ReceivableSummaryCards() {
  const { data, isLoading } = useReceivableStatsQuery()
  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <Card
        label="Total Sisa Piutang"
        value={formatRupiah(data?.total_remaining ?? 0)}
        hint="belum lunas"
        isLoading={isLoading}
      />
      <Card
        label="Pelanggan Berutang"
        value={String(data?.customer_count ?? 0)}
        hint="jumlah pelanggan"
        isLoading={isLoading}
      />
      <Card
        label="Lewat Jatuh Tempo"
        value={formatRupiah(data?.overdue_remaining ?? 0)}
        hint={`${data?.overdue_count ?? 0} piutang perlu ditagih`}
        accent="red"
        isLoading={isLoading}
      />
    </div>
  )
}
