import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useSyncQueueQuery } from '../sync.api'
import type { SyncQueueItem } from '../sync.types'
import { buildSyncQueueColumns } from './SyncQueueTableColumns'

export function SyncQueueTable() {
  const { page, pageSize, onPageChange, onPageSizeChange } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useSyncQueueQuery({ page, limit: pageSize })
  const rows = data?.data ?? []
  const total = data?.total ?? 0

  const columns = buildSyncQueueColumns()

  return (
    <DataTable<SyncQueueItem & Record<string, unknown>>
      columns={columns}
      data={rows as (SyncQueueItem & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Antrian sync kosong"
      emptyDescription="Item yang menunggu diproses akan muncul di sini."
      pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
    />
  )
}
