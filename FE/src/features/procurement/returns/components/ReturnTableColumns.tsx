import { Eye, Trash2 } from 'lucide-react'

import { StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SupplierReturn } from '../returns.types'

export interface ReturnColumnHandlers {
  onDetail: (row: SupplierReturn) => void
  onDelete: (row: SupplierReturn) => void
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}


export function buildReturnColumns(handlers: ReturnColumnHandlers): ColumnDef<SupplierReturn>[] {
  const { onDetail, onDelete } = handlers

  return [
    {
      key: 'return_date',
      header: 'Tanggal',
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
      cell: (row) => <span className="text-sm">{row.supplier_name}</span>,
    },
    {
      key: 'total_return_amount',
      header: 'Total Retur',
      align: 'right',
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
      cell: (row) => <StatusBadge status={row.status} />,
    },
    {
      key: 'id',
      header: 'Aksi',
      align: 'center',
      cell: (row) => (
        <div className="flex gap-1 justify-center">
          <Button variant="ghost" size="sm" onClick={() => onDetail(row)} title="Detail">
            <Eye className="h-4 w-4" />
          </Button>
          {row.status === 'pending' && (
            <Button variant="ghost" size="sm" onClick={() => onDelete(row)} title="Hapus">
              <Trash2 className="h-4 w-4 text-red-500" />
            </Button>
          )}
        </div>
      ),
    },
  ]
}
