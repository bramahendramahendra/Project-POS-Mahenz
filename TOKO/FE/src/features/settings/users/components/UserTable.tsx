import { useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { useAuth } from '@/features/auth'

import {
  useDeleteUserMutation,
  useToggleUserStatusMutation,
  useUserListQuery,
} from '../users.api'
import type { AppUser, UserListFilter } from '../users.types'
import { UserFilterBar } from './UserFilterBar'
import { UserFormModal } from './UserFormModal'
import { buildUserColumns } from './UserTableColumns'

export function UserTable() {
  const { user: currentUser } = useAuth()

  const [filter, setFilter] = useState<UserListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: toggleOpen, open: openToggle, close: closeToggle } = useDisclosure()

  const [editingUser, setEditingUser] = useState<AppUser | null>(null)
  const [deletingUser, setDeletingUser] = useState<AppUser | null>(null)
  const [togglingUser, setTogglingUser] = useState<AppUser | null>(null)

  const { data: userData, isLoading } = useUserListQuery({ ...filter, page, limit: pageSize })
  const users = userData?.data ?? []
  const total = userData?.total ?? 0

  const { mutate: deleteUser, isPending: isDeleting } = useDeleteUserMutation()
  const { mutate: toggleStatus, isPending: isToggling } = useToggleUserStatusMutation()

  const handleFilterChange = (newFilter: UserListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    setSortState(undefined)
    reset()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    reset()
  }

  const handleOpenEdit = (user: AppUser) => {
    setEditingUser(user)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingUser(null)
  }

  const handleOpenDelete = (user: AppUser) => {
    setDeletingUser(user)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingUser(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingUser) return
    deleteUser(deletingUser.id, { onSuccess: () => handleCloseDelete() })
  }

  const handleOpenToggle = (user: AppUser) => {
    setTogglingUser(user)
    openToggle()
  }

  const handleCloseToggle = () => {
    closeToggle()
    setTogglingUser(null)
  }

  const handleConfirmToggle = () => {
    if (!togglingUser) return
    const willActivate = !togglingUser.is_active
    toggleStatus(
      { id: togglingUser.id, willActivate },
      { onSuccess: () => handleCloseToggle() }
    )
  }

  const hasFilter = filter.search || filter.is_active !== undefined

  const columns = buildUserColumns({
    currentUserId: currentUser?.id,
    onEdit: handleOpenEdit,
    onToggleStatus: handleOpenToggle,
    onDelete: handleOpenDelete,
  })

  return (
    <div className="space-y-4">
      <UserFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<AppUser & Record<string, unknown>>
        columns={columns}
        data={users as (AppUser & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'User tidak ditemukan' : 'Belum ada user'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'User muncul di sini setelah pertama kali masuk ke Toko.'
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

      <UserFormModal
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) handleCloseForm()
        }}
        user={editingUser}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) handleCloseDelete()
        }}
        title="Hapus User"
        description={`Yakin ingin menghapus user "${deletingUser?.full_name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />

      <ConfirmDialog
        open={toggleOpen}
        onOpenChange={(open) => {
          if (!open) handleCloseToggle()
        }}
        title={togglingUser?.is_active ? 'Nonaktifkan User' : 'Aktifkan User'}
        description={`User "${togglingUser?.full_name}" akan di${togglingUser?.is_active ? 'nonaktifkan' : 'aktifkan'}. Lanjutkan?`}
        confirmLabel={togglingUser?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={togglingUser?.is_active ? 'destructive' : 'default'}
        isLoading={isToggling}
        onConfirm={handleConfirmToggle}
      />
    </div>
  )
}
