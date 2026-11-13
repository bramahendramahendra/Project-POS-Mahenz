import type { ReconSummary } from '../stock-reconciliation.types'

interface ReconSummaryCardProps {
  summary: ReconSummary | undefined
  isLoading: boolean
}

interface CardProps {
  label: string
  value: string
  isLoading: boolean
  tone?: 'default' | 'warn' | 'ok'
}

function Card({ label, value, isLoading, tone = 'default' }: CardProps) {
  const toneClass =
    tone === 'warn' ? 'text-amber-600' : tone === 'ok' ? 'text-emerald-600' : 'text-gray-800'
  return (
    <div className="space-y-1 rounded-lg border bg-white p-4">
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-20 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className={`text-xl font-bold ${toneClass}`}>{value}</p>
      )}
    </div>
  )
}

export function ReconSummaryCard({ summary, isLoading }: ReconSummaryCardProps) {
  return (
    <div className="space-y-2">
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <Card label="Total Produk" value={String(summary?.total_products ?? 0)} isLoading={isLoading} />
        <Card label="Perlu Ditinjau" value={String(summary?.needs_review ?? 0)} isLoading={isLoading} tone="warn" />
        <Card label="Cocok" value={String(summary?.matched ?? 0)} isLoading={isLoading} tone="ok" />
      </div>
      {summary && !summary.backup_available && (
        <p className="text-xs text-amber-600">
          Catatan: tabel stok lama (products_stock_backup) tidak tersedia — kolom Stok Lama &amp;
          Selisih ditampilkan sebagai &quot;-&quot;.
        </p>
      )}
    </div>
  )
}
