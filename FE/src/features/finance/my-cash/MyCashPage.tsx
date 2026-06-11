import { PageHeader, DataTable } from '@/shared/components'
import { formatRupiah } from '@/shared/utils'

import { useMyCashQuery } from './my-cash.api'
import type { MyCashTransaction } from './my-cash.types'
import { buildMyCashColumns } from './components/MyCashTableColumns'

export function MyCashPage() {
  const { data, isLoading } = useMyCashQuery()
  const myCash = data
  const transactions: MyCashTransaction[] = myCash?.transactions ?? []
  const columns = buildMyCashColumns()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Saya"
        breadcrumbs={[{ label: 'Finance' }, { label: 'Kas Saya' }]}
      />

      <div className="rounded-xl border bg-white p-6 shadow-sm">
        <p className="text-sm text-gray-500 mb-1">Saldo Kas Saat Ini</p>
        {isLoading ? (
          <div className="h-10 w-40 animate-pulse rounded bg-gray-100" />
        ) : (
          <p className="text-4xl font-bold text-gray-900">
            {formatRupiah(myCash?.balance ?? 0)}
          </p>
        )}
      </div>

      <div className="space-y-3">
        <h2 className="font-semibold text-gray-700">Riwayat Kas</h2>
        <DataTable<MyCashTransaction & Record<string, unknown>>
          columns={columns}
          data={transactions as (MyCashTransaction & Record<string, unknown>)[]}
          isLoading={isLoading}
          emptyMessage="Belum ada riwayat kas"
        />
      </div>
    </div>
  )
}
