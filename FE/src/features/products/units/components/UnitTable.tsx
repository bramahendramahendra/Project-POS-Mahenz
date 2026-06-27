import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

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
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

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

  const handleFilterChange = (newFilter: UnitListFilter) => {
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

  const handleOpenEdit = (unit: Unit) => {
    setEditingUnit(unit)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingUnit(null)
  }

  const handleOpenDelete = (unit: Unit) => {
    setDeletingUnit(unit)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deletingUnit) return
    deleteUnit(deletingUnit.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingUnit(null)
      },
    })
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Satuan berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const hasFilter = !!filter.search || filter.is_active !== undefined

  const columns = buildUnitColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })


  return (
    <div className="space-y-4">
      <UnitFilterBar 
        filter={filter} 
        onChange={handleFilterChange} 
        onReset={handleReset} 
      />

      <DataTable<Unit & Record<string, unknown>>
        columns={columns}
        data={units as (Unit & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Satuan tidak ditemukan' : 'Belum ada satuan'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'Tambah satuan pertama Anda untuk memulai.'
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

      <UnitFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) handleCloseForm() }}
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
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
