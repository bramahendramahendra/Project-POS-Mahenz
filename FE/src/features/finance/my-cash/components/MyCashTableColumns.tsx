import { Badge } from '@/shared/components/ui/badge'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { MyCashTransaction } from '../my-cash.types'

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function buildMyCashColumns(): ColumnDef<MyCashTransaction>[] {
  return [
    {
      key: 'created_at',
      header: 'Tanggal',
      cell: (row) => <span className="text-gray-600">{formatDateTime(row.created_at)}</span>,
    },
    {
      key: 'type',
      header: 'Jenis',
      width: '120px',
      cell: (row) =>
        row.type === 'receive' ? (
          <Badge variant="default">Terima</Badge>
        ) : (
          <Badge variant="secondary">Kembalikan</Badge>
        ),
    },
    {
      key: 'amount',
      header: 'Jumlah',
      align: 'right',
      width: '140px',
      cell: (row) => (
        <span
          className={`font-medium ${row.type === 'receive' ? 'text-green-600' : 'text-red-600'}`}
        >
          {row.type === 'receive' ? '+' : '-'}
          {formatRupiah(row.amount)}
        </span>
      ),
    },
    {
      key: 'notes',
      header: 'Catatan',
      cell: (row) => <span className="text-gray-500">{row.notes ?? '-'}</span>,
    },
    {
      key: 'created_by_name',
      header: 'Oleh',
      cell: (row) => <span className="text-gray-500">{row.created_by_name}</span>,
    },
  ]
}
