import type { ReactNode } from 'react'

import { useAuth } from '@/features/auth/hooks/useAuth'
import { useMenuPermission } from '@/shared/hooks/useMenuPermission'
import type { Role } from '@/shared/types'

interface RoleGuardProps {
  children: ReactNode
  fallback?: ReactNode
  // Mode lama — hardcoded role names (backward compatible)
  allowedRoles?: Role[]
  // Mode baru — dynamic permission dari role_menu_access DB
  menuKey?: string
  action?: 'can_view' | 'can_create' | 'can_edit' | 'can_delete'
}

function PermissionGuard({
  menuKey,
  action = 'can_view',
  children,
  fallback,
}: {
  menuKey: string
  action?: 'can_view' | 'can_create' | 'can_edit' | 'can_delete'
  children: ReactNode
  fallback: ReactNode
}) {
  const perm = useMenuPermission(menuKey)
  const allowed = perm[action]
  return allowed ? <>{children}</> : <>{fallback}</>
}

export function RoleGuard({ allowedRoles, menuKey, action = 'can_view', children, fallback = null }: RoleGuardProps) {
  const { user } = useAuth()

  // Mode baru: gunakan dynamic permission dari store
  if (menuKey) {
    return (
      <PermissionGuard menuKey={menuKey} action={action} fallback={fallback}>
        {children}
      </PermissionGuard>
    )
  }

  // Mode lama: hardcoded role check (backward compatible)
  if (!user || !allowedRoles || !allowedRoles.includes(user.roleName)) {
    return <>{fallback}</>
  }

  return <>{children}</>
}
