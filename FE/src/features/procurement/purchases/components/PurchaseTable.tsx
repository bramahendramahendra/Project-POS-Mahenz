import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { monthStart, todayStr } from '@/shared/utils'

import { useSupplierListQuery } from '@/features/procurement/suppliers'

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
  const [filter, setFilter] = useState<SupplierPurchaseFilter>({ page: 1, limit: 10, start_date: monthStart(), end_date: todayStr() })
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
  const [deletingPurchase, setDeletingPurchase] = useState<SupplierPurchase | null>(null)

  const { data: suppliersData } = useSupplierListQuery({ page: 1, limit: 200, search: '' })
  const suppliers = suppliersData?.data ?? []

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

  const handleOpenEdit = (purchase: SupplierPurchase) => {
    setEditingPurchase(purchase)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingPurchase(null)
  }

  const handleOpenDetail = (purchase: SupplierPurchase) => {
    setDetailId(purchase.id)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailId(null)
  }

  const handlePay = (purchase: SupplierPurchase) => {
    setPayingPurchase(purchase)
    openPay()
  }
  const handleOpenDelete = (purchase: SupplierPurchase) => {
    setDeletingPurchase(purchase)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deletingPurchase) return
    deletePurchase(deletingPurchase.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingPurchase(null)
      },
    })
  }

  const hasFilter = !!filter.supplier_id || !!filter.payment_status

  const columns = buildPurchaseColumns({
    onDetail: handleOpenDetail,
    onEdit: handleOpenEdit,
    onPay: handlePay,
    onDelete: handleOpenDelete,
  })

  return (
    <div className="space-y-4">
      <PurchaseFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        suppliers={suppliers}
      />

      <DataTable<SupplierPurchase & Record<string, unknown>>
        columns={columns}
        data={purchases as (SupplierPurchase & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Data pembelian tidak ditemukan' : 'Belum ada data pembelian'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah filter pencarian Anda.'
            : 'Data pembelian supplier akan muncul di sini.'
        }
        pagination={{ 
          page, 
          pageSize, 
          total, 
          onPageChange, 
          onPageSizeChange, 
          pageSizeOptions, 
        }}
      />

      <PurchaseFormModal
        open={formOpen}
        onOpenChange={(o) => { if (!o) handleCloseForm() }}
        initialData={editingPurchase}
      />

      <PurchaseDetailModal
        open={detailOpen}
        onOpenChange={(open) => { if (!open) handleCloseDetail() }}
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
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingPurchase(null) } }}
        title="Hapus Pembelian"
        description={`Yakin ingin menghapus pembelian "${deletingPurchase?.invoice_number}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
