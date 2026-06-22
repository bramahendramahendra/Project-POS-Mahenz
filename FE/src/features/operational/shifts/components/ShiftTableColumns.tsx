import { Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Shift } from '../shifts.types'
import { formatShiftTime } from '../shifts.utils'

interface ShiftColumnHandlers {
  onEdit: (shift: Shift) => void
  onDelete: (shift: Shift) => void
  onToggleStatus: (id: number, isActive: boolean) => void
}

export function buildShiftColumns(handlers: ShiftColumnHandlers): ColumnDef<Shift>[] {
  const { onEdit, onDelete, onToggleStatus } = handlers

  return [
    {
      key: 'name',
      header: 'Nama Shift',
      sortable: true,
      cell: (row) => (
        <span className="font-medium text-gray-800">{row.name}</span>
      ),
    },
    {
      key: 'start_time',
      header: 'Jam Operasional',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatShiftTime(row.start_time, row.end_time)}</span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      sortable: true,
      cell: (row) => (
        <StatusBadge status={row.is_active ? 'active' : 'inactive'} />
      ),
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
              onClick={() => onEdit(row)}
              title="Edit"
            >
              <Pencil size={14} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className={`h-7 w-7 ${row.is_active ? 'text-gray-500 hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}`}
              onClick={() => onToggleStatus(row.id, row.is_active)}
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
              onClick={() => onDelete(row)}
              title="Hapus"
            >
              <Trash2 size={14} />
            </Button>
          </RoleGuard>
        </div>
      ),
    },
  ]
}
