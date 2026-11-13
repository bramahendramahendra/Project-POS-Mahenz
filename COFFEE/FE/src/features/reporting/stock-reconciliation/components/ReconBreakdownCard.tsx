import { formatStockNumber } from '@/features/products/products'

import type { ReconBreakdown } from '../stock-reconciliation.types'

interface ReconBreakdownCardProps {
  breakdown: ReconBreakdown
  baseUnit: string
}

interface RowProps {
  op: '+' | '-' | '±'
  label: string
  value: number
}

function Row({ op, label, value }: RowProps) {
  const opColor = op === '+' ? 'text-emerald-600' : op === '-' ? 'text-red-600' : 'text-gray-500'
  return (
    <div className="flex items-center justify-between border-b border-gray-100 px-3 py-2 text-sm last:border-b-0">
      <span className="flex items-center gap-2 text-gray-700">
        <span className={`w-4 text-center font-bold ${opColor}`}>{op}</span>
        {label}
      </span>
      <span className="font-medium tabular-nums">{formatStockNumber(value)}</span>
    </div>
  )
}

export function ReconBreakdownCard({ breakdown, baseUnit }: ReconBreakdownCardProps) {
  return (
    <div>
      <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">
        Rincian Perhitungan Stok Baru{' '}
        <span className="font-normal normal-case text-gray-400">(satuan dasar: {baseUnit || '-'})</span>
      </p>
      <div className="overflow-hidden rounded-lg border">
        <Row op="+" label="Pembelian (in)" value={breakdown.purchase_in} />
        <Row op="-" label="Pembelian dibatalkan" value={breakdown.purchase_void} />
        <Row op="-" label="Penjualan (out)" value={breakdown.sale_out} />
        <Row op="+" label="Penjualan dibatalkan" value={breakdown.sale_void} />
        <Row op="-" label="Retur supplier" value={breakdown.supplier_return} />
        <Row op="-" label="Kadaluarsa/rusak" value={breakdown.expired} />
        <Row op="±" label="Koreksi manual" value={breakdown.adjustment} />
        <div className="flex items-center justify-between border-t-2 bg-gray-50 px-3 py-2.5 text-sm font-bold">
          <span>= Stok Baru</span>
          <span className="text-base tabular-nums">{formatStockNumber(breakdown.computed_new_stock)}</span>
        </div>
      </div>
    </div>
  )
}
