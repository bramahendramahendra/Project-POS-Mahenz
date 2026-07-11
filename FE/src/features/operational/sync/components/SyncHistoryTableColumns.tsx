import { formatDateTime } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SyncHistoryItem } from '../sync.types'

const STATUS_STYLE: Record<string, string> = {
  success: 'bg-green-100 text-green-700',
  partial: 'bg-yellow-100 text-yellow-700',
  failed: 'bg-red-100 text-red-700',
}

const STATUS_LABEL: Record<string, string> = {
  success: 'Sukses',
  partial: 'Sebagian',
  failed: 'Gagal',
}

export function buildSyncHistoryColumns(): ColumnDef<SyncHistoryItem>[] {
  return [
    {
      key: 'started_at',
      header: 'Waktu Sync',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDateTime(row.started_at)}</span>
      ),
    },
    {
      key: 'device_id',
      header: 'Perangkat',
      cell: (row) => (
        <span className="font-medium text-sm">
          {row.device_id}
          <span className="ml-1.5 text-xs text-gray-400 font-normal">({row.device_type})</span>
        </span>
      ),
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
      key: 'synced_items',
      header: 'Tersync',
      align: 'right',
      cell: (row) => <span className="text-green-600 font-medium">{row.synced_items}</span>,
    },
    {
      key: 'conflict_items',
      header: 'Konflik',
      align: 'right',
      cell: (row) => (
        <span className={row.conflict_items > 0 ? 'text-orange-600 font-medium' : 'text-gray-400'}>
          {row.conflict_items}
        </span>
      ),
    },
    {
      key: 'failed_items',
      header: 'Gagal',
      align: 'right',
      cell: (row) => (
        <span className={row.failed_items > 0 ? 'text-red-600 font-medium' : 'text-gray-400'}>
          {row.failed_items}
        </span>
      ),
    },
    {
      key: 'pending_items',
      header: 'Menunggu',
      align: 'right',
      cell: (row) => (
        <span className={row.pending_items > 0 ? 'text-blue-600 font-medium' : 'text-gray-400'}>
          {row.pending_items}
        </span>
      ),
    },
    {
      key: 'duration_ms',
      header: 'Durasi',
      align: 'right',
      cell: (row) => <span className="text-xs text-gray-500">{row.duration_ms} ms</span>,
    },
  ]
}
