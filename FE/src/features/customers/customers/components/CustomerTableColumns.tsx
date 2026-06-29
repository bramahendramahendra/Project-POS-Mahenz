import { Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'

import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Customer } from '../customers.types'

export interface CustomerColumnHandlers {
  onEdit: (customer: Customer) => void
  onDelete: (customer: Customer) => void
  onToggleStatus: (id: number, isActive: boolean) => void
}

export function buildCustomerColumns(handlers: CustomerColumnHandlers): ColumnDef<Customer>[] {
  const { onEdit, onDelete, onToggleStatus } = handlers

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
      sortable: true,
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
      width: '120px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey="pelanggan.pelanggan" action="can_edit">
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
                <Button variant="ghost" size="icon" className={`h-7 w-7 ${row.is_active ? 'text-gray-500 hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}`} onClick={() => onToggleStatus(row.id, row.is_active)}>
                  {row.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
                </Button>
              </TooltipTrigger>
              <TooltipContent>{row.is_active ? 'Nonaktifkan' : 'Aktifkan'}</TooltipContent>
            </Tooltip>
          </RoleGuard>
          <RoleGuard menuKey="pelanggan.pelanggan" action="can_delete">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-red-600" onClick={() => onDelete(row)}>
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
