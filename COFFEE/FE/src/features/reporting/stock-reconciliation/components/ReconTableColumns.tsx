import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { formatStockNumber } from '@/features/products/products'

import type { ReconListItem } from '../stock-reconciliation.types'

// File ini adalah builder kolom tabel (buildReconColumns), bukan modul komponen
// untuk Fast Refresh. StatusBadge cuma helper internal yang dipakai di dalam
// definisi kolom, jadi aturan react-refresh/only-export-components tidak relevan.
// eslint-disable-next-line react-refresh/only-export-components
function StatusBadge({ item }: { item: ReconListItem }) {
  if (item.needs_stock_review) {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-semibold text-amber-700">
        <span className="h-1.5 w-1.5 rounded-full bg-amber-500" />
        Perlu Ditinjau
      </span>
    )
  }
  return (
    <span className="inline-flex items-center gap-1 rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-semibold text-emerald-700">
      <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
      Cocok
    </span>
  )
}

export function buildReconColumns(onOpen: (item: ReconListItem) => void): ColumnDef<ReconListItem>[] {
  return [
    {
      key: 'product_code',
      header: 'Kode',
      mobileHidden: true,
      cell: (r) => <span className="text-xs font-mono text-gray-500">{r.product_code || '-'}</span>,
    },
    {
      key: 'product_name',
      header: 'Produk',
      mobileLabel: true,
      cell: (r) => <span className="text-sm font-medium">{r.product_name}</span>,
    },
    {
      key: 'old_stock',
      header: 'Stok Lama',
      align: 'right',
      cell: (r) =>
        r.old_stock_available ? (
          <span className="text-sm tabular-nums">{formatStockNumber(r.old_stock)}</span>
        ) : (
          <span className="text-sm text-gray-400">-</span>
        ),
    },
    {
      key: 'new_stock',
      header: 'Stok Baru',
      align: 'right',
      cell: (r) => <span className="text-sm tabular-nums">{formatStockNumber(r.new_stock)}</span>,
    },
    {
      key: 'diff',
      header: 'Selisih',
      align: 'right',
      cell: (r) => {
        if (!r.old_stock_available) return <span className="text-sm text-gray-400">-</span>
        if (r.diff === 0) return <span className="text-sm tabular-nums text-gray-400">0</span>
        return (
          <span className="text-sm font-semibold tabular-nums text-red-600">
            {r.diff > 0 ? '+' : ''}
            {formatStockNumber(r.diff)}
          </span>
        )
      },
    },
    {
      key: 'status',
      header: 'Status',
      cell: (r) => <StatusBadge item={r} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'right',
      cell: (r) => (
        <button
          type="button"
          onClick={() => onOpen(r)}
          className="text-sm font-semibold text-sky-600 hover:text-sky-700"
        >
          Tinjau ›
        </button>
      ),
    },
  ]
}
