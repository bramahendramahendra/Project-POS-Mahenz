import { Eye, Pencil, Trash2 } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierPurchase } from '../purchases.types'

export interface PurchaseColumnHandlers {
  onDetail: (purchase: SupplierPurchase) => void
  onEdit: (purchase: SupplierPurchase) => void
  onPay: (purchase: SupplierPurchase) => void
  onDelete: (purchase: SupplierPurchase) => void
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

export function buildPurchaseColumns(handlers: PurchaseColumnHandlers): ColumnDef<SupplierPurchase>[] {
  const { onDetail, onEdit, onPay, onDelete } = handlers

  return [
    {
      key: 'purchase_code',
      header: 'Kode PO',
      cell: (row) => (
        <span className="text-sm font-mono font-medium text-blue-700">{row.purchase_code}</span>
      ),
    },
    {
      key: 'purchase_date',
      header: 'Tanggal',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDate(row.purchase_date)}</span>
      ),
    },
    {
      key: 'invoice_number',
      header: 'No. Faktur',
      cell: (row) => <span className="text-sm font-medium">{row.invoice_number}</span>,
    },
    {
      key: 'supplier_name',
      header: 'Supplier',
      cell: (row) => (
        <span className="text-sm">{row.supplier_name || '-'}</span>
      ),
    },
    {
      key: 'total_amount',
      header: 'Total',
      align: 'right',
      cell: (row) => (
        <span className="text-sm font-semibold">{formatRupiah(row.total_amount)}</span>
      ),
    },
    {
      key: 'payment_status',
      header: 'Status',
      align: 'center',
      cell: (row) => <StatusBadge status={row.payment_status} />,
    },
    {
      key: 'remaining_amount',
      header: 'Sisa Hutang',
      align: 'right',
      cell: (row) => (
        <span
          className={`text-sm ${row.remaining_amount > 0 ? 'text-red-600 font-medium' : 'text-gray-400'}`}
        >
          {formatRupiah(row.remaining_amount)}
        </span>
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '140px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 text-gray-500 hover:text-blue-600"
                onClick={() => onDetail(row)}
              >
                <Eye size={14} />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Detail</TooltipContent>
          </Tooltip>
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            {row.paid_amount === 0 && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 text-gray-500 hover:text-indigo-600"
                    onClick={() => onEdit(row)}
                  >
                    <Pencil size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Edit</TooltipContent>
              </Tooltip>
            )}
            {row.payment_status !== 'paid' && (
              <Button
                variant="outline"
                size="sm"
                className="h-7 px-2 text-xs"
                onClick={() => onPay(row)}
              >
                Bayar
              </Button>
            )}
            {row.paid_amount === 0 && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 text-gray-500 hover:text-red-600"
                    onClick={() => onDelete(row)}
                  >
                    <Trash2 size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Hapus</TooltipContent>
              </Tooltip>
            )}
          </RoleGuard>
        </div>
      ),
    },
  ]
}
