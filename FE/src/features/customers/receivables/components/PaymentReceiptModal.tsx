import { Printer } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { formatRupiah, formatDate } from '@/shared/utils'
import { useStoreProfileQuery } from '@/features/settings/store'

// Data bukti pembayaran cicilan piutang.
export interface PaymentReceiptData {
  customerName: string
  transactionCode: string
  amount: number
  paymentDate: string
  remainingAfter: number
  notes?: string
}

interface PaymentReceiptModalProps {
  open: boolean
  onClose: () => void
  data: PaymentReceiptData | null
}

export function PaymentReceiptModal({ open, onClose, data }: PaymentReceiptModalProps) {
  const { data: store } = useStoreProfileQuery()

  const footer = (
    <div className="no-print border-t px-5 py-3 flex items-center justify-end gap-2 shrink-0">
      <Button variant="outline" size="sm" onClick={onClose}>
        Tutup
      </Button>
      <Button size="sm" className="gap-1" onClick={() => window.print()}>
        <Printer size={14} />
        Cetak
      </Button>
    </div>
  )

  return (
    <ActionModal
      open={open}
      onOpenChange={(v) => { if (!v) onClose() }}
      title="Bukti Pembayaran"
      description="Bukti pembayaran cicilan piutang"
      contentClassName="max-w-sm"
      headerClassName="border-b px-5 py-3 shrink-0 no-print"
      footer={footer}
    >
      {data && (
        <div className="print-root p-5">
          <div className="text-center mb-3">
            <p className="font-bold text-base">{store?.name ?? 'POS System'}</p>
            {store?.address && <p className="text-xs text-gray-500">{store.address}</p>}
            {store?.phone && <p className="text-xs text-gray-500">{store.phone}</p>}
          </div>

          <p className="text-center text-sm font-semibold border-y py-1.5 mb-3">
            BUKTI PEMBAYARAN PIUTANG
          </p>

          <div className="space-y-1.5 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-500">Tanggal</span>
              <span className="font-medium">{formatDate(data.paymentDate)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">Pelanggan</span>
              <span className="font-medium">{data.customerName}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">No. Transaksi</span>
              <span className="font-mono text-xs">{data.transactionCode}</span>
            </div>
          </div>

          <div className="border-t mt-3 pt-3 space-y-1.5 text-sm">
            <div className="flex justify-between text-base">
              <span className="font-semibold">Dibayar</span>
              <span className="font-bold text-green-600">{formatRupiah(data.amount)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">Sisa Piutang</span>
              <span className={data.remainingAfter > 0 ? 'text-red-600 font-medium' : 'text-gray-500'}>
                {formatRupiah(data.remainingAfter)}
              </span>
            </div>
            {data.remainingAfter <= 0 && (
              <p className="text-center text-green-600 font-semibold pt-1">— LUNAS —</p>
            )}
          </div>

          {data.notes && (
            <p className="text-xs text-gray-400 mt-3 border-t pt-2">Catatan: {data.notes}</p>
          )}

          <p className="text-center text-xs text-gray-400 mt-4">Terima kasih</p>
        </div>
      )}
    </ActionModal>
  )
}
