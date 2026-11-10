import { useState } from 'react'
import { AlertCircle, Loader2 } from 'lucide-react'

import { FormModal } from '@/shared/components'
import { Checkbox } from '@/shared/components/ui/checkbox'

import { useSetUserAccessMutation, useUserAccessQuery } from '../users.api'
import type { AdminUser } from '../users.types'

interface UserAccessModalProps {
  user: AdminUser | null
  onClose: () => void
}

export function UserAccessModal({ user, onClose }: UserAccessModalProps) {
  const { data, isLoading, isError, error } = useUserAccessQuery(user?.id ?? null)
  const { mutate, isPending, error: saveError, reset } = useSetUserAccessMutation()
  // Perubahan lokal admin; null = belum diubah, pakai data dari server.
  // Komponen di-remount per user (lihat `key` di halaman), jadi tidak perlu reset manual.
  const [edited, setEdited] = useState<Set<number> | null>(null)
  const selected = edited ?? new Set((data ?? []).filter((a) => a.granted).map((a) => a.app_id))

  const toggle = (appId: number) => {
    const next = new Set(selected)
    if (next.has(appId)) next.delete(appId)
    else next.add(appId)
    setEdited(next)
    reset()
  }

  const handleSubmit = () => {
    if (!user) return
    mutate({ userId: user.id, appIds: [...selected] }, { onSuccess: onClose })
  }

  return (
    <FormModal
      open={!!user}
      onOpenChange={(v) => !v && onClose()}
      title="Akses aplikasi"
      description={`Pilih aplikasi yang boleh dibuka oleh ${user?.full_name ?? ''} (@${user?.username ?? ''}).`}
      size="sm"
      isLoading={isPending}
      submitDisabled={isLoading || isError}
      onSubmit={handleSubmit}
      submitLabel="Simpan akses"
    >
      {isLoading && (
        <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
          <Loader2 size={16} className="animate-spin" /> Memuat...
        </div>
      )}

      {(isError || saveError) && (
        <div
          role="alert"
          className="mb-3 flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700"
        >
          <AlertCircle size={16} className="mt-0.5 shrink-0" />
          <span>{(error ?? saveError)?.message}</span>
        </div>
      )}

      {data && data.length === 0 && (
        <p className="py-6 text-center text-sm text-muted-foreground">
          Belum ada aplikasi. Tambahkan dulu di menu Kelola Aplikasi.
        </p>
      )}

      {data && data.length > 0 && (
        <ul className="space-y-2">
          {data.map((app) => {
            const id = `access-${app.app_id}`
            return (
              <li key={app.app_id}>
                <label
                  htmlFor={id}
                  className="flex cursor-pointer items-center gap-3 rounded-xl border border-border px-4 py-3 transition-colors hover:bg-secondary/60 has-[[data-state=checked]]:border-brand/40 has-[[data-state=checked]]:bg-accent/60"
                >
                  <Checkbox
                    id={id}
                    checked={selected.has(app.app_id)}
                    onCheckedChange={() => toggle(app.app_id)}
                  />
                  <span className="min-w-0 flex-1">
                    <span className="block text-sm font-medium text-foreground">{app.name}</span>
                    <span className="block font-mono text-xs text-muted-foreground">{app.code}</span>
                  </span>
                  {!app.is_active && (
                    <span className="rounded-full bg-secondary px-2 py-0.5 text-[11px] text-muted-foreground">
                      Nonaktif
                    </span>
                  )}
                </label>
              </li>
            )
          })}
        </ul>
      )}
    </FormModal>
  )
}
