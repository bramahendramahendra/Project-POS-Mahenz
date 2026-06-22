import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ROLES } from '@/shared/constants/roles'

import { ShiftTable } from './components/ShiftTable'
import type { ShiftTableHandle } from './components/ShiftTable'

export function ShiftsPage() {
  const tableRef = useRef<ShiftTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Shift"
        breadcrumbs={[{ label: 'Operasional' }, { label: 'Shift' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Shift
            </Button>
          </RoleGuard>
        }
      />
      <ShiftTable ref={tableRef} />
    </div>
  )
}
