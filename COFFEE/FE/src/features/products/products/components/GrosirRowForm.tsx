import { useForm, Controller, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Sparkles } from 'lucide-react'

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
import { formatRupiah } from '@/shared/utils'

import { grosirSchema } from '../products.schema'
import type { GrosirFormValues } from '../products.schema'

export interface PackageRefOption {
  /** ID asli (mode edit) atau temp_id (mode tambah produk, 0 = paket dasar) */
  id: number
  label: string
  /** Harga per satuan acuan ini — dipakai untuk menghitung Ref. Harga Beli/Jual satuan baru. */
  purchase_price?: number
  selling_price?: number
}

interface GrosirRowFormProps {
  /** Pilihan untuk dropdown "Merujuk ke Satuan" — sudah termasuk paket dasar. */
  refOptions: PackageRefOption[]
  availableUnits: { id: number; name: string; abbreviation: string }[]
  initialValues?: GrosirFormValues
  isSaving?: boolean
  onSave: (values: GrosirFormValues) => void
  onCancel: () => void
}

const defaultGrosirValues: GrosirFormValues = {
  unit_id: 0,
  package_name: '',
  ref_package_id: 0,
  qty: 1,
  ref_qty: 1,
  purchase_price: 0,
  selling_price: 0,
}

export function GrosirRowForm({
  refOptions,
  availableUnits,
  initialValues,
  isSaving,
  onSave,
  onCancel,
}: GrosirRowFormProps) {
  // Default "Merujuk ke Satuan" ke opsi pertama yang tersedia — di mode tambah produk itu
  // selalu id 0 (placeholder anchor), tapi di mode edit produk anchor punya ID asli
  // (bukan 0), jadi tidak bisa di-hardcode 0 di sini.
  const {
    register,
    handleSubmit,
    setValue: setGrosirValue,
    control: grosirControl,
    formState: { errors },
  } = useForm<GrosirFormValues>({
    resolver: zodResolver(grosirSchema),
    defaultValues: initialValues ?? { ...defaultGrosirValues, ref_package_id: refOptions[0]?.id ?? 0 },
  })

  const unitIdVal = useWatch({ control: grosirControl, name: 'unit_id' })
  const refPackageIdVal = useWatch({ control: grosirControl, name: 'ref_package_id' })
  const qtyVal = useWatch({ control: grosirControl, name: 'qty' }) || 0
  const refQtyVal = useWatch({ control: grosirControl, name: 'ref_qty' }) || 0

  const refSelected = refOptions.find((p) => p.id === refPackageIdVal)
  const selectedUnitName = availableUnits.find((u) => u.id === unitIdVal)?.name ?? '(satuan)'

  // Harga per satuan baru = harga per satuan acuan × (ref_qty / qty), sesuai rasio yang diisi
  // (contoh: 12 Batang = 1 Pack → harga per Batang = harga per Pack / 12)
  const hasRatio = qtyVal > 0 && refQtyVal > 0
  const refPurchasePricePerUnit = hasRatio && refSelected?.purchase_price != null
    ? (refSelected.purchase_price * refQtyVal) / qtyVal
    : null
  const refSellingPricePerUnit = hasRatio && refSelected?.selling_price != null
    ? (refSelected.selling_price * refQtyVal) / qtyVal
    : null

  return (
    <div className="rounded-md border border-blue-100 bg-blue-50/40 p-3 space-y-3">
      {/* Contoh konkret — supaya pemakai pertama kali langsung paham cara isinya */}
      <div className="rounded-md bg-white border border-blue-100 px-3 py-2 text-xs text-gray-600">
        <span className="font-medium text-gray-700">Contoh:</span> kalau 1 Dus isi 24 Pack, isi{' '}
        <span className="font-medium">Satuan Baru</span> = Dus, <span className="font-medium">Merujuk ke Satuan</span> = Pack,
        lalu kolom kiri diisi <span className="font-mono">1</span> dan kolom kanan diisi{' '}
        <span className="font-mono">24</span>.
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div className="space-y-1">
          <Label className="text-xs">Satuan Baru *</Label>
          <Select
            value={unitIdVal > 0 ? String(unitIdVal) : ''}
            onValueChange={(v) => setGrosirValue('unit_id', Number(v))}
          >
            <SelectTrigger className={`h-8 text-sm bg-white ${errors.unit_id ? 'border-red-500' : ''}`}>
              <SelectValue placeholder="Pilih satuan yang mau ditambah" />
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
          <Label className="text-xs">Nama Paket (opsional)</Label>
          <Input
            {...register('package_name')}
            placeholder="Misal: Lama, Baru, Promo"
            className="h-8 text-sm bg-white"
          />
          <p className="text-[11px] text-gray-500">
            Isi kalau nanti ada lebih dari satu paket untuk satuan yang sama (mis. rasio beda per supplier/promo).
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div className="space-y-1">
          <Label className="text-xs">Merujuk ke Satuan *</Label>
          <Select
            value={String(refPackageIdVal ?? 0)}
            onValueChange={(v) => setGrosirValue('ref_package_id', Number(v))}
          >
            <SelectTrigger className={`h-8 text-sm bg-white ${errors.ref_package_id ? 'border-red-500' : ''}`}>
              <SelectValue placeholder="Pilih satuan yang sudah ada" />
            </SelectTrigger>
            <SelectContent>
              {refOptions.map((p) => (
                <SelectItem key={p.id} value={String(p.id)}>
                  {p.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.ref_package_id && <p className="text-xs text-red-500">{errors.ref_package_id.message}</p>}
        </div>
        <div className="space-y-1">
          <Label className="text-xs">Perbandingan Jumlah *</Label>
          <div className="flex items-center gap-1.5">
            <Input
              type="number"
              min={0.001}
              step="any"
              {...register('qty', { valueAsNumber: true })}
              className="h-8 text-sm bg-white"
            />
            <span className="text-xs text-gray-500 whitespace-nowrap">{selectedUnitName} =</span>
            <Input
              type="number"
              min={0.001}
              step="any"
              {...register('ref_qty', { valueAsNumber: true })}
              className="h-8 text-sm bg-white"
            />
          </div>
          {errors.qty && <p className="text-xs text-red-500">{errors.qty.message}</p>}
          {errors.ref_qty && <p className="text-xs text-red-500">{errors.ref_qty.message}</p>}
        </div>
      </div>

      {unitIdVal > 0 && refSelected && qtyVal > 0 && refQtyVal > 0 && (
        <div className="rounded-md bg-green-50 border border-green-200 px-3 py-2 text-sm text-green-700 font-medium">
          {qtyVal} {selectedUnitName} = {refQtyVal} {refSelected.label}
        </div>
      )}

      {(refPurchasePricePerUnit != null || refSellingPricePerUnit != null) && (
        <div className="rounded-md bg-indigo-50 border border-indigo-200 px-3 py-2.5 space-y-2">
          <div className="flex items-center gap-1.5 text-xs font-medium text-indigo-700">
            <Sparkles size={13} />
            Estimasi harga per {selectedUnitName}, dihitung dari {refSelected?.label}
          </div>
          <div className="flex flex-wrap gap-2">
            {refPurchasePricePerUnit != null && (
              <button
                type="button"
                onClick={() => setGrosirValue('purchase_price', Math.round(refPurchasePricePerUnit))}
                className="group inline-flex items-center gap-2 rounded-full border border-indigo-200 bg-white pl-2.5 pr-1 py-1 text-xs text-gray-700 transition-colors hover:border-indigo-400 hover:bg-indigo-100"
                title="Klik untuk memakai nilai ini sebagai Harga Beli"
              >
                <span className="text-gray-500">Beli</span>
                <span className="font-semibold text-indigo-700">
                  {formatRupiah(Math.round(refPurchasePricePerUnit))}
                </span>
                <span className="rounded-full bg-indigo-600 px-1.5 py-0.5 text-[10px] font-medium text-white group-hover:bg-indigo-700">
                  Pakai
                </span>
              </button>
            )}
            {refSellingPricePerUnit != null && (
              <button
                type="button"
                onClick={() => setGrosirValue('selling_price', Math.round(refSellingPricePerUnit))}
                className="group inline-flex items-center gap-2 rounded-full border border-indigo-200 bg-white pl-2.5 pr-1 py-1 text-xs text-gray-700 transition-colors hover:border-indigo-400 hover:bg-indigo-100"
                title="Klik untuk memakai nilai ini sebagai Harga Jual"
              >
                <span className="text-gray-500">Jual</span>
                <span className="font-semibold text-indigo-700">
                  {formatRupiah(Math.round(refSellingPricePerUnit))}
                </span>
                <span className="rounded-full bg-indigo-600 px-1.5 py-0.5 text-[10px] font-medium text-white group-hover:bg-indigo-700">
                  Pakai
                </span>
              </button>
            )}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div className="space-y-1">
          <Label className="text-xs">Harga Beli</Label>
          <Controller
            name="purchase_price"
            control={grosirControl}
            render={({ field }) => (
              <RupiahInput value={field.value} onChange={field.onChange} className="h-8 text-sm bg-white" />
            )}
          />
        </div>
        <div className="space-y-1">
          <Label className="text-xs">Harga Jual *</Label>
          <Controller
            name="selling_price"
            control={grosirControl}
            render={({ field }) => (
              <RupiahInput value={field.value} onChange={field.onChange} className="h-8 text-sm bg-white" />
            )}
          />
          {errors.selling_price && <p className="text-xs text-red-500">{errors.selling_price.message}</p>}
        </div>
      </div>

      <div className="flex gap-2 justify-end">
        <Button type="button" variant="outline" size="sm" onClick={onCancel} disabled={isSaving}>
          Batal
        </Button>
        <Button type="button" size="sm" onClick={handleSubmit(onSave)} disabled={isSaving}>
          {isSaving ? 'Menyimpan...' : 'Simpan Paket'}
        </Button>
      </div>
    </div>
  )
}
