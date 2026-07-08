import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, PrerequisiteGuard, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { PurchaseTable } from './components/PurchaseTable'
import type { PurchaseTableHandle } from './components/PurchaseTable'
import { usePurchasePrerequisites } from './hooks/usePurchasePrerequisites'

export function PurchasesPage() {
  const tableRef = useRef<PurchaseTableHandle>(null)
  const { isLoading, items } = usePurchasePrerequisites()

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
      <PrerequisiteGuard
        isLoading={isLoading}
        title="Belum bisa menambah pembelian"
        description="Sebelum menambah pembelian, pastikan data berikut sudah tersedia:"
        items={items}
      >
        <PurchaseTable ref={tableRef} />
      </PrerequisiteGuard>
    </div>
  )
}
