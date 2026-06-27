import { forwardRef, useImperativeHandle, useState } from 'react'
import { Truck } from 'lucide-react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { useSupplierListQuery } from '@/features/procurement/suppliers'
import { useProductListQuery } from '@/features/products/products'

import { monthStart, todayStr } from '@/shared/utils'

import { useSupplierPurchasesQuery, useDeleteSupplierPurchaseMutation } from '../purchases.api'
import type { SupplierPurchase, SupplierPurchaseFilter } from '../purchases.types'
import { PurchaseFilterBar } from './PurchaseFilterBar'
import { PurchaseFormModal } from './PurchaseFormModal'
import { PurchaseDetailModal } from './PurchaseDetailModal'
import { PaymentModal } from './PaymentModal'
import { buildPurchaseColumns } from './PurchaseTableColumns'

export interface PurchaseTableHandle {
  openAdd: () => void
}

export const PurchaseTable = forwardRef<PurchaseTableHandle, object>(function PurchaseTable(_, ref) {
  const [filter, setFilter] = useState<SupplierPurchaseFilter>({
    page: 1,
    limit: 10,
    start_date: monthStart(),
    end_date: todayStr(),
  })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: payOpen, open: openPay, close: closePay } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingPurchase, setEditingPurchase] = useState<SupplierPurchase | null>(null)
  const [payingPurchase, setPayingPurchase] = useState<SupplierPurchase | null>(null)
  const [detailId, setDetailId] = useState<number | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 200, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })
  const suppliers = suppliersData?.data ?? []
  const hasSuppliers = suppliers.length > 0
  const hasProducts = (productsData?.total ?? 0) > 0

  const { data, isLoading } = useSupplierPurchasesQuery({ 
    ...filter, 
    page, 
    limit: pageSize, 
  })
  const purchases = data?.data ?? []
  const total = data?.total ?? 0

  const { mutate: deletePurchase   , isPending: isDeleting } = useDeleteSupplierPurchaseMutation()

  const handleOpenAdd = () => {
    setEditingPurchase(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleEdit = (purchase: SupplierPurchase) => {
    setEditingPurchase(purchase)
    openForm()
  }

  const handlePay = (purchase: SupplierPurchase) => {
    setPayingPurchase(purchase)
    openPay()
  }

  const handleDelete = (purchase: SupplierPurchase) => {
    setDeletingId(purchase.id)
    openDelete()
  }

  const handleDetail = (purchase: SupplierPurchase) => {
    setDetailId(purchase.id)
    openDetail()
  }

  const handleFormClose = () => {
    closeForm()
    setEditingPurchase(null)
  }

  const handleFilterChange = (newFilter: SupplierPurchaseFilter) => {
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

  const handleConfirmDelete = () => {
    if (!deletingId) return
    deletePurchase(deletingId, {
      onSuccess: () => {
        closeDelete()
        setDeletingId(null)
      },
    })
  }

  const columns = buildPurchaseColumns({
    onDetail: handleDetail,
    onEdit: handleEdit,
    onPay: handlePay,
    onDelete: handleDelete,
  })

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
        <div className="mb-4">
          <div className="rounded-full p-3 bg-amber-50">
            <Truck size={24} className="text-amber-500" />
          </div>
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">
          Belum bisa menambah pembelian
        </h3>
        <p className="mb-1 text-sm text-gray-500">
          Sebelum menambah pembelian, pastikan data berikut sudah tersedia:
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

  return (
    <div className="space-y-4">
      <PurchaseFilterBar
        filter={filter}
        suppliers={suppliers}
        onChange={handleFilterChange}
        onReset={handleReset}
      />

      <DataTable<SupplierPurchase & Record<string, unknown>>
        columns={columns}
        data={purchases as (SupplierPurchase & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage="Belum ada data pembelian"
        emptyDescription="Data pembelian supplier akan muncul sesuai filter yang dipilih."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <PurchaseFormModal
        open={formOpen}
        onOpenChange={(o) => { if (!o) handleFormClose() }}
        initialData={editingPurchase}
      />

      <PurchaseDetailModal
        open={detailOpen}
        onOpenChange={(o) => {
          if (!o) {
            closeDetail()
            setDetailId(null)
          }
        }}
        purchaseId={detailId}
      />

      <PaymentModal
        open={payOpen}
        onOpenChange={(o) => {
          if (!o) {
            closePay()
            setPayingPurchase(null)
          }
        }}
        purchase={payingPurchase}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(o) => {
          if (!o) {
            closeDelete()
            setDeletingId(null)
          }
        }}
        title="Hapus Pembelian"
        description="Data pembelian yang dihapus tidak bisa dikembalikan. Yakin ingin melanjutkan?"
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
