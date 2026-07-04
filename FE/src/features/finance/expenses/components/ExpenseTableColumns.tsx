import { Pencil, Trash2 } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import type { Expense } from '../expenses.types'

export interface ExpenseColumnHandlers {
  onEdit: (expense: Expense) => void
  onDelete: (expense: Expense) => void
}

const CATEGORY_LABEL: Record<string, string> = {
  operasional: 'Operasional',
  pembelian: 'Pembelian',
  gaji: 'Gaji',
  lainnya: 'Lainnya',
}

const PAYMENT_METHOD_LABEL: Record<string, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  card: 'Kartu',
  qris: 'QRIS',
  kredit: 'Kredit',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

export function buildExpenseColumns(handlers: ExpenseColumnHandlers): ColumnDef<Expense>[] {
  const { onEdit, onDelete } = handlers

  return [
    {
      key: 'expense_date',
      header: 'Tanggal',
      sortable: true,
      cell: (row) => <span className="text-gray-600 text-sm">{formatDate(row.expense_date)}</span>,
    },
    {
      key: 'category',
      header: 'Kategori',
      cell: (row) => (
        <span className="inline-flex items-center rounded-full bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700">
          {CATEGORY_LABEL[row.category] ?? row.category}
        </span>
      ),
    },
    {
      key: 'description',
      header: 'Keterangan',
      cell: (row) => <span className="text-sm text-gray-700">{row.description}</span>,
    },
    {
      key: 'payment_method',
      header: 'Metode',
      cell: (row) => (
        <span className="text-sm text-gray-600">
          {PAYMENT_METHOD_LABEL[row.payment_method] ?? row.payment_method}
        </span>
      ),
    },
    {
      key: 'amount',
      header: 'Jumlah',
      align: 'right',
      sortable: true,
      cell: (row) => (
        <span className="font-semibold text-red-600 text-sm">{formatRupiah(row.amount)}</span>
      ),
    },
    {
      key: 'user_name',
      header: 'Kasir',
      sortable: true,
      cell: (row) => <span className="text-sm text-gray-500">{row.user_name}</span>,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '100px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey="keuangan.pengeluaran" action="can_edit">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-blue-600" onClick={() => onEdit(row)}>
                  <Pencil size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Edit</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-7 w-7 text-gray-500 hover:text-red-600" onClick={() => onDelete(row)}>
                  <Trash2 size={14} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Hapus</TooltipContent>
            </Tooltip>
          </RoleGuard>
        </div>
      ),
    },
  ]
}
