import { PageHeader } from '@/shared/components'

import { CashDrawerTable } from './components/CashDrawerTable'
import { ShiftPrerequisiteGuard } from './components/ShiftPrerequisiteGuard'

export function CashDrawerPage() {
  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Harian"
        breadcrumbs={[{ label: 'Keuangan' }, { label: 'Kas Harian' }]}
      />

      <ShiftPrerequisiteGuard>
        <CashDrawerTable />
      </ShiftPrerequisiteGuard>
    </div>
  )
}
