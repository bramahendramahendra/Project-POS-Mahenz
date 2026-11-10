import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertCircle, Dices } from 'lucide-react'

import { FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form'
import { Input } from '@/shared/components/ui/input'

import { generatePassword } from '../password'
import { useResetPasswordMutation } from '../users.api'
import { resetPasswordSchema, type ResetPasswordFormValues } from '../users.schema'
import type { AdminUser } from '../users.types'

interface ResetPasswordModalProps {
  user: AdminUser | null
  onClose: () => void
}

export function ResetPasswordModal({ user, onClose }: ResetPasswordModalProps) {
  const { mutate, isPending, error, reset } = useResetPasswordMutation()

  const form = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { password: '' },
  })

  useEffect(() => {
    if (user) {
      reset()
      form.reset({ password: '' })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user])

  const onSubmit = form.handleSubmit((values) => {
    if (!user) return
    mutate({ id: user.id, password: values.password }, { onSuccess: onClose })
  })

  return (
    <FormModal
      open={!!user}
      onOpenChange={(v) => !v && onClose()}
      title="Reset password"
      description={`Tetapkan password sementara untuk @${user?.username ?? ''}. Sesinya akan dihentikan dan ia wajib menggantinya saat login berikutnya.`}
      size="sm"
      isLoading={isPending}
      onSubmit={() => void onSubmit()}
      submitLabel="Tetapkan password"
    >
      <Form {...form}>
        <form onSubmit={(e) => void onSubmit(e)} className="space-y-4" noValidate>
          {error && (
            <div
              role="alert"
              className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700"
            >
              <AlertCircle size={16} className="mt-0.5 shrink-0" />
              <span>{error.message}</span>
            </div>
          )}

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password sementara</FormLabel>
                <div className="flex gap-2">
                  <FormControl>
                    <Input {...field} autoComplete="new-password" className="font-mono" />
                  </FormControl>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => form.setValue('password', generatePassword(), { shouldValidate: true })}
                  >
                    <Dices size={16} /> Acak
                  </Button>
                </div>
                <FormDescription>Minimal 8 karakter. Sampaikan langsung ke user.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <button type="submit" className="hidden" aria-hidden="true" tabIndex={-1} />
        </form>
      </Form>
    </FormModal>
  )
}
