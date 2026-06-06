import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { SupplierTable } from './components/SupplierTable'
import type { SupplierTableHandle } from './components/SupplierTable'

export function SuppliersPage() {
  const tableRef = useRef<SupplierTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Supplier"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Supplier' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Supplier
            </Button>
          </RoleGuard>
        }
      />
      <SupplierTable ref={tableRef} />
    </div>
  )
}
