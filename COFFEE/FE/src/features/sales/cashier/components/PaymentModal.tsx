import { useEffect, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertTriangle } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'

import { createPaymentSchema, type PaymentFormValues } from '../cashier.schema'
import { useCheckoutMutation, useCustomerCreditQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'
import type { CartItem, CartSummary, CheckoutResponse, Discount, PaymentMethod, PaymentPayload, Tax } from '../cashier.types'
import { ORDER_TYPE_LABEL, PAYMENT_METHOD_LABEL } from '../cashier.types'
import { calcChange } from '../cashier.utils'
import { usePayBillMutation } from '../openBill.api'
import { useActiveBill } from '../useActiveBill'
import { transactionReceipt } from '@/features/sales/transactions/transactions.utils'
import { ReceiptPrint } from './ReceiptPrint'

interface PaymentModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const PAYMENT_METHODS: { key: PaymentFormValues['payment_method']; label: string }[] = [
  { key: 'cash', label: PAYMENT_METHOD_LABEL.cash },
  { key: 'transfer', label: PAYMENT_METHOD_LABEL.transfer },
  { key: 'qris', label: PAYMENT_METHOD_LABEL.qris },
  { key: 'card', label: PAYMENT_METHOD_LABEL.card },
  { key: 'kredit', label: PAYMENT_METHOD_LABEL.kredit },
]

function buildRoundedOptions(grandTotal: number): number[] {
  const thresholds = [5_000, 10_000, 25_000, 50_000, 100_000, 500_000]
  const results: number[] = []
  for (const t of thresholds) {
    const rounded = Math.ceil(grandTotal / t) * t
    if (rounded > grandTotal && !results.includes(rounded)) {
      results.push(rounded)
      if (results.length >= 3) break
    }
  }
  return results
}

interface ReceiptSnapshot {
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
}

/** Pembayaran pesanan bayar sekarang (tahap 4.7b). Harga dan total final dihitung server. */
export function PaymentModal({ open, onOpenChange }: PaymentModalProps) {
  const [receiptData, setReceiptData] = useState<CheckoutResponse | null>(null)
  const [receiptSnapshot, setReceiptSnapshot] = useState<ReceiptSnapshot | null>(null)

  const { cart, discount, tax, selectedCustomer, orderType, table, orderName, clearCart, closePaymentModal } = useCashierStore()
  // Pesanan terbuka (4.8): total dari item tersimpan; pesanan baru: dari keranjang.
  const { bill, billOpen, summary } = useActiveBill()
  const { mutate: payBill, isPending: isPayingBill } = usePayBillMutation()
  const { mutate: checkout, isPending: isCheckingOut } = useCheckoutMutation()
  const isPending = isCheckingOut || isPayingBill
  const { data: currentDrawer } = useCashDrawerCurrentQuery()
  const kasOpen = currentDrawer?.status === 'open'

  const { data: customerCredit } = useCustomerCreditQuery(selectedCustomer?.id ?? null)
  const hasCustomer = selectedCustomer !== null

  const { handleSubmit, control, setValue, reset, formState: { errors } } = useForm<PaymentFormValues>({
    resolver: (values, context, options) =>
      zodResolver(createPaymentSchema(summary.grandTotal, values.payment_method === 'kredit'))(values, context, options),
    defaultValues: { payment_method: 'cash', amount_paid: 0 },
  })

  const paymentMethod = useWatch({ control, name: 'payment_method' })
  const amountPaid = useWatch({ control, name: 'amount_paid' }) ?? 0
  const isKredit = paymentMethod === 'kredit'
  const change = calcChange(summary.grandTotal, amountPaid)
  const sufficient = isKredit || summary.grandTotal === 0 || amountPaid >= summary.grandTotal

  const creditLimit = customerCredit?.credit_limit ?? 0
  const outstanding = customerCredit?.outstanding_amount ?? 0
  const exceedsLimit = isKredit && creditLimit > 0 && summary.grandTotal > creditLimit - outstanding

  useEffect(() => {
    if (open) reset({ payment_method: 'cash', amount_paid: 0 })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  useEffect(() => {
    if (isKredit) setValue('amount_paid', 0)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isKredit])

  const onSubmit = (values: PaymentFormValues) => {
    const paid = isKredit ? 0 : (values.amount_paid ?? 0)
    if (billOpen && bill) {
      payBill(
        {
          id: bill.id,
          shift_id: currentDrawer?.shift_id ?? undefined,
          discount: summary.discountAmount,
          tax: summary.taxAmount,
          payment_method: values.payment_method,
          payment_amount: paid,
          customer_id: selectedCustomer?.id,
          is_credit: isKredit,
        },
        {
          onSuccess: (res) => {
            const r = transactionReceipt(res)
            setReceiptSnapshot({ cart: r.cart, summary: r.summary, discount: r.discount, tax: { ...tax, amount: res.tax }, paymentMethod: values.payment_method, amountPaid: paid, customerName: selectedCustomer?.name })
            setReceiptData(r.checkoutData)
            closePaymentModal()
          },
        },
      )
      return
    }
    const payload: PaymentPayload = {
      order_type: orderType,
      table_id: orderType === 'dine_in' ? table?.id : undefined,
      order_name: orderName.trim() || undefined,
      customer_id: selectedCustomer?.id,
      shift_id: currentDrawer?.shift_id ?? undefined,
      is_credit: isKredit,
      device_source: 'web',
      items: cart.map((i) => ({
        variant_id: i.variant_id,
        quantity: i.qty,
        modifier_option_ids: i.modifiers.map((m) => m.option_id),
        notes: i.notes,
        discount_item: i.discount_amount ?? 0,
      })),
      discount: summary.discountAmount,
      tax: summary.taxAmount,
      payment_method: values.payment_method,
      payment_amount: paid,
    }
    const snapshot: ReceiptSnapshot = {
      cart: [...cart],
      summary: { ...summary },
      discount: { ...discount },
      tax: { ...tax },
      paymentMethod: values.payment_method,
      amountPaid: paid,
      customerName: selectedCustomer?.name,
    }
    checkout(payload, {
      onSuccess: (res) => {
        // Total dari server dipakai di struk (sumber kebenaran harga).
        setReceiptSnapshot({
          ...snapshot,
          summary: { ...snapshot.summary, grandTotal: res.total_amount },
        })
        setReceiptData(res)
        closePaymentModal()
      },
    })
  }

  const handleReceiptClose = () => {
    clearCart()
    setReceiptData(null)
    setReceiptSnapshot(null)
  }

  return (
    <>
      <ActionModal
        open={open}
        onOpenChange={onOpenChange}
        title="Pembayaran"
        description={billOpen && bill ? `Pesanan terbuka antrean ${bill.queue_label}${bill.table_name ? ` · ${bill.table_name}` : ''}` : `${ORDER_TYPE_LABEL[orderType]}${table && orderType === 'dine_in' ? ` · ${table.name}` : ''}${orderName.trim() ? ` · ${orderName.trim()}` : ''}`}
        contentClassName="max-w-md"
        isLoading={isPending}
        asForm
        onFormSubmit={handleSubmit(onSubmit)}
        submitLabel={isPending ? 'Memproses...' : 'Proses Bayar'}
        submitDisabled={!sufficient || isPending || !kasOpen}
      >
        <ScrollArea style={{ maxHeight: '70vh' }}>
          <div className="space-y-5 px-6 py-4 touch:[&_button:not([role=checkbox])]:min-h-10">
            {!kasOpen && (
              <div className="flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2.5 text-sm text-amber-700">
                <AlertTriangle size={15} className="mt-0.5 shrink-0" />
                <span>Kas belum dibuka. Buka kas terlebih dahulu sebelum memproses transaksi.</span>
              </div>
            )}

            <div className="rounded-lg bg-gray-50 px-4 py-3 text-center">
              <p className="mb-0.5 text-xs uppercase tracking-wide text-gray-500">Total</p>
              <p className="text-2xl font-bold text-gray-900">{formatRupiah(summary.grandTotal)}</p>
            </div>

            <div className="space-y-2">
              <Label className="text-sm">Metode Pembayaran</Label>
              <div className="grid grid-cols-3 gap-2 sm:grid-cols-5">
                {PAYMENT_METHODS.filter(({ key }) => key !== 'kredit' || hasCustomer).map(({ key, label }) => (
                  <button
                    key={key}
                    type="button"
                    aria-pressed={paymentMethod === key}
                    onClick={() => setValue('payment_method', key)}
                    className={`rounded-lg border-2 py-2 text-xs font-medium transition-all ${
                      paymentMethod === key ? 'border-[#2c3e50] bg-[#2c3e50] text-white' : 'border-gray-200 text-gray-600 hover:border-gray-300'
                    }`}
                  >
                    {label}
                  </button>
                ))}
              </div>
              {!hasCustomer && <p className="text-xs text-gray-500">Pilih pelanggan untuk mencatat sebagai utang.</p>}
            </div>

            {isKredit && selectedCustomer && (
              <div className="space-y-1.5 rounded-lg border border-amber-100 bg-amber-50 p-3 text-sm">
                <p className="text-xs font-medium uppercase tracking-wide text-amber-800">Piutang yang akan terbuat</p>
                <div className="flex justify-between font-semibold text-amber-900">
                  <span>Jumlah utang</span>
                  <span>{formatRupiah(summary.grandTotal)}</span>
                </div>
                <p className="text-xs text-amber-700">Pelanggan: {selectedCustomer.name}</p>
                {exceedsLimit && (
                  <div className="mt-1 flex items-start gap-1.5 rounded-md border border-red-200 bg-red-50 px-2.5 py-2 text-xs text-red-700">
                    <AlertTriangle size={13} className="mt-0.5 shrink-0" />
                    <span>Total melebihi sisa limit kredit. Transaksi tetap dapat diproses.</span>
                  </div>
                )}
              </div>
            )}

            {!isKredit && (
              <>
                <div className="space-y-1.5">
                  <Label htmlFor="amount-paid">Jumlah Bayar</Label>
                  <RupiahInput
                    id="amount-paid"
                    value={amountPaid}
                    onChange={(v) => setValue('amount_paid', v, { shouldValidate: true })}
                    className={`h-11 text-lg ${errors.amount_paid ? 'border-red-500' : ''}`}
                    autoFocus={false}
                  />
                  {errors.amount_paid && <p className="text-xs text-red-500">{errors.amount_paid.message}</p>}
                </div>

                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => setValue('amount_paid', summary.grandTotal, { shouldValidate: true })}
                    className="rounded-md border border-blue-200 bg-blue-50 px-3 py-1.5 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100"
                  >
                    Uang Pas
                  </button>
                  {buildRoundedOptions(summary.grandTotal).map((amt) => (
                    <button
                      key={amt}
                      type="button"
                      onClick={() => setValue('amount_paid', amt, { shouldValidate: true })}
                      className="flex-1 rounded-md border border-gray-200 py-1.5 text-xs text-gray-600 transition-colors hover:bg-gray-50"
                    >
                      {formatRupiah(amt)}
                    </button>
                  ))}
                </div>

                <div className={`flex justify-between rounded-lg px-4 py-3 ${sufficient ? 'bg-green-50' : 'bg-red-50'}`}>
                  <span className={`font-medium ${sufficient ? 'text-green-700' : 'text-red-600'}`}>Kembalian</span>
                  <span className={`text-lg font-bold ${sufficient ? 'text-green-700' : 'text-red-600'}`}>
                    {sufficient ? formatRupiah(change) : `Kurang ${formatRupiah(summary.grandTotal - amountPaid)}`}
                  </span>
                </div>
              </>
            )}
          </div>
        </ScrollArea>
      </ActionModal>

      {receiptData && receiptSnapshot && (
        <ReceiptPrint
          open
          onClose={handleReceiptClose}
          checkoutData={receiptData}
          cart={receiptSnapshot.cart}
          summary={receiptSnapshot.summary}
          discount={receiptSnapshot.discount}
          tax={receiptSnapshot.tax}
          paymentMethod={receiptSnapshot.paymentMethod}
          amountPaid={receiptSnapshot.amountPaid}
          customerName={receiptSnapshot.customerName}
        />
      )}
    </>
  )
}
