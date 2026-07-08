import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'

import { useCreateMenuMutation, useMenuDetailQuery, useUpdateMenuMutation } from '@/features/menu/menu.api'
import type { MenuOption } from '@/features/menu/menu.types'
import {
  createMenuSchema,
  editMenuSchema,
  type CreateMenuFormValues,
  type EditMenuFormValues,
} from '../menus.schema'

const createDefaults: CreateMenuFormValues = {
  parent_id: null,
  key_name: '',
  label: '',
  icon: '',
  path: '',
  order_index: 0,
}

const editDefaults: EditMenuFormValues = {
  parent_id: null,
  label: '',
  icon: '',
  path: '',
  order_index: 0,
}

interface MenuFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  menuId?: number
  parentOptions: MenuOption[]
}

function ParentSelect({
  value,
  onChange,
  options,
}: {
  value: number | null | undefined
  onChange: (v: number | null) => void
  options: MenuOption[]
}) {
  return (
    <Select
      value={value?.toString() ?? 'none'}
      onValueChange={(v) => onChange(v === 'none' ? null : Number(v))}
    >
      <SelectTrigger>
        <SelectValue placeholder="Tidak ada (root)" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="none">Tidak ada (root)</SelectItem>
        {options.map((m) => (
          <SelectItem key={m.id} value={m.id.toString()}>{m.label}</SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}

function CreateMenuForm({
  open,
  onOpenChange,
  parentOptions,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  parentOptions: MenuOption[]
}) {
  const { mutate: create, isPending } = useCreateMenuMutation()

  const { register, handleSubmit, reset, watch, setValue, formState: { errors } } =
    useForm<CreateMenuFormValues>({ resolver: zodResolver(createMenuSchema), defaultValues: createDefaults })

  useEffect(() => { if (!open) reset(createDefaults) }, [open, reset])

  const onSubmit = (values: CreateMenuFormValues) => {
    create(values, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Tambah Menu"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Tambah"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>Key Name <span className="text-red-500">*</span></Label>
          <Input
            {...register('key_name')}
            placeholder="contoh: inventory.warehouse"
            className={errors.key_name ? 'border-red-500' : ''}
          />
          {errors.key_name && <p className="text-xs text-red-500">{errors.key_name.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>Parent Menu</Label>
          <ParentSelect
            value={watch('parent_id')}
            onChange={(v) => setValue('parent_id', v)}
            options={parentOptions}
          />
        </div>

        <div className="space-y-1.5">
          <Label>Label <span className="text-red-500">*</span></Label>
          <Input
            {...register('label')}
            placeholder="contoh: Gudang"
            className={errors.label ? 'border-red-500' : ''}
          />
          {errors.label && <p className="text-xs text-red-500">{errors.label.message}</p>}
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <Label>Icon (Lucide)</Label>
            <Input {...register('icon')} placeholder="contoh: Warehouse" />
          </div>
          <div className="space-y-1.5">
            <Label>Path URL</Label>
            <Input {...register('path')} placeholder="contoh: /warehouse" />
          </div>
        </div>

        <div className="space-y-1.5">
          <Label>Urutan</Label>
          <Input type="number" {...register('order_index', { valueAsNumber: true })} />
        </div>
      </div>
    </FormModal>
  )
}

function EditMenuForm({
  open,
  onOpenChange,
  menuId,
  parentOptions,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  menuId: number
  parentOptions: MenuOption[]
}) {
  const { data: detail } = useMenuDetailQuery(open ? menuId : 0)
  const { mutate: update, isPending } = useUpdateMenuMutation()

  const { register, handleSubmit, reset, watch, setValue, formState: { errors } } =
    useForm<EditMenuFormValues>({ resolver: zodResolver(editMenuSchema), defaultValues: editDefaults })

  useEffect(() => { if (!open) reset(editDefaults) }, [open, reset])

  useEffect(() => {
    if (detail) {
      reset({
        parent_id:   detail.parent_id,
        label:       detail.label,
        icon:        detail.icon ?? '',
        path:        detail.path ?? '',
        order_index: detail.order_index,
      })
    }
  }, [detail, reset])

  const onSubmit = (values: EditMenuFormValues) => {
    update({ id: menuId, ...values }, { onSuccess: () => onOpenChange(false) })
  }

  const filteredParentOptions = parentOptions.filter((m) => m.id !== menuId)

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Edit Menu"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Simpan"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>Parent Menu</Label>
          <ParentSelect
            value={watch('parent_id')}
            onChange={(v) => setValue('parent_id', v)}
            options={filteredParentOptions}
          />
        </div>

        <div className="space-y-1.5">
          <Label>Label <span className="text-red-500">*</span></Label>
          <Input
            {...register('label')}
            placeholder="contoh: Gudang"
            className={errors.label ? 'border-red-500' : ''}
          />
          {errors.label && <p className="text-xs text-red-500">{errors.label.message}</p>}
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <Label>Icon (Lucide)</Label>
            <Input {...register('icon')} placeholder="contoh: Warehouse" />
          </div>
          <div className="space-y-1.5">
            <Label>Path URL</Label>
            <Input {...register('path')} placeholder="contoh: /warehouse" />
          </div>
        </div>

        <div className="space-y-1.5">
          <Label>Urutan</Label>
          <Input type="number" {...register('order_index', { valueAsNumber: true })} />
        </div>
      </div>
    </FormModal>
  )
}

export function MenuFormModal({ open, onOpenChange, menuId, parentOptions }: MenuFormModalProps) {
  if (menuId !== undefined) {
    return <EditMenuForm open={open} onOpenChange={onOpenChange} menuId={menuId} parentOptions={parentOptions} />
  }
  return <CreateMenuForm open={open} onOpenChange={onOpenChange} parentOptions={parentOptions} />
}
