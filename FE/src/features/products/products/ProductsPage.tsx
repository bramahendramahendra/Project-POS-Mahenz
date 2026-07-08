import { useRef } from 'react'
import { Plus, Upload } from 'lucide-react'

import { PageHeader, PrerequisiteGuard, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { ProductTable } from './components/ProductTable'
import type { ProductTableHandle } from './components/ProductTable'
import { useProductPrerequisites } from './hooks/useProductPrerequisites'

export function ProductsPage() {
  const tableRef = useRef<ProductTableHandle>(null)
  const { isLoading, items } = useProductPrerequisites()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Produk"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Produk' }]}
        actions={
          <RoleGuard menuKey="produk.produk" action="can_create">
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
      <PrerequisiteGuard
        isLoading={isLoading}
        title="Belum bisa menambah produk"
        description="Sebelum menambah produk, pastikan data berikut sudah tersedia:"
        items={items}
      >
        <ProductTable ref={tableRef} />
      </PrerequisiteGuard>
    </div>
  )
}
