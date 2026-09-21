import { useEffect, useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { formatRupiah, formatDate } from '@/shared/utils'

import { useResolveCostMutation, useSkipCostMutation } from '../capital-reconciliation.api'
import type { CapitalReconListItem } from '../capital-reconciliation.types'

interface Props {
  item: CapitalReconListItem | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ResolveCostModal({ item, open, onOpenChange }: Props) {
  const { mutate: resolve, isPending: isResolving } = useResolveCostMutation()
  const { mutate: skip, isPending: isSkipping } = useSkipCostMutation()

  const [cost, setCost] = useState('')
  const [note, setNote] = useState('')

  // Isi ulang form tiap ganti baris. Set-state di effect disengaja: sinkron dgn item baru.
  useEffect(() => {
    if (item) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setCost('')
      setNote('')
    }
  }, [item])

  const busy = isResolving || isSkipping

  const handleResolve = () => {
    if (!item) return
    const value = Number(cost)
    if (!cost || Number.isNaN(value) || value < 0) return
    resolve(
      { review_id: item.review_id, correct_cost: value, note: note || undefined },
      { onSuccess: () => onOpenChange(false) },
    )
  }

  const handleSkip = () => {
    if (!item) return
    skip(
      { review_id: item.review_id, note: note || undefined },
      { onSuccess: () => onOpenChange(false) },
    )
  }

  const costValue = Number(cost)
  const costInvalid = cost !== '' && (Number.isNaN(costValue) || costValue < 0)
  const costTooHigh = item != null && !costInvalid && cost !== '' && costValue >= item.sell_price

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Koreksi Modal (HPP)</DialogTitle>
          <DialogDescription>
            Tetapkan modal per satuan yang benar untuk baris penjualan ini.
          </DialogDescription>
        </DialogHeader>

        {item && (
          <div className="space-y-4">
            {/* Info baris */}
            <div className="rounded-lg border bg-gray-50 p-3 text-sm">
              <div className="font-semibold">{item.product_name}</div>
              <div className="mt-1 grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-gray-600">
                <span>Transaksi</span>
                <span className="text-right font-mono">{item.transaction_code}</span>
                <span>Tanggal</span>
                <span className="text-right">{formatDate(item.transaction_date)}</span>
                <span>Jumlah</span>
                <span className="text-right">
                  {item.quantity} {item.unit}
                </span>
                <span>Harga jual / satuan</span>
                <span className="text-right">{formatRupiah(item.sell_price)}</span>
                <span>Modal tercatat sekarang</span>
                <span className="text-right text-red-600">{formatRupiah(item.old_purchase_price)}</span>
                <span>Perkiraan sistem</span>
                <span className="text-right">{formatRupiah(item.suggested_price)}</span>
              </div>
              <div className="mt-2 flex items-start gap-1.5 rounded bg-amber-50 p-2 text-xs text-amber-700">
                <AlertTriangle size={13} className="mt-0.5 shrink-0" />
                <span>{item.reason}</span>
              </div>
            </div>

            {/* Input modal benar */}
            <div className="space-y-1">
              <Label htmlFor="correct-cost" className="text-xs text-gray-600">
                Modal yang benar (per {item.unit || 'satuan'})
              </Label>
              <Input
                id="correct-cost"
                type="number"
                min={0}
                step="any"
                value={cost}
                onChange={(e) => setCost(e.target.value)}
                placeholder="Masukkan modal per satuan"
                className={costInvalid || costTooHigh ? 'border-red-400' : ''}
              />
              {costInvalid && <p className="text-xs text-red-500">Nilai modal tidak valid.</p>}
              {costTooHigh && (
                <p className="text-xs text-amber-600">
                  Modal ini lebih besar/sama dengan harga jual — pastikan memang benar (rugi).
                </p>
              )}
            </div>

            {/* Catatan */}
            <div className="space-y-1">
              <Label htmlFor="note" className="text-xs text-gray-600">
                Catatan (opsional)
              </Label>
              <Textarea
                id="note"
                value={note}
                onChange={(e) => setNote(e.target.value)}
                placeholder="Contoh: modal riil sesuai nota beli"
                className="min-h-[56px]"
              />
            </div>
          </div>
        )}

        <DialogFooter className="flex items-center justify-between gap-2 sm:justify-between">
          <Button variant="ghost" onClick={handleSkip} disabled={busy || !item}>
            {isSkipping && <Loader2 size={14} className="mr-1 animate-spin" />}
            Tandai Selesai (tanpa ubah)
          </Button>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)} disabled={busy}>
              Batal
            </Button>
            <Button onClick={handleResolve} disabled={busy || !item || cost === '' || costInvalid}>
              {isResolving && <Loader2 size={14} className="mr-1 animate-spin" />}
              Simpan Modal
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
