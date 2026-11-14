import { useState } from 'react'
import { Lock, LockOpen, Pencil, Plus, Search, Trash2 } from 'lucide-react'

import { ConfirmDialog, DataTable, FormModal, PageHeader, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'
import type { ColumnDef } from '@/shared/components/DataTable/DataTable.types'
import { useDebounce } from '@/shared/hooks'

import type { MasterApi } from '../menu-items.api'

interface MasterRow {
  id: number
  name: string
  sort_order: number
  is_active: boolean
}

interface SimpleMasterPageProps<T extends MasterRow> {
  title: string
  description: string
  itemLabel: string // mis. "Kategori", "Ukuran"
  menuKey: string
  masterApi: MasterApi<T>
  usage: (row: T) => number
  usageLabel: (count: number) => string
  namePlaceholder: string
}

// Halaman master sederhana (nama + urutan + status) untuk Kategori Menu dan Ukuran.
export function SimpleMasterPage<T extends MasterRow>({
  title,
  description,
  itemLabel,
  menuKey,
  masterApi,
  usage,
  usageLabel,
  namePlaceholder,
}: SimpleMasterPageProps<T>) {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)
  const [isActive, setIsActive] = useState<boolean | undefined>(undefined)
  const [editing, setEditing] = useState<T | 'new' | null>(null)
  const [toggling, setToggling] = useState<T | null>(null)
  const [deleting, setDeleting] = useState<T | null>(null)

  const { data = [], isLoading } = masterApi.useList({ search: debouncedSearch ?? '', is_active: isActive })
  const rows = data
  const { mutate: toggle, isPending: isToggling } = masterApi.useToggle()
  const { mutate: remove, isPending: isDeleting } = masterApi.useDelete()

  const columns: ColumnDef<T>[] = [
    { key: 'name', header: `Nama ${itemLabel}`, mobileLabel: true, cell: (r) => <span className="font-medium text-gray-800">{r.name}</span> },
    { key: 'sort_order', header: 'Urutan', align: 'center', width: '90px', cell: (r) => <span className="text-gray-700">{r.sort_order}</span> },
    { key: 'usage', header: 'Dipakai', cell: (r) => <span className="text-gray-700">{usageLabel(usage(r))}</span> },
    { key: 'is_active', header: 'Status', align: 'center', width: '100px', cell: (r) => <StatusBadge status={r.is_active ? 'active' : 'inactive'} /> },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'center',
      width: '130px',
      cell: (r) => (
        <div className="flex items-center justify-center gap-1">
          <RoleGuard menuKey={menuKey} action="can_edit">
            <ActionButton label="Edit" onClick={() => setEditing(r)}>
              <Pencil size={14} />
            </ActionButton>
            <ActionButton label={r.is_active ? 'Nonaktifkan' : 'Aktifkan'} onClick={() => setToggling(r)}>
              {r.is_active ? <Lock size={14} /> : <LockOpen size={14} />}
            </ActionButton>
          </RoleGuard>
          {usage(r) === 0 && (
            <RoleGuard menuKey={menuKey} action="can_delete">
              <ActionButton label="Hapus" onClick={() => setDeleting(r)} danger>
                <Trash2 size={14} />
              </ActionButton>
            </RoleGuard>
          )}
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <PageHeader
        title={title}
        description={description}
        breadcrumbs={[{ label: 'Menu' }, { label: title }]}
        actions={
          <RoleGuard menuKey={menuKey} action="can_create">
            <Button onClick={() => setEditing('new')} className="gap-1">
              <Plus size={16} />
              Tambah {itemLabel}
            </Button>
          </RoleGuard>
        }
      />

      <div className="flex flex-wrap items-center gap-2 rounded-lg border bg-white p-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            placeholder={`Cari ${itemLabel.toLowerCase()}...`}
            aria-label={`Cari ${itemLabel.toLowerCase()}`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-8 h-9 text-sm"
          />
        </div>
        <ToggleGroup
          type="single"
          value={isActive === undefined ? 'all' : isActive ? 'active' : 'inactive'}
          onValueChange={(v) => v && setIsActive(v === 'all' ? undefined : v === 'active')}
        >
          <ToggleGroupItem value="all">Semua</ToggleGroupItem>
          <ToggleGroupItem value="active">Aktif</ToggleGroupItem>
          <ToggleGroupItem value="inactive">Nonaktif</ToggleGroupItem>
        </ToggleGroup>
      </div>

      <DataTable<T & Record<string, unknown>>
        columns={columns as ColumnDef<T & Record<string, unknown>>[]}
        data={rows as (T & Record<string, unknown>)[]}
        isLoading={isLoading}
        emptyMessage={`Belum ada ${itemLabel.toLowerCase()}`}
        emptyDescription={`Tambah ${itemLabel.toLowerCase()} untuk dipakai di menu.`}
      />

      <MasterFormModal
        key={editing === 'new' ? 'new' : (editing?.id ?? 'closed')}
        open={editing != null}
        onOpenChange={(open) => {
          if (!open) setEditing(null)
        }}
        row={editing === 'new' ? null : editing}
        itemLabel={itemLabel}
        masterApi={masterApi}
        namePlaceholder={namePlaceholder}
      />

      <ConfirmDialog
        open={toggling != null}
        onOpenChange={(open) => {
          if (!open) setToggling(null)
        }}
        title={`${toggling?.is_active ? 'Nonaktifkan' : 'Aktifkan'} ${itemLabel}`}
        description={`${itemLabel} "${toggling?.name}" akan di${toggling?.is_active ? 'nonaktifkan' : 'aktifkan'}. Lanjutkan?`}
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
        title={`Hapus ${itemLabel}`}
        description={`Yakin ingin menghapus ${itemLabel.toLowerCase()} "${deleting?.name}"?`}
        confirmLabel="Ya, Hapus"
        variant="destructive"
        isLoading={isDeleting}
        onConfirm={() => deleting && remove(deleting.id, { onSuccess: () => setDeleting(null) })}
      />
    </div>
  )
}

function ActionButton({ label, onClick, danger, children }: { label: string; onClick: () => void; danger?: boolean; children: React.ReactNode }) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          aria-label={label}
          className={`h-7 w-7 text-gray-500 ${danger ? 'hover:text-red-600' : 'hover:text-blue-600'}`}
          onClick={onClick}
        >
          {children}
        </Button>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}

function MasterFormModal({
  open,
  onOpenChange,
  row,
  itemLabel,
  masterApi,
  namePlaceholder,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  row: MasterRow | null
  itemLabel: string
  masterApi: MasterApi<MasterRow>
  namePlaceholder: string
}) {
  // Nilai awal dari baris yang diedit; komponen dibuat ulang (key) setiap kali form dibuka.
  const [name, setName] = useState(row?.name ?? '')
  const [sortOrder, setSortOrder] = useState(String(row?.sort_order ?? 0))
  const [error, setError] = useState('')
  const { mutate: save, isPending } = masterApi.useSave()

  const submit = () => {
    const trimmed = name.trim()
    const order = Number(sortOrder)
    if (!trimmed) return setError(`Nama ${itemLabel.toLowerCase()} wajib diisi`)
    if (!Number.isInteger(order) || order < 0 || order > 9999) return setError('Urutan harus angka 0 sampai 9999')
    save({ id: row?.id, name: trimmed, sort_order: order }, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={`${row ? 'Edit' : 'Tambah'} ${itemLabel}`}
      size="sm"
      isLoading={isPending}
      onSubmit={submit}
      submitLabel="Simpan"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="master-name">
            Nama {itemLabel} <span className="text-red-500">*</span>
          </Label>
          <Input id="master-name" value={name} onChange={(e) => setName(e.target.value)} placeholder={namePlaceholder} />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="master-order">Urutan Tampil</Label>
          <Input id="master-order" inputMode="numeric" value={sortOrder} onChange={(e) => setSortOrder(e.target.value.replace(/[^\d]/g, ''))} />
          <p className="text-xs text-gray-500">Angka kecil tampil lebih dulu.</p>
        </div>
        {error && <p className="text-xs text-red-500">{error}</p>}
      </div>
    </FormModal>
  )
}
