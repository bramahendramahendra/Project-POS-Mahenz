import { Eye, Trash2 } from 'lucide-react'

import { StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Supplier } from '../../suppliers/suppliers.types'
import type { SupplierPurchase } from '../purchases.types'

export interface PurchaseColumnHandlers {
  onDetail: (purchase: SupplierPurchase) => void
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


export function buildPurchaseColumns(
  handlers: PurchaseColumnHandlers,
  suppliers: Supplier[]
): ColumnDef<SupplierPurchase>[] {
  const { onDetail, onPay, onDelete } = handlers
  const supplierMap = new Map(suppliers.map((s) => [s.id, s.name]))

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
        <span className="text-sm">
          {row.supplier_id ? (supplierMap.get(row.supplier_id) ?? '-') : '-'}
        </span>
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
      key: 'id',
      header: 'Aksi',
      align: 'center',
      cell: (row) => (
        <div className="flex gap-1 justify-center">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="ghost" size="sm" onClick={() => onDetail(row)}>
                <Eye className="h-4 w-4" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Detail</TooltipContent>
          </Tooltip>
          {row.payment_status !== 'paid' && (
            <Button variant="outline" size="sm" onClick={() => onPay(row)} className="text-xs h-7 px-2">
              Bayar
            </Button>
          )}
          {row.paid_amount === 0 && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="sm" onClick={() => onDelete(row)}>
                  <Trash2 className="h-4 w-4 text-red-500" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Hapus</TooltipContent>
            </Tooltip>
          )}
        </div>
      ),
    },
  ]
}
