import { DataTable } from '@/shared/components'
import type { PaginationProps } from '@/shared/components/DataTable/DataTable.types'

import type { Shift } from '../shifts.types'
import { buildShiftColumns } from './ShiftTableColumns'

interface ShiftTableProps {
  data: Shift[]
  isLoading: boolean
  pagination: PaginationProps
  onClose: (shift: Shift) => void
}

export function ShiftTable({ data, isLoading, pagination, onClose }: ShiftTableProps) {
  const columns = buildShiftColumns({ onClose })

  return (
    <DataTable<Shift & Record<string, unknown>>
      columns={columns}
      data={data as (Shift & Record<string, unknown>)[]}
      isLoading={isLoading}
      emptyMessage="Belum ada data shift"
      emptyDescription="Shift akan muncul setelah kasir membuka sesi kerja."
      pagination={pagination}
    />
  )
}
