import { useRef } from 'react'
import { Plus, Upload } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { ProductTable } from './components/ProductTable'
import type { ProductTableHandle } from './components/ProductTable'

export function ProductsPage() {
  const tableRef = useRef<ProductTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Produk"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Produk' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => tableRef.current?.openImport()} className="gap-1">
                <Upload size={16} />
                Import Produk
              </Button>
              <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
                <Plus size={16} />
                Tambah Produk
              </Button>
            </div>
          </RoleGuard>
        }
      />
      <ProductTable ref={tableRef} />
    </div>
  )
}
