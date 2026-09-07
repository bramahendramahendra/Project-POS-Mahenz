/**
 * Payment modal for backdate cashier — adds time picker (HH:mm) and shows
 * readonly date/kasir info from the active backdate cash drawer.
 */
import { useEffect, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { AlertTriangle, CalendarClock, User } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Checkbox } from '@/shared/components/ui/checkbox'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'

import { useCustomerCreditQuery } from '../../cashier/cashier.api'
import { calcCartSummary } from '../../cashier/cashier.utils'
import type { PaymentMethod } from '../../cashier/cashier.types'
import { useBackdateCashierStore } from '../backdate-cashier.store'
import { useBackdateCheckoutMutation } from '../backdate-cashier.api'
import type { BackdateCashDrawer } from '../../backdate-cash/backdate-cash.types'
import type { BackdatePaymentPayload } from '../backdate-cashier.types'

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

const paymentMethodEnum = z.enum(['cash', 'transfer', 'qris', 'card', 'kredit', 'balance'] as const)

function createBackdatePaymentSchema(effectiveTotal: number, isKredit: boolean) {
  const base = z.object({
    payment_method: paymentMethodEnum,
    amount_paid: z.number().default(0),
    transaction_time: z.string().min(1, 'Jam transaksi wajib diisi'),
  })
  if (effectiveTotal === 0 || isKredit) return base
  return base.refine((d) => d.amount_paid >= effectiveTotal, {
    message: 'Jumlah bayar kurang dari total',
    path: ['amount_paid'],
  })
}

type FormValues = { payment_method: PaymentMethod; amount_paid: number; transaction_time: string }

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  backdateDrawer: BackdateCashDrawer | null | undefined
}

export function BackdatePaymentModal({ open, onOpenChange, backdateDrawer }: Props) {
  const { cart, discount, tax, selectedCustomer, clearCart, closePaymentModal } = useBackdateCashierStore()
  const summary = calcCartSummary(cart, discount, tax)
  const { mutate: checkout, isPending } = useBackdateCheckoutMutation()

  const { data: creditData } = useCustomerCreditQuery(selectedCustomer?.id ?? null)
  const hasCustomer = selectedCustomer !== null
  const customerBalance = creditData?.balance ?? 0
  const hasSaldo = hasCustomer && customerBalance > 0
  const [useSaldo, setUseSaldo] = useState(false)
  const balanceUsed = useSaldo ? Math.min(customerBalance, summary.grandTotal) : 0
  const effectiveTotal = summary.grandTotal - balanceUsed

  const kasOpen = backdateDrawer != null && backdateDrawer.status === 'open'
  const drawerDate = backdateDrawer ? new Date(backdateDrawer.open_time).toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' }) : ''

  const {
    handleSubmit,
    control,
    setValue,
    register,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: (values, context, options) =>
      zodResolver(createBackdatePaymentSchema(effectiveTotal, values.payment_method === 'kredit'))(values, context, options),
    defaultValues: { payment_method: 'cash', amount_paid: 0, transaction_time: '' },
  })

  const paymentMethod = useWatch({ control, name: 'payment_method' })
  const amountPaid = useWatch({ control, name: 'amount_paid' }) ?? 0
  const isKredit = paymentMethod === 'kredit'
  const roundedOptions = buildRoundedOptions(effectiveTotal > 0 ? effectiveTotal : summary.grandTotal)
  const sufficient = isKredit || effectiveTotal === 0 || amountPaid >= effectiveTotal

  useEffect(() => {
    // Reset form + ephemeral UI state saat modal dibuka. set-state di effect
    // memang disengaja di sini (menyinkronkan state form ke kondisi "baru dibuka"),
    // bukan cascading render yang tidak disengaja.
    if (open) {
      reset({ payment_method: 'cash', amount_paid: 0, transaction_time: '' })
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setUseSaldo(false)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  useEffect(() => { if (isKredit) setValue('amount_paid', 0) }, [isKredit, setValue])

  const onSubmit = (values: FormValues) => {
    if (!backdateDrawer) return
    const finalPaymentMethod = effectiveTotal === 0 && useSaldo ? 'balance' as PaymentMethod : values.payment_method
    const effectiveAmountPaid = isKredit ? 0 : (effectiveTotal === 0 ? 0 : values.amount_paid)

    const payload: BackdatePaymentPayload = {
      transaction_time: values.transaction_time,
      shift_id: backdateDrawer.shift_id ?? undefined,
      customer_id: selectedCustomer?.id,
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
      payment_method: finalPaymentMethod,
      payment_amount: effectiveAmountPaid,
      change_amount: isKredit ? 0 : Math.max(0, effectiveAmountPaid - effectiveTotal),
    }

    checkout(payload, {
      onSuccess: () => {
        clearCart()
        closePaymentModal()
      },
    })
  }

  return (
    <ActionModal
      open={open}
      onOpenChange={onOpenChange}
      title="Pembayaran Historis"
      description="Form pembayaran transaksi historis"
      contentClassName="max-w-md"
      isLoading={isPending}
      asForm
      onFormSubmit={handleSubmit(onSubmit)}
      submitLabel={isPending ? 'Memproses...' : '✓ Proses Bayar'}
      submitDisabled={!sufficient || isPending || !kasOpen}
    >
      <ScrollArea style={{ maxHeight: '70vh' }}>
        <div className="px-6 py-4 space-y-5">
          {/* Guard */}
          {!kasOpen && (
            <div className="flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2.5 text-sm text-amber-700">
              <AlertTriangle size={15} className="shrink-0 mt-0.5" />
              <span>Kas historis belum dibuka. Buka di menu Kas Historis terlebih dahulu.</span>
            </div>
          )}

          {/* Info tanggal + kasir */}
          {backdateDrawer && (
            <div className="rounded-lg bg-blue-50 border border-blue-200 px-4 py-2.5 space-y-1">
              <div className="flex items-center gap-2 text-sm text-blue-700">
                <CalendarClock size={14} />
                <span className="font-medium">Tanggal: {drawerDate}</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-blue-600">
                <User size={14} />
                <span>Kasir: {backdateDrawer.user_name}</span>
              </div>
            </div>
          )}

          {/* Jam transaksi */}
          <div className="space-y-1.5">
            <Label>
              Jam Transaksi <span className="text-red-500">*</span>
            </Label>
            <Input type="time" {...register('transaction_time')} className={errors.transaction_time ? 'border-red-500' : ''} />
            {errors.transaction_time && <p className="text-xs text-red-500">{errors.transaction_time.message}</p>}
          </div>

          {/* Grand total */}
          <div className="rounded-lg bg-gray-50 px-4 py-3 text-center">
            <p className="text-xs text-gray-500 uppercase tracking-wide mb-0.5">Total Belanja</p>
            <p className="text-2xl font-bold text-gray-900">{formatRupiah(summary.grandTotal)}</p>
          </div>

          {/* Saldo */}
          {hasSaldo && (
            <div className={`rounded-lg border-2 p-3 transition-colors ${useSaldo ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white'}`}>
              <div className="flex items-center gap-2.5">
                <Checkbox id="use-saldo-bd" checked={useSaldo} onCheckedChange={(v) => setUseSaldo(v === true)} />
                <label htmlFor="use-saldo-bd" className="flex-1 cursor-pointer">
                  <span className="text-sm font-semibold text-gray-800">Gunakan Saldo</span>
                </label>
                <span className="text-sm font-bold text-blue-700">{formatRupiah(customerBalance)}</span>
              </div>
              {useSaldo && effectiveTotal > 0 && (
                <div className="mt-2 pt-2 border-t border-blue-200 flex justify-between text-xs font-semibold text-red-600">
                  <span>Sisa bayar</span>
                  <span>{formatRupiah(effectiveTotal)}</span>
                </div>
              )}
            </div>
          )}

          {/* Jika saldo cukup */}
          {effectiveTotal === 0 && useSaldo ? (
            <div className="flex justify-between rounded-lg px-4 py-3 bg-green-50">
              <span className="font-medium text-green-700">Status</span>
              <span className="font-bold text-green-700">Lunas via Saldo ✓</span>
            </div>
          ) : (
            <>
              {/* Payment method */}
              <div className="space-y-2">
                <Label className="text-sm">Metode Pembayaran</Label>
                <div className="grid grid-cols-3 gap-2 sm:grid-cols-5">
                  {PAYMENT_METHODS.filter(({ key }) => key !== 'kredit' || hasCustomer).map(({ key, label }) => (
                    <button
                      key={key}
                      type="button"
                      onClick={() => setValue('payment_method', key)}
                      className={`rounded-lg border-2 py-2 text-xs font-medium transition-all ${
                        paymentMethod === key ? 'border-[#2c3e50] bg-[#2c3e50] text-white' : 'border-gray-200 text-gray-600 hover:border-gray-300'
                      }`}
                    >
                      {label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Amount */}
              {!isKredit && effectiveTotal > 0 && (
                <div className="space-y-2">
                  <Label className="text-sm">Jumlah Bayar</Label>
                  <RupiahInput value={amountPaid} onChange={(v) => setValue('amount_paid', v)} className={errors.amount_paid ? 'border-red-500' : ''} />
                  {errors.amount_paid && <p className="text-xs text-red-500">{errors.amount_paid.message}</p>}

                  {/* Quick amounts */}
                  <div className="flex flex-wrap gap-2">
                    <button type="button" onClick={() => setValue('amount_paid', effectiveTotal)} className="rounded-md border px-3 py-1.5 text-xs font-medium text-blue-700 border-blue-200 bg-blue-50 hover:bg-blue-100">
                      Uang Pas
                    </button>
                    {roundedOptions.map((val) => (
                      <button key={val} type="button" onClick={() => setValue('amount_paid', val)} className="rounded-md border px-3 py-1.5 text-xs font-medium text-gray-700 border-gray-200 hover:bg-gray-50">
                        {formatRupiah(val)}
                      </button>
                    ))}
                  </div>

                  {/* Change */}
                  {amountPaid > effectiveTotal && (
                    <div className="flex justify-between rounded-lg bg-green-50 px-4 py-2.5">
                      <span className="text-sm text-green-700">Kembalian</span>
                      <span className="text-sm font-bold text-green-800">{formatRupiah(amountPaid - effectiveTotal)}</span>
                    </div>
                  )}
                </div>
              )}
            </>
          )}
        </div>
      </ScrollArea>
    </ActionModal>
  )
}
