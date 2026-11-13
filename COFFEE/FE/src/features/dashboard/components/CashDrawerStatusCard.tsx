import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { PrerequisiteGuard } from '@/shared/components'
import { ROUTES } from '@/shared/constants/routes'
import { useCashDrawerCurrentQuery, OpenCashDrawerModal } from '@/features/finance/cash-drawer'
import { useCashDrawerPrerequisites } from '@/features/finance/cash-drawer/hooks/useCashDrawerPrerequisites'

export function CashDrawerStatusCard() {
  const navigate = useNavigate()
  const { data: currentDrawer, isLoading } = useCashDrawerCurrentQuery()
  const { items: prerequisiteItems, isLoading: prerequisiteLoading } = useCashDrawerPrerequisites()
  const [openModalOpen, setOpenModalOpen] = useState(false)

  const isOpen = currentDrawer?.status === 'open'

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-4 pb-4">
          <div className="h-14 animate-pulse rounded-lg bg-gray-100" />
        </CardContent>
      </Card>
    )
  }

  if (isOpen) {
    return (
      <Card>
        <CardContent className="pt-4 pb-4 flex items-center justify-between flex-wrap gap-3">
          <div className="flex items-center gap-2">
            <Badge variant="default" className="bg-green-600">● Kas Buka</Badge>
            {currentDrawer?.shift_name && (
              <span className="text-sm text-gray-500">
                {currentDrawer.shift_name}
                {currentDrawer.shift_start && currentDrawer.shift_end
                  ? ` (${currentDrawer.shift_start} – ${currentDrawer.shift_end})`
                  : ''}
              </span>
            )}
          </div>
          <Button size="sm" onClick={() => navigate(ROUTES.KASIR)}>
            Mulai Transaksi
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <PrerequisiteGuard
      isLoading={prerequisiteLoading}
      title="Belum bisa membuka kas"
      description="Sebelum membuka kas, pastikan data berikut sudah tersedia:"
      items={prerequisiteItems}
    >
      <Card>
        <CardContent className="pt-4 pb-4 flex items-center justify-between flex-wrap gap-3">
          <Badge variant="secondary">● Kas Belum Dibuka</Badge>
          <Button size="sm" onClick={() => setOpenModalOpen(true)}>
            Buka Kas
          </Button>
        </CardContent>
      </Card>

      <OpenCashDrawerModal open={openModalOpen} onOpenChange={setOpenModalOpen} />
    </PrerequisiteGuard>
  )
}
