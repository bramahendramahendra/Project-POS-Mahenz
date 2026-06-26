import { useEffect, useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'

import { useCloseCashDrawerMutation } from '../cash-drawer.api'
import { closeCashDrawerSchema, type CloseCashDrawerFormValues } from '../cash-drawer.schema'

const defaultValues: CloseCashDrawerFormValues = {
  closing_balance: 0,
  notes: '',
}

interface CloseCashDrawerModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  cashDrawerId: number | null
}

export function CloseCashDrawerModal({ open, onOpenChange, cashDrawerId }: CloseCashDrawerModalProps) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<CloseCashDrawerFormValues | null>(null)

  const { mutate: closeDrawer, isPending } = useCloseCashDrawerMutation()

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm<CloseCashDrawerFormValues>({
    resolver: zodResolver(closeCashDrawerSchema),
    defaultValues,
  })

  useEffect(() => {
    if (!open) return
    reset(defaultValues)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: CloseCashDrawerFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues || cashDrawerId === null) return

    closeDrawer(
      {
        id: cashDrawerId,
        closing_balance: pendingValues.closing_balance,
        notes: pendingValues.notes || undefined,
      },
      {
        onSuccess: () => {
          toast.success('Kas berhasil ditutup')
          handleClose()
        },
      }
    )
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title="Tutup Kas Hari Ini"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Tutup Kas"
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="close-balance">Saldo Penutupan (Rp)</Label>
            <Controller
              name="closing_balance"
              control={control}
              render={({ field }) => (
                <RupiahInput
                  id="close-balance"
                  placeholder="0"
                  value={field.value}
                  onChange={field.onChange}
                  className={errors.closing_balance ? 'border-red-500' : ''}
                />
              )}
            />
            {errors.closing_balance && (
              <p className="text-xs text-red-500">{errors.closing_balance.message}</p>
            )}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="close-notes">Catatan (opsional)</Label>
            <Input
              id="close-notes"
              {...register('notes')}
              placeholder="Masukkan catatan penutupan kas..."
            />
          </div>
        </div>
      </FormModal>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) {
            setIsConfirming(false)
            setPendingValues(null)
          }
        }}
        title="Tutup Kas"
        description="Yakin ingin menutup kas hari ini? Tindakan ini tidak dapat dibatalkan."
        confirmLabel="Ya, Tutup Kas"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
