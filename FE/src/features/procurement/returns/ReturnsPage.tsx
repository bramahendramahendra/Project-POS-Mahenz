import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { ReturnTable } from './components/ReturnTable'
import type { ReturnTableHandle } from './components/ReturnTable'

export function ReturnsPage() {
  const tableRef = useRef<ReturnTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Retur Pembelian"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Retur' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Retur
            </Button>
          </RoleGuard>
        }
      />
      <ReturnTable ref={tableRef} />
    </div>
  )
}
