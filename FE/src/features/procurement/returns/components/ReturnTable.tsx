import { forwardRef, useImperativeHandle, useState } from 'react'
import { Truck } from 'lucide-react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { monthStart, todayStr } from '@/shared/utils'
import { useSupplierListQuery } from '@/features/procurement/suppliers'
import { useProductListQuery } from '@/features/products/products'

import {
  useSupplierReturnsQuery,
  useDeleteSupplierReturnMutation,
} from '../returns.api'
import type { SupplierReturn, SupplierReturnFilter } from '../returns.types'
import { buildReturnColumns } from './ReturnTableColumns'
import { ReturnFilterBar } from './ReturnFilterBar'
import { ReturnFormModal } from './ReturnFormModal'
import { ReturnDetailModal } from './ReturnDetailModal'

export interface ReturnTableHandle {
  openAdd: () => void
}

export const ReturnTable = forwardRef<ReturnTableHandle, object>(function ReturnTable(_, ref) {
  const [filter, setFilter] = useState<SupplierReturnFilter>({
    start_date: monthStart(),
    end_date: todayStr(),
  })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()

  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [detailId, setDetailId] = useState<number | null>(null)

  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 200, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })
  const suppliers = suppliersData?.data ?? []
  const hasSuppliers = suppliers.length > 0
  const hasProducts = (productsData?.total ?? 0) > 0

  const { data, isLoading } = useSupplierReturnsQuery({ ...filter, page, limit: pageSize })
  const { mutate: deleteReturn, isPending: isDeleting } = useDeleteSupplierReturnMutation()

  const returns = data?.data ?? []
  const total = data?.total ?? 0

  useImperativeHandle(ref, () => ({ openAdd: openForm }))

  const handleFilterChange = (newFilter: SupplierReturnFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ start_date: monthStart(), end_date: todayStr() })
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const handleDetail = (row: SupplierReturn) => {
    setDetailId(row.id)
    openDetail()
  }

  const handleDelete = (id: number) => {
    setDeletingId(id)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deletingId) return
    deleteReturn(deletingId, {
      onSuccess: () => {
        closeDelete()
        setDeletingId(null)
      },
    })
  }

  if (isSuppliersLoading || isProductsLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  if (!hasSuppliers || !hasProducts) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4 flex gap-3">
          <div className="rounded-full p-3 bg-amber-50">
            <Truck size={24} className="text-amber-500" />
          </div>
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa menambah retur</h3>
        <p className="mb-1 text-sm text-gray-500">
          Sebelum menambah retur, pastikan data berikut sudah tersedia:
        </p>
        <ul className="mb-6 space-y-1 text-sm">
          {!hasSuppliers && (
            <li className="flex items-center gap-2 text-amber-600">
              <span>!</span>
              Belum ada supplier — tambahkan di menu Supplier
            </li>
          )}
          {!hasProducts && (
            <li className="flex items-center gap-2 text-amber-600">
              <span>!</span>
              Belum ada produk — tambahkan di menu Produk
            </li>
          )}
        </ul>
      </div>
    )
  }

  const columns = buildReturnColumns({ onDetail: handleDetail, onDelete: handleDelete })

  return (
    <>
      <ReturnFilterBar
        filter={filter}
        suppliers={suppliers}
        onChange={handleFilterChange}
        onReset={handleReset}
      />

      <DataTable<SupplierReturn & Record<string, unknown>>
        columns={columns}
        data={returns as (SupplierReturn & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data retur"
        emptyDescription="Data retur pembelian akan muncul sesuai filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <ReturnFormModal open={formOpen} onOpenChange={(o) => !o && closeForm()} />

      <ReturnDetailModal
        returnId={detailId}
        open={detailOpen}
        onOpenChange={(o) => {
          if (!o) {
            closeDetail()
            setDetailId(null)
          }
        }}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(o) => {
          if (!o) {
            closeDelete()
            setDeletingId(null)
          }
        }}
        title="Hapus Retur"
        description="Data retur yang dihapus tidak bisa dikembalikan. Yakin ingin melanjutkan?"
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </>
  )
})
