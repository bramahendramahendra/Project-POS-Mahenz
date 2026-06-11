import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useCustomerListQuery, useDeleteCustomerMutation } from '../customers.api'
import type { Customer, CustomerListFilter } from '../customers.types'
import { CustomerFilterBar } from './CustomerFilterBar'
import { CustomerFormModal } from './CustomerFormModal'
import { buildCustomerColumns } from './CustomerTableColumns'

export interface CustomerTableHandle {
  openAdd: () => void
}

export const CustomerTable = forwardRef<CustomerTableHandle, object>(function CustomerTable(_, ref) {
  const [filter, setFilter] = useState<CustomerListFilter>({ page: 1, limit: 10, search: '' })

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

  const handleOpenAdd = () => {
    setEditingCustomer(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleOpenEdit = (customer: Customer) => {
    setEditingCustomer(customer)
    openForm()
  }

  const handleOpenDelete = (customer: Customer) => {
    setDeletingCustomer(customer)
    openDelete()
  }

  const handleFilterChange = (newFilter: CustomerListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    resetPage()
  }

  const handleDelete = () => {
    if (!deletingCustomer) return
    deleteCustomer(deletingCustomer.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingCustomer(null)
      },
    })
  }

  const columns = buildCustomerColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
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
        emptyMessage={hasFilter ? 'Pelanggan tidak ditemukan' : 'Belum ada pelanggan'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah pelanggan pertama Anda untuk memulai.'
        }
      />

      <CustomerFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) { closeForm(); setEditingCustomer(null) } }}
        customer={editingCustomer}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingCustomer(null) } }}
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
