import { PageHeader, PrerequisiteGuard } from '@/shared/components'

import { CashDrawerTable } from './components/CashDrawerTable'
import { useCashDrawerPrerequisites } from './hooks/useCashDrawerPrerequisites'

export function CashDrawerPage() {
  const { isLoading, items } = useCashDrawerPrerequisites()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Harian"
        breadcrumbs={[{ label: 'Keuangan' }, { label: 'Kas Harian' }]}
      />

      <PrerequisiteGuard
        isLoading={isLoading}
        title="Belum bisa membuka Kas Harian"
        description="Sebelum menggunakan Kas Harian, pastikan data berikut sudah tersedia:"
        items={items}
      >
        <CashDrawerTable />
      </PrerequisiteGuard>
    </div>
  )
}
