import { useState } from 'react'

import { PageHeader } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useReceivableListQuery } from './receivables.api'
import type { Receivable, ReceivableListFilter } from './receivables.types'
import { PaymentRecordModal } from './components/PaymentRecordModal'
import { ReceivableFilterBar } from './components/ReceivableFilterBar'
import { ReceivableTable } from './components/ReceivableTable'

export function ReceivablesPage() {
  const [filter, setFilter] = useState<ReceivableListFilter>({ page: 1, limit: 10 })
  const [payTarget, setPayTarget] = useState<Receivable | null>(null)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useReceivableListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const receivables = data?.data ?? []
  const total = data?.total ?? 0

  return (
    <div className="space-y-4">
      <PageHeader title="Piutang" breadcrumbs={[{ label: 'Finance' }, { label: 'Piutang' }]} />

      <ReceivableFilterBar filter={filter} onChange={setFilter} onReset={reset} />

      <ReceivableTable
        data={receivables}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        onPay={(r) => setPayTarget(r)}
      />

      <PaymentRecordModal
        open={!!payTarget}
        onOpenChange={(open) => { if (!open) setPayTarget(null) }}
        receivable={payTarget}
      />
    </div>
  )
}
