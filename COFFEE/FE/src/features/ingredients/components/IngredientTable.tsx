import { forwardRef, useImperativeHandle, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { ConfirmDialog, DataTable } from '@/shared/components'
import { useDisclosure, usePagination, usePageSizeOptions } from '@/shared/hooks'
import type { SortState } from '@/shared/components/DataTable/DataTable.types'
import { ROUTES } from '@/shared/constants/routes'

import {
  useDeleteIngredientMutation,
  useIngredientListQuery,
  useToggleIngredientStatusMutation,
} from '../ingredients.api'
import type { Ingredient, IngredientListFilter } from '../ingredients.types'
import { IngredientFilterBar } from './IngredientFilterBar'
import { IngredientFormModal } from './IngredientFormModal'
import { buildIngredientColumns } from './IngredientTableColumns'
import { RecordStockModal, type RecordMode } from './RecordStockModal'

export interface IngredientTableHandle {
  openAdd: () => void
}

const initialFilter: IngredientListFilter = { page: 1, limit: 10, search: '' }

export const IngredientTable = forwardRef<IngredientTableHandle, object>(function IngredientTable(_, ref) {
  const navigate = useNavigate()
  const [filter, setFilter] = useState<IngredientListFilter>(initialFilter)
  const [sortState, setSortState] = useState<SortState | undefined>(undefined)
  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()
  const [editing, setEditing] = useState<Ingredient | null>(null)
  const [recording, setRecording] = useState<{ mode: RecordMode; ingredient: Ingredient } | null>(null)
  const [toggling, setToggling] = useState<Ingredient | null>(null)
  const [deleting, setDeleting] = useState<Ingredient | null>(null)

  const { data, isLoading } = useIngredientListQuery({ ...filter, page, limit: pageSize })
  const rows = data?.data ?? []
  const total = data?.total ?? 0

  const { mutate: toggleStatus, isPending: isToggling } = useToggleIngredientStatusMutation()
  const { mutate: deleteIngredient, isPending: isDeleting } = useDeleteIngredientMutation()

  const openAdd = () => {
    setEditing(null)
    openForm()
  }
  useImperativeHandle(ref, () => ({ openAdd }))

  const handleFilterChange = (next: IngredientListFilter) => {
    setFilter(next)
    resetPage()
  }

  const handleReset = () => {
    setFilter(initialFilter)
    setSortState(undefined)
    resetPage()
  }

  const handleSort = (sort: SortState) => {
    setSortState(sort)
    setFilter((prev) => ({ ...prev, sort_by: sort.key, sort_order: sort.order }))
    resetPage()
  }

  const columns = buildIngredientColumns({
    onUsage: (ingredient) => setRecording({ mode: 'usage', ingredient }),
    onWaste: (ingredient) => setRecording({ mode: 'waste', ingredient }),
    onHistory: (ingredient) => navigate(`${ROUTES.INGREDIENTS_HISTORY}?ingredient_id=${ingredient.id}`),
    onEdit: (ingredient) => {
      setEditing(ingredient)
      openForm()
    },
    onToggleStatus: setToggling,
    onDelete: setDeleting,
  })

  const hasFilter = !!filter.search || filter.is_active !== undefined || !!filter.low_stock

  return (
    <div className="space-y-4">
      <IngredientFilterBar filter={filter} onChange={handleFilterChange} onReset={handleReset} />

      <DataTable<Ingredient & Record<string, unknown>>
        columns={columns}
        data={rows as (Ingredient & Record<string, unknown>)[]}
        isLoading={isLoading}
        currentSort={sortState}
        onSort={handleSort}
        emptyMessage={hasFilter ? 'Bahan tidak ditemukan' : 'Belum ada bahan'}
        emptyDescription={
          hasFilter ? 'Coba ubah kata kunci atau filter pencarian Anda.' : 'Tambah bahan pertama untuk mulai mencatat stok.'
        }
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <IngredientFormModal
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) {
            closeForm()
            setEditing(null)
          }
        }}
        ingredient={editing}
      />

      <RecordStockModal
        mode={recording?.mode ?? 'usage'}
        ingredient={recording?.ingredient ?? null}
        onOpenChange={(open) => {
          if (!open) setRecording(null)
        }}
      />

      <ConfirmDialog
        open={toggling != null}
        onOpenChange={(open) => {
          if (!open) setToggling(null)
        }}
        title={toggling?.is_active ? 'Nonaktifkan Bahan' : 'Aktifkan Bahan'}
        description={
          toggling?.is_active
            ? `Bahan "${toggling?.name}" akan dinonaktifkan dan tidak bisa dicatat pemakaiannya. Riwayat tetap tersimpan. Lanjutkan?`
            : `Bahan "${toggling?.name}" akan diaktifkan kembali. Lanjutkan?`
        }
        confirmLabel={toggling?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={toggling?.is_active ? 'destructive' : 'default'}
        isLoading={isToggling}
        onConfirm={() =>
          toggling && toggleStatus({ id: toggling.id, willActivate: !toggling.is_active }, { onSuccess: () => setToggling(null) })
        }
      />

      <ConfirmDialog
        open={deleting != null}
        onOpenChange={(open) => {
          if (!open) setDeleting(null)
        }}
        title="Hapus Bahan"
        description={`Yakin ingin menghapus bahan "${deleting?.name}"? Hanya bahan tanpa riwayat stok yang bisa dihapus.`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={() => deleting && deleteIngredient(deleting.id, { onSuccess: () => setDeleting(null) })}
      />
    </div>
  )
})
