import { useState } from 'react'

import { DataTable, PageHeader } from '@/shared/components'
import { usePageSizeOptions, usePagination } from '@/shared/hooks'

import {
  useCapitalReconListQuery,
  useCapitalReconSummaryQuery,
} from './capital-reconciliation.api'
import type { CapitalReconFilter, CapitalReconListItem } from './capital-reconciliation.types'
import { CapitalReconFilterBar } from './components/CapitalReconFilterBar'
import { CapitalReconSummaryCard } from './components/CapitalReconSummaryCard'
import { buildCapitalReconColumns } from './components/CapitalReconTableColumns'
import { ResolveCostModal } from './components/ResolveCostModal'

export function CapitalReconciliationPage() {
  const [filter, setFilter] = useState<CapitalReconFilter>({ status: 'pending' })
  const [selected, setSelected] = useState<CapitalReconListItem | null>(null)
  const [modalOpen, setModalOpen] = useState(false)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: listData, isLoading: listLoading } = useCapitalReconListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const { data: summary, isLoading: summaryLoading } = useCapitalReconSummaryQuery()

  const items: CapitalReconListItem[] = listData?.data ?? []
  const total = listData?.total ?? 0

  const handleFilterChange = (next: CapitalReconFilter) => {
    setFilter(next)
    reset()
  }

  const handleReset = () => {
    setFilter({ status: 'pending' })
    reset()
  }

  const handleOpen = (item: CapitalReconListItem) => {
    setSelected(item)
    setModalOpen(true)
  }

  const columns = buildCapitalReconColumns(handleOpen)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Rekonsiliasi Modal"
        breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Rekonsiliasi Modal' }]}
      />

      <CapitalReconSummaryCard summary={summary} isLoading={summaryLoading} />

      <CapitalReconFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<CapitalReconListItem & Record<string, unknown>>
        columns={columns}
        data={items as (CapitalReconListItem & Record<string, unknown>)[]}
        isLoading={listLoading}
        emptyMessage="Tidak ada baris"
        emptyDescription="Tidak ada baris modal yang cocok dengan filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <ResolveCostModal item={selected} open={modalOpen} onOpenChange={setModalOpen} />
    </div>
  )
}
