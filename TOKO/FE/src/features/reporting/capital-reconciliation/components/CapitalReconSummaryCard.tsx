import type { CapitalReconSummary } from '../capital-reconciliation.types'

interface Props {
  summary: CapitalReconSummary | undefined
  isLoading: boolean
}

function Card({
  label,
  value,
  isLoading,
  tone = 'default',
}: {
  label: string
  value: string
  isLoading: boolean
  tone?: 'default' | 'warn' | 'ok'
}) {
  const toneClass =
    tone === 'warn' ? 'text-amber-600' : tone === 'ok' ? 'text-emerald-600' : 'text-gray-800'
  return (
    <div className="space-y-1 rounded-lg border bg-white p-4">
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-16 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className={`text-xl font-bold ${toneClass}`}>{value}</p>
      )}
    </div>
  )
}

export function CapitalReconSummaryCard({ summary, isLoading }: Props) {
  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <Card label="Perlu Ditinjau" value={String(summary?.pending ?? 0)} isLoading={isLoading} tone="warn" />
      <Card label="Sudah Dikoreksi" value={String(summary?.resolved ?? 0)} isLoading={isLoading} tone="ok" />
      <Card label="Ditandai Selesai" value={String(summary?.skipped ?? 0)} isLoading={isLoading} />
    </div>
  )
}
