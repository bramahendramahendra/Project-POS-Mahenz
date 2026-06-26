import { Eye, XCircle } from 'lucide-react'

import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { CashDrawer } from '../cash-drawer.types'

export interface CashDrawerColumnHandlers {
  onRowClick: (row: CashDrawer) => void
  onForceClose?: (row: CashDrawer) => void
  canForceClose?: boolean
}

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function buildCashDrawerColumns(
  handlers: CashDrawerColumnHandlers
): ColumnDef<CashDrawer>[] {
  const { onRowClick, onForceClose, canForceClose } = handlers

  return [
    {
      key: 'open_time',
      header: 'Waktu Buka',
      cell: (row) => <span className="text-gray-600 text-sm">{formatDateTime(row.open_time)}</span>,
    },
    {
      key: 'shift_name',
      header: 'Shift',
      align: 'center',
      cell: (row) => (
        <span className="text-sm text-gray-600">{row.shift_name ?? '—'}</span>
      ),
    },
    {
      key: 'user_name',
      header: 'Kasir',
      cell: (row) => <span className="text-sm font-medium">{row.user_name}</span>,
    },
    {
      key: 'opening_balance',
      header: 'Saldo Awal Tunai',
      align: 'right',
      cell: (row) => <span className="text-sm">{formatRupiah(row.opening_balance)}</span>,
    },
    {
      key: 'total_cash_sales',
      header: 'Total Masuk',
      align: 'right',
      cell: (row) => (
        <span className="text-green-600 font-medium text-sm">{formatRupiah(row.total_cash_sales)}</span>
      ),
    },
    {
      key: 'total_expenses',
      header: 'Total Keluar',
      align: 'right',
      cell: (row) => (
        <span className="text-red-600 font-medium text-sm">{formatRupiah(row.total_expenses)}</span>
      ),
    },
    {
      key: 'closing_balance',
      header: 'Saldo Akhir Tunai',
      align: 'right',
      cell: (row) => (
        <span className="font-semibold text-sm">
          {row.status === 'closed' && row.closing_balance != null
            ? formatRupiah(row.closing_balance)
            : '—'}
        </span>
      ),
    },
    {
      key: 'difference',
      header: 'Selisih',
      align: 'right',
      cell: (row) => {
        const diff = row.difference ?? 0
        return (
          <span
            className={`text-sm font-medium ${
              diff === 0 ? 'text-gray-500' : diff > 0 ? 'text-green-600' : 'text-red-600'
            }`}
          >
            {row.status === 'closed' ? `${diff > 0 ? '+' : ''}${formatRupiah(diff)}` : '—'}
          </span>
        )
      },
    },
    {
      key: 'status',
      header: 'Status',
      align: 'center',
      cell: (row) =>
        row.status === 'closed' ? (
          <Badge variant="secondary">Tutup</Badge>
        ) : (
          <Badge variant="default" className="bg-green-600">Buka</Badge>
        ),
    },
    {
      key: 'id',
      header: 'Aksi',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 text-gray-500 hover:text-blue-600"
                onClick={() => onRowClick(row)}
              >
                <Eye size={14} />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Lihat Detail</TooltipContent>
          </Tooltip>
          {canForceClose && row.status === 'open' && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 text-gray-500 hover:text-red-600"
                  onClick={() => onForceClose?.(row)}
                >
                  <XCircle size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Tutup Paksa</TooltipContent>
            </Tooltip>
          )}
        </div>
      ),
    },
  ]
}
