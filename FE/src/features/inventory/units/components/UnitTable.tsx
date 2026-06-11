import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useUnitListQuery, useDeleteUnitMutation, useToggleUnitStatusMutation } from '../units.api'
import type { Unit, UnitListFilter } from '../units.types'
import { UnitFilterBar } from './UnitFilterBar'
import { UnitFormModal } from './UnitFormModal'
import { buildUnitColumns } from './UnitTableColumns'

export interface UnitTableHandle {
  openAdd: () => void
}

export const UnitTable = forwardRef<UnitTableHandle, object>(function UnitTable(_, ref) {
  const [filter, setFilter] = useState<UnitListFilter>({ page: 1, limit: 10, search: '' })

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingUnit, setEditingUnit] = useState<Unit | null>(null)
  const [deletingUnit, setDeletingUnit] = useState<Unit | null>(null)

  const { data: unitData, isLoading } = useUnitListQuery({ ...filter, page, limit: pageSize })
  const units = unitData?.data ?? []
  const total = unitData?.total ?? 0

  const { mutate: deleteUnit, isPending: isDeleting } = useDeleteUnitMutation()
  const { mutate: toggleStatus } = useToggleUnitStatusMutation()

  const handleOpenAdd = () => {
    setEditingUnit(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleOpenEdit = (unit: Unit) => {
    setEditingUnit(unit)
    openForm()
  }

  const handleOpenDelete = (unit: Unit) => {
    setDeletingUnit(unit)
    openDelete()
  }

  const handleFilterChange = (newFilter: UnitListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    resetPage()
  }

  const handleDelete = () => {
    if (!deletingUnit) return
    deleteUnit(deletingUnit.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingUnit(null)
      },
    })
  }

  const handleToggleStatus = (row: Unit) => {
    toggleStatus(row.id, {
      onSuccess: () =>
        toast.success(`Satuan berhasil ${row.is_active ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const columns = buildUnitColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  const hasFilter = !!filter.search

  return (
    <div className="space-y-4">
      <UnitFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Unit & Record<string, unknown>>
        columns={columns}
        data={units as (Unit & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        emptyMessage={hasFilter ? 'Satuan tidak ditemukan' : 'Belum ada satuan'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah satuan pertama Anda untuk memulai.'
        }
      />

      <UnitFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) { closeForm(); setEditingUnit(null) } }}
        unit={editingUnit}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingUnit(null) } }}
        title="Hapus Satuan"
        description={`Yakin ingin menghapus satuan "${deletingUnit?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleDelete}
      />
    </div>
  )
})
