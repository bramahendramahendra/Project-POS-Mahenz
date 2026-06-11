import { useEffect, useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

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
    if (open) {
      reset(defaultValues)
    } else {
      setPendingValues(null)
      setIsConfirming(false)
    }
  }, [open, reset])

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
          setIsConfirming(false)
          onOpenChange(false)
        },
      }
    )
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (isConfirming) return
          onOpenChange(val)
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
                <Input
                  id="close-balance"
                  type="number"
                  min={0}
                  placeholder="0"
                  value={field.value === 0 ? '' : field.value}
                  onChange={(e) => field.onChange(Number(e.target.value) || 0)}
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
