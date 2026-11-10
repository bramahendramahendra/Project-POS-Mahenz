import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from 'react-router-dom'
import { AlertCircle, Info, KeyRound, Loader2 } from 'lucide-react'

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form'
import { PasswordInput } from '@/shared/components/ui/password-input'
import { ROUTES } from '@/shared/constants'

import { useChangePasswordMutation } from './auth.api'
import { changePasswordSchema, type ChangePasswordFormValues } from './auth.schema'
import { useAuth } from './hooks/useAuth'

export function ChangePasswordPage() {
  const navigate = useNavigate()
  const { mustChangePassword } = useAuth()
  // Disimpan agar teks tidak berubah di tengah proses setelah flag dihapus oleh mutation.
  const [forced] = useState(mustChangePassword)
  const { mutate, isPending, error, reset } = useChangePasswordMutation()

  const form = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: { oldPassword: '', newPassword: '', confirmPassword: '' },
  })

  const onSubmit = (values: ChangePasswordFormValues) => {
    mutate(
      { old_password: values.oldPassword, new_password: values.newPassword },
      { onSuccess: () => navigate(ROUTES.HOME, { replace: true }) }
    )
  }

  return (
    <div className="mx-auto max-w-lg animate-fade-up">
      <div className="mb-6 flex items-center gap-3">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent text-accent-foreground">
          <KeyRound size={20} />
        </span>
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Ganti password</h1>
          <p className="text-sm text-muted-foreground">Gunakan password yang sulit ditebak.</p>
        </div>
      </div>

      {forced && (
        <div
          role="status"
          className="mb-5 flex items-start gap-2 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800"
        >
          <Info size={16} className="mt-0.5 shrink-0" />
          <span>
            Password Anda bersifat sementara. Ganti dulu sebelum menggunakan Portal.
          </span>
        </div>
      )}

      <div className="rounded-3xl border border-border/70 bg-white p-6 shadow-soft sm:p-8">
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            onChange={() => error && reset()}
            className="space-y-5"
            noValidate
          >
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
              name="oldPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password saat ini</FormLabel>
                  <FormControl>
                    <PasswordInput
                      {...field}
                      autoComplete="current-password"
                      className="h-11 rounded-xl"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password baru</FormLabel>
                  <FormControl>
                    <PasswordInput
                      {...field}
                      autoComplete="new-password"
                      className="h-11 rounded-xl"
                    />
                  </FormControl>
                  <FormDescription>Minimal 8 karakter.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Ulangi password baru</FormLabel>
                  <FormControl>
                    <PasswordInput
                      {...field}
                      autoComplete="new-password"
                      className="h-11 rounded-xl"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex items-center justify-end gap-3 pt-2">
              {!forced && (
                <button
                  type="button"
                  onClick={() => navigate(ROUTES.HOME)}
                  className="h-11 rounded-xl px-5 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  Batal
                </button>
              )}
              <button
                type="submit"
                disabled={isPending}
                className="flex h-11 items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-brand to-brand-violet px-6 text-sm font-semibold text-white shadow-glow transition-all hover:brightness-110 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isPending && <Loader2 size={16} className="animate-spin" />}
                Simpan password
              </button>
            </div>
          </form>
        </Form>
      </div>
    </div>
  )
}
