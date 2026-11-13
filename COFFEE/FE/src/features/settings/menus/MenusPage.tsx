import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { MenuTable } from './components/MenuTable'
import type { MenuTableHandle } from './components/MenuTable'

export function MenusPage() {
  const tableRef = useRef<MenuTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Menu"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen Menu' }]}
        actions={
          <RoleGuard menuKey="sistem.menus" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={14} />
              Tambah Menu
            </Button>
          </RoleGuard>
        }
      />
      <MenuTable ref={tableRef} />
    </div>
  )
}
