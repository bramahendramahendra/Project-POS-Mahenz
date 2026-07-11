import { formatDateTime } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SyncQueueItem } from '../sync.types'

const STATUS_STYLE: Record<string, string> = {
  pending: 'bg-gray-100 text-gray-700',
  syncing: 'bg-blue-100 text-blue-700',
  synced: 'bg-green-100 text-green-700',
  failed: 'bg-red-100 text-red-700',
}

const STATUS_LABEL: Record<string, string> = {
  pending: 'Menunggu',
  syncing: 'Sinkronisasi...',
  synced: 'Tersync',
  failed: 'Gagal',
}

const ACTION_LABEL: Record<string, string> = {
  create: 'Tambah',
  update: 'Ubah',
  delete: 'Hapus',
}

const ENTITY_TYPE_LABEL: Record<string, string> = {
  product: 'Produk',
  transaction: 'Transaksi',
  expense: 'Pengeluaran',
  cash_drawer: 'Kas Harian',
}

export function buildSyncQueueColumns(): ColumnDef<SyncQueueItem>[] {
  return [
    {
      key: 'created_at',
      header: 'Waktu',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDateTime(row.created_at)}</span>
      ),
    },
    {
      key: 'device_id',
      header: 'Perangkat',
      cell: (row) => <span className="font-medium text-sm">{row.device_id}</span>,
    },
    {
      key: 'entity_type',
      header: 'Entitas',
      cell: (row) => (
        <span className="text-sm">
          {ENTITY_TYPE_LABEL[row.entity_type] ?? row.entity_type}
          {row.entity_id ? <span className="text-gray-400"> #{row.entity_id}</span> : null}
        </span>
      ),
    },
    {
      key: 'action',
      header: 'Aksi',
      cell: (row) => <span className="text-sm">{ACTION_LABEL[row.action] ?? row.action}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      align: 'center',
      cell: (row) => (
        <span
          className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${STATUS_STYLE[row.status] ?? ''}`}
        >
          {STATUS_LABEL[row.status] ?? row.status}
        </span>
      ),
    },
    {
      key: 'retry_count',
      header: 'Percobaan',
      align: 'right',
      cell: (row) => (
        <span className={row.retry_count > 0 ? 'text-orange-600 font-medium' : 'text-gray-400'}>
          {row.retry_count}
        </span>
      ),
    },
  ]
}
