import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { IngredientTable, type IngredientTableHandle } from './components/IngredientTable'

// Daftar Bahan (tahap 4.3). Stok = pembelian - terpakai - terbuang, tanpa opname (keputusan 84).
export function IngredientsPage() {
  const tableRef = useRef<IngredientTableHandle>(null)

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
      <IngredientTable ref={tableRef} />
    </div>
  )
}
