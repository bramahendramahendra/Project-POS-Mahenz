import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions, useTableSelection } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import {
  useProductListQuery,
  useDeleteProductMutation,
  useBulkToggleProductStatusMutation,
  useToggleProductStatusMutation,
} from '../products.api'
import { useCategoryOptionsQuery } from '@/features/products/categories'
import type { Product, ProductListFilter } from '../products.types'
import { exportProductsToExcel } from '../products.utils'
import { buildProductColumns } from './ProductTableColumns'
import { ImportCsvModal } from './ImportCsvModal'
import { LabelPrintModal } from './LabelPrintModal'
import { ProductDetailModal } from './ProductDetailModal'
import { ProductBulkActionBar } from './ProductBulkActionBar'
import { ProductFilterBar } from './ProductFilterBar'
import { ProductFormModal } from './ProductFormModal'

export interface ProductTableHandle {
  openAdd: () => void
  openImport: () => void
}

export const ProductTable = forwardRef<ProductTableHandle, object>(function ProductTable(_, ref) {
  const [filter, setFilter] = useState<ProductListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)
  
  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: importOpen, open: openImport, close: closeImport } = useDisclosure()
  const { isOpen: labelOpen, open: openLabel, close: closeLabel } = useDisclosure()

  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [deletingProduct, setDeletingProduct] = useState<Product | null>(null)
  const [detailProduct, setDetailProduct] = useState<Product | null>(null)
  const [singleLabelProduct, setSingleLabelProduct] = useState<Product | null>(null)

  const { data: productData, isLoading } = useProductListQuery({ ...filter, page, limit: pageSize })
  const products = productData?.data ?? []
  const total = productData?.total ?? 0

  const { data: categories = [] } = useCategoryOptionsQuery()

  const { mutate: deleteProduct, isPending: isDeleting } = useDeleteProductMutation()
  const { mutate: toggleStatus } = useToggleProductStatusMutation()
  const { mutate: bulkToggleStatus, isPending: isBulkToggling } = useBulkToggleProductStatusMutation()
  const { selectedKeys, toggle, selectAll, clearSelection, hasSelection, count } =
    useTableSelection<Product & { id: number }>()

  const selectedProducts = products.filter((p) => selectedKeys.has(p.id))
  const allActive = selectedProducts.length > 0 && selectedProducts.every((p) => p.is_active)
  const allInactive = selectedProducts.length > 0 && selectedProducts.every((p) => !p.is_active)

  const handleOpenAdd = () => {
    setEditingProduct(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd, openImport }))

  const handleFilterChange = (newFilter: ProductListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    setSortState(undefined)
    reset()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    reset()
  }

  const handleOpenEdit = (product: Product) => {
    setEditingProduct(product)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingProduct(null)
  }

  const handleOpenDetail = (product: Product) => {
    setDetailProduct(product)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailProduct(null)
  }


  const handleOpenDelete = (product: Product) => {
    setDeletingProduct(product)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deletingProduct) return
    deleteProduct(deletingProduct.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingProduct(null)
      },
    })
  }

  function handleExportExcel() {
    exportProductsToExcel(selectedProducts)
  }

  function handleBulkToggleStatus() {
    const ids = selectedProducts.map((p) => p.id)
    const label = allActive ? 'dinonaktifkan' : 'diaktifkan'
    bulkToggleStatus(ids, {
      onSuccess: () => {
        toast.success(`${ids.length} produk berhasil ${label}`)
        clearSelection()
      },
    })
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Produk berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const hasFilter = filter.search || filter.category_id || filter.is_active !== undefined

  const columns = buildProductColumns({
    onDetail: handleOpenDetail,
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onLabel: (product) => {
      setSingleLabelProduct(product)
      openLabel()
    },
    onToggleStatus: handleToggleStatus,
  })

  return (
    <div className="space-y-4">
      <ProductFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        categories={categories}
      />

      {hasSelection && (
        <ProductBulkActionBar
          count={count}
          allActive={allActive}
          allInactive={allInactive}
          isBulkToggling={isBulkToggling}
          onToggleStatus={handleBulkToggleStatus}
          onExport={handleExportExcel}
          onPrintLabel={openLabel}
          onClear={clearSelection}
        />
      )}

      <DataTable<Product & Record<string, unknown>>
        columns={columns}
        data={products as (Product & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Produk tidak ditemukan' : 'Belum ada produk'}
        emptyDescription={
          hasFilter 
            ? 'Coba ubah filter atau kata kunci pencarian Anda.' 
            : 'Tambah produk pertama Anda untuk memulai.'
        }
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        rowSelection={{
          enabled: true,
          rowKey: 'id',
          selectedKeys,
          onSelectionChange: (keys) => {
            if (keys.size === 0) {
              clearSelection()
            } else if (keys.size >= products.length) {
              selectAll(products as (Product & { id: number })[])
            } else {
              const added = [...keys].find((k) => !selectedKeys.has(k))
              const removed = [...selectedKeys].find((k) => !keys.has(k))
              if (added !== undefined) toggle(added)
              else if (removed !== undefined) toggle(removed)
            }
          },
        }}
      />

      <ProductFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        product={editingProduct}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingProduct(null) } }}
        title="Hapus Produk"
        description={`Yakin ingin menghapus produk "${deletingProduct?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />

      <ImportCsvModal
        open={importOpen}
        onOpenChange={(open) => { if (!open) closeImport() }}
      />

      <LabelPrintModal
        open={labelOpen}
        onOpenChange={(open) => { if (!open) { closeLabel(); setSingleLabelProduct(null) } }}
        products={singleLabelProduct ? [singleLabelProduct] : selectedProducts}
      />

      <ProductDetailModal
        open={detailOpen}
        onOpenChange={(open) => { if (!open) handleCloseDetail() }}
        productId={detailProduct?.id}
      />
    </div>
  )
})
