import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

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

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetailModal, close: closeDetailModal } = useDisclosure()

  const [editingSupplier, setEditingSupplier] = useState<Supplier | null>(null)
  const [deletingSupplier, setDeletingSupplier] = useState<Supplier | null>(null)
  const [detailSupplierId, setDetailSupplierId] = useState<number | null>(null)

  const { data: supplierData, isLoading } = useSupplierListQuery({ ...filter, page, limit: pageSize })
  const suppliers = supplierData?.data ?? []
  const total = supplierData?.total ?? 0

  const { mutate: deleteSupplier, isPending: isDeleting } = useDeleteSupplierMutation()
  const { mutate: toggleStatus } = useToggleSupplierStatusMutation()

  const handleOpenAdd = () => {
    setEditingSupplier(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleOpenEdit = (supplier: Supplier) => {
    setEditingSupplier(supplier)
    openForm()
  }

  const handleOpenDetail = (id: number) => {
    setDetailSupplierId(id)
    openDetailModal()
  }

  const handleOpenDelete = (supplier: Supplier) => {
    setDeletingSupplier(supplier)
    openDelete()
  }

  const handleFilterChange = (newFilter: SupplierListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    resetPage()
  }

  const handleDelete = () => {
    if (!deletingSupplier) return
    deleteSupplier(deletingSupplier.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingSupplier(null)
      },
    })
  }

  const handleToggleStatus = (row: Supplier) => {
    toggleStatus(row.id, {
      onSuccess: () =>
        toast.success(`Supplier berhasil ${row.is_active ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const columns = buildSupplierColumns({
    onDetail: handleOpenDetail,
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  const hasFilter = !!filter.search || filter.is_active !== undefined

  return (
    <div className="space-y-4">
      <SupplierFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Supplier & Record<string, unknown>>
        columns={columns}
        data={suppliers as (Supplier & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        emptyMessage={hasFilter ? 'Supplier tidak ditemukan' : 'Belum ada supplier'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah supplier pertama Anda untuk memulai.'
        }
      />

      <SupplierFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) { closeForm(); setEditingSupplier(null) } }}
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
        onConfirm={handleDelete}
      />

      <SupplierDetailModal
        open={detailOpen}
        onOpenChange={(open) => { if (!open) { closeDetailModal(); setDetailSupplierId(null) } }}
        supplierId={detailSupplierId ?? undefined}
      />
    </div>
  )
})
