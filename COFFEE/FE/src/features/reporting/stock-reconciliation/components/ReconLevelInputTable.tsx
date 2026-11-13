import { Input } from '@/shared/components/ui/input'
import { formatStockNumber } from '@/features/products/products'

import type { ReconPackageLevel } from '../stock-reconciliation.types'

interface ReconLevelInputTableProps {
  packages: ReconPackageLevel[]
  values: Record<number, string>
  onChange: (packageId: number, value: string) => void
  baseUnit: string
}

function parseNum(v: string): number {
  const n = Number(v)
  return Number.isFinite(n) ? n : 0
}

export function ReconLevelInputTable({
  packages,
  values,
  onChange,
  baseUnit,
}: ReconLevelInputTableProps) {
  const totalBase = packages.reduce((sum, p) => {
    const val = parseNum(values[p.package_id] ?? String(p.current_stock))
    return sum + val * p.factor_to_base
  }, 0)

  return (
    <div>
      <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">
        Koreksi Manual{' '}
        <span className="font-normal normal-case text-gray-400">(input stok yang benar per level satuan)</span>
      </p>
      <div className="overflow-hidden rounded-lg border">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr className="text-[11px] uppercase tracking-wide text-gray-500">
              <th className="px-3 py-2 text-left font-semibold">Level</th>
              <th className="px-3 py-2 text-right font-semibold">Stok Sekarang</th>
              <th className="px-3 py-2 text-right font-semibold">Stok Benar</th>
              <th className="px-3 py-2 text-right font-semibold">Setara Dasar</th>
            </tr>
          </thead>
          <tbody>
            {packages.map((p) => {
              const val = values[p.package_id] ?? String(p.current_stock)
              const eqv = parseNum(val) * p.factor_to_base
              return (
                <tr key={p.package_id} className="border-t border-gray-100">
                  <td className="px-3 py-2 text-sm">
                    <span className="font-medium">{p.unit_name}</span>
                    {p.package_name && <span className="ml-1 text-xs text-gray-400">({p.package_name})</span>}
                    {p.is_default ? (
                      <span className="ml-1 text-xs font-semibold text-sky-600">(dasar)</span>
                    ) : (
                      <span className="ml-1 text-xs text-gray-400">(×{formatStockNumber(p.factor_to_base)})</span>
                    )}
                  </td>
                  <td className="px-3 py-2 text-right text-sm tabular-nums text-gray-500">
                    {formatStockNumber(p.current_stock)}
                  </td>
                  <td className="px-3 py-2 text-right">
                    <Input
                      type="number"
                      min={0}
                      step="any"
                      value={val}
                      onChange={(e) => onChange(p.package_id, e.target.value)}
                      className="ml-auto h-8 w-28 text-right tabular-nums"
                    />
                  </td>
                  <td className="px-3 py-2 text-right text-xs tabular-nums text-gray-500">
                    {formatStockNumber(eqv)}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        <div className="flex items-center justify-between border-t-2 bg-gray-50 px-3 py-2.5 text-sm">
          <span>Total setara satuan dasar (baru)</span>
          <b className="text-base tabular-nums">
            {formatStockNumber(totalBase)} {baseUnit}
          </b>
        </div>
      </div>
    </div>
  )
}
