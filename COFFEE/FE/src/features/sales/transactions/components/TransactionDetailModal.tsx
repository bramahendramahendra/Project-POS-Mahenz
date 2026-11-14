import { useState } from 'react'
import { Printer } from 'lucide-react'

import { ActionModal, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { ReceiptPrint } from '@/features/sales/cashier/components/ReceiptPrint'
import { ORDER_TYPE_LABEL } from '@/features/sales/cashier/cashier.types'

import { useTransactionDetailQuery } from '../transactions.api'
import { PAYMENT_LABELS, TRANSACTION_STATUS_BADGE, formatDateTimeLong, transactionReceipt } from '../transactions.utils'
import { VoidTransactionDialog } from './VoidTransactionDialog'

interface TransactionDetailModalProps {
  transactionId: number | null
  onClose: () => void
}

export function TransactionDetailModal({ transactionId, onClose }: TransactionDetailModalProps) {
  const [voidOpen, setVoidOpen] = useState(false)
  const [receiptOpen, setReceiptOpen] = useState(false)

  const { data: transaction, isLoading } = useTransactionDetailQuery(transactionId ?? 0)

  const open = transactionId !== null
  // Dialog anak (void, struk cetak ulang) tampil bergantian dengan modal detail, bukan bertumpuk.
  const isChildDialogOpen = voidOpen || receiptOpen

  const footer = transaction ? (
    <div className="flex shrink-0 items-center gap-2 border-t px-5 py-3">
      <Button variant="outline" size="sm" className="gap-1" onClick={() => setReceiptOpen(true)}>
        <Printer size={14} />
        Cetak Ulang
      </Button>
      <RoleGuard menuKey="penjualan.transaksi" action="can_delete">
        {transaction.status === 'completed' && (
          <Button variant="destructive" size="sm" className="ml-auto" onClick={() => setVoidOpen(true)}>
            Void Transaksi
          </Button>
        )}
      </RoleGuard>
    </div>
  ) : null

  const meta = transaction
    ? [
        ['Tanggal', formatDateTimeLong(transaction.transaction_date)],
        ['Antrean', transaction.queue_label || '—'],
        ['Pesanan', transaction.order_type ? ORDER_TYPE_LABEL[transaction.order_type] : '—'],
        ['Meja', transaction.table_name || '—'],
        ['Nama Pemesan', transaction.order_name || '—'],
        ['Dibuat oleh', transaction.created_by_name || '—'],
        ['Kasir penerima', transaction.status === 'completed' || transaction.status === 'void' ? transaction.kasir_name : '—'],
        ['Diselesaikan oleh', transaction.settled_by_name || '—'],
        ['Pelanggan', transaction.customer_name || '—'],
        ['Metode Bayar', transaction.status === 'completed' || transaction.status === 'void' ? PAYMENT_LABELS[transaction.payment_method] : '—'],
      ]
    : []

  return (
    <>
      <ActionModal
        open={open && !isChildDialogOpen}
        onOpenChange={(val) => {
          if (!val && !isChildDialogOpen) onClose()
        }}
        title="Detail Transaksi"
        description={transaction?.transaction_code}
        descriptionClassName="font-mono text-xs"
        contentClassName="max-w-lg max-h-[90vh]"
        headerClassName="border-b px-5 py-3 shrink-0"
        footer={footer}
      >
        <ScrollArea className="flex-1">
          <div className="p-5">
            {isLoading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-8 animate-pulse rounded bg-gray-100" />
                ))}
              </div>
            ) : transaction ? (
              <div className="space-y-5">
                <div className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                  {meta.map(([label, value]) => (
                    <div key={label}>
                      <p className="text-xs text-gray-500">{label}</p>
                      <p className="font-medium">{value}</p>
                    </div>
                  ))}
                  <div>
                    <p className="text-xs text-gray-500">Status</p>
                    <StatusBadge {...TRANSACTION_STATUS_BADGE[transaction.status]} />
                  </div>
                </div>

                <div>
                  <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">Item</p>
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b text-xs text-gray-500">
                        <th className="pb-1 text-left">Menu</th>
                        <th className="pb-1 text-center">Qty</th>
                        <th className="pb-1 text-right">Harga</th>
                        <th className="pb-1 text-right">Subtotal</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y">
                      {transaction.items.map((item) => (
                        <tr key={item.id}>
                          <td className="py-1.5 align-top">
                            <p className="font-medium">
                              {item.product_name}
                              {item.variant_name && <span className="font-normal text-gray-500"> · {item.variant_name}</span>}
                            </p>
                            {item.modifiers.map((m) => (
                              <p key={`${m.group_name}-${m.option_name}`} className="text-xs text-gray-500">
                                + {m.option_name}
                                {m.price > 0 && ` ${formatRupiah(m.price)}`}
                              </p>
                            ))}
                            {item.notes && <p className="text-xs text-gray-500">Catatan: {item.notes}</p>}
                            {item.product_id !== null && <p className="text-xs text-blue-600">Barang jadi</p>}
                            {item.discount_item > 0 && <p className="text-xs text-green-600">Disc -{formatRupiah(item.discount_item)}</p>}
                          </td>
                          <td className="py-1.5 text-center align-top">{item.quantity}</td>
                          <td className="py-1.5 text-right align-top">{formatRupiah(item.price)}</td>
                          <td className="py-1.5 text-right align-top font-medium">{formatRupiah(item.subtotal)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {transaction.status === 'cancelled' && (
                  <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">Dibatalkan: {transaction.cancel_reason}</p>
                )}
                {transaction.status === 'merged' && transaction.merged_into_id && (
                  <p className="rounded-md border bg-gray-50 px-3 py-2 text-sm text-gray-600">Digabung ke pesanan ID {transaction.merged_into_id}.</p>
                )}
                {!!transaction.cancellations?.length && (
                  <div>
                    <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">Item dihapus / dibatalkan</p>
                    <ul className="space-y-1 text-xs text-gray-600">
                      {transaction.cancellations.map((c, i) => (
                        <li key={i}>
                          {c.quantity}x {c.item_name}
                          {c.variant_name && ` · ${c.variant_name}`}
                          {c.stock_action === 'waste' ? ' (dibuang)' : c.stock_action === 'return' ? ' (kembali ke stok)' : ''}
                          {' · '}
                          {c.user_name}
                          {c.reason && ` · ${c.reason}`}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                <div className="space-y-1.5 rounded-lg bg-gray-50 p-4 text-sm">
                  <div className="flex justify-between text-gray-600">
                    <span>Subtotal</span>
                    <span>{formatRupiah(transaction.subtotal)}</span>
                  </div>
                  {transaction.discount > 0 && (
                    <div className="flex justify-between text-green-600">
                      <span>Diskon</span>
                      <span>-{formatRupiah(transaction.discount)}</span>
                    </div>
                  )}
                  {transaction.tax > 0 && (
                    <div className="flex justify-between text-gray-600">
                      <span>Pajak</span>
                      <span>+{formatRupiah(transaction.tax)}</span>
                    </div>
                  )}
                  <div className="flex justify-between border-t pt-2 font-bold text-gray-900">
                    <span>TOTAL</span>
                    <span>{formatRupiah(transaction.total_amount)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Bayar</span>
                    <span>{formatRupiah(transaction.payment_amount)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Kembalian</span>
                    <span>{formatRupiah(transaction.change_amount)}</span>
                  </div>
                </div>
              </div>
            ) : null}
          </div>
        </ScrollArea>
      </ActionModal>

      {transaction && voidOpen && (
        <VoidTransactionDialog
          transaction={transaction}
          onClose={() => setVoidOpen(false)}
          onDone={() => {
            setVoidOpen(false)
            onClose()
          }}
        />
      )}

      {transaction &&
        receiptOpen &&
        (() => {
          const r = transactionReceipt(transaction)
          return (
            <ReceiptPrint
              open={receiptOpen}
              onClose={() => setReceiptOpen(false)}
              checkoutData={r.checkoutData}
              cart={r.cart}
              summary={r.summary}
              discount={r.discount}
              tax={r.tax}
              paymentMethod={transaction.payment_method}
              amountPaid={transaction.payment_amount}
              customerName={transaction.customer_name}
              mode="reprint"
            />
          )
        })()}
    </>
  )
}
