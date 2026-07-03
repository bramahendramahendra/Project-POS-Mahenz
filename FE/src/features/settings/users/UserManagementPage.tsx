import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { UserTable } from './components/UserTable'
import type { UserTableHandle } from './components/UserTable'

export function UserManagementPage() {
  const tableRef = useRef<UserTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen User"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen User' }]}
        actions={
          <RoleGuard menuKey="sistem.users" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah User
            </Button>
          </RoleGuard>
        }
      />
      <UserTable ref={tableRef} />
    </div>
  )
}
