import { useState } from 'react'
import { Minus, Plus } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { formatRupiah } from '@/shared/utils'
import type { Transaction, TransactionItem } from '@/features/sales/transactions/transactions.types'

import { useCancelBillMutation, useMoveBillTableMutation, useRemoveBillItemMutation, useTableOptionsQuery } from '../openBill.api'
import { StockChoice, type StockAction } from './StockChoice'

const itemLabel = (i: TransactionItem) => `${i.product_name}${i.variant_name ? ` · ${i.variant_name}` : ''}`

/** Hapus atau kurangi item pesanan terbuka (keputusan 105): barang jadi pilih kembali ke stok/dibuang. */
export function RemoveBillItemDialog({ bill, item, onClose }: { bill: Transaction; item: TransactionItem; onClose: () => void }) {
  const [qty, setQty] = useState(item.quantity)
  const [action, setAction] = useState<StockAction>('return')
  const [reason, setReason] = useState('')
  const { mutate, isPending } = useRemoveBillItemMutation()
  const isStock = item.product_id !== null
  const reasonMissing = isStock && action === 'waste' && reason.trim() === ''

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title="Hapus Item"
      description={`${itemLabel(item)} · antrean ${bill.queue_label}`}
      contentClassName="max-w-md"
      isLoading={isPending}
      submitLabel={qty >= item.quantity ? 'Hapus Item' : `Kurangi ${qty}`}
      submitDisabled={reasonMissing || isPending}
      onSubmit={() =>
        mutate(
          { id: bill.id, item_id: item.id, quantity: qty >= item.quantity ? 0 : qty, action: isStock ? action : undefined, reason: reason.trim() || undefined },
          { onSuccess: onClose },
        )
      }
    >
      <div className="space-y-4 px-6 py-4">
        <div className="flex items-center justify-between gap-3">
          <Label>Jumlah dihapus</Label>
          <div className="flex items-center gap-1.5">
            <Button type="button" variant="outline" size="icon" className="h-9 w-9" aria-label="Kurangi jumlah dihapus" onClick={() => setQty((q) => Math.max(1, q - 1))}>
              <Minus size={14} />
            </Button>
            <span className="w-14 text-center text-sm font-semibold" aria-label="Jumlah dihapus">
              {qty} / {item.quantity}
            </span>
            <Button type="button" variant="outline" size="icon" className="h-9 w-9" aria-label="Tambah jumlah dihapus" onClick={() => setQty((q) => Math.min(item.quantity, q + 1))}>
              <Plus size={14} />
            </Button>
          </div>
        </div>
        {isStock && (
          <div className="flex items-center justify-between gap-3">
            <Label>Barang jadi</Label>
            <StockChoice value={action} onChange={setAction} label={`Stok ${item.product_name}`} />
          </div>
        )}
        <div className="space-y-1.5">
          <Label htmlFor="remove-reason">
            Alasan {isStock && action === 'waste' ? <span className="text-red-500">*</span> : <span className="text-xs font-normal text-gray-500">(opsional)</span>}
          </Label>
          <Textarea id="remove-reason" rows={2} maxLength={255} value={reason} onChange={(e) => setReason(e.target.value)} className="resize-none" placeholder="Contoh: salah pesan" />
          {reasonMissing && <p className="text-xs text-red-500">Alasan wajib diisi untuk barang yang dibuang.</p>}
        </div>
      </div>
    </ActionModal>
  )
}

/** Pindah meja: hanya ke meja aktif yang kosong (keputusan 103). */
export function MoveTableDialog({ bill, onClose }: { bill: Transaction; onClose: () => void }) {
  const { data: tables = [] } = useTableOptionsQuery()
  const { mutate, isPending } = useMoveBillTableMutation()
  const free = tables.filter((t) => t.open_bill_id === null && t.id !== bill.table_id)

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title="Pindah Meja"
      description={`Antrean ${bill.queue_label} · sekarang ${bill.table_name}`}
      contentClassName="max-w-md"
      isLoading={isPending}
      hideFooter
    >
      <div className="grid max-h-[60vh] grid-cols-2 gap-2 overflow-y-auto px-5 py-4 sm:grid-cols-3">
        {free.length === 0 && <p className="col-span-full py-4 text-center text-sm text-gray-500">Tidak ada meja kosong.</p>}
        {free.map((t) => (
          <Button key={t.id} variant="outline" className="h-12" disabled={isPending} onClick={() => mutate({ id: bill.id, table_id: t.id }, { onSuccess: onClose })}>
            {t.name}
          </Button>
        ))}
      </div>
    </ActionModal>
  )
}

/** Batal pesanan terbuka: alasan wajib; barang jadi per item kembali ke stok (bawaan) atau dibuang. */
export function CancelBillDialog({ bill, onClose, onDone }: { bill: Transaction; onClose: () => void; onDone: () => void }) {
  const stockItems = bill.items.filter((i) => i.product_id !== null)
  const [actions, setActions] = useState<Record<number, StockAction>>(() => Object.fromEntries(stockItems.map((i) => [i.id, 'return'])))
  const [reason, setReason] = useState('')
  const { mutate, isPending } = useCancelBillMutation()

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title="Batalkan Pesanan"
      description={`Antrean ${bill.queue_label} · ${formatRupiah(bill.subtotal)}`}
      contentClassName="max-w-lg"
      isLoading={isPending}
      submitLabel="Ya, Batalkan"
      submitDisabled={reason.trim() === '' || isPending}
      onSubmit={() =>
        mutate(
          { id: bill.id, reason: reason.trim(), items: stockItems.map((i) => ({ item_id: i.id, action: actions[i.id] ?? 'return' })) },
          { onSuccess: onDone },
        )
      }
    >
      <div className="max-h-[65vh] space-y-4 overflow-y-auto px-6 py-4">
        {stockItems.map((item) => (
          <div key={item.id} className="flex flex-wrap items-center justify-between gap-2 rounded-lg border p-2.5">
            <span className="text-sm text-gray-800">
              {itemLabel(item)} x{item.quantity}
            </span>
            <StockChoice value={actions[item.id] ?? 'return'} onChange={(v) => setActions((p) => ({ ...p, [item.id]: v }))} label={`Stok ${item.product_name}`} />
          </div>
        ))}
        <div className="space-y-1.5">
          <Label htmlFor="cancel-reason">
            Alasan batal <span className="text-red-500">*</span>
          </Label>
          <Textarea id="cancel-reason" rows={2} maxLength={255} value={reason} onChange={(e) => setReason(e.target.value)} className="resize-none" placeholder="Contoh: pelanggan tidak jadi" />
        </div>
      </div>
    </ActionModal>
  )
}
