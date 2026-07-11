import { PageHeader } from '@/shared/components'

import { BackupTab } from './components/BackupTab'

export function BackupPage() {
  return (
    <div className="space-y-4">
      <PageHeader
        title="Backup & Restore"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Backup & Restore' }]}
      />
      <BackupTab />
    </div>
  )
}
