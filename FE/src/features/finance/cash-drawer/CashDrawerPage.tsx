import { useState } from 'react'

import { PageHeader } from '@/shared/components'
import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useCashDrawerCurrentQuery, useCashDrawerListQuery } from './cash-drawer.api'
import type { CashDrawer, CashDrawerListFilter } from './cash-drawer.types'
import { CashDrawerFilterBar } from './components/CashDrawerFilterBar'
import { CashDrawerTable } from './components/CashDrawerTable'
import { CashDrawerDetailModal } from './components/CashDrawerDetailModal'
import { OpenCashDrawerModal } from './components/OpenCashDrawerModal'
import { CloseCashDrawerModal } from './components/CloseCashDrawerModal'

const SHIFT_LABELS: Record<string, string> = {
  pagi: 'Pagi',
  siang: 'Siang',
  malam: 'Malam',
}

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

export function CashDrawerPage() {
  const { user } = useAuthStore()
  const isAdminOrOwner = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const [filter, setFilter] = useState<CashDrawerListFilter>({
    page: 1,
    limit: 10,
    start_date: monthStartString(),
    end_date: todayString(),
  })
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [openModalOpen, setOpenModalOpen] = useState(false)
  const [closeModalOpen, setCloseModalOpen] = useState(false)

  const { page, pageSize, onPageChange, onPageSizeChange } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useCashDrawerListQuery({ ...filter, page, limit: pageSize })
  const { data: currentData } = useCashDrawerCurrentQuery()

  const items: CashDrawer[] = data?.data ?? []
  const total = data?.total ?? 0

  const currentDrawer = currentData ?? null
  const isOpen = currentDrawer?.status === 'open'

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Harian"
        breadcrumbs={[{ label: 'Keuangan' }, { label: 'Kas Harian' }]}
      />

      <Card>
        <CardContent className="pt-4 pb-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <p className="text-sm font-medium text-gray-500">Status Kas Hari Ini</p>
              <div className="flex items-center gap-2">
                {isOpen ? (
                  <Badge variant="default" className="bg-green-600">● Buka</Badge>
                ) : (
                  <Badge variant="secondary">● Tutup</Badge>
                )}
                {isOpen && currentDrawer?.shift && (
                  <span className="text-sm text-gray-500">
                    Shift: {SHIFT_LABELS[currentDrawer.shift] ?? currentDrawer.shift}
                  </span>
                )}
              </div>
            </div>
            <div>
              {isOpen && isAdminOrOwner ? (
                <Button onClick={() => setCloseModalOpen(true)}>Tutup Kas</Button>
              ) : !isOpen ? (
                <Button onClick={() => setOpenModalOpen(true)}>Buka Kas</Button>
              ) : null}
            </div>
          </div>
        </CardContent>
      </Card>

      <CashDrawerFilterBar
        filter={filter}
        onChange={setFilter}
        onReset={() => setFilter({ page: 1, limit: 10, start_date: monthStartString(), end_date: todayString() })}
      />

      <CashDrawerTable
        data={items}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        onRowClick={(row) => setSelectedId(row.id)}
      />

      <CashDrawerDetailModal cashDrawerId={selectedId} onClose={() => setSelectedId(null)} />

      <OpenCashDrawerModal
        open={openModalOpen}
        onOpenChange={(val) => setOpenModalOpen(val)}
      />

      <CloseCashDrawerModal
        open={closeModalOpen}
        onOpenChange={(val) => setCloseModalOpen(val)}
        cashDrawerId={currentDrawer?.id ?? null}
      />
    </div>
  )
}
