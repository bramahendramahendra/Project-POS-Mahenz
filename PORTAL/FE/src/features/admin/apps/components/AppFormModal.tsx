import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertCircle } from 'lucide-react'

import { FormModal } from '@/shared/components'
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

import { useCreateAppMutation, useUpdateAppMutation } from '../apps.api'
import { createAppSchema, type CreateAppFormValues } from '../apps.schema'
import type { AdminApp, CreateAppResult } from '../apps.types'

const ICON_OPTIONS = [
  { value: 'app', label: 'Umum' },
  { value: 'store', label: 'Toko' },
  { value: 'coffee', label: 'Kopi' },
] as const

const EMPTY: CreateAppFormValues = {
  code: '',
  name: '',
  description: '',
  path: '',
  icon: 'app',
  sort_order: 0,
}

interface AppFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** Bila diisi: mode ubah. Bila kosong: mode tambah. */
  app?: AdminApp | null
  onCreated: (result: CreateAppResult) => void
}

export function AppFormModal({ open, onOpenChange, app, onCreated }: AppFormModalProps) {
  const isEdit = !!app
  const create = useCreateAppMutation()
  const update = useUpdateAppMutation()
  const mutation = isEdit ? update : create

  const form = useForm<CreateAppFormValues>({
    resolver: zodResolver(createAppSchema),
    defaultValues: EMPTY,
  })

  useEffect(() => {
    if (!open) return
    create.reset()
    update.reset()
    form.reset(
      app
        ? {
            code: app.code,
            name: app.name,
            description: app.description,
            path: app.path,
            icon: app.icon || 'app',
            sort_order: app.sort_order,
          }
        : EMPTY
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, app])

  const onSubmit = form.handleSubmit((values) => {
    if (app) {
      const { code: _code, ...payload } = values
      void _code
      update.mutate({ id: app.id, payload }, { onSuccess: () => onOpenChange(false) })
      return
    }
    create.mutate(values, {
      onSuccess: (result) => {
        onOpenChange(false)
        onCreated(result)
      },
    })
  })

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Ubah Aplikasi' : 'Tambah Aplikasi'}
      description={
        isEdit
          ? 'Kode aplikasi tidak dapat diubah.'
          : 'Kunci aplikasi dibuat otomatis dan hanya ditampilkan sekali setelah disimpan.'
      }
      isLoading={mutation.isPending}
      onSubmit={() => void onSubmit()}
      submitLabel={isEdit ? 'Simpan perubahan' : 'Buat aplikasi'}
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

          <FormField
            control={form.control}
            name="code"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Kode aplikasi</FormLabel>
                <FormControl>
                  <Input {...field} disabled={isEdit} placeholder="mis. coffee" autoComplete="off" />
                </FormControl>
                <FormDescription>
                  Dipakai sebagai pengenal (app_id) saat login SSO. Tidak dapat diubah.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nama</FormLabel>
                <FormControl>
                  <Input {...field} placeholder="mis. POS Coffee" autoComplete="off" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Deskripsi</FormLabel>
                <FormControl>
                  <Input {...field} placeholder="Keterangan singkat di kartu aplikasi" autoComplete="off" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="path"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Alamat aplikasi</FormLabel>
                <FormControl>
                  <Input {...field} placeholder="/coffee/" autoComplete="off" />
                </FormControl>
                <FormDescription>
                  Path di domain yang sama (mis. /coffee/) atau URL penuh (mis. http://localhost:3002/).
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="grid grid-cols-2 gap-4">
            <FormField
              control={form.control}
              name="icon"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Ikon</FormLabel>
                  <Select value={field.value || 'app'} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {ICON_OPTIONS.map((o) => (
                        <SelectItem key={o.value} value={o.value}>
                          {o.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="sort_order"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Urutan</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      value={field.value}
                      onChange={(e) => field.onChange(e.target.valueAsNumber || 0)}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          {/* Tombol submit tersembunyi agar tombol Enter di input ikut menyimpan */}
          <button type="submit" className="hidden" aria-hidden="true" tabIndex={-1} />
        </form>
      </Form>
    </FormModal>
  )
}
