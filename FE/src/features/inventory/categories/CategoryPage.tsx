import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { CategoryTable } from './components/CategoryTable'
import type { CategoryTableHandle } from './components/CategoryTable'

export function CategoryPage() {
  const tableRef = useRef<CategoryTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kategori Produk"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Kategori' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Kategori
            </Button>
          </RoleGuard>
        }
      />
      <CategoryTable ref={tableRef} />
    </div>
  )
}
