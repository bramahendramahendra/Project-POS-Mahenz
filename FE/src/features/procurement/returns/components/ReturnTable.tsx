import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { monthStart, todayStr } from '@/shared/utils'
import { useSupplierListQuery } from '@/features/procurement/suppliers'

import { useSupplierReturnsQuery, useDeleteSupplierReturnMutation } from '../returns.api'
import type { SupplierReturn, SupplierReturnFilter } from '../returns.types'
import { buildReturnColumns } from './ReturnTableColumns'
import { ReturnFilterBar } from './ReturnFilterBar'
import { ReturnFormModal } from './ReturnFormModal'
import { ReturnDetailModal } from './ReturnDetailModal'

export interface ReturnTableHandle {
  openAdd: () => void
}

export const ReturnTable = forwardRef<ReturnTableHandle, object>(function ReturnTable(_, ref) {
  const [filter, setFilter] = useState<SupplierReturnFilter>({ page: 1, limit: 10, start_date: monthStart(), end_date: todayStr() })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()

  const [deletingReturn, setDeletingReturn] = useState<SupplierReturn | null>(null)
  const [detailReturn, setDetailReturn] = useState<SupplierReturn | null>(null)

  const { data: suppliersData } = useSupplierListQuery({ page: 1, limit: 200, search: '' })
  const suppliers = suppliersData?.data ?? []

  const { data, isLoading } = useSupplierReturnsQuery({ ...filter, page, limit: pageSize })
  const returns = data?.data ?? []
  const total = data?.total ?? 0
  
  const { mutate: deleteReturn, isPending: isDeleting } = useDeleteSupplierReturnMutation()

  useImperativeHandle(ref, () => ({ openAdd: openForm }))

  const handleFilterChange = (newFilter: SupplierReturnFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, start_date: monthStart(), end_date: todayStr() })
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleOpenDetail = (supplierReturn: SupplierReturn) => {
    setDetailReturn(supplierReturn)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailReturn(null)
  }

  const handleOpenDelete = (supplierReturn: SupplierReturn) => {
    setDeletingReturn(supplierReturn)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingReturn(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingReturn) return
    deleteReturn(deletingReturn.id, {
      onSuccess: () => handleCloseDelete(),
    })
  }

  const hasFilter = !!filter.supplier_id || !!filter.status

  const columns = buildReturnColumns({
    onDetail: handleOpenDetail,
    onDelete: handleOpenDelete,
  })

  return (
    <div className="space-y-4">
      <ReturnFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        suppliers={suppliers}
      />

      <DataTable<SupplierReturn & Record<string, unknown>>
        columns={columns}
        data={returns as (SupplierReturn & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Data retur tidak ditemukan' : 'Belum ada data retur'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah filter pencarian Anda.'
            : 'Data retur pembelian akan muncul di sini.'
        }
        pagination={{ 
          page, 
          pageSize, 
          total, 
          onPageChange, 
          onPageSizeChange, 
          pageSizeOptions 
        }}
      />

      <ReturnFormModal 
        open={formOpen} 
        onOpenChange={(o) => !o && closeForm()} 
      />

      <ReturnDetailModal
        returnId={detailReturn?.id ?? null}
        open={detailOpen}
        onOpenChange={(open) => { if (!open) handleCloseDetail() }}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Retur"
        description={`Yakin ingin menghapus retur "${deletingReturn?.return_code}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
