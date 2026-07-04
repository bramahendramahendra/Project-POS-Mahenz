import { DataTable } from '@/shared/components'
import type { PaginationProps, SortState } from '@/shared/components/DataTable/DataTable.types'

import type { Expense } from '../expenses.types'
import { buildExpenseColumns } from './ExpenseTableColumns'

interface ExpenseTableProps {
  data: Expense[]
  isLoading: boolean
  pagination: PaginationProps
  currentSort?: SortState
  onSort?: (sort: SortState) => void
  onEdit: (expense: Expense) => void
  onDelete: (expense: Expense) => void
}

export function ExpenseTable({
  data,
  isLoading,
  pagination,
  currentSort,
  onSort,
  onEdit,
  onDelete,
}: ExpenseTableProps) {
  const columns = buildExpenseColumns({ onEdit, onDelete })

  return (
    <DataTable<Expense & Record<string, unknown>>
      columns={columns}
      data={data as (Expense & Record<string, unknown>)[]}
      isLoading={isLoading}
      currentSort={currentSort}
      onSort={onSort}
      emptyMessage="Belum ada data pengeluaran"
      emptyDescription="Data pengeluaran akan muncul sesuai filter periode yang dipilih."
      pagination={pagination}
    />
  )
}
