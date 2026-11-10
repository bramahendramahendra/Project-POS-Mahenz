import { useEffect } from 'react'
import { useForm, type Resolver } from 'react-hook-form'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { generatePassword } from '../password'
import {
  createUserSchema,
  updateUserSchema,
  type CreateUserFormValues,
} from '../users.schema'
import {
  useCreateUserMutation,
  useRoleOptionsQuery,
  useUpdateUserMutation,
} from '../users.api'
import type { AdminUser } from '../users.types'

const EMPTY: CreateUserFormValues = { username: '', password: '', full_name: '', role_id: 0 }

interface UserFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** Bila diisi: mode ubah. Bila kosong: mode tambah. Gunakan `key` agar form di-reset saat mode berganti. */
  user?: AdminUser | null
}

export function UserFormModal({ open, onOpenChange, user }: UserFormModalProps) {
  const isEdit = !!user
  const { data: roles } = useRoleOptionsQuery()
  const create = useCreateUserMutation()
  const update = useUpdateUserMutation()
  const mutation = isEdit ? update : create

  const form = useForm<CreateUserFormValues>({
    resolver: (isEdit
      ? zodResolver(updateUserSchema)
      : zodResolver(createUserSchema)) as unknown as Resolver<CreateUserFormValues>,
    defaultValues: EMPTY,
  })

  useEffect(() => {
    if (!open) return
    create.reset()
    update.reset()
    form.reset(
      user
        ? { ...EMPTY, full_name: user.full_name, role_id: user.role_id }
        : EMPTY
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user])

  const onSubmit = form.handleSubmit((values) => {
    if (user) {
      update.mutate(
        { id: user.id, payload: { full_name: values.full_name, role_id: values.role_id } },
        { onSuccess: () => onOpenChange(false) }
      )
      return
    }
    create.mutate(values, { onSuccess: () => onOpenChange(false) })
  })

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Ubah User' : 'Tambah User'}
      description={
        isEdit
          ? `Mengubah data @${user?.username}. Username tidak dapat diubah.`
          : 'User baru wajib mengganti password saat login pertama.'
      }
      isLoading={mutation.isPending}
      onSubmit={() => void onSubmit()}
      submitLabel={isEdit ? 'Simpan perubahan' : 'Tambah user'}
    >
      <Form {...form}>
        <form onSubmit={(e) => void onSubmit(e)} className="space-y-4" noValidate>
          {mutation.error && (
            <div
              role="alert"
              className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700"
            >
              <AlertCircle size={16} className="mt-0.5 shrink-0" />
              <span>{mutation.error.message}</span>
            </div>
          )}

          {!isEdit && (
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input {...field} autoComplete="off" placeholder="mis. budi01" />
                  </FormControl>
                  <FormDescription>Huruf dan angka saja, 3-50 karakter.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          <FormField
            control={form.control}
            name="full_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nama lengkap</FormLabel>
                <FormControl>
                  <Input {...field} autoComplete="off" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="role_id"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Role</FormLabel>
                <Select
                  value={field.value ? String(field.value) : ''}
                  onValueChange={(v) => field.onChange(Number(v))}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Pilih role" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {(roles ?? []).map((r) => (
                      <SelectItem key={r.id} value={String(r.id)}>
                        {r.display_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          {!isEdit && (
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password awal</FormLabel>
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
                  <FormDescription>
                    Minimal 8 karakter. Sampaikan ke user; ia akan diminta menggantinya.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          <button type="submit" className="hidden" aria-hidden="true" tabIndex={-1} />
        </form>
      </Form>
    </FormModal>
  )
}
