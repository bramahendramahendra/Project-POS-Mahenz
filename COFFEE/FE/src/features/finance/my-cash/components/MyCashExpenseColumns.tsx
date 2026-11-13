import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { CashDrawerExpenseItem } from '../my-cash.types'

export function buildMyCashExpenseColumns(): ColumnDef<CashDrawerExpenseItem>[] {
  return [
    {
      key: 'category',
      header: 'Kategori',
      cell: (row) => (
        <span className="text-sm font-medium">{row.category || '—'}</span>
      ),
    },
    {
      key: 'description',
      header: 'Keterangan',
      cell: (row) => (
        <span className="text-sm text-gray-600">{row.description || '—'}</span>
      ),
    },
    {
      key: 'amount',
      header: 'Jumlah',
      align: 'right',
      cell: (row) => (
        <span className="text-sm font-semibold text-red-600">
          {formatRupiah(row.amount)}
        </span>
      ),
    },
  ]
}
