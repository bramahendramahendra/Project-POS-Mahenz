import { Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Category } from '../categories.types'

export interface CategoryColumnHandlers {
  onEdit: (category: Category) => void
  onDelete: (category: Category) => void
  onToggleStatus: (category: Category) => void
}

export function buildCategoryColumns(handlers: CategoryColumnHandlers): ColumnDef<Category>[] {
  const { onEdit, onDelete, onToggleStatus } = handlers

  return [
    {
      key: 'code',
      header: 'Kode',
      width: '80px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Kategori',
      cell: (row) => (
        <span className="font-medium text-gray-800">{row.name}</span>
      ),
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
      key: 'product_count',
      header: 'Jumlah Produk',
      align: 'center',
      width: '130px',
      cell: (row) => (
        <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600">
          {row.product_count} produk
        </span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
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
              onClick={() => onToggleStatus(row)}
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
