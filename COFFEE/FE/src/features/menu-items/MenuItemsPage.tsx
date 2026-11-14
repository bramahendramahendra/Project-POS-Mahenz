import { useState } from 'react'
import { Ban, CircleCheck, Lock, LockOpen, Pencil, Plus, RotateCcw, Search, Trash2 } from 'lucide-react'

import { ConfirmDialog, DataTable, PageHeader, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'
import { Switch } from '@/shared/components/ui/switch'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { useDebounce, usePagination, usePageSizeOptions } from '@/shared/hooks'

import { MenuItemFormModal } from './components/MenuItemFormModal'
import {
  menuCategoryApi,
  useDeleteMenuItemMutation,
  useMenuItemListQuery,
  useSetSoldOutMutation,
  useToggleMenuItemMutation,
} from './menu-items.api'
import type { MenuItem, MenuItemListFilter, StockMode } from './menu-items.types'
import { STOCK_MODE_LABEL, formatPriceRange } from './menu-items.utils'

const MENU_KEY = 'menu_jual.daftar'
const initialFilter: MenuItemListFilter = { page: 1, limit: 10, search: '' }

// Daftar Menu (tahap 4.4). Kasir hanya melihat dan menandai habis/tersedia (keputusan 88).
export function MenuItemsPage() {
  const [filter, setFilter] = useState<MenuItemListFilter>(initialFilter)
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)
  const { page, pageSize, onPageChange, onPageSizeChange, reset: resetPage } = usePagination({ initialPageSize: 10 })
  const pageSizeOptions = usePageSizeOptions()

  const [formItemId, setFormItemId] = useState<number | null | undefined>(undefined) // undefined = tertutup, null = baru
  const [toggling, setToggling] = useState<MenuItem | null>(null)
  const [deleting, setDeleting] = useState<MenuItem | null>(null)

  const { data: categories = [] } = menuCategoryApi.useOptions()
  const { data, isLoading } = useMenuItemListQuery({ ...filter, search: debouncedSearch ?? '', page, limit: pageSize })
  const rows = data?.data ?? []
  const total = data?.total ?? 0
  const { mutate: setSoldOut } = useSetSoldOutMutation()
  const { mutate: toggle, isPending: isToggling } = useToggleMenuItemMutation()
  const { mutate: remove, isPending: isDeleting } = useDeleteMenuItemMutation()

  const changeFilter = (patch: Partial<MenuItemListFilter>) => {
    setFilter((prev) => ({ ...prev, ...patch }))
    resetPage()
  }

  const columns: ColumnDef<MenuItem>[] = [
    {
      key: 'name',
      header: 'Nama Menu',
      mobileLabel: true,
      cell: (m) => (
        <div className="flex flex-wrap items-center gap-1.5">
          <span className="font-medium text-gray-800">{m.name}</span>
          {m.is_sold_out && <span className="rounded-full bg-red-100 px-2 py-0.5 text-[11px] font-medium text-red-700">Habis</span>}
        </div>
      ),
    },
    { key: 'category_name', header: 'Kategori', cell: (m) => <span className="text-gray-700">{m.category_name}</span> },
    {
      key: 'stock_mode',
      header: 'Mode',
      cell: (m) => (
        <div>
          <span className="text-gray-700">{STOCK_MODE_LABEL[m.stock_mode]}</span>
          {m.product_name && <p className="text-xs text-gray-500">{m.product_name}</p>}
        </div>
      ),
    },
    { key: 'price', header: 'Harga', align: 'right', cell: (m) => <span className="font-medium text-gray-800">{formatPriceRange(m.price_min, m.price_max)}</span> },
    {
      key: 'variant_count',
      header: 'Varian',
      cell: (m) => <span className="text-gray-700">{m.has_variants ? `${m.variant_count} kombinasi` : 'Tanpa varian'}</span>,
    },
    { key: 'is_active', header: 'Status', align: 'center', width: '100px', cell: (m) => <StatusBadge status={m.is_active ? 'active' : 'inactive'} /> },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '170px',
      cell: (m) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey={MENU_KEY} action="can_view">
            <IconBtn
              label={m.is_sold_out ? 'Tandai tersedia' : 'Tandai habis'}
              onClick={() => setSoldOut({ id: m.id, isSoldOut: !m.is_sold_out })}
              className={m.is_sold_out ? 'hover:text-green-600' : 'hover:text-red-600'}
            >
              {m.is_sold_out ? <CircleCheck size={14} /> : <Ban size={14} />}
            </IconBtn>
          </RoleGuard>
          <RoleGuard menuKey={MENU_KEY} action="can_edit">
            <IconBtn label="Edit" onClick={() => setFormItemId(m.id)}>
              <Pencil size={14} />
            </IconBtn>
            <IconBtn label={m.is_active ? 'Nonaktifkan' : 'Aktifkan'} onClick={() => setToggling(m)}>
              {m.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
            </IconBtn>
          </RoleGuard>
          <RoleGuard menuKey={MENU_KEY} action="can_delete">
            <IconBtn label="Hapus" onClick={() => setDeleting(m)} className="hover:text-red-600">
              <Trash2 size={14} />
            </IconBtn>
          </RoleGuard>
        </div>
      ),
    },
  ]

  const hasFilter = !!search || !!filter.category_id || !!filter.stock_mode || filter.is_active !== undefined || !!filter.is_sold_out

  return (
    <div className="space-y-4">
      <PageHeader
        title="Daftar Menu"
        description="Menu yang dijual di kasir: harga, varian, resep, dan kelompok tambahan."
        breadcrumbs={[{ label: 'Menu' }, { label: 'Daftar Menu' }]}
        actions={
          <RoleGuard menuKey={MENU_KEY} action="can_create">
            <Button onClick={() => setFormItemId(null)} className="gap-1">
              <Plus size={16} />
              Tambah Menu
            </Button>
          </RoleGuard>
        }
      />

      <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[200px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input placeholder="Cari menu..." aria-label="Cari menu" value={search}
            onChange={(e) => {
              setSearch(e.target.value)
              resetPage()
            }}
            className="pl-8 h-9 text-sm"
          />
        </div>
        <Select value={filter.category_id ? String(filter.category_id) : 'all'} onValueChange={(v) => changeFilter({ category_id: v === 'all' ? undefined : Number(v) })}>
          <SelectTrigger className="h-9 w-44" aria-label="Filter kategori">
            <SelectValue placeholder="Semua kategori" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua kategori</SelectItem>
            {categories.map((c) => (
              <SelectItem key={c.id} value={String(c.id)}>
                {c.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={filter.stock_mode ?? 'all'} onValueChange={(v) => changeFilter({ stock_mode: v === 'all' ? undefined : (v as StockMode) })}>
          <SelectTrigger className="h-9 w-40" aria-label="Filter mode stok">
            <SelectValue placeholder="Semua mode" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua mode</SelectItem>
            <SelectItem value="recipe">{STOCK_MODE_LABEL.recipe}</SelectItem>
            <SelectItem value="direct">{STOCK_MODE_LABEL.direct}</SelectItem>
          </SelectContent>
        </Select>
        <label className="flex items-center gap-2 text-sm text-gray-700">
          <Switch checked={!!filter.is_sold_out} onCheckedChange={(v) => changeFilter({ is_sold_out: v || undefined })} aria-label="Hanya menu habis" />
          Habis
        </label>
        <Button
          variant="outline"
          size="sm"
          className="h-9 gap-1"
          onClick={() => {
            setSearch('')
            setFilter(initialFilter)
            resetPage()
          }}
        >
          <RotateCcw size={13} />
          Reset
        </Button>
      </div>

      <DataTable<MenuItem & Record<string, unknown>>
        columns={columns as ColumnDef<MenuItem & Record<string, unknown>>[]}
        data={rows as (MenuItem & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={hasFilter ? 'Menu tidak ditemukan' : 'Belum ada menu'}
        emptyDescription={hasFilter ? 'Coba ubah kata kunci atau filter.' : 'Buat kategori dan ukuran dulu, lalu tambah menu.'}
        pagination={{ page, pageSize, total, onPageChange, onPageSizeChange, pageSizeOptions }}
      />

      <MenuItemFormModal
        open={formItemId !== undefined}
        onOpenChange={(open) => {
          if (!open) setFormItemId(undefined)
        }}
        itemId={formItemId ?? null}
      />

      <ConfirmDialog
        open={toggling != null}
        onOpenChange={(open) => {
          if (!open) setToggling(null)
        }}
        title={toggling?.is_active ? 'Nonaktifkan Menu' : 'Aktifkan Menu'}
        description={`Menu "${toggling?.name}" akan di${toggling?.is_active ? 'nonaktifkan dan tidak tampil di kasir' : 'aktifkan'}. Lanjutkan?`}
        confirmLabel={toggling?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={toggling?.is_active ? 'destructive' : 'default'}
        isLoading={isToggling}
        onConfirm={() => toggling && toggle(toggling.id, { onSuccess: () => setToggling(null) })}
      />

      <ConfirmDialog
        open={deleting != null}
        onOpenChange={(open) => {
          if (!open) setDeleting(null)
        }}
        title="Hapus Menu"
        description={`Yakin ingin menghapus menu "${deleting?.name}" beserta varian dan resepnya?`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={() => deleting && remove(deleting.id, { onSuccess: () => setDeleting(null) })}
      />
    </div>
  )
}

function IconBtn({ label, onClick, className, children }: { label: string; onClick: () => void; className?: string; children: React.ReactNode }) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button variant="ghost" size="icon" aria-label={label} className={`h-7 w-7 text-gray-500 ${className ?? 'hover:text-blue-600'}`} onClick={onClick}>
          {children}
        </Button>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}
