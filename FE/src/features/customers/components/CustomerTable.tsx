import { forwardRef, useImperativeHandle, useState } from 'react'
import { Pencil, RotateCcw, Search, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, DataTable, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { useDebounce, useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import {
  useCustomerListQuery,
  useCreateCustomerMutation,
  useUpdateCustomerMutation,
  useDeleteCustomerMutation,
} from '../customers.api'
import type { Customer } from '../customers.types'
import { CustomerFormModal } from './CustomerFormModal'
import type { CustomerFormValues } from '../customers.schema'

export interface CustomerTableHandle {
  openAdd: () => void
}

export const CustomerTable = forwardRef<CustomerTableHandle, object>(function CustomerTable(_, ref) {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: confirmOpen, open: openConfirm, close: closeConfirm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null)
  const [deletingCustomer, setDeletingCustomer] = useState<Customer | null>(null)
  const [pendingAction, setPendingAction] = useState<{ values: CustomerFormValues; customer: Customer | null } | null>(null)

  const { data: customerData, isLoading } = useCustomerListQuery({
    page,
    limit: pageSize,
    search: debouncedSearch,
  })
  const customers = customerData?.data ?? []
  const total = customerData?.total ?? 0

  const { mutate: createCustomer, isPending: isCreating } = useCreateCustomerMutation()
  const { mutate: updateCustomer, isPending: isUpdating } = useUpdateCustomerMutation()
  const { mutate: deleteCustomer, isPending: isDeleting } = useDeleteCustomerMutation()

  const isPending = isCreating || isUpdating

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

  const handleCloseForm = () => {
    closeForm()
    setEditingCustomer(null)
  }

  const onFormSubmit = (values: CustomerFormValues) => {
    setPendingAction({ values, customer: editingCustomer })
    closeForm()
    openConfirm()
  }

  const handleConfirmCancel = () => {
    closeConfirm()
    if (pendingAction) {
      setEditingCustomer(pendingAction.customer)
      openForm()
    }
    setPendingAction(null)
  }

  const handleConfirmSave = () => {
    if (!pendingAction) return
    if (pendingAction.customer !== null) {
      updateCustomer(
        { id: pendingAction.customer.id, ...pendingAction.values },
        {
          onSuccess: () => {
            closeConfirm()
            setEditingCustomer(null)
            setPendingAction(null)
          },
        }
      )
    } else {
      createCustomer(pendingAction.values, {
        onSuccess: () => {
          closeConfirm()
          setPendingAction(null)
        },
      })
    }
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

  const columns: ColumnDef<Customer>[] = [
    {
      key: 'customer_code',
      header: 'Kode',
      width: '90px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.customer_code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Pelanggan',
      cell: (row) => <span className="font-medium text-gray-800">{row.name}</span>,
    },
    {
      key: 'phone',
      header: 'Telepon',
      cell: (row) =>
        row.phone ? (
          <span className="font-mono text-sm">{row.phone}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'address',
      header: 'Alamat',
      cell: (row) =>
        row.address ? (
          <span className="text-sm text-gray-600">{row.address}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-gray-500 hover:text-blue-600"
              onClick={() => handleOpenEdit(row)}
              title="Edit"
            >
              <Pencil size={14} />
            </Button>
          </RoleGuard>
          <RoleGuard allowedRoles={[ROLES.OWNER]}>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-gray-500 hover:text-red-600"
              onClick={() => handleOpenDelete(row)}
              title="Hapus"
            >
              <Trash2 size={14} />
            </Button>
          </RoleGuard>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      {/* Search */}
      <div className="flex items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            placeholder="Cari pelanggan..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); resetPage() }}
            className="pl-8 h-9 text-sm"
          />
        </div>
        {search && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => { setSearch(''); resetPage() }}
            className="h-9 gap-1"
          >
            <RotateCcw size={13} />
            Reset
          </Button>
        )}
      </div>

      {/* Table */}
      <DataTable<Customer & Record<string, unknown>>
        columns={columns}
        data={customers as (Customer & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        emptyMessage={debouncedSearch ? 'Pelanggan tidak ditemukan' : 'Belum ada pelanggan'}
        emptyDescription={
          debouncedSearch
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah pelanggan pertama Anda untuk memulai.'
        }
      />

      {/* Form Modal */}
      <CustomerFormModal
        open={formOpen}
        onOpenChange={(open) => {
          if (!open && pendingAction !== null) return
          if (!open) handleCloseForm()
        }}
        customer={editingCustomer}
        onSubmit={onFormSubmit}
        isLoading={isPending}
      />

      {/* Confirm Save */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={(open) => { if (!open) handleConfirmCancel() }}
        title={pendingAction?.customer !== null ? 'Update Pelanggan' : 'Tambah Pelanggan'}
        description={`Yakin ingin ${pendingAction?.customer !== null ? 'mengupdate' : 'menambahkan'} pelanggan "${pendingAction?.values.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmSave}
      />

      {/* Confirm Delete */}
      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) { closeDelete(); setDeletingCustomer(null) }
        }}
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
