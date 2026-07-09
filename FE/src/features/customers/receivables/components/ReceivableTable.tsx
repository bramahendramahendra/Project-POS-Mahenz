import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useReceivableListQuery } from '../receivables.api'
import type { Receivable, ReceivableListFilter } from '../receivables.types'
import { PaymentRecordModal } from './PaymentRecordModal'
import { ReceivableFilterBar } from './ReceivableFilterBar'
import { buildReceivableColumns } from './ReceivableTableColumns'

export function ReceivableTable() {
  const [filter, setFilter] = useState<ReceivableListFilter>({ page: 1, limit: 10, search: '' })
  const [payTarget, setPayTarget] = useState<Receivable | null>(null)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useReceivableListQuery({ ...filter, page, limit: pageSize })
  const receivables = data?.data ?? []
  const total = data?.total ?? 0

  const columns = buildReceivableColumns({ onPay: (r) => setPayTarget(r) })

  return (
    <div className="space-y-4">
      <ReceivableFilterBar filter={filter} onChange={setFilter} onReset={reset} />

      <DataTable<Receivable & Record<string, unknown>>
        columns={columns}
        data={receivables as (Receivable & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada piutang"
        emptyDescription="Piutang akan muncul saat transaksi dilakukan dengan metode kredit."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <PaymentRecordModal
        open={!!payTarget}
        onOpenChange={(open) => { if (!open) setPayTarget(null) }}
        receivable={payTarget}
      />
    </div>
  )
}
