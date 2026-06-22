import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

import { useCreateUnitMutation, useUpdateUnitMutation } from '../units.api'
import { unitSchema, type UnitFormValues } from '../units.schema'
import type { Unit } from '../units.types'

interface UnitFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  unit?: Unit | null
}

const defaultValues: UnitFormValues = {
  name: '',
  abbreviation: '',
}

export function UnitFormModal({ open, onOpenChange, unit }: UnitFormModalProps) {
  const isEdit = unit != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<UnitFormValues | null>(null)

  const { mutate: createUnit, isPending: isCreating } = useCreateUnitMutation()
  const { mutate: updateUnit, isPending: isUpdating } = useUpdateUnitMutation()
  const isPending = isCreating || isUpdating

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<UnitFormValues>({
    resolver: zodResolver(unitSchema),
    defaultValues,
  })

  useEffect(() => {
    if (!open) return
    if (unit) {
      reset({ name: unit.name, abbreviation: unit.abbreviation })
    } else {
      reset(defaultValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, unit])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: UnitFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && unit) {
      updateUnit(
        { id: unit.id, ...pendingValues },
        {
          onSuccess: () => {
            toast.success('Satuan berhasil diperbarui')
            handleClose()
          },
          onError: (error) => toast.error(error.message),
        }
      )
    } else {
      createUnit(pendingValues, {
        onSuccess: () => {
          toast.success('Satuan berhasil ditambahkan')
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
        title={isEdit ? 'Edit Satuan' : 'Tambah Satuan'}
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label htmlFor="unit-name">
              Nama Satuan <span className="text-red-500">*</span>
            </Label>
            <Input
              id="unit-name"
              {...register('name')}
              placeholder="Nama satuan (contoh: Pieces, Lusin, Kardus)"
              className={errors.name ? 'border-red-500' : ''}
            />
            {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="unit-abbreviation">
              Singkatan <span className="text-red-500">*</span>
            </Label>
            <Input
              id="unit-abbreviation"
              {...register('abbreviation')}
              placeholder="Singkatan (contoh: Pcs, Lsn, Kds)"
              className={errors.abbreviation ? 'border-red-500' : ''}
            />
            {errors.abbreviation && (
              <p className="text-xs text-red-500">{errors.abbreviation.message}</p>
            )}
          </div>
        </div>
      </FormModal>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title={isEdit ? 'Update Satuan' : 'Tambah Satuan'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} satuan "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
