import { formatDateTime } from '@/shared/utils'
import { formatStockNumber } from '@/features/products/products'

import type { ReconKartuStokRow } from '../stock-reconciliation.types'

const TYPE_LABEL: Record<string, { label: string; cls: string }> = {
  in: { label: 'Pembelian', cls: 'bg-emerald-100 text-emerald-700' },
  void_purchase: { label: 'Void Beli', cls: 'bg-orange-100 text-orange-700' },
  out: { label: 'Penjualan', cls: 'bg-red-100 text-red-700' },
  void: { label: 'Void Jual', cls: 'bg-sky-100 text-sky-700' },
  return: { label: 'Retur', cls: 'bg-orange-100 text-orange-700' },
  expired: { label: 'Kadaluarsa', cls: 'bg-red-100 text-red-700' },
  adjustment: { label: 'Koreksi', cls: 'bg-violet-100 text-violet-700' },
}

// Jenis mutasi yang MENAMBAH stok -> qty ditampilkan bertanda +
const ADD_TYPES = new Set(['in', 'void'])

function TypeTag({ type }: { type: string }) {
  const t = TYPE_LABEL[type] ?? { label: type, cls: 'bg-gray-100 text-gray-600' }
  return <span className={`rounded px-1.5 py-0.5 text-xs font-semibold ${t.cls}`}>{t.label}</span>
}

export function ReconKartuStok({ rows }: { rows: ReconKartuStokRow[] }) {
  return (
    <div>
      <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">
        Kartu Stok <span className="font-normal normal-case text-gray-400">(riwayat mutasi, terbaru dulu)</span>
      </p>
      <div className="overflow-hidden rounded-lg border">
        <div className="max-h-64 overflow-y-auto">
          <table className="w-full">
            <thead className="sticky top-0 bg-gray-50">
              <tr className="text-[11px] uppercase tracking-wide text-gray-500">
                <th className="px-3 py-2 text-left font-semibold">Tanggal</th>
                <th className="px-3 py-2 text-left font-semibold">Jenis</th>
                <th className="px-3 py-2 text-right font-semibold">Qty</th>
                <th className="px-3 py-2 text-right font-semibold">Sebelum</th>
                <th className="px-3 py-2 text-right font-semibold">Sesudah</th>
                <th className="px-3 py-2 text-left font-semibold">Ref</th>
                <th className="px-3 py-2 text-left font-semibold">Oleh</th>
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-3 py-6 text-center text-sm text-gray-400">
                    Belum ada riwayat mutasi
                  </td>
                </tr>
              ) : (
                rows.map((r) => {
                  const isAdd = ADD_TYPES.has(r.mutation_type)
                  return (
                    <tr key={r.id} className="border-t border-gray-100 text-xs">
                      <td className="whitespace-nowrap px-3 py-2">{formatDateTime(r.created_at)}</td>
                      <td className="px-3 py-2">
                        <TypeTag type={r.mutation_type} />
                      </td>
                      <td className={`px-3 py-2 text-right tabular-nums ${isAdd ? 'text-emerald-600' : 'text-red-600'}`}>
                        {isAdd ? '+' : '-'}
                        {formatStockNumber(r.quantity)}
                      </td>
                      <td className="px-3 py-2 text-right tabular-nums text-gray-500">{formatStockNumber(r.stock_before)}</td>
                      <td className="px-3 py-2 text-right tabular-nums text-gray-500">{formatStockNumber(r.stock_after)}</td>
                      <td className="px-3 py-2">
                        {r.reference_type ? (
                          <span className="text-sky-600">
                            {r.reference_type}
                            {r.reference_id ? ` #${r.reference_id}` : ''}
                          </span>
                        ) : (
                          <span className="text-gray-400">-</span>
                        )}
                      </td>
                      <td className="px-3 py-2 text-gray-500">{r.user_name || '-'}</td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
