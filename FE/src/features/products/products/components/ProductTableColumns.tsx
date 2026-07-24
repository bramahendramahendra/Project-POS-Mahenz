import { Eye, Lock, LockOpen, Pencil, Printer, Trash2, TriangleAlert } from 'lucide-react'

import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import { calcMargin } from '../products.utils'
import type { Product } from '../products.types'
import type { ProductExpirySeverity } from '../expiry-batches.types'

export interface ProductColumnHandlers {
  onDetail: (product: Product) => void
  onEdit: (product: Product) => void
  onDelete: (product: Product) => void
  onLabel: (product: Product) => void
  onToggleStatus: (id: number, isActive: boolean) => void
  onOpenExpiry: (product: Product) => void
  expirySeverityMap: Record<number, ProductExpirySeverity>
}

export function buildProductColumns(handlers: ProductColumnHandlers): ColumnDef<Product>[] {
  const { onDetail, onEdit, onDelete, onLabel, onToggleStatus, onOpenExpiry, expirySeverityMap } = handlers

  return [
    {
      key: 'name',
      header: 'Nama Produk',
      sortable: true,
      mobileLabel: true,
      cell: (row) => {
        const expiry = expirySeverityMap[row.id]
        return (
          <div className="flex items-center gap-1.5">
            <span className="font-medium text-gray-800">{row.name}</span>
            {expiry && (
              <button
                type="button"
                onClick={() => onOpenExpiry(row)}
                title={`${expiry.warning_count} batch perlu dicek`}
                className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium ${
                  expiry.severity === 'expired'
                    ? 'bg-red-100 text-red-600 hover:bg-red-200'
                    : 'bg-amber-100 text-amber-700 hover:bg-amber-200'
                }`}
              >
                <TriangleAlert size={10} />
                {expiry.severity === 'expired' ? 'Expired' : 'Mendekati Expired'}
              </button>
            )}
          </div>
        )
      },
    },
    {
      key: 'barcode',
      header: 'Barcode / SKU',
      mobileHidden: true,
      cell: (row) => (
        <div className="flex flex-col gap-0.5">
          {row.barcode ? (
            <code className="text-xs text-gray-700">{row.barcode}</code>
          ) : (
            <span className="text-gray-400 text-xs">—</span>
          )}
          {row.sku ? (
            <span className="text-xs text-gray-500">{row.sku}</span>
          ) : (
            <span className="text-gray-300 text-xs">—</span>
          )}
        </div>
      ),
    },
    {
      key: 'category_name',
      header: 'Kategori',
      cell: (row) =>
        row.category_name ? (
          <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600">
            {row.category_name}
          </span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'purchase_price',
      header: 'Harga Beli',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className="text-sm">{formatRupiah(row.purchase_price)}</span>
      ),
    },
    {
      key: 'selling_price',
      header: 'Harga Jual',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className="font-medium">{formatRupiah(row.selling_price)}</span>
      ),
    },
    {
      key: 'margin',
      header: 'Margin',
      align: 'center',
      width: '80px',
      cell: (row) => {
        const m = calcMargin(row.purchase_price, row.selling_price)
        return (
          <span
            className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
              m >= 30
                ? 'bg-green-100 text-green-700'
                : m >= 15
                  ? 'bg-amber-100 text-amber-700'
                  : 'bg-red-100 text-red-600'
            }`}
          >
            {m}%
          </span>
        )
      },
    },
    {
      key: 'stock',
      header: 'Stok',
      align: 'right',
      width: '80px',
      sortable: true,
      cell: (row) => (
        <span
          className={`font-medium ${
            row.stock === 0
              ? 'text-red-600'
              : row.stock < row.min_stock
                ? 'text-amber-600'
                : 'text-gray-800'
          }`}
        >
          {row.stock}
        </span>
      ),
    },
    {
      key: 'unit_name',
      header: 'Satuan',
      width: '80px',
      mobileHidden: true,
      cell: (row) =>
        row.unit_name ? (
          <span className="text-sm text-gray-600">{row.unit_name}</span>
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
      cell: (row) => (
        <StatusBadge status={row.is_active ? 'active' : 'inactive'} />
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
              <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-blue-600" onClick={() => onDetail(row)}>
                <Eye size={14} />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Lihat Detail</TooltipContent>
          </Tooltip>
          <RoleGuard menuKey="produk.produk" action="can_edit">
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
                <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-indigo-600" onClick={() => onLabel(row)}>
                  <Printer size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Cetak Label</TooltipContent>
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
          <RoleGuard menuKey="produk.produk" action="can_delete">
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
