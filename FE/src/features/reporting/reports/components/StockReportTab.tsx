import { useState } from 'react'

import { DataTable } from '@/shared/components'
import { formatRupiah } from '@/shared/utils'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { useStockReportQuery } from '../reports.api'
import type { StockReport, StockReportFilter } from '../reports.types'
import { StockReportFilterBar } from './StockReportFilterBar'
import { buildStockReportColumns } from './StockReportTableColumns'

export function StockReportTab() {
  const [filter, setFilter] = useState<StockReportFilter>({})
  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useStockReportQuery({ ...filter, page, page_size: pageSize })
  const items: StockReport[] = data?.items ?? []
  const total = data?.total ?? 0
  const totalStockValue = data?.total_stock_value ?? 0

  const columns = buildStockReportColumns()

  return (
    <div className="space-y-4">
      <StockReportFilterBar filter={filter} onChange={setFilter} onReset={reset} />

      <div className="grid grid-cols-2 gap-3 max-w-sm">
        <div className="rounded-lg border bg-white p-4">
          <p className="text-xs text-gray-500">Total Item Produk</p>
          <p className="text-xl font-bold text-gray-800">{total}</p>
        </div>
        <div className="rounded-lg border bg-white p-4">
          <p className="text-xs text-gray-500">Total Nilai Stok</p>
          <p className="text-xl font-bold text-gray-800">{formatRupiah(totalStockValue)}</p>
        </div>
      </div>

      <DataTable<StockReport & Record<string, unknown>>
        columns={columns}
        data={items as (StockReport & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data stok"
        emptyDescription="Data stok produk akan muncul sesuai filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
    </div>
  )
}
