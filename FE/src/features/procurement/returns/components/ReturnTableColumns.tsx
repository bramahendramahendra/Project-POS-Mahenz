import { Eye, Trash2 } from 'lucide-react'

import { StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatDate, formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierReturn } from '../returns.types'

export interface ReturnColumnHandlers {
  onDetail: (row: SupplierReturn) => void
  onDelete: (row: SupplierReturn) => void
}

export function buildReturnColumns(handlers: ReturnColumnHandlers): ColumnDef<SupplierReturn>[] {
  const { onDetail, onDelete } = handlers

  return [
    {
      key: 'return_date',
      header: 'Tanggal',
      sortable: true,
      cell: (row) => <span className="text-sm text-gray-600">{formatDate(row.return_date)}</span>,
    },
    {
      key: 'return_code',
      header: 'Kode Retur',
      cell: (row) => (
        <span className="text-sm font-mono font-medium text-blue-700">{row.return_code}</span>
      ),
    },
    {
      key: 'supplier_name',
      header: 'Supplier',
      sortable: true,
      cell: (row) => <span className="text-sm">{row.supplier_name}</span>,
    },
    {
      key: 'total_return_amount',
      header: 'Total Retur',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className="text-sm font-semibold text-red-600">
          {formatRupiah(row.total_return_amount)}
        </span>
      ),
    },
    {
      key: 'reason',
      header: 'Alasan',
      cell: (row) => <span className="text-sm text-gray-600">{row.reason}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      align: 'center',
      sortable: true,
      cell: (row) => <StatusBadge status={row.status} />,
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
          {row.status === 'pending' && (
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
