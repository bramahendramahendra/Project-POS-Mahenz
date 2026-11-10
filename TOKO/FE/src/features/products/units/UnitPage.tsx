import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { UnitTable } from './components/UnitTable'
import type { UnitTableHandle } from './components/UnitTable'

export function UnitPage() {
  const tableRef = useRef<UnitTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Satuan Produk"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Satuan' }]}
        actions={
          <RoleGuard menuKey="produk.unit" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Satuan
            </Button>
          </RoleGuard>
        }
      />
      <UnitTable ref={tableRef} />
    </div>
  )
}
