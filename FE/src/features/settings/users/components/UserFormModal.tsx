import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { PasswordInput } from '@/shared/components/ui/password-input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useRoleListQuery } from '@/features/settings/roles/roles.api'

import { useCreateUserMutation, useUpdateUserMutation } from '../users.api'
import { createUserSchema, updateUserSchema } from '../users.schema'
import type { CreateUserFormValues, UpdateUserFormValues } from '../users.schema'
import type { AppUser } from '../users.types'

function RoleSelect({
  value,
  onChange,
  error,
}: {
  value: number | undefined
  onChange: (v: number) => void
  error?: string
}) {
  const { data: roles } = useRoleListQuery({ is_active: true })
  const roleList = roles ?? []

  return (
    <div className="space-y-1.5">
      <Label>
        Role <span className="text-red-500">*</span>
      </Label>
      <Select value={value ? String(value) : ''} onValueChange={(v) => onChange(Number(v))}>
        <SelectTrigger>
          <SelectValue placeholder="Pilih role..." />
        </SelectTrigger>
        <SelectContent>
          {roleList.map((r) => (
            <SelectItem key={r.id} value={String(r.id)}>
              {r.display_name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {error && <p className="text-xs text-red-500">{error}</p>}
    </div>
  )
}

// ─── Create ─────────────────────────────────────────────────────────────────

const createDefaults: CreateUserFormValues = {
  username: '',
  password: '',
  full_name: '',
  role_id: undefined as unknown as number,
}

function CreateUserForm({ open, onOpenChange }: { open: boolean; onOpenChange: (v: boolean) => void }) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<CreateUserFormValues | null>(null)

  const { mutate: createUser, isPending } = useCreateUserMutation()
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateUserFormValues>({
    resolver: zodResolver(createUserSchema),
    defaultValues: createDefaults,
  })

  useEffect(() => {
    if (open) reset(createDefaults)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: CreateUserFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return
    createUser(pendingValues, {
      onSuccess: () => {
        toast.success('User berhasil ditambahkan')
        handleClose()
      },
    })
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title="Tambah User"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Tambah"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label>
              Username <span className="text-red-500">*</span>
            </Label>
            <Input {...register('username')} className={errors.username ? 'border-red-500' : ''} />
            {errors.username && <p className="text-xs text-red-500">{errors.username.message}</p>}
          </div>
          <div className="space-y-1.5">
            <Label>
              Password <span className="text-red-500">*</span>
            </Label>
            <PasswordInput
              autoComplete="new-password"
              {...register('password')}
              className={errors.password ? 'border-red-500' : ''}
            />
            {errors.password && <p className="text-xs text-red-500">{errors.password.message}</p>}
          </div>
          <div className="space-y-1.5">
            <Label>
              Nama Lengkap <span className="text-red-500">*</span>
            </Label>
            <Input {...register('full_name')} className={errors.full_name ? 'border-red-500' : ''} />
            {errors.full_name && <p className="text-xs text-red-500">{errors.full_name.message}</p>}
          </div>
          <RoleSelect
            value={watch('role_id')}
            onChange={(v) => setValue('role_id', v)}
            error={errors.role_id?.message}
          />
        </div>
      </FormModal>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Tambah User"
        description={`Yakin ingin menambahkan user "${pendingValues?.full_name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}

// ─── Edit ───────────────────────────────────────────────────────────────────

function EditUserForm({
  open,
  onOpenChange,
  user,
}: {
  open: boolean
  onOpenChange: (v: boolean) => void
  user: AppUser
}) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<UpdateUserFormValues | null>(null)

  const { mutate: updateUser, isPending } = useUpdateUserMutation()
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<UpdateUserFormValues>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: { full_name: user.full_name, role_id: user.role_id },
  })

  useEffect(() => {
    if (open) reset({ full_name: user.full_name, role_id: user.role_id })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user])

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  const onSubmit = (values: UpdateUserFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return
    updateUser(
      { id: user.id, ...pendingValues },
      {
        onSuccess: () => {
          toast.success('User berhasil diperbarui')
          handleClose()
        },
      }
    )
  }

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title="Edit User"
        size="sm"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label>
              Nama Lengkap <span className="text-red-500">*</span>
            </Label>
            <Input {...register('full_name')} className={errors.full_name ? 'border-red-500' : ''} />
            {errors.full_name && <p className="text-xs text-red-500">{errors.full_name.message}</p>}
          </div>
          <RoleSelect
            value={watch('role_id')}
            onChange={(v) => setValue('role_id', v)}
            error={errors.role_id?.message}
          />
        </div>
      </FormModal>

      <ConfirmDialog
        open={isConfirming}
        onOpenChange={(val) => {
          if (!val) handleClose()
        }}
        title="Update User"
        description={`Yakin ingin mengupdate user "${pendingValues?.full_name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}

// ─── Public wrapper ───────────────────────────────────────────────────────────

interface UserFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  user?: AppUser | null
}

export function UserFormModal({ open, onOpenChange, user }: UserFormModalProps) {
  if (user) {
    return <EditUserForm open={open} onOpenChange={onOpenChange} user={user} />
  }
  return <CreateUserForm open={open} onOpenChange={onOpenChange} />
}
