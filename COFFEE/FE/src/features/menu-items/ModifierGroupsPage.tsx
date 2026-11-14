import { useState } from 'react'
import { Lock, LockOpen, Pencil, Plus, Search, Trash2 } from 'lucide-react'

import { ConfirmDialog, DataTable, PageHeader, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { useDebounce } from '@/shared/hooks'
import { formatRupiah } from '@/shared/utils'

import { ModifierGroupFormModal } from './components/ModifierGroupFormModal'
import {
  useDeleteModifierGroupMutation,
  useModifierGroupListQuery,
  useToggleModifierGroupMutation,
} from './menu-items.api'
import type { ModifierGroup } from './menu-items.types'

const MENU_KEY = 'menu_jual.tambahan'

// Kelompok tambahan (keputusan 54, 87), dipasang ke menu lewat form menu.
export function ModifierGroupsPage() {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)
  const [editing, setEditing] = useState<ModifierGroup | 'new' | null>(null)
  const [toggling, setToggling] = useState<ModifierGroup | null>(null)
  const [deleting, setDeleting] = useState<ModifierGroup | null>(null)

  const { data = [], isLoading } = useModifierGroupListQuery({ search: debouncedSearch ?? '' })
  const { mutate: toggle, isPending: isToggling } = useToggleModifierGroupMutation()
  const { mutate: remove, isPending: isDeleting } = useDeleteModifierGroupMutation()

  const columns: ColumnDef<ModifierGroup>[] = [
    { key: 'name', header: 'Nama Kelompok', mobileLabel: true, cell: (g) => <span className="font-medium text-gray-800">{g.name}</span> },
    {
      key: 'selection_type',
      header: 'Aturan',
      cell: (g) => (
        <div className="text-sm text-gray-700">
          <p>{g.selection_type === 'single' ? 'Pilih satu' : 'Boleh banyak'}</p>
          <p className="text-xs text-gray-500">{g.is_required ? 'Wajib dipilih' : 'Opsional'}</p>
        </div>
      ),
    },
    {
      key: 'options',
      header: 'Pilihan',
      cell: (g) => (
        <div className="flex max-w-[360px] flex-wrap gap-1">
          {g.options.map((o) => (
            <span
              key={o.id}
              className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-xs ${o.is_active ? 'text-gray-700' : 'text-gray-400 line-through'} ${o.is_default ? 'border-blue-300 bg-blue-50' : ''}`}
            >
              {o.name}
              {o.price > 0 && <span className="text-gray-500">+{formatRupiah(o.price)}</span>}
              {o.is_default && <span className="text-blue-600">(bawaan)</span>}
            </span>
          ))}
        </div>
      ),
    },
    { key: 'menu_count', header: 'Dipakai', cell: (g) => <span className="text-gray-700">{g.menu_count > 0 ? `${g.menu_count} menu` : 'Belum dipakai'}</span> },
    { key: 'is_active', header: 'Status', align: 'center', width: '100px', cell: (g) => <StatusBadge status={g.is_active ? 'active' : 'inactive'} /> },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '130px',
      cell: (g) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey={MENU_KEY} action="can_edit">
            <IconBtn label="Edit" onClick={() => setEditing(g)}>
              <Pencil size={14} />
            </IconBtn>
            <IconBtn label={g.is_active ? 'Nonaktifkan' : 'Aktifkan'} onClick={() => setToggling(g)}>
              {g.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
            </IconBtn>
          </RoleGuard>
          {g.menu_count === 0 && (
            <RoleGuard menuKey={MENU_KEY} action="can_delete">
              <IconBtn label="Hapus" onClick={() => setDeleting(g)}>
                <Trash2 size={14} />
              </IconBtn>
            </RoleGuard>
          )}
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kelompok Tambahan"
        description="Pilihan tambahan untuk menu, mis. level gula (gratis) atau topping (berbayar)."
        breadcrumbs={[{ label: 'Menu' }, { label: 'Kelompok Tambahan' }]}
        actions={
          <RoleGuard menuKey={MENU_KEY} action="can_create">
            <Button onClick={() => setEditing('new')} className="gap-1">
              <Plus size={16} />
              Tambah Kelompok
            </Button>
          </RoleGuard>
        }
      />

      <div className="flex items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input placeholder="Cari kelompok..." aria-label="Cari kelompok" value={search} onChange={(e) => setSearch(e.target.value)} className="pl-8 h-9 text-sm" />
        </div>
      </div>

      <DataTable<ModifierGroup & Record<string, unknown>>
        columns={columns as ColumnDef<ModifierGroup & Record<string, unknown>>[]}
        data={data as (ModifierGroup & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage="Belum ada kelompok tambahan"
        emptyDescription="Tambah kelompok seperti Level Gula atau Topping."
      />

      <ModifierGroupFormModal
        key={editing === 'new' ? 'new' : (editing?.id ?? 'closed')}
        open={editing != null}
        onOpenChange={(open) => {
          if (!open) setEditing(null)
        }}
        group={editing === 'new' ? null : editing}
      />

      <ConfirmDialog
        open={toggling != null}
        onOpenChange={(open) => {
          if (!open) setToggling(null)
        }}
        title={toggling?.is_active ? 'Nonaktifkan Kelompok' : 'Aktifkan Kelompok'}
        description={`Kelompok "${toggling?.name}" akan di${toggling?.is_active ? 'nonaktifkan' : 'aktifkan'}. Lanjutkan?`}
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
        title="Hapus Kelompok"
        description={`Yakin ingin menghapus kelompok "${deleting?.name}" beserta pilihannya?`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={() => deleting && remove(deleting.id, { onSuccess: () => setDeleting(null) })}
      />
    </div>
  )
}

function IconBtn({ label, onClick, children }: { label: string; onClick: () => void; children: React.ReactNode }) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button variant="ghost" size="icon" aria-label={label} className="h-7 w-7 text-gray-500 hover:text-blue-600" onClick={onClick}>
          {children}
        </Button>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}
