import { ROLES } from '@/shared/constants'
import type { Role } from '@/shared/types'

import { useAuthStore } from '../auth.store'

export const useAuth = () => {
  const { user, isAuthenticated, accessToken, setSession, clearSession } = useAuthStore()

  return {
    user,
    isAuthenticated,
    accessToken,
    mustChangePassword: user?.mustChangePassword ?? false,
    isAdmin: user?.roleName === ROLES.ADMIN,
    isUser: user?.roleName === ROLES.USER,
    hasRole: (roles: Role[]) => !!user && roles.includes(user.roleName),
    setSession,
    clearSession,
  }
}
