import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import {
  useCustomerListQuery,
  useDeleteCustomerMutation,
  useToggleCustomerStatusMutation,
} from '../customers.api'
import type { Customer, CustomerListFilter } from '../customers.types'
import { CustomerFilterBar } from './CustomerFilterBar'
import { CustomerFormModal } from './CustomerFormModal'
import { buildCustomerColumns } from './CustomerTableColumns'

export interface CustomerTableHandle {
  openAdd: () => void
}

export const CustomerTable = forwardRef<CustomerTableHandle, object>(function CustomerTable(_, ref) {
  const [filter, setFilter] = useState<CustomerListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null)
  const [deletingCustomer, setDeletingCustomer] = useState<Customer | null>(null)

  const { data: customerData, isLoading } = useCustomerListQuery({ ...filter, page, limit: pageSize })
  const customers = customerData?.data ?? []
  const total = customerData?.total ?? 0

  const { mutate: deleteCustomer, isPending: isDeleting } = useDeleteCustomerMutation()
  const { mutate: toggleStatus } = useToggleCustomerStatusMutation()

  const handleOpenAdd = () => {
    setEditingCustomer(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleCloseForm = () => {
    closeForm()
    setEditingCustomer(null)
  }

  const handleOpenEdit = (customer: Customer) => {
    setEditingCustomer(customer)
    openForm()
  }

  const handleOpenDelete = (customer: Customer) => {
    setDeletingCustomer(customer)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingCustomer(null)
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Pelanggan berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const handleFilterChange = (newFilter: CustomerListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    setSortState(undefined)
    resetPage()
  }

  const handleDelete = () => {
    if (!deletingCustomer) return
    deleteCustomer(deletingCustomer.id, {
      onSuccess: () => handleCloseDelete(),
    })
  }

  const columns = buildCustomerColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })

  const hasFilter = !!filter.search || filter.is_active !== undefined

  return (
    <div className="space-y-4">
      <CustomerFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Customer & Record<string, unknown>>
        columns={columns}
        data={customers as (Customer & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Pelanggan tidak ditemukan' : 'Belum ada pelanggan'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah pelanggan pertama Anda untuk memulai.'
        }
      />

      <CustomerFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) handleCloseForm() }}
        customer={editingCustomer}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Pelanggan"
        description={`Yakin ingin menghapus pelanggan "${deletingCustomer?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleDelete}
      />
    </div>
  )
})
