import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { CustomerTable } from './components/CustomerTable'
import type { CustomerTableHandle } from './components/CustomerTable'

export function CustomersPage() {
  const tableRef = useRef<CustomerTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Pelanggan"
        breadcrumbs={[{ label: 'Pelanggan' }]}
        actions={
          <RoleGuard menuKey="pelanggan.pelanggan" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Pelanggan
            </Button>
          </RoleGuard>
        }
      />
      <CustomerTable ref={tableRef} />
    </div>
  )
}
