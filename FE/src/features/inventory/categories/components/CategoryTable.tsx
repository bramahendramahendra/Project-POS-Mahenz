import { forwardRef, useImperativeHandle, useState } from 'react'
import { Lock, LockOpen, Pencil, Trash2 } from 'lucide-react'
import { toast } from 'sonner'

import { ROLES } from '@/shared/constants'
import { ConfirmDialog, DataTable, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'

import {
  useCategoryListQuery,
  useDeleteCategoryMutation,
  useToggleCategoryStatusMutation,
} from '../categories.api'
import type { Category, CategoryListFilter } from '../categories.types'
import { CategoryFilterBar } from './CategoryFilterBar'
import { CategoryFormModal } from './CategoryFormModal'

export interface CategoryTableHandle {
  openAdd: () => void
}

export const CategoryTable = forwardRef<CategoryTableHandle, object>(function CategoryTable(_, ref) {
  const [filter, setFilter] = useState<CategoryListFilter>({ page: 1, limit: 10, search: '' })

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()

  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [deletingCategory, setDeletingCategory] = useState<Category | null>(null)

  const { data: categoryData, isLoading } = useCategoryListQuery({
    ...filter,
    page,
    limit: pageSize,
  })
  const categories = categoryData?.data ?? []
  const total = categoryData?.total ?? 0

  const { mutate: deleteCategory, isPending: isDeleting } = useDeleteCategoryMutation()
  const { mutate: toggleStatus } = useToggleCategoryStatusMutation()

  const handleOpenAdd = () => {
    setEditingCategory(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleOpenEdit = (category: Category) => {
    setEditingCategory(category)
    openForm()
  }

  const handleOpenDelete = (category: Category) => {
    setDeletingCategory(category)
    openDelete()
  }

  const handleFilterChange = (newFilter: CategoryListFilter) => {
    setFilter(newFilter)
    resetPage()
  }

  const handleReset = () => {
    setFilter({ page: 1, limit: 10, search: '' })
    resetPage()
  }

  const handleDelete = () => {
    if (!deletingCategory) return
    deleteCategory(deletingCategory.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingCategory(null)
      },
    })
  }

  const handleToggleStatus = (row: Category) => {
    toggleStatus(row.id, {
      onSuccess: () =>
        toast.success(`Kategori berhasil ${row.is_active ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
  }

  const columns: ColumnDef<Category>[] = [
    {
      key: 'code',
      header: 'Kode',
      width: '80px',
      cell: (row) => (
        <span className="font-mono text-xs font-semibold text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded">
          {row.code}
        </span>
      ),
    },
    {
      key: 'name',
      header: 'Nama Kategori',
      cell: (row) => (
        <span className="font-medium text-gray-800">
          {row.name}
        </span>
      ),
    },
    {
      key: 'description',
      header: 'Deskripsi',
      cell: (row) =>
        row.description ? (
          <span className="text-sm text-gray-600">{row.description}</span>
        ) : (
          <span className="text-gray-400 text-sm">—</span>
        ),
    },
    {
      key: 'product_count',
      header: 'Jumlah Produk',
      align: 'center',
      width: '130px',
      cell: (row) => (
        <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs text-gray-600">
          {row.product_count} produk
        </span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      align: 'center',
      width: '100px',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '120px',
      cell: (row) => (
        <div className="flex items-center justify-center gap-1">
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
              onClick={() => handleToggleStatus(row)}
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

  const hasFilter = filter.search || filter.is_active !== undefined

  return (
    <div className="space-y-4">
      <CategoryFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
      />

      <DataTable<Category & Record<string, unknown>>
        columns={columns}
        data={categories as (Category & Record<string, unknown>)[]}
        isLoading={isLoading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange,
          onPageSizeChange,
          pageSizeOptions,
        }}
        emptyMessage={hasFilter ? 'Kategori tidak ditemukan' : 'Belum ada kategori'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'Tambah kategori pertama Anda untuk memulai.'
        }
      />

      <CategoryFormModal
        open={formOpen}
        onOpenChange={(val) => { if (!val) { closeForm(); setEditingCategory(null) } }}
        category={editingCategory}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => {
          if (!open) { closeDelete(); setDeletingCategory(null) }
        }}
        title="Hapus Kategori"
        description={`Yakin ingin menghapus kategori "${deletingCategory?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleDelete}
      />
    </div>
  )
})
