import { useState } from 'react'
import { ArrowRightLeft, Combine, Printer, Trash2, X, XCircle } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { formatRupiah } from '@/shared/utils'
import type { Transaction, TransactionItem } from '@/features/sales/transactions/transactions.types'
import { TRANSACTION_STATUS_BADGE, transactionReceipt } from '@/features/sales/transactions/transactions.utils'

import { useCashierStore } from '../cashier.store'
import { ORDER_TYPE_LABEL } from '../cashier.types'
import { useMergeBillMutation } from '../openBill.api'
import { CancelBillDialog, MoveTableDialog, RemoveBillItemDialog } from './BillDialogs'
import { OpenBillsDialog } from './OpenBillsDialog'
import { ReceiptPrint } from './ReceiptPrint'

type Dialog = 'move' | 'merge' | 'cancel' | 'bill' | null

/**
 * Pesanan terbuka yang sedang dibuka di kasir (tahap 4.8): info pesanan, aksi (pindah meja, gabung,
 * tagihan, batal, tutup) dan item yang sudah tersimpan (bisa dihapus/dikurangi).
 */
export function ActiveBillPanel({ bill }: { bill: Transaction }) {
  const { setActiveBill } = useCashierStore()
  const [dialog, setDialog] = useState<Dialog>(null)
  const [removing, setRemoving] = useState<TransactionItem | null>(null)
  const { mutate: merge } = useMergeBillMutation()

  if (bill.status !== 'open') {
    return (
      <div className="shrink-0 space-y-2 border-b bg-amber-50 px-4 py-3 text-sm text-amber-800">
        <p>
          Pesanan antrean {bill.queue_label} sudah {TRANSACTION_STATUS_BADGE[bill.status].label.toLowerCase()}.
        </p>
        <Button size="sm" variant="outline" onClick={() => setActiveBill(null)}>
          Kembali ke pesanan baru
        </Button>
      </div>
    )
  }

  const receipt = transactionReceipt(bill)
  const action = (label: string, icon: React.ReactNode, onClick: () => void, danger = false) => (
    <Button type="button" variant="outline" size="sm" onClick={onClick} className={`h-8 gap-1 px-2 text-xs ${danger ? 'text-red-600 hover:text-red-700' : ''}`}>
      {icon}
      {label}
    </Button>
  )

  return (
    <div className="shrink-0 border-b">
      <div className="space-y-2 bg-blue-50 px-4 py-3">
        <div className="flex items-start gap-3">
          <div className="text-center">
            <p className="text-[10px] uppercase tracking-wide text-blue-700">Antrean</p>
            <p className="text-2xl font-bold leading-none text-blue-900" aria-label="Antrean pesanan terbuka">
              {bill.queue_label}
            </p>
          </div>
          <div className="min-w-0 flex-1 text-sm">
            <p className="font-semibold text-gray-800">
              {bill.order_type ? ORDER_TYPE_LABEL[bill.order_type] : ''}
              {bill.table_name && ` · ${bill.table_name}`}
            </p>
            <p className="text-xs text-gray-600">
              {bill.order_name || 'Tanpa nama'} · dibuat {bill.created_by_name || '—'}
            </p>
          </div>
          <button type="button" aria-label="Tutup pesanan terbuka" onClick={() => setActiveBill(null)} className="text-gray-400 hover:text-gray-600">
            <X size={16} />
          </button>
        </div>
        {bill.is_overdue && (
          <p className="rounded-md bg-amber-100 px-2 py-1.5 text-xs text-amber-800">
            Pesanan tertunda dari hari sebelumnya: hanya bisa dibayar, dijadikan utang, atau dibatalkan.
          </p>
        )}
        <div className="flex flex-wrap gap-1.5">
          {!bill.is_overdue && bill.order_type === 'dine_in' && action('Pindah Meja', <ArrowRightLeft size={12} />, () => setDialog('move'))}
          {!bill.is_overdue && action('Gabung', <Combine size={12} />, () => setDialog('merge'))}
          {action('Tagihan', <Printer size={12} />, () => setDialog('bill'))}
          {action('Batalkan', <XCircle size={12} />, () => setDialog('cancel'), true)}
        </div>
      </div>

      <ul className="max-h-48 divide-y divide-gray-100 overflow-y-auto" aria-label="Item pesanan tersimpan">
        {bill.items.length === 0 && <li className="px-4 py-3 text-xs text-gray-500">Belum ada item tersimpan.</li>}
        {bill.items.map((item) => (
          <li key={item.id} className="flex items-start gap-2 px-4 py-2" aria-label={`Item tersimpan ${item.product_name}`}>
            <div className="min-w-0 flex-1 text-sm">
              <p className="text-gray-800">
                {item.quantity}x {item.product_name}
                {item.variant_name && <span className="text-gray-500"> · {item.variant_name}</span>}
              </p>
              {item.modifiers.map((m) => (
                <p key={`${m.group_name}-${m.option_name}`} className="text-xs text-gray-500">
                  + {m.option_name}
                </p>
              ))}
              {item.notes && <p className="text-xs text-gray-500">Catatan: {item.notes}</p>}
            </div>
            <span className="text-sm font-medium text-gray-800">{formatRupiah(item.subtotal)}</span>
            {!bill.is_overdue && (
              <button type="button" aria-label={`Hapus ${item.product_name}`} onClick={() => setRemoving(item)} className="mt-0.5 text-gray-300 hover:text-red-500">
                <Trash2 size={14} />
              </button>
            )}
          </li>
        ))}
      </ul>

      {removing && <RemoveBillItemDialog bill={bill} item={removing} onClose={() => setRemoving(null)} />}
      {dialog === 'move' && <MoveTableDialog bill={bill} onClose={() => setDialog(null)} />}
      {dialog === 'cancel' && (
        <CancelBillDialog
          bill={bill}
          onClose={() => setDialog(null)}
          onDone={() => {
            setDialog(null)
            setActiveBill(null)
          }}
        />
      )}
      {dialog === 'merge' && (
        <OpenBillsDialog
          title={`Gabung ke antrean ${bill.queue_label}`}
          description="Item pesanan yang dipilih pindah ke pesanan ini; meja pesanan itu menjadi kosong."
          actionLabel="Gabung ke sini"
          excludeId={bill.id}
          onClose={() => setDialog(null)}
          onPick={(src) => merge({ id: bill.id, source_id: src.id }, { onSuccess: () => setDialog(null) })}
        />
      )}
      {dialog === 'bill' && (
        <ReceiptPrint
          open
          mode="bill"
          onClose={() => setDialog(null)}
          checkoutData={receipt.checkoutData}
          cart={receipt.cart}
          summary={receipt.summary}
          discount={receipt.discount}
          tax={receipt.tax}
          paymentMethod="cash"
          amountPaid={0}
        />
      )}
    </div>
  )
}
