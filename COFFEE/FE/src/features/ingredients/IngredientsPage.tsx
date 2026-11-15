import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, PrerequisiteGuard, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { IngredientTable, type IngredientTableHandle } from './components/IngredientTable'
import { useIngredientPrerequisites } from './hooks/useIngredientPrerequisites'

// Daftar Bahan (tahap 4.3). Stok = pembelian - terpakai - terbuang, tanpa opname (keputusan 84).
export function IngredientsPage() {
  const tableRef = useRef<IngredientTableHandle>(null)
  const { isLoading, items } = useIngredientPrerequisites()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Daftar Bahan"
        description="Stok bertambah dari pembelian dan berkurang dari pemakaian atau barang terbuang yang dicatat."
        breadcrumbs={[{ label: 'Bahan Baku' }, { label: 'Daftar Bahan' }]}
        actions={
          <RoleGuard menuKey="bahan.daftar" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Bahan
            </Button>
          </RoleGuard>
        }
      />
      <PrerequisiteGuard
        isLoading={isLoading}
        title="Belum bisa menambah bahan"
        description="Sebelum menambah bahan, pastikan data berikut sudah tersedia:"
        items={items}
      >
        <IngredientTable ref={tableRef} />
      </PrerequisiteGuard>
    </div>
  )
}
