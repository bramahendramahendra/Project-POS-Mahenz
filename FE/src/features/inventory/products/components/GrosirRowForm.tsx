import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { grosirSchema } from '../products.schema'
import type { GrosirFormValues } from '../products.schema'

interface GrosirRowFormProps {
  baseUnitName: string
  basePurchase: number
  baseSelling: number
  availableUnits: { id: number; name: string; abbreviation: string }[]
  initialValues?: GrosirFormValues
  onSave: (v: GrosirFormValues) => void
  onCancel: () => void
}

export function GrosirRowForm({
  baseUnitName,
  basePurchase,
  baseSelling,
  availableUnits,
  initialValues,
  onSave,
  onCancel,
}: GrosirRowFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    setValue: setGrosirValue,
    control: grosirControl,
    formState: { errors },
  } = useForm<GrosirFormValues>({
    resolver: zodResolver(grosirSchema),
    defaultValues: initialValues ?? { unit_id: 0, package_name: '', conversion_qty: 1, purchase_price: 0, selling_price: 0 },
  })

  const convQty = watch('conversion_qty') || 0
  const unitIdVal = watch('unit_id')
  const refPurchase = basePurchase * convQty
  const refSelling = baseSelling * convQty

  return (
    <div className="rounded-md border bg-gray-50 p-3 space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <Label className="text-xs">Satuan *</Label>
          <Select
            value={unitIdVal > 0 ? String(unitIdVal) : ''}
            onValueChange={(v) => setGrosirValue('unit_id', Number(v))}
          >
            <SelectTrigger className={`h-8 text-sm ${errors.unit_id ? 'border-red-500' : ''}`}>
              <SelectValue placeholder="Pilih Satuan" />
            </SelectTrigger>
            <SelectContent>
              {availableUnits.filter((u) => u.id > 0).map((u) => (
                <SelectItem key={u.id} value={String(u.id)}>
                  {u.name} ({u.abbreviation})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.unit_id && <p className="text-xs text-red-500">{errors.unit_id.message}</p>}
        </div>
        <div className="space-y-1">
          <Label className="text-xs">Nama Paket</Label>
          <Input
            {...register('package_name')}
            placeholder="Opsional, misal: 1 Dus, 3 Botol"
            className="h-8 text-sm"
          />
        </div>
      </div>

      <div className="space-y-1">
        <Label className="text-xs">Konversi ke {baseUnitName || 'Satuan Dasar'} *</Label>
        <Input
          type="number"
          min={1}
          {...register('conversion_qty', { valueAsNumber: true })}
          className="h-8 text-sm"
        />
        {convQty > 0 && baseUnitName && (
          <p className="text-xs text-blue-600">1 paket ini = {convQty} {baseUnitName}</p>
        )}
        {errors.conversion_qty && <p className="text-xs text-red-500">{errors.conversion_qty.message}</p>}
      </div>

      {convQty > 0 && (basePurchase > 0 || baseSelling > 0) && (
        <div className="rounded bg-blue-50 border border-blue-100 px-3 py-1.5 text-xs text-blue-700 flex gap-4">
          <span>Ref. Harga Beli: <strong>Rp {refPurchase.toLocaleString('id-ID')}</strong></span>
          <span>Ref. Harga Jual: <strong>Rp {refSelling.toLocaleString('id-ID')}</strong></span>
        </div>
      )}

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <Label className="text-xs">Harga Beli Aktual</Label>
          <Controller
            name="purchase_price"
            control={grosirControl}
            render={({ field }) => (
              <RupiahInput value={field.value} onChange={field.onChange} className="h-8 text-sm" />
            )}
          />
        </div>
        <div className="space-y-1">
          <Label className="text-xs">Harga Jual Aktual *</Label>
          <Controller
            name="selling_price"
            control={grosirControl}
            render={({ field }) => (
              <RupiahInput value={field.value} onChange={field.onChange} className="h-8 text-sm" />
            )}
          />
          {errors.selling_price && <p className="text-xs text-red-500">{errors.selling_price.message}</p>}
        </div>
      </div>

      <div className="flex gap-2 justify-end">
        <Button type="button" variant="outline" size="sm" onClick={onCancel}>Batal</Button>
        <Button type="button" size="sm" onClick={handleSubmit(onSave)}>Simpan Paket</Button>
      </div>
    </div>
  )
}
