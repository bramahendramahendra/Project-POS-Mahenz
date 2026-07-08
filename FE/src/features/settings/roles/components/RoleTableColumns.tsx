import { Pencil, Settings, Trash2 } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import { Switch } from '@/shared/components/ui/switch'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Role } from '../roles.types'

export interface RoleColumnHandlers {
  onManageAccess: (role: Role) => void
  onEdit: (role: Role) => void
  onDelete: (role: Role) => void
  onToggleStatus: (id: number) => void
}

export function buildRoleColumns(handlers: RoleColumnHandlers): ColumnDef<Role>[] {
  const { onManageAccess, onEdit, onDelete, onToggleStatus } = handlers

  return [
    {
      key: 'name',
      header: 'Nama Role',
      cell: (row) => <span className="font-mono text-xs font-medium">{row.name}</span>,
    },
    {
      key: 'display_name',
      header: 'Label',
      cell: (row) => <span className="font-medium text-gray-800">{row.display_name}</span>,
    },
    {
      key: 'description',
      header: 'Deskripsi',
      cell: (row) =>
        row.description ? (
          <span className="text-sm text-gray-600">{row.description}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'is_system',
      header: 'Tipe',
      align: 'center',
      width: '110px',
      cell: (row) =>
        row.is_system ? <Badge variant="secondary">Sistem</Badge> : <Badge variant="outline">Custom</Badge>,
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <Switch checked={row.is_active} onCheckedChange={() => onToggleStatus(row.id)} disabled={row.is_system} />
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '130px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button size="icon" variant="ghost" className="h-7 w-7" onClick={() => onManageAccess(row)}>
                <Settings size={14} />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Atur Akses Menu</TooltipContent>
          </Tooltip>
          {!row.is_system && (
            <RoleGuard menuKey="sistem.roles" action="can_edit">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button size="icon" variant="ghost" className="h-7 w-7" onClick={() => onEdit(row)}>
                    <Pencil size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Edit</TooltipContent>
              </Tooltip>
            </RoleGuard>
          )}
          {!row.is_system && (
            <RoleGuard menuKey="sistem.roles" action="can_delete">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-red-500 hover:text-red-600"
                    onClick={() => onDelete(row)}
                  >
                    <Trash2 size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Hapus</TooltipContent>
              </Tooltip>
            </RoleGuard>
          )}
        </div>
      ),
    },
  ]
}
