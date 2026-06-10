import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'

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
  onSubmit: (values: CustomerFormValues) => void
  isLoading?: boolean
}

export function CustomerFormModal({
  open,
  onOpenChange,
  customer,
  onSubmit,
  isLoading,
}: CustomerFormModalProps) {
  const isEdit = customer != null

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
    }
  }, [open, customer, reset])

  const creditLimit = watch('credit_limit') ?? 0

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Edit Pelanggan' : 'Tambah Pelanggan'}
      size="md"
      isLoading={isLoading}
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
  )
}
