import type { ReactNode } from 'react'

import { cn } from '@/shared/utils'

interface ScrollableTableProps {
  children: ReactNode
  minWidth?: number
  className?: string
  tableClassName?: string
}

export function ScrollableTable({
  children,
  minWidth = 560,
  className,
  tableClassName,
}: ScrollableTableProps) {
  return (
    <div className="space-y-1.5">
      <p className="text-[11px] text-gray-400 sm:hidden">
        ← Geser tabel ke samping untuk lihat semua kolom →
      </p>
      <div className={cn('rounded-lg border bg-white overflow-x-auto', className)}>
        <table className={cn('w-full text-sm', tableClassName)} style={{ minWidth }}>
          {children}
        </table>
      </div>
    </div>
  )
}
