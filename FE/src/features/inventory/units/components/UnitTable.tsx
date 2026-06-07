import { forwardRef, useImperativeHandle, useState } from 'react'
import { Lock, LockOpen, Pencil, RotateCcw, Search, Trash2 } from 'lucide-react'
import { toast } from 'sonner'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, DataTable, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { useDebounce, useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import {
  useUnitListQuery,
  useCreateUnitMutation,
  useUpdateUnitMutation,
  useDeleteUnitMutation,
  useToggleUnitStatusMutation,
} from '../units.api'
import type { Unit } from '../units.types'
import { UnitFormModal } from './UnitFormModal'
import type { UnitFormValues } from './UnitFormModal'

export interface UnitTableHandle {
  openAdd: () => void
}

export const UnitTable = forwardRef<UnitTableHandle, object>(function UnitTable(_, ref) {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: confirmOpen, open: openConfirm, close: closeConfirm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingUnit, setEditingUnit] = useState<Unit | null>(null)
  const [deletingUnit, setDeletingUnit] = useState<Unit | null>(null)
  const [pendingAction, setPendingAction] = useState<{ values: UnitFormValues; unit: Unit | null } | null>(null)

  const { data: unitData, isLoading } = useUnitListQuery({
    page,
    limit: pageSize,
    search: debouncedSearch,
  })
  const units = unitData?.data ?? []
  const total = unitData?.total ?? 0

  const { mutate: createUnit, isPending: isCreating } = useCreateUnitMutation()
  const { mutate: updateUnit, isPending: isUpdating } = useUpdateUnitMutation()
  const { mutate: deleteUnit, isPending: isDeleting } = useDeleteUnitMutation()
  const { mutate: toggleStatus } = useToggleUnitStatusMutation()

  const isPending = isCreating || isUpdating

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

  const handleCloseForm = () => {
    closeForm()
    setEditingUnit(null)
  }

  const onFormSubmit = (values: UnitFormValues) => {
    setPendingAction({ values, unit: editingUnit })
    closeForm()
    openConfirm()
  }

  const handleConfirmCancel = () => {
    closeConfirm()
    if (pendingAction) {
      setEditingUnit(pendingAction.unit)
      openForm()
    }
    setPendingAction(null)
  }

  const handleConfirmSave = () => {
    if (!pendingAction) return
    if (pendingAction.unit !== null) {
      updateUnit(
        { id: pendingAction.unit.id, ...pendingAction.values },
        {
          onSuccess: () => {
            closeConfirm()
            setEditingUnit(null)
            setPendingAction(null)
          },
        }
      )
    } else {
      createUnit(pendingAction.values, {
        onSuccess: () => {
          closeConfirm()
          setPendingAction(null)
        },
      })
    }
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

  const columns: ColumnDef<Unit>[] = [
    {
      key: 'name',
      header: 'Nama Satuan',
      cell: (row) => (
        <span className="font-medium text-gray-800">
          {row.name}
        </span>
      ),
    },
    {
      key: 'abbreviation',
      header: 'Singkatan',
      width: '100px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.abbreviation}
        </span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '120px',
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
      {/* Search */}
      <div className="flex items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            placeholder="Cari satuan..."
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
      <DataTable<Unit & Record<string, unknown>>
        columns={columns}
        data={units as (Unit & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        emptyMessage={debouncedSearch ? 'Satuan tidak ditemukan' : 'Belum ada satuan'}
        emptyDescription={
          debouncedSearch
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah satuan pertama Anda untuk memulai.'
        }
      />

      {/* Form Modal */}
      <UnitFormModal
        open={formOpen}
        onOpenChange={(open) => {
          if (!open && pendingAction !== null) return
          if (!open) handleCloseForm()
        }}
        unit={editingUnit}
        onSubmit={onFormSubmit}
        isLoading={isPending}
      />

      {/* Confirm Save */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={(open) => { if (!open) handleConfirmCancel() }}
        title={pendingAction?.unit !== null ? 'Update Satuan' : 'Tambah Satuan'}
        description={`Yakin ingin ${pendingAction?.unit !== null ? 'mengupdate' : 'menambahkan'} satuan "${pendingAction?.values.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmSave}
      />

      {/* Confirm Delete */}
      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) { closeDelete(); setDeletingUnit(null) }
        }}
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
