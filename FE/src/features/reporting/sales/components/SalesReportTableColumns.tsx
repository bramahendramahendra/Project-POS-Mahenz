import { formatRupiah, formatDate } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { SalesReport } from '../sales.types'

export function buildSalesReportColumns(): ColumnDef<SalesReport>[] {
  return [
    {
      key: 'transaction_date',
      header: 'Tanggal',
      sortable: true,
      cell: (r) => <span className="text-sm text-gray-600">{formatDate(r.transaction_date)}</span>,
    },
    {
      key: 'transaction_code',
      header: 'Kode Transaksi',
      sortable: true,
      cell: (r) => <span className="text-sm font-mono font-medium">{r.transaction_code}</span>,
    },
    {
      key: 'cashier_name',
      header: 'Kasir',
      sortable: true,
      cell: (r) => <span className="text-sm">{r.cashier_name}</span>,
    },
    {
      key: 'customer_name',
      header: 'Customer',
      sortable: true,
      cell: (r) => <span className="text-sm text-gray-500">{r.customer_name ?? '-'}</span>,
    },
    {
      key: 'total_amount',
      header: 'Total',
      align: 'right',
      sortable: true,
      cell: (r) => <span className="text-sm font-semibold">{formatRupiah(r.total_amount)}</span>,
    },
    {
      key: 'payment_method',
      header: 'Metode Bayar',
      sortable: true,
      cell: (r) => <span className="text-sm capitalize">{r.payment_method}</span>,
    },
    {
      key: 'status',
      header: 'Status',
      align: 'center',
      sortable: true,
      cell: (r) =>
        r.status === 'completed' ? (
          <span className="inline-flex rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-700">
            Selesai
          </span>
        ) : (
          <span className="inline-flex rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-medium text-red-700">
            Void
          </span>
        ),
    },
  ]
}
