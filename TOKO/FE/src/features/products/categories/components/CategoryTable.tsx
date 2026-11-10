import { forwardRef, useImperativeHandle, useState } from 'react'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'

import {
  useCategoryListQuery,
  useDeleteCategoryMutation,
  useToggleCategoryStatusMutation,
} from '../categories.api'
import type { Category, CategoryListFilter } from '../categories.types'
import { CategoryFilterBar } from './CategoryFilterBar'
import { CategoryFormModal } from './CategoryFormModal'
import { buildCategoryColumns } from './CategoryTableColumns'

export interface CategoryTableHandle {
  openAdd: () => void
}

export const CategoryTable = forwardRef<CategoryTableHandle, object>(function CategoryTable(_, ref) {
  const [filter, setFilter] = useState<CategoryListFilter>({ page: 1, limit: 10, search: '' })
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)

  const { page, pageSize, onPageChange, onPageSizeChange, reset } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const { isOpen: deleteOpen, open: openDelete, close: closeDelete } = useDisclosure()
  const { isOpen: toggleOpen, open: openToggle, close: closeToggle } = useDisclosure()

  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [deletingCategory, setDeletingCategory] = useState<Category | null>(null)
  const [togglingCategory, setTogglingCategory] = useState<Category | null>(null)

  const { data: categoryData, isLoading } = useCategoryListQuery({ ...filter, page, limit: pageSize })
  const categories = categoryData?.data ?? []
  const total = categoryData?.total ?? 0

  const { mutate: deleteCategory, isPending: isDeleting } = useDeleteCategoryMutation()
  const { mutate: toggleStatus, isPending: isToggling } = useToggleCategoryStatusMutation()

  const handleOpenAdd = () => {
    setEditingCategory(null)
    openForm()
  }

  useImperativeHandle(ref, () => ({ openAdd: handleOpenAdd }))

  const handleFilterChange = (newFilter: CategoryListFilter) => {
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

  const handleOpenEdit = (category: Category) => {
    setEditingCategory(category)
    openForm()
  }

  const handleCloseForm = () => {
    closeForm()
    setEditingCategory(null)
  }

  const handleOpenDelete = (category: Category) => {
    setDeletingCategory(category)
    openDelete()
  }

  const handleCloseDelete = () => {
    closeDelete()
    setDeletingCategory(null)
  }

  const handleConfirmDelete = () => {
    if (!deletingCategory) return
    deleteCategory(deletingCategory.id, {
      onSuccess: () => handleCloseDelete(),
    })
  }

  const handleToggleStatus = (id: number) => {
    const category = categories.find((c) => c.id === id)
    if (!category) return
    setTogglingCategory(category)
    openToggle()
  }

  const handleCloseToggle = () => {
    closeToggle()
    setTogglingCategory(null)
  }

  const handleConfirmToggle = () => {
    if (!togglingCategory) return
    toggleStatus(
      { id: togglingCategory.id, isActive: togglingCategory.is_active },
      { onSuccess: () => handleCloseToggle() }
    )
  }

  const hasFilter = filter.search || filter.is_active !== undefined

  const columns = buildCategoryColumns({
    onEdit: handleOpenEdit,
    onDelete: handleOpenDelete,
    onToggleStatus: handleToggleStatus,
  })


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
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Kategori tidak ditemukan' : 'Belum ada kategori'}
        emptyDescription={
          hasFilter
            ? 'Coba ubah kata kunci atau filter pencarian Anda.'
            : 'Tambah kategori pertama Anda untuk memulai.'
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

      <CategoryFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
        category={editingCategory}
      />

      <ConfirmDialog
        open={deleteOpen}
        onOpenChange={(open) => { if (!open) handleCloseDelete() }}
        title="Hapus Kategori"
        description={`Yakin ingin menghapus kategori "${deletingCategory?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />

      <ConfirmDialog
        open={toggleOpen}
        onOpenChange={(open) => { if (!open) handleCloseToggle() }}
        title={togglingCategory?.is_active ? 'Nonaktifkan Kategori' : 'Aktifkan Kategori'}
        description={`Kategori "${togglingCategory?.name}" akan di${togglingCategory?.is_active ? 'nonaktifkan' : 'aktifkan'}. Lanjutkan?`}
        confirmLabel={togglingCategory?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={togglingCategory?.is_active ? 'destructive' : 'default'}
        isLoading={isToggling}
        onConfirm={handleConfirmToggle}
      />
    </div>
  )
})
