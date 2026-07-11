import { useEffect, useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'

import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/components/ui/alert-dialog'
import { Button } from '@/shared/components/ui/button'
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

  useEffect(() => {
    if (!open) setTyped('')
  }, [open])

  const handleOpenChange = (val: boolean) => {
    if (isLoading) return
    onOpenChange(val)
  }

  const canConfirm = typed === CONFIRM_TEXT && !isLoading

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <div className="flex items-center gap-3">
            <AlertTriangle size={22} className="text-red-500" />
            <AlertDialogTitle>Restore Database</AlertDialogTitle>
          </div>
          <AlertDialogDescription className="pt-1 space-y-2">
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
          </AlertDialogDescription>
        </AlertDialogHeader>

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

        <AlertDialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isLoading}
          >
            Batal
          </Button>
          <Button type="button" variant="destructive" onClick={onConfirm} disabled={!canConfirm}>
            {isLoading && <Loader2 size={14} className="animate-spin" />}
            Ya, Timpa Semua Data
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
