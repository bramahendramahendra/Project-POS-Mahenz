import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { CashierPerformance } from '../cashier-performance.types'

export function buildCashierPerformanceColumns(): ColumnDef<CashierPerformance>[] {
  return [
    {
      key: 'cashier_name',
      header: 'Nama Kasir',
      cell: (r) => <span className="font-medium">{r.cashier_name}</span>,
    },
    {
      key: 'total_transactions',
      header: 'Jml Transaksi',
      align: 'right',
      cell: (r) => <span>{r.total_transactions}</span>,
    },
    {
      key: 'total_sales',
      header: 'Total Penjualan',
      align: 'right',
      cell: (r) => <span className="font-semibold text-green-600">{formatRupiah(r.total_sales)}</span>,
    },
    {
      key: 'avg_per_transaction',
      header: 'Rata-rata/Transaksi',
      align: 'right',
      cell: (r) => <span className="text-gray-600">{formatRupiah(r.avg_per_transaction)}</span>,
    },
    {
      key: 'void_count',
      header: 'Void',
      align: 'right',
      cell: (r) =>
        r.void_count > 0 ? (
          <span className="inline-flex rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium text-red-600">
            {r.void_count}
          </span>
        ) : (
          <span className="text-gray-400">0</span>
        ),
    },
  ]
}
