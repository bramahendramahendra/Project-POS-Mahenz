import { useState } from 'react'
import { Clock, Plus } from 'lucide-react'

import { PageHeader } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useActiveShiftQuery, useShiftListQuery } from './shifts.api'
import type { Shift, ShiftListFilter, ShiftStatus } from './shifts.types'
import { formatDateTime } from './shifts.utils'
import { CloseShiftModal } from './components/CloseShiftModal'
import { OpenShiftModal } from './components/OpenShiftModal'
import { ShiftFilterBar } from './components/ShiftFilterBar'
import { ShiftTable } from './components/ShiftTable'

export function ShiftsPage() {
  const [status, setStatus] = useState<ShiftStatus | 'all'>('all')
  const [openShiftOpen, setOpenShiftOpen] = useState(false)
  const [closeTarget, setCloseTarget] = useState<Shift | null>(null)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()
  const { data: activeShift } = useActiveShiftQuery()

  const filter: ShiftListFilter = {
    page,
    limit: pageSize,
    status: status === 'all' ? undefined : status,
  }

  const { data, isLoading } = useShiftListQuery(filter)
  const shifts = data?.data ?? []
  const total = data?.total ?? 0

  const handleStatusChange = (val: ShiftStatus | 'all') => {
    setStatus(val)
    reset()
  }

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Shift"
        breadcrumbs={[{ label: 'Shifts' }]}
        actions={
          !activeShift ? (
            <Button onClick={() => setOpenShiftOpen(true)} className="gap-1">
              <Plus size={16} />
              Buka Shift
            </Button>
          ) : undefined
        }
      />

      {activeShift && (
        <div className="flex items-center justify-between rounded-xl border border-green-200 bg-green-50 px-4 py-3">
          <div className="flex items-center gap-2 text-green-700 text-sm">
            <Clock size={16} />
            <span>
              Shift sedang berjalan sejak{' '}
              <span className="font-semibold">{formatDateTime(activeShift.opened_at)}</span> oleh{' '}
              {activeShift.kasir_name}
            </span>
          </div>
          <Button
            size="sm"
            variant="outline"
            className="text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
            onClick={() => setCloseTarget(activeShift)}
          >
            Tutup Shift
          </Button>
        </div>
      )}

      <ShiftFilterBar
        status={status}
        onChange={handleStatusChange}
        onReset={() => { setStatus('all'); reset() }}
      />

      <ShiftTable
        data={shifts}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        onClose={(shift) => setCloseTarget(shift)}
      />

      <OpenShiftModal open={openShiftOpen} onOpenChange={setOpenShiftOpen} />
      <CloseShiftModal
        open={!!closeTarget}
        onOpenChange={(open) => { if (!open) setCloseTarget(null) }}
        shift={closeTarget}
      />
    </div>
  )
}
