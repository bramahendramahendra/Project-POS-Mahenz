import { Pencil, Trash2 } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { MenuResponse } from '@/features/menu/menu.types'

export interface MenuColumnHandlers {
  parentMap: Record<number, string>
  onEdit: (menu: MenuResponse) => void
  onDelete: (menu: MenuResponse) => void
}

export function buildMenuColumns(handlers: MenuColumnHandlers): ColumnDef<MenuResponse>[] {
  const { parentMap, onEdit, onDelete } = handlers

  return [
    {
      key: 'key_name',
      header: 'Key',
      cell: (row) => <span className="font-mono text-xs">{row.key_name}</span>,
    },
    {
      key: 'label',
      header: 'Label',
      cell: (row) => <span className="font-medium text-gray-800">{row.label}</span>,
    },
    {
      key: 'parent_id',
      header: 'Parent',
      cell: (row) =>
        row.parent_id ? (
          <span className="text-gray-500 text-xs">{parentMap[row.parent_id] ?? '-'}</span>
        ) : (
          <span className="text-blue-500 text-xs">Root</span>
        ),
    },
    {
      key: 'path',
      header: 'Path',
      cell: (row) =>
        row.path ? (
          <span className="font-mono text-xs text-gray-500">{row.path}</span>
        ) : (
          <span className="text-gray-400 text-xs">-</span>
        ),
    },
    {
      key: 'order_index',
      header: 'Urutan',
      align: 'center',
      width: '90px',
      cell: (row) => <span className="text-gray-500">{row.order_index}</span>,
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <Badge variant={row.is_active ? 'default' : 'secondary'}>
          {row.is_active ? 'Aktif' : 'Nonaktif'}
        </Badge>
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey="sistem.menus" action="can_edit">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button size="icon" variant="ghost" className="h-7 w-7" onClick={() => onEdit(row)}>
                  <Pencil size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Edit</TooltipContent>
            </Tooltip>
          </RoleGuard>
          <RoleGuard menuKey="sistem.menus" action="can_delete">
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
        </div>
      ),
    },
  ]
}
