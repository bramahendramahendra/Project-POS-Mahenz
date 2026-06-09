import { forwardRef, useImperativeHandle, useState } from 'react'
import { Eye, FileDown, Lock, LockOpen, Pencil, Printer, Ruler, Tag, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import * as XLSX from 'xlsx'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, DataTable, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { useDisclosure, usePagination, usePageSizeOptions, useTableSelection } from '@/shared/hooks'
import { formatRupiah } from '@/shared/utils'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import { calcMargin } from '../products.utils'
import {
  useBulkToggleProductStatusMutation,
  useDeleteProductMutation,
  useProductListQuery,
  useToggleProductStatusMutation,
} from '../products.api'
import { useCategoryOptionsQuery } from '@/features/inventory/categories'
import { useUnitOptionsQuery } from '@/features/inventory/units'
import type { Product, ProductListFilter } from '../products.types'
import { ImportCsvModal } from './ImportCsvModal'
import { LabelPrintModal } from './LabelPrintModal'
import { ProductDetailModal } from './ProductDetailModal'
import { ProductFilterBar } from './ProductFilter'
import { ProductFormModal } from './ProductFormModal'

export interface ProductTableHandle {
  openAdd: () => void
  openImport: () => void
}

export const ProductTable = forwardRef<ProductTableHandle, object>(function ProductTable(_, ref) {
  const [filter, setFilter] = useState<ProductListFilter>({ page: 1, limit: 10, search: '' })
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [detailProduct, setDetailProduct] = useState<Product | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Product | null>(null)
  const [singleLabelProduct, setSingleLabelProduct] = useState<Product | null>(null)

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: detailOpen, open: openDetail, close: closeDetail } = useDisclosure()
  const { isOpen: importOpen, open: openImport, close: closeImport } = useDisclosure()
  const { isOpen: labelOpen, open: openLabel, close: closeLabel } = useDisclosure()

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination()
  const pageSizeOptions = usePageSizeOptions()

  const { data: productData, isLoading } = useProductListQuery({ ...filter, page, limit: pageSize })
  const { data: categories = [], isLoading: isCatLoading } = useCategoryOptionsQuery()
  const { data: units = [], isLoading: isUnitLoading } = useUnitOptionsQuery()
  const { mutate: deleteProduct, isPending: isDeleting } = useDeleteProductMutation()
  const { mutate: toggleStatus } = useToggleProductStatusMutation()
  const { mutate: bulkToggleStatus, isPending: isBulkToggling } = useBulkToggleProductStatusMutation()
  const { selectedKeys, toggle, selectAll, clearSelection, hasSelection, count } =
    useTableSelection<Product & { id: number }>()

  const products = productData?.data ?? []
  const total = productData?.total ?? 0

  const selectedProducts = products.filter((p) => selectedKeys.has(p.id))
  const allActive = selectedProducts.length > 0 && selectedProducts.every((p) => p.is_active)
  const allInactive = selectedProducts.length > 0 && selectedProducts.every((p) => !p.is_active)
  const showBulkStatus = allActive || allInactive

  const handleOpenAdd = () => {
    setEditingProduct(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd, openImport }))

  const handleOpenEdit = (product: Product) => {
    setEditingProduct(product)
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

  const handleFilterChange = (newFilter: ProductListFilter) => {
    setFilter(newFilter)
    reset()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    reset()
  }

  function handleExportExcel() {
    const rows = selectedProducts.map((p) => ({
      'Nama Produk': p.name,
      Barcode: p.barcode ?? '',
      SKU: p.sku ?? '',
      Kategori: p.category_name ?? '',
      'Harga Beli': p.purchase_price,
      'Harga Jual': p.selling_price,
      Stok: p.stock,
      'Stok Minimum': p.min_stock,
      Satuan: p.unit_name ?? '',
      Status: p.is_active ? 'Aktif' : 'Nonaktif',
    }))
    const ws = XLSX.utils.json_to_sheet(rows)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Produk')
    XLSX.writeFile(wb, `produk-export-${Date.now()}.xlsx`)
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

  const hasCategories = categories.length > 0
  const hasActiveUnits = units.length > 0

  if (isCatLoading || isUnitLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  if (!hasCategories || !hasActiveUnits) {
    return (
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
    )
  }

  const columns: ColumnDef<Product>[] = [
    {
      key: 'name',
      header: 'Nama Produk',
      sortable: true,
      cell: (row) => <span className="font-medium text-gray-800">{row.name}</span>,
    },
    {
      key: 'barcode',
      header: 'Barcode / SKU',
      cell: (row) => (
        <div className="flex flex-col gap-0.5">
          {row.barcode ? (
            <code className="text-xs text-gray-700">{row.barcode}</code>
          ) : (
            <span className="text-gray-400 text-xs">—</span>
          )}
          {row.sku ? (
            <span className="text-xs text-gray-500">{row.sku}</span>
          ) : (
            <span className="text-gray-300 text-xs">—</span>
          )}
        </div>
      ),
    },
    {
      key: 'category_name',
      header: 'Kategori',
      cell: (row) =>
        row.category_name ? (
          <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600">
            {row.category_name}
          </span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'purchase_price',
      header: 'Harga Beli',
      align: 'right',
      cell: (row) => <span className="text-sm">{formatRupiah(row.purchase_price)}</span>,
    },
    {
      key: 'selling_price',
      header: 'Harga Jual',
      align: 'right',
      cell: (row) => <span className="font-medium">{formatRupiah(row.selling_price)}</span>,
    },
    {
      key: 'margin',
      header: 'Margin',
      align: 'center',
      width: '80px',
      cell: (row) => {
        const m = calcMargin(row.purchase_price, row.selling_price)
        return (
          <span
            className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
              m >= 30
                ? 'bg-green-100 text-green-700'
                : m >= 15
                  ? 'bg-amber-100 text-amber-700'
                  : 'bg-red-100 text-red-600'
            }`}
          >
            {m}%
          </span>
        )
      },
    },
    {
      key: 'stock',
      header: 'Stok',
      align: 'right',
      width: '80px',
      cell: (row) => (
        <span
          className={`font-medium ${
            row.stock === 0
              ? 'text-red-600'
              : row.stock < row.min_stock
                ? 'text-amber-600'
                : 'text-gray-800'
          }`}
        >
          {row.stock}
        </span>
      ),
    },
    {
      key: 'unit_name',
      header: 'Satuan',
      width: '80px',
      cell: (row) =>
        row.unit_name ? (
          <span className="text-sm text-gray-600">{row.unit_name}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '90px',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '130px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-gray-500 hover:text-blue-600"
            onClick={() => handleOpenDetail(row)}
            title="Lihat Detail"
          >
            <Eye size={14} />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-gray-500 hover:text-indigo-600"
            onClick={() => {
              setSingleLabelProduct(row)
              openLabel()
            }}
            title="Cetak Label"
          >
            <Printer size={14} />
          </Button>
          <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-gray-500 hover:text-blue-600"
              onClick={() => handleOpenEdit(row)}
              title="Edit"
            >
              <Pencil size={14} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className={`h-7 w-7 ${row.is_active ? 'text-gray-500 hover:text-amber-600' : 'text-gray-400 hover:text-green-600'}`}
              onClick={() =>
                toggleStatus(row.id, {
                  onSuccess: () =>
                    toast.success(`Produk berhasil ${row.is_active ? 'dinonaktifkan' : 'diaktifkan'}`),
                })
              }
              title={row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
            >
              {row.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
            </Button>
          </RoleGuard>
          <RoleGuard allowedRoles={[ROLES.OWNER]}>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-gray-500 hover:text-red-600"
              onClick={() => handleOpenDelete(row)}
              title="Hapus"
            >
              <Trash2 size={14} />
            </Button>
          </RoleGuard>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <ProductFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        categories={categories}
      />

      {hasSelection && (
        <div className="flex items-center gap-3 rounded-lg border bg-blue-50 px-4 py-2 text-sm">
          <span className="font-medium text-blue-700">{count} produk dipilih</span>
          <div className="ml-auto flex gap-2">
            {showBulkStatus && (
              <Button
                variant="outline"
                size="sm"
                disabled={isBulkToggling}
                onClick={handleBulkToggleStatus}
                className={allActive ? 'text-amber-600 hover:text-amber-700' : 'text-green-600 hover:text-green-700'}
              >
                {isBulkToggling ? 'Memproses...' : allActive ? 'Nonaktifkan' : 'Aktifkan'}
              </Button>
            )}
            <Button variant="outline" size="sm" onClick={handleExportExcel} className="gap-1">
              <FileDown size={14} />
              Export Excel
            </Button>
            <Button variant="outline" size="sm" onClick={openLabel}>
              Cetak Label
            </Button>
            <Button variant="outline" size="sm" onClick={() => clearSelection()}>
              Batalkan Pilihan
            </Button>
          </div>
        </div>
      )}

      <DataTable<Product & Record<string, unknown>>
        columns={columns}
        data={products as (Product & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada produk"
        emptyDescription="Tambah produk pertama Anda untuk memulai."
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
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
