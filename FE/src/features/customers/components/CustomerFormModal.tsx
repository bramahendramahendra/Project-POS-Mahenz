import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'

import { useCreateCustomerMutation, useUpdateCustomerMutation } from '../customers.api'
import type { Customer } from '../customers.types'
import { customerSchema, type CustomerFormValues } from '../customers.schema'

const defaultValues: CustomerFormValues = {
  name: '',
  phone: '',
  address: '',
  credit_limit: 0,
  notes: '',
}

interface CustomerFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  customer?: Customer | null
}

export function CustomerFormModal({ open, onOpenChange, customer }: CustomerFormModalProps) {
  const isEdit = customer != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<CustomerFormValues | null>(null)

  const { mutate: createCustomer, isPending: isCreating } = useCreateCustomerMutation()
  const { mutate: updateCustomer, isPending: isUpdating } = useUpdateCustomerMutation()
  const isPending = isCreating || isUpdating

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CustomerFormValues>({
    resolver: zodResolver(customerSchema),
    defaultValues,
  })

  useEffect(() => {
    if (open) {
      if (customer) {
        reset({
          name: customer.name,
          phone: customer.phone ?? '',
          address: customer.address ?? '',
          credit_limit: customer.credit_limit ?? 0,
          notes: customer.notes ?? '',
        })
      } else {
        reset(defaultValues)
      }
    } else {
      setPendingValues(null)
      setIsConfirming(false)
    }
  }, [open, customer, reset])

  const creditLimit = watch('credit_limit') ?? 0

  const onSubmit = (values: CustomerFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && customer) {
      updateCustomer(
        { id: customer.id, ...pendingValues },
        {
          onSuccess: () => {
            setIsConfirming(false)
            onOpenChange(false)
          },
        }
      )
    } else {
      createCustomer(pendingValues, {
        onSuccess: () => {
          setIsConfirming(false)
          onOpenChange(false)
        },
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
        title={isEdit ? 'Edit Pelanggan' : 'Tambah Pelanggan'}
        size="md"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label htmlFor="cust-name">
              Nama Pelanggan <span className="text-red-500">*</span>
            </Label>
            <Input
              id="cust-name"
              {...register('name')}
              placeholder="Nama pelanggan"
              className={errors.name ? 'border-red-500' : ''}
            />
            {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="cust-phone">Telepon</Label>
              <Input id="cust-phone" {...register('phone')} placeholder="08xx-xxxx-xxxx" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="cust-credit">Limit Kredit</Label>
              <RupiahInput
                id="cust-credit"
                value={creditLimit}
                onChange={(val) => setValue('credit_limit', val)}
                placeholder="0"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="cust-address">Alamat</Label>
            <Textarea
              id="cust-address"
              {...register('address')}
              placeholder="Alamat lengkap"
              className="resize-none"
              rows={2}
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="cust-notes">Catatan</Label>
            <Textarea
              id="cust-notes"
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
        title={isEdit ? 'Update Pelanggan' : 'Tambah Pelanggan'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} pelanggan "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
