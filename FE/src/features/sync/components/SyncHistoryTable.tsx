import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useSyncHistoryQuery } from '../sync.api'
import type { SyncHistoryItem } from '../sync.types'
import { buildSyncHistoryColumns } from './SyncHistoryTableColumns'

export function SyncHistoryTable() {
  const [expandedId, setExpandedId] = useState<number | null>(null)
  const { page, pageSize, onPageChange, onPageSizeChange } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useSyncHistoryQuery({ page, page_size: pageSize })
  const rows = data?.data ?? []
  const total = data?.total ?? 0

  const handleToggleExpand = (id: number) => {
    setExpandedId((prev) => (prev === id ? null : id))
  }

  const columns = buildSyncHistoryColumns({ expandedId, onToggleExpand: handleToggleExpand })

  return (
    <DataTable<SyncHistoryItem & Record<string, unknown>>
      columns={columns}
      data={rows as (SyncHistoryItem & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada riwayat sinkronisasi"
      emptyDescription="Riwayat akan muncul setelah perangkat melakukan sinkronisasi."
      pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
    />
  )
}
