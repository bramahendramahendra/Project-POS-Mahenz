import { Pencil, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Customer } from '../customers.types'

export interface CustomerColumnHandlers {
  onEdit: (customer: Customer) => void
  onDelete: (customer: Customer) => void
}

export function buildCustomerColumns(handlers: CustomerColumnHandlers): ColumnDef<Customer>[] {
  const { onEdit, onDelete } = handlers

  return [
    {
      key: 'customer_code',
      header: 'Kode',
      width: '90px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.customer_code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Pelanggan',
      cell: (row) => <span className="font-medium text-gray-800">{row.name}</span>,
    },
    {
      key: 'phone',
      header: 'Telepon',
      cell: (row) =>
        row.phone ? (
          <span className="font-mono text-sm">{row.phone}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'address',
      header: 'Alamat',
      cell: (row) =>
        row.address ? (
          <span className="text-sm text-gray-600">{row.address}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '100px',
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
