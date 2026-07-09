import { useRef } from 'react'
import { Plus } from 'lucide-react'

import { PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { ExpenseTable } from './components/ExpenseTable'
import type { ExpenseTableHandle } from './components/ExpenseTable'

export function ExpensesPage() {
  const tableRef = useRef<ExpenseTableHandle>(null)

  return (
    <div className="space-y-4">
      <PageHeader
        title="Pengeluaran"
        breadcrumbs={[{ label: 'Finance' }, { label: 'Pengeluaran' }]}
        actions={
          <RoleGuard menuKey="keuangan.pengeluaran" action="can_create">
            <Button onClick={() => tableRef.current?.openAdd()} className="gap-1">
              <Plus size={16} />
              Tambah Pengeluaran
            </Button>
          </RoleGuard>
        }
      />
      <ExpenseTable ref={tableRef} />
    </div>
  )
}
