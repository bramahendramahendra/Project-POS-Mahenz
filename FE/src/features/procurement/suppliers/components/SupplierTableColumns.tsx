import { Eye, Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Supplier } from '../suppliers.types'

export interface SupplierColumnHandlers {
  onDetail: (id: number) => void
  onEdit: (supplier: Supplier) => void
  onDelete: (supplier: Supplier) => void
  onToggleStatus: (id: number, isActive: boolean) => void
}

export function buildSupplierColumns(handlers: SupplierColumnHandlers): ColumnDef<Supplier>[] {
  const { onDetail, onEdit, onDelete, onToggleStatus } = handlers

  return [
    {
      key: 'supplier_code',
      header: 'Kode',
      width: '110px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.supplier_code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Supplier',
      sortable: true,
      cell: (row) => (
        <span className="font-medium text-gray-800">{row.name}</span>
      ),
    },
    {
      key: 'contact_person',
      header: 'Nama Kontak',
      cell: (row) =>
        row.contact_person ? (
          <span>{row.contact_person}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
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
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '90px',
      sortable: true,
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '140px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-gray-500 hover:text-indigo-600"
            onClick={() => onDetail(row.id)}
            title="Detail"
          >
            <Eye size={14} />
          </Button>
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
