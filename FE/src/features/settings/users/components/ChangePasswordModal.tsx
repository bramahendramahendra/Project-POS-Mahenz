import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { FormModal } from '@/shared/components'
import { Label } from '@/shared/components/ui/label'
import { PasswordInput } from '@/shared/components/ui/password-input'

import { useChangePasswordMutation } from '../users.api'
import { changePasswordSchema } from '../users.schema'
import type { ChangePasswordFormValues } from '../users.schema'
import type { AppUser } from '../users.types'

interface ChangePasswordModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  user: AppUser | null
}

const defaultValues: ChangePasswordFormValues = { password: '', confirm_password: '' }

export function ChangePasswordModal({ open, onOpenChange, user }: ChangePasswordModalProps) {
  const { mutate: changePassword, isPending } = useChangePasswordMutation()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues,
  })

  useEffect(() => {
    if (open) reset(defaultValues)
  }, [open, reset])

  const onSubmit = (values: ChangePasswordFormValues) => {
    if (!user) return
    changePassword(
      { id: user.id, password: values.password },
      { onSuccess: () => onOpenChange(false) }
    )
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={`Ganti Password${user ? ` — ${user.full_name}` : ''}`}
      size="sm"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel="Simpan"
    >
      <div className="space-y-3">
        <div className="space-y-1.5">
          <Label>
            Password Baru <span className="text-red-500">*</span>
          </Label>
          <PasswordInput
            autoComplete="new-password"
            {...register('password')}
            placeholder="Masukkan password baru"
            className={errors.password ? 'border-red-500' : ''}
          />
          {errors.password && <p className="text-xs text-red-500">{errors.password.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label>
            Konfirmasi Password <span className="text-red-500">*</span>
          </Label>
          <PasswordInput
            autoComplete="new-password"
            {...register('confirm_password')}
            placeholder="Ulangi password baru"
            className={errors.confirm_password ? 'border-red-500' : ''}
          />
          {errors.confirm_password && (
            <p className="text-xs text-red-500">{errors.confirm_password.message}</p>
          )}
        </div>
      </div>
    </FormModal>
  )
}
