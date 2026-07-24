import { Banknote, Eye, PackagePlus, Pencil, Trash2, Ban } from 'lucide-react'

import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatDate, formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierPurchase } from '../purchases.types'

export interface PurchaseColumnHandlers {
  onDetail: (purchase: SupplierPurchase) => void
  onEdit: (purchase: SupplierPurchase) => void
  onPay: (purchase: SupplierPurchase) => void
  onDelete: (purchase: SupplierPurchase) => void
  onVoid: (purchase: SupplierPurchase) => void
  onAddItems: (purchase: SupplierPurchase) => void
}

export function buildPurchaseColumns(handlers: PurchaseColumnHandlers): ColumnDef<SupplierPurchase>[] {
  const { onDetail, onEdit, onPay, onDelete, onVoid, onAddItems } = handlers

  return [
    {
      key: 'purchase_code',
      header: 'Kode PO',
      mobileLabel: true,
      cell: (row) => (
        <span className="text-sm font-mono font-medium text-blue-700">{row.purchase_code}</span>
      ),
    },
    {
      key: 'purchase_date',
      header: 'Tanggal',
      sortable: true,
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDate(row.purchase_date)}</span>
      ),
    },
    {
      key: 'invoice_number',
      header: 'No. Faktur',
      mobileHidden: true,
      cell: (row) => <span className="text-sm font-medium">{row.invoice_number}</span>,
    },
    {
      key: 'supplier_name',
      header: 'Supplier',
      sortable: true,
      cell: (row) => (
        <span className="text-sm">{row.supplier_name || '-'}</span>
      ),
    },
    {
      key: 'total_amount',
      header: 'Total',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className="text-sm font-semibold">{formatRupiah(row.total_amount)}</span>
      ),
    },
    {
      key: 'payment_status',
      header: 'Status',
      align: 'center',
      sortable: true,
      cell: (row) => (
        <StatusBadge status={row.status === 'void' ? 'void' : row.payment_status} />
      ),
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
      cell: (row) => {
        const isVoid = row.status === 'void'
        return (
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

            {!isVoid && (
              <RoleGuard menuKey="pengadaan.pembelian" action="can_edit">
                {row.paid_amount === 0 && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-gray-500 hover:text-blue-600"
                        onClick={() => onEdit(row)}
                      >
                        <Pencil size={14} />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Edit</TooltipContent>
                  </Tooltip>
                )}
                {row.payment_status !== 'unpaid' && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-gray-500 hover:text-teal-600"
                        onClick={() => onAddItems(row)}
                      >
                        <PackagePlus size={14} />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Tambah Item</TooltipContent>
                  </Tooltip>
                )}
                {row.payment_status !== 'paid' && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-gray-500 hover:text-green-600"
                        onClick={() => onPay(row)}
                      >
                        <Banknote size={14} />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Bayar</TooltipContent>
                  </Tooltip>
                )}
              </RoleGuard>
            )}

            {!isVoid && row.paid_amount === 0 && (
              <RoleGuard menuKey="pengadaan.pembelian" action="can_delete">
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
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-gray-500 hover:text-orange-600"
                      onClick={() => onVoid(row)}
                    >
                      <Ban size={14} />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Void</TooltipContent>
                </Tooltip>
              </RoleGuard>
            )}
          </div>
        )
      },
    },
  ]
}
