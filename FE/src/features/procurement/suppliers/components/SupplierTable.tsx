import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import { useSupplierListQuery, useDeleteSupplierMutation, useToggleSupplierStatusMutation } from '../suppliers.api'
import type { Supplier, SupplierListFilter } from '../suppliers.types'
import { SupplierFilterBar } from './SupplierFilterBar'
import { SupplierFormModal } from './SupplierFormModal'
import { SupplierDetailModal } from './SupplierDetailModal'
import { buildSupplierColumns } from './SupplierTableColumns'

export interface SupplierTableHandle {
  openAdd: () => void
}

export const SupplierTable = forwardRef<SupplierTableHandle, object>(function SupplierTable(_, ref) {
  const [filter, setFilter] = useState<SupplierListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()

  const [editingSupplier, setEditingSupplier] = useState<Supplier | null>(null)
  const [deletingSupplier, setDeletingSupplier] = useState<Supplier | null>(null)
  const [detailSupplier, setDetailSupplier] = useState<Supplier | null>(null)

  const { data: supplierData, isLoading } = useSupplierListQuery({ 
    ...filter, 
    page, 
    limit: pageSize 
  })
  const suppliers = supplierData?.data ?? []
  const total = supplierData?.total ?? 0

  const { mutate: deleteSupplier, isPending: isDeleting } = useDeleteSupplierMutation()
  const { mutate: toggleStatus } = useToggleSupplierStatusMutation()

  const handleOpenAdd = () => {
    setEditingSupplier(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleFilterChange = (newFilter: SupplierListFilter) => {
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

  const handleOpenEdit = (supplier: Supplier) => {
    setEditingSupplier(supplier)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingSupplier(null)
  }

  const handleOpenDetail = (supplier: Supplier) => {
    setDetailSupplier(supplier)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailSupplier(null)
  }

  const handleOpenDelete = (supplier: Supplier) => {
    setDeletingSupplier(supplier)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deletingSupplier) return
    deleteSupplier(deletingSupplier.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingSupplier(null)
      },
    })
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Supplier berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const hasFilter = !!filter.search || filter.is_active !== undefined

  const columns = buildSupplierColumns({
    onDetail: handleOpenDetail,
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  return (
    <div className="space-y-4">
      <SupplierFilterBar 
        filter={filter} 
        onChange={handleFilterChange} 
        onReset={handleReset} 
      />

      <DataTable<Supplier & Record<string, unknown>>
        columns={columns}
        data={suppliers as (Supplier & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Supplier tidak ditemukan' : 'Belum ada supplier'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'Tambah supplier pertama Anda untuk memulai.'
        }
        pagination={{ 
          page, 
          pageSize, 
          total, 
          onPageChange, 
          onPageSizeChange, 
          pageSizeOptions 
        }}
      />

      <SupplierFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        supplier={editingSupplier}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingSupplier(null) } }}
        title="Hapus Supplier"
        description={`Yakin ingin menghapus supplier "${deletingSupplier?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />

      <SupplierDetailModal
        open={detailOpen}
        onOpenChange={(open) => { if (!open) handleCloseDetail() }}
        supplierId={detailSupplier?.id}
      />
    </div>
  )
})
