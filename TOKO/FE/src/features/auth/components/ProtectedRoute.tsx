import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'

import { ErrorBoundary } from '@/shared/components'
import { AppLayout } from '@/shared/components/layouts'
import { ROUTES } from '@/shared/constants/routes'
import { useMyMenusQuery } from '@/features/menu/menu.api'
import { useMenuStore } from '@/features/menu/menu.store'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'

import { useAuthStore } from '../auth.store'

interface ProtectedRouteProps {
  menuKey: string
}

export function ProtectedRoute({ menuKey }: ProtectedRouteProps) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const user = useAuthStore((s) => s.user)

  const menus = useMenuStore((s) => s.menus)
  const isLoaded = useMenuStore((s) => s.isLoaded)
  const setMenus = useMenuStore((s) => s.setMenus)
  const hasAccess = useMenuStore((s) => s.hasAccess)

  const { data } = useMyMenusQuery(isAuthenticated && !isLoaded)

  useEffect(() => {
    if (data && !isLoaded) setMenus(data)
  }, [data, isLoaded, setMenus])

  if (!isAuthenticated || !user) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  if (user.apps !== 'web') {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  // Menu (dan permission-nya) belum selesai dimuat — jangan putuskan akses dulu,
  // supaya tidak salah redirect saat refresh browser (menu store tidak di-persist).
  if (!isLoaded) {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-gray-400">
        Memuat...
      </div>
    )
  }

  if (!hasAccess(menuKey)) {
    return <Navigate to={getDefaultRoute(menus)} replace />
  }

  return (
    <AppLayout>
      <ErrorBoundary>
        <Outlet />
      </ErrorBoundary>
    </AppLayout>
  )
}
