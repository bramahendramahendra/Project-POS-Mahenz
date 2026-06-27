import { formatDate, formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { CashflowItem } from '../finance.types'

export function buildFinanceColumns(): ColumnDef<CashflowItem>[] {
  return [
    {
      key: 'date',
      header: 'Tanggal',
      cell: (row) => <span className="text-gray-600 text-sm">{formatDate(row.date)}</span>,
    },
    {
      key: 'category',
      header: 'Kategori',
      cell: (row) => <span className="font-medium text-sm">{row.category}</span>,
    },
    {
      key: 'description',
      header: 'Deskripsi',
      cell: (row) => <span className="text-gray-600 text-sm">{row.description}</span>,
    },
    {
      key: 'type',
      header: 'Tipe',
      align: 'center',
      cell: (row) =>
        row.type === 'income' ? (
          <span className="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-700">
            Pemasukan
          </span>
        ) : (
          <span className="inline-flex items-center rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-medium text-red-700">
            Pengeluaran
          </span>
        ),
    },
    {
      key: 'amount',
      header: 'Nominal',
      align: 'right',
      cell: (row) => (
        <span className={row.type === 'income' ? 'font-semibold text-green-600' : 'font-semibold text-red-600'}>
          {row.type === 'income' ? '+' : '-'}
          {formatRupiah(row.amount)}
        </span>
      ),
    },
  ]
}
