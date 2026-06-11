import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { formatRupiah } from '@/shared/utils'

import { useCloseShiftMutation } from '../shifts.api'
import { closeShiftSchema, type CloseShiftFormValues } from '../shifts.schema'
import type { Shift } from '../shifts.types'

interface CloseShiftModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  shift: Shift | null
}

export function CloseShiftModal({ open, onOpenChange, shift }: CloseShiftModalProps) {
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [pendingValues, setPendingValues] = useState<CloseShiftFormValues | null>(null)

  const { mutate: closeShift, isPending } = useCloseShiftMutation()

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<CloseShiftFormValues>({
    resolver: zodResolver(closeShiftSchema),
    defaultValues: { closing_balance: 0, notes: '' },
  })

  useEffect(() => {
    if (!open) reset({ closing_balance: 0, notes: '' })
  }, [open, reset])

  const closingBalance = watch('closing_balance') || 0
  const expected = (shift?.opening_balance ?? 0) + (shift?.total_revenue ?? 0)
  const selisih = closingBalance - expected

  const onSubmit = (values: CloseShiftFormValues) => {
    setPendingValues(values)
    setConfirmOpen(true)
  }

  const handleConfirm = () => {
    if (!shift || !pendingValues) return
    closeShift(
      {
        id: shift.id,
        payload: { closing_balance: pendingValues.closing_balance, notes: pendingValues.notes || undefined },
      },
      {
        onSuccess: () => {
          setConfirmOpen(false)
          onOpenChange(false)
        },
      }
    )
  }

  return (
    <>
      <FormModal
        open={open}
        onOpenChange={onOpenChange}
        title="Tutup Shift"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Tutup Shift"
      >
        {shift && (
          <div className="space-y-4">
            <div className="rounded-lg bg-gray-50 p-3 space-y-1.5 text-sm">
              <div className="flex justify-between text-gray-600">
                <span>Kasir</span>
                <span className="font-medium">{shift.kasir_name}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Total Transaksi</span>
                <span>{shift.total_transactions}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Total Pendapatan</span>
                <span className="text-green-600 font-medium">{formatRupiah(shift.total_revenue)}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Modal Awal</span>
                <span>{formatRupiah(shift.opening_balance)}</span>
              </div>
              <div className="flex justify-between border-t pt-1.5 font-semibold text-gray-700">
                <span>Ekspektasi Kas</span>
                <span>{formatRupiah(expected)}</span>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="closing-balance">
                Uang di Laci Akhir Shift <span className="text-red-500">*</span>
              </Label>
              <Input
                id="closing-balance"
                type="number"
                min={0}
                {...register('closing_balance', { valueAsNumber: true })}
                className={errors.closing_balance ? 'border-red-500' : ''}
                placeholder="0"
              />
              {errors.closing_balance && (
                <p className="text-xs text-red-500">{errors.closing_balance.message}</p>
              )}
            </div>

            <div
              className={`rounded-lg p-3 text-sm flex justify-between font-semibold ${
                selisih < 0 ? 'bg-red-50 text-red-600' : 'bg-green-50 text-green-700'
              }`}
            >
              <span>Selisih</span>
              <span>
                {selisih >= 0 ? '+' : ''}
                {formatRupiah(selisih)}
              </span>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="close-notes">Catatan</Label>
              <Input id="close-notes" {...register('notes')} placeholder="Catatan (opsional)" />
            </div>
          </div>
        )}
      </FormModal>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Tutup Shift"
        description={`Shift kasir ${shift?.kasir_name ?? ''} akan ditutup. Tindakan ini tidak dapat dibatalkan. Lanjutkan?`}
        confirmLabel="Ya, Tutup Shift"
        variant="destructive"
        onConfirm={handleConfirm}
        isLoading={isPending}
      />
    </>
  )
}
