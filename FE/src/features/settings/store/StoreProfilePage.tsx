import { useState } from 'react'
import { Check, Copy, Pencil } from 'lucide-react'

import { PageHeader } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { RoleGuard } from '@/shared/components/RoleGuard/RoleGuard'
import { ROLES } from '@/shared/constants/roles'

import { useStoreProfileQuery } from '../settings.api'
import { StoreProfileForm } from './components/StoreProfileForm'

function CopyField({ label, value }: { label: string; value?: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    if (!value) return
    navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="rounded-lg border bg-white p-4">
      <p className="mb-1 text-xs text-gray-500">{label}</p>
      <div className="flex items-center justify-between gap-2">
        <p className="text-sm font-medium text-gray-900">{value || <span className="text-gray-400 italic">Belum diisi</span>}</p>
        {value && (
          <button
            onClick={handleCopy}
            className="shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
            title={`Salin ${label}`}
          >
            {copied ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
          </button>
        )}
      </div>
    </div>
  )
}

function ReadOnlyField({ label, value }: { label: string; value?: string | number }) {
  return (
    <div className="rounded-lg border bg-white p-4">
      <p className="mb-1 text-xs text-gray-500">{label}</p>
      <p className="text-sm font-medium text-gray-900">
        {value !== undefined && value !== '' ? value : <span className="text-gray-400 italic">Belum diisi</span>}
      </p>
    </div>
  )
}

export function StoreProfilePage() {
  const [isEditing, setIsEditing] = useState(false)
  const { data: profile, isLoading } = useStoreProfileQuery()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Profil Toko"
        breadcrumbs={[{ label: 'Sistem' }, { label: 'Profil Toko' }]}
        actions={
          !isEditing && (
            <RoleGuard allowedRoles={[ROLES.OWNER, ROLES.ADMIN]}>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(true)}>
                <Pencil className="mr-2 h-4 w-4" />
                Edit
              </Button>
            </RoleGuard>
          )
        }
      />

      {isEditing ? (
        <StoreProfileForm
          onCancel={() => setIsEditing(false)}
          onSuccess={() => setIsEditing(false)}
        />
      ) : (
        <div className="max-w-lg space-y-3">
          {isLoading ? (
            Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-lg bg-gray-100" />
            ))
          ) : (
            <>
              <CopyField label="Nama Toko" value={profile?.name} />
              <CopyField label="Alamat" value={profile?.address} />
              <CopyField label="Nomor Telepon" value={profile?.phone} />
              <CopyField label="Email" value={profile?.email} />
              <ReadOnlyField
                label="Pajak Default Kasir (%)"
                value={profile?.tax_default !== undefined ? `${profile.tax_default}%` : undefined}
              />
            </>
          )}
        </div>
      )}
    </div>
  )
}
