import { forwardRef, useImperativeHandle, useState } from 'react'
import { Eye, Lock, LockOpen, Pencil, RotateCcw, Search, Trash2 } from 'lucide-react'
import { toast } from 'sonner'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, DataTable, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useDebounce, useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import {
  useSupplierListQuery,
  useCreateSupplierMutation,
  useUpdateSupplierMutation,
  useDeleteSupplierMutation,
  useToggleSupplierStatusMutation,
} from '../suppliers.api'
import type { Supplier } from '../suppliers.types'
import { SupplierFormModal } from './SupplierFormModal'
import { SupplierDetailModal } from './SupplierDetailModal'
import type { SupplierFormValues } from '../suppliers.schema'

export interface SupplierTableHandle {
  openAdd: () => void
}

export const SupplierTable = forwardRef<SupplierTableHandle, object>(function SupplierTable(_, ref) {
  const [search, setSearch] = useState('')
  const [isActiveFilter, setIsActiveFilter] = useState<'all' | 'true' | 'false'>('all')
  const debouncedSearch = useDebounce(search, 300)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: confirmOpen, open: openConfirm, close: closeConfirm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetailModal, close: closeDetailModal } = useDisclosure()

  const [editingSupplier, setEditingSupplier] = useState<Supplier | null>(null)
  const [deletingSupplier, setDeletingSupplier] = useState<Supplier | null>(null)
  const [detailSupplierId, setDetailSupplierId] = useState<number | null>(null)
  const [pendingAction, setPendingAction] = useState<{ values: SupplierFormValues; supplier: Supplier | null} | null>(null)

  const isActiveValue = isActiveFilter === 'all' ? undefined : isActiveFilter === 'true'

  const { data: supplierData, isLoading } = useSupplierListQuery({
    page,
    limit: pageSize,
    search: debouncedSearch,
    is_active: isActiveValue,
  })
  const suppliers = supplierData?.data ?? []
  const total = supplierData?.total ?? 0

  const { mutate: createSupplier, isPending: isCreating } = useCreateSupplierMutation()
  const { mutate: updateSupplier, isPending: isUpdating } = useUpdateSupplierMutation()
  const { mutate: deleteSupplier, isPending: isDeleting } = useDeleteSupplierMutation()
  const { mutate: toggleStatus } = useToggleSupplierStatusMutation()

  const isPending = isCreating || isUpdating

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

  const handleCloseForm = () => {
    closeForm()
    setEditingSupplier(null)
  }

  const onFormSubmit = (values: SupplierFormValues) => {
    setPendingAction({ values, supplier: editingSupplier })
    closeForm()
    openConfirm()
  }

  const handleConfirmCancel = () => {
    closeConfirm()
    if (pendingAction) {
      setEditingSupplier(pendingAction.supplier)
      openForm()
    }
    setPendingAction(null)
  }

  const handleConfirmSave = () => {
    if (!pendingAction) return
    if (pendingAction.supplier !== null) {
      updateSupplier(
        { id: pendingAction.supplier.id, ...pendingAction.values },
        {
          onSuccess: () => {
            closeConfirm()
            setEditingSupplier(null)
            setPendingAction(null)
          },
        }
      )
    } else {
      createSupplier(pendingAction.values, {
        onSuccess: () => {
          closeConfirm()
          setPendingAction(null)
        },
      })
    }
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

  const columns: ColumnDef<Supplier>[] = [
    {
      key: 'supplier_code',
      header: 'Kode',
      width: '110px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.supplier_code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Supplier',
      cell: (row) => (
        <span className="font-medium text-gray-800">
          {row.name}
        </span>
      ),
    },
    {
      key: 'contact_person',
      header: 'Nama Kontak',
      cell: (row) =>
        row.contact_person ? (
          <span>{row.contact_person}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
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
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '90px',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '140px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-gray-500 hover:text-indigo-600"
            onClick={() => handleOpenDetail(row.id)}
            title="Detail"
          >
            <Eye size={14} />
          </Button>
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
            <Button
              variant="ghost"
              size="icon"
              className={`h-7 w-7 ${row.is_active ? 'text-gray-500 hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}`}
              onClick={() => handleToggleStatus(row)}
              title={row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
            >
              {row.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
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
      {/* Filter */}
      <div className="flex items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            placeholder="Cari supplier..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); resetPage() }}
            className="pl-8 h-9 text-sm"
          />
        </div>
        <Select
          value={isActiveFilter}
          onValueChange={(v) => { setIsActiveFilter(v as 'all' | 'true' | 'false'); resetPage() }}
        >
          <SelectTrigger className="w-36 h-9">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Status</SelectItem>
            <SelectItem value="true">Aktif</SelectItem>
            <SelectItem value="false">Nonaktif</SelectItem>
          </SelectContent>
        </Select>
        {(search || isActiveFilter !== 'all') && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => { setSearch(''); setIsActiveFilter('all'); resetPage() }}
            className="h-9 gap-1"
          >
            <RotateCcw size={13} />
            Reset
          </Button>
        )}
      </div>

      {/* Table */}
      <DataTable<Supplier & Record<string, unknown>>
        columns={columns}
        data={suppliers as (Supplier & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        emptyMessage={debouncedSearch ? 'Supplier tidak ditemukan' : 'Belum ada supplier'}
        emptyDescription={
          debouncedSearch
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah supplier pertama Anda untuk memulai.'
        }
      />

      {/* Form Modal */}
      <SupplierFormModal
        open={formOpen}
        onOpenChange={(open) => {
          if (!open && pendingAction !== null) return
          if (!open) handleCloseForm()
        }}
        supplier={editingSupplier}
        onSubmit={onFormSubmit}
        isLoading={isPending}
      />

      {/* Confirm Save */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={(open) => { if (!open) handleConfirmCancel() }}
        title={pendingAction?.supplier !== null ? 'Update Supplier' : 'Tambah Supplier'}
        description={`Yakin ingin ${pendingAction?.supplier !== null ? 'mengupdate' : 'menambahkan'} supplier "${pendingAction?.values.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmSave}
      />

      {/* Confirm Delete */}
      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) { closeDelete(); setDeletingSupplier(null) }
        }}
        title="Hapus Supplier"
        description={`Yakin ingin menghapus supplier "${deletingSupplier?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleDelete}
      />

      {/* Detail Modal */}
      <SupplierDetailModal
        open={detailOpen}
        onOpenChange={(open) => {
          if (!open) { closeDetailModal(); setDetailSupplierId(null) }
        }}
        supplierId={detailSupplierId ?? undefined}
      />
    </div>
  )
})
