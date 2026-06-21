import { useEffect, useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { useOpenCashDrawerMutation } from '../cash-drawer.api'
import { openCashDrawerSchema, type OpenCashDrawerFormValues } from '../cash-drawer.schema'
import type { ShiftType } from '../cash-drawer.types'

const defaultValues: OpenCashDrawerFormValues = {
  opening_balance: 0,
  shift: undefined,
  notes: '',
}

const SHIFT_OPTIONS: { value: ShiftType; label: string }[] = [
  { value: 'pagi', label: 'Pagi' },
  { value: 'siang', label: 'Siang' },
  { value: 'malam', label: 'Malam' },
]

interface OpenCashDrawerModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OpenCashDrawerModal({ open, onOpenChange }: OpenCashDrawerModalProps) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<OpenCashDrawerFormValues | null>(null)

  const { mutate: openDrawer, isPending } = useOpenCashDrawerMutation()

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm<OpenCashDrawerFormValues>({
    resolver: zodResolver(openCashDrawerSchema),
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

  const onSubmit = (values: OpenCashDrawerFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    openDrawer(
      {
        opening_balance: pendingValues.opening_balance,
        shift: pendingValues.shift,
        notes: pendingValues.notes || undefined,
      },
      {
        onSuccess: () => {
          toast.success('Kas berhasil dibuka')
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
        title="Buka Kas"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Buka Kas"
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="open-balance">Saldo Awal (Rp)</Label>
            <Controller
              name="opening_balance"
              control={control}
              render={({ field }) => (
                <Input
                  id="open-balance"
                  type="number"
                  min={0}
                  placeholder="0"
                  value={field.value === 0 ? '' : field.value}
                  onChange={(e) => field.onChange(Number(e.target.value) || 0)}
                  className={errors.opening_balance ? 'border-red-500' : ''}
                />
              )}
            />
            {errors.opening_balance && (
              <p className="text-xs text-red-500">{errors.opening_balance.message}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label>Shift (opsional)</Label>
            <Controller
              name="shift"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value ?? ''}
                  onValueChange={(v) => field.onChange(v === '' ? undefined : (v as ShiftType))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Pilih shift..." />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Tidak ada shift</SelectItem>
                    {SHIFT_OPTIONS.map((opt) => (
                      <SelectItem key={opt.value} value={opt.value}>
                        {opt.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="open-notes">Catatan (opsional)</Label>
            <Input
              id="open-notes"
              {...register('notes')}
              placeholder="Catatan pembukaan kas..."
            />
          </div>
        </div>
      </FormModal>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Buka Kas"
        description="Yakin ingin membuka kas sekarang?"
        confirmLabel="Ya, Buka Kas"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
