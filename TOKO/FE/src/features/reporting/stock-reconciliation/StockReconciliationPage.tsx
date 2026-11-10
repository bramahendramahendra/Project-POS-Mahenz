import { useState } from 'react'

import { DataTable, PageHeader } from '@/shared/components'
import { usePageSizeOptions, usePagination } from '@/shared/hooks'

import { useReconListQuery, useReconSummaryQuery } from './stock-reconciliation.api'
import type { ReconFilter, ReconListItem } from './stock-reconciliation.types'
import { ReconDetailModal } from './components/ReconDetailModal'
import { ReconFilterBar } from './components/ReconFilterBar'
import { ReconSummaryCard } from './components/ReconSummaryCard'
import { buildReconColumns } from './components/ReconTableColumns'

export function StockReconciliationPage() {
  // default: hanya tampilkan produk yang perlu ditinjau
  const [filter, setFilter] = useState<ReconFilter>({ only_review: true })
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [modalOpen, setModalOpen] = useState(false)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: listData, isLoading: listLoading } = useReconListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const { data: summary, isLoading: summaryLoading } = useReconSummaryQuery()

  const items: ReconListItem[] = listData?.data ?? []
  const total = listData?.total ?? 0

  const handleFilterChange = (next: ReconFilter) => {
    setFilter(next)
    reset()
  }

  const handleReset = () => {
    setFilter({ only_review: true })
    reset()
  }

  const handleOpen = (item: ReconListItem) => {
    setSelectedId(item.product_id)
    setModalOpen(true)
  }

  const columns = buildReconColumns(handleOpen)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Rekonsiliasi Stok"
        breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Rekonsiliasi Stok' }]}
      />

      <ReconFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <ReconSummaryCard summary={summary} isLoading={summaryLoading} />

      <DataTable<ReconListItem & Record<string, unknown>>
        columns={columns}
        data={items as (ReconListItem & Record<string, unknown>)[]}
        isLoading={listLoading}
        emptyMessage="Tidak ada produk"
        emptyDescription="Tidak ada produk yang cocok dengan filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <ReconDetailModal productId={selectedId} open={modalOpen} onOpenChange={setModalOpen} />
    </div>
  )
}
