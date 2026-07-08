import { useRef } from 'react'
import { Shield } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { RoleTable } from './components/RoleTable'
import type { RoleTableHandle } from './components/RoleTable'

export function RolesPage() {
  const tableRef = useRef<RoleTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Role"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen Role' }]}
        actions={
          <RoleGuard menuKey="sistem.roles" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Shield size={14} />
              Tambah Role
            </Button>
          </RoleGuard>
        }
      />
      <RoleTable ref={tableRef} />
    </div>
  )
}
