import { useState } from 'react'
import { Printer } from 'lucide-react'
import { toast } from 'sonner'

import { ConfirmDialog, RoleGuard, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { formatRupiah } from '@/shared/utils'
import { ReceiptPrint } from '@/features/sales/cashier/components/ReceiptPrint'
import type { CartItem, CartSummary, Discount, Tax } from '@/features/sales/cashier'

import { useTransactionDetailQuery, useVoidTransactionMutation } from '../transactions.api'
import type { Transaction } from '../transactions.types'
import { PAYMENT_LABELS, formatDateTimeLong } from '../transactions.utils'

interface TransactionDetailModalProps {
  transactionId: number | null
  onClose: () => void
}

function buildReceiptData(t: Transaction): {
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
} {
  const cart: CartItem[] = t.items.map((item) => ({
    product_id: item.product_id,
    product_name: item.product_name,
    unit_id: item.id,
    unit_name: item.unit,
    conversion_qty: item.conversion_qty ?? 1,
    qty: item.quantity,
    price: item.price,
    subtotal: item.subtotal,
    discount_amount: item.discount_item > 0 ? item.discount_item : undefined,
  }))

  const summary: CartSummary = {
    subtotal: t.subtotal,
    discountAmount: t.discount,
    taxAmount: t.tax,
    grandTotal: t.total_amount,
  }

  const discount: Discount = { type: 'amount', value: t.discount, amount: t.discount }
  const tax: Tax = { percent: 0, amount: t.tax }

  return { cart, summary, discount, tax }
}

export function TransactionDetailModal({ transactionId, onClose }: TransactionDetailModalProps) {
  const [voidConfirmOpen, setVoidConfirmOpen] = useState(false)
  const [receiptOpen, setReceiptOpen] = useState(false)

  const { data: transaction, isLoading } = useTransactionDetailQuery(transactionId ?? 0)
  const { mutate: voidTransaction, isPending: isVoiding } = useVoidTransactionMutation()

  const open = transactionId !== null

  const handleVoid = () => {
    if (!transactionId) return
    voidTransaction(transactionId, {
      onSuccess: () => {
        toast.success('Transaksi berhasil dibatalkan')
        setVoidConfirmOpen(false)
        onClose()
      },
    })
  }

  return (
    <>
      <Dialog open={open} onOpenChange={(val) => { if (!val) onClose() }}>
        <DialogContent className="flex flex-col gap-0 p-0 max-w-lg max-h-[90vh]">
          <DialogHeader className="border-b px-5 py-3 shrink-0">
            <DialogTitle>Detail Transaksi</DialogTitle>
            {transaction && (
              <DialogDescription className="font-mono text-xs">
                {transaction.transaction_code}
              </DialogDescription>
            )}
          </DialogHeader>

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
                {/* Meta info */}
                <div className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                  <div>
                    <p className="text-xs text-gray-500">Tanggal</p>
                    <p className="font-medium">{formatDateTimeLong(transaction.transaction_date)}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Status</p>
                    <StatusBadge
                      status={transaction.status === 'completed' ? 'success' : 'error'}
                    />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Kasir</p>
                    <p className="font-medium">{transaction.kasir_name}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Pelanggan</p>
                    <p className="font-medium">{transaction.customer_name ?? '—'}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Metode Bayar</p>
                    <p className="font-medium">{PAYMENT_LABELS[transaction.payment_method]}</p>
                  </div>
                </div>

                {/* Items table */}
                <div>
                  <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">
                    Item
                  </p>
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
                      {transaction.items.map((item, i) => (
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
                </div>

                {/* Summary */}
                <div className="rounded-lg bg-gray-50 p-4 space-y-1.5 text-sm">
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

          {transaction && (
            <div className="border-t px-5 py-3 flex items-center gap-2 shrink-0">
              <Button
                variant="outline"
                size="sm"
                className="gap-1"
                onClick={() => setReceiptOpen(true)}
              >
                <Printer size={14} />
                Cetak Ulang
              </Button>
              <RoleGuard menuKey="penjualan.transaksi" action="can_delete">
                {transaction.status === 'completed' && (
                  <Button
                    variant="destructive"
                    size="sm"
                    className="ml-auto"
                    onClick={() => setVoidConfirmOpen(true)}
                  >
                    Void Transaksi
                  </Button>
                )}
              </RoleGuard>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={voidConfirmOpen}
        onOpenChange={(val) => { if (!val) setVoidConfirmOpen(false) }}
        title="Batalkan Transaksi"
        description="Transaksi yang dibatalkan tidak dapat dikembalikan. Yakin ingin melanjutkan?"
        confirmLabel="Ya, Batalkan"
        variant="destructive"
        isLoading={isVoiding}
        onConfirm={handleVoid}
      />

      {transaction && receiptOpen && (() => {
        const { cart, summary, discount, tax } = buildReceiptData(transaction)
        return (
          <ReceiptPrint
            open={receiptOpen}
            onClose={() => setReceiptOpen(false)}
            checkoutData={{
              id: transaction.id,
              transaction_code: transaction.transaction_code,
              total_amount: transaction.total_amount,
              payment_amount: transaction.payment_amount,
              change_amount: transaction.change_amount,
              transaction_date: transaction.transaction_date,
            }}
            cart={cart}
            summary={summary}
            discount={discount}
            tax={tax}
            paymentMethod={transaction.payment_method}
            amountPaid={transaction.payment_amount}
            customerName={transaction.customer_name}
          />
        )
      })()}
    </>
  )
}
