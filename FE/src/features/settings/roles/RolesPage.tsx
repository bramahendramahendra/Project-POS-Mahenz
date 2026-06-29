import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Pencil, Settings, Shield, Trash2 } from 'lucide-react'

import { ConfirmDialog, PageHeader } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import { Switch } from '@/shared/components/ui/switch'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { Input } from '@/shared/components/ui/input'
import { useDebounce, useDisclosure } from '@/shared/hooks'
import { ROUTES } from '@/shared/constants/routes'

import { useDeleteRoleMutation, useRoleListQuery, useToggleRoleStatusMutation } from './roles.api'
import type { Role } from './roles.types'
import { RoleFormModal } from './components/RoleFormModal'

export function RolesPage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [deletingRole, setDeletingRole] = useState<Role | null>(null)

  const { data: roles = [], isLoading } = useRoleListQuery(
    debouncedSearch ? { search: debouncedSearch } : undefined
  )
  const { mutate: deleteRole, isPending: isDeleting } = useDeleteRoleMutation()
  const { mutate: toggleStatus } = useToggleRoleStatusMutation()

  const handleOpenAdd = () => {
    setEditingRole(null)
    openForm()
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

  return (
    <div className="space-y-4">
      <PageHeader
        title="Manajemen Role"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Manajemen Role' }]}
        actions={
          <Button onClick={handleOpenAdd}>
            <Shield size={14} className="mr-2" />
            Tambah Role
          </Button>
        }
      />

      <div className="flex items-center gap-3">
        <Input
          placeholder="Cari role..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-xs"
        />
      </div>

      <div className="rounded-lg border bg-white overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Nama Role</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Label</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Deskripsi</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Tipe</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Status</th>
              <th className="text-center px-4 py-3 font-medium text-gray-600">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {isLoading && (
              <tr><td colSpan={6} className="text-center py-8 text-gray-400">Memuat...</td></tr>
            )}
            {!isLoading && roles.length === 0 && (
              <tr><td colSpan={6} className="text-center py-8 text-gray-400">Belum ada role</td></tr>
            )}
            {roles.map((role) => (
              <tr key={role.id} className="border-b last:border-0 hover:bg-gray-50">
                <td className="px-4 py-3 font-mono text-xs font-medium">{role.name}</td>
                <td className="px-4 py-3 font-medium">{role.display_name}</td>
                <td className="px-4 py-3 text-gray-500">{role.description ?? '-'}</td>
                <td className="px-4 py-3 text-center">
                  {role.is_system
                    ? <Badge variant="secondary">Sistem</Badge>
                    : <Badge variant="outline">Custom</Badge>
                  }
                </td>
                <td className="px-4 py-3 text-center">
                  <Switch
                    checked={role.is_active}
                    onCheckedChange={() => toggleStatus(role.id)}
                    disabled={role.is_system}
                  />
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center justify-center gap-1">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() => navigate(ROUTES.SETTINGS_ROLES_ACCESS.replace(':id', String(role.id)))}
                        >
                          <Settings size={14} />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Atur Akses Menu</TooltipContent>
                    </Tooltip>
                    {!role.is_system && (
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button size="icon" variant="ghost" onClick={() => handleOpenEdit(role)}>
                            <Pencil size={14} />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Edit</TooltipContent>
                      </Tooltip>
                    )}
                    {!role.is_system && (
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            size="icon"
                            variant="ghost"
                            className="text-red-500 hover:text-red-600"
                            onClick={() => handleOpenDelete(role)}
                          >
                            <Trash2 size={14} />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Hapus</TooltipContent>
                      </Tooltip>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

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
}
