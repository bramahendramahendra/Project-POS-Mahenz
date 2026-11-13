import { useEffect } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'

import { useRecordUsageMutation, useRecordWasteMutation } from '../ingredients.api'
import { recordStockSchema, recordWasteSchema, type RecordStockFormValues } from '../ingredients.schema'
import type { Ingredient } from '../ingredients.types'
import { formatQty } from '../ingredients.utils'
import { DecimalInput } from './DecimalInput'

export type RecordMode = 'usage' | 'waste'

interface RecordStockModalProps {
  mode: RecordMode
  ingredient: Ingredient | null
  onOpenChange: (open: boolean) => void
}

const defaultValues: RecordStockFormValues = { qty: undefined as unknown as number, notes: '' }

// Catat bahan terpakai atau terbuang, satu bahan per input (keputusan 14.10). Stok boleh minus
// dengan peringatan (keputusan 14.9). Barang terbuang wajib alasan (keputusan 83).
export function RecordStockModal({ mode, ingredient, onOpenChange }: RecordStockModalProps) {
  const isWaste = mode === 'waste'
  const { mutate: recordUsage, isPending: usagePending } = useRecordUsageMutation()
  const { mutate: recordWaste, isPending: wastePending } = useRecordWasteMutation()

  const {
    control,
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<RecordStockFormValues>({
    resolver: zodResolver(isWaste ? recordWasteSchema : recordStockSchema),
    defaultValues,
  })

  useEffect(() => {
    if (ingredient) reset(defaultValues)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ingredient, mode])

  const qty = useWatch({ control, name: 'qty' })
  const willBeNegative = ingredient != null && typeof qty === 'number' && ingredient.stock - qty < 0

  const onSubmit = (values: RecordStockFormValues) => {
    if (!ingredient) return
    const done = { onSuccess: () => onOpenChange(false) }
    if (isWaste) recordWaste({ ingredient_id: ingredient.id, qty: values.qty, reason: values.notes }, done)
    else recordUsage({ ingredient_id: ingredient.id, qty: values.qty, notes: values.notes }, done)
  }

  return (
    <FormModal
      open={ingredient != null}
      onOpenChange={onOpenChange}
      title={isWaste ? 'Catat Barang Terbuang' : 'Catat Bahan Terpakai'}
      size="sm"
      isLoading={usagePending || wastePending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Catat"
    >
      {ingredient && (
        <div className="space-y-3">
          <div className="rounded-md bg-gray-50 px-3 py-2 text-sm">
            <p className="font-medium text-gray-800">{ingredient.name}</p>
            <p className="text-gray-600">
              Stok sisa: <span className="font-semibold">{formatQty(ingredient.stock, ingredient.unit_abbr)}</span>
            </p>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="record-qty">
              Jumlah ({ingredient.unit_abbr}) <span className="text-red-500">*</span>
            </Label>
            <Controller
              control={control}
              name="qty"
              render={({ field }) => (
                <DecimalInput
                  id="record-qty"
                  value={field.value}
                  onChange={(v) => field.onChange(v ?? undefined)}
                  placeholder="Contoh: 1 atau 0,5"
                  className={errors.qty ? 'border-red-500' : ''}
                />
              )}
            />
            {errors.qty && <p className="text-xs text-red-500">{errors.qty.message}</p>}
            {willBeNegative && !errors.qty && (
              <p className="text-xs text-amber-600">
                Stok akan menjadi minus. Tetap bisa dicatat, periksa kembali jumlahnya.
              </p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="record-notes">
              {isWaste ? (
                <>
                  Alasan <span className="text-red-500">*</span>
                </>
              ) : (
                'Catatan'
              )}
            </Label>
            <Textarea
              id="record-notes"
              rows={2}
              {...register('notes')}
              placeholder={isWaste ? 'Contoh: tumpah, basi, rusak' : 'Opsional, contoh: buka bungkus baru'}
              className={errors.notes ? 'border-red-500' : ''}
            />
            {errors.notes && <p className="text-xs text-red-500">{errors.notes.message}</p>}
          </div>
        </div>
      )}
    </FormModal>
  )
}
