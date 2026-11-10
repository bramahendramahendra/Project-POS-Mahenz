import { useState } from 'react'
import { KeyRound } from 'lucide-react'

import { useLoginOptionsQuery, useSetDirectLoginMutation } from '@/features/auth/auth.api'
import { useAuth } from '@/features/auth/hooks/useAuth'
import { ConfirmDialog } from '@/shared/components'
import { Switch } from '@/shared/components/ui/switch'

// Pengaturan login langsung (username + password di aplikasi ini). Hanya terlihat dan
// hanya bisa diubah oleh role admin; BE menegakkannya juga (role lain mendapat 403).
export function DirectLoginCard() {
  const { isAdmin } = useAuth()
  const { data: options, isLoading } = useLoginOptionsQuery()
  const { mutate, isPending } = useSetDirectLoginMutation()
  const [confirmOff, setConfirmOff] = useState(false)

  if (!isAdmin || isLoading || !options) return null

  const enabled = options.direct_login_enabled
  const portalReady = options.portal_login_enabled
  const forced = options.direct_login_forced

  const handleChange = (next: boolean) => {
    if (next) {
      mutate(true)
    } else {
      setConfirmOff(true)
    }
  }

  let hint = 'Bila dimatikan, semua orang masuk lewat Portal dan form username + password disembunyikan.'
  if (!portalReady) {
    hint = 'Login lewat Portal belum diatur di server ini, sehingga login langsung tidak dapat dimatikan.'
  } else if (forced) {
    hint = 'Dipaksa aktif oleh pengaturan darurat di server. Setelan ini baru berlaku setelah pengaturan darurat dimatikan.'
  }

  return (
    <div className="flex items-start justify-between gap-4 rounded-lg border bg-white p-4">
      <div className="flex gap-3">
        <KeyRound size={20} className="mt-0.5 text-gray-500" />
        <div>
          <p className="text-sm font-medium text-gray-800">Login langsung (username dan password)</p>
          <p className="mt-0.5 text-sm text-gray-500">{hint}</p>
        </div>
      </div>
      <Switch
        checked={enabled}
        onCheckedChange={handleChange}
        disabled={isPending || (!portalReady && enabled)}
        aria-label="Login langsung"
      />

      <ConfirmDialog
        open={confirmOff}
        onOpenChange={setConfirmOff}
        title="Matikan login langsung?"
        description="Form username dan password akan disembunyikan dan login langsung ditolak. Pastikan Portal berfungsi dan akun Anda bisa masuk lewat Portal. Anda bisa menyalakannya lagi dari sini."
        confirmLabel="Matikan"
        variant="destructive"
        isLoading={isPending}
        onConfirm={() => mutate(false, { onSuccess: () => setConfirmOff(false) })}
      />
    </div>
  )
}
