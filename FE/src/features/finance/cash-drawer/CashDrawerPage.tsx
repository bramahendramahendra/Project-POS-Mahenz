import { useState } from 'react'

import { PageHeader } from '@/shared/components'
import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'

import { useShiftOptionsQuery } from '@/features/operational/shifts'

import { useCashDrawerCurrentQuery } from './cash-drawer.api'
import { CashDrawerTable } from './components/CashDrawerTable'
import { OpenCashDrawerModal } from './components/OpenCashDrawerModal'
import { CloseCashDrawerModal } from './components/CloseCashDrawerModal'
import { ShiftPrerequisiteGuard } from './components/ShiftPrerequisiteGuard'

export function CashDrawerPage() {
  const { user } = useAuthStore()
  const isAdminOrOwner = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const [openModalOpen, setOpenModalOpen] = useState(false)
  const [closeModalOpen, setCloseModalOpen] = useState(false)

  const { data: currentData } = useCashDrawerCurrentQuery()
  const currentDrawer = currentData ?? null
  const isOpen = currentDrawer?.status === 'open'

  const { data: shiftsRaw } = useShiftOptionsQuery()
  const hasShifts = (shiftsRaw ?? []).length > 0

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
                {isOpen && currentDrawer?.shift_name && (
                  <span className="text-sm text-gray-500">
                    {currentDrawer.shift_name}
                    {currentDrawer.shift_start && currentDrawer.shift_end
                      ? ` (${currentDrawer.shift_start} – ${currentDrawer.shift_end})`
                      : ''}
                  </span>
                )}
              </div>
            </div>
            <div>
              {isOpen && isAdminOrOwner ? (
                <Button onClick={() => setCloseModalOpen(true)}>Tutup Kas</Button>
              ) : !isOpen ? (
                <Button onClick={() => setOpenModalOpen(true)} disabled={!hasShifts}>
                  Buka Kas
                </Button>
              ) : null}
            </div>
          </div>
        </CardContent>
      </Card>

      <ShiftPrerequisiteGuard>
        <CashDrawerTable />
      </ShiftPrerequisiteGuard>

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
