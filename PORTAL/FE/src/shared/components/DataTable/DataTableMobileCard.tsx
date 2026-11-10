import { Checkbox } from '@/shared/components/ui/checkbox'

import type { ColumnDef } from './DataTable.types'

interface DataTableMobileCardProps {
  row: Record<string, unknown>
  columns: ColumnDef<Record<string, unknown>>[]
  rowSelection?: {
    isSelected: boolean
    onToggle: () => void
  }
}

export function DataTableMobileCard({ row, columns, rowSelection }: DataTableMobileCardProps) {
  const titleCol = columns.find((c) => c.mobileLabel) ?? columns[0]
  const actionsCol = columns.find((c) => c.key === 'actions')
  const bodyCols = columns.filter(
    (c) => c !== titleCol && c !== actionsCol && !c.mobileHidden
  )

  return (
    <div className="rounded-md border bg-white p-3">
      <div className="flex items-start justify-between gap-2">
        {rowSelection && (
          <Checkbox checked={rowSelection.isSelected} onCheckedChange={rowSelection.onToggle} />
        )}
        <div className="flex-1 min-w-0">
          {titleCol.cell ? titleCol.cell(row) : String(row[titleCol.key] ?? '')}
        </div>
      </div>

      {bodyCols.length > 0 && (
        <dl className="mt-2 flex flex-col gap-1.5 border-t pt-2">
          {bodyCols.map((col) => (
            <div key={col.key} className="flex items-center justify-between gap-3 text-sm">
              <dt className="text-gray-500 shrink-0">{col.header}</dt>
              <dd className="text-right min-w-0">
                {col.cell ? col.cell(row) : String(row[col.key] ?? '')}
              </dd>
            </div>
          ))}
        </dl>
      )}

      {actionsCol && (
        <div className="mt-2 flex flex-wrap items-center justify-end gap-1 border-t pt-2">
          {actionsCol.cell ? actionsCol.cell(row) : null}
        </div>
      )}
    </div>
  )
}
