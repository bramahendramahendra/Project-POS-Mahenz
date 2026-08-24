import { useEffect } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertTriangle } from 'lucide-react'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { useState } from 'react'

import { ActionModal } from '@/shared/components'
import { Checkbox } from '@/shared/components/ui/checkbox'
import { Label } from '@/shared/components/ui/label'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'

import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'
import { createPaymentSchema, type PaymentFormValues } from '../cashier.schema'
import { useCheckoutMutation, useCustomerCreditQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'
import type {
  CartSummary,
  CheckoutResponse,
  Discount,
  PaymentMethod,
  PaymentPayload,
  Tax,
} from '../cashier.types'
import type { CartItem } from '../cashier.types'
import { calcCartSummary, calcChange } from '../cashier.utils'
import { ReceiptPrint } from './ReceiptPrint'

interface PaymentModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const PAYMENT_METHODS: { key: PaymentMethod; label: string }[] = [
  { key: 'cash', label: 'Tunai' },
  { key: 'transfer', label: 'Transfer' },
  { key: 'qris', label: 'QRIS' },
  { key: 'card', label: 'Kartu' },
  { key: 'kredit', label: 'Hutang' },
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

function buildPayload(
  cart: CartItem[],
  summary: CartSummary,
  _discount: Discount,
  _tax: Tax,
  customerId: number | undefined,
  shiftId: number | undefined,
  paymentMethod: PaymentMethod,
  amountPaid: number,
  balanceUsed: number = 0,
): PaymentPayload {
  const isKredit = paymentMethod === 'kredit'
  return {
    customer_id: customerId,
    shift_id: shiftId,
    is_credit: isKredit,
    balance_used: balanceUsed > 0 ? balanceUsed : undefined,
    device_source: 'web',
    items: cart.map((i) => ({
      product_id: i.product_id,
      product_name: i.product_name,
      unit_id: i.unit_id,
      unit: i.unit_name,
      conversion_qty: Number(i.conversion_qty),
      quantity: Number(i.qty),
      price: Number(i.price),
      subtotal: Number(i.subtotal),
      discount_item: Number(i.discount_amount ?? 0),
    })),
    subtotal: summary.subtotal,
    discount: summary.discountAmount,
    tax: summary.taxAmount,
    total_amount: summary.grandTotal,
    payment_method: paymentMethod,
    payment_amount: isKredit ? 0 : amountPaid,
    change_amount: isKredit ? 0 : Math.max(0, amountPaid - summary.grandTotal),
  }
}

export function PaymentModal({ open, onOpenChange }: PaymentModalProps) {
  const [receiptOpen, setReceiptOpen] = useState(false)
  const [receiptData, setReceiptData] = useState<CheckoutResponse | null>(null)
  const [receiptSnapshot, setReceiptSnapshot] = useState<{
    cart: CartItem[]
    summary: CartSummary
    discount: Discount
    tax: Tax
    paymentMethod: PaymentMethod
    amountPaid: number
    customerName?: string
  } | null>(null)

  const { cart, discount, tax, selectedCustomer, clearCart, closePaymentModal } = useCashierStore()
  const summary = calcCartSummary(cart, discount, tax)
  const { mutate: checkout, isPending } = useCheckoutMutation()
  const { data: currentDrawer } = useCashDrawerCurrentQuery()
  const kasOpen = currentDrawer?.status === 'open'

  const { data: creditData } = useCustomerCreditQuery(selectedCustomer?.id ?? null)
  const customerCredit = creditData
  const hasCustomer = selectedCustomer !== null

  // Saldo pelanggan
  const [useSaldo, setUseSaldo] = useState(false)
  const customerBalance = customerCredit?.balance ?? 0
  const hasSaldo = hasCustomer && customerBalance > 0

  const {
    handleSubmit,
    control,
    setValue,
    reset,
    formState: { errors },
  } = useForm<PaymentFormValues>({
    resolver: (values, context, options) =>
      zodResolver(createPaymentSchema(summary.grandTotal, values.payment_method === 'kredit', effectiveTotal))(
        values,
        context,
        options
      ),
    defaultValues: { payment_method: 'cash', amount_paid: 0 },
  })

  const paymentMethod = useWatch({ control, name: 'payment_method' })
  const amountPaid = useWatch({ control, name: 'amount_paid' }) ?? 0
  const isKredit = paymentMethod === 'kredit'
  const change = calcChange(summary.grandTotal, amountPaid)
  const roundedOptions = buildRoundedOptions(summary.grandTotal)

  // Saldo calculations
  const balanceUsed = useSaldo ? Math.min(customerBalance, summary.grandTotal) : 0
  const effectiveTotal = summary.grandTotal - balanceUsed
  const changeAfterSaldo = calcChange(effectiveTotal, amountPaid)
  const sufficient = isKredit || effectiveTotal === 0 || amountPaid >= effectiveTotal

  const creditLimit = customerCredit?.credit_limit ?? 0
  const outstanding = customerCredit?.outstanding_amount ?? 0
  const remainingLimit = creditLimit > 0 ? creditLimit - outstanding : Infinity
  const exceedsLimit = isKredit && creditLimit > 0 && summary.grandTotal > remainingLimit

  useEffect(() => {
    if (open) {
      reset({ payment_method: 'cash', amount_paid: 0 })
      setUseSaldo(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  useEffect(() => {
    if (isKredit) setValue('amount_paid', 0)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isKredit])

  const onSubmit = (values: PaymentFormValues) => {
    const effectiveAmountPaid = isKredit ? 0 : (effectiveTotal === 0 ? 0 : (values.amount_paid ?? 0))
    const finalPaymentMethod = effectiveTotal === 0 && useSaldo ? 'balance' as PaymentMethod : values.payment_method
    const snapshot = {
      cart: [...cart],
      summary: { ...summary },
      discount: { ...discount },
      tax: { ...tax },
      paymentMethod: finalPaymentMethod,
      amountPaid: effectiveAmountPaid,
      customerName: selectedCustomer?.name,
    }
    const payload = buildPayload(
      cart, summary, discount, tax,
      selectedCustomer?.id, currentDrawer?.shift_id ?? undefined,
      finalPaymentMethod, effectiveAmountPaid, balanceUsed,
    )
    checkout(payload, {
      onSuccess: (data) => {
        const res = data as unknown as CheckoutResponse
        setReceiptSnapshot(snapshot)
        setReceiptData(res)
        closePaymentModal()
        setReceiptOpen(true)
      },
    })
  }

  const handleReceiptClose = () => {
    clearCart()
    setReceiptOpen(false)
    setReceiptData(null)
    setReceiptSnapshot(null)
  }

  return (
    <>
      <ActionModal
        open={open}
        onOpenChange={onOpenChange}
        title="Pembayaran"
        description="Form pembayaran transaksi"
        contentClassName="max-w-md"
        isLoading={isPending}
        asForm
        onFormSubmit={handleSubmit(onSubmit)}
        submitLabel={isPending ? 'Memproses...' : '✓ Proses Bayar'}
        submitDisabled={!sufficient || isPending || !kasOpen}
      >
            <ScrollArea style={{ maxHeight: '70vh' }}>
            <div className="px-6 py-4 space-y-5">
              {/* Guard kas belum buka */}
              {!kasOpen && (
                <div className="flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2.5 text-sm text-amber-700">
                  <AlertTriangle size={15} className="shrink-0 mt-0.5" />
                  <span>Kas belum dibuka. Hubungi admin/owner untuk membuka kas sebelum memproses transaksi.</span>
                </div>
              )}

              {/* Grand total */}
              <div className="rounded-lg bg-gray-50 px-4 py-3 text-center">
                <p className="text-xs text-gray-500 uppercase tracking-wide mb-0.5">Total Belanja</p>
                <p className="text-2xl font-bold text-gray-900">{formatRupiah(summary.grandTotal)}</p>
              </div>

              {/* Gunakan Saldo — hanya tampil jika pelanggan dipilih DAN saldo > 0 */}
              {hasSaldo && (
                <div className={`rounded-lg border-2 p-3 transition-colors ${useSaldo ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white'}`}>
                  <div className="flex items-center gap-2.5">
                    <Checkbox
                      id="use-saldo"
                      checked={useSaldo}
                      onCheckedChange={(v) => setUseSaldo(v === true)}
                    />
                    <label htmlFor="use-saldo" className="flex-1 cursor-pointer">
                      <span className="text-sm font-semibold text-gray-800">Gunakan Saldo</span>
                    </label>
                    <span className="text-sm font-bold text-blue-700">{formatRupiah(customerBalance)}</span>
                  </div>
                  {useSaldo && (
                    <div className="mt-2.5 pt-2 border-t border-blue-200 space-y-1">
                      <div className="flex justify-between text-xs text-gray-600">
                        <span>Saldo terpakai</span>
                        <span className="font-semibold text-gray-900">{formatRupiah(balanceUsed)}{balanceUsed >= customerBalance ? ' (habis)' : ''}</span>
                      </div>
                      {balanceUsed < summary.grandTotal && (
                        <div className="flex justify-between text-xs text-gray-600">
                          <span>Sisa saldo setelah</span>
                          <span className="font-semibold">{formatRupiah(customerBalance - balanceUsed)}</span>
                        </div>
                      )}
                      <div className="flex justify-between text-xs font-semibold text-red-600">
                        <span>Sisa bayar</span>
                        <span>{formatRupiah(effectiveTotal)}</span>
                      </div>
                    </div>
                  )}
                </div>
              )}

              {/* Jika saldo cukup, tidak perlu pilih metode */}
              {effectiveTotal === 0 && useSaldo ? (
                <div className="flex justify-between rounded-lg px-4 py-3 bg-green-50">
                  <span className="font-medium text-green-700">Status</span>
                  <span className="font-bold text-green-700">Lunas via Saldo ✓</span>
                </div>
              ) : (
                <>
                  {/* Payment method — hidden Hutang if no customer */}
                  <div className="space-y-2">
                    <Label className="text-sm">
                      Metode Pembayaran{effectiveTotal > 0 && effectiveTotal < summary.grandTotal ? ` (sisa ${formatRupiah(effectiveTotal)})` : ''}
                    </Label>
                    <div className="grid grid-cols-3 gap-2 sm:grid-cols-5">
                      {PAYMENT_METHODS.filter(({ key }) => {
                        // Hide Hutang if no customer
                        if (key === 'kredit' && !hasCustomer) return false
                        return true
                      }).map(({ key, label }) => {
                        const isKreditOption = key === 'kredit'
                        const disabled = isKreditOption && !hasCustomer
                        return (
                          <button
                            key={key}
                            type="button"
                            disabled={disabled}
                            onClick={() => setValue('payment_method', key)}
                            title={disabled ? 'Pilih pelanggan terlebih dahulu' : undefined}
                            className={`rounded-lg border-2 py-2 text-xs font-medium transition-all ${
                              paymentMethod === key
                                ? 'border-[#2c3e50] bg-[#2c3e50] text-white'
                                : disabled
                                  ? 'border-gray-100 text-gray-300 cursor-not-allowed'
                                  : 'border-gray-200 text-gray-600 hover:border-gray-300'
                            }`}
                          >
                            {label}
                          </button>
                        )
                      })}
                    </div>
                  </div>

                  {/* Kredit/Hutang info */}
                  {isKredit && selectedCustomer && (
                    <div className="rounded-lg border border-amber-100 bg-amber-50 p-3 space-y-1.5 text-sm">
                      <p className="font-medium text-amber-800 text-xs uppercase tracking-wide">
                        ⚠️ Piutang yang akan terbuat
                      </p>
                      <div className="flex justify-between text-amber-900 font-semibold">
                        <span>Jumlah hutang</span>
                        <span>{formatRupiah(effectiveTotal)}</span>
                      </div>
                      <p className="text-xs text-amber-700">Pelanggan: {selectedCustomer.name}</p>
                      {exceedsLimit && (
                        <div className="flex items-start gap-1.5 rounded-md bg-red-50 border border-red-200 px-2.5 py-2 text-xs text-red-700 mt-1">
                          <AlertTriangle size={13} className="shrink-0 mt-0.5" />
                          <span>Total melebihi sisa limit kredit. Transaksi tetap dapat diproses.</span>
                        </div>
                      )}
                    </div>
                  )}

                  {/* Amount paid — hide for kredit */}
                  {!isKredit && (
                    <>
                      <div className="space-y-1.5">
                        <Label htmlFor="amount-paid">Jumlah Bayar</Label>
                        <RupiahInput
                          id="amount-paid"
                          value={amountPaid}
                          onChange={(v) => setValue('amount_paid', v, { shouldValidate: true })}
                          className={`text-lg h-11 ${errors.amount_paid ? 'border-red-500' : ''}`}
                          autoFocus={false}
                        />
                        {errors.amount_paid && (
                          <p className="text-xs text-red-500">{errors.amount_paid.message}</p>
                        )}
                      </div>

                      {roundedOptions.length > 0 && (
                        <div className="flex gap-2">
                          {buildRoundedOptions(effectiveTotal).map((amt) => (
                            <button
                              key={amt}
                              type="button"
                              onClick={() => setValue('amount_paid', amt, { shouldValidate: true })}
                              className="flex-1 rounded-md border border-gray-200 py-1.5 text-xs text-gray-600 hover:bg-gray-50 transition-colors"
                            >
                              {formatRupiah(amt)}
                            </button>
                          ))}
                        </div>
                      )}

                      <div className={`flex justify-between rounded-lg px-4 py-3 ${sufficient ? 'bg-green-50' : 'bg-red-50'}`}>
                        <span className={`font-medium ${sufficient ? 'text-green-700' : 'text-red-600'}`}>
                          Kembalian
                        </span>
                        <span className={`font-bold text-lg ${sufficient ? 'text-green-700' : 'text-red-600'}`}>
                          {sufficient
                            ? formatRupiah(changeAfterSaldo)
                            : `Kurang ${formatRupiah(effectiveTotal - amountPaid)}`}
                        </span>
                      </div>
                    </>
                  )}
                </>
              )}
            </div>
            </ScrollArea>
      </ActionModal>

      {receiptOpen && receiptData && receiptSnapshot && (
        <ReceiptPrint
          open={receiptOpen}
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
