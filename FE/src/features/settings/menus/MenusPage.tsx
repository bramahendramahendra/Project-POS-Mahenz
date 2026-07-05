import { useState } from 'react'
import { Pencil, Plus, Trash2 } from 'lucide-react'

import { ConfirmDialog, PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import { Input } from '@/shared/components/ui/input'
import { useDebounce, useDisclosure } from '@/shared/hooks'

import { useDeleteMenuMutation, useMenuListQuery } from '@/features/menu/menu.api'
import type { MenuResponse } from '@/features/menu/menu.types'

import { MenuFormModal } from './components/MenuFormModal'

export function MenusPage() {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingMenu, setEditingMenu] = useState<MenuResponse | null>(null)
  const [deletingMenu, setDeletingMenu] = useState<MenuResponse | null>(null)

  const { data: menus = [], isLoading } = useMenuListQuery(
    debouncedSearch ? { search: debouncedSearch } : undefined
  )
  const { mutate: deleteMenu, isPending: isDeleting } = useDeleteMenuMutation()

  const handleOpenAdd = () => {
    setEditingMenu(null)
    openForm()
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

  const parentMap = menus.reduce<Record<number, string>>((acc, m) => {
    acc[m.id] = m.label
    return acc
  }, {})

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Menu"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen Menu' }]}
        actions={
          <RoleGuard menuKey="sistem.menus" action="can_create">
            <Button onClick={handleOpenAdd}>
              <Plus size={14} className="mr-2" />
              Tambah Menu
            </Button>
          </RoleGuard>
        }
      />

      <div className="flex items-center gap-3">
        <Input
          placeholder="Cari menu..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-xs"
        />
      </div>

      <div className="rounded-lg border bg-white overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Key</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Label</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Parent</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Path</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Urutan</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Status</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {isLoading && (
              <tr><td colSpan={7} className="text-center py-8 text-gray-400">Memuat...</td></tr>
            )}
            {!isLoading && menus.length === 0 && (
              <tr><td colSpan={7} className="text-center py-8 text-gray-400">Belum ada menu</td></tr>
            )}
            {menus.map((menu) => (
              <tr key={menu.id} className="border-b last:border-0 hover:bg-gray-50">
                <td className="px-4 py-3 font-mono text-xs">{menu.key_name}</td>
                <td className="px-4 py-3 font-medium">{menu.label}</td>
                <td className="px-4 py-3 text-gray-500 text-xs">
                  {menu.parent_id ? parentMap[menu.parent_id] ?? '-' : <span className="text-blue-500">Root</span>}
                </td>
                <td className="px-4 py-3 text-gray-500 font-mono text-xs">{menu.path ?? '-'}</td>
                <td className="px-4 py-3 text-center text-gray-500">{menu.order_index}</td>
                <td className="px-4 py-3 text-center">
                  <Badge variant={menu.is_active ? 'default' : 'secondary'}>
                    {menu.is_active ? 'Aktif' : 'Nonaktif'}
                  </Badge>
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center justify-center gap-1">
                    <RoleGuard menuKey="sistem.menus" action="can_edit">
                      <Button size="icon" variant="ghost" onClick={() => handleOpenEdit(menu)}>
                        <Pencil size={14} />
                      </Button>
                    </RoleGuard>
                    <RoleGuard menuKey="sistem.menus" action="can_delete">
                      <Button
                        size="icon"
                        variant="ghost"
                        className="text-red-500 hover:text-red-600"
                        onClick={() => handleOpenDelete(menu)}
                      >
                        <Trash2 size={14} />
                      </Button>
                    </RoleGuard>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <MenuFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        menuId={editingMenu?.id}
        allMenus={menus}
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
}
