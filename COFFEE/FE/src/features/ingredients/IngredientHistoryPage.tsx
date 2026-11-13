import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'

import { DataTable, PageHeader } from '@/shared/components'
import { usePagination, usePageSizeOptions } from '@/shared/hooks'

import { CancelMutationModal } from './components/CancelMutationModal'
import { MutationFilterBar } from './components/MutationFilterBar'
import { buildMutationColumns } from './components/MutationTableColumns'
import { useMutationListQuery } from './ingredients.api'
import type { IngredientMutation, MutationListFilter } from './ingredients.types'

// Riwayat Stok Bahan: siapa, kapan, berapa, catatan. Entri terpakai/terbuang bisa dibatalkan
// oleh admin/owner (keputusan 81).
export function IngredientHistoryPage() {
  const [searchParams] = useSearchParams()
  const initialIngredient = Number(searchParams.get('ingredient_id')) || undefined
  const [filter, setFilter] = useState<MutationListFilter>({ page: 1, limit: 20, ingredient_id: initialIngredient })
  const [cancelling, setCancelling] = useState<IngredientMutation | null>(null)

  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 20 })
  const pageSizeOptions = usePageSizeOptions()

  const { data, isLoading } = useMutationListQuery({ ...filter, page, limit: pageSize })
  const rows = data?.data ?? []
  const total = data?.total ?? 0

  const hasFilter = !!filter.ingredient_id || !!filter.type || !!filter.start_date || !!filter.end_date

  return (
    <div className="space-y-4">
      <PageHeader
        title="Riwayat Stok Bahan"
        description="Catatan terpakai, terbuang, dan pembatalan. Entri tidak bisa diubah, salah input dibatalkan oleh admin."
        breadcrumbs={[{ label: 'Bahan Baku' }, { label: 'Riwayat Stok Bahan' }]}
      />
      <MutationFilterBar
        filter={filter}
        onChange={(next) => {
          setFilter(next)
          resetPage()
        }}
        onReset={() => {
          setFilter({ page: 1, limit: 20 })
          resetPage()
        }}
      />
      <DataTable<IngredientMutation & Record<string, unknown>>
        columns={buildMutationColumns(setCancelling)}
        data={rows as (IngredientMutation & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={hasFilter ? 'Riwayat tidak ditemukan' : 'Belum ada riwayat stok bahan'}
        emptyDescription={hasFilter ? 'Coba ubah filter.' : 'Riwayat muncul setelah pemakaian atau barang terbuang dicatat.'}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />
      <CancelMutationModal
        mutation={cancelling}
        onOpenChange={(open) => {
          if (!open) setCancelling(null)
        }}
      />
    </div>
  )
}
