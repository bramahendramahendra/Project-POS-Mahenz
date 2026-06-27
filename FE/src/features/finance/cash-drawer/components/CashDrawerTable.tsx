import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'
import { formatRupiah, monthStart, todayStr } from '@/shared/utils'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card'

import { useCashDrawerListQuery, useCashDrawerSummaryQuery } from '../cash-drawer.api'
import type { CashDrawer, CashDrawerListFilter } from '../cash-drawer.types'
import { CashDrawerFilterBar } from './CashDrawerFilterBar'
import { CashDrawerDetailModal } from './CashDrawerDetailModal'
import { CloseCashDrawerModal } from './CloseCashDrawerModal'
import { buildCashDrawerColumns } from './CashDrawerTableColumns'

interface SummaryCardProps {
  title: string
  value: number
  valueClass?: string
  isLoading?: boolean
}

function SummaryCard({ title, value, valueClass = '', isLoading }: SummaryCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-gray-500">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="h-6 w-32 animate-pulse rounded bg-gray-200" />
        ) : (
          <p className={`text-lg font-semibold ${valueClass}`}>{formatRupiah(value)}</p>
        )}
      </CardContent>
    </Card>
  )
}

const defaultFilter: CashDrawerListFilter = {
  start_date: monthStart(),
  end_date: todayStr(),
}

export function CashDrawerTable() {
  const [filter, setFilter] = useState<CashDrawerListFilter>(defaultFilter)
  const [detailDrawerId, setDetailDrawerId] = useState<number | undefined>(undefined)
  const [forceCloseDrawerId, setForceCloseDrawerId] = useState<number | undefined>(undefined)

  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: forceCloseOpen, open: openForceClose, close: closeForceClose } = useDisclosure()

  const { user } = useAuthStore()
  const canForceClose = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useCashDrawerListQuery({ ...filter, page, limit: pageSize })
  const { data: summary, isLoading: summaryLoading } = useCashDrawerSummaryQuery(filter)
  const items: CashDrawer[] = data?.data ?? []
  const total = data?.total ?? 0

  const handleFilterChange = (newFilter: CashDrawerListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter(defaultFilter)
    reset()
  }

  const handleOpenDetail = (row: CashDrawer) => {
    setDetailDrawerId(row.id)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailDrawerId(undefined)
  }

  const handleOpenForceClose = (row: CashDrawer) => {
    setForceCloseDrawerId(row.id)
    openForceClose()
  }

  const handleCloseForceClose = () => {
    closeForceClose()
    setForceCloseDrawerId(undefined)
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

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <SummaryCard
          title="Total Saldo Awal Tunai"
          value={summary?.total_opening ?? 0}
          isLoading={summaryLoading}
        />
        <SummaryCard
          title="Total Saldo Akhir Tunai"
          value={summary?.total_closing ?? 0}
          isLoading={summaryLoading}
        />
        <SummaryCard
          title="Total Pengeluaran"
          value={summary?.total_expenses ?? 0}
          valueClass="text-red-600"
          isLoading={summaryLoading}
        />
        <SummaryCard
          title="Selisih Bersih"
          value={summary?.net ?? 0}
          valueClass={(summary?.net ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'}
          isLoading={summaryLoading}
        />
      </div>

      <DataTable<CashDrawer & Record<string, unknown>>
        columns={columns}
        data={items as (CashDrawer & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data kas harian"
        emptyDescription="Data kas harian akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <CashDrawerDetailModal
        open={detailOpen}
        onOpenChange={(val) => { if (!val) handleCloseDetail() }}
        cashDrawerId={detailDrawerId}
      />

      <CloseCashDrawerModal
        open={forceCloseOpen}
        onOpenChange={(val) => { if (!val) handleCloseForceClose() }}
        cashDrawerId={forceCloseDrawerId ?? null}
      />
    </div>
  )
}
