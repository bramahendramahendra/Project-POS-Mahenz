import { PageHeader } from '@/shared/components'

import { FinanceTable } from './components/FinanceTable'

export function FinancePage() {
  return (
    <div className="space-y-4">
      <PageHeader title="Keuangan" breadcrumbs={[{ label: 'Finance' }, { label: 'Keuangan' }]} />
      <FinanceTable />
    </div>
  )
}
