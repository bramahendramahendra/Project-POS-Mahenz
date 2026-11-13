import { useState } from 'react'
import { Printer } from 'lucide-react'

import { ActionModal, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah, formatDate } from '@/shared/utils'
import { ReceiptPrint } from '@/features/sales/cashier/components/ReceiptPrint'
import type { CartItem, CartSummary, Discount, Tax } from '@/features/sales/cashier'

import { useReceivableDetailQuery } from '../receivables.api'
import type { Receivable, ReceivableDetail, ReceivableStatus } from '../receivables.types'

interface ReceivableDetailModalProps {
  receivable: Receivable | null
  onClose: () => void
  onPay: (receivable: Receivable) => void
}

// Status piutang -> tipe StatusBadge (keduanya memakai istilah yang sama).
const STATUS_BADGE: Record<ReceivableStatus, 'unpaid' | 'partial' | 'paid' | 'void'> = {
  unpaid: 'unpaid',
  partial: 'partial',
  paid: 'paid',
  void: 'void',
}

// Bangun data untuk komponen ReceiptPrint dari detail piutang (pola sama
// TransactionDetailModal). unit_id/conversion_qty tidak esensial untuk tampilan
// struk cetak-ulang, jadi diberi nilai aman.
function buildReceiptData(d: ReceivableDetail): {
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
} {
  const cart: CartItem[] = d.items.map((item, idx) => ({
    product_id: idx,
    product_name: item.product_name,
    unit_id: idx,
    unit_name: item.unit,
    conversion_qty: 1,
    qty: item.quantity,
    price: item.price,
    subtotal: item.subtotal,
    discount_amount: item.discount_item > 0 ? item.discount_item : undefined,
  }))

  const itemsSubtotal = d.items.reduce((s, it) => s + it.subtotal, 0)
  const summary: CartSummary = {
    subtotal: itemsSubtotal,
    discountAmount: 0,
    taxAmount: 0,
    grandTotal: d.total_amount,
  }
  const discount: Discount = { type: 'amount', value: 0, amount: 0 }
  const tax: Tax = { percent: 0, amount: 0 }
  return { cart, summary, discount, tax }
}

export function ReceivableDetailModal({ receivable, onClose, onPay }: ReceivableDetailModalProps) {
  const [receiptOpen, setReceiptOpen] = useState(false)

  const id = receivable?.id ?? 0
  const open = receivable !== null
  const { data: detail, isLoading } = useReceivableDetailQuery(id)

  const isChildOpen = receiptOpen
  const canPay = receivable && receivable.status !== 'paid' && receivable.status !== 'void'

  const footer = detail ? (
    <div className="border-t px-5 py-3 flex items-center gap-2 shrink-0">
      <Button variant="outline" size="sm" className="gap-1" onClick={() => setReceiptOpen(true)}>
        <Printer size={14} />
        Cetak Nota
      </Button>
      {canPay && receivable && (
        <Button size="sm" className="ml-auto" onClick={() => onPay(receivable)}>
          Catat Pembayaran
        </Button>
      )}
    </div>
  ) : null

  return (
    <>
      <ActionModal
        open={open && !isChildOpen}
        onOpenChange={(val) => { if (!val && !isChildOpen) onClose() }}
        title="Detail Piutang"
        description={receivable?.transaction_code}
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
            ) : detail ? (
              <div className="space-y-5">
                {/* Bagian 1: Ringkasan */}
                <div className="grid grid-cols-1 gap-x-4 gap-y-2 text-sm sm:grid-cols-2">
                  <div>
                    <p className="text-xs text-gray-500">Pelanggan</p>
                    <p className="font-medium">{detail.customer_name}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Status</p>
                    <StatusBadge status={STATUS_BADGE[detail.status]} />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Tanggal Transaksi</p>
                    <p className="font-medium">{formatDate(detail.created_at)}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Jatuh Tempo</p>
                    <p className="font-medium">{detail.due_date ? formatDate(detail.due_date) : '—'}</p>
                  </div>
                </div>

                {/* Ringkasan nominal */}
                <div className="rounded-lg bg-gray-50 p-4 space-y-1.5 text-sm">
                  <div className="flex justify-between text-gray-600">
                    <span>Total Piutang</span>
                    <span className="font-medium">{formatRupiah(detail.total_amount)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Sudah Dibayar</span>
                    <span className="text-green-600">{formatRupiah(detail.paid_amount)}</span>
                  </div>
                  <div className="flex justify-between border-t pt-1.5 font-semibold">
                    <span>Sisa</span>
                    <span className={detail.remaining_amount > 0 ? 'text-red-600' : 'text-gray-500'}>
                      {formatRupiah(detail.remaining_amount)}
                    </span>
                  </div>
                </div>

                {/* Bagian 2: Rincian Barang */}
                <div>
                  <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">
                    Rincian Barang
                  </p>
                  {detail.items.length > 0 ? (
                    <table className="w-full text-sm">
                      <thead>
                        <tr className="border-b text-xs text-gray-500">
                          <th className="text-left pb-1">Produk</th>
                          <th className="text-center pb-1">Qty</th>
                          <th className="text-right pb-1">Harga</th>
                          <th className="text-right pb-1">Subtotal</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y">
                        {detail.items.map((item, i) => (
                          <tr key={i}>
                            <td className="py-1.5">
                              <p className="font-medium">{item.product_name}</p>
                              <p className="text-xs text-gray-400">{item.unit}</p>
                              {item.discount_item > 0 && (
                                <p className="text-xs text-green-600">
                                  Disc -{formatRupiah(item.discount_item)}
                                </p>
                              )}
                            </td>
                            <td className="text-center py-1.5">{item.quantity}</td>
                            <td className="text-right py-1.5">{formatRupiah(item.price)}</td>
                            <td className="text-right py-1.5 font-medium">
                              {formatRupiah(item.subtotal)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  ) : (
                    <p className="text-sm text-gray-400">
                      Rincian barang tidak tersedia (transaksi asal tidak ditemukan).
                    </p>
                  )}
                </div>

                {/* Bagian 3: Riwayat Cicilan */}
                <div>
                  <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">
                    Riwayat Pembayaran
                  </p>
                  {detail.payments.length > 0 ? (
                    <div className="space-y-2">
                      {detail.payments.map((p) => (
                        <div
                          key={p.id}
                          className="flex items-start justify-between rounded-lg border p-2.5 text-sm"
                        >
                          <div>
                            <p className="font-medium text-green-600">{formatRupiah(p.amount)}</p>
                            <p className="text-xs text-gray-500">
                              {formatDate(p.payment_date)}
                              {p.user_name ? ` • ${p.user_name}` : ''}
                            </p>
                            {p.notes && <p className="text-xs text-gray-400">{p.notes}</p>}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-sm text-gray-400">Belum ada pembayaran.</p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-400">Data piutang tidak ditemukan.</p>
            )}
          </div>
        </ScrollArea>
      </ActionModal>

      {detail && receiptOpen && (() => {
        const { cart, summary, discount, tax } = buildReceiptData(detail)
        return (
          <ReceiptPrint
            open={receiptOpen}
            onClose={() => setReceiptOpen(false)}
            checkoutData={{
              id: detail.transaction_id,
              transaction_code: detail.transaction_code,
              total_amount: detail.total_amount,
              payment_amount: 0,
              change_amount: 0,
              transaction_date: detail.created_at,
            }}
            cart={cart}
            summary={summary}
            discount={discount}
            tax={tax}
            paymentMethod="kredit"
            amountPaid={0}
            customerName={detail.customer_name}
            mode="reprint"
          />
        )
      })()}
    </>
  )
}
