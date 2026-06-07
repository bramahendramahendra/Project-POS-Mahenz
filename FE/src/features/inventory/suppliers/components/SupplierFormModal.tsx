import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Textarea } from '@/shared/components/ui/textarea'
import { Label } from '@/shared/components/ui/label'

import type { Supplier } from '../suppliers.types'
import { supplierSchema, type SupplierFormValues } from '../suppliers.schema'

const defaultValues: SupplierFormValues = {
  name: '',
  contact_person: '',
  phone: '',
  email: '',
  address: '',
  notes: '',
}

interface SupplierFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  supplier?: Supplier | null
  onSubmit: (values: SupplierFormValues) => void
  isLoading?: boolean
}

export function SupplierFormModal({
  open,
  onOpenChange,
  supplier,
  onSubmit,
  isLoading,
}: SupplierFormModalProps) {
  const isEdit = supplier != null

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<SupplierFormValues>({
    resolver: zodResolver(supplierSchema),
    defaultValues,
  })

  useEffect(() => {
    if (open) {
      if (supplier) {
        reset({
          name: supplier.name,
          contact_person: supplier.contact_person ?? '',
          phone: supplier.phone ?? '',
          email: supplier.email ?? '',
          address: supplier.address ?? '',
          notes: supplier.notes ?? '',
        })
      } else {
        reset(defaultValues)
      }
    }
  }, [open, supplier, reset])

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Edit Supplier' : 'Tambah Supplier'}
      size="md"
      isLoading={isLoading}
      onSubmit={handleSubmit(onSubmit)}
    >
      <div className="space-y-4">
        <div className="space-y-1.5">
          <Label htmlFor="sup-name">
            Nama Supplier <span className="text-red-500">*</span>
          </Label>
          <Input
            id="sup-name"
            {...register('name')}
            placeholder="Nama supplier"
            className={errors.name ? 'border-red-500' : ''}
          />
          {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <Label htmlFor="sup-contact">Nama Kontak</Label>
            <Input id="sup-contact" {...register('contact_person')} placeholder="Nama kontak" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="sup-phone">No. Telepon</Label>
            <Input id="sup-phone" {...register('phone')} placeholder="08xx-xxxx-xxxx" />
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="sup-email">Email</Label>
          <Input
            id="sup-email"
            type="email"
            {...register('email')}
            placeholder="email@supplier.com"
            className={errors.email ? 'border-red-500' : ''}
          />
          {errors.email && <p className="text-xs text-red-500">{errors.email.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="sup-address">Alamat</Label>
          <Textarea
            id="sup-address"
            {...register('address')}
            placeholder="Alamat lengkap supplier"
            className="resize-none"
            rows={2}
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="sup-notes">Catatan</Label>
          <Textarea
            id="sup-notes"
            {...register('notes')}
            placeholder="Catatan tambahan (opsional)"
            className="resize-none"
            rows={2}
          />
        </div>
      </div>
    </FormModal>
  )
}
