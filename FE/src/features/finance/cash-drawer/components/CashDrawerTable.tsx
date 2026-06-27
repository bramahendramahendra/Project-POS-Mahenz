import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'
import { monthStart, todayStr } from '@/shared/utils'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useCashDrawerListQuery, useCashDrawerSummaryQuery } from '../cash-drawer.api'
import type { CashDrawer, CashDrawerListFilter } from '../cash-drawer.types'
import { CashDrawerFilterBar } from './CashDrawerFilterBar'
import { CashDrawerDetailModal } from './CashDrawerDetailModal'
import { CloseCashDrawerModal } from './CloseCashDrawerModal'
import { CashDrawerSummaryCard } from './CashDrawerSummaryCard'
import { buildCashDrawerColumns } from './CashDrawerTableColumns'

export function CashDrawerTable() {
  const [filter, setFilter] = useState<CashDrawerListFilter>({
    start_date: monthStart(),
    end_date: todayStr(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const [detailDrawerId, setDetailDrawerId] = useState<number | null>(null)
  const [forceCloseDrawerId, setForceCloseDrawerId] = useState<number | null>(null)

  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: forceCloseOpen, open: openForceClose, close: closeForceClose } = useDisclosure()

  const { user } = useAuthStore()
  const canForceClose = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const { data, isLoading } = useCashDrawerListQuery({ ...filter, page, limit: pageSize })
  const { data: summary, isLoading: summaryLoading } = useCashDrawerSummaryQuery(filter)
  const items: CashDrawer[] = data?.data ?? []
  const total = data?.total ?? 0

  const handleFilterChange = (newFilter: CashDrawerListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ start_date: monthStart(), end_date: todayStr() })
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleOpenDetail = (cashDrawer: CashDrawer) => {
    setDetailDrawerId(cashDrawer.id)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailDrawerId(null)
  }

  const handleOpenForceClose = (cashDrawer: CashDrawer) => {
    setForceCloseDrawerId(cashDrawer.id)
    openForceClose()
  }

  const handleCloseForceClose = () => {
    closeForceClose()
    setForceCloseDrawerId(null)
  }

  const columns = buildCashDrawerColumns({
    onRowClick: handleOpenDetail,
    onForceClose: handleOpenForceClose,
    canForceClose,
  })

  return (
    <div className="space-y-4">
      <CashDrawerFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        showKasirFilter={canForceClose}
      />

      <CashDrawerSummaryCard summary={summary} isLoading={summaryLoading} />

      <DataTable<CashDrawer & Record<string, unknown>>
        columns={columns}
        data={items as (CashDrawer & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data kas harian"
        emptyDescription="Data kas harian akan muncul sesuai filter periode yang dipilih."
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions
        }}
      />

      <CashDrawerDetailModal
        open={detailOpen}
        onOpenChange={(val) => { if (!val) handleCloseDetail() }}
        cashDrawerId={detailDrawerId ?? undefined}
      />

      <CloseCashDrawerModal
        open={forceCloseOpen}
        onOpenChange={(val) => { if (!val) handleCloseForceClose() }}
        cashDrawerId={forceCloseDrawerId}
      />
    </div>
  )
}
