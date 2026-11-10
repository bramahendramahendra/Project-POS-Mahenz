import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import { monthStart, todayStr } from '@/shared/utils'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useExpenseListQuery, useDeleteExpenseMutation } from '../expenses.api'
import type { Expense, ExpenseListFilter } from '../expenses.types'
import { ExpenseFilterBar } from './ExpenseFilterBar'
import { ExpenseFormModal } from './ExpenseFormModal'
import { buildExpenseColumns } from './ExpenseTableColumns'

export interface ExpenseTableHandle {
  openAdd: () => void
}

const DEFAULT_FILTER: ExpenseListFilter = {
  page: 1,
  limit: 10,
  start_date: monthStart(),
  end_date: todayStr(),
}

export const ExpenseTable = forwardRef<ExpenseTableHandle, object>(function ExpenseTable(_, ref) {
  const [filter, setFilter] = useState<ExpenseListFilter>(DEFAULT_FILTER)
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingExpense, setEditingExpense] = useState<Expense | null>(null)
  const [deletingExpense, setDeletingExpense] = useState<Expense | null>(null)

  const { data, isLoading } = useExpenseListQuery({ ...filter, page, limit: pageSize })
  const expenses = data?.data ?? []
  const total = data?.total ?? 0

  const { mutate: deleteExpense, isPending: isDeleting } = useDeleteExpenseMutation()

  const handleOpenAdd = () => {
    setEditingExpense(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleFilterChange = (newFilter: ExpenseListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter(DEFAULT_FILTER)
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleOpenEdit = (expense: Expense) => {
    setEditingExpense(expense)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingExpense(null)
  }

  const handleOpenDelete = (expense: Expense) => {
    setDeletingExpense(expense)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingExpense(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingExpense) return
    deleteExpense(deletingExpense.id, { onSuccess: () => handleCloseDelete() })
  }

  const columns = buildExpenseColumns({ onEdit: handleOpenEdit, onDelete: handleOpenDelete })

  return (
    <div className="space-y-4">
      <ExpenseFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Expense & Record<string, unknown>>
        columns={columns}
        data={expenses as (Expense & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data pengeluaran"
        emptyDescription="Data pengeluaran akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <ExpenseFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        expense={editingExpense}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Pengeluaran"
        description={`Yakin ingin menghapus pengeluaran "${deletingExpense?.description}"? Tindakan ini tidak dapat dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
