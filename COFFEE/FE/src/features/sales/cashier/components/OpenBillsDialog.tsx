import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { formatRupiah } from '@/shared/utils'

import { ORDER_TYPE_LABEL } from '../cashier.types'
import { useOpenBillsQuery, type OpenBillSummary } from '../openBill.api'

interface OpenBillsDialogProps {
  title: string
  description?: string
  actionLabel: string
  /** Pesanan yang tidak ditampilkan (mis. pesanan aktif saat memilih bill untuk digabung). */
  excludeId?: number | null
  onPick: (bill: OpenBillSummary) => void
  onClose: () => void
}

function timeLabel(iso: string): string {
  return new Date(iso).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

/** Daftar pesanan terbuka (semua kasir, keputusan 106), dipakai untuk membuka atau menggabung bill. */
export function OpenBillsDialog({ title, description, actionLabel, excludeId, onPick, onClose }: OpenBillsDialogProps) {
  const { data = [], isLoading } = useOpenBillsQuery()
  const bills = data.filter((b) => b.id !== excludeId)

  return (
    <ActionModal
      open
      onOpenChange={(open) => {
        if (!open) onClose()
      }}
      title={title}
      description={description}
      contentClassName="max-w-lg max-h-[90vh]"
      hideFooter
    >
      <div className="max-h-[70vh] space-y-2 overflow-y-auto px-5 py-4">
        {isLoading ? (
          [1, 2, 3].map((i) => <div key={i} className="h-16 animate-pulse rounded-lg bg-gray-100" />)
        ) : bills.length === 0 ? (
          <p className="py-6 text-center text-sm text-gray-500">Tidak ada pesanan terbuka.</p>
        ) : (
          bills.map((b) => (
            <div key={b.id} className="flex items-center gap-3 rounded-lg border p-3" aria-label={`Pesanan terbuka ${b.queue_label}`}>
              <span className="w-12 shrink-0 text-center text-xl font-bold text-gray-800">{b.queue_label}</span>
              <div className="min-w-0 flex-1 text-sm">
                {b.is_overdue && (
                  <span className="mb-0.5 inline-block rounded bg-amber-100 px-1.5 py-0.5 text-[11px] font-semibold text-amber-700">Tertunda · {b.bill_date}</span>
                )}
                <p className="font-medium text-gray-800">
                  {b.order_type ? ORDER_TYPE_LABEL[b.order_type] : '—'}
                  {b.table_name && ` · ${b.table_name}`}
                  {b.order_name && ` · ${b.order_name}`}
                </p>
                <p className="text-xs text-gray-500">
                  {b.item_count} item · {formatRupiah(b.subtotal)} · {timeLabel(b.created_at)} oleh {b.created_by_name || '—'}
                </p>
              </div>
              <Button size="sm" onClick={() => onPick(b)}>
                {actionLabel}
              </Button>
            </div>
          ))
        )}
      </div>
    </ActionModal>
  )
}
