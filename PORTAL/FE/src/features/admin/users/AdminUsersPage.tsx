import { useState } from 'react'
import {
  AppWindow,
  KeyRound,
  MoreHorizontal,
  Pencil,
  Plus,
  Power,
  Search,
  Trash2,
} from 'lucide-react'

import { useAuth } from '@/features/auth/hooks/useAuth'
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
import { Input } from '@/shared/components/ui/input'
import { useDebounce } from '@/shared/hooks'

import { UserAccessModal } from './components/UserAccessModal'
import { ResetPasswordModal } from './components/ResetPasswordModal'
import { UserFormModal } from './components/UserFormModal'
import { useDeleteUserMutation, useToggleUserMutation, useUsersQuery } from './users.api'
import type { AdminUser } from './users.types'

const PAGE_SIZE_OPTIONS = [10, 20, 50]

export function AdminUsersPage() {
  const { user: me } = useAuth()

  const [limit, setLimit] = useState(10)
  const [searchInput, setSearchInput] = useState('')
  const search = useDebounce(searchInput.trim(), 350)

  // Halaman kembali ke 1 otomatis saat kata kunci / ukuran halaman berubah (tanpa effect).
  const queryKey = `${search}|${limit}`
  const [pageInfo, setPageInfo] = useState({ page: 1, key: queryKey })
  const page = pageInfo.key === queryKey ? pageInfo.page : 1
  const setPage = (p: number) => setPageInfo({ page: p, key: queryKey })

  const { data, isLoading } = useUsersQuery({ page, limit, search: search || undefined })
  const toggleUser = useToggleUserMutation()
  const deleteUser = useDeleteUserMutation()

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<AdminUser | null>(null)
  const [resetting, setResetting] = useState<AdminUser | null>(null)
  const [accessUser, setAccessUser] = useState<AdminUser | null>(null)
  const [toggling, setToggling] = useState<AdminUser | null>(null)
  const [deleting, setDeleting] = useState<AdminUser | null>(null)

  const openCreate = () => {
    setEditing(null)
    setFormOpen(true)
  }
  const openEdit = (u: AdminUser) => {
    setEditing(u)
    setFormOpen(true)
  }

  const columns: ColumnDef<AdminUser>[] = [
    {
      key: 'full_name',
      header: 'Pengguna',
      mobileLabel: true,
      cell: (row) => (
        <div>
          <p className="font-medium text-foreground">
            {row.full_name}
            {row.id === me?.id && (
              <span className="ml-2 rounded-full bg-accent px-2 py-0.5 text-[11px] font-medium text-accent-foreground">
                Anda
              </span>
            )}
          </p>
          <p className="text-xs text-muted-foreground">@{row.username}</p>
        </div>
      ),
    },
    {
      key: 'role_name',
      header: 'Role',
      cell: (row) => (
        <span className="rounded-full border border-border bg-secondary px-2.5 py-0.5 text-xs font-medium capitalize">
          {row.role_name}
        </span>
      ),
    },
    {
      key: 'is_active',
      header: 'Status',
      cell: (row) => (
        <div className="flex flex-wrap items-center gap-1.5">
          <StatusBadge status={row.is_active ? 'active' : 'inactive'} />
          {row.must_change_password && (
            <StatusBadge status="warning" label="Wajib ganti password" size="sm" />
          )}
        </div>
      ),
    },
    {
      key: 'actions',
      header: '',
      width: '56px',
      align: 'right',
      cell: (row) => {
        const isSelf = row.id === me?.id
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" aria-label={`Aksi untuk ${row.username}`}>
                <MoreHorizontal size={18} />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem onSelect={() => openEdit(row)}>
                <Pencil size={16} className="mr-2" /> Ubah
              </DropdownMenuItem>
              <DropdownMenuItem onSelect={() => setAccessUser(row)}>
                <AppWindow size={16} className="mr-2" /> Akses aplikasi
              </DropdownMenuItem>
              {!isSelf && (
                <DropdownMenuItem onSelect={() => setResetting(row)}>
                  <KeyRound size={16} className="mr-2" /> Reset password
                </DropdownMenuItem>
              )}
              {!isSelf && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onSelect={() => setToggling(row)}>
                    <Power size={16} className="mr-2" /> {row.is_active ? 'Nonaktifkan' : 'Aktifkan'}
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onSelect={() => setDeleting(row)}
                    className="text-destructive focus:text-destructive"
                  >
                    <Trash2 size={16} className="mr-2" /> Hapus
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]

  return (
    <div>
      <PageHeader
        title="Kelola User"
        description="Akun yang dapat login ke Portal dan aplikasi yang boleh dibuka."
        actions={
          <Button onClick={openCreate} className="gap-2">
            <Plus size={16} /> Tambah User
          </Button>
        }
      />

      <div className="relative mb-4 max-w-sm">
        <Search
          size={16}
          className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          placeholder="Cari nama atau username"
          aria-label="Cari user"
          className="h-10 rounded-xl bg-white pl-10"
        />
      </div>

      <DataTable<AdminUser>
        columns={columns}
        data={data?.data ?? []}
        isLoading={isLoading}
        emptyMessage={search ? 'User tidak ditemukan' : 'Belum ada user'}
        emptyDescription={search ? 'Coba kata kunci lain.' : undefined}
        pagination={{
          page,
          pageSize: limit,
          total: data?.total ?? 0,
          onPageChange: setPage,
          pageSizeOptions: PAGE_SIZE_OPTIONS,
          onPageSizeChange: setLimit,
        }}
      />

      <UserFormModal
        key={editing?.id ?? 'new'}
        open={formOpen}
        onOpenChange={setFormOpen}
        user={editing}
      />
      <ResetPasswordModal user={resetting} onClose={() => setResetting(null)} />
      <UserAccessModal
        key={accessUser?.id ?? 'none'}
        user={accessUser}
        onClose={() => setAccessUser(null)}
      />

      <ConfirmDialog
        open={!!toggling}
        onOpenChange={(v) => !v && setToggling(null)}
        title={toggling?.is_active ? 'Nonaktifkan user?' : 'Aktifkan user?'}
        description={
          toggling?.is_active
            ? `@${toggling.username} tidak dapat login lagi sampai diaktifkan kembali.`
            : `@${toggling?.username ?? ''} dapat login kembali.`
        }
        confirmLabel={toggling?.is_active ? 'Nonaktifkan' : 'Aktifkan'}
        variant={toggling?.is_active ? 'destructive' : 'default'}
        isLoading={toggleUser.isPending}
        onConfirm={() => toggling && toggleUser.mutate(toggling.id, { onSuccess: () => setToggling(null) })}
      />

      <ConfirmDialog
        open={!!deleting}
        onOpenChange={(v) => !v && setDeleting(null)}
        title="Hapus user?"
        description={`@${deleting?.username ?? ''} akan dihapus permanen beserta akses aplikasinya. Tindakan ini tidak dapat dibatalkan.`}
        confirmLabel="Hapus"
        variant="destructive"
        isLoading={deleteUser.isPending}
        onConfirm={() => deleting && deleteUser.mutate(deleting.id, { onSuccess: () => setDeleting(null) })}
      />
    </div>
  )
}
