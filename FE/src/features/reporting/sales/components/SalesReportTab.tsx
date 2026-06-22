import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { formatRupiah, monthStart, todayStr } from '@/shared/utils'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useSalesReportQuery } from '../sales.api'
import type { SalesReport, SalesReportFilter } from '../sales.types'
import { SalesReportFilterBar } from './SalesReportFilterBar'
import { buildSalesReportColumns } from './SalesReportTableColumns'

interface SummaryCardProps {
  label: string
  value: string
  isLoading: boolean
}

function SummaryCard({ label, value, isLoading }: SummaryCardProps) {
  return (
    <div className="rounded-lg border bg-white p-4 space-y-1">
      <p className="text-xs text-gray-500">{label}</p>
      {isLoading ? (
        <div className="h-7 w-28 animate-pulse rounded bg-gray-100" />
      ) : (
        <p className="text-xl font-bold text-gray-800">{value}</p>
      )}
    </div>
  )
}

export function SalesReportTab() {
  const [filter, setFilter] = useState<SalesReportFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useSalesReportQuery({ ...filter, page, page_size: pageSize })
  const items: SalesReport[] = data?.items ?? []
  const total = data?.total ?? 0
  const summary = data?.summary

  const columns = buildSalesReportColumns()

  return (
    <div className="space-y-4">
      <SalesReportFilterBar filter={filter} onChange={setFilter} onReset={reset} exportData={items} />

      <div className="grid grid-cols-3 gap-3">
        <SummaryCard
          label="Total Transaksi"
          value={String(summary?.total_transactions ?? 0)}
          isLoading={isLoading}
        />
        <SummaryCard
          label="Total Pendapatan"
          value={formatRupiah(summary?.total_revenue ?? 0)}
          isLoading={isLoading}
        />
        <SummaryCard
          label="Rata-rata per Transaksi"
          value={formatRupiah(summary?.avg_per_transaction ?? 0)}
          isLoading={isLoading}
        />
      </div>

      <DataTable<SalesReport & Record<string, unknown>>
        columns={columns}
        data={items as (SalesReport & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data penjualan"
        emptyDescription="Data penjualan akan muncul sesuai filter periode yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
    </div>
  )
}
