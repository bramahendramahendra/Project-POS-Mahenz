import { Navigate, Outlet, useLocation } from 'react-router-dom'

import { ErrorBoundary } from '@/shared/components'
import { PortalLayout } from '@/shared/components/layouts'
import { ROUTES } from '@/shared/constants/routes'
import type { Role } from '@/shared/types'

import { useAuthStore } from '../auth.store'

interface ProtectedRouteProps {
  /** Bila diisi, hanya role ini yang boleh membuka route di dalamnya. */
  roles?: Role[]
}

export function ProtectedRoute({ roles }: ProtectedRouteProps) {
  const location = useLocation()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const user = useAuthStore((s) => s.user)

  if (!isAuthenticated || !user || user.apps !== 'web') {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  // Password sementara: hanya halaman ganti password yang boleh dibuka.
  if (user.mustChangePassword && location.pathname !== ROUTES.CHANGE_PASSWORD) {
    return <Navigate to={ROUTES.CHANGE_PASSWORD} replace />
  }

  if (roles && !roles.includes(user.roleName)) {
    return <Navigate to={ROUTES.HOME} replace />
  }

  return (
    <PortalLayout>
      <ErrorBoundary>
        <Outlet />
      </ErrorBoundary>
    </PortalLayout>
  )
}
