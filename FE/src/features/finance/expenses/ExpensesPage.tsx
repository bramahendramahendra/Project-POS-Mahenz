import { useState } from 'react'
import { Plus } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useExpensesQuery, useDeleteExpenseMutation } from './expenses.api'
import type { Expense, ExpenseListFilter } from './expenses.types'
import { ExpenseFilterBar } from './components/ExpenseFilterBar'
import { ExpenseTable } from './components/ExpenseTable'
import { ExpenseFormModal } from './components/ExpenseFormModal'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

export function ExpensesPage() {
  const [filter, setFilter] = useState<ExpenseListFilter>({
    page: 1,
    limit: 10,
    start_date: monthStartString(),
    end_date: todayString(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()
  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)
  const [deletingExpense, setDeletingExpense] = useState<Expense | null>(null)

  const { data, isLoading } = useExpensesQuery({ ...filter, page, limit: pageSize })
  const { mutate: deleteExpense, isPending: isDeleting } = useDeleteExpenseMutation()

  const items: Expense[] = data?.data ?? []
  const total = data?.total ?? 0

  function handleEdit(expense: Expense) {
    setEditingExpense(expense)
    openForm()
  }

  function handleDelete(expense: Expense) {
    setDeletingExpense(expense)
    openDelete()
  }

  function handleFilterChange(newFilter: ExpenseListFilter) {
    setFilter(newFilter)
    reset()
  }

  function confirmDelete() {
    if (!deletingExpense) return
    deleteExpense(deletingExpense.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingExpense(null)
      },
    })
  }

  return (
    <div className="space-y-4">
      <PageHeader
        title="Pengeluaran"
        breadcrumbs={[{ label: 'Finance' }, { label: 'Pengeluaran' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={openForm} className="gap-1">
              <Plus size={16} />
              Tambah Pengeluaran
            </Button>
          </RoleGuard>
        }
      />

      <ExpenseFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={() => {
          setFilter({ page: 1, limit: 10, start_date: monthStartString(), end_date: todayString() })
          reset()
        }}
      />

      <ExpenseTable
        data={items}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />

      <ExpenseFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) { closeForm(); setEditingExpense(null) } }}
        expense={editingExpense}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) {
            closeDelete()
            setDeletingExpense(null)
          }
        }}
        title="Hapus Pengeluaran"
        description={`Yakin ingin menghapus pengeluaran "${deletingExpense?.description}"? Tindakan ini tidak dapat dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={confirmDelete}
      />
    </div>
  )
}
