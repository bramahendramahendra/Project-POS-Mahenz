import { useEffect, useMemo, useState } from 'react'
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
import { useShiftOptionsQuery } from '@/features/operational/shifts'

import { useOpenCashDrawerMutation } from '../cash-drawer.api'
import { openCashDrawerSchema, type OpenCashDrawerFormValues } from '../cash-drawer.schema'

function detectShiftId(shifts: { id: number; start_time: string; end_time: string }[]): number | undefined {
  if (!shifts.length) return undefined

  const now = new Date()
  const currentMinutes = now.getHours() * 60 + now.getMinutes()

  const toMinutes = (time: string) => {
    const [h, m] = time.split(':').map(Number)
    return h * 60 + m
  }

  // Cari shift yang jam sekarang ada di dalam range-nya
  const matched = shifts.find((s) => {
    const start = toMinutes(s.start_time)
    const end = toMinutes(s.end_time)
    if (start <= end) return currentMinutes >= start && currentMinutes < end
    // Overnight shift (misal 22:00 – 06:00)
    return currentMinutes >= start || currentMinutes < end
  })

  if (matched) return matched.id

  // Tidak ada yang cocok — pilih shift yang paling dekat (start_time terdekat setelah sekarang)
  const sorted = [...shifts].sort((a, b) => {
    const diffA = (toMinutes(a.start_time) - currentMinutes + 1440) % 1440
    const diffB = (toMinutes(b.start_time) - currentMinutes + 1440) % 1440
    return diffA - diffB
  })
  return sorted[0]?.id
}

const defaultValues: OpenCashDrawerFormValues = {
  shift_id: 0,
  opening_balance: 0,
  notes: '',
}

interface OpenCashDrawerModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OpenCashDrawerModal({ open, onOpenChange }: OpenCashDrawerModalProps) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<OpenCashDrawerFormValues | null>(null)

  const { mutate: openDrawer, isPending } = useOpenCashDrawerMutation()
  const { data: shiftOptionsRaw } = useShiftOptionsQuery()
  const shiftOptions = useMemo(() => shiftOptionsRaw ?? [], [shiftOptionsRaw])

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
    const autoId = detectShiftId(shiftOptions)
    reset({ ...defaultValues, shift_id: autoId ?? 0 })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, shiftOptions])

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
        shift_id: pendingValues.shift_id,
        opening_balance: pendingValues.opening_balance,
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

  const selectedShift = shiftOptions.find((s) => s.id === (pendingValues?.shift_id ?? 0))

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
            <Label>
              Shift <span className="text-red-500">*</span>
            </Label>
            <Controller
              name="shift_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value ? String(field.value) : ''}
                  onValueChange={(v) => field.onChange(Number(v))}
                >
                  <SelectTrigger className={errors.shift_id ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Pilih shift..." />
                  </SelectTrigger>
                  <SelectContent>
                    {shiftOptions.map((s) => (
                      <SelectItem key={s.id} value={String(s.id)}>
                        {s.name} ({s.start_time} – {s.end_time})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {errors.shift_id && (
              <p className="text-xs text-red-500">{errors.shift_id.message}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="open-balance">
              Saldo Awal (Rp) <span className="text-red-500">*</span>
            </Label>
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
          if (!val) {
            setIsConfirming(false)
            setPendingValues(null)
          }
        }}
        title="Buka Kas"
        description={`Buka kas untuk shift ${selectedShift ? `${selectedShift.name} (${selectedShift.start_time} – ${selectedShift.end_time})` : ''}?`}
        confirmLabel="Ya, Buka Kas"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
