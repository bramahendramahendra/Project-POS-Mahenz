import { forwardRef, useImperativeHandle, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import { ROUTES } from '@/shared/constants/routes'

import { useDeleteRoleMutation, useRoleListQuery, useToggleRoleStatusMutation } from '../roles.api'
import type { Role, RoleListFilter } from '../roles.types'
import { RoleFilterBar } from './RoleFilterBar'
import { RoleFormModal } from './RoleFormModal'
import { buildRoleColumns } from './RoleTableColumns'

export interface RoleTableHandle {
  openAdd: () => void
}

export const RoleTable = forwardRef<RoleTableHandle, object>(function RoleTable(_, ref) {
  const navigate = useNavigate()
  const [filter, setFilter] = useState<RoleListFilter>({ page: 1, limit: 10, search: '' })

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [deletingRole, setDeletingRole] = useState<Role | null>(null)

  const { data: roleData, isLoading } = useRoleListQuery({ ...filter, page, limit: pageSize })
  const roles = roleData?.data ?? []
  const total = roleData?.total ?? 0

  const { mutate: deleteRole, isPending: isDeleting } = useDeleteRoleMutation()
  const { mutate: toggleStatus } = useToggleRoleStatusMutation()

  const handleOpenAdd = () => {
    setEditingRole(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleFilterChange = (newFilter: RoleListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleManageAccess = (role: Role) => {
    navigate(ROUTES.SETTINGS_ROLES_ACCESS.replace(':id', String(role.id)))
  }

  const handleOpenEdit = (role: Role) => {
    setEditingRole(role)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingRole(null)
  }

  const handleOpenDelete = (role: Role) => {
    setDeletingRole(role)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingRole(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingRole) return
    deleteRole(deletingRole.id, { onSuccess: () => handleCloseDelete() })
  }

  const hasFilter = !!filter.search

  const columns = buildRoleColumns({
    onManageAccess: handleManageAccess,
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: (id) => toggleStatus(id),
  })

  return (
    <div className="space-y-4">
      <RoleFilterBar filter={filter} onChange={handleFilterChange} />

      <DataTable<Role & Record<string, unknown>>
        columns={columns}
        data={roles as (Role & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={hasFilter ? 'Role tidak ditemukan' : 'Belum ada role'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah role pertama Anda untuk memulai.'
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

      <RoleFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        roleId={editingRole?.id}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Role"
        description={`Yakin ingin menghapus role "${deletingRole?.display_name}"? Semua user dengan role ini harus dipindahkan terlebih dahulu.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
