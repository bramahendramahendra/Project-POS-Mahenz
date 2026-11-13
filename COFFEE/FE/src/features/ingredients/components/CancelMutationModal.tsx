import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'

import { useCancelMutationMutation } from '../ingredients.api'
import { cancelMutationSchema, type CancelMutationFormValues } from '../ingredients.schema'
import type { IngredientMutation } from '../ingredients.types'
import { MUTATION_TYPE_LABEL, formatQty } from '../ingredients.utils'

interface CancelMutationModalProps {
  mutation: IngredientMutation | null
  onOpenChange: (open: boolean) => void
}

// Pembatalan entri terpakai/terbuang oleh admin/owner (keputusan 81). Entri asli tetap terlihat
// bertanda dibatalkan; stok dan total dikembalikan.
export function CancelMutationModal({ mutation, onOpenChange }: CancelMutationModalProps) {
  const { mutate: cancel, isPending } = useCancelMutationMutation()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CancelMutationFormValues>({ resolver: zodResolver(cancelMutationSchema), defaultValues: { reason: '' } })

  useEffect(() => {
    if (mutation) reset({ reason: '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mutation])

  return (
    <FormModal
      open={mutation != null}
      onOpenChange={onOpenChange}
      title="Batalkan Entri"
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit((v) => mutation && cancel({ id: mutation.id, reason: v.reason }, { onSuccess: () => onOpenChange(false) }))}
      submitLabel="Batalkan Entri"
    >
      {mutation && (
        <div className="space-y-3">
          <div className="rounded-md bg-gray-50 px-3 py-2 text-sm text-gray-700">
            <p className="font-medium text-gray-800">{mutation.ingredient_name}</p>
            <p>
              {MUTATION_TYPE_LABEL[mutation.type]} {formatQty(Math.abs(mutation.qty_change), mutation.unit_abbr)} oleh{' '}
              {mutation.created_by_name}
            </p>
            <p className="mt-1 text-xs text-gray-500">Stok dan total akan dikembalikan. Entri ini tetap terlihat di riwayat.</p>
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="cancel-reason">
              Alasan <span className="text-red-500">*</span>
            </Label>
            <Textarea
              id="cancel-reason"
              rows={2}
              {...register('reason')}
              placeholder="Contoh: salah input, seharusnya 1"
              className={errors.reason ? 'border-red-500' : ''}
            />
            {errors.reason && <p className="text-xs text-red-500">{errors.reason.message}</p>}
          </div>
        </div>
      )}
    </FormModal>
  )
}
