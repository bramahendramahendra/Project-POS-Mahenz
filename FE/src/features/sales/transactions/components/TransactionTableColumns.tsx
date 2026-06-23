import { Ban, Eye } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Transaction } from '../transactions.types'
import { PAYMENT_LABELS, formatDateTimeShort } from '../transactions.utils'

export interface TransactionColumnHandlers {
  onDetail: (transaction: Transaction) => void
  onVoid: (transaction: Transaction) => void
}

export function buildTransactionColumns({ onDetail, onVoid }: TransactionColumnHandlers): ColumnDef<Transaction>[] {
  return [
    {
      key: 'transaction_code',
      header: 'Kode',
      cell: (row) => (
        <span className="font-mono font-semibold text-gray-800 text-sm">
          {row.transaction_code}
        </span>
      ),
    },
    {
      key: 'transaction_date',
      header: 'Tanggal',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDateTimeShort(row.transaction_date)}</span>
      ),
    },
    {
      key: 'customer_name',
      header: 'Pelanggan',
      cell: (row) =>
        row.customer_name ? (
          <span className="text-sm">{row.customer_name}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'kasir_name',
      header: 'Kasir',
      cell: (row) => <span className="text-sm">{row.kasir_name}</span>,
    },
    {
      key: 'total_amount',
      header: 'Total',
      align: 'right',
      cell: (row) => <span className="font-semibold">{formatRupiah(row.total_amount)}</span>,
    },
    {
      key: 'payment_method',
      header: 'Metode',
      align: 'center',
      cell: (row) => (
        <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600">
          {PAYMENT_LABELS[row.payment_method]}
        </span>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      align: 'center',
      cell: (row) => <StatusBadge status={row.status === 'completed' ? 'success' : 'error'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '90px',
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
          <RoleGuard allowedRoles={[ROLES.OWNER]}>
            {row.status === 'completed' && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 text-gray-500 hover:text-red-600"
                    onClick={() => onVoid(row)}
                  >
                    <Ban size={14} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Void</TooltipContent>
              </Tooltip>
            )}
          </RoleGuard>
        </div>
      ),
    },
  ]
}
