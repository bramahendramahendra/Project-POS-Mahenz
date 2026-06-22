import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Shift } from '../shifts.types'
import { formatShiftTime } from '../shifts.utils'

interface ShiftColumnHandlers {
  onEdit: (shift: Shift) => void
  onDelete: (shift: Shift) => void
  onToggleStatus: (shift: Shift) => void
}

export function buildShiftColumns({ onEdit, onDelete, onToggleStatus }: ShiftColumnHandlers): ColumnDef<Shift>[] {
  return [
    {
      key: 'name',
      header: 'Nama Shift',
      cell: (row) => <span className="font-medium">{row.name}</span>,
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
      cell: (row) =>
        row.is_active ? (
          <Badge variant="default" className="bg-green-600">Aktif</Badge>
        ) : (
          <Badge variant="secondary">Nonaktif</Badge>
        ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '160px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Button
            size="sm"
            variant="outline"
            className="h-7 px-3 text-xs"
            onClick={() => onEdit(row)}
          >
            Edit
          </Button>
          <Button
            size="sm"
            variant="outline"
            className={`h-7 px-3 text-xs ${
              row.is_active
                ? 'text-yellow-600 border-yellow-200 hover:bg-yellow-50'
                : 'text-green-600 border-green-200 hover:bg-green-50'
            }`}
            onClick={() => onToggleStatus(row)}
          >
            {row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
          </Button>
          <Button
            size="sm"
            variant="outline"
            className="h-7 px-3 text-xs text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
            onClick={() => onDelete(row)}
          >
            Hapus
          </Button>
        </div>
      ),
    },
  ]
}
