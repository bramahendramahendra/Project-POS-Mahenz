import { DataTable } from '@/shared/components'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { TopProductItem } from '../dashboard.types'

interface TopProductsTableProps {
  data: TopProductItem[]
  isLoading: boolean
}

export function TopProductsTable({ data, isLoading }: TopProductsTableProps) {
  const rows = data.map((item, i) => ({ ...item, rank: i + 1 }))

  const columns: ColumnDef<TopProductItem & { rank: number }>[] = [
    {
      key: 'rank',
      header: '#',
      align: 'center',
      width: '32px',
      cell: (row) => <span className="text-gray-400 font-mono">{row.rank}</span>,
    },
    {
      key: 'product_name',
      header: 'Produk',
      cell: (row) => <span className="font-medium text-gray-800">{row.product_name}</span>,
    },
    {
      key: 'total_qty',
      header: 'Qty',
      align: 'right',
      cell: (row) => <span className="text-gray-600">{row.total_qty}</span>,
    },
    {
      key: 'total_value',
      header: 'Revenue',
      align: 'right',
      cell: (row) => <span className="font-semibold text-blue-600">{formatRupiah(row.total_value)}</span>,
    },
  ]

  return (
    <div className="rounded-lg border bg-white overflow-hidden">
      <div className="px-4 py-3 border-b">
        <h3 className="font-semibold text-gray-700 text-sm">Top Produk Terlaris</h3>
      </div>
      <DataTable<(TopProductItem & { rank: number }) & Record<string, unknown>>
        columns={columns}
        data={rows as ((TopProductItem & { rank: number }) & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada data"
      />
    </div>
  )
}
