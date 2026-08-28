import { useEffect, useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { Textarea } from '@/shared/components/ui/textarea'
import { RoleGuard } from '@/shared/components'
import { formatStockNumber } from '@/features/products/products'

import {
  useAdjustStockMutation,
  useMarkReviewedMutation,
  useReconDetailQuery,
} from '../stock-reconciliation.api'
import { RECON_MENU_KEY } from '../stock-reconciliation.constants'
import { ReconBreakdownCard } from './ReconBreakdownCard'
import { ReconKartuStok } from './ReconKartuStok'
import { ReconLevelInputTable } from './ReconLevelInputTable'

interface ReconDetailModalProps {
  productId: number | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

function Stat({ label, value, unit, tone }: { label: string; value: string; unit: string; tone?: 'danger' }) {
  return (
    <div className="rounded-lg border bg-white px-3 py-2.5">
      <p className="text-[11px] text-gray-500">{label}</p>
      <p className={`text-xl font-bold ${tone === 'danger' ? 'text-red-600' : 'text-gray-800'}`}>{value}</p>
      <p className="text-[11px] font-medium text-gray-400">{unit}</p>
    </div>
  )
}

export function ReconDetailModal({ productId, open, onOpenChange }: ReconDetailModalProps) {
  const { data: detail, isLoading } = useReconDetailQuery(open ? productId : null)
  const { mutate: adjust, isPending: isAdjusting } = useAdjustStockMutation()
  const { mutate: markReviewed, isPending: isMarking } = useMarkReviewedMutation()

  const [levelValues, setLevelValues] = useState<Record<number, string>>({})
  const [note, setNote] = useState('')

  // reset form tiap ganti produk / buka ulang
  useEffect(() => {
    if (detail) {
      const init: Record<number, string> = {}
      detail.packages.forEach((p) => {
        init[p.package_id] = String(p.current_stock)
      })
      setLevelValues(init)
      setNote('')
    }
  }, [detail])

  const busy = isAdjusting || isMarking

  const handleLevelChange = (packageId: number, value: string) => {
    setLevelValues((prev) => ({ ...prev, [packageId]: value }))
  }

  const handleSubmit = () => {
    if (!detail) return
    const levels = detail.packages.map((p) => ({
      package_id: p.package_id,
      new_stock: Number(levelValues[p.package_id] ?? p.current_stock) || 0,
    }))
    adjust(
      { product_id: detail.product_id, levels, note: note || undefined },
      { onSuccess: () => onOpenChange(false) },
    )
  }

  const handleMarkReviewed = () => {
    if (!detail) return
    markReviewed(
      { product_id: detail.product_id, note: note || undefined },
      { onSuccess: () => onOpenChange(false) },
    )
  }

  return (
    <Dialog open={open} onOpenChange={(v) => !busy && onOpenChange(v)}>
      <DialogContent
        className="flex max-h-[88vh] max-w-[820px] flex-col gap-0 overflow-hidden p-0"
        onInteractOutside={(e) => busy && e.preventDefault()}
      >
        <DialogHeader className="border-b px-6 py-4">
          <DialogTitle>{detail?.product_name ?? 'Detail Rekonsiliasi'}</DialogTitle>
          <DialogDescription>
            {detail
              ? `${detail.product_code || '-'} · Satuan dasar: ${detail.base_unit || '-'}`
              : 'Memuat detail...'}
          </DialogDescription>
        </DialogHeader>

        <ScrollArea className="min-h-0 flex-1">
          <div className="space-y-5 px-6 py-4">
            {isLoading || !detail ? (
              <div className="flex items-center justify-center py-16 text-sm text-gray-400">
                <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Memuat...
              </div>
            ) : (
              <>
                {/* Headline lama vs baru */}
                <div className="grid grid-cols-3 gap-3">
                  <Stat
                    label="Stok Lama"
                    value={detail.old_stock_available ? formatStockNumber(detail.old_stock) : '-'}
                    unit={detail.base_unit}
                  />
                  <Stat label="Stok Baru" value={formatStockNumber(detail.new_stock)} unit={detail.base_unit} />
                  <Stat
                    label="Selisih"
                    value={detail.old_stock_available ? formatStockNumber(detail.diff) : '-'}
                    unit={detail.base_unit}
                    tone={detail.old_stock_available && detail.diff !== 0 ? 'danger' : undefined}
                  />
                </div>

                {detail.needs_stock_review && detail.stock_review_note && (
                  <div className="flex gap-2 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2.5 text-sm text-amber-800">
                    <AlertTriangle size={16} className="mt-0.5 flex-shrink-0" />
                    <span>
                      <b>Perlu Ditinjau:</b> {detail.stock_review_note}
                    </span>
                  </div>
                )}

                <ReconBreakdownCard breakdown={detail.breakdown} baseUnit={detail.base_unit} />

                <ReconKartuStok rows={detail.kartu_stok} />

                <RoleGuard
                  menuKey={RECON_MENU_KEY}
                  action="can_edit"
                  fallback={
                    <p className="text-xs text-gray-400">Anda tidak memiliki izin untuk mengoreksi stok.</p>
                  }
                >
                  <ReconLevelInputTable
                    packages={detail.packages}
                    values={levelValues}
                    onChange={handleLevelChange}
                    baseUnit={detail.base_unit}
                  />
                  <div>
                    <label className="mb-1 block text-xs text-gray-500">Catatan koreksi</label>
                    <Textarea
                      value={note}
                      onChange={(e) => setNote(e.target.value)}
                      placeholder="Contoh: Hasil hitung fisik gudang"
                      className="min-h-[56px]"
                    />
                  </div>
                </RoleGuard>
              </>
            )}
          </div>
        </ScrollArea>

        <div className="flex items-center justify-between gap-2 border-t px-6 py-4">
          <RoleGuard menuKey={RECON_MENU_KEY} action="can_edit">
            <Button variant="ghost" onClick={handleMarkReviewed} disabled={busy || !detail}>
              {isMarking && <Loader2 size={14} className="mr-1 animate-spin" />}
              Tandai Selesai
            </Button>
          </RoleGuard>
          <div className="ml-auto flex gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)} disabled={busy}>
              Batal
            </Button>
            <RoleGuard menuKey={RECON_MENU_KEY} action="can_edit">
              <Button onClick={handleSubmit} disabled={busy || !detail}>
                {isAdjusting && <Loader2 size={14} className="mr-1 animate-spin" />}
                Simpan Koreksi Stok
              </Button>
            </RoleGuard>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
