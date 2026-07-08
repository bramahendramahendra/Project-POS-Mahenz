import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'

import {
  useDeleteMenuMutation,
  useMenuListQuery,
  useMenuOptionsQuery,
} from '@/features/menu/menu.api'
import type { MenuListFilter, MenuResponse } from '@/features/menu/menu.types'

import { MenuFilterBar } from './MenuFilterBar'
import { MenuFormModal } from './MenuFormModal'
import { buildMenuColumns } from './MenuTableColumns'

export interface MenuTableHandle {
  openAdd: () => void
}

export const MenuTable = forwardRef<MenuTableHandle, object>(function MenuTable(_, ref) {
  const [filter, setFilter] = useState<MenuListFilter>({ page: 1, limit: 10, search: '' })

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingMenu, setEditingMenu] = useState<MenuResponse | null>(null)
  const [deletingMenu, setDeletingMenu] = useState<MenuResponse | null>(null)

  const { data: menuData, isLoading } = useMenuListQuery({ ...filter, page, limit: pageSize })
  const menus = menuData?.data ?? []
  const total = menuData?.total ?? 0

  const { data: rootOptions = [] } = useMenuOptionsQuery()
  const parentMap = rootOptions.reduce<Record<number, string>>((acc, m) => {
    acc[m.id] = m.label
    return acc
  }, {})

  const { mutate: deleteMenu, isPending: isDeleting } = useDeleteMenuMutation()

  const handleOpenAdd = () => {
    setEditingMenu(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleFilterChange = (newFilter: MenuListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleOpenEdit = (menu: MenuResponse) => {
    setEditingMenu(menu)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingMenu(null)
  }

  const handleOpenDelete = (menu: MenuResponse) => {
    setDeletingMenu(menu)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingMenu(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingMenu) return
    deleteMenu(deletingMenu.id, { onSuccess: () => handleCloseDelete() })
  }

  const hasFilter = !!filter.search

  const columns = buildMenuColumns({
    parentMap,
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
  })

  return (
    <div className="space-y-4">
      <MenuFilterBar filter={filter} onChange={handleFilterChange} />

      <DataTable<MenuResponse & Record<string, unknown>>
        columns={columns}
        data={menus as (MenuResponse & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={hasFilter ? 'Menu tidak ditemukan' : 'Belum ada menu'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci pencarian Anda.'
            : 'Tambah menu pertama Anda untuk memulai.'
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

      <MenuFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        menuId={editingMenu?.id}
        parentOptions={rootOptions}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Menu"
        description={`Yakin ingin menghapus menu "${deletingMenu?.label}"? Sub-menu di bawahnya akan kehilangan parent.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
