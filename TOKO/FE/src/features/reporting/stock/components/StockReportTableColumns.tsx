import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { formatStockNumber } from '@/features/products/products'

import type { StockReport, StockStatus } from '../stock.types'

// Tampilan badge per status: Habis (merah), Menipis (kuning), Aman & diragukan (tanpa badge).
const STATUS_BADGE: Record<StockStatus, { label: string; className: string; textClass: string } | null> = {
  out: {
    label: 'Habis',
    className: 'bg-red-100 text-red-700',
    textClass: 'text-red-600',
  },
  low: {
    label: 'Menipis',
    className: 'bg-amber-100 text-amber-700',
    textClass: 'text-amber-600',
  },
  ok: null,
  '': null,
}

export function buildStockReportColumns(): ColumnDef<StockReport>[] {
  return [
    {
      key: 'product_code',
      header: 'Kode',
      sortable: true,
      mobileHidden: true,
      cell: (r) => <span className="text-xs font-mono text-gray-500">{r.product_code}</span>,
    },
    {
      key: 'product_name',
      header: 'Nama Produk',
      sortable: true,
      mobileLabel: true,
      cell: (r) => <span className="text-sm font-medium">{r.product_name}</span>,
    },
    {
      key: 'category_name',
      header: 'Kategori',
      sortable: true,
      cell: (r) => <span className="text-sm text-gray-500">{r.category_name}</span>,
    },
    {
      key: 'unit',
      header: 'Satuan',
      mobileHidden: true,
      cell: (r) => <span className="text-sm">{r.unit}</span>,
    },
    {
      key: 'current_stock',
      header: 'Stok Saat Ini',
      align: 'right',
      sortable: true,
      cell: (r) => {
        const badge = STATUS_BADGE[r.stock_status]
        return (
          <div className="flex items-center justify-end gap-2">
            <span className={`text-sm font-semibold ${badge ? badge.textClass : ''}`}>
              {formatStockNumber(r.current_stock)}
            </span>
            {badge && (
              <span
                className={`inline-flex rounded-full px-1.5 py-0.5 text-xs font-medium ${badge.className}`}
              >
                {badge.label}
              </span>
            )}
          </div>
        )
      },
    },
    {
      key: 'stock_value',
      header: 'Nilai Stok',
      align: 'right',
      sortable: true,
      cell: (r) => <span className="text-sm font-medium">{formatRupiah(r.stock_value)}</span>,
    },
  ]
}
