import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/shared/components/ui/collapsible'
import { formatDateTime } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SyncHistoryItem } from '../sync.types'

const STATUS_STYLE: Record<string, string> = {
  success: 'bg-green-100 text-green-700',
  partial: 'bg-yellow-100 text-yellow-700',
  error: 'bg-red-100 text-red-700',
}

const STATUS_LABEL: Record<string, string> = {
  success: 'Sukses',
  partial: 'Sebagian',
  error: 'Error',
}

interface SyncHistoryColumnHandlers {
  expandedId: number | null
  onToggleExpand: (id: number) => void
}

export function buildSyncHistoryColumns({ expandedId, onToggleExpand }: SyncHistoryColumnHandlers): ColumnDef<SyncHistoryItem>[] {
  return [
    {
      key: 'created_at',
      header: 'Waktu Sync',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDateTime(row.created_at)}</span>
      ),
    },
    {
      key: 'device_info',
      header: 'Perangkat',
      cell: (row) => <span className="font-medium text-sm">{row.device_info}</span>,
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
      key: 'synced_count',
      header: 'Tersync',
      align: 'right',
      cell: (row) => <span className="text-green-600 font-medium">{row.synced_count}</span>,
    },
    {
      key: 'error_count',
      header: 'Error',
      align: 'right',
      cell: (row) => (
        <span className={row.error_count > 0 ? 'text-red-600 font-medium' : 'text-gray-400'}>
          {row.error_count}
        </span>
      ),
    },
    {
      key: 'message',
      header: 'Pesan',
      cell: (row) =>
        row.message ? (
          <Collapsible open={expandedId === row.id} onOpenChange={() => onToggleExpand(row.id)}>
            <CollapsibleTrigger className="text-xs text-blue-600 hover:underline">
              {expandedId === row.id ? 'Sembunyikan' : 'Lihat pesan'}
            </CollapsibleTrigger>
            <CollapsibleContent>
              <p className="mt-1 text-xs text-gray-600 bg-gray-50 rounded p-2 max-w-xs">
                {row.message}
              </p>
            </CollapsibleContent>
          </Collapsible>
        ) : (
          <span className="text-xs text-gray-400">—</span>
        ),
    },
  ]
}
