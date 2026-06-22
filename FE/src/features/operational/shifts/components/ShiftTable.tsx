import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

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
  const [filter, setFilter] = useState<ShiftListFilter>({ page: 1, limit: 10 })

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingShift, setEditingShift] = useState<Shift | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Shift | null>(null)

  const { data: shiftData, isLoading } = useShiftListQuery({ ...filter, page, limit: pageSize })
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
    reset()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10 })
    reset()
  }

  const handleOpenDelete = (shift: Shift) => {
    setDeleteTarget(shift)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deleteTarget) return
    deleteShift(deleteTarget.id, {
      onSuccess: () => {
        toast.success('Shift berhasil dihapus')
        closeDelete()
        setDeleteTarget(null)
      },
    })
  }

  const handleToggleStatus = (shift: Shift) => {
    toggleStatus(shift.id, {
      onSuccess: () =>
        toast.success(`Shift berhasil ${shift.is_active ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const hasFilter = !!filter.search

  const columns = buildShiftColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  return (
    <div className="space-y-4">
      <ShiftFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Shift & Record<string, unknown>>
        columns={columns}
        data={shifts as (Shift & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={hasFilter ? 'Shift tidak ditemukan' : 'Belum ada data shift'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah shift pertama untuk memulai.'
        }
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <ShiftFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        shift={editingShift}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeleteTarget(null) } }}
        title="Hapus Shift"
        description={`Yakin ingin menghapus shift "${deleteTarget?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
