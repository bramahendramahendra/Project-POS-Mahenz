import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, PrerequisiteGuard, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { ReturnTable } from './components/ReturnTable'
import type { ReturnTableHandle } from './components/ReturnTable'
import { useReturnPrerequisites } from './hooks/useReturnPrerequisites'

export function ReturnsPage() {
  const tableRef = useRef<ReturnTableHandle>(null)
  const { isLoading, items } = useReturnPrerequisites()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Retur Pembelian"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Retur' }]}
        actions={
          <RoleGuard menuKey="pengadaan.retur" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Retur
            </Button>
          </RoleGuard>
        }
      />
      <PrerequisiteGuard
        isLoading={isLoading}
        title="Belum bisa menambah retur"
        description="Retur pembelian dibuat berdasarkan faktur pembelian yang sudah ada."
        items={items}
      >
        <ReturnTable ref={tableRef} />
      </PrerequisiteGuard>
    </div>
  )
}
