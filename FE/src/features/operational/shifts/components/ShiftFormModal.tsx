import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

import { useCreateShiftMutation, useUpdateShiftMutation } from '../shifts.api'
import { shiftFormSchema, type ShiftFormValues } from '../shifts.schema'
import type { Shift } from '../shifts.types'

interface ShiftFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  shift?: Shift | null
}

const defaultValues: ShiftFormValues = {
  name: '',
  start_time: '',
  end_time: '',
}

export function ShiftFormModal({ open, onOpenChange, shift }: ShiftFormModalProps) {
  const isEdit = shift != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<ShiftFormValues | null>(null)

  const { mutate: createShift, isPending: isCreating } = useCreateShiftMutation()
  const { mutate: updateShift, isPending: isUpdating } = useUpdateShiftMutation()
  const isPending = isCreating || isUpdating

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ShiftFormValues>({
    resolver: zodResolver(shiftFormSchema),
    defaultValues,
  })

  useEffect(() => {
    if (!open) return
    if (shift) {
      reset({ name: shift.name, start_time: shift.start_time, end_time: shift.end_time })
    } else {
      reset(defaultValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, shift])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: ShiftFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && shift) {
      updateShift(
        { id: shift.id, ...pendingValues },
        {
          onSuccess: () => {
            toast.success('Shift berhasil diperbarui')
            handleClose()
          },
          onError: (error) => toast.error(error.message),
        }
      )
    } else {
      createShift(pendingValues, {
        onSuccess: () => {
          toast.success('Shift berhasil ditambahkan')
          handleClose()
        },
        onError: (error) => toast.error(error.message),
      })
    }
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title={isEdit ? 'Edit Shift' : 'Tambah Shift'}
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label htmlFor="shift-name">
              Nama Shift <span className="text-red-500">*</span>
            </Label>
            <Input
              id="shift-name"
              {...register('name')}
              placeholder="contoh: Shift Pagi"
              className={errors.name ? 'border-red-500' : ''}
            />
            {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="shift-start">
                Jam Mulai <span className="text-red-500">*</span>
              </Label>
              <Input
                id="shift-start"
                type="time"
                {...register('start_time')}
                className={errors.start_time ? 'border-red-500' : ''}
              />
              {errors.start_time && (
                <p className="text-xs text-red-500">{errors.start_time.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="shift-end">
                Jam Selesai <span className="text-red-500">*</span>
              </Label>
              <Input
                id="shift-end"
                type="time"
                {...register('end_time')}
                className={errors.end_time ? 'border-red-500' : ''}
              />
              {errors.end_time && (
                <p className="text-xs text-red-500">{errors.end_time.message}</p>
              )}
            </div>
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
        title={isEdit ? 'Update Shift' : 'Tambah Shift'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} shift "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
