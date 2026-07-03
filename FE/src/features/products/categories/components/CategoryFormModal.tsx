import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'

import { useCreateCategoryMutation, useUpdateCategoryMutation } from '../categories.api'
import { categorySchema, type CategoryFormValues } from '../categories.schema'
import type { Category } from '../categories.types'

interface CategoryFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  category?: Category | null
}

const defaultValues: CategoryFormValues = {
  name: '',
  description: '',
}

export function CategoryFormModal({ open, onOpenChange, category }: CategoryFormModalProps) {
  const isEdit = category != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<CategoryFormValues | null>(null)

  const { mutate: createCategory, isPending: isCreating } = useCreateCategoryMutation()
  const { mutate: updateCategory, isPending: isUpdating } = useUpdateCategoryMutation()
  const isPending = isCreating || isUpdating

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
    if (!open) return
    if (category) {
      reset({ name: category.name, description: category.description ?? '' })
    } else {
      reset(defaultValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, category])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: CategoryFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && category) {
      updateCategory(
        { id: category.id, ...pendingValues },
        {
          onSuccess: () => {
            toast.success('Kategori berhasil diperbarui')
            handleClose()
          },
          onError: (error) => {
            toast.error(error.message)
            setIsConfirming(false)
          },
        }
      )
    } else {
      createCategory(pendingValues, {
        onSuccess: () => {
          toast.success('Kategori berhasil ditambahkan')
          handleClose()
        },
        onError: (error) => {
          toast.error(error.message)
          setIsConfirming(false)
        },
      })
    }
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title={isEdit ? 'Edit Kategori' : 'Tambah Kategori'}
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
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

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title={isEdit ? 'Update Kategori' : 'Tambah Kategori'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} kategori "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
