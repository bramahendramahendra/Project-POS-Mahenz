import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

import { useOpenShiftMutation } from '../shifts.api'
import { openShiftSchema, type OpenShiftFormValues } from '../shifts.schema'

interface OpenShiftModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OpenShiftModal({ open, onOpenChange }: OpenShiftModalProps) {
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [pendingValues, setPendingValues] = useState<OpenShiftFormValues | null>(null)

  const { mutate: openShift, isPending } = useOpenShiftMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<OpenShiftFormValues>({
    resolver: zodResolver(openShiftSchema),
    defaultValues: { opening_balance: 0, notes: '' },
  })

  useEffect(() => {
    if (!open) reset({ opening_balance: 0, notes: '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const handleClose = () => {
    setConfirmOpen(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: OpenShiftFormValues) => {
    setPendingValues(values)
    setConfirmOpen(true)
  }

  const handleConfirm = () => {
    if (!pendingValues) return
    openShift(
      { opening_balance: pendingValues.opening_balance, notes: pendingValues.notes || undefined },
      {
        onSuccess: () => {
          toast.success('Shift berhasil dibuka')
          handleClose()
        },
      }
    )
  }

  return (
    <>
      <FormModal
        open={open && !confirmOpen}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Buka Shift Baru"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Buka Shift"
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="opening-balance">
              Modal Awal (uang di laci kasir) <span className="text-red-500">*</span>
            </Label>
            <Input
              id="opening-balance"
              type="number"
              min={0}
              {...register('opening_balance', { valueAsNumber: true })}
              className={errors.opening_balance ? 'border-red-500' : ''}
              placeholder="0"
            />
            {errors.opening_balance && (
              <p className="text-xs text-red-500">{errors.opening_balance.message}</p>
            )}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="open-notes">Catatan</Label>
            <Input id="open-notes" {...register('notes')} placeholder="Catatan (opsional)" />
          </div>
        </div>
      </FormModal>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Buka Shift Baru"
        description={`Mulai shift baru dengan modal awal Rp ${(pendingValues?.opening_balance ?? 0).toLocaleString('id-ID')}?`}
        confirmLabel="Ya, Buka Shift"
        onConfirm={handleConfirm}
        isLoading={isPending}
      />
    </>
  )
}
