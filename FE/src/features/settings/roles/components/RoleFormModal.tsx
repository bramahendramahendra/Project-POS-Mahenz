import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'

import { useCreateRoleMutation, useRoleDetailQuery, useUpdateRoleMutation } from '../roles.api'
import { createRoleSchema, editRoleSchema, type CreateRoleFormValues, type EditRoleFormValues } from '../roles.schema'

const createDefaults: CreateRoleFormValues = { name: '', display_name: '', description: '' }
const editDefaults: EditRoleFormValues = { display_name: '', description: '' }

interface RoleFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  roleId?: number
}

function CreateRoleForm({
  open,
  onOpenChange,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const { mutate: create, isPending } = useCreateRoleMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateRoleFormValues>({
    resolver: zodResolver(createRoleSchema),
    defaultValues: createDefaults,
  })

  useEffect(() => {
    if (!open) reset(createDefaults)
  }, [open, reset])

  const onSubmit = (values: CreateRoleFormValues) => {
    create(values, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Tambah Role"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Tambah"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>Nama Role (slug) <span className="text-red-500">*</span></Label>
          <Input
            {...register('name')}
            placeholder="contoh: supervisor"
            className={errors.name ? 'border-red-500' : ''}
          />
          {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
          <p className="text-xs text-gray-400">Hanya huruf kecil, angka, dan underscore. Tidak bisa diubah setelah dibuat.</p>
        </div>

        <div className="space-y-1.5">
          <Label>Label Tampil <span className="text-red-500">*</span></Label>
          <Input
            {...register('display_name')}
            placeholder="contoh: Supervisor Gudang"
            className={errors.display_name ? 'border-red-500' : ''}
          />
          {errors.display_name && <p className="text-xs text-red-500">{errors.display_name.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>Deskripsi</Label>
          <Textarea {...register('description')} placeholder="Deskripsi role ini..." rows={3} />
        </div>
      </div>
    </FormModal>
  )
}

function EditRoleForm({
  open,
  onOpenChange,
  roleId,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  roleId: number
}) {
  const { data: detail } = useRoleDetailQuery(open ? roleId : 0)
  const { mutate: update, isPending } = useUpdateRoleMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<EditRoleFormValues>({
    resolver: zodResolver(editRoleSchema),
    defaultValues: editDefaults,
  })

  useEffect(() => {
    if (!open) {
      reset(editDefaults)
    }
  }, [open, reset])

  useEffect(() => {
    if (detail) {
      reset({ display_name: detail.display_name, description: detail.description ?? '' })
    }
  }, [detail, reset])

  const onSubmit = (values: EditRoleFormValues) => {
    update({ id: roleId, ...values }, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Edit Role"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Simpan"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>Label Tampil <span className="text-red-500">*</span></Label>
          <Input
            {...register('display_name')}
            placeholder="contoh: Supervisor Gudang"
            className={errors.display_name ? 'border-red-500' : ''}
          />
          {errors.display_name && <p className="text-xs text-red-500">{errors.display_name.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>Deskripsi</Label>
          <Textarea {...register('description')} placeholder="Deskripsi role ini..." rows={3} />
        </div>
      </div>
    </FormModal>
  )
}

export function RoleFormModal({ open, onOpenChange, roleId }: RoleFormModalProps) {
  if (roleId !== undefined) {
    return <EditRoleForm open={open} onOpenChange={onOpenChange} roleId={roleId} />
  }
  return <CreateRoleForm open={open} onOpenChange={onOpenChange} />
}
