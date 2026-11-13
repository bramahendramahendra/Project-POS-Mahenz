import { PageHeader } from '@/shared/components'

import { ReceivableTable } from './components/ReceivableTable'

export function ReceivablesPage() {
  return (
    <div className="space-y-4">
      <PageHeader title="Piutang" breadcrumbs={[{ label: 'Finance' }, { label: 'Piutang' }]} />
      <ReceivableTable />
    </div>
  )
}
