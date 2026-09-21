import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { formatRupiah, formatDate } from '@/shared/utils'

import type { CapitalReconListItem } from '../capital-reconciliation.types'

// eslint-disable-next-line react-refresh/only-export-components
function StatusBadge({ status }: { status: CapitalReconListItem['status'] }) {
  if (status === 'pending') {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-semibold text-amber-700">
        <span className="h-1.5 w-1.5 rounded-full bg-amber-500" />
        Perlu Ditinjau
      </span>
    )
  }
  if (status === 'resolved') {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-emerald-100 px-2 py-0.5 text-xs font-semibold text-emerald-700">
        <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
        Sudah Dikoreksi
      </span>
    )
  }
  return (
    <span className="inline-flex items-center gap-1 rounded-full bg-gray-100 px-2 py-0.5 text-xs font-semibold text-gray-600">
      <span className="h-1.5 w-1.5 rounded-full bg-gray-400" />
      Selesai
    </span>
  )
}

export function buildCapitalReconColumns(
  onOpen: (item: CapitalReconListItem) => void,
): ColumnDef<CapitalReconListItem>[] {
  return [
    {
      key: 'transaction_date',
      header: 'Tanggal',
      mobileHidden: true,
      cell: (r) => <span className="text-xs text-gray-500">{formatDate(r.transaction_date)}</span>,
    },
    {
      key: 'transaction_code',
      header: 'Kode',
      mobileHidden: true,
      cell: (r) => <span className="font-mono text-xs text-gray-500">{r.transaction_code}</span>,
    },
    {
      key: 'product_name',
      header: 'Produk',
      mobileLabel: true,
      cell: (r) => (
        <div>
          <div className="text-sm font-medium">{r.product_name}</div>
          <div className="text-xs text-gray-400">
            {r.quantity} {r.unit}
          </div>
        </div>
      ),
    },
    {
      key: 'sell_price',
      header: 'Harga Jual',
      align: 'right',
      cell: (r) => <span className="text-sm tabular-nums">{formatRupiah(r.sell_price)}</span>,
    },
    {
      key: 'old_purchase_price',
      header: 'Modal Sekarang',
      align: 'right',
      cell: (r) => (
        <span className="text-sm tabular-nums text-red-600">{formatRupiah(r.old_purchase_price)}</span>
      ),
    },
    {
      key: 'reason',
      header: 'Alasan',
      mobileHidden: true,
      cell: (r) => <span className="text-xs text-gray-500">{r.reason}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      cell: (r) => <StatusBadge status={r.status} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'right',
      cell: (r) =>
        r.status === 'pending' ? (
          <button
            type="button"
            onClick={() => onOpen(r)}
            className="text-sm font-semibold text-sky-600 hover:text-sky-700"
          >
            Tinjau ›
          </button>
        ) : (
          <span className="text-xs text-gray-400">
            {r.resolved_price != null ? formatRupiah(r.resolved_price) : '—'}
          </span>
        ),
    },
  ]
}
