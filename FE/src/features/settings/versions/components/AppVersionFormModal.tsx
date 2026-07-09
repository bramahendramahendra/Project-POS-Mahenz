import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Switch } from '@/shared/components/ui/switch'
import { Textarea } from '@/shared/components/ui/textarea'

import { useCreateAppVersionMutation } from '../versions.api'
import { createAppVersionSchema, type CreateAppVersionFormValues } from '../versions.schema'

const defaultValues: CreateAppVersionFormValues = {
  version:       '',
  download_url:  '',
  release_notes: '',
  is_mandatory:  false,
}

interface AppVersionFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function AppVersionFormModal({ open, onOpenChange }: AppVersionFormModalProps) {
  const { mutate: create, isPending } = useCreateAppVersionMutation()

  const { register, handleSubmit, reset, watch, setValue, formState: { errors } } =
    useForm<CreateAppVersionFormValues>({ resolver: zodResolver(createAppVersionSchema), defaultValues })

  useEffect(() => { if (!open) reset(defaultValues) }, [open, reset])

  const onSubmit = (values: CreateAppVersionFormValues) => {
    create(values, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Tambah Versi Aplikasi"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Tambah"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>Versi <span className="text-red-500">*</span></Label>
          <Input
            {...register('version')}
            placeholder="contoh: 1.2.3"
            className={errors.version ? 'border-red-500' : ''}
          />
          {errors.version && <p className="text-xs text-red-500">{errors.version.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>URL Download <span className="text-red-500">*</span></Label>
          <Input
            {...register('download_url')}
            placeholder="https://..."
            className={errors.download_url ? 'border-red-500' : ''}
          />
          {errors.download_url && <p className="text-xs text-red-500">{errors.download_url.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>Catatan Rilis</Label>
          <Textarea
            {...register('release_notes')}
            placeholder="Deskripsi perubahan di versi ini..."
            rows={3}
            className="resize-none"
          />
        </div>

        <div className="flex items-center justify-between py-1">
          <div>
            <p className="text-sm font-medium text-gray-800">Update Wajib</p>
            <p className="text-xs text-gray-500">Paksa pengguna memperbarui ke versi ini</p>
          </div>
          <Switch
            checked={watch('is_mandatory')}
            onCheckedChange={(v) => setValue('is_mandatory', v)}
          />
        </div>
      </div>
    </FormModal>
  )
}
