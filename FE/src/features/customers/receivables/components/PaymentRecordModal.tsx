import { useEffect, useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { Label } from '@/shared/components/ui/label'
import { formatRupiah, todayStr } from '@/shared/utils'

import { useAddPaymentMutation } from '../receivables.api'
import type { Receivable } from '../receivables.types'
import { createPaymentSchema, type PaymentFormValues } from '../receivables.schema'

interface PaymentRecordModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  receivable: Receivable | null
}

export function PaymentRecordModal({ open, onOpenChange, receivable }: PaymentRecordModalProps) {
  const remaining = receivable?.remaining_amount ?? 0
  const receivableId = receivable?.id ?? 0

  const [confirmOpen, setConfirmOpen] = useState(false)
  const [pendingValues, setPendingValues] = useState<PaymentFormValues | null>(null)

  const paymentSchema = createPaymentSchema(remaining)

  const { mutate: addPayment, isPending } = useAddPaymentMutation(receivableId)

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    control,
    formState: { errors },
  } = useForm<PaymentFormValues>({
    resolver: zodResolver(paymentSchema),
    defaultValues: { amount: 0, payment_date: todayStr(), notes: '' },
  })

  useEffect(() => {
    if (!open) return
    reset({ amount: 0, payment_date: todayStr(), notes: '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const handleClose = () => {
    setConfirmOpen(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: PaymentFormValues) => {
    setPendingValues(values)
    setConfirmOpen(true)
  }

  const handleConfirm = () => {
    if (!pendingValues) return
    addPayment(
      {
        amount: pendingValues.amount,
        payment_date: pendingValues.payment_date,
        notes: pendingValues.notes || undefined,
      },
      {
        onSuccess: () => handleClose(),
      }
    )
  }

  return (
    <>
      <FormModal
        open={open && !confirmOpen}
        onOpenChange={(val) => {
          if (!val && !confirmOpen) handleClose()
        }}
        title="Catat Pembayaran Piutang"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan Pembayaran"
      >
        {receivable && (
          <div className="space-y-4">
            <div className="rounded-lg bg-gray-50 p-3 space-y-1.5 text-sm">
              <div className="flex justify-between text-gray-600">
                <span>Pelanggan</span>
                <span className="font-medium">{receivable.customer_name}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Total Piutang</span>
                <span>{formatRupiah(receivable.total_amount)}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Sudah Dibayar</span>
                <span className="text-green-600">{formatRupiah(receivable.paid_amount)}</span>
              </div>
              <div className="flex justify-between border-t pt-1.5 font-semibold">
                <span>Sisa</span>
                <span className="text-red-600">{formatRupiah(remaining)}</span>
              </div>
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <Label htmlFor="pay-amount">
                  Jumlah Bayar <span className="text-red-500">*</span>
                </Label>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="h-6 text-xs"
                  onClick={() => setValue('amount', remaining, { shouldValidate: true })}
                >
                  Bayar Lunas
                </Button>
              </div>
              <Controller
                control={control}
                name="amount"
                render={({ field }) => (
                  <RupiahInput
                    id="pay-amount"
                    value={field.value}
                    onChange={field.onChange}
                    className={errors.amount ? 'border-red-500' : ''}
                  />
                )}
              />
              {errors.amount && <p className="text-xs text-red-500">{errors.amount.message}</p>}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="pay-date">
                Tanggal Pembayaran <span className="text-red-500">*</span>
              </Label>
              <Input
                id="pay-date"
                type="date"
                {...register('payment_date')}
                className={errors.payment_date ? 'border-red-500' : ''}
              />
              {errors.payment_date && (
                <p className="text-xs text-red-500">{errors.payment_date.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="pay-notes">Catatan</Label>
              <Input id="pay-notes" {...register('notes')} placeholder="Catatan (opsional)" />
            </div>
          </div>
        )}
      </FormModal>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Konfirmasi Pembayaran"
        description={`Catat pembayaran ${formatRupiah(pendingValues?.amount ?? 0)} untuk piutang ${receivable?.customer_name ?? ''}?`}
        confirmLabel="Ya, Simpan"
        onConfirm={handleConfirm}
        isLoading={isPending}
      />
    </>
  )
}
