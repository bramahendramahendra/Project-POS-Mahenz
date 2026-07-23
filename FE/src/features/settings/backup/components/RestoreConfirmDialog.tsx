import { useState } from 'react'

import { ExtendedConfirmDialog } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

const CONFIRM_TEXT = 'RESTORE'

interface RestoreConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  filename: string | null
  createdAt: string | null
  isLoading?: boolean
  onConfirm: () => void
}

export function RestoreConfirmDialog({
  open,
  onOpenChange,
  filename,
  createdAt,
  isLoading,
  onConfirm,
}: RestoreConfirmDialogProps) {
  const [typed, setTyped] = useState('')

  const handleOpenChange = (val: boolean) => {
    if (!val) setTyped('')
    onOpenChange(val)
  }

  const canConfirm = typed === CONFIRM_TEXT && !isLoading

  return (
    <ExtendedConfirmDialog
      open={open}
      onOpenChange={handleOpenChange}
      title="Restore Database"
      variant="destructive"
      isLoading={isLoading}
      confirmLabel="Ya, Timpa Semua Data"
      confirmDisabled={!canConfirm}
      onConfirm={onConfirm}
      description={
        <span className="block space-y-2">
          <span className="block">
            Tindakan ini akan <strong className="text-red-600">menimpa seluruh data</strong> yang
            ada saat ini dengan isi file backup berikut. Semua transaksi, produk, dan data lain
            yang dibuat setelah tanggal backup ini akan{' '}
            <strong className="text-red-600">hilang permanen</strong> dan tidak bisa dikembalikan.
          </span>
          {filename && (
            <span className="block rounded-md bg-gray-50 border px-3 py-2 text-xs font-mono text-gray-700">
              {filename}
              {createdAt && (
                <span className="block text-gray-400 mt-0.5">
                  Dibuat: {new Date(createdAt).toLocaleString('id-ID')}
                </span>
              )}
            </span>
          )}
        </span>
      }
    >
      <div className="space-y-1.5">
        <Label htmlFor="restore-confirm-text" className="text-xs">
          Ketik <span className="font-mono font-semibold">{CONFIRM_TEXT}</span> untuk konfirmasi
        </Label>
        <Input
          id="restore-confirm-text"
          value={typed}
          onChange={(e) => setTyped(e.target.value)}
          placeholder={CONFIRM_TEXT}
          autoComplete="off"
          disabled={isLoading}
        />
      </div>
    </ExtendedConfirmDialog>
  )
}
