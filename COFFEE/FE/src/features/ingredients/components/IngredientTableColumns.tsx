import { History, Lock, LockOpen, MinusCircle, PackageX, Pencil, Trash2, TriangleAlert } from 'lucide-react'

import { RoleGuard, StatusBadge } from '@/shared/components'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Ingredient } from '../ingredients.types'
import { formatQty } from '../ingredients.utils'
import { IconAction } from './IconAction'

export interface IngredientColumnHandlers {
  onUsage: (row: Ingredient) => void
  onWaste: (row: Ingredient) => void
  onHistory: (row: Ingredient) => void
  onEdit: (row: Ingredient) => void
  onToggleStatus: (row: Ingredient) => void
  onDelete: (row: Ingredient) => void
}

export function buildIngredientColumns(handlers: IngredientColumnHandlers): ColumnDef<Ingredient>[] {
  const { onUsage, onWaste, onHistory, onEdit, onToggleStatus, onDelete } = handlers

  return [
    {
      key: 'name',
      header: 'Nama Bahan',
      sortable: true,
      mobileLabel: true,
      cell: (row) => (
        <div className="flex flex-wrap items-center gap-1.5">
          <span className="font-medium text-gray-800">{row.name}</span>
          {row.is_low_stock && (
            <span className="inline-flex items-center gap-1 rounded-full bg-amber-100 px-2 py-0.5 text-[11px] font-medium text-amber-700">
              <TriangleAlert size={10} />
              Menipis
            </span>
          )}
        </div>
      ),
    },
    {
      key: 'stock',
      header: 'Stok Sisa',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className={`font-semibold ${row.stock < 0 ? 'text-red-600' : 'text-gray-800'}`}>
          {formatQty(row.stock, row.unit_abbr)}
        </span>
      ),
    },
    {
      key: 'total_used',
      header: 'Total Terpakai',
      align: 'right',
      sortable: true,
      cell: (row) => <span className="text-gray-700">{formatQty(row.total_used, row.unit_abbr)}</span>,
    },
    {
      key: 'total_wasted',
      header: 'Total Terbuang',
      align: 'right',
      sortable: true,
      cell: (row) => <span className="text-gray-700">{formatQty(row.total_wasted, row.unit_abbr)}</span>,
    },
    {
      key: 'min_stock',
      header: 'Batas Menipis',
      align: 'right',
      mobileHidden: true,
      cell: (row) =>
        row.min_stock == null ? (
          <span className="text-gray-400 text-sm">—</span>
        ) : (
          <span className="text-gray-600">{formatQty(row.min_stock, row.unit_abbr)}</span>
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
      width: '200px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          {row.is_active && (
            <RoleGuard menuKey="bahan.riwayat" action="can_create">
              <IconAction label="Catat Terpakai" onClick={() => onUsage(row)} className="hover:text-blue-600">
                <MinusCircle size={14} />
              </IconAction>
              <IconAction label="Catat Terbuang" onClick={() => onWaste(row)} className="hover:text-red-600">
                <PackageX size={14} />
              </IconAction>
            </RoleGuard>
          )}
          <RoleGuard menuKey="bahan.riwayat" action="can_view">
            <IconAction label="Riwayat Stok" onClick={() => onHistory(row)} className="hover:text-blue-600">
              <History size={14} />
            </IconAction>
          </RoleGuard>
          <RoleGuard menuKey="bahan.daftar" action="can_edit">
            <IconAction label="Edit" onClick={() => onEdit(row)} className="hover:text-blue-600">
              <Pencil size={14} />
            </IconAction>
            <IconAction
              label={row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
              onClick={() => onToggleStatus(row)}
              className={row.is_active ? 'hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}
            >
              {row.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
            </IconAction>
          </RoleGuard>
          {!row.has_history && (
            <RoleGuard menuKey="bahan.daftar" action="can_delete">
              <IconAction label="Hapus Bahan" onClick={() => onDelete(row)} className="hover:text-red-600">
                <Trash2 size={14} />
              </IconAction>
            </RoleGuard>
          )}
        </div>
      ),
    },
  ]
}
