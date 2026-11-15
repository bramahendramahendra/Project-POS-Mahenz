import { PageHeader, PrerequisiteGuard } from '@/shared/components'

import { StockReportTab } from './components/StockReportTab'
import { useStockReportPrerequisites } from './hooks/useStockReportPrerequisites'

export function StockReportPage() {
  const { isLoading, items } = useStockReportPrerequisites()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Laporan Stok"
        breadcrumbs={[{ label: 'Pelaporan' }, { label: 'Stok' }]}
      />
      <PrerequisiteGuard
        isLoading={isLoading}
        title="Laporan stok belum tersedia"
        description="Laporan stok membutuhkan data berikut:"
        items={items}
      >
        <StockReportTab />
      </PrerequisiteGuard>
    </div>
  )
}
