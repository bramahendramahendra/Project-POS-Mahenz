import { Navigate } from 'react-router-dom'

import { ROUTES } from '@/shared/constants'

import { useAuthStore } from '../auth.store'

export function RootRedirect() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const mustChangePassword = useAuthStore((s) => s.user?.mustChangePassword ?? false)

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  return <Navigate to={mustChangePassword ? ROUTES.CHANGE_PASSWORD : ROUTES.HOME} replace />
}
