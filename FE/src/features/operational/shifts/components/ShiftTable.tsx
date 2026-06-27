import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import {
  useShiftListQuery,
  useDeleteShiftMutation,
  useToggleShiftStatusMutation,
} from '../shifts.api'
import type { Shift, ShiftListFilter } from '../shifts.types'
import { ShiftFilterBar } from './ShiftFilterBar'
import { ShiftFormModal } from './ShiftFormModal'
import { buildShiftColumns } from './ShiftTableColumns'

export interface ShiftTableHandle {
  openAdd: () => void
}

export const ShiftTable = forwardRef<ShiftTableHandle, object>(function ShiftTable(_, ref) {
  const [filter, setFilter] = useState<ShiftListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingShift, setEditingShift] = useState<Shift | null>(null)
  const [deletingShift, setDeletingShift] = useState<Shift | null>(null)

  const { data: shiftData, isLoading } = useShiftListQuery({ 
    ...filter, 
    page, 
    limit: pageSize,
  })
  const shifts = shiftData?.data ?? []
  const total = shiftData?.total ?? 0

  const { mutate: deleteShift, isPending: isDeleting } = useDeleteShiftMutation()
  const { mutate: toggleStatus } = useToggleShiftStatusMutation()

  const handleOpenAdd = () => {
    setEditingShift(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleOpenEdit = (shift: Shift) => {
    setEditingShift(shift)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingShift(null)
  }

  const handleFilterChange = (newFilter: ShiftListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleOpenDelete = (shift: Shift) => {
    setDeletingShift(shift)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingShift(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingShift) return
    deleteShift(deletingShift.id, {
      onSuccess: () => handleCloseDelete(),
    })
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Shift berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const hasFilter = !!filter.search || filter.is_active !== undefined

  const columns = buildShiftColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  return (
    <div className="space-y-4">
      <ShiftFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
      />

      <DataTable<Shift & Record<string, unknown>>
        columns={columns}
        data={shifts as (Shift & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Shift tidak ditemukan' : 'Belum ada data shift'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'Tambah shift pertama untuk memulai.'
        }
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
      />

      <ShiftFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        shift={editingShift}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Shift"
        description={`Yakin ingin menghapus shift "${deletingShift?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
