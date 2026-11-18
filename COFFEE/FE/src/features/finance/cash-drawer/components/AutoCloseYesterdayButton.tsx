import { useState } from 'react'
import { Lock } from 'lucide-react'

import { ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import { useAutoCloseYesterdayMutation } from '../cash-drawer.api'

/**
 * Tombol admin untuk menutup semua kas hari sebelumnya yang masih terbuka (keputusan 133), sama seperti penutupan
 * otomatis tengah malam. Berguna bila server mati saat tengah malam, atau kas tertunda karena pesanan terbuka yang
 * kini sudah selesai. Untuk mencatat selisih uang per kas tetap memakai Tutup Paksa.
 */
export function AutoCloseYesterdayButton() {
  const [open, setOpen] = useState(false)
  const { mutate, isPending } = useAutoCloseYesterdayMutation()

  return (
    <>
      <Button variant="outline" onClick={() => setOpen(true)} className="gap-1">
        <Lock size={16} />
        Tutup Kas Hari Sebelumnya
      </Button>
      <ConfirmDialog
        open={open}
        onOpenChange={setOpen}
        title="Tutup kas hari sebelumnya?"
        description="Semua kas dari hari sebelumnya yang masih terbuka akan ditutup. Saldo akhir dianggap sama dengan saldo seharusnya (selisih Rp 0). Kas pada hari yang masih punya pesanan terbuka tidak ikut ditutup."
        confirmLabel="Ya, tutup"
        isLoading={isPending}
        onConfirm={() => mutate(undefined, { onSuccess: () => setOpen(false) })}
      />
    </>
  )
}
