import { formatDate, todayStr } from '@/shared/utils'

import type { SalesReport } from './sales.types'

export function exportSalesCSV(data: SalesReport[]): void {
  const headers = ['Tanggal', 'Kode Transaksi', 'Kasir', 'Customer', 'Total', 'Diskon', 'Metode Bayar', 'Status']
  const rows = data.map((r) => [
    formatDate(r.transaction_date),
    r.transaction_code,
    r.cashier_name,
    r.customer_name || '-',
    r.total_amount,
    r.discount,
    r.payment_method,
    r.status,
  ])
  const csv = [headers, ...rows].map((row) => row.join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `laporan-penjualan-${todayStr()}.csv`
  a.click()
  URL.revokeObjectURL(url)
}
