import { PageHeader } from '@/shared/components'

import { TransactionTable } from './components/TransactionTable'

export function TransactionsPage() {
  return (
    <div className="space-y-4">
      <PageHeader
        title="Transaksi"
        breadcrumbs={[{ label: 'Penjualan' }, { label: 'Transaksi' }]}
      />
      <TransactionTable />
    </div>
  )
}
