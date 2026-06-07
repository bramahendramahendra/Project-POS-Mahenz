import { useRef, useState } from 'react'
import { Plus, Upload, Tag, Ruler } from 'lucide-react'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, PageHeader, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { usePagination, useDisclosure, usePageSizeOptions } from '@/shared/hooks'

import { useDeleteProductMutation, useProductListQuery } from './products.api'
import { useUnitOptionsQuery } from '@/features/inventory/units'
import { useCategoryOptionsQuery } from '@/features/inventory/categories'
import type { Product, ProductListFilter } from './products.types'
import { ImportCsvModal } from './components/ImportCsvModal'
import { LabelPrintModal } from './components/LabelPrintModal'
import { ProductDetailModal } from './components/ProductDetailModal'
import { ProductFilterBar } from './components/ProductFilter'
import { ProductFormModal } from './components/ProductFormModal'
import type { ProductTableHandle } from './components/ProductTable'
import { ProductTable } from './components/ProductTable'

export function ProductsPage() {
  const tableRef = useRef<ProductTableHandle>(null)
  const [filter, setFilter] = useState<ProductListFilter>({ page: 1, limit: 10, search: '' })
  const [selectedProducts, setSelectedProducts] = useState<Product[]>([])
  const [singleLabelProduct, setSingleLabelProduct] = useState<Product | null>(null)

  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [detailProduct, setDetailProduct] = useState<Product | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Product | null>(null)

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: importOpen, open: openImport, close: closeImport } = useDisclosure()
  const { isOpen: labelOpen, open: openLabel, close: closeLabel } = useDisclosure()

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: productData, isLoading } = useProductListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const { data: categories = [], isLoading: isCatLoading } = useCategoryOptionsQuery()
  const { data: units = [], isLoading: isUnitLoading } = useUnitOptionsQuery()
  const { mutate: deleteProduct, isPending: isDeleting } = useDeleteProductMutation()

  const hasCategories = categories.length > 0
  const hasActiveUnits = units.length > 0

  const products = productData?.data ?? []
  const total = productData?.total ?? 0

  const handleFilterChange = (newFilter: ProductListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    reset()
  }

  const handleOpenEdit = (product: Product) => {
    setEditingProduct(product)
    openForm()
  }

  const handleOpenAdd = () => {
    setEditingProduct(null)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingProduct(null)
  }

  const handleOpenDelete = (product: Product) => {
    setDeleteTarget(product)
    openDelete()
  }

  const handleConfirmDelete = () => {
    if (!deleteTarget) return
    deleteProduct(deleteTarget.id, {
      onSuccess: () => {
        closeDelete()
        setDeleteTarget(null)
      },
    })
  }

  const handleOpenDetail = (product: Product) => {
    setDetailProduct(product)
    openDetail()
  }

  const handleCloseDetail = () => {
    closeDetail()
    setDetailProduct(null)
  }

  if (isCatLoading || isUnitLoading) {
    return (
      <div className="space-y-4">
        <PageHeader
          title="Produk"
          breadcrumbs={[{ label: 'Inventori' }, { label: 'Produk' }]}
        />
        <div className="space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
          ))}
        </div>
      </div>
    )
  }

  if (!hasCategories || !hasActiveUnits) {
    return (
      <div className="space-y-4">
        <PageHeader
          title="Produk"
          breadcrumbs={[{ label: 'Inventori' }, { label: 'Produk' }]}
        />
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
          <div className="mb-4 flex gap-3">
            <div className={`rounded-full p-3 ${!hasCategories ? 'bg-amber-50' : 'bg-green-50'}`}>
              <Tag size={24} className={!hasCategories ? 'text-amber-500' : 'text-green-500'} />
            </div>
            <div className={`rounded-full p-3 ${!hasActiveUnits ? 'bg-amber-50' : 'bg-green-50'}`}>
              <Ruler size={24} className={!hasActiveUnits ? 'text-amber-500' : 'text-green-500'} />
            </div>
          </div>
          <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa menambah produk</h3>
          <p className="mb-1 text-sm text-gray-500">
            Sebelum menambah produk, pastikan data berikut sudah tersedia:
          </p>
          <ul className="mb-6 text-sm">
            <li className={`flex items-center gap-2 ${hasCategories ? 'text-green-600' : 'text-amber-600'}`}>
              <span>{hasCategories ? '✓' : '!'}</span>
              {hasCategories ? 'Kategori sudah tersedia' : 'Belum ada kategori — tambahkan di tab Kategori'}
            </li>
            <li className={`flex items-center gap-2 ${hasActiveUnits ? 'text-green-600' : 'text-amber-600'}`}>
              <span>{hasActiveUnits ? '✓' : '!'}</span>
              {hasActiveUnits ? 'Satuan sudah tersedia' : 'Belum ada satuan aktif — tambahkan di tab Satuan'}
            </li>
          </ul>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <PageHeader
        title="Produk"
        breadcrumbs={[{ label: 'Inventori' }, { label: 'Produk' }]}
        actions={
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <div className="flex gap-2">
              <Button variant="outline" onClick={openImport} className="gap-1">
                <Upload size={16} />
                Import Produk
              </Button>
              <Button onClick={handleOpenAdd} className="gap-1">
                <Plus size={16} />
                Tambah Produk
              </Button>
            </div>
          </RoleGuard>
        }
      />

      <div className="space-y-3">
        <ProductFilterBar
          filter={filter}
          onChange={handleFilterChange}
          onReset={handleReset}
          categories={categories}
        />
        <ProductTable
          ref={tableRef}
          data={products}
          isLoading={isLoading}
          pagination={{
            page,
            pageSize,
            total,
            onPageChange,
            onPageSizeChange,
            pageSizeOptions,
          }}
          onSelectionChange={setSelectedProducts}
          onEdit={handleOpenEdit}
          onDelete={handleOpenDelete}
          onDetail={handleOpenDetail}
          onPrintLabel={() => openLabel()}
          onPrintSingleLabel={(product) => {
            setSingleLabelProduct(product)
            openLabel()
          }}
        />
      </div>

      <ProductFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        product={editingProduct}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeleteTarget(null) } }}
        title="Hapus Produk"
        description={`Yakin ingin menghapus produk "${deleteTarget?.name}"? Tindakan ini tidak bisa dibatalkan.`}
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
        onOpenChange={(open) => {
          if (!open) { closeLabel(); setSingleLabelProduct(null) }
        }}
        products={singleLabelProduct ? [singleLabelProduct] : selectedProducts}
      />

      <ProductDetailModal
        open={detailOpen}
        onOpenChange={(open) => { if (!open) handleCloseDetail() }}
        productId={detailProduct?.id}
      />
    </div>
  )
}
