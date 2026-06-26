import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { NonCashTransaction } from '../my-cash.types'

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function buildMyCashNonCashTransactionColumns(): ColumnDef<NonCashTransaction>[] {
  return [
    {
      key: 'transaction_date',
      header: 'Waktu',
      cell: (row) => (
        <span className="text-sm text-gray-600">{formatDateTime(row.transaction_date)}</span>
      ),
    },
    {
      key: 'transaction_code',
      header: 'Kode Transaksi',
      cell: (row) => (
        <span className="text-sm font-medium">{row.transaction_code}</span>
      ),
    },
    {
      key: 'customer_name',
      header: 'Pelanggan',
      cell: (row) => (
        <span className="text-sm text-gray-600">{row.customer_name || '—'}</span>
      ),
    },
    {
      key: 'payment_method_label',
      header: 'Metode Bayar',
      cell: (row) => (
        <span className="text-sm text-gray-700">{row.payment_method_label}</span>
      ),
    },
    {
      key: 'total_amount',
      header: 'Jumlah',
      align: 'right',
      cell: (row) => (
        <span className="text-sm font-semibold text-blue-600">
          {formatRupiah(row.total_amount)}
        </span>
      ),
    },
  ]
}
