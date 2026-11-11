import { useEffect, useRef, useState } from 'react'
import { ArrowDown, ArrowUp, ArrowUpDown } from 'lucide-react'

import { useBreakpoint } from '@/shared/hooks'
import { Checkbox } from '@/shared/components/ui/checkbox'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table'
import { cn } from '@/shared/utils'

import type { ColumnDef, DataTableProps, SortState } from './DataTable.types'
import { DataTableEmpty } from './DataTableEmpty'
import { DataTableMobileCard } from './DataTableMobileCard'
import { DataTablePagination } from './DataTablePagination'
import { DataTableSkeleton } from './DataTableSkeleton'

/** Tombol ikon aksi diperbesar di perangkat sentuh (tampilan laptop tidak berubah). */
const ACTIONS_TOUCH = 'touch:[&_button]:min-h-10 touch:[&_button]:min-w-10'

export function DataTable<TData extends Record<string, unknown>>({
  columns,
  data,
  isLoading,
  emptyMessage,
  emptyDescription,
  pagination,
  rowSelection,
  currentSort,
  onSort,
  className,
}: DataTableProps<TData>) {
  const isDesktopTable = useBreakpoint('md')

  // ─── Kolom aksi dibekukan di kanan ──────────────────────────────────────
  // Tabel lebar (tablet / laptop kecil) harus digeser; kolom 'actions' tetap terlihat di kanan.
  // Bayangan di sisi kiri kolom aksi hanya tampil bila masih ada kolom yang tertutup.
  const tableRef = useRef<HTMLTableElement>(null)
  const [hasHiddenRight, setHasHiddenRight] = useState(false)
  useEffect(() => {
    const scroller = tableRef.current?.parentElement
    if (!scroller) return
    const update = () =>
      setHasHiddenRight(scroller.scrollWidth - scroller.clientWidth - scroller.scrollLeft > 1)
    update()
    scroller.addEventListener('scroll', update, { passive: true })
    const observer = new ResizeObserver(update)
    observer.observe(scroller)
    return () => {
      scroller.removeEventListener('scroll', update)
      observer.disconnect()
    }
  }, [isDesktopTable, data, columns])

  const stickyActions = (isHeader: boolean) =>
    cn(
      'sticky right-0 bg-white',
      isHeader ? 'z-[2]' : 'z-[1]',
      // Latar hover/terpilih baris ikut tampil di kolom beku (latar sel harus padat agar isi di bawahnya tertutup)
      !isHeader &&
        "before:pointer-events-none before:absolute before:inset-0 before:content-[''] group-hover:before:bg-muted/50 group-data-[state=selected]:before:bg-muted",
      hasHiddenRight && 'shadow-[-8px_0_8px_-8px_rgba(0,0,0,0.25)]'
    )

  // ─── Row selection helpers ──────────────────────────────────────────────
  const allRowKeys = data.map((row) => row[rowSelection?.rowKey as keyof TData] as string | number)
  const selectedKeys = rowSelection?.selectedKeys ?? new Set<string | number>()
  const allSelected = allRowKeys.length > 0 && allRowKeys.every((k) => selectedKeys.has(k))
  const someSelected = allRowKeys.some((k) => selectedKeys.has(k)) && !allSelected

  const handleSelectAll = () => {
    if (!rowSelection) return
    if (allSelected) {
      const next = new Set(selectedKeys)
      allRowKeys.forEach((k) => next.delete(k))
      rowSelection.onSelectionChange(next)
    } else {
      const next = new Set(selectedKeys)
      allRowKeys.forEach((k) => next.add(k))
      rowSelection.onSelectionChange(next)
    }
  }

  const handleSelectRow = (key: string | number) => {
    if (!rowSelection) return
    const next = new Set(selectedKeys)
    if (next.has(key)) next.delete(key)
    else next.add(key)
    rowSelection.onSelectionChange(next)
  }

  // ─── Sort handler ───────────────────────────────────────────────────────
  const handleSort = (key: string) => {
    if (!onSort) return
    const order: SortState['order'] =
      currentSort?.key === key && currentSort.order === 'asc' ? 'desc' : 'asc'
    onSort({ key, order })
  }

  const getSortIcon = (key: string) => {
    if (!currentSort || currentSort.key !== key) return <ArrowUpDown size={12} className="text-gray-400" />
    return currentSort.order === 'asc'
      ? <ArrowUp size={12} className="text-blue-500" />
      : <ArrowDown size={12} className="text-blue-500" />
  }

  // ─── Loading ────────────────────────────────────────────────────────────
  if (isLoading) {
    const colCount = columns.length + (rowSelection?.enabled ? 1 : 0)
    return <DataTableSkeleton columns={colCount} />
  }

  // ─── Empty ──────────────────────────────────────────────────────────────
  if (data.length === 0) {
    return (
      <div className={cn('rounded-md border bg-white', className)}>
        <DataTableEmpty message={emptyMessage} description={emptyDescription} />
      </div>
    )
  }

  // ─── Mobile card view ───────────────────────────────────────────────────
  if (!isDesktopTable) {
    return (
      <div className={cn('flex flex-col gap-3', className)}>
        {data.map((row, rowIdx) => {
          const rowKey = rowSelection
            ? (row[rowSelection.rowKey as keyof TData] as string | number)
            : rowIdx
          const isSelected = rowSelection?.selectedKeys?.has(rowKey) ?? false

          return (
            <DataTableMobileCard
              key={String(rowKey)}
              row={row}
              columns={columns as ColumnDef<Record<string, unknown>>[]}
              rowSelection={
                rowSelection?.enabled
                  ? { isSelected, onToggle: () => handleSelectRow(rowKey) }
                  : undefined
              }
            />
          )
        })}
        {pagination && <DataTablePagination {...pagination} />}
      </div>
    )
  }

  return (
    <div className={cn('rounded-md border bg-white', className)}>
      <Table ref={tableRef}>
        <TableHeader className="sticky top-0 bg-white z-10">
          <TableRow>
            {rowSelection?.enabled && (
              <TableHead className="w-10">
                <Checkbox
                  checked={allSelected ? true : someSelected ? 'indeterminate' : false}
                  onCheckedChange={handleSelectAll}
                />
              </TableHead>
            )}
            {columns.map((col) => (
              <TableHead
                key={col.key}
                style={{ width: col.width, textAlign: col.align ?? 'left' }}
                className={cn(
                  col.sortable && 'cursor-pointer select-none',
                  col.key === 'actions' && stickyActions(true)
                )}
                onClick={col.sortable ? () => handleSort(col.key) : undefined}
              >
                <span className="inline-flex items-center gap-1">
                  {col.header}
                  {col.sortable && getSortIcon(col.key)}
                </span>
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.map((row, rowIdx) => {
            const rowKey = rowSelection
              ? (row[rowSelection.rowKey as keyof TData] as string | number)
              : rowIdx
            const isSelected = rowSelection?.selectedKeys?.has(rowKey) ?? false

            return (
              <TableRow key={String(rowKey)} data-state={isSelected ? 'selected' : undefined} className="group">
                {rowSelection?.enabled && (
                  <TableCell className="w-10">
                    <Checkbox
                      checked={isSelected}
                      onCheckedChange={() => handleSelectRow(rowKey)}
                    />
                  </TableCell>
                )}
                {columns.map((col) => (
                  <TableCell
                    key={col.key}
                    style={{ textAlign: col.align ?? 'left' }}
                    className={cn(col.key === 'actions' && [stickyActions(false), ACTIONS_TOUCH])}
                  >
                    {col.key === 'actions' ? (
                      // Isi di atas lapisan hover kolom beku
                      <div className="relative">{col.cell ? col.cell(row) : null}</div>
                    ) : col.cell ? (
                      col.cell(row)
                    ) : (
                      String(row[col.key as keyof TData] ?? '')
                    )}
                  </TableCell>
                ))}
              </TableRow>
            )
          })}
        </TableBody>
      </Table>

      {pagination && <DataTablePagination {...pagination} />}
    </div>
  )
}
