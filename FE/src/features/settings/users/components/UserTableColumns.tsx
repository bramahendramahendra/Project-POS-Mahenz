import { KeyRound, Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'

import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { AppUser } from '../users.types'

const ROLE_BADGE_COLORS: Record<string, string> = {
  owner: 'bg-purple-100 text-purple-700',
  admin: 'bg-blue-100 text-blue-700',
  kasir: 'bg-green-100 text-green-700',
}

function roleBadgeClass(roleName: string): string {
  return ROLE_BADGE_COLORS[roleName] ?? 'bg-gray-100 text-gray-700'
}

export interface UserColumnHandlers {
  currentUserId?: number
  onEdit: (user: AppUser) => void
  onChangePassword: (user: AppUser) => void
  onToggleStatus: (user: AppUser) => void
  onDelete: (user: AppUser) => void
}

export function buildUserColumns(handlers: UserColumnHandlers): ColumnDef<AppUser>[] {
  const { currentUserId, onEdit, onChangePassword, onToggleStatus, onDelete } = handlers

  return [
    {
      key: 'username',
      header: 'Username',
      sortable: true,
      cell: (row) => <span className="font-mono text-sm text-gray-700">{row.username}</span>,
    },
    {
      key: 'full_name',
      header: 'Nama',
      sortable: true,
      cell: (row) => <span className="font-medium text-gray-800">{row.full_name}</span>,
    },
    {
      key: 'role_name',
      header: 'Role',
      sortable: true,
      cell: (row) => (
        <span
          className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${roleBadgeClass(row.role_name)}`}
        >
          {row.role_name}
        </span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      sortable: true,
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '160px',
      cell: (row) => {
        const isSelf = row.id === currentUserId
        return (
          <div className="flex items-center justify-center gap-1">
            <RoleGuard menuKey="sistem.users" action="can_edit">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-blue-600" onClick={() => onEdit(row)}>
                    <Pencil size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Edit</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-blue-600" onClick={() => onChangePassword(row)}>
                    <KeyRound size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Ganti Password</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className={`h-7 w-7 ${row.is_active ? 'text-gray-500 hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}`}
                    disabled={isSelf}
                    onClick={() => onToggleStatus(row)}
                  >
                    {row.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{row.is_active ? 'Nonaktifkan' : 'Aktifkan'}</TooltipContent>
              </Tooltip>
            </RoleGuard>
            <RoleGuard menuKey="sistem.users" action="can_delete">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 text-gray-500 hover:text-red-600"
                    disabled={isSelf}
                    onClick={() => onDelete(row)}
                  >
                    <Trash2 size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Hapus</TooltipContent>
              </Tooltip>
            </RoleGuard>
          </div>
        )
      },
    },
  ]
}
