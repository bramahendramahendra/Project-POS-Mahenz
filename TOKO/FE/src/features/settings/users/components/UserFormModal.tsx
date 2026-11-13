import { useEffect, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useRoleOptionsQuery } from '@/features/settings/roles'

import { useUpdateUserMutation } from '../users.api'
import { updateUserSchema } from '../users.schema'
import type { UpdateUserFormValues } from '../users.schema'
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
  const { data: roles } = useRoleOptionsQuery()
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

interface UserFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  user: AppUser | null
}

// Edit user Toko: hanya role. Username, nama, dan password dikelola di Portal
// (keputusan 69, 76); nama disalin otomatis setiap kali user login.
export function UserFormModal({ open, onOpenChange, user }: UserFormModalProps) {
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<UpdateUserFormValues | null>(null)

  const { mutate: updateUser, isPending } = useUpdateUserMutation()
  const {
    handleSubmit,
    reset,
    setValue,
    control,
    formState: { errors },
  } = useForm<UpdateUserFormValues>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: { role_id: user?.role_id },
  })

  const watchRoleId = useWatch({ control, name: 'role_id' })

  useEffect(() => {
    if (open && user) reset({ role_id: user.role_id })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user])

  if (!user) return null

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
        onSuccess: () => handleClose(),
        onError: () => setIsConfirming(false),
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
            <Label htmlFor="user-username">Username</Label>
            <Input id="user-username" value={user.username} disabled readOnly />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="user-full-name">Nama Lengkap</Label>
            <Input id="user-full-name" value={user.full_name} disabled readOnly />
            <p className="text-xs text-gray-500">Nama dan password diatur di Portal.</p>
          </div>
          <RoleSelect
            value={watchRoleId}
            onChange={(v) => setValue('role_id', v)}
            error={errors.role_id?.message}
          />
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
        title="Update User"
        description={`Yakin ingin mengubah role user "${user.full_name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}
