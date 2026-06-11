import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Textarea } from '@/shared/components/ui/textarea'
import { Label } from '@/shared/components/ui/label'

import { useCreateSupplierMutation, useUpdateSupplierMutation } from '../suppliers.api'
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
}

export function SupplierFormModal({ open, onOpenChange, supplier }: SupplierFormModalProps) {
  const isEdit = supplier != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<SupplierFormValues | null>(null)

  const { mutate: createSupplier, isPending: isCreating } = useCreateSupplierMutation()
  const { mutate: updateSupplier, isPending: isUpdating } = useUpdateSupplierMutation()
  const isPending = isCreating || isUpdating

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
    } else {
      setPendingValues(null)
      setIsConfirming(false)
    }
  }, [open, supplier, reset])

  const onSubmit = (values: SupplierFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && supplier) {
      updateSupplier(
        { id: supplier.id, ...pendingValues },
        {
          onSuccess: () => {
            toast.success('Supplier berhasil diperbarui')
            setIsConfirming(false)
            onOpenChange(false)
          },
          onError: (error) => toast.error(error.message),
        }
      )
    } else {
      createSupplier(pendingValues, {
        onSuccess: () => {
          toast.success('Supplier berhasil ditambahkan')
          setIsConfirming(false)
          onOpenChange(false)
        },
        onError: (error) => toast.error(error.message),
      })
    }
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (isConfirming) return
          onOpenChange(val)
        }}
        title={isEdit ? 'Edit Supplier' : 'Tambah Supplier'}
        size="md"
        isLoading={isPending}
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

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) {
            setIsConfirming(false)
            setPendingValues(null)
          }
        }}
        title={isEdit ? 'Update Supplier' : 'Tambah Supplier'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} supplier "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
