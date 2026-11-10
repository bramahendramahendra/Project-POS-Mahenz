import { useState } from 'react'
import { KeyRound, MoreHorizontal, Pencil, Plus, Power } from 'lucide-react'

import { ConfirmDialog, DataTable, PageHeader, StatusBadge } from '@/shared/components'
import type { ColumnDef } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'

import { useAdminAppsQuery, useRotateKeyMutation, useToggleAppMutation } from './apps.api'
import type { AdminApp } from './apps.types'
import { AppFormModal } from './components/AppFormModal'
import { AppKeyModal } from './components/AppKeyModal'

interface KeyInfo {
  key: string
  appName: string
  isRotation: boolean
}

export function AdminAppsPage() {
  const { data, isLoading } = useAdminAppsQuery()
  const toggle = useToggleAppMutation()
  const rotate = useRotateKeyMutation()

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<AdminApp | null>(null)
  const [keyInfo, setKeyInfo] = useState<KeyInfo | null>(null)
  const [toggling, setToggling] = useState<AdminApp | null>(null)
  const [rotating, setRotating] = useState<AdminApp | null>(null)

  const openCreate = () => {
    setEditing(null)
    setFormOpen(true)
  }
  const openEdit = (app: AdminApp) => {
    setEditing(app)
    setFormOpen(true)
  }

  const columns: ColumnDef<AdminApp>[] = [
    {
      key: 'name',
      header: 'Aplikasi',
      mobileLabel: true,
      cell: (row) => (
        <div>
          <p className="font-medium text-foreground">{row.name}</p>
          <p className="font-mono text-xs text-muted-foreground">{row.code}</p>
        </div>
      ),
    },
    {
      key: 'path',
      header: 'Alamat',
      cell: (row) => <span className="font-mono text-xs text-muted-foreground">{row.path}</span>,
    },
    {
      key: 'user_count',
      header: 'Pengguna',
      align: 'center',
      cell: (row) => <span className="tabular-nums">{row.user_count}</span>,
    },
    {
      key: 'has_key',
      header: 'Kunci',
      cell: (row) =>
        row.has_key ? (
          <StatusBadge status="success" label="Terpasang" size="sm" />
        ) : (
          <StatusBadge status="warning" label="Belum ada" size="sm" />
        ),
    },
    {
      key: 'is_active',
      header: 'Status',
      cell: (row) => <StatusBadge status={row.is_active ? 'active' : 'inactive'} />,
    },
    {
      key: 'actions',
      header: '',
      width: '56px',
      align: 'right',
      cell: (row) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" aria-label={`Aksi untuk ${row.name}`}>
              <MoreHorizontal size={18} />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-52">
            <DropdownMenuItem onSelect={() => openEdit(row)}>
              <Pencil size={16} className="mr-2" /> Ubah
            </DropdownMenuItem>
            <DropdownMenuItem onSelect={() => setRotating(row)}>
              <KeyRound size={16} className="mr-2" /> Buat kunci baru
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onSelect={() => setToggling(row)}>
              <Power size={16} className="mr-2" /> {row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ]

  return (
    <div>
      <PageHeader
        title="Kelola Aplikasi"
        description="Daftar aplikasi yang bisa dibuka dari Portal beserta kunci SSO-nya."
        actions={
          <Button onClick={openCreate} className="gap-2">
            <Plus size={16} /> Tambah Aplikasi
          </Button>
        }
      />

      <DataTable<AdminApp>
        columns={columns}
        data={data ?? []}
        isLoading={isLoading}
        emptyMessage="Belum ada aplikasi"
        emptyDescription="Tambahkan aplikasi pertama agar bisa dibuka dari Portal."
      />

      <AppFormModal
        open={formOpen}
        onOpenChange={setFormOpen}
        app={editing}
        onCreated={(result) =>
          setKeyInfo({ key: result.app_key, appName: result.app.name, isRotation: false })
        }
      />

      <AppKeyModal
        appKey={keyInfo?.key ?? null}
        appName={keyInfo?.appName ?? ''}
        isRotation={keyInfo?.isRotation}
        onClose={() => setKeyInfo(null)}
      />

      <ConfirmDialog
        open={!!toggling}
        onOpenChange={(v) => !v && setToggling(null)}
        title={toggling?.is_active ? 'Nonaktifkan aplikasi?' : 'Aktifkan aplikasi?'}
        description={
          toggling?.is_active
            ? `${toggling.name} tidak akan muncul di beranda user dan login SSO ke aplikasi ini ditolak.`
            : `${toggling?.name ?? ''} akan kembali muncul bagi user yang punya akses.`
        }
        confirmLabel={toggling?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={toggling?.is_active ? 'destructive' : 'default'}
        isLoading={toggle.isPending}
        onConfirm={() => toggling && toggle.mutate(toggling.id, { onSuccess: () => setToggling(null) })}
      />

      <ConfirmDialog
        open={!!rotating}
        onOpenChange={(v) => !v && setRotating(null)}
        title="Buat kunci baru?"
        description={`Kunci lama ${rotating?.name ?? ''} langsung tidak berlaku. Login SSO ke aplikasi ini gagal sampai kunci baru dipasang di aplikasi tujuan.`}
        confirmLabel="Buat kunci baru"
        variant="destructive"
        isLoading={rotate.isPending}
        onConfirm={() =>
          rotating &&
          rotate.mutate(rotating.id, {
            onSuccess: (res) => {
              setKeyInfo({ key: res.app_key, appName: rotating.name, isRotation: true })
              setRotating(null)
            },
          })
        }
      />
    </div>
  )
}
