import { useState } from 'react'
import { CheckCircle2, Trash2, TriangleAlert } from 'lucide-react'

import { ConfirmDialog, ExtendedConfirmDialog, FormModal, RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Textarea } from '@/shared/components/ui/textarea'

import {
  useConfirmExpiryBatchMutation,
  useExpiryWarningsQuery,
  useWriteOffExpiryBatchMutation,
} from '../expiry-batches.api'
import type { ExpiryWarning } from '../expiry-batches.types'
import type { Product } from '../products.types'

interface ExpiryWarningModalProps {
  product: Product | null
  onOpenChange: (open: boolean) => void
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' })
}

function daysLeftLabel(warning: ExpiryWarning): string {
  if (warning.severity === 'expired') return `Lewat ${Math.abs(warning.days_left)} hari`
  if (warning.days_left === 0) return 'Expired hari ini'
  return `${warning.days_left} hari lagi`
}

export function ExpiryWarningModal({ product, onOpenChange }: ExpiryWarningModalProps) {
  const open = product !== null
  const [writeOffTarget, setWriteOffTarget] = useState<ExpiryWarning | null>(null)
  const [writeOffNotes, setWriteOffNotes] = useState('')
  const [confirmTarget, setConfirmTarget] = useState<ExpiryWarning | null>(null)
  // Sequential, bukan tumpuk — pola sama seperti FormModal + ConfirmDialog di seluruh
  // aplikasi (mis. ProductFormModal): modal induk digerbang tertutup selagi dialog
  // konfirmasi anak (write-off/confirm) tampil, bukan tetap terbuka di belakangnya.
  const isWriteOffConfirming = writeOffTarget !== null
  const isConfirmConfirming = confirmTarget !== null

  const { data } = useExpiryWarningsQuery()
  const { mutate: confirmBatch, isPending: isConfirming } = useConfirmExpiryBatchMutation()
  const { mutate: writeOff, isPending: isWritingOff } = useWriteOffExpiryBatchMutation()

  const warnings = (data?.warnings ?? []).filter((w) => w.product_id === product?.id)

  const handleWriteOffConfirm = () => {
    if (!writeOffTarget) return
    writeOff(
      { id: writeOffTarget.id, notes: writeOffNotes },
      {
        onSuccess: () => {
          setWriteOffTarget(null)
          setWriteOffNotes('')
        },
      }
    )
  }

  const handleConfirmBatch = () => {
    if (!confirmTarget) return
    confirmBatch({ id: confirmTarget.id }, { onSuccess: () => setConfirmTarget(null) })
  }

  return (
    <>
      <FormModal
        open={open && !isWriteOffConfirming && !isConfirmConfirming}
        onOpenChange={(val) => { if (!val && !isWriteOffConfirming && !isConfirmConfirming) onOpenChange(val) }}
        title="Batch Perlu Dicek"
        description={product?.name}
        size="md"
        hideFooter
      >
        <div className="space-y-2">
          {warnings.length === 0 ? (
            <p className="py-6 text-center text-sm text-gray-400">
              Tidak ada batch yang perlu dicek untuk produk ini.
            </p>
          ) : (
            warnings.map((w) => (
              <div key={w.id} className="rounded-lg border p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1.5">
                    <TriangleAlert
                      size={13}
                      className={w.severity === 'expired' ? 'text-red-500' : 'text-amber-500'}
                    />
                    <span className="text-sm font-medium">{w.qty} unit</span>
                    <span className="text-xs text-gray-400">· expired {formatDate(w.expired_date)}</span>
                  </div>
                  <span
                    className={`text-[11px] font-medium ${
                      w.severity === 'expired' ? 'text-red-600' : 'text-amber-600'
                    }`}
                  >
                    {daysLeftLabel(w)}
                  </span>
                </div>
                <p className="text-xs text-gray-500">
                  Cek fisik rak: kalau masih ada barangnya, musnahkan. Kalau sudah habis terjual, konfirmasi aman.
                </p>
                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="h-7 gap-1 text-xs"
                    disabled={isConfirming}
                    onClick={() => setConfirmTarget(w)}
                  >
                    <CheckCircle2 size={13} />
                    Sudah Dicek, Aman
                  </Button>
                  <RoleGuard menuKey="produk.produk" action="can_delete">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="h-7 gap-1 text-xs text-red-600 border-red-200 hover:bg-red-50"
                      onClick={() => setWriteOffTarget(w)}
                    >
                      <Trash2 size={13} />
                      Musnahkan
                    </Button>
                  </RoleGuard>
                </div>
              </div>
            ))
          )}
        </div>
      </FormModal>

      <ExtendedConfirmDialog
        open={writeOffTarget !== null}
        onOpenChange={(val) => {
          if (!val) {
            setWriteOffTarget(null)
            setWriteOffNotes('')
          }
        }}
        title="Musnahkan Stok Expired"
        variant="destructive"
        confirmLabel="Ya, Musnahkan"
        isLoading={isWritingOff}
        onConfirm={handleWriteOffConfirm}
        description={
          writeOffTarget
            ? `${writeOffTarget.qty} unit dengan expired ${formatDate(writeOffTarget.expired_date)} akan dikurangi dari stok dan tercatat sebagai kerugian. Tindakan ini tidak bisa dibatalkan.`
            : ''
        }
      >
        <Textarea
          value={writeOffNotes}
          onChange={(e) => setWriteOffNotes(e.target.value)}
          placeholder="Catatan (opsional)"
          className="resize-none"
          rows={2}
        />
      </ExtendedConfirmDialog>

      <ConfirmDialog
        open={confirmTarget !== null}
        onOpenChange={(val) => { if (!val) setConfirmTarget(null) }}
        title="Konfirmasi Batch Aman"
        confirmLabel="Ya, Sudah Dicek"
        isLoading={isConfirming}
        onConfirm={handleConfirmBatch}
        description={
          confirmTarget
            ? `${confirmTarget.qty} unit dengan expired ${formatDate(confirmTarget.expired_date)} akan ditandai sudah dicek dan aman. Pastikan Anda sudah memeriksa fisik rak.`
            : ''
        }
      />
    </>
  )
}
