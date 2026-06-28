import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save } from 'lucide-react'

import { PageHeader } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Checkbox } from '@/shared/components/ui/checkbox'
import { ROUTES } from '@/shared/constants/routes'

import { useRoleDetailQuery, useRoleMenuAccessQuery, useSetRoleAccessMutation } from './roles.api'
import type { RoleMenuAccessItem } from './roles.types'

interface AccessState {
  [menuId: number]: {
    can_view: boolean
    can_create: boolean
    can_edit: boolean
    can_delete: boolean
  }
}

type AccessField = keyof AccessState[number]

interface CheckboxCellProps {
  menuId: number
  field: AccessField
  accessState: AccessState
  onToggle: (menuId: number, field: AccessField) => void
}

function CheckboxCell({ menuId, field, accessState, onToggle }: CheckboxCellProps) {
  const checked = accessState[menuId]?.[field] ?? false
  return (
    <td className="text-center px-3 py-2.5">
      <Checkbox
        checked={checked}
        onCheckedChange={() => onToggle(menuId, field)}
        className="cursor-pointer"
      />
    </td>
  )
}

interface AccessRowProps {
  item: RoleMenuAccessItem
  isChild?: boolean
  accessState: AccessState
  onToggle: (menuId: number, field: AccessField) => void
}

function AccessRow({ item, isChild = false, accessState, onToggle }: AccessRowProps) {
  return (
    <tr key={item.menu_id} className="border-b last:border-0 hover:bg-gray-50">
      <td className="px-4 py-2.5">
        <span className={`${isChild ? 'ml-6 text-gray-600' : 'font-medium'} text-sm`}>
          {isChild && <span className="text-gray-300 mr-2">└</span>}
          {item.label}
          <span className="ml-2 text-xs text-gray-400 font-mono">{item.key_name}</span>
        </span>
      </td>
      <CheckboxCell menuId={item.menu_id} field="can_view" accessState={accessState} onToggle={onToggle} />
      <CheckboxCell menuId={item.menu_id} field="can_create" accessState={accessState} onToggle={onToggle} />
      <CheckboxCell menuId={item.menu_id} field="can_edit" accessState={accessState} onToggle={onToggle} />
      <CheckboxCell menuId={item.menu_id} field="can_delete" accessState={accessState} onToggle={onToggle} />
    </tr>
  )
}

export function RoleAccessPage() {
  const { id } = useParams<{ id: string }>()
  const roleId = Number(id)
  const navigate = useNavigate()

  const { data: role } = useRoleDetailQuery(roleId)
  const { data: accessItems = [], isLoading } = useRoleMenuAccessQuery(roleId)
  const { mutate: saveAccess, isPending: isSaving } = useSetRoleAccessMutation(roleId)

  const [accessState, setAccessState] = useState<AccessState>({})

  useEffect(() => {
    if (accessItems.length > 0) {
      const initial: AccessState = {}
      accessItems.forEach((item) => {
        initial[item.menu_id] = {
          can_view:   item.can_view,
          can_create: item.can_create,
          can_edit:   item.can_edit,
          can_delete: item.can_delete,
        }
      })
      setAccessState(initial)
    }
  }, [accessItems])

  const handleToggle = (menuId: number, field: AccessField) => {
    setAccessState((prev) => {
      const current = prev[menuId] ?? { can_view: false, can_create: false, can_edit: false, can_delete: false }
      const updated = { ...current, [field]: !current[field] }
      if (field === 'can_view' && !updated.can_view) {
        updated.can_create = false
        updated.can_edit   = false
        updated.can_delete = false
      }
      if (field !== 'can_view' && updated[field]) {
        updated.can_view = true
      }
      return { ...prev, [menuId]: updated }
    })
  }

  const handleSave = () => {
    const accesses = Object.entries(accessState).map(([menuId, perm]) => ({
      menu_id:    Number(menuId),
      can_view:   perm.can_view,
      can_create: perm.can_create,
      can_edit:   perm.can_edit,
      can_delete: perm.can_delete,
    }))
    saveAccess({ accesses })
  }

  const parents = accessItems.filter((i) => i.parent_id === null)
  const childrenOf = (parentId: number) => accessItems.filter((i) => i.parent_id === parentId)

  return (
    <div className="space-y-4">
      <PageHeader
        title={`Akses Menu — ${role?.display_name ?? '...'}`}
        breadcrumbs={[
          { label: 'Sistem' },
          { label: 'Manajemen Role', path: ROUTES.SETTINGS_ROLES },
          { label: 'Akses Menu' },
        ]}
        actions={
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => navigate(ROUTES.SETTINGS_ROLES)}>
              <ArrowLeft size={14} className="mr-2" />
              Kembali
            </Button>
            <Button onClick={handleSave} disabled={isSaving}>
              <Save size={14} className="mr-2" />
              {isSaving ? 'Menyimpan...' : 'Simpan Akses'}
            </Button>
          </div>
        }
      />

      <div className="rounded-lg border bg-white overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Menu</th>
              <th className="text-center px-3 py-3 font-medium text-gray-600">Lihat</th>
              <th className="text-center px-3 py-3 font-medium text-gray-600">Tambah</th>
              <th className="text-center px-3 py-3 font-medium text-gray-600">Edit</th>
              <th className="text-center px-3 py-3 font-medium text-gray-600">Hapus</th>
            </tr>
          </thead>
          <tbody>
            {isLoading && (
              <tr><td colSpan={5} className="text-center py-8 text-gray-400">Memuat...</td></tr>
            )}
            {!isLoading && parents.map((parent) => (
              <React.Fragment key={parent.menu_id}>
                <AccessRow item={parent} accessState={accessState} onToggle={handleToggle} />
                {childrenOf(parent.menu_id).map((child) => (
                  <AccessRow key={child.menu_id} item={child} isChild accessState={accessState} onToggle={handleToggle} />
                ))}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
