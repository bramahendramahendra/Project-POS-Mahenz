import { useEffect, useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'
import { useUnitOptionsQuery } from '@/features/products/units/units.api'

import { useCreateIngredientMutation, useUpdateIngredientMutation } from '../ingredients.api'
import { ingredientSchema, type IngredientFormValues } from '../ingredients.schema'
import type { Ingredient } from '../ingredients.types'
import { DecimalInput } from './DecimalInput'

interface IngredientFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  ingredient: Ingredient | null
}

const defaultValues: IngredientFormValues = {
  name: '',
  unit_id: undefined as unknown as number,
  min_stock: null,
}

// Tambah/ubah bahan. Stok bahan baru selalu 0 (keputusan 83); stok bertambah lewat pembelian.
export function IngredientFormModal({ open, onOpenChange, ingredient }: IngredientFormModalProps) {
  const isEdit = ingredient != null
  // Satuan terkunci bila bahan sudah punya riwayat stok (BE juga menolak).
  const unitLocked = isEdit && !!ingredient.has_history

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<IngredientFormValues | null>(null)

  const { data: units = [] } = useUnitOptionsQuery()
  const { mutate: createIngredient, isPending: isCreating } = useCreateIngredientMutation()
  const { mutate: updateIngredient, isPending: isUpdating } = useUpdateIngredientMutation()
  const isPending = isCreating || isUpdating

  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<IngredientFormValues>({ resolver: zodResolver(ingredientSchema), defaultValues })

  useEffect(() => {
    if (!open) return
    reset(
      ingredient
        ? { name: ingredient.name, unit_id: ingredient.unit_id, min_stock: ingredient.min_stock }
        : defaultValues
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, ingredient])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: IngredientFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return
    const callbacks = { onSuccess: () => handleClose(), onError: () => setIsConfirming(false) }
    if (isEdit && ingredient) updateIngredient({ id: ingredient.id, ...pendingValues }, callbacks)
    else createIngredient(pendingValues, callbacks)
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title={isEdit ? 'Edit Bahan' : 'Tambah Bahan'}
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label htmlFor="ingredient-name">
              Nama Bahan <span className="text-red-500">*</span>
            </Label>
            <Input
              id="ingredient-name"
              {...register('name')}
              placeholder="Contoh: Kopi Arabika, Susu UHT"
              className={errors.name ? 'border-red-500' : ''}
            />
            {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
          </div>

          <div className="space-y-1.5">
            <Label>
              Satuan Stok <span className="text-red-500">*</span>
            </Label>
            <Controller
              control={control}
              name="unit_id"
              render={({ field }) =>
                unitLocked ? (
                  <Input value={`${ingredient?.unit_name} (${ingredient?.unit_abbr})`} disabled readOnly />
                ) : (
                  <Select value={field.value ? String(field.value) : ''} onValueChange={(v) => field.onChange(Number(v))}>
                    <SelectTrigger className="w-full" aria-label="Satuan stok">
                      <SelectValue placeholder="Pilih satuan..." />
                    </SelectTrigger>
                    <SelectContent>
                      {units.map((u) => (
                        <SelectItem key={u.id} value={String(u.id)}>
                          {u.name} ({u.abbreviation})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )
              }
            />
            {!unitLocked && units.length === 0 && (
              <p className="text-xs text-amber-600">Belum ada satuan aktif. Tambahkan dulu di menu Produk &gt; Unit.</p>
            )}
            {unitLocked && (
              <p className="text-xs text-gray-500">Satuan tidak bisa diubah karena bahan sudah punya riwayat stok.</p>
            )}
            {errors.unit_id && <p className="text-xs text-red-500">{errors.unit_id.message}</p>}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="ingredient-min-stock">Batas Stok Menipis</Label>
            <Controller
              control={control}
              name="min_stock"
              render={({ field }) => (
                <DecimalInput
                  id="ingredient-min-stock"
                  value={field.value}
                  onChange={field.onChange}
                  placeholder="Opsional, contoh: 2"
                  className={errors.min_stock ? 'border-red-500' : ''}
                />
              )}
            />
            <p className="text-xs text-gray-500">Bahan ditandai "Menipis" bila stok sama atau di bawah angka ini.</p>
            {errors.min_stock && <p className="text-xs text-red-500">{errors.min_stock.message}</p>}
          </div>

          {!isEdit && (
            <p className="rounded-md bg-gray-50 px-3 py-2 text-xs text-gray-600">
              Stok bahan baru dimulai dari 0. Stok bertambah lewat menu Pembelian.
            </p>
          )}
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
        title={isEdit ? 'Update Bahan' : 'Tambah Bahan'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} bahan "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
