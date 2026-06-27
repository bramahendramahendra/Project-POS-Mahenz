import { forwardRef, useImperativeHandle, useState } from 'react'
import { toast } from 'sonner'

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

  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [deletingCategory, setDeletingCategory] = useState<Category | null>(null)

  const { data: categoryData, isLoading } = useCategoryListQuery({ ...filter, page, limit: pageSize })
  const categories = categoryData?.data ?? []
  const total = categoryData?.total ?? 0

  const { mutate: deleteCategory, isPending: isDeleting } = useDeleteCategoryMutation()
  const { mutate: toggleStatus } = useToggleCategoryStatusMutation()

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


  const handleConfirmDelete = () => {
    if (!deletingCategory) return
    deleteCategory(deletingCategory.id, {
      onSuccess: () => {
        closeDelete()
        setDeletingCategory(null)
      },
    })
  }

  const handleToggleStatus = (id: number, isActive: boolean) => {
    toggleStatus(id, {
      onSuccess: () =>
        toast.success(`Kategori berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`),
    })
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
        onOpenChange={(open) => { if (!open) { closeDelete(); setDeletingCategory(null) } }}
        title="Hapus Kategori"
        description={`Yakin ingin menghapus kategori "${deletingCategory?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </div>
  )
})
