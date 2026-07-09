import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'

import { useStoreProfileQuery, useUpdateStoreProfileMutation } from '../store.api'
import { storeProfileSchema, type StoreProfileFormValues } from '../store.schema'

interface StoreProfileFormProps {
  onCancel: () => void
  onSuccess: () => void
}

export function StoreProfileForm({ onCancel, onSuccess }: StoreProfileFormProps) {
  const { data: profile, isLoading } = useStoreProfileQuery()
  const { mutate: updateProfile, isPending } = useUpdateStoreProfileMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<StoreProfileFormValues>({
    resolver: zodResolver(storeProfileSchema),
    defaultValues: { name: '', address: '', phone: '', email: '', tax_default: 0 },
  })

  useEffect(() => {
    if (profile) {
      reset({
        name: profile.name,
        address: profile.address ?? '',
        phone: profile.phone ?? '',
        email: profile.email ?? '',
        tax_default: profile.tax_default ?? 0,
      })
    }
  }, [profile, reset])

  const onSubmit = (values: StoreProfileFormValues) => {
    updateProfile(
      {
        name: values.name,
        address: values.address || undefined,
        phone: values.phone || undefined,
        email: values.email || undefined,
        tax_default: values.tax_default,
      },
      { onSuccess },
    )
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-10 animate-pulse rounded-lg bg-gray-100" />
        ))}
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 max-w-lg">
      <div className="space-y-1.5">
        <Label htmlFor="store-name">
          Nama Toko <span className="text-red-500">*</span>
        </Label>
        <Input
          id="store-name"
          {...register('name')}
          className={errors.name ? 'border-red-500' : ''}
        />
        {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="store-address">Alamat</Label>
        <Input id="store-address" {...register('address')} placeholder="Alamat toko (opsional)" />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="store-phone">Nomor Telepon</Label>
        <Input id="store-phone" {...register('phone')} placeholder="08xx..." />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="store-email">Email</Label>
        <Input
          id="store-email"
          type="email"
          {...register('email')}
          className={errors.email ? 'border-red-500' : ''}
          placeholder="toko@email.com"
        />
        {errors.email && <p className="text-xs text-red-500">{errors.email.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="store-tax">Pajak Default Kasir (%)</Label>
        <Input
          id="store-tax"
          type="number"
          min={0}
          max={100}
          {...register('tax_default', { valueAsNumber: true })}
          className={errors.tax_default ? 'border-red-500' : ''}
          placeholder="0"
        />
        {errors.tax_default && <p className="text-xs text-red-500">{errors.tax_default.message}</p>}
        <p className="text-xs text-gray-400">Dipakai sebagai nilai awal pajak di halaman kasir</p>
      </div>

      <div className="flex gap-2">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isPending}>
          Batal
        </Button>
        <Button type="submit" disabled={isPending}>
          {isPending ? 'Menyimpan...' : 'Simpan Perubahan'}
        </Button>
      </div>
    </form>
  )
}
