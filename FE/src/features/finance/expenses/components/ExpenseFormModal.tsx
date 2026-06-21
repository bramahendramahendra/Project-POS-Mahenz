import { useEffect, useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { useCreateExpenseMutation, useUpdateExpenseMutation } from '../expenses.api'
import type { Expense } from '../expenses.types'
import {
  expenseSchema,
  type ExpenseFormValues,
  EXPENSE_CATEGORIES,
  EXPENSE_PAYMENT_METHODS,
} from '../expenses.schema'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

const defaultValues: ExpenseFormValues = {
  expense_date: todayString(),
  category: 'operasional',
  description: '',
  amount: 0,
  payment_method: 'cash',
  notes: '',
}

interface ExpenseFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  expense?: Expense | null
}

export function ExpenseFormModal({ open, onOpenChange, expense }: ExpenseFormModalProps) {
  const isEdit = expense != null

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<ExpenseFormValues | null>(null)

  const { mutate: create, isPending: isCreating } = useCreateExpenseMutation()
  const { mutate: update, isPending: isUpdating } = useUpdateExpenseMutation()
  const isPending = isCreating || isUpdating

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm<ExpenseFormValues>({
    resolver: zodResolver(expenseSchema),
    defaultValues,
  })

  useEffect(() => {
    if (!open) return
    if (expense) {
      reset({
        expense_date: expense.expense_date,
        category: expense.category,
        description: expense.description,
        amount: expense.amount,
        payment_method: expense.payment_method,
        notes: expense.notes ?? '',
      })
    } else {
      reset({ ...defaultValues, expense_date: todayString() })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, expense])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: ExpenseFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return
    const payload = { ...pendingValues, notes: pendingValues.notes || undefined }

    if (isEdit && expense) {
      update(
        { id: expense.id, ...payload },
        {
          onSuccess: () => {
            toast.success('Pengeluaran berhasil diperbarui')
            handleClose()
          },
        }
      )
    } else {
      create(payload, {
        onSuccess: () => {
          toast.success('Pengeluaran berhasil ditambahkan')
          handleClose()
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
        title={isEdit ? 'Edit Pengeluaran' : 'Tambah Pengeluaran'}
        size="md"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        <div className="space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="exp-date">
                Tanggal <span className="text-red-500">*</span>
              </Label>
              <Input
                id="exp-date"
                type="date"
                {...register('expense_date')}
                className={errors.expense_date ? 'border-red-500' : ''}
              />
              {errors.expense_date && (
                <p className="text-xs text-red-500">{errors.expense_date.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="exp-category">
                Kategori <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="category"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className={errors.category ? 'border-red-500' : ''}>
                      <SelectValue placeholder="Pilih kategori" />
                    </SelectTrigger>
                    <SelectContent>
                      {EXPENSE_CATEGORIES.map((cat) => (
                        <SelectItem key={cat.value} value={cat.value}>
                          {cat.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.category && (
                <p className="text-xs text-red-500">{errors.category.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="exp-description">
              Keterangan <span className="text-red-500">*</span>
            </Label>
            <Input
              id="exp-description"
              {...register('description')}
              placeholder="Masukkan keterangan pengeluaran"
              className={errors.description ? 'border-red-500' : ''}
            />
            {errors.description && (
              <p className="text-xs text-red-500">{errors.description.message}</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="exp-amount">
                Jumlah <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="amount"
                control={control}
                render={({ field }) => (
                  <RupiahInput
                    id="exp-amount"
                    value={field.value}
                    onChange={field.onChange}
                    className={errors.amount ? 'border-red-500' : ''}
                  />
                )}
              />
              {errors.amount && (
                <p className="text-xs text-red-500">{errors.amount.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="exp-payment">
                Metode Bayar <span className="text-red-500">*</span>
              </Label>
              <Controller
                name="payment_method"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className={errors.payment_method ? 'border-red-500' : ''}>
                      <SelectValue placeholder="Pilih metode" />
                    </SelectTrigger>
                    <SelectContent>
                      {EXPENSE_PAYMENT_METHODS.map((m) => (
                        <SelectItem key={m.value} value={m.value}>
                          {m.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.payment_method && (
                <p className="text-xs text-red-500">{errors.payment_method.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="exp-notes">Catatan</Label>
            <Textarea
              id="exp-notes"
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
          if (!val) handleClose()
        }}
        title={isEdit ? 'Update Pengeluaran' : 'Tambah Pengeluaran'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} pengeluaran "${pendingValues?.description}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
