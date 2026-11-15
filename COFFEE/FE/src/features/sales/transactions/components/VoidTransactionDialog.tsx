import { useState } from 'react'
import { AlertTriangle } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { formatRupiah } from '@/shared/utils'
import { StockChoice, type StockAction } from '@/features/sales/cashier/components/StockChoice'

import { useVoidTransactionMutation } from '../transactions.api'
import type { Transaction } from '../transactions.types'

interface VoidTransactionDialogProps {
  transaction: Transaction
  onClose: () => void
  onDone: () => void
}

/**
 * Void transaksi lunas (keputusan 100): tiap item barang jadi dipilih "kembali ke stok" (bawaan)
 * atau "dibuang"; alasan wajib bila ada yang dibuang. Racikan tidak punya pilihan stok.
 */
export function VoidTransactionDialog({ transaction, onClose, onDone }: VoidTransactionDialogProps) {
  const stockItems = transaction.items.filter((i) => i.product_id !== null)
  const [actions, setActions] = useState<Record<number, StockAction>>(() => Object.fromEntries(stockItems.map((i) => [i.id, 'return'])))
  const [reason, setReason] = useState('')
  const { mutate: voidTransaction, isPending } = useVoidTransactionMutation()

  const hasWaste = Object.values(actions).includes('waste')
  const reasonMissing = hasWaste && reason.trim() === ''

  const submit = () => {
    if (reasonMissing) return
    voidTransaction(
      {
        id: transaction.id,
        items: stockItems.map((i) => ({ item_id: i.id, action: actions[i.id] ?? 'return' })),
        reason: hasWaste ? reason.trim() : undefined,
      },
      { onSuccess: onDone },
    )
  }

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title="Void Transaksi"
      description={`${transaction.transaction_code}${transaction.queue_label ? ` · antrean ${transaction.queue_label}` : ''}`}
      contentClassName="max-w-lg"
      isLoading={isPending}
      onSubmit={submit}
      submitLabel="Ya, Void"
      submitDisabled={reasonMissing || isPending}
    >
      <div className="max-h-[65vh] space-y-4 overflow-y-auto px-6 py-4">
        <div className="flex items-start gap-2 rounded-md border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700">
          <AlertTriangle size={15} className="mt-0.5 shrink-0" />
          <span>
            Transaksi {formatRupiah(transaction.total_amount)} dibatalkan dan tidak bisa dikembalikan. Kas tunai dan piutang ikut
            disesuaikan.
          </span>
        </div>

        {stockItems.length > 0 ? (
          <div className="space-y-2">
            <Label>Barang jadi</Label>
            {stockItems.map((item) => (
              <div key={item.id} className="flex flex-wrap items-center justify-between gap-2 rounded-lg border p-2.5">
                <span className="text-sm text-gray-800">
                  {item.product_name}
                  {item.variant_name && <span className="text-gray-500"> · {item.variant_name}</span>} x{item.quantity}
                </span>
                <StockChoice
                  value={actions[item.id]}
                  onChange={(a) => setActions((prev) => ({ ...prev, [item.id]: a }))}
                  label={`Stok ${item.product_name}`}
                />
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-gray-500">Tidak ada barang jadi di transaksi ini (racikan tidak mengubah stok).</p>
        )}

        {hasWaste && (
          <div className="space-y-1.5">
            <Label htmlFor="void-reason">
              Alasan dibuang <span className="text-red-500">*</span>
            </Label>
            <Textarea id="void-reason" value={reason} maxLength={255} onChange={(e) => setReason(e.target.value)} placeholder="Contoh: roti jatuh" rows={2} className="resize-none" />
            {reasonMissing && <p className="text-xs text-red-500">Alasan wajib diisi untuk barang yang dibuang.</p>}
          </div>
        )}
      </div>
    </ActionModal>
  )
}
