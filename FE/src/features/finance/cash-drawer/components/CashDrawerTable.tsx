import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'

import { useCashDrawerListQuery } from '../cash-drawer.api'
import type { CashDrawer, CashDrawerListFilter } from '../cash-drawer.types'
import { CashDrawerFilterBar } from './CashDrawerFilterBar'
import { CashDrawerDetailModal } from './CashDrawerDetailModal'
import { CloseCashDrawerModal } from './CloseCashDrawerModal'
import { buildCashDrawerColumns } from './CashDrawerTableColumns'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

const defaultFilter: CashDrawerListFilter = {
  page: 1,
  limit: 10,
  start_date: monthStartString(),
  end_date: todayString(),
}

export function CashDrawerTable() {
  const [filter, setFilter] = useState<CashDrawerListFilter>(defaultFilter)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [forceCloseId, setForceCloseId] = useState<number | null>(null)

  const { user } = useAuthStore()
  const canForceClose = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useCashDrawerListQuery({ ...filter, page, limit: pageSize })
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

  const columns = buildCashDrawerColumns({
    onRowClick: (row) => setSelectedId(row.id),
    onForceClose: (row) => setForceCloseId(row.id),
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

      <DataTable<CashDrawer & Record<string, unknown>>
        columns={columns}
        data={items as (CashDrawer & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data kas harian"
        emptyDescription="Data kas harian akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <CashDrawerDetailModal cashDrawerId={selectedId} onClose={() => setSelectedId(null)} />

      <CloseCashDrawerModal
        open={forceCloseId !== null}
        onOpenChange={(val) => { if (!val) setForceCloseId(null) }}
        cashDrawerId={forceCloseId}
      />
    </div>
  )
}
