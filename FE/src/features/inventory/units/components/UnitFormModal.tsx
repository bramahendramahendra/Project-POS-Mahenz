import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

import type { Unit } from '../units.types'
import { unitSchema, type UnitFormValues } from '../units.schema'

const defaultValues: UnitFormValues = {
  name: '',
  abbreviation: '',
}

interface UnitFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  unit?: Unit | null
  onSubmit: (values: UnitFormValues) => void
  isLoading?: boolean
}

export function UnitFormModal({
  open,
  onOpenChange,
  unit,
  onSubmit,
  isLoading,
}: UnitFormModalProps) {
  const isEdit = unit != null

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
    if (open) {
      if (unit) {
        reset({
          name: unit.name,
          abbreviation: unit.abbreviation,
        })
      } else {
        reset(defaultValues)
      }
    }
  }, [open, unit, reset])

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Edit Satuan' : 'Tambah Satuan'}
      size="sm"
      isLoading={isLoading}
      onSubmit={handleSubmit(onSubmit)}
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
  )
}
