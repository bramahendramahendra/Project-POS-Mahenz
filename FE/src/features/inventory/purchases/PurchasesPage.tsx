import { useState } from 'react'
import { Plus, Truck } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import { useSupplierListQuery } from '@/features/inventory/suppliers'
import { useProductListQuery } from '@/features/inventory/products'

import {
  useSupplierPurchasesQuery,
  useDeleteSupplierPurchaseMutation,
} from './purchases.api'
import type { SupplierPurchase, SupplierPurchaseFilter } from './purchases.types'
import { PurchaseTable } from './components/PurchaseTable'
import { PurchaseFilterBar } from './components/PurchaseFilterBar'
import { PurchaseFormModal } from './components/PurchaseFormModal'
import { PurchaseDetailModal } from './components/PurchaseDetailModal'
import { PaymentModal } from './components/PaymentModal'

function monthStartString() {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

function todayString() {
  return new Date().toISOString().split('T')[0]
}

export function PurchasesPage() {
  const [filter, setFilter] = useState<SupplierPurchaseFilter>({
    start_date: monthStartString(),
    end_date: todayString(),
  })

  const { page, pageSize, onPageChange, onPageSizeChange } = usePagination()
  const pageSizeOptions = usePageSizeOptions()
  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: payOpen, open: openPay, close: closePay } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [payingPurchase, setPayingPurchase] = useState<SupplierPurchase | null>(null)
  const [detailId, setDetailId] = useState<number | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const { data: suppliersData, isLoading: isSuppliersLoading } = useSupplierListQuery({ page: 1, limit: 200, search: '' })
  const { data: productsData, isLoading: isProductsLoading } = useProductListQuery({ page: 1, limit: 1, search: '' })
  const suppliers = suppliersData?.data ?? []
  const hasSuppliers = suppliers.length > 0
  const hasProducts = (productsData?.total ?? 0) > 0

  const { data, isLoading } = useSupplierPurchasesQuery({ ...filter, page, limit: pageSize })
  const { mutate: deletePurchase, isPending: isDeleting } = useDeleteSupplierPurchaseMutation()

  const purchases = data?.data ?? []
  const total = data?.total ?? 0

  function handlePay(purchase: SupplierPurchase) {
    setPayingPurchase(purchase)
    openPay()
  }

  function handleDelete(purchase: SupplierPurchase) {
    setDeletingId(purchase.id)
    openDelete()
  }

  function handleDetail(purchase: SupplierPurchase) {
    setDetailId(purchase.id)
    openDetail()
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

  if (isSuppliersLoading || isProductsLoading) {
    return (
      <div className="space-y-4">
        <PageHeader
          title="Pembelian Supplier"
          breadcrumbs={[{ label: 'Inventori' }, { label: 'Pembelian' }]}
        />
        <div className="space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
          ))}
        </div>
      </div>
    )
  }

  if (!hasSuppliers || !hasProducts) {
    return (
      <div className="space-y-4">
        <PageHeader
          title="Pembelian Supplier"
          breadcrumbs={[{ label: 'Inventori' }, { label: 'Pembelian' }]}
        />
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
          <div className="mb-4 flex gap-3">
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
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <PageHeader
        title="Pembelian Supplier"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Pembelian' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button onClick={openForm} className="gap-1">
              <Plus size={16} />
              Tambah Pembelian
            </Button>
          </RoleGuard>
        }
      />

      <PurchaseFilterBar filter={filter} suppliers={suppliers} onChange={setFilter} />

      <PurchaseTable
        data={purchases}
        isLoading={isLoading}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
        suppliers={suppliers}
        onDetail={handleDetail}
        onPay={handlePay}
        onDelete={handleDelete}
      />

      <PurchaseFormModal open={formOpen} onOpenChange={(o) => !o && closeForm()} />

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
}
