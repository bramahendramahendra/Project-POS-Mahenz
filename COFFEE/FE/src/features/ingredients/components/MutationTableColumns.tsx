import { Undo2 } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { formatDateTime } from '@/shared/utils'

import type { IngredientMutation } from '../ingredients.types'
import { MUTATION_TYPE_CLASS, MUTATION_TYPE_LABEL, formatQty, formatQtyChange } from '../ingredients.utils'

export function buildMutationColumns(onCancel: (row: IngredientMutation) => void): ColumnDef<IngredientMutation>[] {
  return [
    {
      key: 'created_at',
      header: 'Waktu',
      width: '150px',
      cell: (row) => <span className="text-sm text-gray-700">{formatDateTime(row.created_at)}</span>,
    },
    {
      key: 'ingredient_name',
      header: 'Bahan',
      mobileLabel: true,
      cell: (row) => <span className="font-medium text-gray-800">{row.ingredient_name}</span>,
    },
    {
      key: 'type',
      header: 'Jenis',
      width: '110px',
      cell: (row) => (
        <span className={`inline-flex rounded-full px-2 py-0.5 text-[11px] font-medium ${MUTATION_TYPE_CLASS[row.type]}`}>
          {MUTATION_TYPE_LABEL[row.type]}
        </span>
      ),
    },
    {
      key: 'qty_change',
      header: 'Jumlah',
      align: 'right',
      cell: (row) => (
        <span className={`font-semibold ${row.is_cancelled ? 'text-gray-400 line-through' : row.qty_change < 0 ? 'text-red-600' : 'text-green-700'}`}>
          {formatQtyChange(row.qty_change, row.unit_abbr)}
        </span>
      ),
    },
    {
      key: 'stock_after',
      header: 'Stok Sesudah',
      width: '120px',
      align: 'right',
      mobileHidden: true,
      cell: (row) => <span className="text-gray-700">{formatQty(row.stock_after, row.unit_abbr)}</span>,
    },
    {
      key: 'notes',
      header: 'Catatan / Alasan',
      cell: (row) => (
        <div className="max-w-[200px] space-y-0.5 text-sm">
          {row.notes ? <p className="text-gray-700 break-words">{row.notes}</p> : <p className="text-gray-400">—</p>}
          {row.is_cancelled && (
            <p className="text-xs text-red-600">
              Dibatalkan oleh {row.cancelled_by_name} ({formatDateTime(row.cancelled_at ?? '')}): {row.cancel_reason}
            </p>
          )}
        </div>
      ),
    },
    {
      key: 'created_by_name',
      header: 'Oleh',
      cell: (row) => <span className="text-sm text-gray-700 break-words">{row.created_by_name || '—'}</span>,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '64px',
      cell: (row) =>
        (row.type === 'usage' || row.type === 'waste') && !row.is_cancelled ? (
          <RoleGuard menuKey="bahan.riwayat" action="can_delete">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Batalkan entri"
                  className="h-7 w-7 text-gray-500 hover:text-red-600"
                  onClick={() => onCancel(row)}
                >
                  <Undo2 size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Batalkan entri</TooltipContent>
            </Tooltip>
          </RoleGuard>
        ) : null,
    },
  ]
}
