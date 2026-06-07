import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'

import type { Category } from '../categories.types'
import { categorySchema, type CategoryFormValues } from '../categories.schema'

const defaultValues: CategoryFormValues = {
  name: '',
  description: '',
}

interface CategoryFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  category?: Category | null
  onSubmit: (values: CategoryFormValues) => void
  isLoading?: boolean
}

export function CategoryFormModal({
  open,
  onOpenChange,
  category,
  onSubmit,
  isLoading,
}: CategoryFormModalProps) {
  const isEdit = category != null

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CategoryFormValues>({
    resolver: zodResolver(categorySchema),
    defaultValues,
  })

  useEffect(() => {
    if (open) {
      if (category) {
        reset({
          name: category.name,
          description: category.description ?? '',
        })
      } else {
        reset(defaultValues)
      }
    }
  }, [open, category, reset])

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Edit Kategori' : 'Tambah Kategori'}
      size="sm"
      isLoading={isLoading}
      onSubmit={handleSubmit(onSubmit)}
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="category-name">
            Nama Kategori <span className="text-red-500">*</span>
          </Label>
          <Input
            id="category-name"
            {...register('name')}
            placeholder="Nama kategori"
            className={errors.name ? 'border-red-500' : ''}
          />
          {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="category-description">Deskripsi</Label>
          <Textarea
            id="category-description"
            {...register('description')}
            placeholder="Deskripsi kategori (opsional)"
            className="resize-none"
            rows={3}
          />
        </div>
      </div>
    </FormModal>
  )
}
