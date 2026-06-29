import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { PurchaseTable } from './components/PurchaseTable'
import type { PurchaseTableHandle } from './components/PurchaseTable'
import { PurchasePrerequisiteGuard } from './components/PurchasePrerequisiteGuard'

export function PurchasesPage() {
  const tableRef = useRef<PurchaseTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Pembelian Supplier"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Pembelian' }]}
        actions={
          <RoleGuard menuKey="pengadaan.pembelian" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Pembelian
            </Button>
          </RoleGuard>
        }
      />
      <PurchasePrerequisiteGuard>
        <PurchaseTable ref={tableRef} />
      </PurchasePrerequisiteGuard>
    </div>
  )
}
