import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { StockReport } from '../stock.types'

export function buildStockReportColumns(): ColumnDef<StockReport>[] {
  return [
    {
      key: 'product_code',
      header: 'Kode',
      cell: (r) => <span className="text-xs font-mono text-gray-500">{r.product_code}</span>,
    },
    {
      key: 'product_name',
      header: 'Nama Produk',
      cell: (r) => <span className="text-sm font-medium">{r.product_name}</span>,
    },
    {
      key: 'category_name',
      header: 'Kategori',
      cell: (r) => <span className="text-sm text-gray-500">{r.category_name}</span>,
    },
    {
      key: 'unit',
      header: 'Satuan',
      cell: (r) => <span className="text-sm">{r.unit}</span>,
    },
    {
      key: 'current_stock',
      header: 'Stok Saat Ini',
      align: 'right',
      cell: (r) => (
        <div className="flex items-center justify-end gap-2">
          <span className={`text-sm font-semibold ${r.current_stock < r.min_stock ? 'text-red-600' : ''}`}>
            {r.current_stock}
          </span>
          {r.current_stock < r.min_stock && (
            <span className="inline-flex rounded-full bg-red-100 px-1.5 py-0.5 text-xs font-medium text-red-700">
              Stok Rendah
            </span>
          )}
        </div>
      ),
    },
    {
      key: 'stock_value',
      header: 'Nilai Stok',
      align: 'right',
      cell: (r) => <span className="text-sm font-medium">{formatRupiah(r.stock_value)}</span>,
    },
  ]
}
